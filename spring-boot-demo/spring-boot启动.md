
#### 1. spring boot 启动代码  
首先, 看一段 spring boot 的启动代码如下. 发现启动分为2大部分:  
* `SpringApplication.run()` 容器启动
* `@SpringBootApplication` 注解
```java
// (1) 标明是 spring boot 应用, 开启自动配置
@SpringBootApplication   
public class Example {
    public static void main(String[] args) {
        // (2) SpringApplication.run() xxx启动 spring boot 应用
        ConfigurableApplicationContext context = SpringApplication.run(Example.class);
    }
}
```

#### 2. 关于 `SpringApplication.run()` 容器启动
看到使用 `SpringApplication.run(Class<?>[] primarySources)` 这句话启动应用. 内部其实返回了  
```java
return new SpringApplication(primarySources).run(args);
```
这个方法的关键在于两点: 
1. `new SpringApplication(primarySources)` 构造器
	```java
	public SpringApplication(ResourceLoader resourceLoader, Class<?>... primarySources) {
		// resourceLoader 是 spring framework 的内容, 参看 spring framework
		this.resourceLoader = resourceLoader;
		Assert.notNull(primarySources, "PrimarySources must not be null");
		this.primarySources = new LinkedHashSet<>(Arrays.asList(primarySources));   // java config 数组
		//(1) 判断是否是 web 应用. 返回 (REACTIVE, NONE, SERVLET) 类型.  内部就是用 try{Class.forName()}, 看几个 class 是否存在
		this.webApplicationType = WebApplicationType.deduceFromClasspath();
		// (2)
		setInitializers((Collection) getSpringFactoriesInstances(ApplicationContextInitializer.class));
		// (3)
		setListeners((Collection) getSpringFactoriesInstances(ApplicationListener.class));
		// (4)
		this.mainApplicationClass = deduceMainApplicationClass();
	}

	private ResourceLoader resourceLoader;  // spring 资源加载器
	private Set<Class<?>> primarySources;   // java config 数组
	private WebApplicationType webApplicationType;   // web 类型


	private List<ApplicationContextInitializer<?>> initializers;
	private List<ApplicationListener<?>> listeners;
	```
    * (2),(3) 两步用到了 `getSpringFactoriesInstances` 方法来初始化字节的 `initializers` 属性和 `listener` 属性. 该方法如下: 
        ```java
        private <T> Collection<T> getSpringFactoriesInstances(Class<T> type, Class<?>[] parameterTypes, Object... args) {
            ClassLoader classLoader = getClassLoader();
            // 解析 'META_INF/spring.factories' 中的配置, 该文件每行是一个 key-value 对; key: 工厂类名, value: 工厂实现类的类名
            // 将 factoryName 和 fatoryImplementationName 解析出来
            Set<String> names = new LinkedHashSet<>(SpringFactoriesLoader.loadFactoryNames(type, classLoader));
            // 内部用反射 constructor 构造对象
            List<T> instances = createSpringFactoriesInstances(type, parameterTypes, classLoader, args, names);
            AnnotationAwareOrderComparator.sort(instances);
            return instances;
        }
        ```
        因此, (2),(3)两步就是反射构造函数去实例化 `META_INF/spring.factories` 文件中配置的2组类 `ApplicationContextInitializer` 和 `ApplicationListener`. 文件内容为: 
        ```properties
        # Initializers
        org.springframework.context.ApplicationContextInitializer=\
        org.springframework.boot.autoconfigure.SharedMetadataReaderFactoryContextInitializer,\
        org.springframework.boot.autoconfigure.logging.ConditionEvaluationReportLoggingListener
        
        # Application Listeners
        org.springframework.context.ApplicationListener=\
        org.springframework.boot.autoconfigure.BackgroundPreinitializer
        ```
        \[说明\]: `META_INF/spring.factories` 文件是 spring-boot 自动配置的关键
         以 mybatis 为例, 和 spring-boot 集成时, 会引入一个 mybatis-spring-boot-autoconfigure.jar 包. 最重要的是在 `META_INF/spring.factories`中定义了一些列工厂类, 这些类本身作为 bean, 又可以加载其它 bean. 像 mybatis 的自动配置入口就注册了 `sqlSessionFactory` 和 `sqlSessionTemplate`: 
         ```java
        @ConditionalOnClass({ SqlSessionFactory.class, SqlSessionFactoryBean.class })
        class MybatisAutoConfiguration implements InitializingBean{
          @Bean
          @ConditionalOnMissingBean
          public SqlSessionTemplate sqlSessionTemplate(SqlSessionFactory sqlSessionFactory) {
            return new SqlSessionTemplate(sqlSessionFactory);
          }
        }
        ```
    * 第(4)步: `deduceMainApplicationClass()` 是尝试获取执行 main 方法的类. 没什么用, 主要用来打日志
        ```java
        // 通过构造一个 RuntimeException, 查找 main 方法的调用栈来决定启动类是什么
        private Class<?> deduceMainApplicationClass() {
                try {
                    StackTraceElement[] stackTrace = new RuntimeException().getStackTrace();
                    for (StackTraceElement stackTraceElement : stackTrace) {
                        // 调用栈发现 main 方法, 获取主类
                        if ("main".equals(stackTraceElement.getMethodName())) {
                            return Class.forName(stackTraceElement.getClassName());
                        }
                    }
                }
                catch (ClassNotFoundException ex) {
                    // Swallow and continue
                }
                return null;
            }
        ```
2. SpringApplication 的 `run` 方法  
   该方法逻辑较多  
    ```java
    public ConfigurableApplicationContext run(String... args) {
        ////////////////////////////////////
        //StopWatch 主要用来统计 run() 方法的启动时长. 后面有掉用 stopWatch.stop() 方法
        StopWatch stopWatch = new StopWatch();
        stopWatch.start();
        ////////////////////////////////////
        ConfigurableApplicationContext context = null;
        Collection<SpringBootExceptionReporter> exceptionReporters = new ArrayList<>();
        // 设置了 java.awt.headless 环境变量, 和 awt 有关, 可以无视
        configureHeadlessProperty();
        ////////////////////////////////////
        // 反射构造器, 实例化 spring.factories 中的 SpringApplicationRunListener 实现类并启动
        SpringApplicationRunListeners listeners = getRunListeners(args);
        listeners.starting();
        ////////////////////////////////////
        try {
            ApplicationArguments applicationArguments = new DefaultApplicationArguments(args);
            ////////////////////////////////////
            // (1) prepareEnvironment 加载配置变量. 该方法, 根据内部会根据 spring-boot 配置的 profile 加载配置.
            //     生成 PropertySource 和 Environment. 这两个类时干什么的, 参见 spring-framework
            ConfigurableEnvironment environment = prepareEnvironment(listeners, applicationArguments);
            configureIgnoreBeanInfo(environment);
            ////////////////////////////////////
            // (2) 打印 banner
            Banner printedBanner = printBanner(environment);
            // (3) 创建 ApplicaitonContext. 详见下方
            context = createApplicationContext();
            // (4) 对 spring.factories 中记录的配置的 SpringBootExceptionReporter 实现类
            //     org.springframework.boot.diagnostics.FailureAnalyzers 做初始化
            exceptionReporters = getSpringFactoriesInstances(SpringBootExceptionReporter.class,
                    new Class[] { ConfigurableApplicationContext.class }, context);
            // (5) 见下面
            prepareContext(context, environment, listeners, applicationArguments, printedBanner);
            // (6) 内部就是替我们手动执行了 applicationContext.refresh()
            refreshContext(context);
            afterRefresh(context, applicationArguments);
            stopWatch.stop();
            if (this.logStartupInfo) {
                new StartupInfoLogger(this.mainApplicationClass).logStarted(getApplicationLog(), stopWatch);
            }
            // 同下(8)
            listeners.started(context);
            // (7) 执行所有类型为 `ApplicationRunner` 和 `CommandLineRunner` 的 bean 的 run() 方法
            callRunners(context, applicationArguments);
        }
        catch (Throwable ex) {
            handleRunFailure(context, ex, exceptionReporters, listeners);
            throw new IllegalStateException(ex);
        }
    
        try {
            // (8) 有关 SpringApplicationRunListeners 
            listeners.running(context);
        }
        catch (Throwable ex) {
            handleRunFailure(context, ex, exceptionReporters, null);
            throw new IllegalStateException(ex);
        }
        return context;
    }
    ```
    * (1)处, `prepareEnvironment` 加载配置变量
    ```java
    private ConfigurableEnvironment prepareEnvironment(SpringApplicationRunListeners listeners,
                ApplicationArguments applicationArguments) {
            // Create and configure the environment
            ConfigurableEnvironment environment = getOrCreateEnvironment();
            configureEnvironment(environment, applicationArguments.getSourceArgs());
            ConfigurationPropertySources.attach(environment);
            listeners.environmentPrepared(environment);
            bindToSpringApplication(environment);
            if (!this.isCustomEnvironment) {
                environment = new EnvironmentConverter(getClassLoader()).convertEnvironmentIfNecessary(environment,
                        deduceEnvironmentClass());
            }
            ConfigurationPropertySources.attach(environment);
            return environment;
        }
    ```
    * (2)处, 打印的 banner 如下  
    ```
      .   ____          _            __ _ _
     /\\ / ___'_ __ _ _(_)_ __  __ _ \ \ \ \
    ( ( )\___ | '_ | '_| | '_ \/ _` | \ \ \ \
     \\/  ___)| |_)| | | | | || (_| |  ) ) ) )
      '  |____| .__|_| |_|_| |_\__, | / / / /
     =========|_|==============|___/=/_/_/_/
     :: Spring Boot ::
    ```
    * (3) `createApplicationContext` 创建 `ConfigurableApplicationContext`    
     该方法根据判断出的 web 类型(普通web, webserver等)通过反射构造器, 创建 ApplicationContext 对象. 相当于帮你手写了 `new AnnotationConfigApplicationContext()` 初始化容器			
    ```java
    protected ConfigurableApplicationContext createApplicationContext() {
        // 根据 webApplicationType 类型，获得 ApplicationContext 类型
        Class<?> contextClass = this.applicationContextClass;
        if (contextClass == null) {
            try {
                switch (this.webApplicationType) {
                case SERVLET:
                    contextClass = Class.forName(DEFAULT_SERVLET_WEB_CONTEXT_CLASS);
                    break;
                case REACTIVE:
                    contextClass = Class.forName(DEFAULT_REACTIVE_WEB_CONTEXT_CLASS);
                    break;
                default:
                    contextClass = Class.forName(DEFAULT_CONTEXT_CLASS);
                }
            } catch (ClassNotFoundException ex) {
                throw new IllegalStateException("Unable create a default ApplicationContext, " + "please specify an ApplicationContextClass", ex);
            }
        }
        // 创建 ApplicationContext 对象
        return (ConfigurableApplicationContext) BeanUtils.instantiateClass(contextClass);
    }
    ```
    * (5) `prepareContext` 准备 ApplicationContext 对象，主要是初始化它的一些属性。
    ```java
    private void prepareContext(ConfigurableApplicationContext context, ConfigurableEnvironment environment,
                SpringApplicationRunListeners listeners, ApplicationArguments applicationArguments, Banner printedBanner) {
        // (1) 给 applicationContext 设置配的变量
        context.setEnvironment(environment);
        postProcessApplicationContext(context);
        // (2) 缓存 ApplicationContextInitializer
        applyInitializers(context);
        listeners.contextPrepared(context);
        if (this.logStartupInfo) {
            logStartupInfo(context.getParent() == null);
            logStartupProfileInfo(context);
        }
        ConfigurableListableBeanFactory beanFactory = context.getBeanFactory();
        //(3) 在 applicationContexdt 的 BeanFactory 中缓存启动参数. 名称为 "springApplicationArguments"
        //    内部是缓存在 SingletonBeanRegistry 的 LinkedHashSet 数据结构中
        beanFactory.registerSingleton("springApplicationArguments", applicationArguments);
        if (printedBanner != null) {
            beanFactory.registerSingleton("springBootBanner", printedBanner);
        }
        if (beanFactory instanceof DefaultListableBeanFactory) {
            ((DefaultListableBeanFactory) beanFactory)
                    .setAllowBeanDefinitionOverriding(this.allowBeanDefinitionOverriding);
        }
        if (this.lazyInitialization) {
            context.addBeanFactoryPostProcessor(new LazyInitializationBeanFactoryPostProcessor());
        }
        Set<Object> sources = getAllSources();
        Assert.notEmpty(sources, "Sources must not be empty");
    
        // (4) 加载一系列的 BeanDefinition. 具体如下方法
        load(context, sources.toArray(new Object[0]));
        listeners.contextLoaded(context);
    }
    // prepareContext() 的第(4)步骤: 分析怎么加载的 BeanDefinition
    protected void load(ApplicationContext context, Object[] sources) {
        if (logger.isDebugEnabled()) {
            logger.debug("Loading source " + StringUtils.arrayToCommaDelimitedString(sources));
        }
        // (1) 创建 BeanDefinitionRegistry 对象. 具体什么是 BeanDefinitionRegistry 和 BeanDefinition, 参见 spring framework
        BeanDefinitionLoader loader = createBeanDefinitionLoader(getBeanDefinitionRegistry(context), sources);
        if (this.beanNameGenerator != null) {
            loader.setBeanNameGenerator(this.beanNameGenerator);
        }
        if (this.resourceLoader != null) {
            loader.setResourceLoader(this.resourceLoader);
        }
        if (this.environment != null) {
            loader.setEnvironment(this.environment);
        }
        // 执行 BeanDefinition 加载. 具体加载步骤参见 spring-framework
        loader.load();
    }
    ```
    
    * (8)是 SpringApplicationRunListeners 相关处理, 包括 `started()` 和 `run()` 方法. 参看`事件设计模式`
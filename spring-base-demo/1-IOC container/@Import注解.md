#### 一. `@Import`注解的功能  
1. 引入其他的 `@Configuration`    
  下例中, 在 `ConfigB` 中声明 @Bean `ServiceB`, 在 `ConfigA` 中引入 `ConfigB`, 并用 `ConfigA` 配置容器. 发现容器中包含 `ServiceB` bean. 证明 @Import 生效, 起到和 <import> 标签相同的功能
    ```java
    ///////////// ConfigA
    @Import(ConfigB.class)
    @Configuration
    public class ConfigA {}
    
    ///////////// ConfigB
    @Configuration
    public class ConfigB {
        @Bean
        public ServiceInterface getServiceB() {
            return new ServiceB();
        }
    }
    
    ///////////// Service
    public interface ServiceInterface {
        void test();
    }
    
    class ServiceB implements ServiceInterface {
        @Override
        public void test() {
            System.out.println("ServiceB");
        }
    }
    
    
    ///////////// Main 函数
    public class Main {
        public static void main(String[] args) {
            AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext();
            ctx.register(ConfigA.class);
            ctx.refresh();
            // out: ServiceB
            ctx.getBean(ServiceInterface.class).test();
        }
    }
    ``` 
   
2. 实例化其它 bean  
  当 `@Import` 指定的类不是 @Configuration 配置类时, 容器会把类当成 bean 导入.  
  例如把上面代码中的 ConfigA 的 @Import 修改为 `@Import(ServiceB.class)`，就会生成 ServiceB 的 Bean 到容器中
    ```java
    @Import(ServiceB.class)
    @Configuration
    public class ConfigA {}
    ```

3. @Import(ImportSelector.class)     
  导入 `ImportSelector` 中返回的 class name 数组, 可以不把要导入的 Configuration 类或 Bean 类写在 @Import 中, 而是传入一个继承 `ImportSelector` 的子类, 覆盖 `String[] selectImports()` 方法, 返回导入的类名数组
    ```java
    public interface ImportSelector {
        String[] selectImports(AnnotationMetadata importingClassMetadata);   // 返回 包名.类名 即可
    }
    ```
    一般的, 框架中如果要基于 AnnotationMetaData 做动态加载类, 一般会把 `@Import(ImportSelector.class)` 写在自定义的注解上, 让自定义注解提供额外的动态信息. 
    ```java
    //////////////// 自定义注解提供额外的 name 属性
    @Retention(RetentionPolicy.RUNTIME)
    @Documented
    @Target(ElementType.TYPE)
    @Import(ServiceImportSelector.class)
    @interface MyEnableService{
        String name();
    }
    
    //////////////// ImportSelector 实现动态导入
    class ServiceImportSelector implements ImportSelector {
        @Override
        public String[] selectImports(AnnotationMetadata metaData) {
            MultiValueMap<String, Object> attributes = metaData.getAllAnnotationAttributes(MyEnableService.class.getName(), true);
            String nameAttr = attributes.get("name").get(0).toString();
            if ("B".equals(nameAttr)){
                return new String[]{"annotation.testImport.ServiceB"};
            } else {
                return new String[]{"annotation.testImport.ServiceA"};
            }
        }
    }
    
    //////////////// 被 @MyEnableService 注解的 Config 类 
    @MyEnableService(name = "B")
    @Configuration
    public class ConfigA {}
    ```
   此外, @Import 也可以填入 `DeferredImportSelector` 接口的实现类, 同 `ImportSelector` 接口效果一样, 只是前者生成 bean 的时机在后
  
4. @import(ImportBeanDefinitionRegistrar.class) - 构造修改 BeanDefinition        
  使用 `ImportBeanDefinitionRegistrar` 子类, 修改 bean 的`BeanDefinition`. 比如, 前面的 ImportSelector , 实例化 bean 时会调用 bean 的无参构造函数. 如果想调用带参构造函数实例化 bean, 要在 BeanDefinition 中加入构造参数  
  如下例子, ServiceC 的构造函数有 name 属性, 要正确实例化, 要在 Import 中构造 BeanDefinition
    ```java
    //////////////// Bean 只有带参构造函数
    class ServiceC implements ServiceInterface {
        private String name;
    
        public ServiceC(String name) {
            this.name = name;
        }
    
        @Override
        public void test() {
            System.out.println("ServiceC (name:)" + this.name);
        }
    }
    
    //////////////// 自定义注解提供额外的 name 属性
    @Retention(RetentionPolicy.RUNTIME)
    @Documented
    @Target(ElementType.TYPE)
    //@Import(ServiceImportSelector.class)
    @Import(ServiceBeanDefinitionRegistrar.class)
    @interface MyEnableService{
        String name();
    }
    
    //////////////// ImportBeanDefinitionRegistrar 实现动态导入
    class ServiceBeanDefinitionRegistrar implements ImportBeanDefinitionRegistrar {
        @Override
        public void registerBeanDefinitions(AnnotationMetadata metaData, BeanDefinitionRegistry registry) {
            MultiValueMap<String, Object> attributes = metaData.getAllAnnotationAttributes(MyEnableService.class.getName(), true);
            String nameAttr = attributes.get("name").get(0).toString();
            if ("C".equals(nameAttr)){
                BeanDefinitionBuilder builder = BeanDefinitionBuilder.rootBeanDefinition(ServiceC.class).
                        addConstructorArgValue("zhangsan");
                registry.registerBeanDefinition("serviceC", builder.getBeanDefinition());
            } else {
                BeanDefinitionBuilder builder = BeanDefinitionBuilder.rootBeanDefinition(ServiceB.class);
                registry.registerBeanDefinition("serviceB", builder.getBeanDefinition());
            }
        }
    }
    
    //////////////// 被 @MyEnableService 注解的 Config 类
    @MyEnableService(name = "C")
    @Configuration
    public class ConfigA {}
   
    //////////////// Main   
    public class Main {
        public static void main(String[] args) {
            AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext();
            ctx.register(ConfigA.class);
            ctx.refresh();
            // out: ServiceC (name:)zhangsan
            ctx.getBean(ServiceInterface.class).test();
        }
    }   
    ```
  
  
#### 二. `@Import` 注解解析时刻  
 加载解析 `@Import ` 注解位于 `BeanFactoryPostProcessor` 被处理的时候：       
  `AbstractApplicationContext` 的 `refresh()` 方法   
* -> `invokeBeanFactoryPostProcessors(beanFactory);`     
* -> `PostProcessorRegistrationDelegate.invokeBeanFactoryPostProcessors(beanFactory, getBeanFactoryPostProcessors());`     
* -> `registryProcessor.postProcessBeanDefinitionRegistry(registry);`       
  
  这里的registryProcessor,我们指 ConfigurationClassPostProcessor       
* -> `ConfigurationClassPostProcessor.postProcessBeanDefinitionRegistry(registry)`       
* -> `processConfigBeanDefinitions(registry)`
    ```java
    // ConfigurationClassPostProcessor
    
    public void processConfigBeanDefinitions(BeanDefinitionRegistry registry) {
        //省略一些配置检查与设置的逻辑
    
        //根据 @Order 注解，排序所有的 @Configuration 类
        configCandidates.sort((bd1, bd2) -> {
                int i1 = ConfigurationClassUtils.getOrder(bd1.getBeanDefinition());
                int i2 = ConfigurationClassUtils.getOrder(bd2.getBeanDefinition());
                return Integer.compare(i1, i2);});
        // 创建 ConfigurationClassParser 解析 @Configuration 类
        ConfigurationClassParser parser = new ConfigurationClassParser(
                this.metadataReaderFactory, this.problemReporter, this.environment,
                this.resourceLoader, this.componentScanBeanNameGenerator, registry);
    
        //剩余没有解析的 @Configuration 类
        Set<BeanDefinitionHolder> candidates = new LinkedHashSet<>(configCandidates);
        //已经解析的 @Configuration 类
        Set<ConfigurationClass> alreadyParsed = new HashSet<>(configCandidates.size());
        do {
            // 关键解析步骤
            parser.parse(candidates);
               ...  
            //   
            this.reader.loadBeanDefinitions(configClasses);
        }
        while (!candidates.isEmpty());
    
        //省略后续清理逻辑
    }
    ```
*  上面的 parser.parse() 最终调用 `doProcessConfigurationClass()`    
  该方法同时处理 `@Component`, `@ComponentScan`, `@Import` 等重要注解

    ```java
    protected final SourceClass doProcessConfigurationClass(
                ConfigurationClass configClass, SourceClass sourceClass, Predicate<String> filter)
                throws IOException {
    
        //处理`@Component`注解的MemberClass相关代码...
        //处理`@PropertySource`注解相关代码...
        //处理`@ComponentScan`注解相关代码...
        //处理`@Import`注解：
        processImports(configClass, sourceClass, getImports(sourceClass), filter, true);
        //处理`@ImportResource`注解相关代码...
        //处理`@Bean`注解相关代码...
        //处理接口方法相关代码...
        //处理父类相关代码...
    }
    ```
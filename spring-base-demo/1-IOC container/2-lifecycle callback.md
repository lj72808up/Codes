#### 6. bean 的实例化控制 (生命周期回调)
1. bean 的初始化后回调
    * 如果想让 bean 在构建完毕后执行一些自定义的操作, 可以让 bean class 实现接口 `InitializingBean` 的 `afterPropertiesSet()` 方法.
    * 不建议使用实现 `InitializingBean` 接口的方法执行回调, 建议在回调函数上使用 `@PostConstruct` 注解, 或在 xml 中配置 `init-method` 属性
    ```xml
    <bean id="exampleInitBean" class="examples.ExampleBean" init-method="init"/>
    ```

2. bean 的销毁前回调
    * 如果想让 bean 在销毁前执行一些自定义的操作, 可以让 bean class 实现接口 `DisposableBean` 的 `destroy` 方法
    * 不建议使用实现 `DisposableBean` ,接口的方法执行回调, 建议在回调函数上使用 `@PreDestroy` 注解, 或在 xml 中配置 `destroy-method` 属性
    ```xml
    <bean id="exampleInitBean" class="examples.ExampleBean" destroy-method="cleanup"/>
    ```

3. 配置全局默认的创建销毁回调函数
    可以再 xml 的 `beans` 标签中增加 `default-init-method` 和 `default-destroy-method` 标签, 指定全局回调函数. 这样可以不必在每个 bean 上配置 `init-method` 和 `destroy-method`
    ```xml
    <beans default-init-method="init">

        <bean id="blogService" class="com.something.DefaultBlogService">
            <property name="blogDao" ref="blogDao" />
        </bean>

    </beans>
    ```
4. 总结, 自定义 bean 的初始化和销毁回调函数有3种方法
    1. 实现 `InitializingBean` 和 `DisposableBean` 接口
    2. 自定义 `init()`` and `destroy()`` 函数
    3. 使用 ``@PostConstruct` 和 ``@PreDestroy` 注解

5. 优雅关闭 ApplicationContext    
    web 版的 ApplicationContext 已经实现优雅关闭, 非 web 版的 ApplicationContext 可以自定义 hook 实现
    ```java
    ConfigurableApplicationContext context = new ClassPathXmlApplicationContext("application1.xml");
    context.registerShutdownHook();   // AbstractApplicationContext 中实现
    ```
6. 其它声明周期的回调    
   * spring 还提供了很多名字带 `Aware` 后缀的接口, bean class 实现这些接口后, 会在不同生命周期上回调自定义的函数. 
   * 需要注意的是, 这种写法已经破坏了 xml , 或 annotation 配置的 IOC 初衷, 因此, 继承接口定义 bean 的做法只适用于少数基础架构方面的类

7. BeanPostProcessor   
    spring 除了提供了单个 bean 级别的生命周期接口, 还提供了整个容器级别的生命周期接口 `BeanPostProcessor`, 用来对容器内所有 bean 的初始化做回调  
    1. `BeanPostProcessor` 的方法   
        `BeanPostProcessor` 只有2个方法 
        * 初始化前回调: `postProcessBeforeInitialization(Object bean, String beanName)` 
        * 初始化后回调: `postProcessAfterInitialization(Object bean, String beanName)`
    2. `BeanPostProcessor` 与 `InitializingBean`, 各种 `aware` 接口的调用顺序   
        这些顺序要看 beanfactory 类里面是怎么调用的, 以 `AbstractAutowireCapableBeanFactory` 为例
        ```java
        protected Object initializeBean(final String beanName, final Object bean, @Nullable RootBeanDefinition mbd) {
        		if (System.getSecurityManager() != null) {
        			AccessController.doPrivileged((PrivilegedAction<Object>) () -> {
        				invokeAwareMethods(beanName, bean);
        				return null;
        			}, getAccessControlContext());
        		}
        		else {
        		    //若果Bean实现了BeanNameAware、BeanClassLoaderAware、BeanFactoryAware则初始化Bean的属性值
        			invokeAwareMethods(beanName, bean);
        		}
        
        		Object wrappedBean = bean;
        		if (mbd == null || !mbd.isSynthetic()) {
        		    //applyBeanPostProcessorsBeforeInitialization 方法内部主要是
        		    //用来调用后置处理器BeanPostProcessor的postProcessBeforeInitialization方法
        			wrappedBean = applyBeanPostProcessorsBeforeInitialization(wrappedBean, beanName);
        		}
        
        		try {
        		    //Bean如果继承了InitializingBean类则会调用afterPropertiesSet方法，如果设置了init-method方法，
        		    //则调用init-method方法，afterPropertiesSet方法在init-method方法之前调用
        			invokeInitMethods(beanName, wrappedBean, mbd);
        		}
        		catch (Throwable ex) {
        			throw new BeanCreationException(
        					(mbd != null ? mbd.getResourceDescription() : null),
        					beanName, "Invocation of init method failed", ex);
        		}
        		if (mbd == null || !mbd.isSynthetic()) {
        		//applyBeanPostProcessorsAfterInitialization 方法内部主要是，
        		//用来调用后置处理器BeanPostProcessor的postProcessAfterInitialization方法
        			wrappedBean = applyBeanPostProcessorsAfterInitialization(wrappedBean, beanName);
        		}
        
        		return wrappedBean;
        	}
       ```
       因此顺序是:   
       `aware` 接口 -> `BeanPostProcessor.postProcessBeforeInitialization()` -> InitializingBean.afterPropertiesSet -> `BeanPostProcessor.postProcessAfterInitialization()`
    3. `BeanPostProcessor` 与 AOP 的关系   
       aop 因为要返回 bean 的代理对象, 而 BeanPostProcessor 的 postProcessAfterInitialization 方法恰恰是在实例化后调用, 所以可以再这个方法内, 返回一个 bean 的代理对象, 从而实现 AOP
       
    4. AutowiredAnnotationBeanPostProcessor 的实现
        AutowiredAnnotationBeanPostProcessor 是 spring 内部的一个 `BeanPostProcessor`, 用来自动装配 bean 上的注解字段, setter 方法调用, 和任何其他的装配方法
             
8. BeanFactoryPostProcessor
   * `BeanFactoryPostProcessor` 是操作 bean 元数据配置的接口. 该接口会在IOC容器初始化之后，Bean初始化之前被调用. 比如在 `AbstractApplicationContext` 中, 就有调用`BeanFactoryPostProcessor` 接口的方法
        ```java
        // BeanFactoryPostProcessor 的入口在 AbstractApplicationContext
        private void invokeBeanFactoryPostProcessors(
                Collection<? extends BeanFactoryPostProcessor> postProcessors, ConfigurableListableBeanFactory beanFactory) {
    
            for (BeanFactoryPostProcessor postProcessor : postProcessors) {
                postProcessor.postProcessBeanFactory(beanFactory);   // 调用入口
            }
        }
        ```
    * spring IOC 让 `BeanFactoryPostProcessor` 读取用户配置的 metadata , 并对 metadata 进行自定义更改
   一个经典案例是 `PropertySourcesPlaceholderConfigurer`, 它常用于数据库数据源 bean 的配置, 数据源的用户名, 密码往往来自 property 文件, `PropertySourcesPlaceholderConfigurer` 会把 xml 中的属性名 jdbc.username 替换为 property 文件中的 username 属性. 其实现逻辑是在 bean 初始化前, 修改 bean 的 `BeanDefinition`.  springboot 也使用 `PropertySourcesPlaceholderConfigurer` 加载属性配置, 默认从 application.properties
    和 application.yml
       ```xml
        <bean class="org.springframework.context.support.PropertySourcesPlaceholderConfigurer">
            <property name="locations" value="classpath:com/something/jdbc.properties"/>
        </bean>
        
        <bean id="dataSource" destroy-method="close"
                class="org.apache.commons.dbcp.BasicDataSource">
            <property name="driverClassName" value="${jdbc.driverClassName}"/>
            <property name="url" value="${jdbc.url}"/>
            <property name="username" value="${jdbc.username}"/>
            <property name="password" value="${jdbc.password}"/>
        </bean>
        ```
9. BeanFactory    
    BeanFactory 是创建 bean 的工厂, 实现 `BeanFactory` 接口, 可以创造自己的创建 bean 逻辑. spring 中有超过 50 种 `BeanFactory` 
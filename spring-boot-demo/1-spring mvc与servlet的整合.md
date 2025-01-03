前后端分离架构中, spring mvc 三个核心组件
  * HandlerMapping
  * HandlerAdapter
  * HandlerExceptionResolver






#### 1. Servlet 2.0 配置 servlet 的方法
1. 早先的 web.xml 配置
  配置中, 分2步骤: ContextLoaderListener 初始化 WebApplicationContext; 和 DispatcherServlet 初始化的子容器.
  接下来分别看这两个容器是怎么初始化的
```xml
<!-- [1] Spring配置
     ContextLoaderListener 会初始化 Root Spring WebApplicationContext 容器 -->
<listener>
    <listener-class>org.springframework.web.context.ContextLoaderListener</listener-class>
</listener>



<!-- [2] Spring MVC配置
     这是一个 javax.servlet.http.HttpServlet 对象，它除了拦截我们制定的 *.do 请求外，也会初始化一个属于它的 Spring WebApplicationContext 容器。
     该容器是 [1] 容器的子容器 -->
<servlet>
    <servlet-name>spring</servlet-name>  <!-- 这个 name 对应下面 servlet-mapping 中的 name -->
    <servlet-class>org.springframework.web.servlet.DispatcherServlet</servlet-class>
    <load-on-startup>1</load-on-startup>
</servlet>

<servlet-mapping>
    <servlet-name>spring</servlet-name>
    <url-pattern>*.do</url-pattern>
</servlet-mapping>
```


#### 2. ContextLoaderListener 如何初始化 WebApplicationContext
1. 定义
`public class ContextLoaderListener extends ContextLoader implements ServletContextListener`

2. ContextLoaderListener 的构造方法中, 会传入 spring 构造的的 WebApplicationContext, 从而实现和 spring 容器的互通
```java
// ContextLoaderListener

public ContextLoaderListener(WebApplicationContext context) {
    super(context);
}
```

3. ContextLoader 可以自动初始化一个 WebApplicationContext (如果没有传入, 就自动创建)
 内部使用 `BeanUtils.instantiateClass` 生成 WebApplicationContext 对象
```java
// ContextLoader

public WebApplicationContext initWebApplicationContext(ServletContext servletContext) {
    ...
    if (this.context == null) {
        this.context = createWebApplicationContext(servletContext);
    }
    ...
}

protected WebApplicationContext createWebApplicationContext(ServletContext sc) {
		Class<?> contextClass = determineContextClass(sc);
		return (ConfigurableWebApplicationContext) BeanUtils.instantiateClass(contextClass);
	}


protected void configureAndRefreshWebApplicationContext(ConfigurableWebApplicationContext wac, ServletContext sc) {
    ...
    //  调用 refresh() , 联动 ApplicationContext
    wac.refresh()
    ...
}
```

#### 3. DispatcherServlet 初始化第二个 ApplicationContext
1. 结构. DispatcherServlet 从父层级, 依次完成:
    * HttpServletBean ，负责将 ServletConfig 设置到当前 Servlet 对象中。
    * FrameworkServlet ，负责初始化 Spring Servlet WebApplicationContext 容器。
    * DispatcherServlet ，负责初始化 Spring MVC 的各个组件，以及处理客户端的请求

2. HttpServletBean
HttpServletBean 因为 实现了 `EnvironmentAware` 接口, 自动就把 Servlet 相关配置整合到 Bean 中
    ```java
    public abstract class HttpServletBean extends HttpServlet implements EnvironmentCapable, EnvironmentAware
    ```

3. FrameworkServlet
 创建自身 Servlet 的 WebApplicationContext, 并将第一步的 Root WebApplicationContext 设为 parent
 父子容器的意义是: 如果两个容器有父子关系, 则查找 bean 先从子容器找, 找不到再去父容器找. 但父容器不能查找子容器的 bean
                这主要是起到一个 bean 隔离的作用. spring mvc 相关的 bean 由自己的容器管理, 但又要访问一些共有 bean, 因此设置了父 ApplicationContext
    ```java
    // FrameworkServlet.java

    protected WebApplicationContext initWebApplicationContext() {
        ...
        cwac.setParent(rootContext)
        ...
    }
    ```

#### 4. servlet 3.0 如何取消了 web.xml
1. servlet 3.0 取消 web.xml 的两种方法
    * (1) 方法1: 引入了2个注解 ``@WebServlet` 和 `@WebFilter`, 替代 xml 的 `<servlet>`, `<filter>` 标签
    * (2) 方法2:
      servlet3.0 增加新借口 `ServletContext` , 支持:
        * 在运行时动态增加 servlet : ` public ServletRegistration.Dynamic addServlet(String servletName, String className);`
        * filter: `public void addFilter(String className);`
        * listener: `public void addListener(String className);`
      servlet3.0 会自动调用 `ServletContainerInitializer` 接口的 `startup()` 方法在 java code 中增加 servlet, 所以 `startup()` 方法里要调用 ServletContext 的 `addListener`, `addFilter` 方法
      servlet3.0 如何知道是哪个类实现了 `ServletContainerInitializer` 接口呢? 通过 SPI 机制, 将实现类写在 `META_INF/services/接口名` 文件中

2. spring 整合 servlet 的方法是第二种
    spring-web 在 `META_INF/services/javax.servlet.ServletContainerInitializer` 文件中配置了该接口的实现类 `SpringServletContainerInitializer`
    ```java
    @HandlesTypes(WebApplicationInitializer.class)    // spring 将 servlet 的添加委派给了 WebApplicationInitializer
    public class SpringServletContainerInitializer implements ServletContainerInitializer
    ```

3. 以上描述, 都是 servlet3 和 spring 取消 web.xml 的做法, 而 spring-boot 有自己的做法, 并没有遵循 servlet3 的规范

#### 5. spring boot如何取消了 web.xml 加载 servlet  
1. spring boot 中如何注册 servlet 和 filter?  
    1. 方式一: Servlet `3.0 注解` + `@ServletComponentScan`    
       Spring Boot 依旧兼容 Servlet 3.0 一系列以 `@Web*` 开头的注解：`@WebServlet`，`@WebFilter`，`@WebListener`
        ```java
        @WebServlet("/hello")
        public class HelloWorldServlet extends HttpServlet{}
        @WebFilter("/hello/*")
        public class HelloWorldFilter implements Filter {}
        
        @SpringBootApplication
        @ServletComponentScan
        public class SpringBootServletApplication {
           public static void main(String[] args) {
              SpringApplication.run(SpringBootServletApplication.class, args);
           }
        }
        ```
        大体原理，@ServletComponentScan 注解上的 @Import(ServletComponentScanRegistrar.class) ，它会将扫描到的 @WebServlet、@WebFilter、@WebListener 的注解对应的类，最终封装成 FilterRegistrationBean、ServletRegistrationBean、ServletListenerRegistrationBean 对象，注册到 Spring 容器中。也就是说，和注册方式二：RegistrationBean 统一了。
   2. 方式二: 注册 `ServletRegistrationBean` 或 `FilterRegistrationBean`
        ```java
        @Bean
        public ServletRegistrationBean helloWorldServlet() {
            ServletRegistrationBean helloWorldServlet = new ServletRegistrationBean();
            myServlet.addUrlMappings("/hello");
            myServlet.setServlet(new HelloWorldServlet());
            return helloWorldServlet;
        }
        
        @Bean
        public FilterRegistrationBean helloWorldFilter() {
            FilterRegistrationBean helloWorldFilter = new FilterRegistrationBean();
            myFilter.addUrlPatterns("/hello/*");
            myFilter.setFilter(new HelloWorldFilter());
            return helloWorldFilter;
        }
        ```
      ServletRegistrationBean 和 FilterRegistrationBean 都集成自 `RegistrationBean`, `RegistrationBean`实现了 `ServletContextInitializer` 接口 
   
2. 为什么 spring-boot 没有使用 servlet3.0 的标注    
    因为在使用 Spring Boot 时，开发阶段一般都是使用内嵌 Tomcat 容器，但部署时却存在两种选择：一种是打成 jar 包，使用 java -jar 的方式运行；另一种是打成 war 包，交给外置容器去运行。
    前者就会导致容器搜索算法出现问题，因为这是 jar 包的运行策略，不会按照 Servlet 3.0 的策略去加载 ServletContainerInitializer！
2. spring boot 如何实现 servlet 的注册?   
    当使用内嵌 tomcat 时, `javax.servlet.ServletContainerInitializer` 的实现类是 `TomcatStarter`, 但是并没有发现有 SPI 文件对应到 TomcatStarter. 
    1. TomcatStarter 既然不是通过 SPI 机制装配的，那是怎么被 Spring 使用的？
     是在 `TomcatServletWebServerFactory`的 `configureContext(Context context, ServletContextInitializer[] initializers)` 方法中 `new TomcatStarter()` 出来的
        ```java
        // TomcatServletWebServerFactory
       
        protected void configureContext(Context context, ServletContextInitializer[] initializers) {
            TomcatStarter starter = new TomcatStarter(initializers);  // 实例化 TomcatStarter
            if (context instanceof TomcatEmbeddedContext) {
            ...
                embeddedContext.setStarter(starter);
            }
            context.addServletContainerInitializer(starter, NO_CLASSES);
            ...
        }
        ```
   2. 那 `TomcatServletWebServerFactory` 又是如何使用的呢?  
   `TomcatWebServerFactoryCustomizerConfiguration` 作为 `@Configuration`, 将 `TomcatWebServerFactoryCustomizer` 作为 bean 加入容器
        ```java
        // EmbeddedWebServerFactoryCustomizerAutoConfiguration.java
        
        @Configuration(proxyBeanMethods = false)
        @ConditionalOnWebApplication
        @EnableConfigurationProperties(ServerProperties.class)
        public class EmbeddedWebServerFactoryCustomizerAutoConfiguration {
        
            /** Tomcat is being used. */
            @Configuration(proxyBeanMethods = false)
            @ConditionalOnClass({ Tomcat.class, UpgradeProtocol.class })
            public static class TomcatWebServerFactoryCustomizerConfiguration {
        
                @Bean
                public TomcatWebServerFactoryCustomizer tomcatWebServerFactoryCustomizer(Environment environment,
                        ServerProperties serverProperties) {
                    return new TomcatWebServerFactoryCustomizer(environment, serverProperties);
                }
        
            }
             ...
        }
        ```
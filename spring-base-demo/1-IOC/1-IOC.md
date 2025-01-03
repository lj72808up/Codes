#### 1. 什么是IoC
* `IoC (Inversion of Control)` , 也可叫做 `dependency injection (DI)`. 是 spring framework 创建对象, 配置对象, 组装对象的一种方式. 可以创建平时" 用构造函数创建的对象" , 也可以创建" 工厂创建的对象"
* IoC 容器在 spring 中由两类接口实现: `org.springframework.beans` (spring-beans) 和 `org.springframework.context` (spring-context).   
  `org.springframework.beans` 接口是最基本的 Ioc 容器, 用来配置对象的依赖属性 (基本类型或对象类型); `org.springframework.context` 接口是 `beans` 接口的子接口, 在 beans 接口基础场增加了几个特性
        * 支持 AOP 
        * 支持国际化
        * 支持 event publication
        * 用于生成对象的配置文件的不同形式 (比如 `WebApplicationContext`)
* 被 Ioc 容器管理的对象, 在 spring 中叫做 `bean`

#### 2. 配置 bean 并启动 spring container
* 要让 spring 创建`bean` 对象, 需要声明2部分:   
    1. POJO 对象
    2. configuration metadata (对象的元数据配置) : XML 文件. java 注解, 或 java 代码形式  
        * xml 格式的元数据配置如下 
            1. id : `bean` 的唯一标识
            2. class : `bean` 所属的类型, 写类全名
            3. 类的成员变量配置: 
                * name : 成员变量名
                * ref : 引用的 `bean` id
            ```xml
            <?xml version="1.0" encoding="UTF-8"?>
            <beans xmlns="http://www.springframework.org/schema/beans"
                xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                xsi:schemaLocation="http://www.springframework.org/schema/beans
                    https://www.springframework.org/schema/beans/spring-beans.xsd">
            
                <!-- services -->
            
                <bean id="petStore" class="org.springframework.samples.jpetstore.services.PetStoreServiceImpl">
                    <property name="accountDao" ref="accountDao"/>
                    <property name="itemDao" ref="itemDao"/>
                    <!-- additional collaborators and configuration for this bean go here -->
                </bean>
            
                <!-- more bean definitions for services go here -->
            
            </beans>
            ```
* 启动 spring container 并使用对象  
  `ApplicationContext` 是可以获取对象的工厂接口, 它的方法 `T getBean(String name, Class<T> requiredType)` 用来获取 bean 对象
    ```java
    // create and configure beans
    ApplicationContext context = new ClassPathXmlApplicationContext("services.xml", "daos.xml");
    
    // retrieve configured instance
    PetStoreService service = context.getBean("petStore", PetStoreService.class);
    
    // use configured instance
    List<String> userList = service.getUsernameList();
    ```
  
#### 3. bean 初始化配置
* 对于 spring container 来说, 配置的 bean 会被表示为 `BeanDefinition` 对象. 这是最基础的初始化方式, 对象没有依赖, 需要空参构造. 包含以下几个元数据
    1. `bean` 实现类的类名
    2. `bean` 的行为配置
    3. 对其它 `bean` 的引用

* `bean` 的初始化方式配置
    1. 配置构造器初始化
        * POJO
        ```java
        public class ExampleBean {
            @Override
            public String toString() {
                return "hello world";
            }
        }
        ```
        * metadata configuration 
        ```xml
          <!-- 空参构造函数生成bean -->
          <bean id="example1" class="xml.ExampleBean">
          </bean>
        ```
        * 使用生成的对象
        ```java
        ApplicationContext context = new ClassPathXmlApplicationContext("application1.xml");
        ExampleBean bean1 = (ExampleBean) context.getBean("example1");
        System.out.println(bean1.toString());
        ```
    2. 配置静态工厂方法初始化
        * POJO
        ```java
        public class FactoryExampleBean {
            private static FactoryExampleBean example = new FactoryExampleBean();
            private FactoryExampleBean() {}
            public static FactoryExampleBean createInstance() {
                return example;
            }
        
            @Override
            public String toString() {
                return "hello static";
            }
        }
        ```
       * metadata
       ```xml
        <!-- 静态工厂生成 bean -->
        <bean id="example2"
              class = "xml.FactoryExampleBean"
              factory-method="createInstance" />
        ```
       * 使用
       ```java
        // 2. 静态工厂初始化对象
        FactoryExampleBean bean2 = (FactoryExampleBean) context.getBean("example2");
        System.out.println(bean2.toString());
        ```
    3. 配置实例工厂方法初始化  
      实例工厂初始化, 是使用另外的工厂对象的某个实例方法产生的对象, 因此不再需要 class 属性, 取而代之的是 factory-bean 属性
        * POJO
        ```java
        public class Example3{
            @Override
            public String toString() {
                return "hello instance class";
            }
        }
        
        // 产生 Example3 的工厂类
        public class InstanceFactory {
            private static InstanceFactory factory = new InstanceFactory();
            private InstanceFactory(){}
            public static InstanceFactory createInstance() {
                return factory;
            }
            public Example3 makeExample() {
                return new Example3();
            }
        }
        ```
       * metadata
       ```xml
        <!-- 实例工厂生成 bean -->
        <bean id="instanceFactory"
              class="xml.instance.InstanceFactory"/>
        <bean id="example3"
              factory-bean="instanceFactory"
              factory-method="makeExample"/>
        ```
       * 使用对象
       ```java
        // 3. 实例工厂初始化对象
        Example3 bean3 = (Example3) context.getBean("example3");
        System.out.println(bean3.toString());
        ```
#### 4. 依赖 (dependencies)
DI (dependencies injection) 是用来组合对象的, 主要有2种形式的注入: `基于构造器注入`(Constructor-based) 和 `基于setter注入`(Setter-based)
1. 基于构造器注入  
 `构造器注入`是通过 spring 容器调用 bean 的构造函数, 并传入一系列构造器参数的形式产生的. 每个构造器参数都是一个依赖(dependency). 调用静态工厂方法和指定构造器参数的形式是等价的.
    * xml形式的构造器注入: `<constructor-arg>`  
        ```java
        package x.y;
        
        public class ThingOne {
           public ThingOne(ThingTwo thingTwo, ThingThree thingThree) {
               // ...
           }
        }
        ```
        ```xml
        <beans>
            <bean id="beanOne" class="x.y.ThingOne">
                <constructor-arg ref="beanTwo"/>
                <constructor-arg ref="beanThree"/>
            </bean>
        
            <bean id="beanTwo" class="x.y.ThingTwo"/>
        
            <bean id="beanThree" class="x.y.ThingThree"/>
        </beans>
        ```
    * xml 对构造器参数的类型的判断  
      如上例所示, 当`constructor-arg`标签使用 `ref` 属性指定容器中的其它 bean 做为构造器依赖时, 容器可以推断出依赖的类型. 但如果要直接注入简单类型的依赖值, 容器则无法推断出具体的类型 (如42, 是字符串还是int). 解决办法有如下3种: 
        1. 需要使用 type 参数声明依赖的类型
            ```java
            package examples;
            // bean
            public class ExampleBean {
            
                private int years;
                private String ultimateAnswer;
            
                public ExampleBean(int years, String ultimateAnswer) {
                    this.years = years;
                    this.ultimateAnswer = ultimateAnswer;
                }
            }
            ```
           ```xml
           <!--bean-->
            <bean id="exampleBean" class="examples.ExampleBean">
                <constructor-arg type="int" value="7500000"/>
                <constructor-arg type="java.lang.String" value="42"/>
            </bean>
            ```
        2. index 属性: 通过参数角标来推断并注入参数类型.   
          不过当 bean 有参数个数相同的构造器时, 会导致注入失败. 演示忽略~
         
        3. name 属性: 通过构造器的参数名指定依赖的类型  
          通过参数名传入构造器的依赖, 需要在定义类时, 使用`@ConstructorProperties`注解声明参数名称列表
            ```java
            package examples;
            public class ExampleBean {
                @ConstructorProperties({"years", "ultimateAnswer"})
                public ExampleBean(int years, String ultimateAnswer) {
                    this.years = years;
                    this.ultimateAnswer = ultimateAnswer;
                }
            }
            ```
            ```xml
           <!--bean-->
            <bean id="exampleBean" class="examples.ExampleBean">
                <constructor-arg name="years" value="7500000"/>
                <constructor-arg name="ultimateAnswer" value="42"/>
            </bean>
            ```
    * 简写形式
    
2. 基于setter的依赖注入  
 容器在调用无参构造, 或无参的静态工厂方法创建 bean 后, 可以使用 setter-based 注入调用对象的 setter 方法. setter-based DI 可以用在已经通过构造器注入完毕的 bean 上.
    * 我该使用 `Constructor-based` 注入还是 `setter-based` 注入?  
     因为可以混合使用 constructor-based 和 setter-based 注入依赖, 所以一个很好的原则是: 
        * 使用构造函数注入强制依赖
        * 使用 setter 方法或其它配置方法( configuration methods) 注入可选依赖.   
        
       Spring 小组提倡使用 constructor 注入, 这可以确保 bean 是不可变对象, 且保证所有的依赖都是非 null 的. 更进一步讲, constructor 注入让 bean 是完全初始化状态的.   
       Setter 注入应主要用于可选依赖 (可以被赋予默认值的依赖). 否则, `非空校验`(not-null) 要在所有使用这个依赖的地方进行. setter 注入的一个好处是让对象在初始化后可以重新配置, 重新注入  
       最后, 使用的 DI 形式主要看类的形式, 有的类没有暴露任何 setter 方法, 那你只能用构造器注入
       
    * `<property>` 标签配置 setter 注入
        ```java
        public class ExampleBean {
        
            private AnotherBean beanOne;
            private YetAnotherBean beanTwo;
            private int i;
        
            public void setBeanOne(AnotherBean beanOne) {
                this.beanOne = beanOne;
            }
        
            public void setBeanTwo(YetAnotherBean beanTwo) {
                this.beanTwo = beanTwo;
            }
        
            public void setIntegerProperty(int i) {
                this.i = i;
            }
        }
        ```
        ```xml
        <bean id="exampleBean" class="examples.ExampleBean">
            <property name="beanOne">
                <ref bean="anotherExampleBean"/>
            </property>
        
            <property name="beanTwo" ref="yetAnotherBean"/>
            <property name="integerProperty" value="1"/>
        </bean>
        
        <bean id="anotherExampleBean" class="examples.AnotherBean"/>
        <bean id="yetAnotherBean" class="examples.YetAnotherBean"/>
        ```
    
3. spring 容器处理依赖 (bean dependency) 的过程
    * 处理流程
        1. 配置元数据:   
         ApplicationContext is created and initialized with configuration metadata that describes all the beans. 元数据的配置形式可以是 XML, Java 代码, 注解
       
        2. 真正创建 bean 时, 注入依赖:    
          For each bean, its dependencies are expressed in the form of properties, constructor arguments, or arguments to the static-factory method (if you use that instead of a normal constructor). These dependencies are provided to the bean, when the bean is actually created.
       
        3. 构造器的参数:   
          Each property or constructor argument is an actual definition of the value to set, or a reference to another bean in the container.  
          Each property or constructor argument that is a value is converted from its specified format to the actual type of that property or constructor argument. By default, Spring can convert a value supplied in string format to all built-in types, such as int, long, String, boolean, and so forth.  
   
    * spring 容器会校验 bean 的配置, 但是 bean 在使用时才真正创建   
     The Spring container validates the configuration of each bean as the container is created. However, the bean properties themselves are not set until the bean is actually created. Beans that are singleton-scoped and set to be pre-instantiated (the default) are created when the container is created. Scopes are defined in Bean Scopes. Otherwise, the bean is created only when it is requested. Creation of a bean potentially causes a graph of beans to be created, as the bean’s dependencies and its dependencies' dependencies (and so on) are created and assigned. Note that resolution mismatches among those dependencies may show up late — that is, on first creation of the affected bean.
   
    * 环形依赖怎么解决? setter-based 注入   
    如果class A 在进行构造器注入时需要 class B的实例, 而 class B在构造器注入时有需要 class A的实例, 就会造成环形依赖. 构成唤醒依赖的原因是只使用了`构造器注入`. 所以解决办法是将其中一个 class 使用 setter 依赖注入

    * 预加载 singleton bean模式, 解决容器配置异常的延迟显示  
        1. 当 bean 被创建时, spring 会尽可能推迟依赖的处理或bean的属性配置. 换种说法是, 一个正确加载元数据配置的 spring 容器, 可能在使用 bean 时, 发现这个对象或其依赖有问题而抛出异常 - 比如, 抛出对象的属性缺失或属性值无效等异常. 这就是为什么 spring 容器采用预先创建每一个 singleton 的 bean, 目的就是避免配置错误的延迟显示. 当然也可以通过覆盖配置, 把默认的 singleton bean 改为懒加载模式, 让 bean 在第一次使用时才创建     
        2. 有2种级别的懒加载: 单个 bean 上, 或整个 spring 容器上
        ```xml
        <!-- 单个 bean 上-->
        <bean id="lazy" class="com.something.ExpensiveToCreateBean" lazy-init="true"/>
        <!-- 容器上 -->
        <beans default-lazy-init="true">
            <!-- no beans will be pre-instantiated... -->
        </beans>
        ```
     
4. Autowiring
    * 控制 bean 是否作为 autowire candidate 的方法

6. 方法注入 (method injection)  
spring 容器会复写被 lookup 注解的方法并返回容器中改名称的 bean. 这种注入方法是通过 cglib 库的动态生成子类复写被 lookup 注解的方法实现的


#### 5. bean scope
一个 bean 的配置, 相当于一个"实例化某个类的对象"的菜单. "菜单"这个概念很重要, 因为按照一个菜单(bean definition)的步骤来做, 可以生成多个对象  
spring 默认提供6种 scope , 其中4种都是只用于 web 开发的
* singleton (默认的): Spring IoC 容器中, 只能实例化出一个对象
* prototype : 可以实例化出任意数量的对象
* request : 一次 HTTP 请求, 只能实例化一个对象 (只用于 web 版本的 Spring ApplicationContext)
* session: 一个完整的 HTTP Session 中, 只能实例化一个对象  (只用于 web 版本的 Spring ApplicationContext)
* application : 在 ServletContext 的整个生命周期中, 只能实例化一个对象 (只用于 web 版本的 Spring ApplicationContext)
* websocket: 一个 WebSocket 的生命周期中, 只能实例化一个对象 (只用于 web 版本的 Spring ApplicationContext)


1. singleton  
  这种 singleton 实例出的单个对象, 会被保存在一个 cache 中, 后续所有其它 bean 的引用都会指向缓存中的这个对象

2. prototype  
	*  如果一个 beanA 引用了一个 prototype 类型的 beanB, 则每次创建 beanA 的对象时, 会新建一个 beanB 的对象. 同理, 执行`getBean()`方法获取一个 prototype 类型的 bean 时, 也会重新创建对象. 
	```xml
	<bean id="beanB" scope="prototype" class="..."/>
	<bean id="beanA" class ="...">
	```
	* 所以 prototype 适用于无状态的 bean , 而 singleton 适用于有状态的bean
	* spring 只负责创建, 不负责销毁 prototype 类型的bean, 需要手动销毁 (或等待执行正常的 gc). 让 spring 实例化 prototype 类型的 bean, 是 new 创建对象的一种替代
 
3. Singleton Beans 依赖 Prototype-bean 的情况
    * 对于DI (dependency injection): singleton Bean 只实例化一次, 所以其依赖的 Prototype-bean 对象也只实例化一次
    * 对于MI (method injection): 以为方法注入依赖, 会被 cglib 重写, 所以这种方法可以让注入不同的 Prototype-bean 对象

4. Request, Session, Application 和 WebSocket 域   
request, session, application, and websocket 域只用于 web 版本的 Spring ApplicationContext 实现 (such as XmlWebApplicationContext).
    1. 初始化 Web 配置   
       想要使用这几个域, 初始化 web 配置, 比如使用 spring 内置的 web 监听器, 他会将 http request 对象绑定到处理请求的线程, 让 request- 和 session-scoped 域下的 bean 在调用链结束后清除
        ```xml
       <!--加在 web 应用的 web.xml文件中-->
        <web-app>
            ...
            <listener>
                <listener-class>
                    org.springframework.web.context.request.RequestContextListener
                </listener-class>
            </listener>
            ...
        </web-app>
        ```
     2. 这4种 scope 的配置
     	* xml
     	```xml
     	<!-- scope="session" 或者 "application" -->
     	<bean id="loginAction" class="com.something.LoginAction" scope="request"/>
     	```
     	* 注解
     	```java
     	@RequestScope  // 或者 @Sessionscope,  @ApplicationScope
     	@Component
     	public class LoginAction {
     	    // ...
     	}
     	```
     
     3. 如何把这4个短生命周期的 bean 注入到长声明周期的 bean 中?  
       比如将 session scope 的 bean 注入到 singleton 的 bean 中. 为了保证每次 session 获取的 bean 不同, 不能直接注入, 需要把 session scope 的 bean 配置成代理对象注入 (该代理对象拥有和真实 bean 相同的接口). 转换代理对象, 通过 spring 的 aop 配置  
       
       ```xml
     	<?xml version="1.0" encoding="UTF-8"?>
     	<beans xmlns="http://www.springframework.org/schema/beans"
     	    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
     	    xmlns:aop="http://www.springframework.org/schema/aop"
     	    xsi:schemaLocation="http://www.springframework.org/schema/beans
     	        https://www.springframework.org/schema/beans/spring-beans.xsd
     	        http://www.springframework.org/schema/aop
     	        https://www.springframework.org/schema/aop/spring-aop.xsd">
     
     	    <!--  HTTP Session-scoped bean 暴露成一个代理对象 -->
     	    <bean id="userPreferences" class="com.something.UserPreferences" scope="session">
     	        <!-- 告诉 spring 用代理包围 bean -->
     	        <aop:scoped-proxy/> 
     	    </bean>
     
     
     	    <bean id="userService" class="com.something.SimpleUserService">
     	        <property name="userPreferences" ref="userPreferences"/>
     	    </bean>
     	</beans>
       ```
       	
      4. `<aop:scoped-proxy/> ` 配置使用 cglib 做代理, 需要额外导入包. 如果想使用 java 的动态代理, 需要配置属性 `proxy-target-class="false"`. java 动态代理需要被代理的 bean 至少实现一个接口
       	```xml
       	<!-- DefaultUserPreferences implements the UserPreferences interface -->
       	<bean id="userPreferences" class="com.stuff.DefaultUserPreferences" scope="session">
       	    <aop:scoped-proxy proxy-target-class="false"/>
       	</bean>
       
       	<bean id="userManager" class="com.stuff.UserManager">
       	    <property name="userPreferences" ref="userPreferences"/>
       	</bean>
       	```
        
5. 自定义 scope
    1. 创建一个 scope  
    所有的 scope 都是 `org.springframework.context.support.Scope` 的实现, 该接口定义了"获取 bean", "超过生命周期销毁 bean" 等方法  (比如 `org.springframework.web.context.request.SessionScope`)
    2. 使用这个 scope
       beanfactory 注册 scope 
       ```java
        Scope threadScope = new SimpleThreadScope();
        beanFactory.registerScope("thread", threadScope);
       ```
       xml 中配置这个 scope
       ```xml
        <bean class="org.springframework.beans.factory.config.CustomScopeConfigurer">
            <property name="scopes">
                <map>
                    <entry key="thread">
                        <bean class="org.springframework.context.support.SimpleThreadScope"/>
                    </entry>
                </map>
            </property>
        </bean>
    
        <bean id="thing2" class="x.y.Thing2" scope="thread">
            <property name="name" value="Rick"/>
            <aop:scoped-proxy/>
        </bean>
    
        <bean id="thing1" class="x.y.Thing1">
            <property name="thing2" ref="thing2"/>
        </bean>
       ```
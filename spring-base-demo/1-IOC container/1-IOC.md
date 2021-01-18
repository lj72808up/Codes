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
          <bean id="example1" class="exampleBeans.ExampleBean">
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
              class = "exampleBeans.FactoryExampleBean"
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
              class="exampleBeans.instance.InstanceFactory"/>
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
DI (dependencies injection) 是用来组合对象的, 有2中形式的注入: `基于构造器注入` 和 `基于setter注入`
* 基于构造器注入
    * 简写形式
* 基于setter注入
    * 简写形式
    
* 两种注入方式的区别
* 环形依赖的解决办法

* Autowiring

* 懒加载
 spring 默认是 singleton 的作用域
#### 5. bean scope



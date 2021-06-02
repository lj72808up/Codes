1. 从`@Profile`看`@Conditional`     
   在 Spring3.1 的版本，为了满足不同环境注册不同的 Bean ，引入了 @Profile 注解。例如：
    ```java
    @Configuration
    public class DataSourceConfiguration {
    
        @Bean
        @Profile("DEV")
        public DataSource devDataSource() {
            // ... 单机 MySQL
        }
    
        @Bean
        @Profile("PROD")
        public DataSource prodDataSource() {
            // ... 集群 MySQL
        }
        
    }
    ```
   * 在`测试环境`下，我们注册单机 MySQL 的 DataSource Bean 。
   * 在`生产环境`下，我们注册集群 MySQL 的 DataSource Bean 。    
   org.springframework.context.annotation.@Profile ，代码如下：
   ```java
   // Profile.java
   
   @Conditional(ProfileCondition.class)
   public @interface Profile {
   
   	/**
   	 * The set of profiles for which the annotated component should be registered.
   	 */
   	String[] value();
   
   }
   ```
   
2. `@Conditional` 注解
    * `@Conditional` 利用 `Condition` 接口声明的规则, 返回是否 match   
      比如 `@Profile` 的定义上, 配合使用了 `@Conditional(ProfileCondition.class)`
    * `Condition` 接口   
       只有一个方法, 返回是否匹配
        ```java
        public interface Condition {
            boolean matches(ConditionContext context, AnnotatedTypeMetadata metadata);
        }
        ```
3. spring 对 `Condition` 接口的扩展 - `ProfileCondition`
    * `@Profile` 注解利用 `ProfileCondition` 返回是否匹配.   
       获取 `@Profile` 的 value 属性，和 environment 是否有匹配的。如果有，则表示匹配。
        ```java
        // Conditional.java
      
        @Conditional(ProfileCondition.class)
        public @interface Profile { ... }
      
        // ProfileCondition.java
        
        class ProfileCondition implements Condition {
            @Override
            public boolean matches(ConditionContext context, AnnotatedTypeMetadata metadata) {
                // 获得 @Profile 注解的属性
                MultiValueMap<String, Object> attrs = metadata.getAllAnnotationAttributes(Profile.class.getName());
                // 如果非空，进行判断
                if (attrs != null) {
                    // 遍历所有 @Profile 的 value 属性
                    for (Object value : attrs.get("value")) {
                        // 判断 environment 有符合的 Profile ，则返回 true ，表示匹配
                        if (context.getEnvironment().acceptsProfiles(Profiles.of((String[]) value))) {
                            return true;
                        }
                    }
                    // 如果没有，则返回 false
                    return false;
                }
                // 如果为空，就表示满足条件
                return true;
            }      
        }
        ```
      
4. spring boot 对 Condition 接口的扩展 - `SpringBootCondition`
    * `SpringBootCondition` 作为基类, 是一个抽象类. 把 match 结果封装成了 `ConditionOutcome`   
        ```java
        public abstract class SpringBootCondition implements Condition {
        
            private final Log logger = LogFactory.getLog(getClass());
        
            @Override
            public final boolean matches(ConditionContext context, AnnotatedTypeMetadata metadata) {
                // <1> 获得注解的是方法名还是类名
                String classOrMethodName = getClassOrMethodName(metadata);
                try {
                    // <2> 条件匹配结果 ConditionOutcome ( getMatchOutcome()是核心方法 ) 
                    ConditionOutcome outcome = getMatchOutcome(context, metadata);
                    logOutcome(classOrMethodName, outcome);
                    recordEvaluation(context, classOrMethodName, outcome);
                    // <3> 返回是否 match
                    return outcome.isMatch();
                }
                catch (NoClassDefFoundError ex) { ... }
                catch (RuntimeException ex) { ... }
            }
        }
        ```
      
5. spring boot 基于基类 `SpringBootCondition` 实现了大量 Condition 子类.   
    * 以 `OnPropertyCondition` 为例  
        ```java
        // ConditionalOnProperty
        
        @Conditional(OnPropertyCondition.class)
        public @interface ConditionalOnProperty { ...}  
        ```
    * 使用方法  
      通过其两个属性 `name` 以及 `havingValue` 来实现的，其中 name 用来从 `application.properties` 中读取某个属性值。   
        如果该值为空，则返回false;    
        如果值不为空，则将该值与havingValue指定的值进行比较，如果一样则返回true;否则返回false。     
        如果返回值为false，则该configuration不生效；为true则生效。     
        ```java
        @Configuration
        //在application.properties配置"mf.assert"，对应的值为true
        @ConditionalOnProperty(prefix="mf",name = "assert", havingValue = "true")
        public class AssertConfig {
            @Autowired
            private HelloServiceProperties helloServiceProperties;
            @Bean
            public HelloService helloService(){
                HelloService helloService = new HelloService();
                helloService.setMsg(helloServiceProperties.getMsg());
                return helloService;
            }
        }
        ```

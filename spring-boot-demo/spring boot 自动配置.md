#### 3. `@SpringBootApplication` 注解开启自动配置  
自动配置和自动装配不同. `自动装配`指 spring bean 的生成, 自动配置是 spring-boot 的主要动能. 首先先看看该注解的定义   
```java
@Inherited
@SpringBootConfiguration
@EnableAutoConfiguration
@ComponentScan(excludeFilters = { @Filter(type = FilterType.CUSTOM, classes = TypeExcludeFilter.class),
		@Filter(type = FilterType.CUSTOM, classes = AutoConfigurationExcludeFilter.class) })
public @interface SpringBootApplication 
```
下面逐步分析 `@SpringBootApplication` 上的每个注解   

1. `@Inherited`  
该注解是 java 自带注解, 其含义是: 当用 `@Inherited` 实现了自定义注解, 如果在 class 上使用自定义注解, 则继承该 class 的子类自动被标记了自定义注解

2. `@SpringBootConfiguration`
   spring boot 的注解. 看到定义时用到了 `@Configuration`, 因此两者含义一致, 表示该类作为 beans.xml 的替代者来声明 bean
	```java
	@Configuration
	public @interface SpringBootConfiguration {}
	```

3. `@ComponentScan`  
  spring 注解, 自动扫描包下`@Componment`、`@Configuration`、`@Service`


4. `@EnableAutoConfiguration` 
  spring boot 实现自动配置的核心注解  
	```java
	@AutoConfigurationPackage
	@Import(AutoConfigurationImportSelector.class)
	public @interface EnableAutoConfiguration {
	```
	* `@AutoConfigurationPackage` 注解 
		```java
		@Import(AutoConfigurationPackages.Registrar.class)
		public @interface AutoConfigurationPackage
		```

	* `@Import` 
		1. 该注解的功能  
			* (1) @Import(OtherConfiguration.class): 导入其它 `@Configuration` 声明的类, 相当于 beans.xml 中 `<import>` 标签的代替
			* (2) @Import(OtherImportSelector.class): 导入 `ImportSelector` 中代码定义的类名   
				```java
				package com.test;
				class ServiceImportSelector implements ImportSelector {
				    @Override
				    public String[] selectImports(AnnotationMetadata importingClassMetadata) {
				        //可以是@Configuration注解修饰的类，也可以是具体的Bean类的全限定名称
				        return new String[]{"com.test.ConfigB"};
				    }
				}

				@Import(ServiceImportSelector.class)
				@Configuration
				class ConfigA {
				    @Bean
				    @ConditionalOnMissingBean
				    public ServiceInterface getServiceA() {
				        return new ServiceA();
				    }
				}
				```
			* (3) @Import(ImportBeanDefinitionRegistrar.class): 用于在导入 bean 时, 重新定义或更改 bean. 例如动态注入属性，改变Bean的类型和Scope等等
				```java
				// 定义ServiceC
				package com.test;
				class ServiceC implements ServiceInterface {

				    private final String name;

				    ServiceC(String name) {
				        this.name = name;
				    }

				    @Override
				    public void test() {
				        System.out.println(name);
				    }
				}

				// 定义ServiceImportBeanDefinitionRegistrar动态注册ServiceC，修改EnableService
				package com.test;

				@Retention(RetentionPolicy.RUNTIME)
				@Documented
				@Target(ElementType.TYPE)
				@Import(ServiceImportBeanDefinitionRegistrar.class)
				@interface EnableService {
				    String name();
				}

				class ServiceImportBeanDefinitionRegistrar implements ImportBeanDefinitionRegistrar {
				    @Override
				    public void registerBeanDefinitions(AnnotationMetadata importingClassMetadata, BeanDefinitionRegistry registry) {
				        Map<String, Object> map = importingClassMetadata.getAnnotationAttributes(EnableService.class.getName(), true);
				        String name = (String) map.get("name");
				        BeanDefinitionBuilder beanDefinitionBuilder = BeanDefinitionBuilder.rootBeanDefinition(ServiceC.class)
				                //增加构造参数
				                .addConstructorArgValue(name);
				        //注册Bean
				        registry.registerBeanDefinition("serviceC", beanDefinitionBuilder.getBeanDefinition());
				    }
				}

				// 使用 @Import 注解 
				package com.test;
				@EnableService(name = "TestServiceC")
				@Configuration
				class ConfigA {}     // 会自动生成

				// main 函数 
				public static void main(String[] args) {
				    ApplicationContext ctx = new AnnotationConfigApplicationContext(ConfigA.class);
				    ServiceInterface bean = ctx.getBean(ServiceInterface.class);
				    bean.test();
				}
				```
			`[注]`: @Import 注解的解析, 参考 [https://zhuanlan.zhihu.com/p/147025312](https://zhuanlan.zhihu.com/p/147025312)

		2. `@Import(AutoConfigurationImportSelector.class)` 是重头戏的开始, 如后文   


#### 4. `@Import(AutoConfigurationImportSelector.class)` 自动配置重头戏
`AutoConfigurationImportSelector` 类定义如下. 可以发现它除了实现几个Aware类接口外，最关键的就是实现了 `DeferredImportSelector` (继承自ImportSelector)接口  [https://blog.csdn.net/dm_vincent/article/details/77619752](https://blog.csdn.net/dm_vincent/article/details/77619752)
```java
class AutoConfigurationImportSelector implements DeferredImportSelector, BeanClassLoaderAware,
		ResourceLoaderAware, BeanFactoryAware, EnvironmentAware, Ordered{

}
```
1. 该类中的 `getCandidateConfigurations` 获取配置类数组
	1. 函数展示
		```java
		// AutoConfigurationImportSelector.java

		protected List<String> getCandidateConfigurations(AnnotationMetadata metadata, AnnotationAttributes attributes) {
		    // <1> 加载指定类型 EnableAutoConfiguration 对应的，在 `META-INF/spring.factories` 里的类名的数组
		    // getSpringFactoriesLoaderFactoryClass(): return EnableAutoConfiguration.class;
		    // SpringFactoriesLoader.loadFactoryNames():  在 `META-INF/spring.factories` 里的类名的数组
			List<String> configurations = SpringFactoriesLoader.loadFactoryNames(getSpringFactoriesLoaderFactoryClass(), getBeanClassLoader());
			return configurations;
		}
		```

	2. `getCandidateConfigurations` 方法如何被调用的  	
		* (1)处: 在上述分析 `.run()` 方法时, 第`(6)`步, 替我们手动执行了 applicationContext.refresh()
		* (3)处: 
		```java
		public Iterable<Group.Entry> getImports() {
			for (DeferredImportSelectorHolder deferredImport : this.deferredImports) {
				// <1> 处理被 @Import 注解的注解
				// this.group 由 getImportGroup() 方法返回, 具体在"该类的 `getImportGroup()`"章节分析
				this.group.process(deferredImport.getConfigurationClass().getMetadata(),
						deferredImport.getImportSelector());
			}
			// <2> 选择需要导入的
			return this.group.selectImports();
		}
		```

2. 该类的 `getImportGroup()`
	1. 函数展示
	```java
	// AutoConfigurationImportSelector.java

	@Override    // 实现自 DeferredImportSelector 接口
	public Class<? extends Group> getImportGroup() {
		return AutoConfigurationGroup.class;
	}
	```

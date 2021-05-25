Spring framework 核心概念
###一. 统一资源加载策略:  
1. Spring 将资源的定义和资源的加载区分出来  
	* 资源描述接口: `Resource`
	* 资源加载规则接口: `ResourceLoader`   
	  用来根据定义的资源文件地址, 返回 `Resource` 对象. 该接口只有 2个方法  
	 ```java
	 Resource getResource(String location);  // 核心方法, 返回Resource
	 ClassLoader getClassLoader();
	 ```

2. 资源和加载接口提供了默认实现  
	* `org.springframework.core.io.AbstractResource` 为 `Resource` 接口的默认抽象实现  
	   它实现了 `Resource` 接口的大部分的公共实现. 实现的默认方法包括: 判断文件是否存在, 最后修改时间等
	* `DefaultResourceLoader`类  为`ResourceLoader` 的默认实现, 每次只能返回单一的资源 
	* `PathMatchingResourcePatternResolver` 是一个集大成者的 `ResourceLoader`. 它还实现了额外的接口 `ResourcePatternResolver`.   
	   既有 `Resource getResource(String location)` 方法，也有 `Resource[] getResources(String locationPattern)` 方法。



### 二. 加载 BeanDefinition   
先看一段 spring 使用的标准代码
```java
// <1> 获取资源
ClassPathResource resource = new ClassPathResource("application1.xml");
// <2> 获取 BeanFactory
DefaultListableBeanFactory factory = new DefaultListableBeanFactory();
// <3> 组合 BeanFactory 创建 BeanDefinitionReader, 该 Reader 为 Resource 的解析器
XmlBeanDefinitionReader reader = new XmlBeanDefinitionReader(factory);
// <4> 装载 Resource
reader.loadBeanDefinitions(resource); 
```
注释中初步给出 spring 容器的使用过程. 该过程可分三步理解: 
* 资源定位: 比如如果用 beans.xml 配置 bean, 则首先要拿到 beans.xml 这个资源文件. 这一步是返回 `Resource` 的过程
* 装载: 
  使用 `BeanDefinitionReader` 解析 `Resource` 文件, 将文件内容转变成 `BeanDefinition` 结构
	* 在 IoC 容器内部维护着一个 `BeanDefinition Map` 的数据结构
	* 在配置文件中每一个 `<bean>` 都对应着一个 `BeanDefinition` 对象。
* 注册  
  将解析出的 BeanDefinition 对象通过 `BeanDefinitionRegistry` 接口来实现注册, 将这些解析的 BeanDefinition 注入到一个 `HashMap` 容器中  
  此时并没有创建 bean. 只是注册了 BeanDefinition

1. 装载 BeanDefinition 在上面第<4>步   
    ```java
    public int loadBeanDefinitions(EncodedResource encodedResource) throws BeanDefinitionStoreException {
        Assert.notNull(encodedResource, "EncodedResource must not be null");
        if (logger.isTraceEnabled()) {
            logger.trace("Loading XML bean definitions from " + encodedResource);
        }
        // this.resourcesCurrentlyBeingLoaded 是一个 ThreadLocal<Set<EncodedResource>>, 记录当前线程注册过的 EncodedResource
        Set<EncodedResource> currentResources = this.resourcesCurrentlyBeingLoaded.get();
    
        if (!currentResources.add(encodedResource)) {
            throw new BeanDefinitionStoreException(
                    "Detected cyclic loading of " + encodedResource + " - check your import definitions!");
        }
    
        try (InputStream inputStream = encodedResource.getResource().getInputStream()) {
            InputSource inputSource = new InputSource(inputStream);
            if (encodedResource.getEncoding() != null) {
                inputSource.setEncoding(encodedResource.getEncoding());
            }
            // <核心流程>: 解析 Resource ,生成 BeanDefinition
            return doLoadBeanDefinitions(inputSource, encodedResource.getResource());
        }
        catch (IOException ex) {
            throw new BeanDefinitionStoreException(
                    "IOException parsing XML document from " + encodedResource.getResource(), ex);
        }
        finally {
            currentResources.remove(encodedResource);
            if (currentResources.isEmpty()) {
                // 当 currentResources 为空时, 证明该 ThreadLocal 已经不会被使用. 这会导致 ThreadLocalMap.Entry 
                // 键值对内部的 ThreadLocal key 软引用在 gc 时被置空. 而 key 为 null 的 value 访问不到, 却又 gc 不了, 导致溢出
                // 要及时把不再使用的 ThreadLocal 调用 remove 删除
                this.resourcesCurrentlyBeingLoaded.remove();
            }
        }
    }
    ```
2. 核心流程 `doLoadBeanDefinitions()`  
    ```java
    // XmlBeanDefinitionReader.java
    protected int doLoadBeanDefinitions(InputSource inputSource, Resource resource)
                throws BeanDefinitionStoreException {
        try {
            // <1> 解析 xml Resource, 返回 xml 的 Document 对象
            Document doc = doLoadDocument(inputSource, resource);
            // <2> 从 Document 对象中解析出 BeanDefinition 并注册
            int count = registerBeanDefinitions(doc, resource);
            if (logger.isDebugEnabled()) {
                logger.debug("Loaded " + count + " bean definitions from " + resource);
            }
            return count;
        }
    }
    ```
	* <2> 步骤 `registerBeanDefinitions(doc, resource)` 返回注册的 BeanDefinition 个数  
	```java
	// AbstractBeanDefinitionReader.java
	private final BeanDefinitionRegistry registry;

	// XmlBeanDefinitionReader.java
	public int registerBeanDefinitions(Document doc, Resource resource) throws BeanDefinitionStoreException {
		// <1> 创建 BeanDefinitionDocumentReader 对象
		BeanDefinitionDocumentReader documentReader = createBeanDefinitionDocumentReader();
		// <2> 获取已注册的 BeanDefinition 数量
		int countBefore = getRegistry().getBeanDefinitionCount();
		// <3> 创建 XmlReaderContext 对象
		// <4> 注册 BeanDefinition
		documentReader.registerBeanDefinitions(doc, createReaderContext(resource));
		// 计算新注册的 BeanDefinition 数量
		return getRegistry().getBeanDefinitionCount() - countBefore;
	}
	```
    * `BeanDefinition` 是去用来描述一个 bean.   
          它的属性对应 xml 中配置的属性. 比如: 生成 bean 的工厂类, 生成 bean 的工厂方法, bean 的构造函数, bean 的 scope 等  
      ```xml
      <!-- 空参构造函数生成bean -->
      <bean id="example1" class="xml.ExampleBean" init-method="init"/>
      
      <!-- 静态工厂生成 bean -->
      <bean id="example2"
            class="xml.FactoryExampleBean"
            factory-method="createInstance"/>
      
      <!-- 实例工厂生成 bean -->
      <bean id="instanceFactory"
            class="xml.instance.InstanceFactory"/>
      <bean id="example3"
            factory-bean="instanceFactory"
            factory-method="makeExample"
      />
      ```

3. 什么是 `BeanDefinition` 的注册   
  就是将`<neanName, BeanDefinition>`和`<beanName,aliasName>`两种映射关系注册到 `BeanFactory` 中
    * 步骤<1>, 注册 beanName 和 BeanDefinition 的方法是定义在 BeanFactory 内的.    
      `BeanFactory` 内部持有 beanName 和 BeanDefinition 映射的 ConcurrentHashMap. 添加映射时, 还要用同步块处理并行的情况
    ```java
    // DefaultListableBeanFactory.java

    // 是否允许同名 beanName 的 BeanDefinition 进行覆盖注册: true
    private boolean allowBeanDefinitionOverriding = true;
    /** 存储 beanName 和 BeanDefinition 映射的 HashMap */
    private final Map<String, BeanDefinition> beanDefinitionMap = new ConcurrentHashMap<>(256);
    /** List of bean definition names, in registration order. */
    private volatile List<String> beanDefinitionNames = new ArrayList<>(256);
    /** List of names of manually registered singletons, in registration order. */
    private volatile Set<String> manualSingletonNames = new LinkedHashSet<>(16);

    @Override
    public void registerBeanDefinition(String beanName, BeanDefinition beanDefinition)
            throws BeanDefinitionStoreException {
        ...  // 省略各种校验和判断是否可以覆盖同 beanName 和 BeanDefinition 的设置
        // beanDefinitionMap 为全局变量，避免并发情况
        // Cannot modify startup-time collection elements anymore (for stable iteration)
        synchronized (this.beanDefinitionMap) {
            // 添加到 BeanDefinition 到 beanDefinitionMap 中。
            this.beanDefinitionMap.put(beanName, beanDefinition);
            // 添加 beanName 到 beanDefinitionNames 中
            List<String> updatedDefinitions = new ArrayList<>(this.beanDefinitionNames.size() + 1);
            updatedDefinitions.addAll(this.beanDefinitionNames);
            updatedDefinitions.add(beanName);
            this.beanDefinitionNames = updatedDefinitions;
            // 从 manualSingletonNames 移除 beanName
            if (this.manualSingletonNames.contains(beanName)) {
                Set<String> updatedSingletons = new LinkedHashSet<>(this.manualSingletonNames);
                updatedSingletons.remove(beanName);
                this.manualSingletonNames = updatedSingletons;
            }
        }       
    }
    ```
    * 步骤 <2> 中的注册 beanName 和 regist name 映射, 也是注册到一个 ConcurrentHashMap 中  
    ```java
    // SimpleAliasRegistry,java

    /** Map from alias to canonical name. */
    private final Map<String, String> aliasMap = new ConcurrentHashMap<>(16);
    @Override
    public void registerAlias(String name, String alias) {
        ... 
        synchronized (this.aliasMap) {
            // (1) 循环检验: 看是否有和 nam->alias 相反的 alias->name 映射
            checkForAliasCircle(name, alias);
            // (2) 加入HashMap
            this.aliasMap.put(alias, name);
        }
    }
    ```







































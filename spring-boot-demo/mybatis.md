#### spring boot 自动配置 Mybatis
mybatis-spring-boot-autoconfigure 中的 `spring.factories` 中， 配置了 `org.mybatis.spring.boot.autoconfigure.MybatisAutoConfiguration` 这个 @Configuration 类。
```java
@org.springframework.context.annotation.Configuration
public class MybatisAutoConfiguration implements InitializingBean {
	@Bean   // SqlSessionFactoryBean 是 SqlSessionFactory 的一个实现
	public SqlSessionFactory sqlSessionFactory(DataSource dataSource) throws Exception {
		SqlSessionFactoryBean factory = new SqlSessionFactoryBean();
		... 设置一堆属性 ... 
		return factory.getObject();
	}
	@Bean   // sqlSessionTemplate 是 SqlSession 的一个实现
	public SqlSessionTemplate sqlSessionTemplate(SqlSessionFactory sqlSessionFactory) {
		return new SqlSessionTemplate(sqlSessionFactory);
	}

}
```

SqlSessionFactory 是怎么创建 SqlSession 的？ （怎么创建SqlSessionTemplate的？）

#### 1. 先来看看 SqlSessionFactoryBean 是什么
```java
public class SqlSessionFactoryBean
implements FactoryBean<SqlSessionFactory>, InitializingBean, ApplicationListener<ApplicationEvent> {
    public SqlSessionFactory getObject() throws Exception {
            if (this.sqlSessionFactory == null) {
                // 内部会返回 new DefaultSqlSessionFactory
                afterPropertiesSet();
            }
            return this.sqlSessionFactory;
        }
    }
}
```
SqlSessionFactoryBean 实现了 `FactoryBean` 接口， 代表 SqlSessionFactoryBean 的实例不只是一个普通的bean对象，还可以产生 一个 bean 放到容器中。 该类的 `getBean` 会返回 `DefaultSqlSessionFactory`

#### 2. `SqlSessionTemplate` 如何产生
  这问题是在SqlSessionTemplate中实现的。 看到上面的 `new SqlSessionTemplate(sqlSessionFactory)` 返回了 SqlSessionTemplate

#### 3. 如何执行对一个的方法   
SqlSessionTemplate 的构造函数会生成内部的 SqlSession 代理对象。 SqlSessionTemplate 内部的所有 `select` 方法均由该代理对象执行
```java
public class SqlSessionTemplate implements SqlSession, DisposableBean {
// 代理对象
private final SqlSession sqlSessionProxy;

public SqlSessionTemplate(SqlSessionFactory sqlSessionFactory, ExecutorType executorType,
PersistenceExceptionTranslator exceptionTranslator) {
    ... ...
    // JDK 代理
    this.sqlSessionProxy = (SqlSession) newProxyInstance(
        SqlSessionFactory.class.getClassLoader(),
        new Class[] { SqlSession.class },
        new SqlSessionInterceptor()    // InvocationHandler 对象
        );
    }
}
```
 InvocationHandler 使用 DefaultSqlSession 执行真正的方法
```java
private class SqlSessionInterceptor implements InvocationHandler {
    @Override
    public Object invoke(Object proxy, Method method, Object[] args) throws Throwable {
        // (1) 这里的 SqlSessionFactory 就是上文由 DefaultSqlSessionFactory 创建的 DefaultSqlSession
        //     getSqlSession() 这个方法, 内部会先从 ThreadLocal 中找是否存在 SqlSession, 
        //     没有则开启事务管理器, 该事务管理器会负责后续jdbc连接管理工作；
        SqlSession sqlSession = getSqlSession(SqlSessionTemplate.this.sqlSessionFactory,
        SqlSessionTemplate.this.executorType, SqlSessionTemplate.this.exceptionTranslator);
        try {
            // (2) 由 DefaultSqlSession 执行 select 等方法
            Object result = method.invoke(sqlSession, args);
            // (3) 是否需要自动提交事务
            if (!isSqlSessionTransactional(sqlSession, SqlSessionTemplate.this.sqlSessionFactory)) {
            sqlSession.commit(true);
            }
            return result;
        } catch (Throwable t) {
        ... ...
        } finally {
            if (sqlSession != null) {
                closeSqlSession(sqlSession, SqlSessionTemplate.this.sqlSessionFactory);
            }
        }
    }
}
```

#### 4. 每次执行 sql 语句开启一个 SqlSession 吗? SqlSession 是线程安全的吗? 
* 上面的代码分析中, SqlSession 会从 ThreadLocal 中先获取, 因此, SqlSession 是线程安全的. 不同线程的 sql 执行才会开启新的 SqlSession;   
* SqlSession 中包含一个Executor, 用来执行 sql 语句; 而 Executor 中包含一个 `Transaction` 用来表示事务. 这个 `Transaction` 是从 `SpringManagedTransaction` 中获取的, `SpringManagedTransaction.commit()` 调用的 `connection.commit()`

#### 5. DefaultSqlSession 的各种 select 方法， 都是使用 `Executor` 接口完成的   
真正的 Sql 执行者： `Executor`.  `Executor` 包装了一级缓存和二级缓存   
* 一级缓存: 共享在一个 sqlSession 中, 如果 sql 语句相同, 会优先命中一级缓存
* 二级缓存: 在多个 SQLSession 中共享, 开启二级缓存会使用 `CachingExecutor` 装饰 Executor



#### 6.SqlSessionFactory 是怎么创建 SqlSession 的？ （怎么创建SqlSessionTemplate的？）

1. 先来看看 SqlSessionFactoryBean 是什么
   ```java
   public class SqlSessionFactoryBean
   implements FactoryBean<SqlSessionFactory>, InitializingBean, ApplicationListener<ApplicationEvent> {

   	public SqlSessionFactory getObject() throws Exception {
    		if (this.sqlSessionFactory == null) {
    			// 内部会返回 new DefaultSqlSessionFactory
      			afterPropertiesSet();
    		}
    		return this.sqlSessionFactory;
   		}
   	}
   }
   ```
   SqlSessionFactoryBean 实现了 `FactoryBean` 接口， 代表 SqlSessionFactoryBean 的实例不只是一个普通的bean对象，还可以产生 一个 bean 放到容器中。 该类的 `getBean` 会返回 `DefaultSqlSessionFactory`

2. `SqlSessionTemplate` 如何产生
  这问题是在SqlSessionTemplate中实现的。 看到上面的 `new SqlSessionTemplate(sqlSessionFactory)` 返回了 SqlSessionTemplate

3. 如何执行对一个的方法   
  SqlSessionTemplate 的构造函数会生成内部的 SqlSession 代理对象。 SqlSessionTemplate 内部的所有 `select` 方法均由该代理对象执行
    ```java
    public class SqlSessionTemplate implements SqlSession, DisposableBean {
    // 代理对象
    private final SqlSession sqlSessionProxy;
    
    public SqlSessionTemplate(SqlSessionFactory sqlSessionFactory, ExecutorType executorType,
    PersistenceExceptionTranslator exceptionTranslator) {
        ... ...
        // JDK 代理
        this.sqlSessionProxy = (SqlSession) newProxyInstance(
            SqlSessionFactory.class.getClassLoader(),
            new Class[] { SqlSession.class },
            new SqlSessionInterceptor()    // InvocationHandler 对象
            );
        }
    }
    ```
     InvocationHandler 使用 DefaultSqlSession 执行真正的方法
    ```java
    private class SqlSessionInterceptor implements InvocationHandler {
        @Override
        public Object invoke(Object proxy, Method method, Object[] args) throws Throwable {
            // (1) 这里的 SqlSessionFactory 就是上文由 DefaultSqlSessionFactory 创建的 DefaultSqlSession
            SqlSession sqlSession = getSqlSession(SqlSessionTemplate.this.sqlSessionFactory,
            SqlSessionTemplate.this.executorType, SqlSessionTemplate.this.exceptionTranslator);
            try {
                // (2) 由 DefaultSqlSession 执行 select 等方法
                Object result = method.invoke(sqlSession, args);
                // (3) 是否需要自动提交事务
                if (!isSqlSessionTransactional(sqlSession, SqlSessionTemplate.this.sqlSessionFactory)) {
                sqlSession.commit(true);
                }
                return result;
            } catch (Throwable t) {
            ... ...
            } finally {
                if (sqlSession != null) {
                    closeSqlSession(sqlSession, SqlSessionTemplate.this.sqlSessionFactory);
                }
            }
        }
    }
    ```

4. DefaultSqlSession 的各种 select 方法， 都是使用 `Executor` 接口完成的   
  真正的 Sql 执行者： `Executor`. 之所以使用 `Executor`, 因为其提供了一级缓存和二级缓存


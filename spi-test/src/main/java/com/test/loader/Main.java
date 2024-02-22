package com.test.loader;

import java.io.IOException;
import java.net.URL;
import java.sql.Driver;
import java.util.Enumeration;
import java.util.Iterator;
import java.util.ServiceLoader;

public class Main {

    public static void test() throws IOException {
//  jar:file:/home/lj/.m2/repository/mysql/mysql-connector-java/8.0.25/mysql-connector-java-8.0.25.jar!/META-INF/services/java.sql.Driver
//  jar:file:/home/lj/.m2/repository/org/postgresql/postgresql/42.2.22/postgresql-42.2.22.jar!/META-INF/services/java.sql.Driver
        String url = "META-INF/services/java.sql.Driver";
        // 用某个 classLoader 对象, 在其查找范围内查找 url 枚举数组
        // Thread 中的都是 AppClassLoader
        Enumeration<URL> enums = Thread.currentThread().getContextClassLoader().getResources(url);
        while (enums.hasMoreElements()) {
            System.out.println(enums.nextElement());
        }
    }

    public static void main(String[] args) throws IOException {
        test();
// new ServiceLoader<>(service, loader);  loader 是 contextClassLoader, Thread.currentThread().getContextClassLoader()
        /**
         * 这个方法的意思: Driver.class 接口的全名 java.sql.Driver, 于是用 AppClassLoader 去所有 jar 包里的 META-INF/services 目录中,
         * 查找名为 java.sql.Driver 的文件, 里面写了对应的实现类类名
         *
         * ## ServiceLoader.java (load()方法)
         * public static <S> ServiceLoader<S> load(Class<S> service) {
         *     ClassLoader cl = Thread.currentThread().getContextClassLoader();  // cl 是 AppClassLoader
         *     return new ServiceLoader<>(Reflection.getCallerClass(), service, cl);
         * }
         */
        ServiceLoader<Driver> driverServices = ServiceLoader.load(Driver.class);
        // ServiceLoader 实例化对象分两步
        // (1) 迭代每个 META-INF/services/java.sql.Driver 文件里的内容
        //     对每一行进行 Class.forName(cn, false, loader) 加载类;   // cn: com.mysql.cj.jdbc.Driver.class, flase: 不实例化,loader: AppClassLoader
        // (2) 调用对应的共有无参构造函数, 实例化对象. 如果没有这个构造函数, 会抛异常
        //  ServiceLoader.java 实例化实现类代码片段:
        //  Constructor<?> ctor = clazz.getConstructor();
        //  ctor.setAccessible(true);
        //  ctor.newInstance(); // 在 ProviderImple.get() 时获取
        Iterator<Driver> driversIterator = driverServices.iterator();

        while (driversIterator.hasNext()) {
            // 这里会遍历文件中的每行
            Driver d = driversIterator.next();
            System.out.println(d.getClass().getName());
        }

    }
}

package com.test.loader;

import java.io.IOException;
import java.net.URL;
import java.sql.Driver;
import java.util.Enumeration;
import java.util.Iterator;
import java.util.ServiceLoader;

public class Main {

    public static void test () throws IOException {
//  jar:file:/home/lj/.m2/repository/mysql/mysql-connector-java/8.0.25/mysql-connector-java-8.0.25.jar!/META-INF/services/java.sql.Driver
//  jar:file:/home/lj/.m2/repository/org/postgresql/postgresql/42.2.22/postgresql-42.2.22.jar!/META-INF/services/java.sql.Driver
        String url = "META-INF/services/java.sql.Driver";
        // 用某个 classLoader 对象, 在其查找范围内查找 url 枚举数组
        Enumeration<URL> enums = Thread.currentThread().getContextClassLoader().getResources(url);
        while(enums.hasMoreElements()){
            System.out.println(enums.nextElement());
        }
    }
    public static void main(String[] args) throws ClassNotFoundException, IOException {
        test();
// new ServiceLoader<>(service, loader);  loader 是 contextClassLoader, Thread.currentThread().getContextClassLoader()
        /**
         * providers = new LinkedHashMap<>();    // 链表, 遍历的时候向其中加入 provider
         *
         * lookupIterator = new LazyIterator(service, loader);  // 以上来, 先进入 lookupIterator 的 hasNext()
         *
         * loader = Thread.currentThread().getContextClassLoader()
         */
        ServiceLoader<Driver> driverServices = ServiceLoader.load(Driver.class);
        /**
         * knownProviders = providers.entrySet().iterator()
         */
        Iterator<Driver> driversIterator = driverServices.iterator();
        try{
            while(driversIterator.hasNext()) {
                Driver d = driversIterator.next();
                System.out.println(d.getClass().newInstance());
            }
        } catch(Throwable t) {
        }
    }
}

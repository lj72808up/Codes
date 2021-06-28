package com.test.loader;

import java.sql.Driver;
import java.util.Iterator;
import java.util.ServiceLoader;

public class Main {
    public static void main(String[] args) throws ClassNotFoundException {
        // todo 走一遍 spi 代码
        // https://blog.csdn.net/shi2huang/article/details/80308531?ops_request_misc=%257B%2522request%255Fid%2522%253A%2522162486057916780255255684%2522%252C%2522scm%2522%253A%252220140713.130102334..%2522%257D&request_id=162486057916780255255684&biz_id=0&utm_medium=distribute.pc_search_result.none-task-blog-2~all~sobaiduend~default-1-80308531.first_rank_v2_pc_rank_v29_1&utm_term=ServiceLoader&spm=1018.2226.3001.4187
//       Class.forName("com.mysql.cj.jdbc.Driver");
//       Class.forName("org.postgresql.Driver");



        ServiceLoader<Driver> driverServices = ServiceLoader.load(Driver.class);
        Iterator<Driver> driversIterator = driverServices.iterator();
        try{
            while(driversIterator.hasNext()) {
                Driver d = driversIterator.next();
                System.out.println(d.getClass());
            }
        } catch(Throwable t) {
        }
    }
}

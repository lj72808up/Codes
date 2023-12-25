package com2.test.proxy.cglib;

import org.springframework.cglib.core.DebuggingClassWriter;

public class TestMain {
    public static void main(String[] args) {
        // 获取当前项目的根目录
        String userDir = System.getProperty("user.dir");
        //System.setProperty("cglib.debugLocation", userDir);
        System.setProperty(DebuggingClassWriter.DEBUG_LOCATION_PROPERTY, userDir);
        // 1. 目标对象
        Service target = new Service();

        // 2. 根据目标对象生成代理对象
        MyInterceptor myInterceptor = new MyInterceptor(target);

        // 获取 CGLIB 代理类
        Service proxyObject = (Service) myInterceptor.getProxy();

        // 调用代理对象的方法
        proxyObject.finalMethod();
        proxyObject.publicMethod();
    }
}

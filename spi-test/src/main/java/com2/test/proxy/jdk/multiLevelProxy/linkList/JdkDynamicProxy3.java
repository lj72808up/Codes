package com2.test.proxy.jdk.multiLevelProxy.linkList;

import java.lang.reflect.InvocationHandler;
import java.lang.reflect.Method;
import java.lang.reflect.Proxy;
import java.util.ArrayList;
import java.util.List;

public class JdkDynamicProxy3 implements InvocationHandler {
    private final Object target;  // 目标对象
    private final List<MyMethodInterceptor> interceptorList = new ArrayList<>();  // 拦截器链

    public JdkDynamicProxy3(Object target) {
        this.target = target;
    }

    public void addMethodInterceptor(MyMethodInterceptor methodInterceptor) {
        this.interceptorList.add(methodInterceptor);
    }

    @Override
    public Object invoke(Object proxy, Method method, Object[] args) throws Throwable {
        // 被代理对象的方法，参数
        TargetMethod targetMethod = new TargetMethod(target,method,args);

        MyMethodInvocationImpl myMethodInvocation = new MyMethodInvocationImpl();
        myMethodInvocation.setTargetMethod(targetMethod);
        myMethodInvocation.setInterceptorList(interceptorList);

        return myMethodInvocation.proceed();
    }

    /**
     * 获取被代理接口实例对象
     * 通过 Proxy.newProxyInstance 可以获得一个代理对象，它实现了 target.getClass().getInterfaces() 接口
     */
    public Object getProxy() {
        return Proxy.newProxyInstance(target.getClass().getClassLoader(),
                target.getClass().getInterfaces(),
                this);
    }
}
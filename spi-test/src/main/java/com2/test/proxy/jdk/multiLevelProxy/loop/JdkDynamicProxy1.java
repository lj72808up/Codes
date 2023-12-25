package com2.test.proxy.jdk.multiLevelProxy.loop;

import java.lang.reflect.InvocationHandler;
import java.lang.reflect.Method;
import java.lang.reflect.Proxy;

public class JdkDynamicProxy1 implements InvocationHandler {
    /**
     *  目标对象（也被称为被代理对象）
     *  Java 代理模式的一个必要要素就是代理对象要能拿到被代理对象的引用
     */
    private final Object target;

    /**
     * 用来标识 InvocationHandler invoke 调用
     */
    public String tag;

    public JdkDynamicProxy1(Object target){
        this.target=target;
    }

    /**
     * 回调方法
     * @param proxy JDK 生成的代理对象
     * @param method 被代理的方法（也就是需要增强的方法）
     * @param args  被代理方法的参数
     */
    @Override
    public Object invoke(Object proxy, Method method, Object[] args) throws Throwable {
        // 不对 hashCode(）方法代理
        if (method.getName().equals("hashCode")) {
            return hashCode();
        }

        System.out.println("JdkDynamicProxy1 invoke 方法执行前---------------" + info(proxy));
        Object object= method.invoke(this.target, args); // 放行 target Object 的方法
        System.out.println("JdkDynamicProxy1 invoke 方法执行后----------------" + info(proxy));
        return object;
    }

    /**
     * 返回代理类的类名和 hash code
     */
    private String info(Object proxy) {
        return proxy.getClass().getName() + ":" + proxy.hashCode() + "----------------" + this.tag;
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
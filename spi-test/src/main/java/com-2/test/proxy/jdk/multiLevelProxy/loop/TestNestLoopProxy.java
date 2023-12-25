package com2.test.proxy.jdk.multiLevelProxy.loop;

import com2.test.proxy.jdk.Animal;
import com2.test.proxy.jdk.Cat;

public class TestNestLoopProxy {
    public static void main(String[] args) {
        // 1. 构造目标对象
        Cat catTarget = new Cat();

        // 2. 根据目标对象生成第一层代理对象
        JdkDynamicProxy1 proxy = new JdkDynamicProxy1(catTarget);
        proxy.tag = "第一个代理类";

        // JDK 动态代理是基于接口的，所以只能转换为 Cat 实现的接口 Animal
        Animal catProxy = (Animal) proxy.getProxy();

        // 3. 根据第一个代理对象生成代理对象
        JdkDynamicProxy1 proxy2 = new JdkDynamicProxy1(catProxy);
        proxy2.tag = "第二个代理类";

        // JDK 动态代理是基于接口的，所以只能转换为 Cat 实现的接口 Animal
        Animal catProxy2 = (Animal) proxy2.getProxy();

        // 调用 catProxy2 的方法
        catProxy2.eat();
    }
}

/**
 * 输出：
 * JdkDynamicProxy1 invoke 方法执行前---------------com.sun.proxy.$Proxy0:723074861----------------第二个代理类
 * JdkDynamicProxy1 invoke 方法执行前---------------com.sun.proxy.$Proxy0:895328852----------------第一个代理类
 * 猫吃鱼
 * JdkDynamicProxy1 invoke 方法执行后----------------com.sun.proxy.$Proxy0:895328852----------------第一个代理类
 * JdkDynamicProxy1 invoke 方法执行后----------------com.sun.proxy.$Proxy0:723074861----------------第二个代理类
 */
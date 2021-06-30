package com2.test.proxy.jdk.sptingAop;

import com2.test.proxy.jdk.Animal;
import com2.test.proxy.jdk.Cat;
import org.aopalliance.intercept.MethodInterceptor;
import org.aopalliance.intercept.MethodInvocation;
import org.springframework.aop.framework.ProxyFactory;


public class TestProxyFactory {
    public static void main(String[] args) {
        // 1. 构造目标对象
        Cat catTarget = new Cat();

        // 2. 通过目标对象，构造 ProxyFactory 对象
        ProxyFactory factory = new ProxyFactory(catTarget);

        // 添加一个方法拦截器
        factory.addAdvice(new MyMethodInterceptor());

        // 3. 根据目标对象生成代理对象
        Object proxy = factory.getProxy();

        Animal cat = (Animal) proxy;
        cat.eat();
    }

    public static class MyMethodInterceptor implements MethodInterceptor {

        @Override
        public Object invoke(MethodInvocation invocation) throws Throwable {
            System.out.println("MyMethodInterceptor invoke 调用 before invocation.proceed");

            Object ret = invocation.proceed();

            System.out.println("MyMethodInterceptor invoke 调用 after invocation.proceed");
            return ret;
        }
    }
}

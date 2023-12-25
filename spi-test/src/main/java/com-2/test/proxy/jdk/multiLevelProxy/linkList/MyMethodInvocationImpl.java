package com2.test.proxy.jdk.multiLevelProxy.linkList;

import java.util.List;

public class MyMethodInvocationImpl implements MyMethodInvocation {

    private TargetMethod targetMethod; // 被代理的目标方法
    private List<MyMethodInterceptor> interceptorList; // 拦截器链链
    private int currentInterceptorIndex = 0;  // 当前调用的拦截器索引


    public void setInterceptorList(List<MyMethodInterceptor> interceptorList) {
        this.interceptorList = interceptorList;
    }

    public void setTargetMethod(TargetMethod targetMethod) {
        this.targetMethod = targetMethod;
    }

    @Override
    public Object proceed() throws Throwable {
        /**
         *  索引值从 0 开始递增，所以如果 currentInterceptorIndex 等于拦截器集合大小，
         *  说明所有的拦截器都执行完毕了, 可调用目标对象的方法了
         */
        if (this.currentInterceptorIndex == this.interceptorList.size()) {
            // 调用目标方法
            return targetMethod.getMethod().invoke(targetMethod.getTarget(), targetMethod.getArgs());
        }

        /**
         * 调用当前拦截器的 invoke 方法，并让拦截器列表指针向后移动
         */
        MyMethodInterceptor methodInterceptor = this.interceptorList.get(this.currentInterceptorIndex);
        this.currentInterceptorIndex++;
        return methodInterceptor.invoke(this);  // 把自己的调用链传过去
    }
}
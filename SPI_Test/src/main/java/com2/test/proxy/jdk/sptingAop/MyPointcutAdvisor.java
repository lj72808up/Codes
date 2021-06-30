package com2.test.proxy.jdk.sptingAop;

import org.aopalliance.aop.Advice;
import org.springframework.aop.Pointcut;
import org.springframework.aop.PointcutAdvisor;

public class MyPointcutAdvisor implements PointcutAdvisor {

    private final Pointcut pointcut;

    private final Advice advice;

    public MyPointcutAdvisor(Advice advice,Pointcut pointcut) {
        this.advice = advice;
        this.pointcut = pointcut;
    }

    @Override
    public Pointcut getPointcut() {
        return this.pointcut;
    }

    @Override
    public Advice getAdvice() {
        return this.advice;
    }

    /**
     * 此方法暂时忽略，不需要理会
     */
    @Override
    public boolean isPerInstance() {
        return false;
    }
}

package comm.test.proxy.jdk.multiLevelProxy.linkList;

public interface MyMethodInvocation {

    /**
     * 进入拦截器链中的下一个拦截器，驱动责任链模式向前运行
     */
    Object proceed() throws Throwable;
}
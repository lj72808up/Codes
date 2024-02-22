package comm.test.proxy.jdk.multiLevelProxy.linkList;


/**
 *  在到达目标方法之前拦截对方法的调用。
 */
public interface MyMethodInterceptor {

    /**
     * 类似于 InvocationHandler 的 invoke 方法，用于对方法做增强处理，并通过 invocation 参数驱动责任链向前运行

     */
    Object invoke(MyMethodInvocation invocation) throws Throwable;
}
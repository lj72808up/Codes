package components;

import java.util.concurrent.Executor;

/**
 * 有了 components.EventBus，components.AsyncEventBus 的实现就非常简单了。为了实现异步非阻塞的观察者模式，它就不能再继续使用 MoreExecutors.directExecutor() 了，
 * 而是需要在构造函数中，由调用者注入线程池。
 */
public class AsyncEventBus extends EventBus {
    public AsyncEventBus(Executor executor) {
        super(executor);
    }
}
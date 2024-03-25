package components;

import com.google.common.util.concurrent.MoreExecutors;

import java.util.List;
import java.util.concurrent.Executor;

/**
 * components.EventBus 实现的是阻塞同步的观察者模式。看代码你可能会有些疑问，这明明就用到了线程池 Executor 啊。
 * 实际上，MoreExecutors.directExecutor() 是 Google Guava 提供的工具类，看似是多线程，实际上是单线程。
 * 之所以要这么实现，主要还是为了跟 components.AsyncEventBus 统一代码逻辑，做到代码复用
 */
public class EventBus {
    private Executor executor;
    // 不再用 List<> 记录观察者, 取而代之是用一个表记录观察者和其观察的消息类型(方法的形参类型). 执行通知时先插表再反射执行
    private ObserverRegistry registry = new ObserverRegistry();

    public EventBus() {
        this(MoreExecutors.directExecutor());
    }

    protected EventBus(Executor executor) {
        this.executor = executor;
    }

    public void register(Object object) {
        registry.register(object);
    }

    public void post(Object event) {
        List<ObserverAction> observerActions = registry.getMatchedObserverActions(event);
        for (ObserverAction observerAction : observerActions) {
            executor.execute(new Runnable() {
                @Override
                public void run() {
                    observerAction.execute(event);
                }
            });
        }
    }
}
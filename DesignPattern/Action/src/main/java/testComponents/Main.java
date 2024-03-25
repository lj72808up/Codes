package testComponents;

import components.AsyncEventBus;
import components.EventBus;
import testComponents.observers.RegNotificationObserver;
import testComponents.observers.RegPromotionObserver;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

/**
 * 有了事件,就往总线上发. 注册在总线上的观察者, 自己的 @Subscribe 方法会自动被总线捕获, 得到参数类型记录在总线自己的 <形参,观察类> 列表中.
 * 总线会自动根据事件类型调用观察者的方法, 调用的执行是在总线自己的线程池中. 所以总线能同时处理的事件条数 = size(总线线程池)/length(注册的观察者列表)
 */
public class Main {
    public static void main(String[] args) {
//        EventBus eventBus = new EventBus();
        ExecutorService ecs = Executors.newFixedThreadPool(1);
        AsyncEventBus eventBus = new AsyncEventBus(ecs);
        eventBus.register(new RegNotificationObserver());
        eventBus.register(new RegPromotionObserver());

        Long userid1 = 200L;    // 和观察者的 @Subscribe 方法形参类型一致
        eventBus.post(userid1);

        Long userid2 = 300L;    // 和观察者的 @Subscribe 方法形参类型一致
        eventBus.post(userid2);
    }
}

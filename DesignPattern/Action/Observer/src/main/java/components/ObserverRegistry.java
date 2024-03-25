package components;

import com.google.common.base.Preconditions;

import java.lang.reflect.Method;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;
import java.util.concurrent.CopyOnWriteArraySet;

/**
 * components.ObserverRegistry 类就是前面讲到的 Observer 注册表，是最复杂的一个类，框架中几乎所有的核心逻辑都在这个类中。
 * 这个类大量使用了 Java 的反射语法，不过代码整体来说都不难理解，其中，一个比较有技巧的地方是 CopyOnWriteArraySet 的使用。
 * CopyOnWriteArraySet，顾名思义，在写入数据的时候，会创建一个新的 set，并且将原始数据 clone 到新的 set 中，在新的 set 中写入数据完成之后，再用新的 set 替换老的 set。
 * 这样就能保证在写入数据的时候，不影响数据的读取操作，以此来解决读写并发问题。
 * 除此之外，CopyOnWriteSet 还通过加锁的方式，避免了并发写冲突。具体的作用你可以去查看一下 CopyOnWriteSet 类的源码，一目了然
 *
 * 至此，我们用了不到 200 行代码，就实现了一个还算凑活能用的 components.EventBus，从功能上来讲，它跟 Google Guava components.EventBus 几乎一样。
 * 不过，如果去查看Google Guava components.EventBus 的源码，你会发现，在实现细节方面，相比我们现在的实现，它其实做了很多优化，
 * 比如优化了在注册表中查找消息可匹配函数的算法。如果有时间的话，建议你去读一下它的源码。
 */
public class ObserverRegistry {
    private ConcurrentMap<Class<?>, CopyOnWriteArraySet<ObserverAction>> registry = new ConcurrentHashMap<>();

    public void register(Object observer) {
        Map<Class<?>, Collection<ObserverAction>> observerActions = findAllObserverActions(observer);
        for (Map.Entry<Class<?>, Collection<ObserverAction>> entry : observerActions.entrySet()) {
            Class<?> eventType = entry.getKey();
            Collection<ObserverAction> eventActions = entry.getValue();
            CopyOnWriteArraySet<ObserverAction> registeredEventActions = registry.get(eventType);
            if (registeredEventActions == null) {
                registry.putIfAbsent(eventType, new CopyOnWriteArraySet<>());
                registeredEventActions = registry.get(eventType);
            }
            registeredEventActions.addAll(eventActions);
        }
    }

    public List<ObserverAction> getMatchedObserverActions(Object event) {
        List<ObserverAction> matchedObservers = new ArrayList<>();
        Class<?> postedEventType = event.getClass();
        for (Map.Entry<Class<?>, CopyOnWriteArraySet<ObserverAction>> entry : registry.entrySet()) {
            Class<?> eventType = entry.getKey();
            Collection<ObserverAction> eventActions = entry.getValue();
            if (postedEventType.isAssignableFrom(eventType)) {
                matchedObservers.addAll(eventActions);
            }
        }
        return matchedObservers;
    }

    private Map<Class<?>, Collection<ObserverAction>> findAllObserverActions(Object observer) {
        Map<Class<?>, Collection<ObserverAction>> observerActions = new HashMap<>();
        Class<?> clazz = observer.getClass();
        for (Method method : getAnnotatedMethods(clazz)) {
            Class<?>[] parameterTypes = method.getParameterTypes();
            Class<?> eventType = parameterTypes[0];
            if (!observerActions.containsKey(eventType)) {
                observerActions.put(eventType, new ArrayList<>());
            }
            observerActions.get(eventType).add(new ObserverAction(observer, method));
        }
        return observerActions;
    }

    private List<Method> getAnnotatedMethods(Class<?> clazz) {
        List<Method> annotatedMethods = new ArrayList<>();
        for (Method method : clazz.getDeclaredMethods()) {  // 遍历 Class 内定义的函数
            if (method.isAnnotationPresent(Subscribe.class)) { // 判断函数是否存在 @Subscribe 注解
                Class<?>[] parameterTypes = method.getParameterTypes();  // 拿到 @Subscribe 注解的方法的参数类型
                Preconditions.checkArgument(parameterTypes.length == 1, // 限定方法只有一个形参
                        "Method %s has @components.Subscribe annotation but has %s parameters."
                                + "Subscriber methods must have exactly 1 parameter.",
                        method, parameterTypes.length);
                annotatedMethods.add(method);
            }
        }
        return annotatedMethods;
    }
}
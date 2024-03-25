package components;

import com.google.common.base.Preconditions;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;

/**
 * components.ObserverAction 类用来表示 @components.Subscribe 注解的方法，
 * target 表示观察者类的对象，method 表示观察者的方法。它主要用在 components.ObserverRegistry 观察者注册表中
 */
public class ObserverAction {
    private final Object target;
    private final Method method;

    public ObserverAction(Object target, Method method) {
        this.target = Preconditions.checkNotNull(target);
        this.method = method;
        this.method.setAccessible(true);
    }

    public void execute(Object event) { // event是method方法的参数
        try {
            method.invoke(target, event);
        } catch (InvocationTargetException | IllegalAccessException e) {
            e.printStackTrace();
        }
    }
}
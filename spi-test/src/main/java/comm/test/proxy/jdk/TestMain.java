package comm.test.proxy.jdk;

import java.lang.reflect.InvocationHandler;
import java.lang.reflect.Method;
import java.lang.reflect.Proxy;

public class TestMain {
    public static void main(String[] args) {
        // 设置此系统属性,让 JVM 生成的 Proxy 类写入文件. 保存路径为项目的根目录
        System.getProperties().put("sun.misc.ProxyGenerator.saveGeneratedFiles", "true");

        // 1. 构造目标对象
        Cat catTarget = new Cat();

        // 2. 根据目标对象生成代理对象
        MyDynamicProxy proxy = new MyDynamicProxy(catTarget);

        // JDK 动态代理是基于接口的，所以只能转换为 Cat 实现的接口 Animal
        Animal catProxy = (Animal) proxy.getProxy();

        // 调用代理对象的方法
        catProxy.eat();
        catProxy.sayHello("喵喵喵");

        catProxy.toString();

//        System.out.println("catTarget 的 hashCode:"+catTarget.hashCode());
        System.out.println("catProxy 和 catTarget 的 hashCode:"+(catProxy.hashCode() == catTarget.hashCode()));
    }
}

class MyDynamicProxy implements InvocationHandler {
    /**
     * 目标对象（也被称为被代理对象）
     * Java 代理模式的一个必要要素就是代理对象要能拿到被代理对象的引用
     */
    private final Object target;

    public MyDynamicProxy(Object target) {
        this.target = target;
    }

    /**
     * 回调方法
     *
     * @param proxy  JDK 生成的代理对象
     * @param method 被代理的方法（也就是需要增强的方法）
     * @param args   被代理方法的参数
     */
    @Override
    public Object invoke(Object proxy, Method method, Object[] args) throws Throwable {
        System.out.println("MyDynamicProxy invoke 方法执行前-------------------------------");
        Object object = method.invoke(target, args);
        System.out.println(target.getClass());
        System.out.println(proxy.getClass());
        System.out.println("MyDynamicProxy invoke 方法执行后-------------------------------");
        return object;
    }

    /**
     * JDK 动态代理的核心，为目标对象创建被代理后的对象
     * 通过 Proxy.newProxyInstance 可以获得一个代理对象，它实现了 target.getClass().getInterfaces() 的所有接口
     */
    public Object getProxy() {
        return Proxy.newProxyInstance(target.getClass().getClassLoader(),
                target.getClass().getInterfaces(),
                this);
    }
}
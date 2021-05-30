import org.springframework.aop.framework.ProxyFactory;

public class TestProxy {
    public static void main(String[] args) {
        ProxyFactory factory = new ProxyFactory(new SimplePojo());
        factory.addInterface(Pojo.class);
        factory.addAdvice(new RetryAdvice());

        factory.setExposeProxy(true);  // 暴露代理对象, 这个是关键
        Pojo pojo = (Pojo) factory.getProxy();
        pojo.foo();   // foo() 中的 this, 就会指向带对象
    }
}

// target object
class SimplePojo implements Pojo {

    public void foo() {
        this.bar();
    }

    public void bar() {
    }
}

interface Pojo {
    public void foo();
    public void bar();
}

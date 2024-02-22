package comm.test.proxy.jdk.springAop;

import comm.test.proxy.jdk.Animal;
import comm.test.proxy.jdk.Cat;
import org.aopalliance.intercept.MethodInterceptor;
import org.aopalliance.intercept.MethodInvocation;
import org.springframework.aop.framework.ProxyFactory;

public class TestPointAdvisor {
    public static void main(String[] args) {
        // 设置此系统属性,让 JVM 生成的 Proxy 类写入文件. 保存路径为项目的根目录
        System.getProperties().put("sun.misc.ProxyGenerator.saveGeneratedFiles", "true");

        // 1. 构造目标对象
        Animal catTarget = new Cat();

        // 2. 通过目标对象，构造 ProxyFactory 对象
        ProxyFactory factory = new ProxyFactory(catTarget);

        // 添加一个 Advice (DefaultPointcutAdvisor)
        factory.addAdvice(new MyMethodInterceptor());

        // 新增代码：添加一个 PointcutAdvisor
        MyPointcutAdvisor myPointcutAdvisor = new MyPointcutAdvisor(
                new MyMethodInterceptor(),  // Advice
                new MyPointcut()            // PointCut
        );
        factory.addAdvisor(myPointcutAdvisor);

        // 3. 根据目标对象生成代理对象
        Animal cat = (Animal) factory.getProxy();
        System.out.println(cat.getClass());
        cat.eat();

        System.out.println("---------------------------------------");

        cat.go();
    }

    public static class MyMethodInterceptor implements MethodInterceptor {

        @Override
        public Object invoke(MethodInvocation invocation) throws Throwable {
            System.out.println("MyMethodInterceptor invoke 调用 before invocation.proceed");

            Object ret = invocation.proceed();

            System.out.println("MyMethodInterceptor invoke 调用 after invocation.proceed");
            return ret;
        }
    }
}

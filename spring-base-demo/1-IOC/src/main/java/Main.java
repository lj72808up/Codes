import xml.instance.Example3;
import xml.ExampleBean;
import xml.FactoryExampleBean;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.context.support.ClassPathXmlApplicationContext;

public class Main {
    public static void main(String[] args) {
        ConfigurableApplicationContext context = new ClassPathXmlApplicationContext("application1.xml");
        context.registerShutdownHook();   // AbstractApplicationContext 中实现
        // 1. 构造器初始化对象
        ExampleBean bean1 = (ExampleBean) context.getBean("example1");
        System.out.println(bean1.toString());

        // 2. 静态工厂初始化对象
        FactoryExampleBean bean2 = (FactoryExampleBean) context.getBean("example2");
        System.out.println(bean2.toString());

        // 3. 实例工厂初始化对象
        Example3 bean3 = (Example3) context.getBean("example3");
        System.out.println(bean3.toString());

        System.out.println(context.getBeanFactory());
    }
}

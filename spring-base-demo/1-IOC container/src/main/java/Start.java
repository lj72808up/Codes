import exampleBeans.instance.Example3;
import exampleBeans.ExampleBean;
import exampleBeans.FactoryExampleBean;
import org.springframework.beans.factory.BeanFactory;
import org.springframework.context.ApplicationContext;
import org.springframework.context.support.ClassPathXmlApplicationContext;

public class Start {
    public static void main(String[] args) {
        ApplicationContext context = new ClassPathXmlApplicationContext("application1.xml");

        // 1. 构造器初始化对象
        ExampleBean bean1 = (ExampleBean) context.getBean("example1");
        System.out.println(bean1.toString());

        // 2. 静态工厂初始化对象
        FactoryExampleBean bean2 = (FactoryExampleBean) context.getBean("example2");
        System.out.println(bean2.toString());

        // 3. 实例工厂初始化对象
        Example3 bean3 = (Example3) context.getBean("example3");
        System.out.println(bean3.toString());
    }
}

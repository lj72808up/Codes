import jdk.internal.org.objectweb.asm.tree.MultiANewArrayInsnNode;
import org.springframework.beans.factory.support.DefaultListableBeanFactory;
import org.springframework.beans.factory.xml.XmlBeanDefinitionReader;
import org.springframework.core.io.ClassPathResource;
import xml.instance.Example3;
import xml.ExampleBean;
import xml.FactoryExampleBean;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.context.support.ClassPathXmlApplicationContext;

public class Main {
    public static void main(String[] args) {
        Main.testFactory();
//        Main.testBeanDefinition();
    }

    /**
     * 加载 BeanDefinition
     */
    public static void testBeanDefinition() {
        // <1> 获取资源
        ClassPathResource resource = new ClassPathResource("application1.xml");
        // <2> 获取 BeanFactory
        DefaultListableBeanFactory factory = new DefaultListableBeanFactory();
        // <3> 组合 BeanFactory 创建 BeanDefinitionReader, 该 Reader 为 Resource 的解析器
        XmlBeanDefinitionReader reader = new XmlBeanDefinitionReader(factory);
        // <4> 装载 Resource
        reader.loadBeanDefinitions(resource);
    }

    /**
     * spring 中各种工厂
     */
    public static void testFactory() {
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

        // DefaultListableBeanFactory
        System.out.println(context.getBeanFactory());
    }
}

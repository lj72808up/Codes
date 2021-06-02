package annotation;


import annotation.testImport.ConfigA;
import annotation.testImport.ServiceInterface;
import org.springframework.context.annotation.AnnotationConfigApplicationContext;

public class Main {
    public static void main(String[] args) {
        AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext();
        ctx.register(ConfigA.class);
        ctx.refresh();
        // out: ServiceB
        ctx.getBean(ServiceInterface.class).test();
    }
    public static void main1(String[] args) {
        AnnotationConfigApplicationContext ctx = new AnnotationConfigApplicationContext();
        ctx.register(AppConfig.class);
//        ctx.scan("annotation");
        ctx.refresh();
        OtherComponent c = ctx.getBean(OtherComponent.class);
        System.out.println(c);

        Foo f = ctx.getBean(Foo.class);
        Bar b = ctx.getBean(Bar.class);
        System.out.println(f.getBar() == b); // 返回 true
    }
}

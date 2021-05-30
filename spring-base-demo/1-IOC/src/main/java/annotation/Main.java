package annotation;


import org.springframework.context.annotation.AnnotationConfigApplicationContext;

public class Main {
    public static void main(String[] args) throws Exception{
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

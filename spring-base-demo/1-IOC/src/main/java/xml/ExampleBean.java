package xml;

public class ExampleBean {
    @Override
    public String toString() {
        return "hello world";
    }

    public void init() {
        System.out.println("初始化回调函数调用");
    }

    public ExampleBean(){
        System.out.println("空参构造完毕");
    }
}

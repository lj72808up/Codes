package comm.test.proxy.jdk.springAop;

public class TestExtend {
    public static void main(String[] args) {
        new B(1);
    }
}

class A{
    public A(){
        System.out.println("A init");
    }
    public A(int a){
        System.out.println("A:"+a);
    }
}

class B extends A{
    public B(int b){
//        super(); 构造函数会先调用父类无参构造. 如果向下面手动调用了父类的构造函数，则只走父类有残构造
//        super(b);
        System.out.println("B:"+b);
    }
}
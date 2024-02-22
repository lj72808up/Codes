package comm.test.proxy.jdk;

public class Cat implements Animal {
    @Override
    public void eat() {
        System.out.println("猫吃鱼");
    }

    @Override
    public void sayHello(String word) {
        System.out.println(word);
    }

    @Override
    public String toString() {
        System.out.println("我是小花猫");
        return "我是小花猫";
    }

    @Override
    public void go() {
        System.out.println("猫在跑");
    }
}
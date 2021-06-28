public class MyTest1 {
    public MyTest1() {
        // 让 MyTest1 由 AppClassLoader 加载
        System.out.println("MyTest1 classLoader is: " + this.getClass().getClassLoader());
//        new MyTest2();
    }
}

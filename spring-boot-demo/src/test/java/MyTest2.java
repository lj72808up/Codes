public class MyTest2 {
    // 让 MyTest2 由自定义加载器加载
    public MyTest2() {
        System.out.println("MyTest2 classLoader is: " + this.getClass().getClassLoader());
//        new MyTest1();
    }
}

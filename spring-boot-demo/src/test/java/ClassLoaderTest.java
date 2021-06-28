import java.io.*;

public class ClassLoaderTest extends ClassLoader {
    public String baseUrl;

    // findClass() 方法只有在 AppClassLoader 找不到时才调用
    // 所以要想自定义类加载器生效, class 文件不能放在 classPath 下
    @Override
    public Class<?> findClass(String className) {
        System.out.println("自定义findClass被调用...");
        String path = baseUrl + className.replace(".", "\\") + ".class";
        try {
            InputStream is = new FileInputStream(path);
            byte data[] = new byte[is.available()];
            is.read(data);
            return defineClass(className, data, 0, data.length);
        } catch (IOException ignored) {}
        return null;
    }

    private void setPath(String baseUrl) {
        this.baseUrl = baseUrl;
    }

    public ClassLoaderTest(){}
    public ClassLoaderTest(ClassLoader parent) {
        super(parent);
    }


    public static void main(String[] args) throws Exception {
        ClassLoaderTest loader2 = new ClassLoaderTest();
        loader2.setPath("/Users/liujie02/IdeaProjects/Codes/spring-boot-demo/target/myTest/");
        Class<?> c2 = loader2.loadClass("MyTest2");
        Object o2 = c2.newInstance();
        System.out.println();

        ClassLoaderTest loader1 = new ClassLoaderTest();
        loader1.setPath(".");//设置自定义类加载器的加载路径
        //被类加载器加载后，得到Class对象
        Class<?> c1 = loader1.loadClass("MyTest1");
        Object o1 = c1.newInstance();//实例化MyTest1
        System.out.println();


    }
}
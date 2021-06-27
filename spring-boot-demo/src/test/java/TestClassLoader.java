import java.io.IOException;
import java.io.InputStream;

public class TestClassLoader {
    public static void testLoaderLevel(){
        ClassLoader myLoader1 = new MyClassLoader();
        System.out.println(myLoader1);

        ClassLoader parent1 = myLoader1.getParent();
        System.out.println(parent1);

        ClassLoader parent2 = parent1.getParent();
        System.out.println(parent2);

        ClassLoader parent3 = parent2.getParent();
        System.out.println(parent3);

        System.out.println(ClassLoader.getSystemClassLoader() == parent1);
    }

    public static void main(String[] args) throws Exception {
        testLoaderLevel();
        System.out.println("===========================");
        ClassLoader myLoader1 = new MyClassLoader();

        Class<?> classA = myLoader1.loadClass("ClassA"); // class ClassA
        Object objA = classA.newInstance();
        System.out.println(objA.getClass());  // MyClassLoader@1218025c
        System.out.println(objA.getClass().getClassLoader());
        System.out.println(objA.getClass().getClassLoader().getParent().getParent().getParent().getParent());
        System.out.println(objA instanceof ClassA);

        ClassLoader myLoader2 = new MyClassLoader();
        Class<?> classA2 = myLoader2.loadClass("ClassA");
        Object objA2 = classA2.newInstance();
        System.out.println(objA2.getClass());
        System.out.println(objA2.getClass().getClassLoader());
        System.out.println(objA2.getClass().getClassLoader().getParent());
    }
}
class MyClassLoader extends ClassLoader{
    @Override
    public Class<?> loadClass(String name) throws ClassNotFoundException {
        if ("ClassA".equals(name)){
            try{
//                System.out.println(getClass().getResource("/ClassA.class"));
                InputStream is = getClass().getResourceAsStream("/ClassA.class");

                byte[] bytes = new byte[is.available()];
                is.read(bytes);
                return defineClass(name,bytes,0,bytes.length);
            }catch (IOException ignored){}
        }else {
            return super.loadClass(name);
        }
        return null;
    };
}
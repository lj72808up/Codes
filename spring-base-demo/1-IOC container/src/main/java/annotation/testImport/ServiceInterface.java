package annotation.testImport;

public interface ServiceInterface {
    void test();
}

class ServiceA implements ServiceInterface {

    @Override
    public void test() {
        System.out.println("ServiceA");
    }
}

class ServiceB implements ServiceInterface {

    @Override
    public void test() {
        System.out.println("ServiceB");
    }
}

class ServiceC implements ServiceInterface {
    private String name;

    public ServiceC(String name) {
        this.name = name;
    }

    @Override
    public void test() {
        System.out.println("ServiceC (name:)" + this.name);
    }
}
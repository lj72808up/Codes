package xml.instance;

public class InstanceFactory {
    private static InstanceFactory factory = new InstanceFactory();

    private InstanceFactory() {
    }

    public static InstanceFactory createInstance() {
        return factory;
    }

    public Example3 makeExample() {
        return new Example3();
    }
}
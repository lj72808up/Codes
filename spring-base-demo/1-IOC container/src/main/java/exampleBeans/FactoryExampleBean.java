package exampleBeans;

public class FactoryExampleBean {
    private static FactoryExampleBean example = new FactoryExampleBean();
    private FactoryExampleBean() {}
    public static FactoryExampleBean createInstance() {
        return example;
    }

    @Override
    public String toString() {
        return "hello static";
    }
}

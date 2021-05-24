import org.springframework.core.DefaultParameterNameDiscoverer;
import org.springframework.core.MethodParameter;
import org.springframework.core.ParameterNameDiscoverer;

import java.lang.reflect.Method;
import java.lang.reflect.Parameter;
import java.util.Arrays;

public class TestParamName {
    public static void main(String[] args) throws NoSuchMethodException {
        Method method = TestParamName.class.getMethod("test1", String.class, Integer.class);
        int parameterCount = method.getParameterCount();
        Parameter[] parameters = method.getParameters();

        // 打印输出：
        System.out.println("方法参数总数：" + parameterCount);
        Arrays.stream(parameters).forEach(
                p -> System.out.println(p.getType() + "----" + p.getName())
        );

    }

    public void test1(String name, Integer age) {
        System.out.println(name + "::" + age);
    }

    public static void main1(String[] args) throws NoSuchMethodException {
        Method method = TestParamName.class.getMethod("test1", String.class, Integer.class);
        MethodParameter nameParameter = new MethodParameter(method, 0);
        MethodParameter ageParameter = new MethodParameter(method, 1);

        // 打印输出：
        // 使用Parameter输出
        Parameter nameOriginParameter = nameParameter.getParameter();
        Parameter ageOriginParameter = ageParameter.getParameter();
        System.out.println("===================源生Parameter结果=====================");
        System.out.println(nameOriginParameter.getType() + "----" + nameOriginParameter.getName());
        System.out.println(ageOriginParameter.getType() + "----" + ageOriginParameter.getName());
        System.out.println("===================MethodParameter结果=====================");
        System.out.println(nameParameter.getParameterType() + "----" + nameParameter.getParameterName());
        System.out.println(ageParameter.getParameterType() + "----" + ageParameter.getParameterName());
        System.out.println("==============设置上ParameterNameDiscoverer后MethodParameter结果===============");
        ParameterNameDiscoverer parameterNameDiscoverer = new DefaultParameterNameDiscoverer();
        nameParameter.initParameterNameDiscovery(parameterNameDiscoverer);
        ageParameter.initParameterNameDiscovery(parameterNameDiscoverer);
        System.out.println(nameParameter.getParameterType() + "----" + nameParameter.getParameterName());
        System.out.println(ageParameter.getParameterType() + "----" + ageParameter.getParameterName());
    }
}

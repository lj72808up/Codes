package comm.test.proxy.jdk.multiLevelProxy.linkList;

import java.lang.reflect.Method;

public class TargetMethod {
    private final Object target;
    private final Method method;
    private final Object[] args;

    public TargetMethod(Object target, Method method, Object[] args) {
        this.target = target;
        this.method = method;
        this.args = args;
    }

    public Object getTarget() {
        return target;
    }

    public Object[] getArgs() {
        return args;
    }

    public Method getMethod() {
        return method;
    }
}

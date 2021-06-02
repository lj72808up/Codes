package annotation.testImport;

import org.springframework.beans.factory.support.BeanDefinitionBuilder;
import org.springframework.beans.factory.support.BeanDefinitionRegistry;
import org.springframework.context.annotation.*;
import org.springframework.core.type.AnnotationMetadata;
import org.springframework.util.MultiValueMap;

import java.lang.annotation.*;

//////////////// 自定义注解提供额外的 name 属性
@Retention(RetentionPolicy.RUNTIME)
@Documented
@Target(ElementType.TYPE)
//@Import(ServiceImportSelector.class)
@Import(ServiceBeanDefinitionRegistrar.class)
@interface MyEnableService{
    String name();
}

//////////////// ImportBeanDefinitionRegistrar 实现动态导入
class ServiceBeanDefinitionRegistrar implements ImportBeanDefinitionRegistrar {
    @Override
    public void registerBeanDefinitions(AnnotationMetadata metaData, BeanDefinitionRegistry registry) {
        MultiValueMap<String, Object> attributes = metaData.getAllAnnotationAttributes(MyEnableService.class.getName(), true);
        String nameAttr = attributes.get("name").get(0).toString();
        if ("C".equals(nameAttr)){
            BeanDefinitionBuilder builder = BeanDefinitionBuilder.rootBeanDefinition(ServiceC.class).
                    addConstructorArgValue("zhangsan");
            registry.registerBeanDefinition("serviceC", builder.getBeanDefinition());
        } else {
            BeanDefinitionBuilder builder = BeanDefinitionBuilder.rootBeanDefinition(ServiceB.class);
            registry.registerBeanDefinition("serviceB", builder.getBeanDefinition());
        }
    }
}

//////////////// 被 @MyEnableService 注解的 Config 类
@MyEnableService(name = "C")
@Configuration
public class ConfigA {}

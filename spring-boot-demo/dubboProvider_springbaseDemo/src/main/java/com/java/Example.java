package com.java;

//import com.java.config.JdbcConfig;
import com.java.config.Global;
import com.java.model.CasbinRule;
import com.java.service.CasbinService;
import org.springframework.boot.*;
import org.springframework.boot.autoconfigure.*;
import org.springframework.context.ApplicationContext;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.context.annotation.Bean;
import org.springframework.transaction.PlatformTransactionManager;
import org.springframework.transaction.annotation.EnableTransactionManagement;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

// @SpringBootApplication = @Configuration + @EnableAutoConfiguration + @ComponentScan
// 主要是 @Import(AutoConfigurationImportSelector.class) , 用于从 classpath 中搜寻所有 META-INF/spring.factories 配置文件，
// 并将其中 org.spring-framework.boot.autoconfigure.EnableAutoConfiguration 对应的配置项通过反射（Java Reflection）实例化为
// 对应的标注了 @Configuration 的 JavaConfig 形式的 IoC 容器配置类，然后汇总为一个并加载到 IoC 容器
@SpringBootApplication
@EnableTransactionManagement  // 开启事务2步: (1) 启动类加入 @EnableTransactionManagement 注解; (2) service 方法上加入 @Transactional 注解
public class Example {

    @Bean  // 因为是 single , 启动时加载一次, 打印 spring IOC 中所有 bean name
    public CommandLineRunner commandLineRunner(ApplicationContext ctx) {
        return args -> {

            System.out.println("Let's inspect the beans provided by Spring Boot:");

            String[] beanNames = ctx.getBeanDefinitionNames();
            Arrays.sort(beanNames);
            for (String beanName : beanNames) {
                System.out.println(beanName);
            }
        };
    }

    @Bean    // 测试 PlatformTransactionManager 集体是哪个
    public Object testBean(PlatformTransactionManager platformTransactionManager) {
        // org.springframework.jdbc.datasource.DataSourceTransactionManager
        System.out.println(">>>>>>>>>>" + platformTransactionManager.getClass().getName());
        return new Object();
    }


    //第三步:  main方法 和 类注解
    public static void main(String[] args) {
        // 运行main方法, 可以直接在ide运行, 或使用 mvn spring-boot:run 在项目根目录运行
        ConfigurableApplicationContext context = SpringApplication.run(Example.class);
        Global.context = context;

        /*CasbinService cas = context.getBean("casbinService", CasbinService.class);
        ArrayList<Integer> ids = new ArrayList<>();
        ids.add(21);
        ids.add(22);
        ids.add(25);
        List<CasbinRule> res = cas.selectByIds(ids);
        System.out.println(res);*/
    }
}
package com.java.mydubbo;

import com.java.service.DubboSayHello;
import com.java.service.MyResult;
import org.apache.dubbo.config.annotation.DubboService;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.PropertySource;

@DubboService(version = "1.0.0")
@Configuration
@PropertySource("classpath:application.yml")
public class Provider implements DubboSayHello {

    /**
     * The default value of ${dubbo.application.name} is ${spring.application.name}
     */
    @Value("${dubbo.application.name}")
    private String serviceName;

    public MyResult sayHello(String name) {
        String cont = String.format("[%s] : Hello, %s", this.serviceName, name);
        System.out.println(cont);
        return new MyResult(cont);
    }
}
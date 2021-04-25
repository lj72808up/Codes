package com.java.mydubbo;

import com.java.service.DubboSayHello;
import org.apache.dubbo.config.annotation.DubboReference;
import org.springframework.boot.ApplicationRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.EnableAutoConfiguration;
import org.springframework.context.annotation.Bean;

@EnableAutoConfiguration
public class ConsumerBootstrap {


    @DubboReference(version = "1.0.0", url = "dubbo://127.0.0.1:12345")
    private DubboSayHello demoService;

    public static void main(String[] args) {
        SpringApplication.run(ConsumerBootstrap.class).close();
    }

    @Bean
    public ApplicationRunner runner() {
        return args -> {
            System.out.println(demoService.sayHello("mercyblitz"));
        };
    }
}
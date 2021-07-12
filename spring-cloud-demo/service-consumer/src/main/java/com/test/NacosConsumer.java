package com.test;

import com.test.loadBalancer.RandomChooser;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.client.discovery.DiscoveryClient;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.web.client.RestTemplate;

@SpringBootApplication
public class NacosConsumer {
    public static void main(String[] args) {
        ConfigurableApplicationContext context = SpringApplication.run(NacosConsumer.class, args);
        RestTemplate restTemplate = context.getBean("restTemplate", RestTemplate.class);
        RestTemplate restTemplate2 = context.getBean("restTemplate2", RestTemplate.class);
        System.out.println(restTemplate);
        System.out.println(restTemplate2);
    }
}

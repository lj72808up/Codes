package com.test.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.cloud.client.ServiceInstance;
import org.springframework.cloud.client.discovery.DiscoveryClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.RestTemplate;

import java.util.List;

@RestController
public class HelloController {
    @Autowired // nacos 服务注册发现的 client
    private DiscoveryClient client;
    @Autowired
    private RestTemplate restTemplate;

    private static final String providerServiceName = "cloud_provider";

    @GetMapping("/hello")  // @RequestParam spring获取参数名然后比对 url 中的参数名
    public String hello(@RequestParam String name) {
        List<ServiceInstance> serviceInstances = client.getInstances(providerServiceName);
        ServiceInstance service = serviceInstances.get(0);
        String serviceUrl = "http://" + service.getHost() + ":" + service.getPort();
        return restTemplate.getForObject(String.format("%s/echo?name=%s",serviceUrl,name) , String.class);
    }
}

package com.test.controller;

import com.test.loadBalancer.RandomChooser;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.cloud.client.ServiceInstance;
import org.springframework.cloud.client.discovery.DiscoveryClient;
import org.springframework.cloud.client.loadbalancer.LoadBalanced;
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

    /**
     * 第一步: 集成 nacos,实现服务发现式调用
     */
    @GetMapping("/hello")  // @RequestParam spring获取参数名然后比对 url 中的参数名
    public String hello(@RequestParam String name) {
        List<ServiceInstance> serviceInstances = client.getInstances(providerServiceName);
        ServiceInstance service = serviceInstances.get(0);
        String serviceUrl = "http://" + service.getHost() + ":" + service.getPort();
        return restTemplate.getForObject(String.format("%s/echo?name=%s",serviceUrl,name) , String.class)+"\n";
    }

    /**
     * 负载均衡调用实现一: 显式使用 ServiceInstanceChooser 接口选择 RestTemplate
     */
    @Autowired
    private RandomChooser chooser;
    @GetMapping("/helloChooser")  // @RequestParam spring获取参数名然后比对 url 中的参数名
    public String helloBalance() {
//        List<ServiceInstance> serviceInstances = client.getInstances(providerServiceName);
        ServiceInstance service = chooser.choose(providerServiceName);
        String serviceUrl = "http://" + service.getHost() + ":" + service.getPort();
        return serviceUrl+"\n";
    }


    /**
     * 负载均衡调用实现二: 用加上 @LoadBalance 的 RestTemplate 自动实现负载均衡
     */
    @Autowired
    @Qualifier("restTemplate2")
    private RestTemplate restTemplate2;
    @GetMapping("/helloChooser2")  // @RequestParam spring获取参数名然后比对 url 中的参数名
    public String helloBalance2() {
//        List<ServiceInstance> serviceInstances = client.getInstances(providerServiceName);
        System.out.println(restTemplate2);
        return restTemplate2.getForObject("http://" + providerServiceName + "/", String.class)+"\n";
    }
}

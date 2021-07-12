package com.test.config;

import com.test.loadBalancer.RandomChooser;
import com.test.loadBalancer.loadBalanceAnnotation.RandomLoadBalancer;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.cloud.client.ServiceInstance;
import org.springframework.cloud.client.discovery.DiscoveryClient;
import org.springframework.cloud.client.loadbalancer.LoadBalanced;
import org.springframework.cloud.loadbalancer.core.ReactorLoadBalancer;
import org.springframework.cloud.loadbalancer.core.ServiceInstanceListSupplier;
import org.springframework.cloud.loadbalancer.support.LoadBalancerClientFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.env.Environment;
import org.springframework.web.client.RestTemplate;

@Configuration
public class BenConfig {
    @Bean  // 不带负载均衡的 RestTemplate, 可以集合自定义的 ServiceChooser 实现负载均衡
    public RestTemplate restTemplate(){
        return new RestTemplate();
    }

    @Bean
    public RandomChooser randomChooser(DiscoveryClient client){
        return new RandomChooser(client);
    }

    @Bean
    @LoadBalanced // 用拦截器实现的自动负载均衡 RestTemplate
    public RestTemplate restTemplate2(){
        return new RestTemplate();
    }

    @Bean
    public ReactorLoadBalancer<ServiceInstance> reactorServiceInstanceLoadBalancer(
            Environment environment,
            LoadBalancerClientFactory loadBalancerClientFactory) {
        String name = environment.getProperty(LoadBalancerClientFactory.PROPERTY_NAME);
        return new RandomLoadBalancer(loadBalancerClientFactory.getLazyProvider(name,
                ServiceInstanceListSupplier.class), name);
    }
}

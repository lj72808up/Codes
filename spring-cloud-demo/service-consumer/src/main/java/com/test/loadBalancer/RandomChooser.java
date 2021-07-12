package com.test.loadBalancer;

import org.springframework.cloud.client.ServiceInstance;
import org.springframework.cloud.client.discovery.DiscoveryClient;
import org.springframework.cloud.client.loadbalancer.ServiceInstanceChooser;

import java.util.List;
import java.util.Random;

/**
 * 这种自己实现 ServiceInstanceChooser 对象的办法太笨
 *
 */
public class RandomChooser implements ServiceInstanceChooser {

    private final DiscoveryClient client;
    private final Random random;

    public DiscoveryClient getClient() {
        return client;
    }

    public RandomChooser(DiscoveryClient client) {
        this.client = client;
        this.random = new Random();
    }

    @Override  // 随机选择一个实例
    public ServiceInstance choose(String serviceId) {
        List<ServiceInstance> serviceInstances = client.getInstances(serviceId);
        int size = serviceInstances.size();
        return serviceInstances.get(random.nextInt(size));
    }
}

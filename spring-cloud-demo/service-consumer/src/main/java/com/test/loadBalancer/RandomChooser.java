package com.test.loadBalancer;

import org.springframework.cloud.client.ServiceInstance;
import org.springframework.cloud.client.loadbalancer.ServiceInstanceChooser;

public class RandomChooser implements ServiceInstanceChooser {
    @Override
    public ServiceInstance choose(String serviceId) {
        return null;
    }
}

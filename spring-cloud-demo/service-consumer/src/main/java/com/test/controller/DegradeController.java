package com.test.controller;

import com.alibaba.cloud.circuitbreaker.sentinel.SentinelCircuitBreakerFactory;
import com.alibaba.cloud.circuitbreaker.sentinel.SentinelConfigBuilder;
import com.alibaba.cloud.sentinel.annotation.SentinelRestTemplate;
import com.alibaba.csp.sentinel.slots.block.RuleConstant;
import com.alibaba.csp.sentinel.slots.block.degrade.DegradeException;
import com.alibaba.csp.sentinel.slots.block.degrade.DegradeRule;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.cloud.client.circuitbreaker.CircuitBreaker;
import org.springframework.cloud.client.circuitbreaker.CircuitBreakerFactory;
import org.springframework.cloud.client.circuitbreaker.Customizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpStatus;
import org.springframework.http.client.ClientHttpResponse;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.ResponseErrorHandler;
import org.springframework.web.client.RestTemplate;

import java.io.IOException;
import java.util.Collections;
import java.util.Map;

@RestController
public class DegradeController {
    @Autowired // 引入 sentinal starter 后会自动创建 CircuitBreakerFactory
    private CircuitBreakerFactory circuitBreakerFactory;
    @Autowired
    private RestTemplate restTemplate;

    /**
     * 降级方式一: 用 CircuitBreaker.run() 调用 restTemplate
     * 需要用 @Bean 注解 Customizer<SentinelCircuitBreakerFactory> 配置降级策略
     *
     * @return
     */
    @GetMapping("degrade")
    public String degrade() {
        CircuitBreaker breaker = circuitBreakerFactory.create("temp");
        final String url = "https://httpbin.org/status/500"; // httpbin.org 经常用来模拟请求
        StringBuilder stringBuilder = new StringBuilder();
        for (int i = 0; i < 10; i++) {
            String res = breaker.run(
                    () -> {
                        return restTemplate.getForObject(url, String.class);
                    },
                    throwable -> {
                        if (throwable instanceof DegradeException) {
                            return "degrade by sentinel";
                        } else {
                            return "exception by url:" + url;
                        }
                    }
            );
            stringBuilder.append(res).append("\n");
        }
        return stringBuilder.toString();
    }

    /**
     * 降级方式二: RestTemplate 构造器中加入 ResponseHandler (未调通)
     * 需要用 @SentinelRestTemplate 注解 RestTemplate, 实现自动降级
     */
    @Autowired
    @Qualifier("restTemplateDegrade")
    private RestTemplate restTemplateDegrade;

    @GetMapping("degrade2")
    public String degrade2() {
        final String url = "https://httpbin.org/status/500"; // httpbin.org 经常用来模拟请求
        StringBuilder stringBuilder = new StringBuilder();
        for (int i = 0; i < 10; i++) {
            String res = restTemplateDegrade.getForObject(url, String.class);
            stringBuilder.append(res).append("\n");
        }
        return stringBuilder.toString();
    }
}

@Configuration
class SentialConfig {
    @Bean
    public Customizer<SentinelCircuitBreakerFactory> customizer() {
        return factory -> {
            // (1) 失败次数降级
            factory.configureDefault(id ->
                    new SentinelConfigBuilder()
                            .resourceName(id)
                            .rules(Collections.singletonList(
                                    new DegradeRule(id).setGrade(RuleConstant.DEGRADE_GRADE_EXCEPTION_COUNT)
                                            .setCount(3).setTimeWindow(10) // 失败次数3, 时间窗口10
                            )).build());
            // (2) 平均响应时间查过100ms, 则接下来5秒内, 都会自动降级
            factory.configure(builder -> {
                builder.rules(Collections.singletonList(
                        new DegradeRule("slow").setGrade(RuleConstant.DEGRADE_GRADE_RT)
                                .setCount(100).setTimeWindow(5) // 平均响应时间查过100ms, 则接下来5秒内, 都会自动降级
                ));
            }, "rt");
        };
    }

    @Bean
    @SentinelRestTemplate
    public RestTemplate restTemplateDegrade() {
        RestTemplate restTemplate = new RestTemplate();
        restTemplate.setErrorHandler(new ResponseErrorHandler() {
            @Override
            public boolean hasError(ClientHttpResponse response) throws IOException {
                return response.getStatusCode() == HttpStatus.INTERNAL_SERVER_ERROR;
            }

            @Override
            public void handleError(ClientHttpResponse response) throws IOException {
                throw new IllegalStateException("illegal state code 500");
            }
        });
        return restTemplate;
    }
}
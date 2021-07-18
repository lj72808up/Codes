package com.test.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.cloud.context.config.annotation.RefreshScope;
import org.springframework.context.ApplicationContext;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.HashMap;
import java.util.Map;

@RestController
@RefreshScope  // 让 bean scope 设置为 refresh, nacos 修改配置后会发布 RefreshEvent 事件, 监听器会销毁所有 scope 为 refresh 的 bean
public class ConfigController {
    @Autowired
    ApplicationContext applicationContext;
    @Value("${book.category:unknown}")
    private String category;

    @GetMapping("config")
    public Map<String,String> config() {
        String configVal = applicationContext.getEnvironment().getProperty("book.category", "unknown");
        HashMap<String, String> res = new HashMap<>();
        res.put("config",configVal);
        res.put("self.attr",this.category);
        return res;
    }
}

package com.test.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;


@RestController
public class EchoController {
    @GetMapping("/echo")
    public String echo(@RequestParam String name) {
        // @RequestParam: 问号后面的参数
        return "echo: " + name;
    }
}

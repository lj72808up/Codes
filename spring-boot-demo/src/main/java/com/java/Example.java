package com.java;

import org.springframework.boot.*;
import org.springframework.boot.autoconfigure.*;
import org.springframework.web.bind.annotation.*;

@RestController     // 声明该类是一个controller, 并将方法的返回值 render 成 json string 后返给调用者
@EnableAutoConfiguration
public class Example {

    @RequestMapping("/")
    public String home(){
        return "hello my friends";
    }

    @RequestMapping("/person")
    public Person p(){
        return new Person("zhangsan",23);
    }

    //第三步:  main方法 和 类注解
    public static void main(String[] args) {
        // 运行main方法, 可以直接在ide运行, 或使用 mvn spring-boot:run 在项目根目录运行
        SpringApplication.run(Example.class);
    }
}

class Person {
    public String name;
    public int age;

    public Person(String name, int age) {
        this.name = name;
        this.age = age;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public int getAge() {
        return age;
    }

    public void setAge(int age) {
        this.age = age;
    }
}
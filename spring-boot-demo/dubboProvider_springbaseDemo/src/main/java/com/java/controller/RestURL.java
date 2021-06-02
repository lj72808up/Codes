package com.java.controller;

import com.java.exception.BizException;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/test")
// 表示使用 spring mvc 处理 web 请求, 并将方法的返回值 render 成 json string 后返给调用者  (@ResponseBody)
public class RestURL {

    /**
     * spring boot 全局异常处理
     */
    @RequestMapping(value = "/")   // 不加 method 参数, 则不区分访问是 get 还是 post
    public String home() {
        System.out.println("fangwen");
        throw new BizException("发生运行时异常 biz");
//        throw new NullPointerException("发生运行时异常 null");
//        return "hello my friends";
    }

    @RequestMapping(value = "/person", method = {RequestMethod.GET})
    public Person p() {
        return new Person("zhangsan", 23);
    }

    @RequestMapping(value = "/person/{id}")   // 获取路径中的参数
    public String method7(@PathVariable("id") int id) {    // 匹配到的参数放到参数中
        return "get id=" + id;
    }

    // @GetMapping = @RequestMapping(method={HttpMethod.Get})
    @RequestMapping(value = "/p")   // 获取问号后面的参数
    public String allUsers(@RequestParam String name, @RequestParam int age) {
        return "name:" + name + ",age:" + age;
    }


    /**
     * 要求请求的 header 中 Content-Type 必须是 `application/json`
     */
    @PostMapping(value = "/postPerson")
    public Person PostPerson(@RequestBody Person p) {   // 从 post 的 body 中获取参数
        return p;
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
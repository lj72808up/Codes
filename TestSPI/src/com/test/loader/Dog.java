package com.test.loader;

public class Dog implements HelloService {

    @Override
    public void sayHello() {
        System.out.println("dag dag dag...");
    }
}
package com.java.service;

import java.io.Serializable;

public class MyResult implements Serializable {
    private String content ;

    public MyResult(String content) {
        this.content = content;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }
}
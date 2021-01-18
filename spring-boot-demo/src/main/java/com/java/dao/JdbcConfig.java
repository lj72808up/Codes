package com.java.dao;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;
import org.springframework.stereotype.Component;

@Configuration
public class JdbcConfig {
    @Value("${jdbc.mysql.user}")
    private String user;
    @Value("${jdbc.mysql.password}")
    private String password;
}

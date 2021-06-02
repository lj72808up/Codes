package com.java.exception;

import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import javax.servlet.http.HttpServletRequest;

@RestControllerAdvice
public class GlobalExceptionHandler {

    /**
     * 处理自定义的业务异常
     */
    @ExceptionHandler(value = BizException.class)
    public RespFacade bizExceptionHandler(HttpServletRequest req, BizException e) {
        System.out.println("发生biz异常！原因是:" + e.getMessage());
        return new RespFacade(1, "biz");
    }

    /**
     * 处理空指针异常
     */
    @ExceptionHandler(value = NullPointerException.class)
    public RespFacade nullExceptionHandler(HttpServletRequest req, NullPointerException e) {
        System.out.println("发生空指针异常！原因是:" + e.getMessage());
        return new RespFacade(1, "空指针");
    }

    /**
     * 处理其它异常
     */
    @ExceptionHandler(value = Exception.class)
    public RespFacade exceptionHandler(HttpServletRequest req, NullPointerException e) {
        System.out.println("发生异常！原因是:" + e.getMessage());
        return new RespFacade(1, "其它");
    }
}

class RespFacade {
    private int code;
    private String msg;

    public int getCode() {
        return code;
    }

    public void setCode(int code) {
        this.code = code;
    }

    public String getMsg() {
        return msg;
    }

    public void setMsg(String msg) {
        this.msg = msg;
    }

    public RespFacade(int code, String msg) {
        this.code = code;
        this.msg = msg;
    }
}

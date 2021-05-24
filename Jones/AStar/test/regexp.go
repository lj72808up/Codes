package test

import (
    "fmt"
    "regexp"
)

func main() {

    buf := "select * from t1 {aaaa} \\n {{}}where dt>={{startTime}} and dt<={{endTime}} c"

    //解析正则表达式，如果成功返回解释器
    reg1 := regexp.MustCompile(`\{\{([^}]+)\}\}`)
    //reg1 := regexp.MustCompile(`\{\{(.*?)\}\}`)
    if reg1 == nil {
        fmt.Println("regexp err")
        return
    }

    //根据规则提取关键信息
    result1 := reg1.FindAllString(buf,-1)
    fmt.Println("result1 = ", result1)
}
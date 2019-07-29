package com.test

import com.test.core.FreemarkerYielder
import java.util

object Test{
  def main(args: Array[String]): Unit = {
    val yielder = new FreemarkerYielder
    val map = new util.HashMap[String,String]
    map.put("magic1","111111111")
    map.put("v2xid1","111111111")
    map.put("imei1","111111111")
    map.put("mac1","111111111")
    val str = yielder.getCongtent("template","html",map,"test.html")
//    EmailUtil.sendEmail("mail2-in.baidu.com","25","liujie32@baidu.com","liujie32@baidu.com","test",str)
  }
}

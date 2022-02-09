package com.tv.sohu.spark.streaming.test

import java.io.{BufferedReader, FileReader, InputStreamReader}

import org.apache.hadoop.conf.Configuration
import org.apache.hadoop.fs.{FileSystem, Path}

/**
  * Created by wentaoli213587 on 2019/8/14.
  */
object FileTest {
  //大数据中心的应用不能在spark程序中读取hdfs文件，除非是yarn-client模式
  def main(args: Array[String]): Unit = {
    val fs = FileSystem.get(new Configuration())
    val path = new Path("/user/vetl/wentaoli213587/xiaomi/data.txt")
    val br = new BufferedReader(new InputStreamReader(fs.open(path), "UTF8"))
    println(br.readLine())
  }

}

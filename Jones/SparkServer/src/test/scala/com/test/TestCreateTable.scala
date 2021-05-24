package com.test

import java.net.{URL, URLClassLoader}

import org.apache.spark.sql.SparkSession
import utils.dag.{CkParser, JdbcOption, MySqlParser}

object TestCreateTable {
  def main1(args: Array[String]): Unit = {
    val spark = SparkSession.builder().master("local")
      .getOrCreate()
    import spark.implicits._
    val df = Seq(
      (1, "zhangsan", 23),
      (2, "lisi", 24)
    ).toDF("id", "name", "age")

//    var url = "jdbc:mysql://10.139.36.81:3306/adtl_test?user=yuhancheng&password=yuhancheng"
//    url = s"$url&useSSL=false&useUnicode=true&characterEncoding=UTF-8"
//    MySqlParser.createTableFromSql(df,url, "test_a",Map[String,String]("id"->"我就是ID"),"我想见一个表",Array[String]("date","hour"))
//    CkParser.createTableFromSql(df, "default.test_0925", Array("name", "age"), Array("name", "id"))
  }

  def main(args: Array[String]): Unit = {
    val myJarFiles = "file:///Users/liujie02/Downloads/sogou-adtl-dw-tool-udf-assembly-0.1.jar"
    val cl = URLClassLoader.newInstance(Array[URL](new URL(myJarFiles)))
    val myClass = cl.loadClass("com.sogou.adtl.dw.tool.udf.UdfTemplate");
    println(myClass.getClass)
  }



}
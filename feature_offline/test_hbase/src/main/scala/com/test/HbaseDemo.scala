package com.test

import org.apache.spark.sql.execution.datasources.hbase.HBaseTableCatalog
import org.apache.spark.sql.SparkSession

/*
object Encoder{
  implicit val itemEncoder = org.apache.spark.sql.Encoders.product[Item]
}

case class Item(
                 cid: String,
                 sha1: String
               )
*/

object HbaseDemo {
  def main(args: Array[String]): Unit = {
    val spark = SparkSession
      .builder()
      .appName("testSparkHbase")
      .config("spark.sql.hive.convertMetastoreOrc", true)
      .config("spark.sql.orc.enabled", true)
      .config("spark.speculation",true)
      .enableHiveSupport()
      .getOrCreate()

    import spark.implicits._
//    val df1 =  spark.sql("select concat_ws('_',clusterid,primarykey) as cid, primarykey as sha1 from xlab_data_mining.lzh_fingerprint_extract_0419 limit 10")
    val df1 = spark.sql("select cid,sha1 from tmp.test_0428_2")
    def catalog = s"""{
                     |"table":{"namespace":"default", "name":"Contacts"},
                     |"rowkey":"key",
                     |"columns":{
                     |"cid":{"cf":"rowkey", "col":"key", "type":"string"},
                     |"sha1":{"cf":"info", "col":"sha1", "type":"string"}
                     |}
                     |}""".stripMargin

    val df2 = Seq(("c1","pk1"),("c2","pk2")).toDF("cid","sha1")

    if (args(0).equals("a")){
      df2.write.options(Map(HBaseTableCatalog.tableCatalog -> catalog, HBaseTableCatalog.newTable -> "5")).format("org.apache.spark.sql.execution.datasources.hbase").save()
    }else{
      df1.write.options(Map(HBaseTableCatalog.tableCatalog -> catalog, HBaseTableCatalog.newTable -> "5")).format("org.apache.spark.sql.execution.datasources.hbase").save()
    }

  }
}



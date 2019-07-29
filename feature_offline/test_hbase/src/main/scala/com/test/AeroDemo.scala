package com.test

import com.aerospike.client.{AerospikeClient, Bin, Key}
import org.apache.spark.sql.SparkSession

object AeroDemo {
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

    df1.foreachPartition(partition => {
      val client = new AerospikeClient("szth-2sys6-hsc-f32-12.szth.baidu.com",3000)
      partition.foreach(line=>{
        val key = new Key("test","test_spark_0428",line.getAs[String]("cid"))
        val bin1 = new Bin("",line.getAs[String]("sha1"))
        // Write a record
        client.put(null, key, bin1)
      })
    })

  }
}

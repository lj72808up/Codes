package com.test
import org.apache.spark.sql.types.{DataType, IntegerType, Metadata, StringType, StructField, StructType}
import org.apache.spark.sql.{Row, SaveMode, SparkSession}
import play.api.libs.json.{Format, Json}


object TestRedis {
  case class Person(name: String, age: Int)

  def main(args: Array[String]): Unit = {

    val spark = SparkSession.builder()
      .appName("SparkReadRedis")
      .master("local[*]")
      .getOrCreate()


  }
}
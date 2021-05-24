package utils

import adtl.platform.Dataproto.HiveField
import org.apache.spark.sql.{Row, SparkSession}

import scala.collection.mutable.Map

case class JsonUtil(spark: SparkSession) {

  case class Field(name: String, ftype: String, comment: String)

  //  private implicit val fieldWrites = new Writes[Field] {
  //    def writes(field: Field) = Json.obj(
  //      "name" -> field.name,
  //      "type" -> field.ftype,
  //      "comment" -> field.comment)
  //  }

  // 返回databases
  def getDatabaseNames: List[String] = {
    spark.sql("show databases").collect().map(_.getAs[String](0)).toList
  }

  // 返回指定database的tables
  def getTableNames(dbname: String): List[String] = {
    spark.sql(s"use $dbname")
    spark.sql("show tables").collect().map(_.getAs[String]("tableName")).toList
  }

  // 返回指定database.table的fields
  def getFieldCases(dbname: String, tbname: String): List[HiveField] = {
    try {
      spark.sql(s"desc $dbname.$tbname").collect.map {
        field =>
          val field_name = field.getAs[String](0)
          val ftype = field.getAs[String](1)
          val comment = if (field.getAs[String](2)==null) "" else field.getAs[String](2)
          HiveField(field_name, ftype, comment)
      }.toList
    } catch {
      case ex: Exception =>
        ex.printStackTrace()
        List[HiveField]()
    }
  }

  // 返回指定database.table的扩展信息()
  def getExpendInfo(dbname: String, tbname: String) : Map[String,String] = {
    val info = Map[String, String]()
    try {
      val res = spark.sql(s"DESCRIBE FORMATTED $dbname.$tbname").collect
      res.foreach(item => {
        val title = item.getAs[String](0).toLowerCase
        if (title.equals("owner")) { // 创建者
          info("owner") = item.getAs[String](1)
        } else if (title.equals("storage properties")) { // 压缩格式
          info("storageProperties") = item.getAs[String](1)
//            .split(",")(0)
//            .replace("[", "")
//            .replace("]", "")
        } else if (title.equals("partition provider")) { // 分区字段
          info("partitionProvider") = item.getAs[String](1)
        } else if (title.equals("location")) { // 存储格式
          info("location") = item.getAs[String](1)
        } else if (title.equals("created time")) { // 创建时间
          info("createdTime") = item.getAs[String](1)
        }
      })

    } catch {
      case ex: Exception => ex.printStackTrace()
    }
    info
  }

  //  def getHiveJson: String = {
  //
  //    Json.toJson(getDatabaseNames).toString
  //  }
  //
  //  def getDbJson(dbname: String): String = {
  //
  //    Json.toJson(getTableNames(dbname)).toString
  //  }
  //
  //  def getTbJson(dbname: String, tbname: String): String = {
  //    Json.toJson(getFieldCases(dbname, tbname)).toString
  //  }
}

package utils

import java.sql.{Connection, DriverManager, ResultSet, Statement}

import scala.util.control.Breaks._
import adtl.platform.Dataproto.{HiveField, HiveTable, HiveTableSet}
import org.apache.spark.sql.SparkSession
import play.api.libs.json.{JsObject, Json, Writes}

import scala.collection.mutable.Map

class HiveMetaLib(spark: SparkSession) {

  var connection: Connection = _
  var stmt: Statement = _

  def conn(): Unit = {
    Class.forName("org.apache.hive.jdbc.HiveDriver")
    connection = DriverManager.getConnection("jdbc:hive2://rsync.master09.saturn.hadoop.js.ted:10005/default;principal=hive/hive@SATURN.HADOOP.COM")
    stmt = connection.createStatement()
  }

  def execSelect(sql: String): ResultSet = {
    stmt.execute(s"Set mapreduce.job.queuename=hive")
    //    stmt.execute(s"Set mapreduce.job.name=${queryInfo.queryId}")
    val resultSet = stmt.executeQuery(sql)
    resultSet
  }

  def close(): Unit = {
    stmt.close()
    connection.close()
  }
}


class HiveMetaDataHelper(spark: SparkSession) {
  private val hive = new HiveMetaLib(spark)
  hive.conn()

  def getTable(dbName: String, tName: String): HiveTable = {
    hive.conn()
    val integratedFields = this.getFieldCases(dbName, tName) // 完整字段
    val fields = integratedFields.map(f => f.fName) // 字段名
    val info = this.getExpendInfo(dbName, tName)
    hive.close()

    HiveTable().withNameTable(tName)
      .withFileds(fields)
      .withExpendInfo(info.toMap)
      .withIntegratedField(integratedFields)
  }

  def getTables(tables: List[String]): HiveTableSet = {
    hive.conn()
    val hiveTables = tables.map(t => {
      var dbName = ""
      var tName = ""

      val expandName = t.split("\\.")
      if (expandName.length == 1) {
        dbName = "default"
        tName = expandName(0)
      } else {
        dbName = expandName(0)
        tName = expandName(1)
      }

      val integratedFields = this.getFieldCases(dbName, tName) // 完整字段
      val fields = integratedFields.map(f => f.fName) // 字段名
      val info = this.getExpendInfo(dbName, tName)

      HiveTable().withNameTable(t)
        .withFileds(fields)
        .withExpendInfo(info.toMap)
        .withIntegratedField(integratedFields)
    })
    hive.close()
    HiveTableSet(hiveTables)
  }


  // 返回指定database.table的fields
  def getFieldCases(dbname: String, tbname: String): List[HiveField] = {
    try {
      val rs = hive.execSelect(s"desc $dbname.$tbname")
      var fields = List[HiveField]()
      breakable {
        while (rs.next()) {
          val fieldName = rs.getString(1)
          if (fieldName.startsWith("#")) {
            break
          }
          if (!(fieldName == null || fieldName.equalsIgnoreCase(""))) {
            val fieldType = rs.getString(2)
            val comment = if (rs.getString(3) == null) "" else rs.getString(3)
            fields = fields :+ HiveField(fieldName, fieldType, comment)
          }
        }
      }
      fields
    } catch {
      case ex: Exception =>
        ex.printStackTrace()
        List[HiveField]()
    }
  }

  // 返回指定database.table的扩展信息()
  def getExpendInfo(dbname: String, tbname: String): Map[String, String] = {
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
        } else if (title.equals("location")) { // 存储格式
          info("location") = item.getAs[String](1)
        } else if (title.equals("created time")) { // 创建时间
          info("createdTime") = item.getAs[String](1)
        }
      })

      // 分区字段
      val partitionInfo = spark.sql(s"SHOW PARTITIONS $dbname.$tbname").take(1)
      var partitions = List[String]()
      if (partitionInfo.length > 0) { // 存在分区信息
        val len = partitionInfo.length
        for (i <- 0 until len) {
          val info = partitionInfo(i).getAs[String](0)
          val cols = info.split("[,/]").map(x => x.split("=")(0).trim).toList
          partitions = partitions ++ cols
        }
      }
      info("partitions") = partitions.mkString(",")
    } catch {
      case ex: Exception => ex.printStackTrace()
    }
    info
  }
}


object HiveMetaDataHelper {
  implicit val HiveFieldWriter: Writes[HiveField] = new Writes[HiveField] {
    def writes(f: HiveField): JsObject = Json.obj(
      "fName" -> f.fName,
      "fType" -> f.fType,
      "fComment" -> f.fComment
    )
  }

  implicit val HiveTableWriter: Writes[HiveTable] = new Writes[HiveTable] {
    def writes(t: HiveTable): JsObject = Json.obj(
      "tName" -> t.nameTable,
      "fields" -> t.integratedField
    )
  }
}
package utils.dag

import airflow.{DagTask, HiveStorage}
import org.apache.spark.sql.{DataFrame, SparkSession}
import org.apache.spark.sql.types.{ArrayType, BooleanType, DataType, DoubleType, FloatType, IntegerType, LongType, StringType}
import sogou.adtl.clickhouseUtil
import utils.{PrintLogger, SparkUtil}

object CkParser extends PrintLogger {
  val spark: SparkSession = SparkUtil.getSparkSession()

  def createTableFromSql(sql: String, tableIdentifier: String, partitions: Array[String], indexes: Array[String]): Unit = {
    val df = spark.sql(sql)
    var db = "datacenter"
    var tb = tableIdentifier
    val splits = tableIdentifier.split("\\.")
    if (splits.length > 1) {
      db = splits(0)
      tb = splits(1)
    }
    Class.forName("ru.yandex.clickhouse.ClickHouseDriver").newInstance()
    val schema = df.schema.map {
      struct =>
        val name = struct.name
        val typeName = tranClickhouseType(struct.dataType)
        s"$name $typeName"
    }.mkString("(", ", ", ")")
    info(schema)
    clickhouseUtil.createTableCluster(db, tb, schema, partitions.mkString("(", ",", ")"), indexes.mkString("(", ",", ")"))
  }

  def createTableFromTasks(spark: SparkSession, tasks: List[DagTask]): Unit = {
    for (i <- tasks.indices) {
      if (tasks(i).newOrOldOutTable == "newTable") {
        createTableFromSql(tasks(i).sql, tasks(i).tableName, tasks(i).clickhouse.partitions.getOrElse(Array[String]()),
          tasks(i).clickhouse.indices.getOrElse(Array[String]()))
      }
    }
  }

  private def tranClickhouseType(dataType: DataType): String = {
    dataType match {
      case StringType =>
        "String"
      case FloatType =>
        "Float32"
      case DoubleType =>
        "Float64"
      case LongType =>
        "Int64"
      case IntegerType =>
        "Int32"
      case BooleanType =>
        "Boolean"
      case ArrayType(StringType, true) =>
        "Array(String)"
      case ArrayType(IntegerType, true) =>
        "Array(UInt32)"
      case ArrayType(LongType, true) =>
        "Array(UInt64)"
      case ArrayType(IntegerType, false) =>
        "Array(UInt32)"
      case ArrayType(LongType, false) =>
        "Array(UInt64)"
      case _ =>
        "String"
    }
  }
}

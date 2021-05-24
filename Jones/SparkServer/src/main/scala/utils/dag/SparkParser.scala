package utils.dag

import adtl.platform.Dataproto.HiveTableSet
import airflow.{DagTask, HiveStorage, TaskTypeInfo}
import org.apache.spark.sql.SparkSession
import utils.{HiveMetaDataHelper, LinageHelper, Node, PrintLogger}

object SparkParser extends PrintLogger {

  def createTableFromSql(spark: SparkSession,
                         task: DagTask,
                         partition: String): Unit = {
    val selectSql = task.sql
    val tableName = task.tableName
    val staticFields = task.sqlFields.getOrElse(Map())
    val storage = task.storageEngine
    val comment = task.tableDesc
    val taskTypeInfo = task.taskTypeInfo
    val inputType = task.inputType
    var fieldList = List[String]()
    if ("hdfs".equals(inputType)) {
      // 输入是hdfs时, 字段来自taskInfo  -- add 20201103
      if (task.taskTypeInfo.hasExtend) {
        // 有extend的话加到开头
        fieldList = fieldList :+ s"`extend` MAP<STRING,STRING>" // 加上extend字段
      }
      val schemas = taskTypeInfo.schema.flatMap(m => m.get("name"))
      schemas.foreach(
        f => {
          val fType = "STRING"
          val fName = f
          val fComment = staticFields.getOrElse(fName, "")
          // 判断列名不等于分区名
          if (!fName.equals(partition)) {
            fieldList = fieldList :+ s"`$fName` $fType comment '$fComment'"
          }
        }
      )
    } else {
      val df = spark.sql(selectSql)
      val fields = df.schema.fields
      fields.foreach(
        f => {
          val fType = f.dataType.sql
          val fName = f.name
          val fComment = staticFields.getOrElse(fName, "")
          // 判断列名不等于分区名
          if (!fName.equals(partition)) {
            fieldList = fieldList :+ s"`$fName` $fType comment '$fComment'"
          }
        }
      )
    }


    var handleTableName = tableName
    if (!tableName.contains(".")) {
      handleTableName = s"datacenter.$tableName"
    }
    var createSql =
      s"""CREATE TABLE IF NOT EXISTS $handleTableName (
         | ${fieldList.mkString(",\n ")}
         |)
         |COMMENT '$comment'
         |PARTITIONED BY ($partition string)""".stripMargin

    ///// hard code 硬编码, 加上 时间和日期分区
    if (partition.equals("hour")) {
      createSql =
        s"""CREATE TABLE IF NOT EXISTS $handleTableName (
           | ${fieldList.mkString(",\n ")}
           |)
           |COMMENT '$comment'
           |PARTITIONED BY (dt string,hour string)""".stripMargin
    }
    /////

    var storageSql = ""
    HiveStorage.withName(storage) match {
      case HiveStorage.Tsv =>
        storageSql =
          s"""ROW FORMAT DELIMITED
             |FIELDS TERMINATED BY '\t'
             |LOCATION '/user/adtd_platform/search/odin/airflow/${handleTableName}'
             |""".stripMargin
      case HiveStorage.Orc_Snappy =>
        storageSql =
          s"""STORED AS ORC
             |LOCATION '/user/adtd_platform/search/odin/airflow/${handleTableName}'
             |TBLPROPERTIES ('orc.compress'='SNAPPY')
             |""".stripMargin
    }
    createSql = s"${createSql}\n${storageSql}"
    info(createSql)
    spark.sql(createSql).foreach(_=>{})
  }

  def createTableFromTasks(spark: SparkSession, tasks: List[DagTask]): Unit = {
    for (i <- tasks.indices) {
      val partition = tasks(i).taskType match {
        case "day" => "dt"
        case "hour" => "hour"
        case _ => throw new Exception("任务类型错误, 只能是天级任务(day), 或小时级任务(hour)")
      }
      if (tasks(i).newOrOldOutTable == "newTable") {
        createTableFromSql(spark, tasks(i), partition)
      }
    }
  }

  def getTableInfo(spark: SparkSession, tasks: List[DagTask]): HiveTableSet = {
    val tables = tasks.filter(t => (t.outTableType == "hive") || (t.outTableType == "spark")).map(x => {
      var tableName = x.tableName
      if (!tableName.contains(".")) {
        tableName = s"datacenter.$tableName"
      }
      tableName
    })
    new HiveMetaDataHelper(spark).getTables(tables)
  }

  def getUpStreamTables(sql: String, rolesStr: String): Array[Node] = {
    val helper = new LinageHelper()
    try {
      helper.parseLinage(sql, "datacenter", rolesStr)
    } catch {
      case ex: Exception => {
        ex.printStackTrace()
        Array()
      }
    }
  }
}

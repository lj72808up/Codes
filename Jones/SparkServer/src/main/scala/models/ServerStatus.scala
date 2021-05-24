package models

import java.util.concurrent.atomic.AtomicBoolean

import adtl.platform.Dataproto.{HiveDatabase, HiveMeta, HiveTable}
import org.apache.spark.sql.SparkSession
import utils.JsonUtil
import adtl.platform.Dataproto._
import models.config.MainConf

object ServerStatus {
  val isRunning: AtomicBoolean = new AtomicBoolean(false)
}

object HiveMetaInfo {

  var protoHiveMeta = HiveMeta()
  var sparkUDFs: Array[String] = _

  /*def genHiveMeta_old(spark:SparkSession):Unit = {
    val jsonUtil = JsonUtil(spark)
    val hiveDatabases = jsonUtil.getDatabaseNames.map{
      database=>
        val hiveTables = jsonUtil.getTableNames(database).map{
          table=>
            val schemas = jsonUtil.getFieldCases(database,table).map(_.name).filterNot(_.contains("#"))
            HiveTable().withNameTable(table).withFileds(schemas)
        }.filterNot(_.fileds.length == 0)
        HiveDatabase().withNameDatabase(database).withTalbes(hiveTables)
    }.filterNot(_.talbes.length == 0)
    protoHiveMeta = HiveMeta().withDatabases(hiveDatabases)
  }
*/
  def genHiveMeta(spark: SparkSession): Unit = {
    val jsonUtil = JsonUtil(spark)
    val hiveDatabases = jsonUtil.getDatabaseNames.map {
      database =>
        val hiveTables = jsonUtil.getTableNames(database).map {
          table =>
            val integratedFields = jsonUtil.getFieldCases(database, table).filterNot(_.fName.contains("#"))
            val fields = integratedFields.map(_.fName)
            val info = jsonUtil.getExpendInfo(database, table)

            HiveTable().withNameTable(table)
              .withFileds(fields)
              .withExpendInfo(info.toMap)
              .withIntegratedField(integratedFields)
        }.filterNot(_.fileds.isEmpty)
        HiveDatabase().withNameDatabase(database).withTalbes(hiveTables)
    }.filterNot(_.talbes.isEmpty)
    protoHiveMeta = HiveMeta().withDatabases(hiveDatabases)
  }

  // 增加udf函数统计
  def genSparkUDFMeta(spark: SparkSession): Unit = {
    val funcs = spark.sql("show functions").collect
    val funcNames = funcs.map(r => r.getAs[String](0))
    this.sparkUDFs = funcNames
  }
}


object ServerInfo {
  var Host: String = ""
  var Port: Int = 0
  var Respath: String = ""
  var Sparkname: String = ""
  var FrontOptions: FrontOptions = _

}

object MQInfo {
  var Host: String = MainConf("basicConf").common.host
  var userName: String = MainConf("basicConf").common.userName
  var password: String = MainConf("basicConf").common.password
  //var jobQueueName = "ASTAR"
  var jobQueueName: String = MainConf("basicConf").common.mq
  var dbLogQueueName = ""
}
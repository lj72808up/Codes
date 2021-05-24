package utils

import adtl.platform.Dataproto.{QueryInfo, SmallFlow}
import models.config.{MainConf, MySqlConf}
import org.apache.spark.sql.SparkSession

object SmallFlowTaskLib {

  val mysqlConf: MySqlConf = MainConf("basicConf").mysql
  val host: String = mysqlConf.host
  val port: String = mysqlConf.port
  val username: String = mysqlConf.user
  val passwd: String = mysqlConf.passwd
  val dbname: String = mysqlConf.db

  val conn_str = s"jdbc:mysql://$host:$port/$dbname?user=$username&password=$passwd"
  val tbName = "small_flow_tasks"

  def saveTask(spark: SparkSession, smallFlow: SmallFlow): Unit = {

    val conn = try {
      MysqlUtil.getConnection(conn_str)
    }
    catch {
      case ex: Exception => throw ex
    }
    try {
      MysqlUtil.executeInsert(conn, s"INSERT INTO $tbName(task_id, task_name, start_date, end_date, session, is_enabled, full_sql, same_query_sql)" +
        s" VALUES (?,?,?,?,?,?,?,?)",smallFlow.taskId, smallFlow.taskName, smallFlow.startDate, smallFlow.endDate, smallFlow.session, "true", smallFlow.sqls(0), smallFlow.sqls(1))
    } catch {
      case ex: Exception => throw ex
    } finally {
      MysqlUtil.closeConnection(conn)
    }

  }

}

package models

import adtl.platform.Dataproto.{QueryInfo, QueryStatus}
import utils.MysqlUtil

object QueryInfoModel {

  import models.config._

  val tbName: String = "query_submit_log"

  //  val mysqlConf: MySqlConf = MainConf("basicConf").mysql
  val mysqlConf: MySqlConf = MainConf("basicConf").mysql
  val host: String = mysqlConf.host
  val port: String = mysqlConf.port
  val username: String = mysqlConf.user
  val passwd: String = mysqlConf.passwd
  val dbname: String = mysqlConf.db

  // useUnicode=true&characterEncoding=UTF-8
  val conn_str = s"jdbc:mysql://$host:$port/$dbname?user=$username&password=$passwd&useSSL=false"

  object implicits {
    // 定义int到QueryStatus的隐式转换
    implicit def int2QueryStatus(x: Int): QueryStatus = {
      val status = x match {
        case 0 => QueryStatus.Waiting
        case 1 => QueryStatus.Running
        case 2 => QueryStatus.Killing
        case 3 => QueryStatus.Finished
        case 4 => QueryStatus.Killed
        case 5 => QueryStatus.SparkReRunning
        case 6 => QueryStatus.SparkReRunFinished
        case 10 => QueryStatus.SyncQueryWaiting
        case 11 => QueryStatus.SyncQueryRunning
        case 12 => QueryStatus.SyncQueryFinished
        case 13 => QueryStatus.SyncQueryFailed
        case 20 => QueryStatus.CustomQueryWaiting
        case 21 => QueryStatus.CustomQueryRunning
        case 22 => QueryStatus.CustomQueryFinished
        case 23 => QueryStatus.CustomQueryFailed
      }
      status
    }
  }

  // 定义queryStatus到int的显示转换，因为这个int值传给了MysqlUtil类，所以不方便隐式转换
  private def queryStatus2Int(x: QueryStatus): Int = {
    val status = x match {
      case QueryStatus.Waiting => 0
      case QueryStatus.Running => 1
      case QueryStatus.Killing => 2
      case QueryStatus.Finished => 3
      case QueryStatus.Killed => 4
      case QueryStatus.SparkReRunning => 5
      case QueryStatus.SparkReRunFinished => 6
      case QueryStatus.SyncQueryWaiting => 10
      case QueryStatus.SyncQueryRunning => 11
      case QueryStatus.SyncQueryFinished => 12
      case QueryStatus.SyncQueryFailed => 13
      case QueryStatus.CustomQueryWaiting => 20
      case QueryStatus.CustomQueryRunning => 21
      case QueryStatus.CustomQueryFinished => 22
      case QueryStatus.CustomQueryFailed => 23
      case _ => throw new Exception("QueryStatus unknown")
    }
    status
  }

  // 获取QueryInfo的方法，通过传入QueryInfo的Id返回数据库中的QueryInfo对象
  def getQueryInfo(info: QueryInfo): QueryInfo = {

    val queryId = info.queryId
    val conn = try {
      MysqlUtil.getConnection(conn_str)
    }
    catch {
      case ex: Exception => throw ex
    }
    var result: QueryInfo = QueryInfo()
    try {
      val rs = MysqlUtil.executeQuery(conn, s"select * from $tbName where query_id = ?", queryId)
      rs.last
      import implicits._
      result = QueryInfo().withQueryId(rs.getString("query_id")).withSession(rs.getString("session")).withSql(Seq(rs.getString("sql_str"))).withStatus(rs.getInt("status")).withLastUpdate(rs.getLong("last_update")).withDataPath(rs.getString("data_path"))
    }
    catch {
      case ex: Exception => throw ex
    } finally {
      MysqlUtil.closeConnection(conn)
    }
    result
  }

  def updateQueryInfoCreateTime(info: QueryInfo): Int = {
    val conn = try {
      MysqlUtil.getConnection(conn_str)
    }
    catch {
      case ex: Exception => throw ex
    }
    try {
      val info_status: Int = queryStatus2Int(info.status)
      MysqlUtil.executeUpdate(conn, s"UPDATE $tbName SET status=?, create_time = ?, last_update=?, data_path=?, exception_msg=? WHERE query_id=?", info_status, info.createTime, info.lastUpdate, info.dataPath, info.exceptionMsg, info.queryId)
    } catch {
      case ex: Exception => throw ex
    } finally {
      MysqlUtil.closeConnection(conn)
    }
  }

  def updateQueryInfoStartTime(info: QueryInfo): Int = {
    val conn = try {
      MysqlUtil.getConnection(conn_str)
    }
    catch {
      case ex: Exception =>
        ex.printStackTrace()
        throw ex
    }
    try {
      val info_status = queryStatus2Int(info.status)
      MysqlUtil.executeUpdate(conn, s"UPDATE $tbName SET status=?, start_time = ?, last_update=?, data_path=?, exception_msg=? WHERE query_id=?", info_status, info.startTime, info.lastUpdate, info.dataPath, info.exceptionMsg, info.queryId)
    } catch {
      case ex: Exception => throw ex
    } finally {
      MysqlUtil.closeConnection(conn)
    }
  }

  // 用传入的QueryInfo对象的值更新数据库中的对应QueryId的列的值
  def updateQueryInfo(info: QueryInfo): Int = {
    val conn = try {

      MysqlUtil.getConnection(conn_str)
    }
    catch {
      case ex: Exception => throw ex
    }
    try {
      val info_status = queryStatus2Int(info.status)
      MysqlUtil.executeUpdate(conn, s"UPDATE $tbName SET status=?, last_update=?, data_path=?, exception_msg=? WHERE query_id=?", info_status, info.lastUpdate, info.dataPath, info.exceptionMsg, info.queryId)
    } catch {
      case ex: Exception => throw ex
    } finally {
      MysqlUtil.closeConnection(conn)
    }
  }
}

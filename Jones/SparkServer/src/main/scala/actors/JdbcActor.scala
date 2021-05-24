package actors

import actors.RabbitMQActor.QueryPublish
import actors.TaskProxyActor.{SQLSubmitCK, SQLSubmitJdbc, SQLSubmitKylin}
import adtl.platform.Dataproto.{QueryEngine, QueryStatus}
import akka.actor.{Actor, ActorRef}
import models.ServerInfo
import org.apache.spark.sql.SparkSession
import utils.DateTimeUtil.getTsNow
import utils.{JdbcClacLib, PrintLogger}

class JdbcActor(sparkSession: SparkSession, dbLog: ActorRef) extends Actor with PrintLogger{
  override def receive: Receive = {
// mysql
    case jdbcSubmit: SQLSubmitJdbc => { // sender() is TaskProxyActor
      sender() ! QueryPublish(jdbcSubmit.query.copy(status = QueryStatus.fromName("Running").get, lastUpdate = getTsNow(), startTime = getTsNow()))
      val finishStatus = "Finished"
      try {
        val respath = JdbcClacLib.execSqlAndSaveResult(sparkSession, "mysql",jdbcSubmit.query, ServerInfo.Respath)
        sender() ! QueryPublish(jdbcSubmit.query.copy(status = QueryStatus.fromName(finishStatus).get, dataPath = respath, lastUpdate = getTsNow()))
      } catch {
        case e: Exception => {
          e.printStackTrace()
          val errMsg = if (e.getMessage==null) "null point exception" else e.getMessage
          sender() ! QueryPublish(jdbcSubmit.query.copy(
            status = QueryStatus.fromName("Killed").get,
            engine = QueryEngine.fromName("jdbc").get,
            lastUpdate = getTsNow(),
            exceptionMsg = errMsg))
        }
      }
    }
// clickhouse
    case ckSubmit: SQLSubmitCK => {
      sender() ! QueryPublish(ckSubmit.query.copy(status = QueryStatus.fromName("Running").get, lastUpdate = getTsNow(), startTime = getTsNow()))
      val finishStatus = "Finished"
      try {
        val respath = JdbcClacLib.execSqlAndSaveResult(sparkSession, "clickhouse",ckSubmit.query, ServerInfo.Respath)
        sender() ! QueryPublish(ckSubmit.query.copy(status = QueryStatus.fromName(finishStatus).get, dataPath = respath, lastUpdate = getTsNow()))
      } catch {
        case e: Exception => {
          e.printStackTrace()
          val errMsg = if (e.getMessage==null) "null point exception" else e.getMessage
          sender() ! QueryPublish(ckSubmit.query.copy(
            status = QueryStatus.fromName("Killed").get,
            engine = QueryEngine.fromName("jdbc_clickhouse").get,
            lastUpdate = getTsNow(),
            exceptionMsg = errMsg))
        }
      }
    }

    case kySubmit: SQLSubmitKylin => {
      sender() ! QueryPublish(kySubmit.query.copy(status = QueryStatus.fromName("Running").get, lastUpdate = getTsNow(), startTime = getTsNow()))
      val finishStatus = "Finished"
      try {
        val respath = JdbcClacLib.execSqlAndSaveResult(sparkSession, "kylin",kySubmit.query, ServerInfo.Respath)
        sender() ! QueryPublish(kySubmit.query.copy(status = QueryStatus.fromName(finishStatus).get, dataPath = respath, lastUpdate = getTsNow()))
      } catch {
        case e: Exception => {
          e.printStackTrace()
          val errMsg = if (e.getMessage==null) "null point exception" else e.getMessage
          /*sender() ! QueryPublish(kySubmit.query.copy(
            status = QueryStatus.fromName("Killed").get,
            engine = QueryEngine.fromName("kylin").get,
            lastUpdate = getTsNow(),
            exceptionMsg = errMsg))*/
          val reRunQueryInfo = kySubmit.query.copy(
            status = QueryStatus.fromName("Waiting").get,
            engine = QueryEngine.fromName("spark").get,
            lastUpdate = getTsNow(),
            exceptionMsg = errMsg
          )
          sender() ! QueryPublish(reRunQueryInfo)
        }
      }
    }

    case _ => info("unknown task (not jdbc)")
  }
}

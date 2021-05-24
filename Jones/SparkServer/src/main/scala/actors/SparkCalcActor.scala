package actors

import adtl.platform.Dataproto.{QueryEngine, QueryStatus}
import akka.actor.{Actor, ActorLogging, ActorRef}
import models.ServerInfo
import org.apache.spark.sql.SparkSession
import utils.{HDFSUtil, JdbcClacLib, SparkCalcLib}


class SparkCalcActor(sparkSession: SparkSession, dbLog: ActorRef) extends Actor with ActorLogging {

  import RabbitMQActor._
  import TaskProxyActor._
  import utils.DateTimeUtil.getTsNow

  override def preStart(): Unit = {
    this.updateTaskExecTime()
  }

  private def updateTaskExecTime(): Unit = {
    val path = self.path.toString
    TaskProxyActor.JobActorTimeMapping.put(path, System.currentTimeMillis())
  }

  override def receive: Receive = {

    case submit: SQLSubmit => {
      // 更新任务的真正开始时间
      // 更新后的逻辑里，dbLogActor就是消息的终点，不需要考虑消息如何返回的问题
      val startStatus = if (submit.query.status == QueryStatus.SparkReRunning) {
        "SparkReRunning"
      } else if (submit.query.status == QueryStatus.CustomQueryWaiting) {
        "CustomQueryRunning"
      } else {
        "Running"
      }
      sender() ! QueryPublish(submit.query.copy(status = QueryStatus.fromName(startStatus).get, lastUpdate = getTsNow, startTime = getTsNow))
      //dbLog ! QueryInfoUpdateStart(submit.query.withStartTime(getTsNow))
      //TestUtil.printQuery(submit.query)

      try {
        this.updateTaskExecTime()
        val respath = SparkCalcLib.RunSQLSaveLocalCSV(sparkSession, submit.query, ServerInfo.Respath)
        HDFSUtil.delHdfsFilePath(sparkSession, respath)
        val finishStatus = if (submit.query.status == QueryStatus.SparkReRunning) "SparkReRunFinished"
        else if (submit.query.status == QueryStatus.CustomQueryRunning) "CustomQueryFinished"
        else "Finished"
        sender() ! QueryPublish(submit.query.copy(status = QueryStatus.fromName(finishStatus).get, dataPath = respath, lastUpdate = getTsNow))
      } catch {
        case e: Exception =>
          e.printStackTrace()
          val errMsg = if (e.getMessage == null) "null point exception" else e.getMessage
          sender() ! QueryPublish(submit.query.copy(status = QueryStatus.fromValue(4), lastUpdate = getTsNow, exceptionMsg = errMsg))
      }
    }
    case cancel: SQLCancel => {
      SparkCalcLib.CancelSQLByGroupId(sparkSession, cancel.query.queryId)
      sender() ! QueryPublish(cancel.query.copy(status = QueryStatus.fromName("Killed").get, lastUpdate = getTsNow, dataPath = "任务取消"))
    }
    case hiveSubmit: SQLSubmitHive => {
      sender() ! QueryPublish(hiveSubmit.query.copy(status = QueryStatus.fromName("Running").get, lastUpdate = getTsNow, startTime = getTsNow))
      //dbLog ! QueryInfoUpdateStart(hiveSubmit.query.withStartTime(getTsNow))
      try {
        this.updateTaskExecTime()
        val respath = SparkCalcLib.RunSQLSaveLocalCSV(sparkSession, hiveSubmit.query, ServerInfo.Respath)
        val finishStatus = if (hiveSubmit.query.status == QueryStatus.SparkReRunning) "SparkReRunFinished" else "Finished"
        sender() ! QueryPublish(hiveSubmit.query.copy(status = QueryStatus.fromName(finishStatus).get, dataPath = respath, lastUpdate = getTsNow))
      } catch {
        case e: Exception =>
          //sender() ! QueryPublish(hiveSubmit.query.copy(status = QueryStatus.fromName("Killed").get, lastUpdate = getTsNow, exceptionMsg = e.getMessage))
          e.printStackTrace()
          val errMsg = if (e.getMessage == null) "null point exception" else e.getMessage
          val reRunQueryInfo = hiveSubmit.query.copy(
            status = QueryStatus.fromName("Waiting").get,
            engine = QueryEngine.fromName("spark").get,
            lastUpdate = getTsNow(),
            exceptionMsg = errMsg
          )
          sender() ! QueryPublish(reRunQueryInfo)
      }
    }

    case _ => {
      log.info("unknown type task")
    }
  }
}


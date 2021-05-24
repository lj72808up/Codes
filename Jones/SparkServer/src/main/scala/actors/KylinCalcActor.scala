package actors

import actors.DBlogActor.QueryInfoUpdateStart
import adtl.platform.Dataproto.{QueryEngine, QueryStatus}
import akka.actor.{Actor, ActorLogging, ActorRef}
import models.ServerInfo
import models.config.{KylinConf, MainConf}
import org.apache.spark.sql.SparkSession
import utils.{KylinCalcUtil, SparkCalcLib, TestUtil}


class KylinCalcActor(dbLog: ActorRef) extends Actor with ActorLogging {

  import RabbitMQActor._
  import TaskProxyActor._
  import utils.DateTimeUtil.getTsNow

  override def receive: Receive = {
    case submit: SQLSubmit =>
      sender() ! QueryPublish(submit.query.copy(status = QueryStatus.fromName("Running").get, lastUpdate = getTsNow))
      dbLog ! QueryInfoUpdateStart(submit.query.withStartTime(getTsNow))
      try {
        val kylinConf = MainConf("basicConf").kylin
        val appConf = MainConf("basicConf").app
        val respath = KylinCalcUtil.apply(kylinConf.jdbcHost, kylinConf.user, kylinConf.passwd).RunSQLSaveLocalCSV(submit.query, appConf.resPath)
        sender() ! QueryPublish(submit.query.copy(status = QueryStatus.fromName("Finished").get, dataPath = respath, lastUpdate = getTsNow))
      } catch {
        case e: Exception =>
          // sender() ! QueryPublish(submit.query.copy(status = QueryStatus.fromName("Killed").get, lastUpdate = getTsNow, exceptionMsg = e.getMessage))

          // 添加kylin报错由spark重跑的逻辑
          // 修改状态为“由spark进行重跑中”
          // 修改引擎为spark/hive
          // 更新最后更改时间
          println("kylin ex")
          val reRunQueryInfo = submit.query.copy(
            status = QueryStatus.fromName("SparkReRunning").get,
            engine = QueryEngine.fromName("hive").get,
            lastUpdate = getTsNow()
          )
          sender() ! QueryPublish(reRunQueryInfo)
      }
  }
}


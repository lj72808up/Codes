package actors

import java.util
import java.util.concurrent.ConcurrentHashMap

import adtl.platform.Dataproto.{QueryEngine, QueryInfo, QueryStatus}
import akka.actor.{Actor, ActorLogging, ActorRef, Props}
import org.apache.spark.sql.SparkSession

import collection.JavaConversions._
import RabbitMQActor._
import actors.routers.IdleFirstPool
import akka.routing.RoundRobinPool
import utils.DateTimeUtil.getTsNow

class TaskProxyActor(sparkSession: SparkSession, jobQueueNum: Int, dbLog: ActorRef) extends Actor with ActorLogging {

  import TaskProxyActor._

  var queryStatusMap = new util.HashMap[String, Int]()
  //  更改手动后计算 idx 切换 actor 为使用 router 选举 actor
  //  var sparkCalcList = new util.LinkedList[ActorRef]()
  //  var jdbcCalcList = new util.LinkedList[ActorRef]()
  var sparkRouter: ActorRef = _
  var jdbcRouter: ActorRef = _

  var sparkCancel: ActorRef = _
  var kylinCalc: ActorRef = _

  var queueIndex = 0
  var jdbcIndex = 0

  def idxIncrease(actorType: String) {
    // todo actor的筛选策略: 哪个actor没任务那个actor执行, 那个actor的平均处理时间短, 哪个执行
    actorType match {
      case "jdbc" => {
        jdbcIndex += 1
        if (jdbcIndex == jobQueueNum) jdbcIndex = 0
      }
      case _ => {
        queueIndex += 1
        if (queueIndex == jobQueueNum) queueIndex = 0
      }
    }
  }

  override def preStart(): Unit = {
    log.info("preStart task proxy actor")
    sparkCancel = context.actorOf(Props(classOf[SparkCalcActor], sparkSession, dbLog))
    kylinCalc = context.actorOf(Props(classOf[KylinCalcActor], dbLog))

    this.jdbcRouter = context.actorOf(RoundRobinPool(jobQueueNum).props(Props(classOf[JdbcActor],
      sparkSession, dbLog)), "jdbcRouter")
    this.sparkRouter = context.actorOf(IdleFirstPool(jobQueueNum).props(Props(classOf[SparkCalcActor],
      sparkSession, dbLog)), "sparkRouter")
  }

  override def postStop(): Unit = {
    context.stop(sparkCancel)
    /*    sparkCalcList.foreach {
          actorRef =>
            context.stop(actorRef)
        }*/
    context.stop(kylinCalc)
  }

  override def receive: Receive = {
    case submit: SQLSubmit =>
      submit.query.engine match {
        /*case engine if engine.iskylin =>
          log.info("KYLIN Actor")
          kylinCalc ! submit*/
        case engine if engine.isspark =>
          sparkRouter ! submit
        case engine if engine.ishive =>
          sparkRouter ! submit
        //          log.info("SPARK Actor " + queueIndex)
        //          val ref = sparkCalcList.get(queueIndex)
        //          ref ! submit
        //          idxIncrease("spark")
        //case engine if engine.ishive =>
        //  log.info("Hive Actor ")
        //  sparkRouter ! SQLSubmitHive(submit.query)
        //          log.info("Hive Actor " + queueIndex)
        //          val ref = sparkCalcList.get(queueIndex)
        //          ref ! SQLSubmitHive(submit.query)
        //          idxIncrease("spark")
        case engine if engine.isjdbc =>
          log.info("JDBC Actor (mysql)")
          jdbcRouter ! SQLSubmitJdbc(submit.query)
        //          log.info("JDBC Actor (mysql)" + jdbcIndex)
        //          val ref = jdbcCalcList.get(jdbcIndex)
        //          ref ! SQLSubmitJdbc(submit.query)
        //          idxIncrease("jdbc")
        case engine if engine.isjdbcspark =>
          log.info("JDBC SPARK Actor ")
          sparkRouter ! submit

        //          log.info("JDBC SPARK Actor " + queueIndex)
        //          val ref = sparkCalcList.get(queueIndex)
        //          ref ! submit
        //          idxIncrease("spark")
        // 使用jdbcActor跑clickhouse
        case engine if engine.isjdbcclickhouse =>
          log.info("JDBC Actor (clickhouse) ")
          jdbcRouter ! SQLSubmitCK(submit.query)
        //          log.info("JDBC Actor (clickhouse) " + jdbcIndex)
        //          val ref = jdbcCalcList.get(jdbcIndex)
        //          ref ! SQLSubmitCK(submit.query)
        //          idxIncrease("jdbc")
        case engine if engine.iskylin =>
          log.info("JDBC Actor (kylin)")
          jdbcRouter ! SQLSubmitKylin(submit.query)
        case x@_ =>
          val errMsg = s"unknown engine type: $x"
          log.error(errMsg)
          context.parent ! QueryPublish(submit.query.copy(
            status = QueryStatus.fromName("Killed").get,
            lastUpdate = getTsNow(),
            exceptionMsg = errMsg))
      }
    case cancel: SQLCancel =>
      sparkCancel ! cancel
    case publish: QueryPublish =>
      context.parent ! publish
  }
}

object TaskProxyActor {

  val JobActorTimeMapping = new ConcurrentHashMap[String, Long]()

  sealed trait TaskCalcAction

  case class SQLSubmit(query: QueryInfo) extends TaskCalcAction

  case class SQLSubmitHive(query: QueryInfo) extends TaskCalcAction

  case class SQLCancel(query: QueryInfo) extends TaskCalcAction

  case class SQLRunning(query: QueryInfo) extends TaskCalcAction

  case class SQLSubmitJdbc(query: QueryInfo) extends TaskCalcAction

  case class SQLSubmitCK(query: QueryInfo) extends TaskCalcAction

  case class SQLSubmitKylin(query: QueryInfo) extends TaskCalcAction

}

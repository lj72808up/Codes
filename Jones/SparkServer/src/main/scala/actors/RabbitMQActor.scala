package actors

import actors.RabbitMQActor._
import adtl.platform.Dataproto.QueryInfo
import akka.actor.SupervisorStrategy.{Escalate, Restart}
import akka.actor.{Actor, ActorLogging, ActorRef, Kill, OneForOneStrategy, Props}
import akka.routing.{RoundRobinPool, Router}
import com.rabbitmq.client._
import models.MQInfo
import org.apache.spark.sql.SparkSession
import utils.DateTimeUtil.getTsNow

import scala.concurrent.duration.Duration

object RabbitMQActor {

  sealed trait RabbitMQAction

  case object JobConsumer extends RabbitMQAction

  case object DBLogConsumer extends RabbitMQAction

  case class DBLogException(exMsg: String) extends RabbitMQAction

  case object KillRabbit extends RabbitMQAction

  case class QueryPublish(query: QueryInfo) extends RabbitMQAction

}

class RabbitMQActor(spark: SparkSession) extends Actor with ActorLogging {

  import TaskProxyActor._
  import DBlogActor._

  var factory: ConnectionFactory = _
  var conn: Connection = _
  var jobChannel: Channel = _
  var dbLogChannel: Channel = _
  var publishChannel: Channel = _
  var dbLog: ActorRef = _
  var jobProxy: ActorRef = _
  var redisRouter: ActorRef = _

  override val supervisorStrategy =
    // 无限次数重启子actor
    OneForOneStrategy(maxNrOfRetries = -1, withinTimeRange = Duration.Inf) {
      case _: NullPointerException => Restart // Restart同样会触发actor的preStart方法
      case _: Exception => Restart            //Escalate
    }

  def InitMQ(): Unit = {
    this.factory = new ConnectionFactory
    this.factory.setHost(MQInfo.Host)
    factory.setUsername(MQInfo.userName)
    factory.setPassword(MQInfo.password)
    // connection that will recover automatically(首次创建connection失败不会重试, channel级别的异常不会触发重试)
    factory.setAutomaticRecoveryEnabled(true)
    factory.setTopologyRecoveryEnabled(true) // 自动恢复exchange,queue,binding对象
    // attempt recovery every 10 seconds
    factory.setNetworkRecoveryInterval(10 * 1000)
    factory.setRequestedHeartbeat(20)
    this.conn = factory.newConnection()
    this.publishChannel = conn.createChannel()
    this.publishChannel.queueDeclare(MQInfo.jobQueueName, false, false, false, null)
    this.dbLog = context.actorOf(Props[DBlogActor], "dbLog")
    this.jobProxy = context.actorOf(Props(classOf[TaskProxyActor], spark, 10, dbLog), "jobProxy")
    initRedisRouter()
  }

  def initRedisRouter():Unit = {
    val jobQueueNum = 10 // 跑任务的有10个actor, 写Redis的也10个
    this.redisRouter = context.actorOf(RoundRobinPool(jobQueueNum).props(Props(classOf[RedisActor])), "redisRouter")
  }

  def finishMQ(): Unit = {
    jobChannel.close()
    //dbLogChannel.close()
    conn.close()
  }

  override def preStart(): Unit = InitMQ()

  override def postStop(): Unit = finishMQ()

  private def sendMsgToUpdateInfoDb(query: QueryInfo) {}


  override def receive: Receive = {

    case JobConsumer => {
      jobChannel = conn.createChannel()
      jobChannel.queueDeclare(MQInfo.jobQueueName, false, false, false, null)
      val consumer = new DefaultConsumer(jobChannel) {
        override def handleDelivery(consumerTag: String, envelope: Envelope, properties: AMQP.BasicProperties, body: Array[Byte]): Unit = {
          val query: QueryInfo = QueryInfo.parseFrom(body)
          try {
            query.status match {

              case status if status.isWaiting => {
                log.info("status: waiting. submitting query to proxy")

                jobProxy ! SQLSubmit(query)
              }
              case status if status.isCustomQueryWaiting => {
                log.info("status: waiting. submitting query to proxy")
                jobProxy ! SQLSubmit(query)
              }

              case status if status.isCustomQueryRunning => {
                log.info("status: running. updating query_submit_log")
                dbLog ! QueryInfoUpdate(query)
              }

              case status if status.isRunning => {
                log.info("status: running. updating query_submit_log")
                dbLog ! QueryInfoUpdate(query)
              }
              case status if status.isKilling => {
                log.info("status: killing")
                jobProxy ! SQLCancel(query)
                dbLog ! QueryInfoUpdate(query)
              }
              case status if status.isKilled => {
                log.info("status: killed. updating query_submit_log")
                dbLog ! QueryInfoUpdate(query)
              }
              case status if status.isFinished => {
                log.info("status: finished. updating query_submit_log")
                redisRouter ! QueryInfoUpdate(query)   // send a message asynchronously and return immediately
                dbLog ! QueryInfoUpdate(query)
              }
              case status if status.isSparkReRunning => {
                log.info("status: SparkRunning. KYLIN falied and now SPARK is ready to rerun")
                jobProxy ! SQLSubmit(query)
              }
              case status if status.isSparkReRunFinished => {
                log.info("status: spark rerun finished. updating query_submit_log")
                dbLog ! QueryInfoUpdate(query)
              }
            }
          } catch {
            case e: Exception => e.printStackTrace(); log.error("未捕获的异常" + e.toString)
          }
        }
      }
      jobChannel.basicConsume(MQInfo.jobQueueName, true, consumer)
    }
    case KillRabbit => {
      context stop self
    }
    case query: QueryPublish => {
      // add try..catch, avoid rabbitmqActor crash
      try {
        val byteData = query.query.toByteArray
        publishChannel.basicPublish("", MQInfo.jobQueueName, null, byteData)
      } catch {
        case e: Exception => {
          e.printStackTrace()
        }
      }
    }
  }
}

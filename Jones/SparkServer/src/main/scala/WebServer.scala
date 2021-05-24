import java.util

import actors.{RabbitMQActor, SparkUdfActor, Start}
import adtl.platform.Dataproto.QueryStatus
import adtl.platform.Dataproto.QueryStatus.{CustomQueryRunning, Killed, Running, SparkReRunning, SyncQueryRunning, Waiting}
import akka.actor.{ActorSystem, _}
import akka.http.scaladsl.Http
import akka.stream.ActorMaterializer
import com.typesafe.config.{Config, ConfigFactory}
import models.config._
import models.{HiveMetaInfo, ServerInfo, ServerStatus}
import org.apache.log4j.{Level, Logger}
import org.apache.spark.sql.SparkSession
import routes.HttpRoute
import utils.{ConfigUtil, JedisHelper, MysqlUtil}

import scala.concurrent.ExecutionContext

object WebServer extends App {
  val logger = Logger.getLogger(getClass)

  def initUdf(spark:SparkSession,actorSystem: ActorSystem): Unit = {
    val udfJar = "com.sogou.adtl.dw.tool.udf.RegisterUtil"
    try {
      val toolClass = Class.forName(udfJar)
      val funcMap = toolClass.getMethod("registerAll", classOf[SparkSession]).
        invoke(toolClass.getClass, spark).asInstanceOf[util.HashMap[String, util.HashMap[String, String]]]

      val udfActor = actorSystem.actorOf(Props(classOf[SparkUdfActor]), "udfManager")
      udfActor ! Start(funcMap)
    } catch {
      case e: Exception => println(s"$udfJar not found, error ${e.getClass} ${e.getMessage}")
    }
  }

  def initSpark(): SparkSession = {
//    return null
    Logger.getLogger("org").setLevel(Level.FATAL)

    val redisConf: RedisConf = MainConf("basicConf").redis
    val host: String = redisConf.host
    val port: Int = redisConf.port
    val passwd: String = redisConf.passwd

    val spark = SparkSession.builder().appName(ServerInfo.Sparkname)
      .config("spark.serializer", "org.apache.spark.serializer.KryoSerializer")
      .config("spark.kryoserializer.buffer.max", "1024m")
      .config("spark.some.models.config.option", "some-value")
      .config("spark.redis.host", host)
      .config("spark.redis.port", port)
      .config("spark.redis.auth", passwd)
      .config("spark.redis.db", "0")
      .enableHiveSupport()
      .getOrCreate()


    //spark.sparkContext.addSparkListener(new SparkJobListener)
    HiveMetaInfo.genHiveMeta(spark)
    HiveMetaInfo.genSparkUDFMeta(spark)

    spark
  }

  // 启动前增加一个钩子, 把数据库中正在运行的job状态改为失败, 原因是系统重启了
  def repairLastRunningJob(): Unit = {
    val conn = MysqlUtil.getConnection("")
    val sql =
      s"""UPDATE query_submit_log
         |SET status=? , exception_msg=?
         |WHERE status in (${Running.value}, ${SparkReRunning.value}, ${SyncQueryRunning.value}, ${CustomQueryRunning.value})
         |""".stripMargin
    val status = Killed.value
    val exceptionMsg = "服务发生异常, 请重新提交任务"

    MysqlUtil.executeUpdate(conn, sql, status, exceptionMsg)
    MysqlUtil.closeConnection(conn)
    logger.info("Task run in last time killed")
  }

  def startApplication(): Unit = {

    // 读取配置文件
    val confPath = "application.conf"
    val confPathTest = "application-test.conf"
    val conf: Config = ConfigUtil.confFromFile(confPath)
    val confTest: Config = ConfigUtil.confFromFile(confPathTest)
//    val conf = ConfigFactory.load()
    MainConf.addConf("basicConf", MainConf(conf))
    //    MainConf.addConf("testConf", MainConf(confTest))
    val appConf = MainConf("basicConf").app

    ServerInfo.Sparkname = appConf.sparkName // adtl_spark
    ServerInfo.Host = appConf.host // 10.134.101.120
    ServerInfo.Port = appConf.port.toInt // 4803
    ServerInfo.Respath = appConf.resPath // search/odin/qs

    repairLastRunningJob()

    ServerStatus.isRunning.set(true)
    val spark = initSpark()
//    val spark = null

    JedisHelper.initPool()
    implicit val actorSystem = ActorSystem("data_server", conf)
    implicit val executor: ExecutionContext = actorSystem.dispatcher
    implicit val materializer: ActorMaterializer = ActorMaterializer()

    initUdf(spark, actorSystem)

    //增加
    val logger = Logger.getLogger(getClass)
    val bindingFuture = Http().newServerAt(ServerInfo.Host, ServerInfo.Port).bind(new HttpRoute(spark, logger).route)

    logger.info(s"Http Server on ${ServerInfo.Host}:${ServerInfo.Port} is running")


    val consumer = actorSystem.actorOf(Props(classOf[RabbitMQActor], spark), "rabbitMq")
    import RabbitMQActor._
    consumer ! JobConsumer
    while (ServerStatus.isRunning.get()) {
      Thread.sleep(1000)
    }
    consumer ! KillRabbit
    bindingFuture
      .flatMap(_.unbind()) // trigger unbinding from the port
      .onComplete(_ => actorSystem.terminate()) // and shutdown when done*/
  }

  startApplication()
}
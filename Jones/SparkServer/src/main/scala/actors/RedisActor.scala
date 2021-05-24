package actors

import actors.DBlogActor.QueryInfoUpdate
import akka.actor.{Actor, ActorLogging}
import models.ServerInfo
import org.apache.hadoop.fs.{FileSystem, Path}
import org.apache.spark.sql.SparkSession
import utils.{JedisHelper, PrintLogger, SparkUtil, XSSFUtil}

class RedisActor extends Actor with PrintLogger {

  val spark: SparkSession = SparkUtil.getSparkSession()

  override def receive: Receive = {
    case x: QueryInfoUpdate => {
      val queryInfo = x.query
      val queryId = queryInfo.queryId

      val respath = ServerInfo.Respath
      val originPath = s"${respath}/${queryInfo.dataPath}"

      val hadoopConf = spark.sparkContext.hadoopConfiguration
      val hdfs = FileSystem.get(hadoopConf)
      val fsPath = new Path(originPath) // hdfs
      val fis = hdfs.open(fsPath)

      info(s"redis cache start:\t ${queryInfo.dataPath}")
      new JedisHelper().saveRecorderFromIs(fis, queryId)
      info(s"redis cache finish:\t ${queryInfo.dataPath}")
    }
  }
}

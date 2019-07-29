package com.test.actors

import akka.actor.{Actor, Props}
import com.test.Constant
import com.test.database.{Item, MarkUtil}
import com.test.log.Logger
import com.test.loggers.PureLogger
import com.test.utils.{CommandUtil, HdfsUtil}

class BashActor extends Actor with Logger {
  override def preStart(): Unit = log.info("BashActor Application started")
  override def postStop(): Unit = log.info("BashActor Application stopped")

  val pureLogger = new PureLogger

  override def receive:Actor.Receive = {
    case x @BashMessage(cmd,workDir) => {
      val res = this.handle()
      if (res._1){
        log.info("can execute")
        log.info("execute bash scripts ")
        val oldPartition = res._2.get._1
        val newPartition = res._2.get._2
        val cu = new CommandUtil
        val combineCmd = s"$cmd ${oldPartition} ${newPartition}"
        log.info("execute cmd: "+ combineCmd)

        val is = cu.executeStrBash(combineCmd,workdir=workDir)
        cu.cmdOutput(is,pureLogger.pureinfo)
        log.info("complete bash scripts ")

        val util = HdfsUtil(Constant.hdfs,Constant.hdfsUser)
        util.appendToFile(res._3.get.toString,Constant.output)
        util.close()

        Constant.finish = true   // 此次处理完毕

      }else{
        log.info("skip this attempt")
      }
    }
    case _ =>
      log.warn("ignore action")
  }


  def handle():(Boolean,Option[Tuple2[String,String]],Option[Item]) = {
    if (Constant.finish) { // 上次梳理完毕
      Constant.finish = false // 此次开始处理
      log.info(Constant.zhHdfs + ":" + Constant.hdfsUser)
      val res = MarkUtil.getPartitions(Constant.output, Constant.zhHdfs, Constant.hdfs, Constant.hdfsUser)
      log.info((res._1, res._2).toString())
      if (!res._1.equals("")) {
        return (true, Some((res._1, res._2)), Some(res._3))
      }
    }
    return (false, None, None)
  }
}


object BashActor {

  def props(): Props = Props(new BashActor)

}



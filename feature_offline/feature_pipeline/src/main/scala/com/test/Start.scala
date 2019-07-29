package com.test

import com.test.actors.BashMessage
import com.test.database.MarkUtil
import com.test.utils.HdfsUtil

object Start {
  def main(args: Array[String]): Unit = {
    val system = Constant.system
    val scheduler = Constant.scheduler
    val actor = Constant.bashActor

//    var cmd = "bash start-feature.sh 0613 0614"
//    var wd = "/home/work/liujie32/test_jar"

    var cmd = Constant.cmdFirst
    var wd = Constant.cmdWD


    if (args.length>0){
      wd = args(0)
    }

    scheduler.schedule("everyDay", actor, new BashMessage(cmd,wd))

  }
}

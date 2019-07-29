package com.test

import akka.actor.ActorSystem
import com.test.actors.BashActor
import com.typesafe.akka.extension.quartz.QuartzSchedulerExtension
import com.typesafe.config.ConfigFactory

object Constant {

  var finish = true

  lazy val system = ActorSystem("Feature-Pipeline")
  lazy val scheduler = QuartzSchedulerExtension(system)
  lazy val bashActor = system.actorOf(BashActor.props(), "bashActor")

  lazy val config = ConfigFactory.load()
  lazy val zhHdfs = config.getString("auto.zhHdfs")
  lazy val dbSite = config.getString("auto.dbSite")
  lazy val output = config.getString("auto.output")
  lazy val hdfs = config.getString("auto.url")
  lazy val hdfsUser = config.getString("auto.user")

  lazy val cmdFirst = config.getString("cmd.first")
  lazy val cmdWD = config.getString("cmd.wd")
}

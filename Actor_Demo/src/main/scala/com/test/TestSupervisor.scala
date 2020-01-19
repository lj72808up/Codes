package com.test

import akka.actor.{Actor, ActorLogging, ActorRef, ActorSystem, OneForOneStrategy, Props, SupervisorStrategy, SupervisorStrategyConfigurator, Terminated}
import com.typesafe.config.ConfigFactory
import akka.actor.SupervisorStrategy._

import scala.concurrent.duration._

object Test {
  def main(args: Array[String]): Unit = {
    val config = ConfigFactory.parseString(
      """
akka.loglevel = "INFO"
akka.actor.debug{
  receive = on
  lifecycle = on
}
akka.actor{
    guardian-supervisor-strategy = "com.test.MySupervisorStrategy"
  //guardian-supervisor-strategy = "akka.actor.DefaultSupervisorStrategy"
}
""")

    val system = ActorSystem("MyTest", config) // akka://MyTest/system ;; akka://MyTest/user
    val father = system.actorOf(Props[Father], "Father")
    while (true) {
      father ! "start"
      Thread.sleep(1 * 1000)
    }
/*    val children = system.actorOf(Props[Children], "children")
    while (true) {
      children ! "do"
      Thread.sleep(1 * 1000)
    }*/
  }
}

class Father extends Actor with ActorLogging {
//  val children: ActorRef = context.actorOf(Props[Children], "children")
  var children: ActorRef = _
/*  children = context.watch(   // 对创建的actor进行watch, 可以监控actor的生命周期 (该子actor被停止后会收到Terminated消息)
    context.actorOf(Props[Children], "children")
  )*/

  override def preStart:Unit = {
    log.info("------ father 重新启动")
    children = context.watch(   // 对创建的actor进行watch, 可以监控actor的生命周期 (该子actor被停止后会收到Terminated消息)
      context.actorOf(Props[Children], "children")
    )
  }

/*  override val supervisorStrategy =
  // 下面的2个参数表示: 如果1个actor在1分钟之内重启次数超过10次, 该actor就会被停止
    OneForOneStrategy(maxNrOfRetries = 10, withinTimeRange = 1 minute) {
//      case _: IllegalNumberException => Stop  // actor Stop后若继续发到该actor, 则会产生死信错误
//      case _: NullPointerException     => Restart
      case _: Exception                => Escalate  // 对于策略无法覆盖到的异常, 可以使用默认的Escalate
    }*/
  override val supervisorStrategy =
  // 重启次数超过1次就会停止子actor
    OneForOneStrategy(maxNrOfRetries=1,withinTimeRange=Duration.Inf){
    case _:IllegalNumberException => Restart // Restart同样会触发actor的preStart方法
    case _:Exception => Escalate
  }


  override def receive: Receive = {
    case "start" =>
      val num = scala.util.Random.nextInt(100)
      if (num % 2 == 0) {
        children ! "do"
      } else {
        throw new Exception("father出错了")
      }

    case Terminated(actorRef) =>
      log.info("子actor停止了:"+ (children==actorRef))
      children = null
      self ! "reconnect"
    case "reconnect" =>
      log.info("actor停止后手动重启")
      preStart()
  }
}

class Children extends Actor with ActorLogging {

  override def preStart(): Unit = {
    log.info("+++++++ Children 启动了")
  }

  override def receive: Receive = {
    case "do" =>
      val num = scala.util.Random.nextInt(100)
      if (num % 2 == 0) {
        log.info("产生数字:" + num)
      } else {
        throw new IllegalNumberException("数字不能是奇数:" + num)
      }
  }
}

class IllegalNumberException(msg: String) extends Exception(msg)


// 顶级actor的监控策略: 覆盖配置文件中
final class MySupervisorStrategy extends SupervisorStrategyConfigurator {
  override def create(): SupervisorStrategy =
    OneForOneStrategy(maxNrOfRetries = -1, withinTimeRange = Duration.Inf) {
      case _: IllegalNumberException => Restart
      case t =>
        defaultDecider.applyOrElse(t, (_: Any) => Escalate)
    }
}
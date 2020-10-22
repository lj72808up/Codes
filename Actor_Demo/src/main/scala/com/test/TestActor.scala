package com.test

import akka.actor.{Actor, ActorRef, ActorSystem, Kill, PoisonPill, Props, Terminated}

import scala.concurrent.Future

/**
 * Actor 基础操作的demo
 * deference: https://doc.akka.io/docs/akka/2.5.26/actors.html
 */
object SwapperApp extends App {
  val system = ActorSystem("SwapperSystem")
  val swap = system.actorOf(Props[Swapper], name = "swapper")
  swap ! Swap // logs Hi
  swap ! Swap // logs Ho
  swap ! Swap // logs Hi
  swap ! Swap // logs Ho
  swap ! Swap // logs Hi
  swap ! Swap // logs Ho
}

case object Swap

class Swapper extends Actor {

  import context._

  def receive = {
    case Swap =>
      println("Hi")
      become({ // push "become包装的动作" 到动作栈的栈顶
        case Swap =>
          println("Ho")
          unbecome() // 将动作栈栈顶的方法弹出, 回复原来的栈顶
      }, discardOld = false) // push on top instead of replace

    case "stop" => context.stop(self)
  }
}


case object CleanAndStop
class P1(watchActor: ActorRef) extends Actor {
  context.watch(watchActor)

  override def receive: Receive = {
    case x: String => println(s"p1:A${x}")
    case x: Terminated => println(s"停止了${x.actor.path}")
  }
}

class P2 extends Actor {
  override def receive: Receive = {
    case x: String => var a = 1.0; for (i <- 1 to 900000000) {
      a = a * 1.01 / 3.1415 + 1.2
    }; println(s"p2:A${x}")
    case CleanAndStop => context.stop(self)
  }
}

object TestWatch {
  def main(args: Array[String]): Unit = {
    val system = ActorSystem("SwapperSystem")
    val p2 = system.actorOf(Props[P2], name = "p2")
    val p1 = system.actorOf(Props(classOf[P1], p2), name = "p1")
    p1 ! "aaa"
    println("发送 bbb")
    p2 ! "bbb"
    println("发送 kill")
    p2 ! CleanAndStop

  }
}
package com.test

import akka.actor.{Actor, ActorRef, ActorSystem, Props, Terminated}
import akka.routing.{ActorRefRoutee, Broadcast, GetRoutees, RoundRobinPool, RoundRobinRoutingLogic, Routees, Router, SmallestMailboxPool}


object TestRouter {

  def main(args: Array[String]): Unit = {
    val system = ActorSystem("MyTest") // akka://MyTest/system  akka://MyTest/user
    val starter = system.actorOf(Props[Starter], "starter-Entrance")
    starter ! 1 // 告诉 master 发送10个消息给worker
  }
}

class Starter extends Actor {
  val master: ActorRef = context.actorOf(Props[Master], "master-Entrance")

  override def receive: Receive = {
    case cnt: Int => {
      println(s"发送${cnt}条消息给master")
      for (_ <- 1 to cnt) {
        master ! Work()
        Thread.sleep(200)
      }
    }
    case x: Routees => println(x.routees)
    case _: MasterReplay => println(s"starter收到master的回复: " + sender().path)
    case _: WorkerReplay => println(s"starter收到worker的回复: " + sender().path)
    case x => println(x)
  }
}

/**
 * 第一种创建router的方式 (qucik start) , 同时测试 router 的消息透传
 */
class Master extends Actor {
  //TODO 第一种创建router的方式
  var router = {
    var i = 1
    val routees = Vector.fill(5) {
      val r = context.actorOf(Props(classOf[Worker], i)) // 子actor
      i = i + 1
      context.watch(r)
      ActorRefRoutee(r)
    }
    Router(RoundRobinRoutingLogic(), routees)
  }

  override def receive: Receive = {
    case w: Work =>
      //      println(s"master 收到消息")
      router.route(Broadcast(w), sender())
    //      println(s"发给 master 的 sender 是${sender().path}")   // sender 是 starter
    case Terminated(a) =>
      router = router.removeRoutee(a)
      val r = context.actorOf(Props[Worker])
      context.watch(r)
      router = router.addRoutee(r)
    case _: WorkerReplay => {
      println("master 收到 worker 的回复")
    }
  }
}

class Worker(idx: Int) extends Actor {
  override def receive: Receive = {
    case _: Work => {
      println(s"worker " + idx + s" 号工作 (${self.path}), " + s"收到${sender().path}的消息")
      //TODO 切换下面两种 routee 回复消息的方式, 体会 router 的透传
//      sender() ! WorkerReplay() // worker 会回复给started, router(master) 在消息回应时被隐藏
            context.parent ! WorkerReplay() // worker 会回复给 router(master)
    }
    case _ =>
  }
}

case class Work()

case class MasterReplay()

case class WorkerReplay()


/**
 * 第二种创建 router 的方式: 像普通actor一样创建
 */
class Parent extends Actor {
  val router: ActorRef =
    context.actorOf(RoundRobinPool(5).props(Props[ParentWorker]), "router2") // RoundRobinPool, SmallestMailboxPool
  override def receive: Receive = {
    case x: Int => {
      for (_ <- 1 to x) {
        router ! Broadcast(Work())
        Thread.sleep(200)
      }
    }

    case "check" => router ! GetRoutees   // Management Messages
    case x:Routees => println(x.routees)  // router会返回GetRoutees的消息
  }
}

class ParentWorker extends Actor {
  override def preStart(): Unit = {
    println(s"${self.path} - 开始了")
  }
  override def receive: Receive = {
    case _: Work => println(s"worker收到${sender().path}的消息 -- [${self.path}]")
    case _ =>
  }
}


object TestRouter2 {
  def main(args: Array[String]): Unit = {
    val system = ActorSystem("TestRouter2")
    val p = system.actorOf(Props[Parent], "parent")
    p ! 2
    p ! "check"
  }
}



// todo 自定义Router, 有限选择没有 message 的actor, 然后选择当前时间最少的actor, 最好有actor空闲时从其他actor拿消息
//  SmallestMailboxRoutingLogic
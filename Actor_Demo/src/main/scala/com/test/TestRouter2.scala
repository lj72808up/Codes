package com.test

import java.util.concurrent.ConcurrentLinkedQueue

import akka.actor.{ActorRef, ActorSystem, Props}
import akka.dispatch.{Envelope, MessageQueue, ProducesMessageQueue, UnboundedMessageQueueSemantics}
import akka.event.LoggerMessageQueueSemantics
import com.typesafe.config.Config

object TestRouter3 {
  def main(args: Array[String]): Unit = {
    val system = ActorSystem("TestSystem")
    /*val a = system.actorOf(
      Props[MyActor].withDispatcher("my-dispatcher")
    )*/

  }
}


case class MonitorEnvelope(queueSize: Int, receiver: String, entryTime: Long, handle: Envelope)

case class MailboxStatistics(queueSize: Int,
                             receiver: String,
                             sender: String,
                             entryTime: Long,
                             exitTime: Long)


class MonitorQueue(val system: ActorSystem) extends MessageQueue
  with UnboundedMessageQueueSemantics
  with LoggerMessageQueueSemantics {
  private final val queue = new ConcurrentLinkedQueue[MonitorEnvelope]()

  override def enqueue(receiver: ActorRef, handle: Envelope): Unit = {
    val env = MonitorEnvelope(
      queueSize = queue.size() + 1,
      receiver = receiver.toString(),
      entryTime = System.currentTimeMillis(),
      handle = handle
    )
    queue.add(env)
  }

  override def dequeue(): Envelope = {
    val monitor = queue.poll()
    if (monitor != null) {
      monitor.handle.message match {
        case stat: MailboxStatistics => //skip message
        case _ => {
          val stat = MailboxStatistics(
            queueSize = monitor.queueSize,
            receiver = monitor.receiver,
            sender = monitor.handle.sender.toString(),
            entryTime = monitor.entryTime,
            exitTime = System.currentTimeMillis())
          system.eventStream.publish(stat)
        }
      }
      monitor.handle
    } else {
      null
    }
  }

  override def numberOfMessages: Int = queue.size()

  override def hasMessages: Boolean = !queue.isEmpty

  override def cleanUp(owner: ActorRef, deadLetters: MessageQueue): Unit = {
    if (hasMessages) {
      var envelope = this.dequeue()
      while (envelope ne null) {
        deadLetters.enqueue(owner, envelope)
        envelope = this.dequeue()
      }
    }
  }
}


class MonitorMailboxType(settings: ActorSystem.Settings, config: Config)
  extends akka.dispatch.MailboxType
    with ProducesMessageQueue[MonitorQueue]{
  final override def create(owner: Option[ActorRef],
                            system: Option[ActorSystem]): MessageQueue = {
    system match {
      case Some(sys) =>
        new MonitorQueue(sys)
      case _ =>
        throw new IllegalArgumentException("requires a system")
    }
  }
}
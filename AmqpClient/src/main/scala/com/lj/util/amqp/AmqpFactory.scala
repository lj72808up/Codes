package com.lj.util.amqp

import com.rabbitmq.client.{Channel, Connection}
import com.rabbitmq.client.ConnectionFactory

object AmqpFactory {
  lazy val consumeConn:Connection = getConnection()
  lazy val produceConn:Connection = getConnection()
  //  1. Don’t share channels between threads
  //     sharing Channel instances between threads is something to be avoided.
  //     should prefer using a Channel per thread instead of sharing the same Channel across multiple threads
  //     sharing Channel will result in incorrect frame interleaving on the wire, double acknowledgements and so on
  //  2. Don’t open and close connections or channels repeatedly
  //  3. Separate connections for publisher and consumer
  //  4. A large number of connections and channels might affect the RabbitMQ management interface performance

  private def getConnection():Connection = {

    val factory = new ConnectionFactory
    // "guest"/"guest" by default, limited to localhost connections
    //    "amqp://userName:password@hostName:portNumber/virtualHost"
    factory.setUri("amqp://admin:admin@10.140.12.131:5672/test")
    // connection that will recover automatically(首次创建connection失败不会重试, channel级别的异常不会触发重试)
    factory.setAutomaticRecoveryEnabled(true)
    factory.setTopologyRecoveryEnabled(true) // 自动恢复exchange,queue,binding对象
    // attempt recovery every 10 seconds
    factory.setNetworkRecoveryInterval(10 * 1000)
    factory.setRequestedHeartbeat(30);
    return factory.newConnection
  }


  def deleteQueue(channel: Channel, queueName:String,when:DeleteWhen):Unit = {
    when match {
//    delete a queue only if it is empty
      case EmptyQueue() => channel.queueDelete(queueName, false, true)
//    or if it is not used (does not have any consumers)
      case NoConsumer() => channel.queueDelete(queueName, true, false)
      case _ =>
    }
  }


}


trait DeleteWhen
case class NoConsumer() extends DeleteWhen
case class EmptyQueue() extends DeleteWhen
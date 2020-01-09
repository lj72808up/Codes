package com.sogou.util.amqp

import com.rabbitmq.client.{Channel, Connection}

class AmqpConnection {
  private var conn:Connection = _
  //  sharing Channel instances between threads is something to be avoided.
  //  should prefer using a Channel per thread instead of sharing the same Channel across multiple threads
  //  sharing Channel will result in incorrect frame interleaving on the wire, double acknowledgements and so on
  //  可以并发的在一个Connection内使用多个Channel进行Publish或者Receive

  def getConnection():Connection = {
    import com.rabbitmq.client.ConnectionFactory
    val factory = new ConnectionFactory
    // "guest"/"guest" by default, limited to localhost connections
    //    "amqp://userName:password@hostName:portNumber/virtualHost"
    factory.setUri("amqp://liujie02:123456@10.134.33.80:5672/db_server")
    // connection that will recover automatically(首次创建connection失败不会重试, channel级别的异常不会触发重试)
    factory.setAutomaticRecoveryEnabled(true)
    factory.setTopologyRecoveryEnabled(true) // 自动恢复exchange,queue,binding对象
    // attempt recovery every 10 seconds
    factory.setNetworkRecoveryInterval(10 * 1000)
    conn = factory.newConnection
    conn
  }

  def close(channel: Channel):Unit = {
    channel.close()
    conn.close()
  }

  // routingKey==bindingKey时, 消息被转发
  def bindDirectExchange(channel: Channel,exchangeName:String, queueName:String, routingKey:String):Unit = {
    channel.exchangeDeclare(exchangeName, "direct", true)
    // 持久化, 飞排它, 非自动删除的queue
    channel.queueDeclare(queueName, true, false, false, null)
    channel.queueBind(queueName, exchangeName, routingKey)
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
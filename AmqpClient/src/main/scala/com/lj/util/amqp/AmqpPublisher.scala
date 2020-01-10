package com.lj.util.amqp

import com.rabbitmq.client.{AMQP, Channel, MessageProperties, ReturnListener}

class AmqpPublisher(channel: Channel, exchangeName: String, routingKey: String, mandatory: Boolean=false) {

  // mandatory: 是否在消息无法被routing到某个队列时, 讲消息返还给队列
  // 在broker重启之后可以继续消费消息, 需要将exchange,queue,消息都设置为持久化的
  def publish(msg: Array[Byte]): Unit = {
    if (mandatory) {
      channel.addReturnListener(new ReturnListener {
        override def handleReturn(replyCode: Int,
                                  replyText: String,
                                  exchange: String,
                                  outingKey: String,
                                  properties: AMQP.BasicProperties,
                                  body: Array[Byte]): Unit = {}
      })
    }
    channel.basicPublish(exchangeName, routingKey, mandatory,
      MessageProperties.PERSISTENT_BASIC,
      msg)
  }

  // routingKey==bindingKey时, 消息被转发
  def bindDirectExchange(channel: Channel,exchangeName:String, queueName:String, routingKey:String):Unit = {
    channel.exchangeDeclare(exchangeName, "direct", true)
    // 持久化, 非排它, 非自动删除的queue
    channel.queueDeclare(queueName, true, false, false, null)
    channel.queueBind(queueName, exchangeName, routingKey)
  }

  def close(): Unit ={
    channel.close()
  }
}

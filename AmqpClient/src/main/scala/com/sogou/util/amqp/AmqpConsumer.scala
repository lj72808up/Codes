package com.sogou.util.amqp

import com.rabbitmq.client.{AMQP, Channel, DefaultConsumer, Envelope}
import java.io.IOException
import java.util.concurrent.{ExecutorService, Executors}

class AmqpConsumer(channel: Channel) {
  //Consumer tags are used to cancel consumers.
  private var consumerTag:String = _
  private val es: ExecutorService = Executors.newFixedThreadPool(10)

  def pushConsume(queueName:String):Unit = {
    val autoAck = false
    consumerTag = channel.basicConsume(queueName, autoAck, new DefaultConsumer(channel) {
      @throws[IOException]
      override def handleDelivery(consumerTag: String,
                                  envelope: Envelope,
                                  properties: AMQP.BasicProperties,
                                  body: Array[Byte]): Unit = {
        val routingKey = envelope.getRoutingKey
        val contentType = properties.getContentType
        val deliveryTag = envelope.getDeliveryTag

        //TODO (process the message components here ...)
        println(s"${Thread.currentThread().getName}:recieve msg: ${new String(body)}")
        Thread.sleep(2000) // 模仿业务消耗
        if(!autoAck){
          // 手动提交
          channel.basicAck(deliveryTag, false)  // 这种情况下无论起多少个线程worker, 只有当1条消息ack后, 下一条才能处理
        }
      }
    })

  }

  def close():Unit = {
    channel.basicCancel(consumerTag)
  }
}

package com.lj.util.amqp

import com.rabbitmq.client.{AMQP, Channel, DefaultConsumer, Envelope}
import java.io.IOException
import java.util.concurrent.{ExecutorService, Executors}

class AmqpConsumer(channel: Channel) {
  //Consumer tags are used to cancel consumers.
  private var consumerTag:String = _

  def pushConsume(queueName:String, name:String , fun: () =>Unit , autoAck:Boolean=false):Unit = {

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
        println(s"${name}-${Thread.currentThread().getName}:recieve msg: ${new String(body)}")
        try{   // 如果不捕获异常, 会导致发生异常后, 回调函数无法继续使用
          fun()
        }catch {
          case e: Exception => e.printStackTrace()
        }finally {
          if (!autoAck) {
            // 手动提交
            channel.basicAck(deliveryTag, false) // 这种情况下无论起多少个线程worker, 只有当1条消息ack后, 下一条才能处理
          }
        }
      }
    })

  }


  // routingKey==bindingKey时, 消息被转发
  def bindDirectExchange(channel: Channel,exchangeName:String, queueName:String, routingKey:String):Unit = {
    channel.exchangeDeclare(exchangeName, "direct", true)
    // 持久化, 非排它, 非自动删除的queue
    channel.queueDeclare(queueName, true, false, false, null)
    channel.queueBind(queueName, exchangeName, routingKey)
  }


  def cancel():Unit = {
    channel.basicCancel(consumerTag)
  }

  def close(): Unit ={
    channel.close()
  }
}

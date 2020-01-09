package com.test

import com.sogou.util.amqp.{AmqpConnection, AmqpConsumer, AmqpPublisher}

object TestAmqp {
  val routingKey:String = "test"
  val consumerNum:Int = 5  // 5个consumer需要设置5个queue (1个consumer只能与1个queue对应起来)
  val amqp = new AmqpConnection()
  val connection = amqp.getConnection()

  def testProducer(): Unit = {
    val channel = connection.createChannel()
    amqp.bindDirectExchange(channel,"test_exchange","test_queue",routingKey)

    val publisher = new AmqpPublisher(channel,"test_exchange",routingKey)

    for (i <- 1 to 100){
      publisher.publish(s"[${Thread.currentThread().getName}]-haha:$i".getBytes)
    }

    println("finish")

  }

  def testConsumer():Unit = {
    val channel = connection.createChannel()
    amqp.bindDirectExchange(channel,"test_exchange","test_queue",routingKey)

    val consumer = new AmqpConsumer(channel)
    consumer.pushConsume("test_queue")
  }

  /**
   * 2个producer可以并行的发送数据
   */
  def testmultiProducer(): Unit = {
    val run1 = new Runnable {
      override def run(): Unit = {
        testProducer()
        println("1 finish")
      }
    }
    val run2 = new Runnable {
      override def run(): Unit = {
        testProducer()
        println("2 finish")
      }
    }

    new Thread(run1).start()
    new Thread(run2).start()

    testConsumer()
  }

/*  def testMultiConsumer():Unit = {
    testProducer()
    val run1 = new Runnable {
      override def run(): Unit = {
        testConsumer()
      }
    }
    val run2 = new Runnable {
      override def run(): Unit = {
        testConsumer()
      }
    }
    new Thread(run1).start()
    new Thread(run2).start()
  }*/

  def main(args: Array[String]): Unit = {
//    testMultiConsumer
//    testmultiProducer()
//    testProducer()
    testConsumer()
  }
}

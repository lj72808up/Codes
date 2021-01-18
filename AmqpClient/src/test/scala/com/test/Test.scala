package com.test

import java.time.LocalTime
import java.util.Date

import com.lj.util.amqp.{AmqpConsumer, AmqpFactory, AmqpPublisher}

object Test {
  def testProduc(): Unit = {
    val conn = AmqpFactory.produceConn

    val th1 = new Thread(new Runnable {
      override def run(): Unit = {
        val channel = conn.createChannel() // 1. 获取连接
        val publisher = new AmqpPublisher(channel, "test_exchange", "")
        publisher.bindDirectExchange(channel, "test_exchange", "test_queue", "") // 2. channel绑定queue

        for (i <- 1 to 100) {
          println(s"xxxx${i}")
          publisher.publish(s"xxxx${i}".getBytes())
        }
        println("线程完毕")
        channel.close()
      }
    })

    th1.start()
//    th1.join()
//    conn.close()
    println("消息生产完毕")
  }


  // 测试不同线程共享同一个Connection, 但使用各自的channel的情况下(手动提交ack), 多线程能否并行消费?
  // 答案:
  //  (1)未加qos之前不能:
  //    a. 只有1个线程的channel可以消费, 只会打印"Thread-A"或"Thread-B";
  //    b. 即使用2个进程分别用各自的Connection脸上队列, 也只有一个进程可以消费
  //  (2)加了qos之后可以:
  //    a. qos代表unack的最大条数, 不设置的话rabbitmq会把queue中的所有消息一股脑的发送给consumer, 造成内存溢出
  //    b. 配置qos后也会让mq有机会把消息分发给不同的consumer
  def testConsume(): Unit ={

    def fun1():Unit={}
    val conn = AmqpFactory.consumeConn
    new Thread(new Runnable {
      override def run(): Unit = {
        val channel = conn.createChannel() // 生产者的channel要先绑定queue, 消费者的不用
        channel.basicQos(10);  // 并行消费的关键设置
        val consumer = new AmqpConsumer(channel)
        consumer.bindDirectExchange(channel, "test_exchange", "test_queue", "") // 2. channel绑定queue
        consumer.pushConsume("test_queue", "Thread-A",fun1)
      }
    }).start()

    /*new Thread(new Runnable {
      override def run(): Unit = {
        val channel = conn.createChannel()  // 避免在线程之间共享channel
        channel.basicQos(2);
        val consumer = new AmqpConsumer(channel)
        consumer.bindDirectExchange(channel, "test_exchange", "test_queue", "") // 2. channel绑定queue
        consumer.pushConsume("test_queue", "Thread-B",fun1)
      }
    }).start()*/
  }

  def main(args: Array[String]): Unit = {
//    Test.testProduc()
    Test.testConsume()
  }
}

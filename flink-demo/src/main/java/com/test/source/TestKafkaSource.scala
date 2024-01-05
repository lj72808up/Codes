package com.test.source

import org.apache.flink.api.common.eventtime.{SerializableTimestampAssigner, WatermarkStrategy}
import org.apache.flink.connector.kafka.source.KafkaSource
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer
import org.apache.flink.connector.kafka.source.reader.deserializer.KafkaRecordDeserializationSchema
import org.apache.flink.streaming.api.scala.{StreamExecutionEnvironment, _}
import org.apache.kafka.clients.producer.ProducerConfig
import org.apache.kafka.common.serialization.StringDeserializer

import java.time.Duration
import java.util.{Date, Properties}

object TestKafkaSource {
  private val producerConfig = getKafkaProperty()

  def main(args: Array[String]): Unit = {
    val env = StreamExecutionEnvironment.getExecutionEnvironment
    //  makeCheckPoint(env)
    //用最新的api 消费kafka，可以指定分区消费。
    val source = KafkaSource.builder[String]()
      .setBootstrapServers("10.13.11.13")
      .setGroupId("test_group")
      .setTopics("test_topic")
//      .setPartitions(new java.util.Set(new TopicPartition("", 1)))
      .setStartingOffsets(OffsetsInitializer.latest())
      .setDeserializer(KafkaRecordDeserializationSchema.valueOnly(classOf[StringDeserializer])).build()
    // 水位线
    val watermarkStrategy = WatermarkStrategy.forBoundedOutOfOrderness[String](Duration.ofSeconds(120)).withTimestampAssigner(new SerializableTimestampAssigner[String]() { // 抽取时间戳的逻辑
      override def extractTimestamp(element: String, recordTimestamp: Long): Long = {
        new Date().getTime
      }
    })

    // source 加入水位线策略。
    env.fromSource[String](source, watermarkStrategy, "coupon_source").rebalance
      .map(log => log)


  }

  private def getKafkaProperty(): Properties = {
    val producerConfig = new Properties()
    producerConfig.setProperty(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, "10.13.13.1")
    producerConfig.setProperty(ProducerConfig.RETRIES_CONFIG, "10")
    producerConfig.setProperty(ProducerConfig.RETRY_BACKOFF_MS_CONFIG, "200")
    producerConfig.setProperty(ProducerConfig.METADATA_MAX_AGE_CONFIG, "300000")
    producerConfig.setProperty(ProducerConfig.COMPRESSION_TYPE_CONFIG, "snappy")

    producerConfig
  }

}

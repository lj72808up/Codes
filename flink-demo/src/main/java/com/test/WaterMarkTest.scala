package com.test

import org.apache.flink.api.common.eventtime.{SerializableTimestampAssigner, Watermark, WatermarkGenerator, WatermarkOutput, WatermarkStrategy}
import org.apache.flink.api.scala.createTypeInformation
import org.apache.flink.streaming.api.functions.source.FileProcessingMode
import org.apache.flink.streaming.api.scala.{DataStream, StreamExecutionEnvironment}
import org.apache.flink.streaming.api.windowing.assigners.TumblingEventTimeWindows
import org.apache.flink.streaming.api.windowing.time.Time

import java.time.Duration

object WaterMarkTest {
  /** 6. WatermarkGenerator
   * 水印策略中的时间戳抽取很简单,不去讨论, 这里讨论水印生成者 WatermarkGenerator. 实现这个接口的2个方法即可. 分两类:
   * 1) onEvent(): 对每个事件产生一次水印  (下面的 BoundedOutOfOrdernessGenerator.class)
   * 2) onPeriodicEmit(): 虽然还是会接受每个事件产生一个水印,但不立刻发送, 会等待一定间隔后延迟发送 (下面的TimeLagWatermarkGenerator.class)
   * 还有一种 Punctuated 的水印生成者, 当遇到包含水印标记的事件时, 才发射水印时间
   * 如下的 PunctuatedAssigner.class
   */
  def main6(args: Array[String]): Unit = {
    // 1) BoundedOutOfOrdernessGenerator.class
    /** 假定消息的到来是无序的, 但只在一个时间段内无序
     */
    class BoundedOutOfOrdernessGenerator extends WatermarkGenerator[Object] {
      val maxOutOfOrderness = 3500L // 3.5 seconds, 最大无序事件
      var currentMaxTimestamp: Long = _

      override def onEvent(element: Object, eventTimestamp: Long, output: WatermarkOutput): Unit = {
        currentMaxTimestamp = Math.max(eventTimestamp, currentMaxTimestamp)
      }

      override def onPeriodicEmit(output: WatermarkOutput): Unit = {
        // emit the watermark as current highest timestamp minus the out-of-orderness bound
        output.emitWatermark(new Watermark(currentMaxTimestamp - maxOutOfOrderness - 1));
      }
    }
    // 2) TimeLagWatermarkGenerator.class
    /** 假定消息时间到达 flink 的处理时间有一个固定时间的延迟
     */
    class TimeLagWatermarkGenerator extends WatermarkGenerator[Object] {
      val maxTimeLag = 5000L // 5 seconds 延迟时间

      override def onEvent(element: Object, eventTimestamp: Long, output: WatermarkOutput): Unit = {
        // don't need to do anything because we work on processing time
      }

      override def onPeriodicEmit(output: WatermarkOutput): Unit = {
        output.emitWatermark(new Watermark(System.currentTimeMillis() - maxTimeLag));
      }
    }
    // 3) PunctuatedAssigner.class
    /** 事件中带有是否为水印的标记
     */
    class PunctuatedAssigner extends WatermarkGenerator[MyEvent] {
      override def onEvent(event: MyEvent, eventTimestamp: Long, output: WatermarkOutput): Unit = {
        if (event.hasWatermarkMarker()) {
          output.emitWatermark(new Watermark(event.getWatermarkTimestamp()))
        }
      }

      override def onPeriodicEmit(output: WatermarkOutput): Unit = {
        // don't need to do anything because we emit in reaction to events above
      }
    }
    abstract class MyEvent() {
      def hasWatermarkMarker(): Boolean

      def getWatermarkTimestamp(): Long
    }
  }

  /** 5.WatermarkStrategy 可以用在2个地方:
   * (1) source算子上 (如下套在source算子上的完整写法)
   * (2) 和非source算子上
   */
  def main(args: Array[String]): Unit = {
    val env = StreamExecutionEnvironment.getExecutionEnvironment

    var stream:DataStream[MyEvent] = env.socketTextStream("localhost", 9999)
                    .map(x=> MyEvent(java.lang.Long.parseLong(x.split(",")(0)),
                        java.lang.Integer.parseInt(x.split(",")(1)))
                    )

    val strategy = WatermarkStrategy.forBoundedOutOfOrderness(Duration.ofSeconds(20))
      .withTimestampAssigner(new SerializableTimestampAssigner[MyEvent] {
        override def extractTimestamp(element: MyEvent, recordTimestamp: Long): Long = {
          println(element + "=====================>" + recordTimestamp)
          element.ts
        }
      })
    val waterMarkStream = stream.assignTimestampsAndWatermarks(strategy)
    waterMarkStream
      .keyBy(_=>1)   // 这里简单设置不分组
      .window(TumblingEventTimeWindows.of(Time.seconds(10)))
      .reduce((a,b) => MyEvent(Math.min(a.ts,b.ts), a.value+b.value))
      .print()

    /*    val withTimestampsAndWatermarks: DataStream[MyEvent] = stream
          .filter( _.severity == WARNING )
          .assignTimestampsAndWatermarks(<watermark strategy>)

            withTimestampsAndWatermarks
            .keyBy( _.getGroup )
            .window(TumblingEventTimeWindows.of(Time.seconds(10)))
            .reduce( (a, b) => a.add(b) )
    //        .addSink(...)*/
    env.execute("watermarkTest")

  }

  /** 4. EventTime
   * (1) 首先,想使用 eventTime, 必须制定记录中哪个字段代表事件的时间戳. 使用 TimestampAssigner 制定获取时间时间戳的方法
   * (2) WatermarkGenerator 接口,随着 eventTime 字段的到来不断地产生 watermark, 告诉 flink 当前处理到的事件时间
   * (3) WatermarkStrategy 类综合了时间戳分配(TimestampAssigner)和水印产生策略(WatermarkGenerator)
   *
   * (4) TimestampAssigner 不用每次都指定, 比如使用 Kafka/Kinesis 时,可以直接从记录中获取时间
   * (5) 空闲分区问题会导致 watermark 无法更新
   * flink 水印的计算,是由所有分区水印的最小值计算而来. 如果某个分区持续空闲, watermark 就无法更新. 这个问题可以用空闲时间来解决
   * ```WatermarkStrategy.forBoundedOutOfOrderness[(Long, String)](Duration.ofSeconds(20))
   * .withIdleness(Duration.ofMinutes(1))```
   */
  def main4(args: Array[String]): Unit = {
    WatermarkStrategy
      // 内置的水印生成策略
      // 和上面 TimeLagWatermarkGenerator.class 作用类似, 指定事件的最大延迟时间
      .forBoundedOutOfOrderness[(Long, String)](Duration.ofSeconds(20))
      // 指定eventTime抽取方法
      .withTimestampAssigner(new SerializableTimestampAssigner[(Long, String)] {
        override def extractTimestamp(element: (Long, String), recordTimestamp: Long): Long = element._1
      })
  }
}

case class MyEvent(ts: Long, value: Int)
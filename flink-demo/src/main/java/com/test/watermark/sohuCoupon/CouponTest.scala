package com.test.watermark.sohuCoupon

import org.apache.flink.api.common.eventtime.{SerializableTimestampAssigner, WatermarkStrategy}
import org.apache.flink.api.common.functions.ReduceFunction
import org.apache.flink.configuration.{Configuration, RestOptions}
import org.apache.flink.streaming.api.functions.sink.PrintSinkFunction
import org.apache.flink.streaming.api.scala.function.ProcessWindowFunction
import org.apache.flink.streaming.api.scala.{StreamExecutionEnvironment, _}
import org.apache.flink.streaming.api.windowing.assigners.TumblingEventTimeWindows
import org.apache.flink.streaming.api.windowing.time.Time
import org.apache.flink.streaming.api.windowing.windows.TimeWindow
import org.apache.flink.util.Collector

import java.text.SimpleDateFormat
import java.time.Duration

object CouponTest {
  val df = new SimpleDateFormat("yyyyMMddHHmmss")
  val duration = 5
  def main(args: Array[String]): Unit = {
    val env = StreamExecutionEnvironment.getExecutionEnvironment
    val source = env.addSource(new MySource)
    // 水位线策略。
/*    val watermarkStrategy = WatermarkStrategy.forBoundedOutOfOrderness[String](Duration.ofSeconds(1))
      //      .withIdleness(Duration.ofSeconds(5))
      .withTimestampAssigner(new SerializableTimestampAssigner[String]() { // 抽取时间戳的逻辑
        override def extractTimestamp(element: String, recordTimestamp: Long): Long = {
          val date = element.split("_")(0) // date: 20231212104501
          df.parse(date).getTime
        }
      })*/

    val watermarkStrategy = new MyWatermarkStrategy()

    // 加入水位线策略。
    source.assignTimestampsAndWatermarks(watermarkStrategy)
      .map(log => {
        val splits = log.split("_")
        val key = splits(1)
        (key, 1)
      })
      .keyBy(_._1)
      /** 每种 window, 都有一个 getDefaultTrigger() 方法, 返回绑定的 triggr
          *使用 EventTime 的窗口, 默认的是 EventTimeTrigger, TumblingEventTimeWindows 的默认 trigger 就是 EventTimeTrigger */
      .window(TumblingEventTimeWindows.of(Time.seconds(5)))
      .trigger(new MyTrigger(2))    // processDelay, 要设置的大于 forBoundedOutOfOrderness 的乱序 delay, 才能让 eventTimer 比 processTimer早触发
      .reduce(new ReduceFunction[(String, Int)] {
        override def reduce(value1: (String, Int), value2: (String, Int)): (String, Int) =
          (value1._1, value1._2 + value2._2)
      },new MyProcessWindowFunction)  // reduce 算子结合 ProcessWindowFunction, 即不用缓存窗口内的所有数据, 又可以获取 context 里的窗口时间
      .addSink(new PrintSinkFunction[String]())


    env.execute("COUPON_IMP_COUNT")


  }
}


// ProcessWindowFunction[IN, OUT, KEY, W <: Window]
class MyProcessWindowFunction extends ProcessWindowFunction[(String, Int), String, String, TimeWindow] {
  private val format = new SimpleDateFormat("yyyyMMddHHmmss")

  /**
   * key – The key for which this window is evaluated.
   * context – main effect of this function is obtain the context in which the window is being evaluated.
   * elements – all elements by key and by window will be cached, so you'd better use this as reduce function's param, or it will be inefficient
   * out – A collector for emitting elements.
   */
  override def process(key: String,
                       context: Context,
                       elements: Iterable[(String, Int)],
                       out: Collector[String]): Unit = {
    //用窗口时间作为入库时间，更加准确
    val windowStart = format.format(context.window.getStart / (CouponTest.duration * 1000) * (CouponTest.duration * 1000))
    elements.foreach { obj =>
      val cnt = obj._2
      //val datetime = new Timestamp(DateUtil.convertStringToLongTime(fiveminutes, "yyyyMMddHHmm") * 1000)
      out.collect(s"${windowStart}\t${cnt}")
    }
  }
}
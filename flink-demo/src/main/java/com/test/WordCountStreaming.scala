package com.test

import org.apache.flink.api.common.eventtime.{SerializableTimestampAssigner, Watermark, WatermarkGenerator, WatermarkOutput, WatermarkStrategy}
import org.apache.flink.streaming.api.datastream.DataStreamUtils
import org.apache.flink.streaming.api.scala.StreamExecutionEnvironment
import org.apache.flink.streaming.api.windowing.assigners.TumblingProcessingTimeWindows
import org.apache.flink.streaming.api.windowing.time.Time
import org.apache.flink.streaming.api.scala._

import java.time.Duration

object WordCountStreaming{
  /** 3. flink 的 streaming 和 batch 模式
   * (1) 效率上讲, 有界数据源应该优先使用 batch 模式, 因为此模式下会对 join, shuffle 等操作有额外优化
   * (2) 推荐使用命令行指定模式
   *    `$ bin/flink run -Dexecution.runtime-mode=BATCH examples/streaming/WordCount.jar`
   * (3) batch 和 streaming 模式的一些区别
   *  a. streaming 模式一旦收到一条消息,会立刻处理后发给 downstream , 而 batch 模式则可以把任务划分成 stage
   *     每个 stage 的所有数据处理完成后再发往下一个 stage
   *  b. streaming 模式下, 靠 StateBackend 控制状态的存储, 和如何 checkpoint
   *     batch 模式不用记录状态
   * (4) Event Time 和 Watermarks
   *     streaming 模式下, flink 悲观的认为事件会乱序的到来, 即时间戳为 t 的事件可能会在时间戳 t+1 的事件后面到来.
   *     因此, flink 使用时间戳为 t 的watermark 表示,后面不会有时间戳小于 t 的事件到来
   * (5) Processing time
   *     这个是一条记录被处理时所在机器的机器时间. 基于 Processing time 的计算结果重算后无法获得相同结果(时间戳变化了)
   * (6) 故障恢复
   *     STREAMING 模式下, flink 使用 checkpoint 进行故障恢复, flink 会在 checkpoint 处重启
   * (7) 算子行为
   *     像 reduce() 或者 sum() 这些算子, 再 streaming 模式下每当一条新记录达到会滚动更新结果. 而 batch 模式则直接输出最终结果
   */
  def main3(args: Array[String]): Unit = {}
  /** 2. 从内存集合构建 DataStream
   * 内存集合的并行度为1 (不能并行处理)
   */
  def main2(args: Array[String]): Unit = {
    val env = StreamExecutionEnvironment.getExecutionEnvironment
    // Create a DataStream from a list of elements
    val myInts = env.fromElements("hehe","haha")

    // Create a DataStream from any Java collection
    val myColl = env.fromCollection( Array("hehe2","haha2"))

    val resDataStream = myInts.flatMap { _.toLowerCase.split("\\W+") filter { _.nonEmpty } }
      .map { (_, 1) }
      .keyBy(_._1)
//      .window(TumblingProcessingTimeWindows.of(Time.seconds(5)))
      .sum(1)

    // 收集到本地, 用于打印测试结果
    val value = DataStreamUtils.collect(resDataStream.javaStream)
    println("======================>"+value)

    // (2) 上述所有都不会执行, 只有最后的 execute 才触发
    env.execute("Window Stream WordCount")
  }
  /** 1. flink 的 wordcount 程序
   *   滚动窗口 5 秒一次
   */
  def main1(args: Array[String]) {

    // (1) 提供的静态方法, 自动识别是本地环境还是集群环境
    val env = StreamExecutionEnvironment.getExecutionEnvironment
    env.setBufferTimeout(5);  // 底层 netty 发送数据的缓冲延时, 缓冲区满或超时才产生网络通信

    val text = env.socketTextStream("localhost", 9999)

    val counts = text.flatMap { _.toLowerCase.split("\\W+") filter { _.nonEmpty } }
      .map { (_, 1) }
      .keyBy(_._1)
      .window(TumblingProcessingTimeWindows.of(Time.seconds(5)))
      .sum(1)

    counts.print()

    // (2) 上述所有都不会执行, 只有最后的 execute 才触发
    env.execute("Window Stream WordCount")

    //todo IterativeStream 的作用
  }

}

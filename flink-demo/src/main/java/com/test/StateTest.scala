package com.test

import org.apache.flink.api.common.functions.RichFlatMapFunction
import org.apache.flink.api.common.state.{ValueState, ValueStateDescriptor}
import org.apache.flink.configuration.Configuration
import org.apache.flink.streaming.api.scala.StreamExecutionEnvironment
import org.apache.flink.streaming.api.scala._
import org.apache.flink.util.Collector

/** 3- 状态和容错 */
object StateTest {
  /** 1.keyed state
   * state 的作用范围,都是限定在起对应的一个 key 上
   * (1)ValueState: 一个 key 一个 state
   * (2)ListState: 一个key一个列表state. 可以 add(T) 或 addAll(List<T>)
   * (3)ReducingState<T>: 和 ListState 类似,也是通过 add(T) 增加状态, 只是需要指定一个 ReduceFunction
   * (4)AggregatingState<IN, OUT>: 和ReducingState类似,只是in和out是两个类型
   * (5)MapState<UK, UV>: 状态是 kv 对
   * 2. 自定义函数中整合 state
   * a. 为了使用 state, 先要创建对应的 StateDescriptor: ValueStateDescriptor/ ListStateDescriptor/ AggregatingStateDescriptor/ ReducingStateDescriptor, or a MapStateDescriptor.
   * b. StateDescriptor 需要通过 RuntimeContext 获取, 所以能用到状态的只有 RichFUnction 实例
   * c. ValueState 的 value() 方法获取当前状态
   *    update() 方法覆盖状态
   *    clear()  方法清空状态
   */
  def main(args: Array[String]): Unit = {
    // 如下例子接收端口输入的一系列数字,每两个数字计算一次平均值
    val env = StreamExecutionEnvironment.getExecutionEnvironment
    val stream = env.socketTextStream("10.19.117.19", 9999)  // nc -lk 9999
    stream.map(x => (1, if (x=="") 0 else Integer.parseInt(x)))
      .keyBy(_._1)
      .flatMap(new CountAveragePerTwoTime)
      .print()

    env.execute("average task")

    class CountAveragePerTwoTime extends RichFlatMapFunction[(Int, Int), (Int, Int)] {
      private var sum: ValueState[(Int, Int)] = _ // 用一个 valueState 记录当前的和
      override def flatMap(input: (Int, Int), out: Collector[(Int, Int)]): Unit = {
        val curSum = if (sum.value() == null) { // ValueSate为初始化时,获取值为null
          (0, 0)
        } else {
          this.sum.value()
        }
        val newSum = (curSum._1 + input._1, curSum._2 + input._2)
        this.sum.update(newSum)    // 更新状态要用 update
        if(newSum._1 >= 2){
          out.collect((input._1, newSum._2/newSum._1))  // 触发输出结果
          this.sum.clear()
        }
      }

      override def open(parameters: Configuration): Unit = {
        this.sum = this.getRuntimeContext.getState(
          new ValueStateDescriptor[(Int,Int)]("average",createTypeInformation[(Int,Int)])
        )
      }
    }
  }
}

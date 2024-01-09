package com.test.state

import org.apache.flink.api.common.functions.RichFlatMapFunction
import org.apache.flink.api.common.state.{ValueState, ValueStateDescriptor}
import org.apache.flink.configuration.Configuration
import org.apache.flink.streaming.api.functions.source.SourceFunction
import org.apache.flink.streaming.api.scala.StreamExecutionEnvironment
import org.apache.flink.streaming.api.scala._
import org.apache.flink.util.Collector

/**
 * 复杂流处理应用都是有状态的,因为要记住多个事件的信息 (比如 processWindowFunction 要记住窗口内的所有数据)
 */
object TestKeyedState {
  /**
   1. keyed state:'状态'是绑定在 KeyedStream 里某个 key 上的一个值.
      1) 如果是单个的覆盖值: ValueState
      2) 可追加的一个列表: ListState
      3) 每次来新的值执行一次 reduce 函数, 保存单个值: ReducingState
      4) 每次来新的值执行一次 aggregate 函数, 保存单个值: AggregatingState<IN, OUT>
      5) 可追加的一个 map: MapState<UK, UV>
   状态变量是保存在 RuntimeContext 中的, 所以要使用状态变量, 就要自己实现 RichFunction 接口,定义计算逻辑.
   该接口有 getRuntimeContext() 方法获取 RuntimeContext
   */
  def main(args: Array[String]): Unit = {
    // 如下例子接收端口输入的一系列数字,每两个数字计算一次平均值
    val env = StreamExecutionEnvironment.getExecutionEnvironment
    //    val stream = env.socketTextStream("10.19.117.19", 9999) // nc -lk 9999
    val source = env.addSource(new MyStatSource)

    val ds:DataStream[(String,Long)] = source.keyBy(_._1)
    // 使用状态变量
    ds.flatMap(new CountAveragePerTwoTime).print()
    env.execute("average task")

    class CountAveragePerTwoTime extends RichFlatMapFunction[(String, Long), (String, Long)] {
      /**
       * 2. 这里声明了一个单值的状态变量(ValueState). 状态不一定存储在内存中, 可以由状态后端配置
       */
      private var sum: ValueState[(Long, Long)] = _

      override def flatMap(in: (String, Long), out: Collector[(String, Long)]): Unit = {
        val curSum:(Long,Long) = if (sum.value() == null) { // ValueSate为初始化时,获取值为null
          (0L, 0L)
        } else {
          this.sum.value()
        }
        val newSum:(Long,Long) = (curSum._1 + 1L, curSum._2 + in._2)
        this.sum.update(newSum) // 更新状态值
        if (newSum._1 >= 2) {   // 逢2触发结果
          out.collect((in._1, newSum._2 / newSum._1))
          this.sum.clear()
        }
      }

      override def open(parameters: Configuration): Unit = {
        this.sum = this.getRuntimeContext.getState(
          /**
           * 3. 状态变量的声明要通过对应的 Descriptor  (ValueState -> ValueStateDescriptor)
           */
          new ValueStateDescriptor[(Long, Long)]("average", createTypeInformation[(Long, Long)])
        )
      }
    }
  }
}



class MyStatSource extends SourceFunction[(String,Long)] {
  private var running = true
  override def run(ctx: SourceFunction.SourceContext[(String,Long)]): Unit = {
    val arr = Array(
      ("a", 3L), ("b", 4L),
      ("a", 5L), ("b", 8L),
      ("c", 2L)
    )
    arr.foreach(x => {
      ctx.collect(x)
      Thread.sleep(1 * 1000)
    })

    while (running) {Thread.sleep(100 * 1000)}
  }

  override def cancel(): Unit = {
    running = false
  }
}

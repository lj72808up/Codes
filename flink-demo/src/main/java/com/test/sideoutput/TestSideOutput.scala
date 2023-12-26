package com.test.sideoutput

import org.apache.flink.streaming.api.functions.ProcessFunction
import org.apache.flink.streaming.api.scala.{DataStream, OutputTag, StreamExecutionEnvironment, createTypeInformation}
import org.apache.flink.util.Collector

/** 旁路输出: 达到一次输入分多路输出
    用 ProcessFunction, 把输出数据 collect 到不同的 tag  */
object TestSideOutput {
  final val tag1 = new OutputTag[String]("side1") {}
  final val tag2 = new OutputTag[String]("side2") {}

  def main(args: Array[String]): Unit = {
    val env = StreamExecutionEnvironment.getExecutionEnvironment
    val source = env.fromElements(1, 2, 3, 4)
    val ds = source.process(new ProcessFunction[Int, Int] {
      override def processElement(value: Int, ctx: ProcessFunction[Int, Int]#Context, out: Collector[Int]): Unit = {
        out.collect(value) //todo 常规输出
        ctx.output(tag1, s"side1-$value")  //todo 旁路输出
        ctx.output(tag2, s"side2-${2 * value}")
      }
    })

    val side1ds:DataStream[String] = ds.getSideOutput(tag1)
    val side2ds:DataStream[String] = ds.getSideOutput(tag2)

    side1ds.print("side1")
    side2ds.print("side2")

    env.execute()
  }
}

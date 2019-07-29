package com.test

import org.apache.flink.api.scala._

object TestFlink {
  def main(args: Array[String]): Unit = {
    val env = ExecutionEnvironment.getExecutionEnvironment
    val initial = env.fromElements(0)
//    initial.iterateDelta()
    val count = initial.iterate(5)({iterationInput:DataSet[Int] =>
      val res = iterationInput.map({i:Int=>{
        i+1
      }})
      res
    })

    count.print()
  }
  def main2(args: Array[String]): Unit = {
    val env = ExecutionEnvironment.getExecutionEnvironment
    val text = env.fromElements(
      "Who's there?",
      "I think I hear them. Stand, ho! Who's there?")

    val counts = text.flatMap { _.toLowerCase.split("\\W+") filter { _.nonEmpty } }
      .map { (_, 1) }
      .groupBy(0)
      .sum(1)

    counts.print()


  }
}



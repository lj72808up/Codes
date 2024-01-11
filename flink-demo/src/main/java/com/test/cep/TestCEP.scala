package com.test.cep
import org.apache.flink.cep.CEP
import org.apache.flink.cep.pattern.Pattern
import org.apache.flink.cep.pattern.conditions.SimpleCondition
import org.apache.flink.streaming.api.scala.{DataStream, StreamExecutionEnvironment, createTypeInformation}

object TestCEP {
  def main(args: Array[String]): Unit = {
    val env = StreamExecutionEnvironment.getExecutionEnvironment
    val input:DataStream[String] = env.fromElements("a1", "c", "b4", "a2", "b2", "a3")

    val pattern: Pattern[String, _] = Pattern.begin("start")   // 模式名
      // 组合条件
      .where(new SimpleCondition[String] {
        override def filter(value: String): Boolean = {
          value.startsWith("a")
        }
      })

    // 执行 Pattern
    val patternStream = CEP.pattern(input, pattern)
  }
}

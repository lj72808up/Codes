package com.test.cep

import com.fasterxml.jackson.databind.ObjectMapper
import org.apache.flink.api.common.eventtime.{SerializableTimestampAssigner, WatermarkStrategy}
import org.apache.flink.cep.functions.PatternProcessFunction
import org.apache.flink.cep.scala.{CEP, PatternStream}
import org.apache.flink.cep.scala.pattern._
import org.apache.flink.streaming.api.scala.StreamExecutionEnvironment
import org.apache.flink.streaming.api.scala._
import org.apache.flink.streaming.api.windowing.assigners.TumblingEventTimeWindows
import org.apache.flink.streaming.api.windowing.time.Time
import org.apache.flink.util.Collector

import java.text.SimpleDateFormat
import java.time.Duration
import java.util

object TestCEP {
  val df = new SimpleDateFormat("yyyyMMdd HHmmss")

  def main(args: Array[String]): Unit = {
    val env = StreamExecutionEnvironment.getExecutionEnvironment

    /**
     * 1) 包含数据的 DataStream 或 KeyedStream
     *  1. Datastream 中用于模式匹配的对象必须实现 equals() 和 hashCode() 方法, flink cep 用来比较和匹配对象
     */
    val source = env.addSource(new MySource())

    val watermarkStrategy = WatermarkStrategy.forBoundedOutOfOrderness[Event](Duration.ofSeconds(1))
      .withTimestampAssigner(new SerializableTimestampAssigner[Event]() { // 抽取时间戳的逻辑
        override def extractTimestamp(element: Event, recordTimestamp: Long): Long = {
          val date = element.dt // date: 20231212104501
          df.parse(date).getTime
        }
      })

    val ds: DataStream[Event] = source.assignTimestampsAndWatermarks(watermarkStrategy).map(x => x)
      .keyBy(_.userId)


    /**
     * 2) 定义一个 Pattern : 接受到behave是order以后，下一个动作必须是pay才算符合这个需求
     *  1. 单个的规则叫 pattern, 其组合成的规则 graph 叫 pattern sequence
     *     2. Pattern 分三类
     *     1) individual pattern
     *      a. Quantifiers 出现次数 : pattern.oneOrMore(), pattern.times(#ofTimes),
     *         pattern.times(#fromTimes, #toTimes), pattern.greedy()\
     *         b. Condition 满足条件 : pattern.where(), pattern.or()/and() or pattern.until()
     *         2)
     */
    //    val pattern1: Pattern[Event, Event] = Pattern.begin[Event]("start").where(_.behave.equals("order"))
    val pattern1: Pattern[Event, Event] = Pattern.begin[Event]("start").where(_ => true)

    /**
     * cep 可以指定2个单独策略之间以什么方式连续
     * 1) next()，指定严格连续: 期望所有匹配的事件严格的一个接一个出现，中间没有任何不匹配的事件
     * 2) followedBy()，指定松散连续: 忽略匹配的事件之间的不匹配的事件
     * 3) followedByAny()，指定不确定的松散连续: 更进一步的松散连续，允许忽略掉一些匹配事件的附加匹配
     *
     * 4) notNext()，如果不想后面直接连着一个特定事件
     * 5) notFollowedBy()，如果不想一个特定事件发生在两个事件之间的任何地方。
     */
        val patternExp = pattern1.next("next").where(_.behave.equals("pay"))
          .within(Time.seconds(3))    // 通过pattern.within()方法指定一个模式应该在10秒内发生。 这种时间模式支持处理时间和事件时间.

    /**
     * 3) 将规则应用到 DataStream 上,生成 PatternStream
     */
    val pattern: PatternStream[Event] = CEP.pattern(ds, patternExp).inEventTime()

    //记录超时的订单
    val outputTag = new OutputTag[String]("myOutput") {}

    /*/**
     * 用 select 筛出满足规则的 DataStream
     */
    val result = pattern.select(outputTag,
      new PatternTimeoutFunction[Event, String] {
        override def timeout(pattern: util.Map[String, util.List[Event]], timeoutTimestamp: Long): String = {
          val events: util.List[Event] = pattern.get("start")
          val event = events.get(0)
          "=====>"+event.toString
        }
      }, new PatternSelectFunction[Event, String] {
        override def select(pattern: util.Map[String, util.List[Event]]): String = {
          val startEvents = pattern.get("start")
          val endEvents = pattern.get("next")
          val endEvent = endEvents.get(0)
          "=====>"+endEvent.toString
        }
      })*/


    val result = pattern.process(new PatternProcessFunction[Event, String] {
      override def processMatch(`match`: util.Map[String, util.List[Event]], ctx: PatternProcessFunction.Context, out: Collector[String]): Unit = {
        val startEvent = `match`.get("start").get(0)
        val nextEvent = `match`.get("next").get(0)
        out.collect(s"find match: start:{$startEvent} =====> next:{$nextEvent}")
      }
    })

    result.print() // 规则匹配的数据
    println(result)
    val sideOutput = result.getSideOutput(outputTag)
    sideOutput.print()

    env.execute("Test CEP");
  }

  case class Event(userId: String, orderId: String, behave: String, dt: String) {
    override def toString: String = {
      s"${userId}:${orderId}:${behave} at '${dt}'"
    }
  }
}


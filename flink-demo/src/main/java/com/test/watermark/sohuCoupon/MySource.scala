package com.test.watermark.sohuCoupon

import org.apache.flink.streaming.api.functions.source.SourceFunction

import java.text.SimpleDateFormat

class MySource extends SourceFunction[String] {
  private var running = true
  // 设置为当前时间之后的1分钟.
  // 如果设置的时间比当前时间小很多, 后面的 MyTrigger 设计思路: 带超时时间的 EventTimeTrigger 里,
  // processTimeTimer 设置的会远小于机器当前时间, 当值窗口所有记录都是因为 proceeTime 超时触发计算的
  val min = new SimpleDateFormat("yyyyMMddHHmm").format(System.currentTimeMillis() + 60 * 1000)
  val key = "aa"

  override def run(ctx: SourceFunction.SourceContext[String]): Unit = {
    val arr = Array(
      s"${min}01_${key}",
      s"${min}02_${key}",
      s"${min}03_${key}",
      s"${min}04_${key}",
      // 后面这两条记录, 如果使用默认的 EventTimeTrigger, 则根本不会触发计算, 因为没有后续记录进来, 没有发送水印,
      // 无法触发 EventTimeTrigger 的 onEventTime()
      s"${min}06_${key}",
      s"${min}07_${key}"
    )
    arr.foreach(x => {
      ctx.collect(x)
      Thread.sleep(1 * 1000)
    })

    while (running) {
      Thread.sleep(100 * 1000)
    }
  }

  override def cancel(): Unit = {
    running = false
  }
}

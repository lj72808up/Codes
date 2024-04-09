package com.test.cep

import org.apache.flink.streaming.api.functions.source.{RichSourceFunction, SourceFunction}
import com.alibaba.fastjson.JSON
import com.test.cep.TestCEP.Event

class MySource extends RichSourceFunction[Event] {

  override def run(ctx: SourceFunction.SourceContext[Event]): Unit = {
    ctx.collect(Event("1","2222","order","20240315 100000"))
//    Thread.sleep(1000)

    ctx.collect(Event("1","2222","pay","20240315 100005"))
//    Thread.sleep(1000)

//    ctx.collect(Event("2","2223","pay"))
////    Thread.sleep(1000)
//
//    ctx.collect(Event("2","2224","order"))
////    Thread.sleep(4000)
//
//    ctx.collect(Event("2","2224","order"))
//    Thread.sleep(4000)
  }

  override def cancel(): Unit = {}
}
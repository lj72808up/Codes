package com.test.watermark.sohuCoupon;

import org.apache.flink.streaming.api.windowing.triggers.Trigger;
import org.apache.flink.streaming.api.windowing.triggers.TriggerResult;
import org.apache.flink.streaming.api.windowing.windows.TimeWindow;

/**
 * 什么时候调用这几个 override 方法, 判断是否应该触发呢?
 * <a href="https://blog.csdn.net/tzs_1041218129/article/details/127681305">...</a>
 * watermark 默认 200ms 发一次, 如果没有新数据打来, 最后一个窗口内的水印时间始终是一样的
 * (在 windowstrategy 定义的 onPeriodicEmit(),发送的时间戳完全相同)(可以参考 MyWatermarkStrategy#onPeriodicEmit() 中发水印的逻辑
 * , 这个就是 BoundedOutOfOrderness 的复制), 因此就无法让 onEventTime 触发窗口计算
 */
public class MyTrigger extends Trigger<Object, TimeWindow> {
    private final long processDelay ;    // 注册一个 t 的 eventTimer, 和一个 t+processDelay 的 processTimer
 
    public MyTrigger(long processDelay) {
        this.processDelay = processDelay;
    }

    @Override
    public TriggerResult onElement(Object element, long timestamp, TimeWindow window, TriggerContext ctx) throws Exception {
//        System.out.println("================ element trigger ==================");
//        System.out.println(window.maxTimestamp());
        if (window.maxTimestamp() <= ctx.getCurrentWatermark()) {
            // if the watermark is already past the window fire immediately
            return TriggerResult.FIRE;
        } else {
            ctx.registerEventTimeTimer(window.maxTimestamp());
            ctx.registerProcessingTimeTimer(window.maxTimestamp() + processDelay);
            return TriggerResult.CONTINUE;
        }
    }

    @Override
    public TriggerResult onProcessingTime(long time, TimeWindow window, TriggerContext ctx) throws Exception {
        System.out.println("==============process trigger==============");
        ctx.deleteEventTimeTimer(window.maxTimestamp());
        return TriggerResult.FIRE;
    }

    @Override
    public TriggerResult onEventTime(long time, TimeWindow window, TriggerContext ctx) throws Exception {
        System.out.println("==============event trigger==============");
        if (time == window.maxTimestamp()) {
            ctx.deleteProcessingTimeTimer(window.maxTimestamp() + processDelay);
            return TriggerResult.FIRE;
        } else {
            return TriggerResult.CONTINUE;
        }
    }

    @Override
    public void clear(TimeWindow window, TriggerContext ctx) throws Exception {
        ctx.deleteEventTimeTimer(window.maxTimestamp());
        ctx.deleteProcessingTimeTimer(window.maxTimestamp() + processDelay);
    }
}
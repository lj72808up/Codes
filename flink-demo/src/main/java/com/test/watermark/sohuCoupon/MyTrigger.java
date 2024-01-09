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

    /**
     * onElement 方法在每次数据进入该 window 时都会触发：首先判断当前的 watermark 是否已经超过了 window 的最大时间（窗口右界时间），如果已经超过，返回触发结果 FIRE；如果尚未超过，则根据窗口右界时间注册一个事件时间定时器，标记触发结果为 CONTINUE
     */
    @Override
    public TriggerResult onElement(Object element, long timestamp, TimeWindow window, TriggerContext ctx) throws Exception {
        if (window.maxTimestamp() <= ctx.getCurrentWatermark()) {
            // if the watermark is already past the window fire immediately
            return TriggerResult.FIRE;
        } else {
            /**
             * registerEventTimeTimer 的动作, 是将定时器放到了一个优先队列 eventTimeTimersQueue 中。优先队列的权重为定时器的时间戳，时间越小越希望先被触发，自然排在队列的前面。
             * onEventTime 方法是在什么时候被调用呢? 一路跟踪来到了 InternalTimerServiceImpl 类，发现只有在 advanceWatermark 的时候会触发 onEventTime，
             * 具体逻辑是：当 watermark 到来时，根据 watermark 携带的时间戳 t，从事件时间定时器队列中出队所有时间戳小于 t 的定时器，然后触发 Trigger#onEventTime()
             */
            ctx.registerEventTimeTimer(window.maxTimestamp());
            ctx.registerProcessingTimeTimer(window.maxTimestamp() + processDelay);
            return TriggerResult.CONTINUE;
        }
    }

    /** 系统根据系统时间定时地从处理时间定时器队列中取出小于当前系统时间的定时器，然后调用 onProcessingTime 方法，FIRE 窗口 */
    @Override
    public TriggerResult onProcessingTime(long time, TimeWindow window, TriggerContext ctx) throws Exception {
        System.out.println("==============process trigger==============");
        ctx.deleteEventTimeTimer(window.maxTimestamp());
        return TriggerResult.FIRE;
    }

    /**
     * onEventTime 方法的具体实现中，首先比较触发时间是否是自己当前窗口的结束时间，是则 FIRE，否则继续 CONTINUE。
     */
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
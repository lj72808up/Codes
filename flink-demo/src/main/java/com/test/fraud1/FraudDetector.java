package com.test.fraud1;

import org.apache.flink.api.common.state.ValueState;
import org.apache.flink.api.common.state.ValueStateDescriptor;
import org.apache.flink.api.common.typeinfo.Types;
import org.apache.flink.configuration.Configuration;
import org.apache.flink.streaming.api.functions.KeyedProcessFunction;
import org.apache.flink.util.Collector;
import org.apache.flink.walkthrough.common.entity.Alert;
import org.apache.flink.walkthrough.common.entity.Transaction;

public class FraudDetector extends KeyedProcessFunction<Long, Transaction, Alert> {
    private static final double SMALL_AMOUNT = 1.00;
    private static final double LARGE_AMOUNT = 500.00;
    private static final long ONE_MINUTE = 60 * 1000;

    private transient ValueState<Boolean> flagState;  // 状态标记: 一笔小额后, 接着一笔大额, 算是欺诈
    private transient ValueState<Long> timerState;    // 时间标记: 时间间隔内的2比"小额,大额",才算是欺诈
    @Override
    public void open(Configuration parameters) throws Exception {
        /** 知识点一:
         * ValueState 是 flink 中对 KeyedStream 中每个 key 状态的记录.
         *   (1) value():  获取当前状态
         *   (2) update(): 更新当前状态
         *   (3) clear():  清空当前状态
         * 由 ValueStateDescriptor 创建并维护, 其中包含 flink 如何管理变量的元数据.
         * flagDescriptor： 当一笔小交易后跟随一笔大交易就可能要报告欺诈
         */
        ValueStateDescriptor<Boolean> flagDescriptor = new ValueStateDescriptor<>(
                "flag",
                Types.BOOLEAN);
        flagState = getRuntimeContext().getState(flagDescriptor);

        ValueStateDescriptor<Long> timerDescriptor = new ValueStateDescriptor<>(
                "timer-state",
                Types.LONG);
        timerState = getRuntimeContext().getState(timerDescriptor);
    }
    @Override
    public void processElement(Transaction value
            , KeyedProcessFunction<Long, Transaction, Alert>.Context ctx
            , Collector<Alert> out) throws Exception {
        // 获取当前 key 的当前状态, 当成启动或调用 clear 后, 状态为 null
        Boolean lastTransactionWasSmall = flagState.value();

        // Check if the flag is set
        if (lastTransactionWasSmall != null) {
            if (value.getAmount() > LARGE_AMOUNT) {
                // Output an alert downstream
                Alert alert = new Alert();
                alert.setId(value.getAccountId());

                out.collect(alert);
            }

            // Clean up our state (报警后重置状态)
            flagState.clear();
        }

        /** 知识点二: 定时器
         * 设置定时器, 将在将来的某个时间点执行回调函数.
         * 这里设置定时器1分钟, 1分钟后清除状态. 表示只有小额交易后的1分钟内, 发现一笔大交易, 才触发报警
         */
        if (value.getAmount() < SMALL_AMOUNT) {
            // Set the flag to true
            flagState.update(true);

            // set the timer and timer state
            long timer = ctx.timerService().currentProcessingTime() + ONE_MINUTE;
            ctx.timerService().registerProcessingTimeTimer(timer);
            timerState.update(timer);
        }
    }

    @Override /** 触发报警 */
    public void onTimer(long timestamp, KeyedProcessFunction<Long, Transaction, Alert>.OnTimerContext ctx, Collector<Alert> out) throws Exception {
        // remove flag after 1 minute
        timerState.clear();
        flagState.clear();
    }

    private void cleanUp(Context ctx) throws Exception {
        // delete timer
        Long timer = timerState.value();
        ctx.timerService().deleteProcessingTimeTimer(timer);

        // clean up all state
        timerState.clear();
        flagState.clear();
    }
}

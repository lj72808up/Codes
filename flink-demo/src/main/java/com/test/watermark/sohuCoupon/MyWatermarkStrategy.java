package com.test.watermark.sohuCoupon;

import org.apache.flink.api.common.eventtime.*;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.time.Duration;
import java.util.Date;

import static org.apache.flink.util.Preconditions.checkArgument;
import static org.apache.flink.util.Preconditions.checkNotNull;

public class MyWatermarkStrategy implements WatermarkStrategy<String> {
    SimpleDateFormat df = new SimpleDateFormat("yyyyMMddHHmmss");
    @Override
    public TimestampAssigner<String> createTimestampAssigner(TimestampAssignerSupplier.Context context) {
        return new SerializableTimestampAssigner<String>() {
            @Override
            public long extractTimestamp(String element, long recordTimestamp) {
                String date = element.split("_")[0]; // date: 20231212104501
                try {
                    return df.parse(date).getTime();
                } catch (ParseException e) {
                    return new Date().getTime();
                }
            }
        };
    }

    @Override
    public WatermarkGenerator<String> createWatermarkGenerator(WatermarkGeneratorSupplier.Context context) {
        //周期生成水位线
//        return new MyPeriodicGenerator();
        return new MyBoundedOutOfOrdernessWatermarks<>(Duration.ofSeconds(1));

    }

    
    static class MyBoundedOutOfOrdernessWatermarks<T> implements WatermarkGenerator<T> {

        /** The maximum timestamp encountered so far. */
        private long maxTimestamp;

        /** The maximum out-of-orderness that this watermark generator assumes. */
        private final long outOfOrdernessMillis;

        /**
         * Creates a new watermark generator with the given out-of-orderness bound.
         *
         * @param maxOutOfOrderness The bound for the out-of-orderness of the event timestamps.
         */
        public MyBoundedOutOfOrdernessWatermarks(Duration maxOutOfOrderness) {
            checkNotNull(maxOutOfOrderness, "maxOutOfOrderness");
            checkArgument(!maxOutOfOrderness.isNegative(), "maxOutOfOrderness cannot be negative");

            this.outOfOrdernessMillis = maxOutOfOrderness.toMillis();

            // start so that our lowest watermark would be Long.MIN_VALUE.
            this.maxTimestamp = Long.MIN_VALUE + outOfOrdernessMillis + 1;
        }

        // ------------------------------------------------------------------------

        @Override
        public void onEvent(T event, long eventTimestamp, WatermarkOutput output) {
            maxTimestamp = Math.max(maxTimestamp, eventTimestamp);
        }

        @Override
        public void onPeriodicEmit(WatermarkOutput output) {
            output.emitWatermark(new Watermark(maxTimestamp - outOfOrdernessMillis - 1));
        }
    }


}


package org.example;

import com.lmax.disruptor.BusySpinWaitStrategy;
import com.lmax.disruptor.RingBuffer;
import com.lmax.disruptor.WaitStrategy;
import com.lmax.disruptor.dsl.Disruptor;
import com.lmax.disruptor.dsl.ProducerType;
import com.lmax.disruptor.util.DaemonThreadFactory;

import java.util.concurrent.ThreadFactory;

public class Main {
    public static void main(String[] args) {
        ThreadFactory threadFactory = DaemonThreadFactory.INSTANCE;

        WaitStrategy waitStrategy = new BusySpinWaitStrategy();
        Disruptor<ValueEvent> disruptor
                = new Disruptor<>(
                ValueEvent.EVENT_FACTORY,
                16,
                threadFactory,
                ProducerType.SINGLE,
                waitStrategy);

        SingleEventPrintConsumer consumer = new SingleEventPrintConsumer();
        disruptor.handleEventsWith(consumer.getEventHandler());

        RingBuffer<ValueEvent> ringBuffer = disruptor.start();

        for (int eventCount = 0; eventCount < 32; eventCount++) {
            long sequenceId = ringBuffer.next();
            ValueEvent valueEvent = ringBuffer.get(sequenceId);
            valueEvent.setValue(eventCount);
            ringBuffer.publish(sequenceId);
        }


    }
}

package logic;

import java.util.concurrent.atomic.AtomicInteger;

public class Counter {
    // Closed 状态进入 Open 状态的错误个数阀值
    private final int failureCount;

    // failureCount 统计时间窗口
    private final long failureTimeInterval;

    // 当前错误次数
    private final AtomicInteger currentFailCount;

    // 上一次调用失败的时间戳
    private volatile long lastTime;

    // Half-Open 状态下成功次数
    private final AtomicInteger halfOpenSuccessCount;

    public Counter(int failureCount, long failureTimeInterval) {
        this.failureCount = failureCount;
        this.failureTimeInterval = failureTimeInterval;
        this.currentFailCount = new AtomicInteger(0);
        this.halfOpenSuccessCount = new AtomicInteger(0);
        this.lastTime = System.currentTimeMillis();
    }

    /**
     * 增加错误次数, (包含两次错误次数间隔是否超过阈值的判断)
     */
    public synchronized int incrFailureCount() {
        long current = System.currentTimeMillis();
        // 超过时间窗口，当前失败次数清空, 进入新一轮错误次数判定
        // 但要记录本次错误时间, 供下次检查时间窗口
        if (current - lastTime > failureTimeInterval) {
            lastTime = current;
            currentFailCount.set(0);
        }
        // 错误次数 +1 后返回
        return currentFailCount.getAndIncrement();
    }

    public int incrSuccessHalfOpenCount() {
        return this.halfOpenSuccessCount.incrementAndGet();
    }

    /**
     * 是否到达错误次数阈值
     */
    public boolean failureThresholdReached() {
        return getCurFailCount() >= failureCount;
    }

    private int getCurFailCount() {
        return currentFailCount.get();
    }

    /**
     * 重置计数器
     */
    public synchronized void reset() {
        halfOpenSuccessCount.set(0);
        currentFailCount.set(0);
    }
}

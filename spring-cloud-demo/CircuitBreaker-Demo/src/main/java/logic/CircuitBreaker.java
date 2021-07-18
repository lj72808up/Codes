package logic;

import exception.DegradeException;
import state.State;

import java.util.function.Function;
import java.util.function.Supplier;

import static state.State.CLOSED;
import static state.State.HALF_OPEN;
import static state.State.OPEN;

public class CircuitBreaker {
    private State state;

    private Config config;

    private Counter counter;

    private long lastOpenedTime;

    /**
     * 熔断器状态默认为关闭
     */
    public CircuitBreaker(Config config) {
        this.counter = new Counter(config.getFailureCount(), config.getFailureTimeInterval());
        this.state = CLOSED;
        this.config = config;
    }

    /**
     * 阻断器执行业务代码的逻辑
     * @param toRun : 业务代码
     * @param fallback: 失败后的回调
     */
    public <T> T run(Supplier<T> toRun, Function<Throwable, T> fallback) {
        try {
            // (1) 在阻断器 open 下执行业务代码, 先判断是否超时进入 half-open
            if (state == OPEN) {
                // 判断 half-open 是否超时
                if (halfOpenTimeout()) {
                    return halfOpenHandle(toRun, fallback);
                }
                return fallback.apply(new DegradeException("degrade by circuit breaker"));
            }
            // (2) closed 状态下执行业务代码
            else if (state == CLOSED) {
                T result = toRun.get();
                closed();
                return result;
            }
            // (3) half-open 状态下执行业务代码
            else {
                return halfOpenHandle(toRun, fallback);
            }
        } catch (Exception e) {
            // (4) 代用出错后, 记录错误次数
            counter.incrFailureCount();
            if (counter.failureThresholdReached()) { // 错误次数达到阀值，进入 open 状态
                open();
            }
            return fallback.apply(e);
        }
    }

    /**
     * 进入 half-open 状态的操作:
     *  (1) half-open 状态下. 成功次数到达阀值，进入 closed 状态
     *  (2) half-open 状态下, 调用发生一次错误就进入 open 状态
     */
    private <T> T halfOpenHandle(Supplier<T> toRun, Function<Throwable, T> fallback) {
        try {
            halfOpen(); // 由 open 状态超时进入 half-open 状态
            T result = toRun.get();  // 调用业务代码
            int halfOpenSuccCount = counter.incrSuccessHalfOpenCount();
            // half-open 状态成功次数到达阀值，进入 closed 状态
            if (halfOpenSuccCount >= this.config.getHalfOpenSuccessCount()) {
                closed();
            }
            return result;
        } catch (Exception e) {
            // half-open 状态下, 调用发生一次错误就进入 open 状态
            open();
            return fallback.apply(new DegradeException("degrade by circuit breaker"));
        }
    }

    private boolean halfOpenTimeout() {
        return System.currentTimeMillis() - lastOpenedTime > config.getHalfOpenTimeout();
    }

    private void closed() {
        counter.reset();
        state = CLOSED;
    }

    /** 熔断器打开时记录时间, 下次调用通过判断时间间隔决定是否转入 half-open 状态*/
    private void open() {
        state = OPEN;
        lastOpenedTime = System.currentTimeMillis();
    }

    private void halfOpen() {
        state = HALF_OPEN;
    }
}

package logic;

/**
 * open -> half-open, 通过是否大于 halfOpenTimeout 判断
 * half-open 状态下,只要发生1次错误调用, 进化为 open
 * closed -> open: 时间窗口内超过失败次数
 * half open -> closed: 调用成功的个数超过阈值
 */
public class Config {
    // Closed 状态进入 Open 状态的错误个数阀值
    private int failureCount = 5;

    // Half-Open 状态进入 Closed 状态的成功个数阀值
    private int halfOpenSuccessCount = 2;

    // failureCount 统计时间窗口
    //   两次调用业务代码调用发生异常后, 如果时间间隔超过这个间隔, 重置失败次数为1
    private long failureTimeInterval = 2 * 1000;

    // Open 状态进入 Half-Open 状态的超时时间
    private int halfOpenTimeout = 5 * 1000;



    public int getFailureCount() {
        return failureCount;
    }

    public void setFailureCount(int failureCount) {
        this.failureCount = failureCount;
    }

    public long getFailureTimeInterval() {
        return failureTimeInterval;
    }

    public void setFailureTimeInterval(long failureTimeInterval) {
        this.failureTimeInterval = failureTimeInterval;
    }

    public int getHalfOpenTimeout() {
        return halfOpenTimeout;
    }

    public void setHalfOpenTimeout(int halfOpenTimeout) {
        this.halfOpenTimeout = halfOpenTimeout;
    }

    public int getHalfOpenSuccessCount() {
        return halfOpenSuccessCount;
    }

    public void setHalfOpenSuccessCount(int halfOpenSuccessCount) {
        this.halfOpenSuccessCount = halfOpenSuccessCount;
    }
}

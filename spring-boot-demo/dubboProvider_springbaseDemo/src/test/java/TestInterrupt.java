import java.util.concurrent.LinkedBlockingDeque;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.ThreadPoolExecutor;
import java.util.concurrent.TimeUnit;

public class TestInterrupt {
    public static void main(String[] args) throws InterruptedException {
        LinkedBlockingQueue<Runnable> jobs = new LinkedBlockingQueue<>(2);
        ThreadPoolExecutor poolExecutor = new ThreadPoolExecutor(1, 4, 10L, TimeUnit.SECONDS, jobs);
        System.out.println("未提交任务时线程池中线程数量：" + poolExecutor.getPoolSize());

        for (int i = 0; i < 4; i++) {
            poolExecutor.execute(new Runnable() {
                @Override
                public void run() {
                    try {
                        Thread.sleep(3 * 1000);
                        System.out.println(Thread.currentThread().getName() + " 执行完毕" + (System.currentTimeMillis() / 1000));
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                }
            });
            Thread.sleep(1000);
            System.out.println("第 " + i + " 个提交完毕");
        }

        System.out.println("提交任务后线程池中线程数量：" + poolExecutor.getPoolSize());

        Thread.sleep(25 * 1000);
        System.out.println("一段时间后线程池中线程数量：" + poolExecutor.getPoolSize());
    }

    public static Runnable getRunnable() {
        return new Runnable() {
            @Override
            public void run() {
                try {
                    Thread.sleep(3 * 1000);
                    System.out.println(Thread.currentThread().getName() + " 执行完毕");
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
        };
    }

    public static void main1(String[] args) {
        Thread t1 = new Thread(new Runnable() {
            @Override
            public void run() {
                // Thread.sleep(10 * 1000);  // sleep 时被 interrupt,会抛出异常
                for (; ; ) {
                }
            }
        });
        t1.start();
        t1.interrupt();
        // Thread # isInterrupted()不会清除中断标记
        System.out.println(t1.isInterrupted());
        System.out.println(t1.isInterrupted());
    }
}

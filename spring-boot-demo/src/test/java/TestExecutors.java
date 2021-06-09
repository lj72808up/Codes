import java.util.concurrent.*;

public class TestExecutors {
    public static void main(String[] args) {
        int COUNT_BITS = Integer.SIZE - 3;
        System.out.println(Integer.SIZE - 3);
        System.out.println(1 << COUNT_BITS);

        // runState is stored in the high-order bits
        System.out.println("RUNNING: " + toBinaryString(-1 << COUNT_BITS));
        System.out.println("SHUTDOWN: " + toBinaryString(0 << COUNT_BITS));
        System.out.println("STOP: " + toBinaryString(1 << COUNT_BITS));
        System.out.println("TIDYING: " + toBinaryString(2 << COUNT_BITS));
        System.out.println("TERMINATED: " + toBinaryString(3 << COUNT_BITS));

        System.out.println("CAPACITY: " + toBinaryString((1 << COUNT_BITS) - 1));

//        // 默认的 LinkedBlockingQueue 长度为 Integer.Max , 下面初始化长度为1的 LinkedBlockingQueue , 就会因为增加任务失败而执行拒绝策略
//        ExecutorService es = Executors.newFixedThreadPool(1);
        ExecutorService es = new ThreadPoolExecutor(1, 1,
                0L, TimeUnit.MILLISECONDS,
                new LinkedBlockingQueue<Runnable>(1));
        es.execute(new Runnable() {
            @Override
            public void run() {
                System.out.println(Thread.currentThread().getName());
                try {
                    Thread.sleep(2 * 1000);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
                System.out.println("finish...");
            }
        });
        es.execute(new Runnable() {
            @Override
            public void run() {
                System.out.println(Thread.currentThread().getName());
                try {
                    Thread.sleep(2 * 1000);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
                System.out.println("finish...");
            }
        });
        // 因为 poolSize = 2, 这个会阻塞, 但不会执行 RejectHandler
        es.execute(new Runnable() {
            @Override
            public void run() {
                System.out.println(Thread.currentThread().getName());
                try {
                    Thread.sleep(2 * 1000);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
                System.out.println("finish...");
            }
        });
    }

    private static String toBinaryString(int i) {

        final char[] digits = {
                '0', '1', '2', '3', '4', '5',
                '6', '7', '8', '9', 'a', 'b',
                'c', 'd', 'e', 'f', 'g', 'h',
                'i', 'j', 'k', 'l', 'm', 'n',
                'o', 'p', 'q', 'r', 's', 't',
                'u', 'v', 'w', 'x', 'y', 'z'
        };

        int shift = 1;
        char[] buf = new char[32];
        for (int j = 0; j < buf.length; j++) {
            buf[j] = '0';
        }
        int charPos = 32;
        int mask = 1;
        do {
            buf[--charPos] = digits[i & mask];
            i >>>= shift;
        } while (i != 0);

        return new String(buf);
    }
}

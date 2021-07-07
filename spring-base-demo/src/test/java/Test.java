import java.util.concurrent.locks.ReentrantLock;

public class Test {

    public static void main(String[] args) {

        Runnable r = new Runnable() {
            // 一定是公平锁才能达到交替打印的目的
            final ReentrantLock lock = new ReentrantLock(true);

            public void run() {
                for (int i = 0; i < 3; i++) {
                    lock.lock();
                    System.out.println(Thread.currentThread().getId() + "(thread)==> " + i);
                    try {
                        Thread.sleep(200);
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                    lock.unlock();
                }
            }
        };


        for (int i = 0; i < 3; i++) {
            new Thread(r).start();
        }
    }


}

class A {

}



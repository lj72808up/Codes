import exception.DegradeException;
import logic.CircuitBreaker;
import logic.Config;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;
import org.springframework.web.client.RestTemplate;

public class Main {

    private CircuitBreaker circuitBreaker;

    @Before
    public void setUp() {
        circuitBreaker = new CircuitBreaker(new Config());
    }

    @Test
    public void testScene1() {
        CircuitBreaker cb = new CircuitBreaker(new Config());
        String bookName = cb.run(() -> {
            return "deep in spring cloud";  // 正常执行, 因为熔断器初始 closed 状态
        }, t -> {
            System.err.println(t);
            return "boom";
        });
        Assert.assertEquals("deep in spring cloud", bookName);
    }

    @Test
    public void testScene2() {
        CircuitBreaker cb = new CircuitBreaker(new Config());
        RestTemplate restTemplate = new RestTemplate();
        // 调用发生 500 错误,执行回调, 显示 boom
        String result = cb.run(() -> {
            return restTemplate.getForObject("https://httpbin.org/status/500", String.class);
        }, t -> {
            System.err.println(t);
            return "boom";
        });
        System.out.println(result);
    }

    @Test
    public void testScene3() throws InterruptedException {
        Config config = new Config();
        config.setFailureCount(2); // 失败2次就会降级, 时间窗口不变
        CircuitBreaker cb = new CircuitBreaker(config);
        RestTemplate restTemplate = new RestTemplate();

        int degradeCount = 0;

        for (int index = 0; index < 10; index++) {
            // 最后一次调用暂停一段时间, 让阻断器进入 half-open, 但 half-open 下错误一次就会重新变为 open, 所以还是会触发降级
            if (index == 9) {
                Thread.sleep(config.getHalfOpenTimeout() + 1 * 1000);
            }
            String result = cb.run(() -> {
                return restTemplate.getForObject("https://httpbin.org/status/500", String.class);
            }, t -> {
                if (t instanceof DegradeException) {
                    System.out.println("degrade");
                    return "degrade";
                }
                System.out.println("boom");
                return "boom";
            });
            if (result.equals("degrade")) {
                degradeCount++;
            }
        }

        System.out.println(degradeCount);
    }

    @Test
    public void testClosedStatus() {
        // make sure in closed status
        try {
            Thread.sleep(10000l);
        } catch (InterruptedException e) {
            // ignore
        }
        String bookName = circuitBreaker.run(() -> {
            return "deep in spring cloud";
        }, t -> {
            System.err.println(t);
            return "boom";
        });
        Assert.assertEquals("deep in spring cloud", bookName);
    }

    @Test
    public void testOpenStatus() {
        // make sure in closed status
        try {
            Thread.sleep(10000l);
        } catch (InterruptedException e) {
            // ignore
        }
        for (int index = 1; index <= 10; index++) {
            int finalIndex = index;
            String bookName = circuitBreaker.run(() -> {
                throw new IllegalStateException("Oops");
            }, t -> {
                System.err.println(t);
                if (finalIndex > 5) {
                    Assert.assertTrue(t instanceof DegradeException);
                }
                return null;
            });
            if (bookName != null) {
                System.out.println(finalIndex + "-" + bookName);
            }
        }
    }

    @Test
    public void testHalfOpen2Open() {
        // make sure in closed status
        try {
            Thread.sleep(10000l);
        } catch (InterruptedException e) {
            // ignore
        }
        for (int index = 1; index <= 10; index++) {
            int finalIndex = index;
            if (finalIndex == 6) {
                try {
                    Thread.sleep(6000l);
                } catch (InterruptedException e) {
                    // ignore
                }
            }
            String bookName = circuitBreaker.run(() -> {
                if (finalIndex == 6) {
                    return "deep in spring cloud";
                }
                throw new IllegalStateException("Oops");
            }, t -> {
                System.err.println(t);
                if (finalIndex > 6) {
                    Assert.assertTrue(t instanceof DegradeException);
                }
                return null;
            });
            if (bookName != null) {
                System.out.println(finalIndex + "-" + bookName);
            }
        }
    }

    @Test
    public void testHalfOpen2Closed() {
        // make sure in closed status
        try {
            Thread.sleep(10000l);
        } catch (InterruptedException e) {
            // ignore
        }
        for (int index = 1; index <= 10; index++) {
            int finalIndex = index;
            if (finalIndex == 6) {
                try {
                    Thread.sleep(6000l);
                } catch (InterruptedException e) {
                    // ignore
                }
            }
            String bookName = circuitBreaker.run(() -> {
                if (finalIndex == 6 || finalIndex == 7) {
                    return "deep in spring cloud";
                }
                throw new IllegalStateException("Oops");
            }, t -> {
                System.err.println(t);
                if (finalIndex > 7) {
                    Assert.assertTrue(t.getMessage().equals("Oops"));
                }
                return null;
            });
            if (bookName != null) {
                System.out.println(finalIndex + "-" + bookName);
            }
        }
    }
}

package pattern.Eventbus.secondExample;
import com.google.common.eventbus.AsyncEventBus;
import com.google.common.eventbus.EventBus;
import com.google.common.eventbus.Subscribe;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.Executors;

/**
 * 观察者模式: 将观察者和被观察者解耦
 * 观察者接口
 */
interface RegObserver {
    void handleRegSuccess(long userId);
}

/**
 * 负责推送消息的观察者
 */
class RegPromotionObserver implements RegObserver {
    private PromotionService promotionService; // 依赖注入

    @Override
    public void handleRegSuccess(long userId) {
        promotionService.issueNewUserExperienceCash(userId);
    }
}

/**
 * 负责通知用户的观察者
 */
class RegNotificationObserver implements RegObserver {
    private NotificationService notificationService;

    @Override
    public void handleRegSuccess(long userId) {
        notificationService.sendInboxMessage(userId, "Welcome...");
    }
}

class UserController {
    private UserService userService; // 依赖注入
    private List<RegObserver> regObservers = new ArrayList<>();

    // 一次性设置好，之后也不可能动态的修改
    public void setRegObservers(List<RegObserver> observers) {
        regObservers.addAll(observers);
    }

    public Long register(String telephone, String password) {
        //省略输入参数的校验代码
        //省略userService.register()异常的try-catch代码
        long userId = userService.register(telephone, password);  // 真实的注册逻辑

        //todo 下面是同步阻塞方式实现观察者和非观察者, 如果在线程池中执行, 就是异步非阻塞模式
        for (RegObserver observer : regObservers) {
            observer.handleRegSuccess(userId);
        }

        return userId;
    }
}




interface PromotionService {
    void issueNewUserExperienceCash(long userId);
}

interface NotificationService{
    void sendInboxMessage(long userId, String msg);
}

interface UserService {
    long register(String telephone, String password);
}


/**
 * 调用者
 */
class UserControllerGuava {
    private UserService userService; // 依赖注入

    private EventBus eventBus;
    private static final int DEFAULT_EVENTBUS_THREAD_POOL_SIZE = 20;

    public UserControllerGuava() {
        //eventBus = new EventBus(); // 同步阻塞模式
        eventBus = new AsyncEventBus(Executors.newFixedThreadPool(DEFAULT_EVENTBUS_THREAD_POOL_SIZE)); // 异步非阻塞模式
    }

    public void setRegObservers(List<Object> observers) {
        for (Object observer : observers) {
            eventBus.register(observer);
        }
    }

    public Long register(String telephone, String password) {
        //省略输入参数的校验代码
        //省略userService.register()异常的try-catch代码
        long userId = userService.register(telephone, password);

        eventBus.post(userId);

        return userId;
    }
}

/**
 * 观察者一
 */
class RegPromotionObserverGuava {
    private PromotionService promotionService; // 依赖注入

    @Subscribe
    public void handleRegSuccess(Long userId) {
        promotionService.issueNewUserExperienceCash(userId);
    }
}

/**
 * 观察者二
 */
class RegNotificationObserverGuava {
    private NotificationService notificationService;

    @Subscribe
    public void handleRegSuccess(Long userId) {
        notificationService.sendInboxMessage(userId, "...");
    }
}
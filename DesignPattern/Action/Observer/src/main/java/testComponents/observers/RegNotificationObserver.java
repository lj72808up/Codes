package testComponents.observers;

import components.Subscribe;

public class RegNotificationObserver {
//    private NotificationService notificationService;

    @Subscribe
    public void handleRegSuccess(Long userId) throws InterruptedException {
        Thread.sleep(2000);
        System.out.println("[" + Thread.currentThread().getName() + "] " + this.getClass().getName() + " : " + userId);
//        notificationService.sendInboxMessage(userId, "...");
    }
}
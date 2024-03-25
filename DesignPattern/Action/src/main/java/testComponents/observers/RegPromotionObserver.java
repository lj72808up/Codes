package testComponents.observers;

import components.Subscribe;

public class RegPromotionObserver {
//    private PromotionService promotionService; // 依赖注入

    @Subscribe
    public void handleRegSuccess(Long userId) throws InterruptedException {
        Thread.sleep(2000);
        System.out.println("["+Thread.currentThread().getName()+"] " +this.getClass().getName() + " : " + userId);
//        promotionService.issueNewUserExperienceCash(userId);
    }
}
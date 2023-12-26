Example
-----------
1. `com.test.watermark.sohuCoupon.CouponTest`  
    本案例是写流量群曝光累加时产生灵感, 抽象出的编程范式. 解决一个场景:
    "使用时间时间时最后一个窗口因没有后续数据到来, 导致 eventTime 水位线不往前推进, 无法触发窗口计算
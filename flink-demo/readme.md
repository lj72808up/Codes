Example
-----------
1. `com.test.watermark.sohuCoupon.CouponTest`  
    本案例是写流量群曝光累加时产生灵感, 抽象出的编程范式. 解决一个场景:
    "使用时间时间时最后一个窗口因没有后续数据到来, 导致 eventTime 水位线不往前推进, 无法触发窗口计算. 虽然这里用自定义 trigger 的方式, 让 processtime scheduler 在真实时间超时后触发窗口, 但这种方案在回溯数据时会导致触发混乱. 因为回溯的数据全都算 processTime 的超时数据, 导致 eventTime 在同一个5分钟内的数据疯狂多次触发.
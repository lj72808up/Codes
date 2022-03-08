object TestShare {
  // 以3月8日上证点数为基准, 3293 点
  // 1) 单个品种控制好比例, 不要一个低估就玩命投, 最多投到 1/5 就不投了
  def main(args: Array[String]): Unit = {
    val initShare = 3000d // 初始本金
    var curShare = initShare
    var sumShare = initShare

    val everyRates = 2   //(%) 每次跌 10 %
    val turnNum = 5 // 下跌轮数
    val newInRate = 0.3 //(倍数) 每次亏损后投入的现金倍数

    for (i <- 1 to turnNum) {
      val deRate = everyRates / 100d       // 本次跌幅 (换算成小数计算)
      // 每次下跌, 投入下跌前的本金 * 下跌率
      val newInCash = newInRate * curShare // 本次投入的金额
      sumShare = sumShare + newInCash
      curShare = curShare * (1 - deRate) + newInCash // 下跌后所剩的本金
      println("本次下跌:" + printPercent(deRate) + "; 本次亏损:" + (curShare * deRate).formatted("%.2f") + "; 本次新投入: " + (newInCash).formatted("%.2f"))

      println("累计投入:" + sumShare.formatted("%.2f") + "; 还剩金额:" + curShare.formatted("%.2f"))
      println("总亏损率:" + printPercent(1 - (curShare / sumShare)))
      println("========================================================")
    }

    println("不补仓的亏损率: "+ printPercent(1-Math.pow(1-everyRates/100d,turnNum)) +" ("+printPercent(turnNum*everyRates/100d)+")")
    println("结束状态为, 大跌之后补仓, 不计算后续跌幅的")
  }

  def printPercent(rate: Double): String = {
    (rate * 100).formatted("%.2f") + "%"
  }
}

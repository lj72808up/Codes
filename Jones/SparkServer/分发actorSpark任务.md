actor 投递本地 spark 任务的策略:   

1. 优化原则:    
  认为已经执行了很长时间的任务, 其继续执行直到完毕所需要的时间更长     
  对于actor中没有执行的spark任务, 预测其执行完毕的时间是5(min) , 一个估计值     

2. 筛选标准:   
   选择权重最小的actor进行投递   

3. actor 打分策略
    (1) 当前没有任务的 actor
         权重 = 0
    (2) 正在处里任务的 actor 
         权重 += 当前任务已经执行的秒数
    (3) 如果当前 actor 的 mailbox 里有待处理的任务
         权重 += 待处理任务个数 * 5(min) 

4. 任务最长等待时间的优化效果    
    (1) 优化前: 2.7 小时 (发生了任务之间互相等待的情况)     
    (2) 优化后: 1分钟 

```scala
this.jdbcRouter = context.actorOf(RoundRobinPool(jobQueueNum).props(Props(classOf[JdbcActor],
      sparkSession, dbLog)), "jdbcRouter")

this.sparkRouter = context.actorOf(IdleFirstPool(jobQueueNum).props(Props(classOf[SparkCalcActor],
      sparkSession, dbLog)), "sparkRouter")
```

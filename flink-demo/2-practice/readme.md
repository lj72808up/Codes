#### 一. 本章重点介绍状态管理和时间
1. 自定义时间流处理  
对于时间，我们关注的是事件的发生顺序，而不是事件被传输或处理的顺序。因此我们会使用记录在事件中的时间戳，而非处理事件的机器时钟的时间。
2. 有状态流处理  
flink 的算子是有状态的，即如何处理事件可能取决于之前所有事件数据的积累结果。这个“之前处理的事件数据结果“就是状态。
flink 算子的状态访问都是在本地进行，有助于提高吞吐量，降低延迟。状态会存储在 JVM 堆上,但如果状态太大， 可以选择将其存储在高速磁盘中。
3. 容错   
flink通过状态快照(state snapshot)和流重放(stream replay)实现 exactly-once 语义。

#### 二. 测试的 DataStream 生成方法
1. java pojo 对象  
cantest.ExampleDataStream 示例

#### 三. EventTime 和 WaterMarks
1. flink 中有三种时间 
   1. 事件时间（event time）：事件发生的时间
   2. 摄取时间（ingestion time）：flink 收到事件的时间
   3. 处理时间（processing time）：flink pipeline中具体算子处理事件的时间
  想要流重放得到相同的结果，必须使用事件时间
2. 窗口
  窗口往往是针对不同的 key 分组。
    ```
    stream.
    .keyBy(<key selector>)
    .window(<window assigner>)
    .reduce|aggregate|process(<window function>)
    ```
 窗口分配器可以定义不同类型的窗口，类型包括三种  
    1. 滚动时间窗口  
     滚动窗口模式下窗口之间不重叠，且窗口长度Size是固定的。
     滚动窗口和滑动窗口都有一个对齐参数offset，默认是0， 比如时长一小时的滚动窗口，你会得到如 1:00:00.000 - 1:59:59.999、2:00:00.000 - 2:59:59.999 等窗口
     每分钟页面浏览量
    ```
    TumblingEventTimeWindows.of(Time.minutes(1))
    ```
    2. 滑动时间窗口
     滑动窗口以一个步长（Slide）不断向前滑动，窗口的长度固定。
     每10秒钟计算前1分钟的页面浏览量
    ```
    SlidingEventTimeWindows.of(Time.minutes(1), Time.seconds(10))
    ```
    3. 会话窗口
     会话窗口模式下，两个窗口之间有一个间隙，叫做 session gap。当窗口在大于 session gap 的时间内没有收到事件，则窗口关闭。下次再来事件开启新的窗口
     每个会话的网页浏览量，其中会话之间的间隔至少为30分钟。
     会话窗口的实现原理是：窗口之间可以聚合。每个事件在初始被消费时，都会分配一个新的窗口，如果窗口之间的时间间隔小于 session gap，则两个窗口会聚合
    ```
    EventTimeSessionWindows.withGap(Time.minutes(30))
    ```
3. 窗口应用函数
 应用在窗口上的算子，就叫窗口应用函数。在窗口上执行聚合，分为全量聚合和增量聚合
   1. 全量聚合`process(ProcessWindowFunction)`     
    全量聚合会缓存所有分配给窗口的事件流，直到窗口结束。这个操作可能相当昂贵
    ```
    DataStream<SensorReading> input = ...
    
    input
    .keyBy(x -> x.key)
    .window(TumblingEventTimeWindows.of(Time.minutes(1)))
    .process(new MyWastefulMax());
    
    public static class MyWastefulMax extends ProcessWindowFunction<
                   SensorReading,                  // 输入类型
                   Tuple3<String, Long, Integer>,  // 输出类型
                   String,                         // 键类型
                   TimeWindow> {                   // 窗口类型 
        @Override  // 关键处理逻辑
        public void process(
                String key,
                Context context,
                Iterable<SensorReading> events,
                Collector<Tuple3<String, Long, Integer>> out) {
    
            int max = 0;
            for (SensorReading event : events) {
                max = Math.max(event.value, max);
            }
            out.collect(Tuple3.of(key, context.window().getEnd(), max));
        }
    }
   ```
   2. 增量聚合  
      窗口当中每加入一条数据，就进行一次统计
      ``` reduce(reduceFunction)```
      ``` aggregate(aggregateFunction) ```
      ``` sum(),min(),max() ```

    ```
    DataStream<SensorReading> input = ...
    
    input
        .keyBy(x -> x.key)
        .window(TumblingEventTimeWindows.of(Time.minutes(1)))
        .reduce(new MyReducingMax(), new MyWindowFunction());
    
    private static class MyReducingMax implements ReduceFunction<SensorReading> {
        public SensorReading reduce(SensorReading r1, SensorReading r2) {
            return r1.value() > r2.value() ? r1 : r2;
        }
    }
    
    private static class MyWindowFunction extends ProcessWindowFunction<
        SensorReading, Tuple3<String, Long, SensorReading>, String, TimeWindow> {
    
        @Override
        public void process(
                String key,
                Context context,
                Iterable<SensorReading> maxReading,
                Collector<Tuple3<String, Long, SensorReading>> out) {
    
            SensorReading max = maxReading.iterator().next();  // 只读取
            out.collect(Tuple3.of(key, context.window().getEnd(), max));
        }
    }
    ```
    Iterable<SensorReading> 将只包含一个读数 – MyReducingMax 计算出的预先汇总的最大值。
    [例子](https://blog.csdn.net/hyy1568786/article/details/106339463?utm_medium=distribute.pc_aggpage_search_result.none-task-blog-2~aggregatepage~first_rank_ecpm_v1~rank_v31_ecpm-2-106339463.pc_agg_new_rank&utm_term=flink%E5%A2%9E%E9%87%8F%E5%92%8C%E5%85%A8%E9%87%8F%E8%AE%A1%E7%AE%97&spm=1000.2123.3001.4430)

4. 数据管道和ETL
   [练习](https://github.com/apache/flink-training/tree/master/rides-and-fares)
5. 事件驱动应用
   [练习](https://github.com/apache/flink-training/tree/master/long-ride-alerts)
6. 
package com.test.source

/**
 * 所有 source 的统一概念：
 * 一个数据 source 包括三个核心组件
  （1）分片（split）：    split 是 source 进行任务分配和数据并行读取的基本粒度。
  （2）SourceReader：    用来读取分配给子的 split。SourceReader
        该组件在 TaskManagers 上的 SourceOperators 并行运行，并产生并行的事件流/记录流。
  （3）SplitEnumerator： 用来划分 split 给 SourceReader。
        该组件在 JobManager 上以单并行度运行，负责对未分配的 split 进行维护，并以均衡的方式将其分配给 reader。
 */
object TestSourceCommon {
  /** 对于 KafkaSource：
   * split：kafka 的一个分区
   * SplitEnumerator： list all topic partitions involved in the subscribed topics. 这个可以配置成定期发现新 partition
   * SourceReader ： 负责读取每个 parititon 的数据
   */
  object T1{}
}

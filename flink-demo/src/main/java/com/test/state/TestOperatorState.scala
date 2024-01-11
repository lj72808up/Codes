package com.test.state

import org.apache.flink.api.common.state.{ListState, ListStateDescriptor}
import org.apache.flink.api.common.typeinfo.{TypeHint, TypeInformation}
import org.apache.flink.runtime.state.{FunctionInitializationContext, FunctionSnapshotContext}
import org.apache.flink.streaming.api.checkpoint.CheckpointedFunction
import org.apache.flink.streaming.api.functions.sink.SinkFunction
import org.apache.flink.streaming.api.functions.sink.SinkFunction.Context

import scala.collection.JavaConverters._
import scala.collection.mutable.ListBuffer

/**
 * 算子状态在 checkpoint 时使用, 用 CheckpointedFunction 接口来使用 operator state。
 * 算子状态是没有绑定到 key 的状态. 常用语 source 和 sink. 算子状态是一个 list, 之所以设计成 list , 是为了在改变并发后进行状态的重新分派， 重新分配的方式2种:
 * Even-split redistribution: 每个算子都保存一个列表形式的状态集合，整个状态由所有的列表拼接而成。当作业恢复或重新分配的时候，整个状态会按照算子的并发度进行均匀分配。
 *                            比如说，算子 A 的并发读为 1，包含两个元素 element1 和 element2，当并发读增加为 2 时，element1 会被分到并发 0 上，element2 则会被分到并发 1 上。
 * Union redistribution: 每个算子保存一个列表形式的状态集合。整个状态由所有的列表拼接而成。当作业恢复或重新分配时，每个算子都将获得所有的状态数据。
 *                       Do not use this feature if your list may have high cardinality.
 *                       Checkpoint metadata will store an offset to each list entry, which could lead to RPC framesize or out-of-memory errors.
 */
class TestOperatorState {

}


//info  如下示例: 这个SinkFunction 在 CheckpointedFunction 中进行数据缓存, 然后统一发送到下游. 算子状态使用 Even-split redistribution 模式保存
class BufferingSink(threshold: Int = 0) extends SinkFunction[(String, Int)] with CheckpointedFunction {
  //todo FlinkKafkaConsumer(source) 也是用 ListState 存储的 offset 状态,为什么声明一个变量, snapshotState 方法中只要操作这个变量, 就能对状态 checkpoint? 恐怕要看 checkpoint 源码才行
  /**
   * info: 要在 checkpoint 时保存算子状态, 就要有2个变量: 1个是 ListState, 另一个是真正的数据 list.
   *   initState 初始化时, 从 ListState 读出数据 list; snapshotState 把数据list 写入到 ListState 中
   *   source 是 kafka 时, FLinkafkaCOnsumer 就是这2个变量交替使用, 对分区 offset 快照
   *    (1) private transient volatile TreeMap<KafkaTopicPartition, Long> restoredState;
   *    (2) private transient ListState<Tuple2<KafkaTopicPartition, Long>> unionOffsetStates
   */
  private var checkpointedState: ListState[(String, Int)] = _
  private val bufferedElements = ListBuffer[(String, Int)]()

  //info: Writes the given value to the sink. This function is called for every record.
  override def invoke(value: (String, Int), context: Context): Unit = {
    bufferedElements += value
    if (bufferedElements.size >= threshold) {
      for (element <- bufferedElements) {
        // send it to the sink
      }
      bufferedElements.clear()
    }
  }

  /**
   * 进行 checkpoint 时会调用 snapshotState()
   */
  override def snapshotState(context: FunctionSnapshotContext): Unit = {
    checkpointedState.clear()
    for (element <- bufferedElements) {
      checkpointedState.add(element)
    }
  }

  /**
   * 用户自定义函数初始化时会调用 initializeState()，该方法要包含最初的初始化和从 checkpoint 恢复时初始化。
   * 因此 initializeState() 不仅是定义不同状态类型初始化的地方，还需要包括状态恢复的逻辑
   */
  override def initializeState(context: FunctionInitializationContext): Unit = {
    val descriptor = new ListStateDescriptor[(String, Int)](
      "buffered-elements",
      TypeInformation.of(new TypeHint[(String, Int)]() {})
    )

    checkpointedState = context.getOperatorStateStore.getListState(descriptor)

    if(context.isRestored) {
      for(element <- checkpointedState.get().asScala) {
        bufferedElements += element
      }
    }
  }

}

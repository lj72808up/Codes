import org.apache.spark._
import org.apache.spark.streaming._
import org.apache.spark.streaming.StreamingContext._
import org.apache.spark.streaming.dstream.DStream

import java.text.SimpleDateFormat

object StreamingTest {

  def testDF(words: DStream[String]): Unit = {
    /* words.foreachRDD { rdd =>
     // Get the singleton instance of SparkSession
      val spark = SparkSession.builder.config(rdd.sparkContext.getConf).getOrCreate()
      import spark.implicits._

      // Convert RDD[String] to DataFrame
      val wordsDataFrame = rdd.toDF("word")

      // Create a temporary view
      wordsDataFrame.createOrReplaceTempView("words")

      // Do word count on DataFrame using SQL and print it
      val wordCountsDataFrame =
        spark.sql("select word, count(*) as total from words group by word")
      wordCountsDataFrame.show()
    }*/
  }


  def main(args: Array[String]): Unit = {
    val conf = new SparkConf().setMaster("local[2]").setAppName("NetworkWordCount")
    // sparkStreaming 会自动从整点对齐时间，这里 batch interval 2分钟。
    // 如果10点21分开启的程序，10点22分第一次打印rdd结果；10点24分，第二次打印结果
    val ssc = new StreamingContext(conf, Minutes(2))

    /** lines 作为 Input DStreams, 通过 Receiver 对象从源头收集的数据, 然后存储在内存中
      * //  (1) 如果想并行的接受数据流, 可以创建多个 DStream
      * //  (2) 处理数据的 core 要比 receiver 的 core 多, 否则会来不及处理接受的数据
      * //  (3) 除了网络端口的 DStream, spark 还支持监听目录生成 DStream. 目录的监听规则是根据他们的修改时间
      */

    // 自定义 Reciever http://spark.apache.org/docs/latest/streaming-custom-receivers.html
    val lines = ssc.socketTextStream("localhost", 9999)

    // Split each line into words
    val words = lines.flatMap(_.split(" "))

    // Count each word in each batch
    val pairs = words.map(word => (word, 1))

    /** DStream 中, 一些有特色的算子
      * (1) UpdateStateByKey:  窗口累加效果的算子 (testUpdateStateByKey) , 需指定 checkPoint 目录
      * (2) Window Operations: 将几个时间窗口内的 RDDs 们合并产生新的 RDD, 并可以指定每次合并划过的窗口数 (testCombineAndSlide)
      * (3) join 操作
      *     1) DStream 可以和 DStream join
      *     2) DStream 可以和 DataSet join
      * (4) 输出操作 foreachRdd(func)
      *  foreachRdd() 操作在 driver 节点上执行
      *  最通用的输出为 foreachRdd(func),foreachRdd 最特殊的在于可以将 RDD 发到外部系统.但要注意几种错误写法
      *     a. 在 driver 节点上创建了 connection 对象, 因为不能被序列化, 所以无法从 driver 发送到 worker
      *        dstream.foreachRDD { rdd =>
      *        val connection = createNewConnection()  // executed at the driver
      *        rdd.foreach { record =>
      *        connection.send(record) // executed at the worker
      *        }
      *        }
      *     b. 为每条记录创建一个 Connection (也是错误的)
      *        dstream.foreachRDD { rdd =>
      *        rdd.foreach { record =>
      *        val connection = createNewConnection()
      *        connection.send(record)
      *        connection.close()
      *        }
      *        }
      *     c. 正确的做法是, 让 RDD 的每个 partition 公用一个 Connection 对象
      *        dstream.foreachRDD { rdd =>
      *        rdd.foreachPartition { partitionOfRecords =>
      *        val connection = createNewConnection()
      *        partitionOfRecords.foreach(record => connection.send(record))
      *        connection.close()
      *        }
      *        }
      *     d. 最终优化: 让每个 executor 上的 rdd partition 从连接池获取连接, 达到连接共用效果
      *        dstream.foreachRDD { rdd =>
      *        rdd.foreachPartition { partitionOfRecords =>
      *        // ConnectionPool is a static, lazily initialized pool of connections
      *        val connection = ConnectionPool.getConnection()
      *        partitionOfRecords.foreach(record => connection.send(record))
      *        ConnectionPool.returnConnection(connection)  // return to the pool for future reuse
      *        }
      *        }
      *
      * (5) DStream 的 foreachRDD() 方法可以操作每个 batch 下的 rdd, 也就可以在方法中使用 DataFrame 的 api
      */
    val wordCounts = pairs.reduceByKey(_ + _)
    val format = new SimpleDateFormat("yyyyMMdd HH:mm:ss")
    // Print the first ten elements of each RDD generated in this DStream to the console
    wordCounts.foreachRDD((rdd,time)=>{
      println(format.format(time.milliseconds)+"=>"+rdd.collect().mkString(","))
    })

    ssc.start() // Start the computation
    ssc.awaitTermination() // Wait for the computation to terminate

    // nc -lp 9999
    /**
      * (1) 每次的结果都是新的 miniBatch 在 5s 时间内统计的新值, 不累加
      */
  }


  // UpdateStateByKey 记住状态的
  def testUpdateStateByKey(args: Array[String]): Unit = {
    val conf = new SparkConf().setMaster("local[2]").setAppName("NetworkWordCount")
    val ssc = new StreamingContext(conf, Seconds(5))
    ssc.checkpoint("mycheck") // updateStateByKey 算子需要指定 checkpoint

    val lines = ssc.socketTextStream("localhost", 9999)
    // 每行可能用空格区分单词
    val words = lines.flatMap(_.split(" "))
    // 生成 (word,1) 的对儿
    val pairs = words.map(word => (word, 1))
    // 累积效果的算子
    val accumulateCnt = pairs.updateStateByKey[Int](
      // newValues: 新 batch 下, 同一个 key 分区后的 value 列表
      // oldCount: 上 batch 中是否有对应 key 的 value
      (newValues: Seq[Int], oldCount: Option[Int]) => Option[Int] {
        val res = newValues.sum + oldCount.getOrElse(0)
        res
      })
    accumulateCnt.print()

    ssc.start()
    ssc.awaitTermination()
  }

  def testCombineAndSlide(args: Array[String]): Unit = {
    val conf = new SparkConf().setMaster("local[2]").setAppName("NetworkWordCount")
    val ssc = new StreamingContext(conf, Seconds(5))
    //    ssc.checkpoint("mycheck") // updateStateByKey 算子需要指定 checkpoint

    val lines = ssc.socketTextStream("localhost", 9999)
    // 每行可能用空格区分单词
    val words = lines.flatMap(_.split(" "))
    // 生成 (word,1) 的对儿
    val pairs = words.map(word => (word, 1))

    val windowRes = pairs.reduceByKeyAndWindow(
      (a: Int, b: Int) => {
        a + b
      }, // reduceFunc
      Seconds(15), // combine duration : 每次执行组和多少 Rdds
      Seconds(10)) // slide duration : 每隔长时间执行一次

    windowRes.print()

    ssc.start()
    ssc.awaitTermination()
  }
}


/**
  * Kafka 与 Streaming 集成
  * (1) kafka 分区和 rdd partition 一一对应, 但这种对应关系在 shuffle 操作后破坏
  * (2) LocationStrategies: spark executor 和 kafka partition 的映射策略
  * (3) streaming 任务需要在幂等输出后存储 offset, 方法:
  *   a) rdd checkpoint (X) : 该操作会附带 rdd 的 offset,
  *   b) 用 kafka 的 commitOffset api 提交, offset 保存在 kafka 的队列中
  */


/**
 * 如何确定 batch interval 的时间 ?
 * 1) 根据业务方，看要求的时间批次
 * 2) 观察Spark UI上的batch处理时间来定。batch处理时间必须小于batch interval时间。
 */

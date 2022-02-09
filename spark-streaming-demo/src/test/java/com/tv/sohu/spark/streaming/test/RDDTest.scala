package com.tv.sohu.spark.streaming.test

import com.tv.sohu.spark.streaming.dm4.accumulate.AccConf
import com.tv.sohu.spark.util.Utils
import org.apache.hadoop.conf.Configuration
import org.apache.hadoop.fs.{FileSystem, Path}
import org.apache.spark.rdd.{CheckpointFileRDD, RDD}
import org.apache.spark.{SparkConf, SparkContext}

import scala.reflect.ClassTag

/**
  * Created by wentaoli213587 on 2019/5/30.
  */
object RDDTest {

  var sc: SparkContext = _

  def main(args: Array[String]): Unit = {
     sc = new SparkContext(new SparkConf().
      setAppName("Realtime_Accumulate").
      set("spark.serializer", "org.apache.spark.serializer.KryoSerializer").
      set("spark.memory.useLegacyMode", "true").
      set("spark.storage.memoryFraction", "0.2").
      set("spark.shuffle.memoryFraction", "0.6")
      .setMaster("local[5]")
    )

    val checkpointDir = "/user/wentaoli213587/spark/checkpoint"

    val Array(a, b, c) = Utils.parseKeyToArray("t\tb\t", '\t')

//    sc.setCheckpointDir(checkpointDir)

    Thread.sleep(10000)

    val rdd = sc.parallelize(List("a", "b", "c"))
    println(rdd.repartition(10).partitioner)
    rdd.count()

    val rdd2 = sc.parallelize(List("a", "b", "c", "9990pppppppppp"))
//    rdd2.checkpoint()
//    rdd2.count()

    println(getLatestCheckpoint(AccConf.ACC_CHECKPOINT_PATH))
    val cpRDD = readCheckpoint[(String, _)](AccConf.ACC_CHECKPOINT_PATH)
    println(cpRDD.partitioner+"............")

  }

  private def readCheckpoint[T: ClassTag](checkpointDir: String): RDD[T] = {
    val latestPath = getLatestCheckpoint(checkpointDir)
    if(latestPath.isDefined)
      new CheckpointFileRDD[T](sc, latestPath.get)
    else
      sc.emptyRDD[T]
  }

  //e.g. /user/dmspark/accumulate/checkpoint/0519936a-5bff-4ecf-a6f0-3854e5952ec9/rdd-689
  def getLatestCheckpoint(checkpointDir: String): Option[String] = {
    val fs = FileSystem.get(new Configuration())
    def findLatestSubDir(path: Path): Path = {
      if(path == null || !fs.exists(path))
        return null
      val fileStatus = fs.listStatus(path).sortBy(_.getModificationTime).reverse
      if(fileStatus.isEmpty) null else fileStatus(0).getPath
    }
    val latestPath = findLatestSubDir(findLatestSubDir(new Path(checkpointDir)))
    if(latestPath != null) Some(latestPath.toString) else None
  }

}

package org.test.cc

import org.apache.log4j.LogManager
import org.apache.spark.sql.{DataFrame, SparkSession}
import org.apache.spark.sql.functions._
import org.apache.spark.storage.StorageLevel
import org.test.ParseUtil
import org.slf4j.LoggerFactory
import org.test.conf.Constant


/**
  * DataFrame实现的连通图算法, 对ConnectedComponentRDD做出的改进
  */
object AlternatingAlgo {

  def main(args: Array[String]): Unit = {

    val log = LogManager.getRootLogger

    val saveTime = ParseUtil.parseSaveTime(args)
    val spark = SparkSession
      .builder()
      .appName("cc")
      .config("spark.sql.hive.convertMetastoreOrc", true)
      .config("spark.sql.orc.enabled", true)
      .config("spark.speculation",true)
      .enableHiveSupport()
      .getOrCreate()

    spark.sparkContext.setLogLevel("WARN")
    import spark.implicits._
    import org.apache.spark.sql.functions._
    import spark.sql

    val fields = Array("t1_primarykey","t2_primarykey")

    var df = sql(s"select ${fields.mkString(",")} from ${Constant.ccSrcTable} where prob>${Constant.threshold}")
    df = df.select(col("t1_primarykey").as("u"),col("t2_primarykey").as("v"))

    def largeStar(df:DataFrame):DataFrame={

      val ds = df.as[(String,String)]
      // 1. 发射
      val emitDF = ds.flatMap(x=>{
        if(x._1==x._2){
          List((x._1,x._2))
        }else{
          List((x._1,x._2),(x._2,x._1))
        }
      }).toDF("u","v")

      // 2.
      /**
        * 包含自己在内的最小邻居节点
        * @param df
        * @param groupField: 分组字段
        * @param aggField: 取那个字段的最小值
        * @param aggName: 最小值列的别名
        * @return: DataFrame("groupField","aggName")
        */
      def getMinNeighborIncludeSelf(df:DataFrame,groupField:String,aggField:String,aggName:String):DataFrame = {
        val minNeighbor = df.groupBy(col(groupField)).agg(min(aggField).as(aggName))  // columns: |sha1|minNeighbor|
        // used by largeStar
        val minNeighborAndSelf = minNeighbor.select(col(groupField),
          when(col(groupField)<col(aggName),col(groupField))
            .otherwise(col(aggName))
            .as(aggName)
        )
        minNeighborAndSelf   //groupField,aggName
      }
      val minNeighborAndSelfDF = getMinNeighborIncludeSelf(emitDF,"u","v","minNeighbor")

      // 3. join
      val joinRes = emitDF.join(minNeighborAndSelfDF,emitDF("u")===minNeighborAndSelfDF("u"),"left")
        .filter($"v">=$"minNeighbor")
        .select($"v".as("u"),$"minNeighbor".as("v"))
        .union(minNeighborAndSelfDF.select($"u",$"minNeighbor".as("v")))
      // 4. 去重
      val largeRes = joinRes.groupBy($"u",$"v").agg($"u",$"v").select($"u",$"v")
      //      largeRes.persist(StorageLevel.MEMORY_AND_DISK)
      largeRes
    }

    def smallStar(df:DataFrame):DataFrame={
      import spark.implicits._
      val ds = df.as[(String,String)]
      // 1. 发射
      val emitDF = ds.map(x=>{
        if(x._1>x._2){
          (x._1,x._2)
        }else{
          (x._2,x._1)
        }
      }).toDF("u","v")

      // 2. 得到最小
      def getMinNeighborIncludeSelf(df:DataFrame,groupField:String,aggField:String,aggName:String):DataFrame = {
        val minNeighbor = df.groupBy(col(groupField)).agg(min(aggField).as(aggName))  // columns: |sha1|minNeighbor|
        // used by largeStar
        val minNeighborAndSelf = minNeighbor.select(col(groupField),
          when(col(groupField)<col(aggName),col(groupField))
            .otherwise(col(aggName))
            .as(aggName)
        )
        minNeighborAndSelf   //groupField,aggName
      }
      val minNeighborAndSelfDF = getMinNeighborIncludeSelf(emitDF,"u","v","minNeighbor")

      // 3. join
      val joinRes = emitDF.join(minNeighborAndSelfDF,emitDF("u")===minNeighborAndSelfDF("u"),"left")
        .select($"v".as("u"),$"minNeighbor".as("v"))
        .union(minNeighborAndSelfDF.select($"u",$"minNeighbor".as("v")))

      val smallRes = joinRes.groupBy($"u",$"v").agg($"u",$"v").select($"u",$"v")
      smallRes
    }

    /**
      *
      * @param df1
      * @param df2
      * @return has diff?
      */
    def diff(df1:DataFrame,df2:DataFrame):Boolean = {
      val count = df1.join(df2,df1("u")===df2("u")and(df1("v")===df2("v")),"full")
        .filter(df1("u").isNull || df2("u").isNull || df1("v").isNull || df2("v").isNull)
        .count()
      !(count == 0)
    }

    var currentIter = 0
    val iterationLimit = Constant.iterationLimit
    df.persist(StorageLevel.MEMORY_AND_DISK)
    var isChange=true

    require(iterationLimit > 0, s"Maximum of iterations must be greater than 0, but got ${iterationLimit}")

    while(currentIter<iterationLimit && isChange){
      val largeRes = largeStar(df)
      val smallRes = smallStar(largeRes)
      //      smallRes.show
      currentIter = currentIter + 1
      isChange = diff(df,smallRes)
      df.unpersist()
      df = smallRes
      df.persist(StorageLevel.MEMORY_AND_DISK)
    }
    if(!isChange){
      log.error(s"converge after ${currentIter-1}, check cost 1 iteration")
    }else{
      log.error(s"reached iterationLimit")
    }

    df.select($"u".as("sha"),$"v").as("finger_new")
    df.createOrReplaceTempView("tmp_df1")

    spark.sql(Constant.createCCSql)
    spark.sql(s"insert overwrite table ${Constant.ccDistTable} partition(stat_date=$saveTime) select * from tmp_df1")



  }
}

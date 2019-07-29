package org.test

import org.apache.spark.sql.{SaveMode, SparkSession}
import org.apache.spark.storage.StorageLevel
import org.test.conf.{Constant, ExtendStructure}
import org.slf4j.LoggerFactory
import org.test.conf.FeatureEncoders.{extendStructureEncoder, featureResultEncoder}
/**
  * stderror : 日志
  * srdout : print语句
  */

object FingerPrintFeature {
  def main(args: Array[String]): Unit = {

    // slf4j
    val Logger = LoggerFactory.getLogger("TestApp")
    val spark = SparkSession
      .builder()
      .appName("testSparkHive")
//    When set to false, Spark SQL will use the Hive SerDe for parquet tables instead of the built in support.
//    .config("spark.sql.hive.convertMetastoreParquet", true)
      .config("spark.sql.hive.convertMetastoreOrc", true)
      .config("spark.sql.orc.enabled", true)
      .config("spark.speculation",true)
      .enableHiveSupport()
      .getOrCreate()

//    spark.conf.set("mapreduce.fileoutputcommitter.algorithm.version", "2")

    import org.apache.spark.sql.functions._
    import spark.sql
    val tableName = Constant.srcTable
    val saveTime = ParseUtil.parseSaveTime(args)

    val publicFields = Array("primaryKey","clusterid","status","magic","zid","burry") // 共有字段
    val features = Array("imei","mac","imsi","iccid","android","sdcid","serial","manufacturer"
      ,"model","brand","mainboard","equipment","hardware","product","android_ver","api_level"
      ,"version","picture_external_sd","picture_internal_sd","music_internal_sd","video_internal_sd"
      ,"music_external_sd","video_external_sd","resolution","size","cpu_core","cpu_kernel"
      ,"cpu_core_num","sys_fingerprint","cpu_fiber","decode_encode")  // feature需要对比的字段

    /*val origionDf = sql("SELECT primaryKey,clusterid,xid,id,zid,sdcid1,sdcid2,magic,burry,imei,mac,imsi,iccid,android,sdcid,serial," +
      "manufacturer,model,brand,mainboard,equipment,hardware,product,android_ver,api_level,version,picture_external_sd," +
      "picture_internal_sd,music_internal_sd,video_internal_sd,music_external_sd,video_external_sd,resolution,size," +
      "cpu_core,cpu_kernel,cpu_core_num,sys_fingerprint,cpu_fiber,decode_encode,status FROM " + tableName +" where save_time='"+saveTime+"'"//+" limit 100"
    )*/
    val origionSql = s"SELECT ${publicFields.mkString(",")},${features.mkString(",")} FROM ${tableName} where save_time='${saveTime}'"
    val origionDf = sql(origionSql)


    val colNames = origionDf.columns


    val df1Cols = colNames.map(x=>origionDf.col(x).as(s"t1_$x"))
    val df2Cols = colNames.map(x=>origionDf.col(x).as(s"t2_$x"))  // df2把字段名改成t2_前缀
    val df1 = origionDf.select(df1Cols:_*)
    val df2 = origionDf.select(df2Cols:_*)

    // 关联形成大表(xid,magic,imei确定一个设备)
    val extendDF = df1.join(df2,df1("t1_clusterid").equalTo(df2("t2_clusterid")))
      .where(df1("t1_primaryKey")<df2("t2_primaryKey"))

   /* var conditions = features.map(x=>{
      when((extendDF(s"t1_$x") === "null" || extendDF(s"t1_$x") === "" || extendDF(s"t1_$x") === "null|null|null|null" || extendDF(s"t1_$x") ==="|||") &&
        (extendDF(s"t2_$x") === "null" || extendDF(s"t2_$x") === "" || extendDF(s"t2_$x") === "null|null|null|null" || extendDF(s"t2_$x") ==="|||"), "0.5")
        .when(extendDF(s"t1_$x").equalTo(extendDF(s"t2_$x")), "1")
        .otherwise("0")
        .as(x)
    })*/

    var conditions = features.map(x=>{
      when((extendDF(s"t1_$x") === "null" || extendDF(s"t1_$x") === "" || extendDF(s"t1_$x") === "null|null|null|null" || extendDF(s"t1_$x") ==="|||") &&
        (extendDF(s"t2_$x") === "null" || extendDF(s"t2_$x") === "" || extendDF(s"t2_$x") === "null|null|null|null" || extendDF(s"t2_$x") ==="|||"), "2")
        .when(((extendDF(s"t1_$x") === "null" || extendDF(s"t1_$x") === "" || extendDF(s"t1_$x") === "null|null|null|null" || extendDF(s"t1_$x") ==="|||") &&
          (extendDF(s"t2_$x") =!= "null" || extendDF(s"t2_$x") =!= "" || extendDF(s"t2_$x") =!= "null|null|null|null" || extendDF(s"t2_$x") =!="|||"))or(
          (extendDF(s"t1_$x") =!= "null" || extendDF(s"t1_$x") =!= "" || extendDF(s"t1_$x") =!= "null|null|null|null" || extendDF(s"t1_$x") =!="|||") &&
            (extendDF(s"t2_$x") === "null" || extendDF(s"t2_$x") === "" || extendDF(s"t2_$x") === "null|null|null|null" || extendDF(s"t2_$x") ==="|||")
        ), "3")
        .when(extendDF(s"t1_$x").equalTo(extendDF(s"t2_$x")), "4")
        .otherwise("1")
        .as(x)
    })

    /*val primeryCol = when(col("t1_id")<col("t2_id"),md5(concat_ws("|",col("t1_id"),col("t2_id"))))
      .otherwise(md5(concat_ws("|",col("t2_id"),col("t1_id")))).as("primarykey")*/

    val allColumns = /*Array(primeryCol)++ */publicFields.map(x=>col(s"t1_$x"))++ conditions++ publicFields.map(x=>col(s"t2_$x"))

    var df4 = extendDF.select(allColumns:_*)

/*    // call chains too long , must checkpoint, avoid stackoverflow
    df4.repartition(500) // handle Size exceeds Integer.MAX_VALUE
    df4.cache()
    df4.checkpoint()*/

    val ds1 = df4.as[ExtendStructure]
    val ds2 = ds1.map(x=>ConvertUtil.convertFeature(x))
    val resulfDF = ds2.toDF()
    Logger.info("create table ..")
    spark.sql(Constant.createSql)
    Logger.info("store data to hive..")

    resulfDF.createOrReplaceTempView("lj_tmp1")
    sql(s"insert overwrite table ${Constant.distTable} partition(save_time=$saveTime) select * from lj_tmp1")
  }
}

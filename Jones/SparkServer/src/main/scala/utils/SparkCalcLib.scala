package utils

import java.io.{File, PrintWriter}
import java.text.SimpleDateFormat
import java.util.Date

import adtl.platform.Dataproto.{FilterOperation, QueryInfo}
import models.config.MainConf
import org.apache.spark.sql.{DataFrame, SparkSession}
import org.apache.spark.storage.StorageLevel
import org.joda.time.DateTime
import play.api.libs.json.{JsArray, Json}
import util.{DynamicUDF, MTrieFun}


object SparkCalcLib extends PrintLogger {
  info("mainConf 中的 key" + MainConf.confMap.keySet.mkString(" ;; "))
  val appConfResPath: String = MainConf("basicConf").app.resPath
  val qwFileRoot = s"/$appConfResPath/${new SimpleDateFormat("yyyyMMdd").format(new Date)}"

  def RunSQL(spark: SparkSession, query: QueryInfo): DataFrame = {
    val tables = SparkUtil.parseTable(query.sql.head)
    try {
      tables.foreach {
        name =>
          val refreshSql = s"REFRESH TABLE $name"
          info(refreshSql)
          spark.sql(refreshSql)
      }
    } catch {
      case e: Exception => e.printStackTrace()
    }

    //注册query_word的udf函数
    val qwPartition = GetQueryWordPartition(query)
    info(s"query words partitions: ${qwPartition}")
    if ((qwPartition != null) && (!"".equals(qwPartition))) {
      GenQueryWordsUdf(spark, query, qwPartition)
    }

    spark.sparkContext.setJobGroup(query.queryId, "", true)
    println("Submit Query: " + query.queryId)
    spark.sql(query.sql.head)
  }

  def GetQueryWordPartition(query: QueryInfo): String = {
    var qwPartition = ""
    val filters = query.filters
    for (filter <- filters) {
      val filterObj = Json.parse(filter)
      val filterObjType = (filterObj \ "type").as[String]
      if (FilterOperation.QW_PART.name.equals(filterObjType)) {
        qwPartition = (filterObj \ "value").as[String]
        return qwPartition
      } else if (FilterOperation.QW_AM.name.equals(filterObjType) || FilterOperation.QW_FM.name.equals(filterObjType)) {
        val valusArray = (filterObj \ "values").as[JsArray]
        val values = valusArray.as[List[String]]
        qwPartition = values(0)
        return qwPartition
      }
    }
    qwPartition
  }

  def GenQueryWordsUdf(spark: SparkSession, query: QueryInfo, qwPartition: String): Unit = {
    def writeFile(strs: Array[String], filename: String): Unit = {
      val write = new PrintWriter(new File(filename))
      strs.foreach {
        line =>
          write.println(line)
      }
      write.close()
    }

    val fileName = s"$qwFileRoot/$qwPartition"
    val saveFile = new File(fileName)
    if (!saveFile.getParentFile.exists()) saveFile.getParentFile.mkdirs()

    info(s"query words file location: $fileName")
    val wordListRdd = spark.sql(s"select word from default.query_word where part='$qwPartition'")
    val wordListRow = wordListRdd.collect()
    val wordList = wordListRow.map(_.getString(0))
    writeFile(wordList, fileName)
    spark.sparkContext.addFile(fileName)
    DynamicUDF.register(spark, MTrieFun.getFun(qwPartition), qwPartition)
    info("spark 函数注册resisted")
  }

  def RunHiveLocalCSV(spark: SparkSession, queryInfo: QueryInfo, respath: String) = {
    val queryName = MysqlUtil.getTaskNameById(queryInfo.queryId)
    val savePathTmp: String = s"${respath}/${DateTime.now.toString("yyyyMMdd")}/${queryName}_${queryInfo.queryId.substring(0, 4)}.xlsx"
    val returnPath: String = s"${DateTime.now.toString("yyyyMMdd")}/${queryName}_${queryInfo.queryId.substring(0, 4)}.xlsx"
    val excel = new ExcelWiter(savePathTmp)
    val helper = new JdbcHelper
    helper.execSql("hive", queryInfo, excel.writeRsMeta, excel.writeRsRow)
    excel.saveExcel(spark)
    //    HDFSUtil.getExcelFile(spark,savePathTmp)
    //    HDFSUtil.hdfsToFtp(spark, savePathTmp)
    returnPath
  }

  def RunSQLSaveLocalCSV(spark: SparkSession, query: QueryInfo, respath: String) = {
    val queryName = MysqlUtil.getTaskNameById(query.queryId)
    val runRes = RunSQL(spark, query)
    val header =
      if (query.outputs.length == 0 && query.dimensions.length == 0)
        runRes.schema.map(_.name)
      else
        query.dimensions.toList ::: query.outputs.toList

    info(s"outputs:${query.outputs}, dimensions:${query.dimensions}")
    info(s"header: ${header}")

    var returnPath: String = s"${DateTime.now.toString("yyyyMMdd")}/${queryName}_${query.queryId.substring(0, 4)}"
    if (query.session == "csv") {
      val savePathTmp: String = s"${respath}/${DateTime.now.toString("yyyyMMdd")}/${queryName}_${query.queryId.substring(0, 4)}"
      runRes.write.option("header", false).option("delimiter", "\t").csv(savePathTmp)
      HDFSUtil.hdfsToFtpCSV(spark, savePathTmp)
    } else {
      val savePathTmp: String = s"${respath}/${DateTime.now.toString("yyyyMMdd")}/${queryName}_${query.queryId.substring(0, 4)}.xlsx"
      returnPath = returnPath + ".xlsx"
      val transSchemaDf = runRes.columns.zipWithIndex.foldLeft(runRes) {
        case (res, (colName, idx)) =>
          res.withColumnRenamed(colName, header(idx))
      }
      info(s"DataFrame schema: ${transSchemaDf.schema.fields.map(x => (x.name, x.dataType.sql)).mkString(",")}")

      transSchemaDf.write.format("com.crealytics.spark.excel").option("header", "true").option("dateFormat", "yy-mmm-d").save(savePathTmp)
      //      HDFSUtil.hdfsToFtp(spark, savePathTmp)
    }
    returnPath
  }

  def CancelSQLByGroupId(sparkSession: SparkSession, groupId: String) = {
    sparkSession.sparkContext.cancelJobGroup(groupId)
  }


  def RunCustomTask(sparkSession: SparkSession, queryInfo: QueryInfo, respath: String) = {
    val savePathTmp: String = s"${respath}/${DateTime.now.toString("yyyyMMdd")}/${queryInfo.queryId}"
    val returnPath: String = s"${DateTime.now.toString("yyyyMMdd")}/${queryInfo.queryId}"
    import sparkSession.implicits._
    val imeidata = sparkSession.sparkContext.makeRDD(queryInfo.sql(0).split(",", -1)).toDF("uid")
    //      .pipe("./test").map{
    //      line=>
    //        val arr = line.split("\t", -1)
    //        (arr(0), arr(1))
    //    }.toDF("uid", "aid")
    imeidata.createOrReplaceTempView(queryInfo.queryId)
    sparkSession.sql(s"REFRESH TABLE adtl.echo_alive")
    val maxDt = sparkSession.sql("select max(distinct dt) as lastest_dt from adtl.echo_alive").rdd.map(_.getString(0)).first()
    val runRes = sparkSession.sql(s"select a.uid, b.imei from ${queryInfo.queryId} a left join adtl.echo_alive b on a.uid=b.uid and b.dt=$maxDt").persist(StorageLevel.MEMORY_AND_DISK_SER)
    val resCount = runRes.count()
    val partitionNum = (resCount / 50000).toInt + 1
    runRes.repartition(partitionNum).write.format("com.databricks.spark.csv")
      .option("delimiter", "\t")
      .save(savePathTmp)
    runRes.unpersist()
    HDFSUtil.downLoadFileWithMeta(sparkSession, savePathTmp, queryInfo.dimensions.toList ::: queryInfo.outputs.toList)
    //    savePathTmp
    returnPath
  }

  def loadCSVtoTable(spark: SparkSession, fileLocation: String, tableName: String): Unit = {
    val df = spark.read.format("csv")
      //      .option("sep", ";")
      .option("inferSchema", "true")
      .option("header", "true")
      .load(fileLocation)

    df.write.saveAsTable(tableName)
  }

  def quicklyQuery(spark: SparkSession, sql: String): List[Map[String, String]] = {
    // 解析出表名, 进行refresh
    val tables = SparkUtil.parseTable(sql)
    tables.foreach(t => {
      val refreshSql = s"refresh $t"
      info(refreshSql)
      spark.sql(refreshSql).foreach(_ => {})
    })

    // 获取数据
    info(sql)
    val df = spark.sql(sql)
    val res = df.collect()
    val cols = df.columns
    val fieldCount = res(0).length
    res.toList.map(row => {
      var map = Map[String, String]()
      for (i <- 0 until fieldCount) {
        map = map + (cols(i) -> row.getAs[String](cols(i)))
      }
      map
    })
  }
}


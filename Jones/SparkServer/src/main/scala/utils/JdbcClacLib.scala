package utils

import java.io._
import java.util

import adtl.platform.Dataproto.QueryInfo
import org.apache.spark.sql.SparkSession
import org.joda.time.DateTime

object JdbcClacLib {

  def execSqlAndSaveResult(spark:SparkSession,dsType: String, queryInfo: QueryInfo, respath: String): String = {
    val helper = new JdbcHelper
    val queryName = MysqlUtil.getTaskNameById(queryInfo.queryId)
    val savePathTmp: String = s"${respath}/${DateTime.now.toString("yyyyMMdd")}/${queryName}_${queryInfo.queryId.substring(0,4)}.xlsx"
    val returnPath: String = s"${DateTime.now.toString("yyyyMMdd")}/${queryName}_${queryInfo.queryId.substring(0,4)}.xlsx"
    val excel = new ExcelWiter(savePathTmp)
    helper.execSql(dsType, queryInfo, excel.writeRsMeta, excel.writeRsRow)
    excel.saveExcel(spark)
//    excel.saveLocal()   // for test
//    HDFSUtil.hdfsToFtp(spark,savePathTmp)
    returnPath
  }
}

object TestJavaSerial {
  def main(args: Array[String]): Unit = {
    val list = new util.ArrayList[Any]()
    list.add(1)
    list.add("a")
    // 序列化
    import java.io.ObjectOutputStream
    val bos = new ByteArrayOutputStream()
    val oos = new ObjectOutputStream(bos)
    oos.writeObject(list)
    val bytes = bos.toByteArray

    // 反序列化
    import java.io.ObjectInputStream
    val ois2 = new ObjectInputStream(new ByteArrayInputStream(bytes))
    val list2 = ois2.readObject.asInstanceOf[util.ArrayList[Any]]
    for (i <- 0 until list2.size()) {
      println(list2.get(i) + ": :" + list2.get(i).getClass)
    }
  }
}
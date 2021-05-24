package com.test

import java.io.FileInputStream
import java.util.function.Consumer

import com.monitorjbl.xlsx.StreamingReader
import models.ServerInfo
import org.apache.poi.ss.usermodel.Cell
import org.apache.poi.xssf.streaming.SXSSFWorkbook
import org.apache.spark.sql.SparkSession
import utils.{SparkUtil, XSSFUtil}

import scala.util.control.Breaks._

object TestPOI {
  def main(args: Array[String]): Unit = {
    val limit = 500
    val file = "/Users/liujie02/Downloads/ef785762a2546d0d30b54665e1644727.xlsx"
    val in = new FileInputStream(file)
    val res = XSSFUtil.printExcelByStream(in,500)
    println (res)

  }
  def main1(args: Array[String]): Unit = {
    import java.io.FileOutputStream
    try {
      val startTime = System.currentTimeMillis
      val NUM_OF_ROWS = 50*10000
      val NUM_OF_COLUMNS = 30
      var wb:SXSSFWorkbook = null
      try {
        wb = new SXSSFWorkbook();
        wb.setCompressTempFiles(true) //压缩临时文件，很重要，否则磁盘很快就会被写满

        var sh = wb.createSheet
        var rowNum = 0
        for (num <- 0 until NUM_OF_ROWS) {
          if (num % 1000000 == 0) {
            sh = wb.createSheet()
            rowNum = 0
          }
          rowNum += 1
          val row = sh.createRow(rowNum)
          for (cellnum <- 0 until NUM_OF_COLUMNS) {
            val cell = row.createCell(cellnum)
            cell.setCellValue(Math.random)
          }
        }
        val out = new FileOutputStream("ooxml-scatter-chart_SXSSFW_" + NUM_OF_ROWS + ".xlsx")
        wb.write(out)
        out.close()
      } catch {
        case ex: Exception =>
          ex.printStackTrace()
      } finally if (wb != null) wb.dispose // 删除临时文件，很重要，否则磁盘可能会被写满
      val endTime = System.currentTimeMillis
      System.out.println("process " + NUM_OF_ROWS + " spent time:" + (endTime - startTime))
    } catch {
      case e: Exception =>
        e.printStackTrace()
        throw e
    }
  }
  def main2(args: Array[String]): Unit = {
    val spark = SparkSession.builder().master("local")
      .getOrCreate()
    println(spark)
    new Thread(
      new Runnable {
        override def run(): Unit = {
          val spark2 = SparkUtil.getSparkSession()
          println(spark2)
        }
      }
    ).start()

    import spark.implicits._
    val df = Seq(("张三", 32),("李四",24)).toDF("name","age")
    df.write
      .format("com.crealytics.spark.excel")
      .option("header", "true")
      .option("dateFormat", "yy-mmm-d") // Optional, default: yy-m-d h:mm
      .option("timestampFormat", "mm-dd-yyyy hh:mm:ss") // Optional, default: yyyy-mm-dd hh:mm:ss.000
      .mode("overwrite") // Optional, default: overwrite.
      .save("/Users/liujie02/Downloads/testPerson.xlsx")
  }
}

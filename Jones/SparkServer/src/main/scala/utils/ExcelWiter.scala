package utils

import java.io.{File, FileOutputStream}
import java.sql.{ResultSet, Types}
import scala.util.control.Breaks._
import org.apache.hadoop.fs.{FileSystem, Path}
import org.apache.poi.xssf.streaming.SXSSFWorkbook
import org.apache.spark.sql.SparkSession
//import shadeio.poi.xssf.usermodel.XSSFWorkbook

class ExcelWiter(outPath: String) extends PrintLogger {
  private val wb = new SXSSFWorkbook()
  wb.setCompressTempFiles(true)

  private var curSheet = wb.createSheet()

  val writeRsMeta = (rsmd: Seq[String]) => {
    val row = curSheet.createRow(0)
    val cnt = rsmd.size
    (0 until cnt).foreach {
      idx =>
        row.createCell(idx).setCellValue(rsmd(idx))
    }
  }

  val writeRsRow = (rs: ResultSet, columnCount: Int, meta: Seq[String]) => {
    var rowCnt = 1 // 第一次写的时候因为有标题, 所以从第二行开始写
    var sheetCnt = 1
    val maxRowCnt = 1048575 // 每个sheet页最多100万条数据, 奖加上表头共1048576
    val sheetLimit = 4 // 最多写4页
    val rsMeta = rs.getMetaData
    (1 to columnCount).foreach {
      idx => info(s"第${idx}个字段类型:  ${rsMeta.getColumnType(idx)}")
    }
    breakable {
      while (rs.next()) {
        if (rowCnt > maxRowCnt) { // 如果超过最大行数, 就创建一个新的sheet页, 行数置为1
          curSheet = wb.createSheet()
          this.writeRsMeta(meta)
          rowCnt = 1
          sheetCnt = sheetCnt + 1
          if (sheetCnt > sheetLimit) { //超过sheet页限制就跳出
            break
          }
        }
        val currentRow = curSheet.createRow(rowCnt)
        rowCnt = rowCnt + 1
        (1 to columnCount).foreach {
          idx =>
            var cellValue = ""
            val fieldObj = rs.getObject(idx)
            if (fieldObj != null) {
              rsMeta.getColumnType(idx) match {
                case Types.SMALLINT =>
                  val cell = rs.getLong(idx)
                  currentRow.createCell(idx - 1).setCellValue(cell)
                case Types.INTEGER =>
                  val cell = rs.getLong(idx)
                  currentRow.createCell(idx - 1).setCellValue(cell)
                case Types.BIGINT =>
                  val cell = rs.getLong(idx)
                  currentRow.createCell(idx - 1).setCellValue(cell)
                case Types.INTEGER =>
                  val cell = rs.getInt(idx)
                  currentRow.createCell(idx - 1).setCellValue(cell)
                case Types.DECIMAL =>
                  val cell = rs.getDouble(idx)
                  currentRow.createCell(idx - 1).setCellValue(cell)
                case Types.FLOAT =>
                  val cell = rs.getDouble(idx)
                  currentRow.createCell(idx - 1).setCellValue(cell)
                case Types.DOUBLE =>
                  val cell = rs.getDouble(idx)
                  currentRow.createCell(idx - 1).setCellValue(cell)
                case _ =>
                  cellValue = fieldObj.toString
                  currentRow.createCell(idx - 1).setCellValue(cellValue)
              }
            } else {
              currentRow.createCell(idx - 1).setCellValue(cellValue)
            }
        }
      }
    }
  }

  def saveExcel(spark: SparkSession) = {
    HDFSUtil.delHdfsFilePath(spark, outPath)
    val hadoopConf = spark.sparkContext.hadoopConfiguration
    val hdfs = FileSystem.get(hadoopConf)
    val fsPath = new Path(outPath) // hdfs
    val out = hdfs.create(fsPath)
    wb.write(out)
    out.flush()
    out.close()
  }

  def saveLocal() = {
    val out = new FileOutputStream("/Users/liujie02/Downloads/123.xlsx")
    wb.write(out)
    out.flush()
    out.close()
    wb.close()
  }
}

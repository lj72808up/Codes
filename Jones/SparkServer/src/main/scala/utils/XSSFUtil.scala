package utils

import java.io.InputStream
import java.util
import java.util.function.Consumer

import scala.util.control.Breaks._
import adtl.platform.Dataproto.{Linedata, TranStream}
import com.monitorjbl.xlsx.StreamingReader
import org.apache.poi.ss.usermodel.Cell
/*import shadeio.poi.ss.usermodel.CellType
import shadeio.poi.xssf.usermodel.{XSSFCell, XSSFWorkbook}*/

import scala.collection.mutable.ArrayBuffer

object XSSFUtil extends PrintLogger {

  /*def getCellValue(cell: XSSFCell): String = {
    cell.getCellType match {
      case CellType.BOOLEAN =>
        cell.getBooleanCellValue.toString
      case CellType.NUMERIC =>
        cell.getNumericCellValue.toString
      case _ =>
        cell.getStringCellValue
    }
  }

  def printExcelFile(in: InputStream, limitCnt: Int): TranStream = {
    val sheet = new XSSFWorkbook(in).getSheetAt(0)
    import collection.JavaConversions._
    var schema = Array[String]()
    var lines = Array[Linedata]()
    var cnt = 1
    breakable{
      sheet.rowIterator().foreach {
        row => {
          val line = row.cellIterator().toArray.map(
            x => getCellValue(x.asInstanceOf[XSSFCell])
          )
          if (cnt == 1) { // schema行
            schema = line
          } else if (cnt <= limitCnt + 1) {
            val lineData = Linedata().withField(line)
            lines = lines :+ lineData
          } else {
            info(s"$limitCnt 条已获取, 不再继续查看更多数据")
            break
          }
          cnt = cnt + 1
        }
      }
    }

    TranStream().withLine(lines).withSchema(schema)
  }*/

  def printExcelByStream(in: InputStream, limit: Int): TranStream = {
    val wk = StreamingReader.builder()
      .rowCacheSize(100)  //缓存到内存中的行数，默认是10
      .bufferSize(4096)      //读取资源时，缓存到内存的字节大小，默认是1024
      .open(in)                        //打开资源，必须，可以是InputStream或者是File，注意：只能打开XLSX格式的文件

    val sheet = wk.getSheetAt(0)
    var schema = Array[String]()
    var lines = Array[Linedata]()      // 读取所有的limit行
    //遍历所有的行
    val iter = sheet.iterator()
    breakable {
      while (iter.hasNext) {
        val row = iter.next()
        val rowNum = row.getRowNum
//        info("开始遍历第" + rowNum + "行数据：")
        if (row.getRowNum > limit) {  // 第一行为列名
          break
        }
        val line = ArrayBuffer[String]()
        row.cellIterator().forEachRemaining(new Consumer[Cell] {
          override def accept(cell: Cell): Unit = {
            line.append(cell.getStringCellValue)
          }
        })
        if (rowNum == 0) {
          schema = line.toArray
        } else {
          val lineData = Linedata().withField(line.toArray[String])
          lines = lines :+ lineData
        }
      }
    }
    in.close()
    wk.close()
    TranStream().withLine(lines).withSchema(schema)
  }
}

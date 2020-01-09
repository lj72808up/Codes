package com.sogou.util

import java.io.{BufferedOutputStream, FileOutputStream, PrintWriter}

class FileOutputUtil {

  def write2File(content:String, file:String):Unit = {
    val fos = new FileOutputStream(file)
    val bos2 = new BufferedOutputStream(fos)
//    val writer = new PrintWriter(bos)
//    writer.write(content)
//    writer.flush()
    bos2.write(content.getBytes())
    bos2.flush()
  }

}


object Test{
  def main(args: Array[String]): Unit = {
    new FileOutputUtil().write2File("bbb","test.file")
  }
}
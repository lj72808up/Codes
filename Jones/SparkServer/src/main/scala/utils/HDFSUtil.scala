package utils

import java.io._
import java.nio.charset.StandardCharsets

import org.apache.commons.io.IOUtils
import org.apache.commons.net.ftp.{FTP, FTPClient}
import org.apache.hadoop.fs.{FSDataInputStream, FileSystem, Path}
import org.apache.spark.sql.SparkSession

object HDFSUtil {
  def downLoadFile(spark: SparkSession, path: String) = {
    val hadoopConf = spark.sparkContext.hadoopConfiguration
    val hdfs = FileSystem.get(hadoopConf)
    val fsPath = new Path(path)
    val saveFile = new File(s"/$path")
    val os = new FileOutputStream(s"/$path")
    if (!saveFile.getParentFile.exists()) saveFile.getParentFile.mkdirs()
    if (saveFile.exists() && saveFile.isFile) saveFile.delete() //若文件存在直接删除
    hdfs.listStatus(fsPath).filter(_.getPath.getName.endsWith(".csv")).foreach {
      fs =>
        val dataPath = fs.getPath
        val in = hdfs.open(dataPath)
        IOUtils.copy(in, os)
    }
    os.close()
  }

  def downLoadFileWithMeta(spark: SparkSession, path: String, Meta: Seq[String]) = {
    val hadoopConf = spark.sparkContext.hadoopConfiguration
    val hdfs = FileSystem.get(hadoopConf)
    val fsPath = new Path(path) // hdfs
    val saveFile = new File(s"/$path") // local
    if (!saveFile.getParentFile.exists()) saveFile.getParentFile.mkdirs() // 创建上层目录
    if (saveFile.exists() && saveFile.isFile) saveFile.delete() //若文件存在直接删除
    val os = new FileOutputStream(s"/$path") // local
    /*// add bom
    val bom = Array[Byte](0xEF.toByte, 0xBB.toByte, 0xBF.toByte)
    os.write(bom)
    // end add
    val schemaStr = Meta.mkString("\t") + "\n"*/
    val schemaStr = new String(Array(0xEF, 0xBB, 0xBF).map(_.toByte), StandardCharsets.UTF_8) + Meta.mkString("\t") + "\n"

    os.write(schemaStr.getBytes)
    hdfs.listStatus(fsPath).filter(_.getPath.getName.endsWith(".csv")).foreach {
      fs =>
        val dataPath = fs.getPath
        val in = hdfs.open(dataPath)
        IOUtils.copy(in, os)
    }
    os.close()
  }

  def delHdfsFilePath(spark: SparkSession, path: String, isRecusrive: Boolean = true) = {
    val hdfs = FileSystem.get(spark.sparkContext.hadoopConfiguration)
    hdfs.delete(new Path(path), isRecusrive)
  }

  def getExcelFile(spark: SparkSession, path: String) = {
    val hadoopConf = spark.sparkContext.hadoopConfiguration
    val hdfs = FileSystem.get(hadoopConf)
    val fsPath = new Path(path) // hdfs
    val saveFile = new File(s"/$path") // local
    if (!saveFile.getParentFile.exists()) saveFile.getParentFile.mkdirs()
    if (saveFile.exists() && saveFile.isFile) saveFile.delete()
    val os = new FileOutputStream(s"/$path")
    val in = hdfs.open(fsPath)
    IOUtils.copy(in, os)
    os.close()
  }

  // TODO: path 使用datacenter开头 ftp使用域名 
  def hdfsToFtp(spark: SparkSession, path: String) = {
    val hadoopConf = spark.sparkContext.hadoopConfiguration
    val hdfs = FileSystem.get(hadoopConf)
    val fsPath = new Path(path) // hdfs
    val in = hdfs.open(fsPath)

    val ftp = new FTPClient()
    ftp.setControlEncoding("utf-8")
    ftp.connect("datacenter.ftp.adtech.sogou") // ROC申请域名并使用域名
    ftp.login("ftp", "")
    ftp.makeDirectory(fsPath.getParent.toString)

    ftp.setFileType(FTP.BINARY_FILE_TYPE)

    println(ftp.storeFile(fsPath.toString, in))
    in.close()
    ftp.logout()
  }

  def hdfsToFtpCSV(spark: SparkSession, path: String) = {
    val hadoopConf = spark.sparkContext.hadoopConfiguration
    val hdfs = FileSystem.get(hadoopConf)
    val fsPath = new Path(path)

    // 目标csv文件的文件名为 xxx.csv
    val targetMergeCsvFile = s"${path}.csv"

    val al = new java.util.ArrayList[FSDataInputStream]
    hdfs.listStatus(fsPath).filter(_.getPath.getName.endsWith(".csv")).foreach {
      fs =>
        val dataPath = fs.getPath
        val in = hdfs.open(dataPath)
        al.add(in)
    }
    val sio = new SequenceInputStream(java.util.Collections.enumeration(al))
    val os = hdfs.create(new Path(targetMergeCsvFile))
    IOUtils.copy(sio,os)

    /*val ftp = new FTPClient()
    ftp.setControlEncoding("utf-8")
    ftp.connect("datacenter.ftp.adtech.sogou") // ROC申请域名并使用域名
    ftp.login("ftp", "")
    ftp.makeDirectory(fsPath.getParent.toString)

    ftp.setFileType(FTP.BINARY_FILE_TYPE)

    val out = ftp.storeFileStream(path)
    var n = 0
    val by =  new Array[Byte](1024 * 4)
    while(n != -1) {
      n = sio.read(by)
      if (n != -1) {
        out.write(by, 0, n)
      }
    }
    out.close()
    ftp.logout()*/

    sio.close()
    os.close()
  }
}

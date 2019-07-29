package com.test.utils


import java.net.URI

import org.apache.hadoop.conf.Configuration
import org.apache.hadoop.fs.{FileSystem, Path}
import java.io.{BufferedReader, InputStreamReader, PrintWriter}

import com.test.log.Logger


class HdfsUtil(fs:FileSystem) extends Logger{
  def getDirectorySize(url:String)={
    val size = fs.getContentSummary(new Path(url)).getLength
    size // bytes
  }
  def getFileContent(url:String, f:String=>Unit):Unit = {
    val stream = fs.open(new Path(url))
    val reader = new BufferedReader(new InputStreamReader(stream))
    var content = reader.readLine()
    while (content != null) {
      if (content.startsWith("INFO")){
        f(content)
      }
      content = reader.readLine()
    }
    reader.close()
  }
  def appendToFile(content:String,dest:String):Unit = {
    val destPath = new Path(dest)
    val fs_append = fs.append(destPath)
    val writer = new PrintWriter(fs_append)
    writer.append(content)
    writer.flush()
    fs_append.hflush()
    writer.close()
    fs_append.close()
    log.info("append file " + dest)
  }

  def close():Unit = {
    fs.close()
  }
}

object HdfsUtil{
  def apply(hdfs:String,user:String): HdfsUtil = {
    val conf = new Configuration()
    conf.setBoolean("dfs.support.append", true)
    conf.set("fs.hdfs.impl", "org.apache.hadoop.hdfs.DistributedFileSystem")
    conf.set("fs.file.impl", "org.apache.hadoop.fs.LocalFileSystem")
    val fs = FileSystem.get(new URI(hdfs), conf, user)
    new HdfsUtil(fs)
  }

  def main(args: Array[String]): Unit = {
    val util = HdfsUtil("","slab")
    util.getDirectorySize("/aaa/README.txt")
  }
}

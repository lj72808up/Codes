package utils.udf

import java.io
import java.net.{URL, URLClassLoader}

import models.config.MainConf
import org.apache.spark.sql.SparkSession
import sourcecode.File
import utils.{PrintLogger, SparkUtil}

object RegistFunctionUtil extends PrintLogger {

  val spark: SparkSession = SparkUtil.getSparkSession()

  // 程序启动时注册所有已存在的jar包
  def registOnStart(): Unit = {

    /*val udfJar = "com.sogou.adtl.dw.tool.udf.RegisterUtil"
        try {
          val toolClass = Class.forName(udfJar)
          toolClass.getMethod("registerAll", classOf[SparkSession]).invoke(toolClass.getClass, spark)
        } catch {
          case e: Exception => println(s"$udfJar not found, error ${e.getClass} ${e.getMessage}")
        }*/

    //    spark.sparkContext.addFile("test")
  }

  // 程序运行时注册用户上传的jar包
  def registOnRunning(file: String, className: String, registMethod: String): Unit = {
    val funcDir = MainConf("basicConf").app.fuctionJars
    //    val baseDir
    val jarFile = new java.io.File(funcDir, file).getPath
    info(s"------------------ jar location:$jarFile----------------------")
    spark.sparkContext.addJar(s"local:/$jarFile")
    val cl = URLClassLoader.newInstance(Array[URL](new URL(s"file://$jarFile")))
    val registClass = cl.loadClass(className)
    registClass.getMethod(registMethod, classOf[SparkSession]).invoke(registClass.getClass, spark)
    info(s"------------------  ${registClass.getName}注册成功 ------------------ ")
  }

  private def getResourceInJar(file: String): String = {
    val clz = SparkUtil.getClassLoader()
    val mainJarPath = clz.getResource(file).getFile
    mainJarPath
  }

  def showRegistedFunction(): Seq[String] = {
    val sql = "show functions"
    val functions = spark.sql(sql).collect.map(row => row.getAs[String](0))
    functions
  }

  def main(args: Array[String]): Unit = {
    //    registOnRunning("sogou-adtl-dw-tool-udf-assembly-0.1.jar", "com.sogou.adtl.dw.tool.udf.RegisterUtil", "registerAll")
    val file = new io.File("/Users/liujie02/IdeaProjects/SparkServer/", "/abcdefg")
    println(file.getPath)
  }
}

package com.test.utils

import java.io._


class CommandUtil {
  /**
    *
    * @param file : 相对路径 , 该文件下是多行shell
    * @param wd : work dir
    * @return
    */
  def compactShell(file:String,wd:String):String = {
    val location = s"${new File(wd).getCanonicalPath}${File.separator}${file}"
    val target = new FileInputStream(location)
    val reader = new BufferedReader(new InputStreamReader(target))
    val builder = new StringBuilder

    var line = reader.readLine()
    while (line!=null){
      if (! line.startsWith("#!")) {
        builder.append(line)
        builder.append(";")
      }
      line = reader.readLine()
      //}
    }
    builder.toString()
  }


  /**
    * For the second implementation by ProcessBuilder. This is preferred over the Runtime approach because we’re able to customize some details.
    *
    * For example we’re able to:
    *
    * change the working directory our shell command is running in using builder.directory()
    * set-up a custom key-value map as environment using builder.environment()
    * redirect input and output streams to custom replacements
    * inherit both of them to the streams of the current JVM process using builder.inheritIO()
    */

  def executeStrBash(cmd:String, waitFor:Boolean=false, workdir:String=System.getProperty("user.dir")):InputStream ={
    executeListBash(cmd.split(" "),waitFor,workdir)
  }

  def executeListBash(cmd:Array[String], waitFor:Boolean=false, workdir:String=System.getProperty("user.dir")):InputStream ={

    val builder = new ProcessBuilder()
    builder.redirectErrorStream(true)
    builder.command(cmd:_*)
    builder.directory(new File(workdir))
    val process = builder.start()
    var code="unknown"
    if(waitFor)
      process.waitFor()
    process.getInputStream
  }

  def cmdOutput[T](ins:InputStream,fun:String=>T):Unit = {
    val reader = new BufferedReader(new InputStreamReader(ins))
    var line:String = reader.readLine()
    while(line!=null){
      fun(line)
      line = reader.readLine()
    }
    ins.close()
    reader.close()
  }
}

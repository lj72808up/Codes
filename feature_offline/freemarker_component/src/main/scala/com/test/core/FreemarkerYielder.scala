package com.test.core
import java.io._
import java.util

import freemarker.template.Configuration

class FreemarkerYielder {

  def getCongtent(templateDir:String,templateFile:String,data:util.HashMap[String,String],distFile:String): String ={

    val configuration = new Configuration(Configuration.VERSION_2_3_28)
    configuration.setTagSyntax(Configuration.AUTO_DETECT_TAG_SYNTAX);
    configuration.setClassLoaderForTemplateLoading(this.getClass.getClassLoader,templateDir)

    val template = configuration.getTemplate(templateFile,"utf-8")

    var writer:Writer = null
    var reader:BufferedReader = null
    val builder = new StringBuilder
    try{
      writer = new FileWriter(distFile);
      template.process(data,writer)

      reader = new BufferedReader(new InputStreamReader(new FileInputStream(distFile)))
      var line:String = reader.readLine()
      while(line != null) {
        println(line)
        builder.append(line)
        line = reader.readLine()
      }
    }finally {
      writer.close()
      reader.close()
    }
    builder.toString
  }
}

package com.test.database

import java.util

import com.test.utils.HdfsUtil


//gen getter and setter
case class Item(`partition`:String,
                     `version`:String,
                     `dataType`:String,
                     `startTime`:String,
                     `endTime`:String=""){
  override def toString: String = {
    s"INFO:${partition},${version},${dataType},${startTime},${endTime}\n"
  }
}


object MarkUtil{
  def getFileContent(content:String):Option[Item] = {
    val arr = content.split(",")
    if (arr.length>=5){
      Some(Item(arr(0).split(":")(1),
        arr(1),
        arr(2),
        arr(3),
        arr(4)
      ))
    }else if(arr.length == 4){
      Some(Item(arr(0).split(":")(1),
        arr(1),
        arr(2),
        arr(3)
      ))
    }else{
      None
    }
  }

  /**
    * @param url1: out的url
    * @param url2: in的url
    * @return
    */
  def getPartitions(url1:String,url2:String,hdfs:String,user:String):(String,String,Item) = {
    val arr1 = new java.util.ArrayList[String]()
    val arr2 = new java.util.ArrayList[String]()
    def f(arr:java.util.ArrayList[String])(content:String):Unit = {
      arr.add(content)
    }

    val util = HdfsUtil(hdfs,user)
    util.getFileContent(url1,f(arr1))
    util.getFileContent(url2,f(arr2))

    println(arr1)
    println(arr2)

    if(arr1.size()==arr2.size()){
      util.close()
      return ("","",null)
    }
    else{
      val oldItem = getFileContent(arr1.get(arr1.size()-1)).get
      val newItem = getFileContent(arr2.get(arr1.size())).get

      val oldPartition = oldItem.partition
      val newPartition = newItem.partition
      util.close()
      return (oldPartition,newPartition,newItem)
    }
  }
  /*def putResultSet[T:ClassTag](rs:ResultSet, obj:Class[T]): util.ArrayList[T] ={
    val arr = new util.ArrayList[T]()
    val metaData = rs.getMetaData
    val fieldCnt = metaData.getColumnCount

    while (rs.next()){
      val newInstance = obj.newInstance()
      for (i<- 1 to fieldCnt){
        val dbField = metaData.getColumnClassName(i).toLowerCase
        val javaField = toJavaField(dbField)
        val firstCharacter = javaField.substring(0,1)
        val invokeField = javaField.replaceFirst(firstCharacter,firstCharacter.toUpperCase)
        val fieldType = obj.getDeclaredField(invokeField).getType
        val method = obj.getMethod("set"+invokeField,fieldType)

        // read field from resultset
        if(fieldType.isAssignableFrom(classOf[String])){
          method.invoke(newInstance, rs.getString(i))
        }else if(fieldType.isAssignableFrom(classOf[java.lang.Integer]) || fieldType.isAssignableFrom(classOf[Int])){
          method.invoke(newInstance, rs.getInt(i))
        }else if(fieldType.isAssignableFrom(classOf[java.lang.Boolean]) || fieldType.isAssignableFrom(classOf[Boolean])){
          method.invoke(newInstance, rs.getBoolean(i))
        }else if(fieldType.isAssignableFrom(classOf[Date])){
          method.invoke(newInstance, rs.getDate(i))
        }
      }
      arr.add(newInstance)
    }
    // 数据库命名格式转java命名格式
    def toJavaField(dbField:String):String = {
      val splits = dbField.split("_")
      val builder = new StringBuilder
      builder.append(splits(0))

      // 一下划线分割首字母转大写
      if (splits.length>1){
        for (i <- 1 until splits.length){
          val iField = splits(i)
          val firstCharacter = iField.substring(0,1)
          val javaFieldFactor = iField.replaceFirst(firstCharacter,firstCharacter.toUpperCase())
          builder.append(javaFieldFactor)
        }
      }
      builder.toString()
    }

    arr
  }*/

  def main(args: Array[String]): Unit = {
//    val c = "INFO:20190620,v3.2.1,incresement,20190618,20190620"
//    println(MarkUtil.getFileContent(c))
  }
}
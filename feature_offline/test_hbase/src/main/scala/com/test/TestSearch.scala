package com.test
import org.apache.hadoop.hbase.TableName
import org.apache.hadoop.hbase.HBaseConfiguration
import org.apache.hadoop.hbase.client.{ConnectionFactory, Get, HTable}

object TestSearch {
  def main(args: Array[String]): Unit = {

    val conf1 = HBaseConfiguration.create
    val connection = ConnectionFactory.createConnection(conf1)
    val table = connection.getTable(TableName.valueOf("Contacts"))
    try{
      for (i <- 1 to 1000){
        val get1 = new Get("115193184_132dbb85b905fa4610c0ab932d424386ecf688e0".getBytes())
        get1.addColumn("info".getBytes(),"sha1".getBytes())
        val res1 = table.get(get1)
        val sha_1 = res1.getValue("info".getBytes(),"sha1".getBytes())
        println(s"info:sha1=> $sha_1")
      }
    }finally {
      table.close
      connection.close()
    }
  }
}

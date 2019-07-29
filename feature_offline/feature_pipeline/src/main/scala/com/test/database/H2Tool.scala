package com.test.database

import java.sql.Connection
import java.sql.DriverManager
import java.sql.Statement
import java.util

import com.test.log.Logger

import scala.reflect.ClassTag

class H2Tool(dbUrl:String) extends Logger{
  
  val JDBC_DRIVER = "org.h2.Driver"
  //val dbUrl = "jdbc:h2:~/test"

  // Database credentials
  val USER = "slab"
  val PASS = ""
  var conn: Connection = _
  var isConnected = false

  def init() = {
    // STEP 1: Register JDBC driver
    Class.forName(JDBC_DRIVER)
    // STEP 2: Open a connection
    System.out.println("Connecting to a database...")
    conn = DriverManager.getConnection(dbUrl, USER, PASS)
    isConnected = true
  }

  def execSql(sql:String):Unit = {
    if (!isConnected) {
      this.init()
    }
    var stmt: Statement = null
    try {
      stmt = conn.createStatement()
      log.info("exec sql:" + sql)
      stmt.execute(sql)
    } catch {
      case ex: Exception => log.error("h2db fail:" + ex.toString)
    } finally {
      if (stmt != null) {
        stmt.close()
      }
    }
  }

  def readRecord(sql:String):util.ArrayList[Item] = {
    if (!isConnected) {
      this.init()
    }
    var stmt: Statement = null
    var arr:util.ArrayList[Item] = new util.ArrayList[Item]()
    try {
      stmt = conn.createStatement()
      log.info("exec sql:" + sql)
      val rs = stmt.executeQuery(sql)

      while(rs.next()){
        val partition = rs.getString("partition")
        val version = rs.getString("version")
        val dataType = rs.getString("dataType")
        val startTime = rs.getString("startTime")
        val endTime = rs.getString("endTime")

        arr.add(Item(partition,version,dataType,startTime,endTime))
      }

    } catch {
      case ex: Exception => log.error("h2db fail:" + ex.toString)
    } finally {
      if (stmt != null) {
        stmt.close()
      }
    }
    arr
  }


  def insertRecord(item:Item,table:String):Unit = {
    val sql = s"insert into ${table}(partition,version,dataType,startTime,endTime) values ('${item.partition}','${item.version}','${item.dataType}','${item.startTime}','${item.endTime}') "
    log.info("insert sql:"+sql)
    execSql(sql)
  }

  def close():Unit = {
    if(conn!=null) {
      conn.close()
      log.info("h2 conn close")
    }
  }

}


object Test{
  def main(args: Array[String]): Unit = {
    val tool = new H2Tool("jdbc:h2:/tmp/test")
    tool.execSql("drop table test_001")
    tool.execSql(
      """create table if not exists test_001(
        | partition varchar(200),
        | version varchar(200),
        | dataType  varchar(200),
        | startTime varchar(200),
        | endTime varchar(200)
        |)""".stripMargin)
    val item = Item("20190112","3.1.2","increase","firstday","nextday")
    tool.insertRecord(item,"test_001")
//    tool.execSql("insert into test_001(partition,version,dataType,startTime,endTime) values (20190112,312,increase,firstday,nextday)")
    println(tool.readRecord("select * from test_001"))
    tool.close()
  }
}
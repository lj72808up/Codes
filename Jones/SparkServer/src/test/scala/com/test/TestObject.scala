package com.test

import java.sql.{Connection, Statement}

object TestObject {
  def main(args: Array[String]): Unit = {
    val a1 = new A
    a1.sayHello()

    val a2 = new A
    a2.sayHello()
  }
}


object A {
  var a: Int = init()

  def init(): Int = {
    println("初始化")
    99
  }
}

class A {
  val aa = A.a
  println("xixi")

  def sayHello(): Unit = {
    println(s"a:${aa}")
  }
}


object TestJdbc {
  // JDBC driver name and database URL// JDBC driver name and database URL
  val JDBC_DRIVER = "com.mysql.jdbc.Driver"
  val DB_URL = "jdbc:mysql://adtd.biz_lvpi.rds.sogou:3306/adtd_biz_lvpi?user=adtd&password=noSafeNoWork2016&characterEncoding=cp1252"

  //  Database credentials
  val USER = "adtd"
  val PASS = "noSafeNoWork2016"

  def main(args: Array[String]): Unit = {
    import java.sql.DriverManager
    import java.sql.SQLException
    var conn: Connection = null
    var stmt: Statement = null
    try {
      //STEP 2: Register JDBC driver
      Class.forName("com.mysql.jdbc.Driver")
      //STEP 3: Open a connection
      System.out.println("Connecting to database...")
      conn = DriverManager.getConnection(DB_URL, USER, PASS)
      //STEP 4: Execute a query
      System.out.println("Creating statement...")
      stmt = conn.createStatement
      var sql = "SELECT type,first_page_pv3 from adtd_biz_lvpi.rpt_adflow_ws_bczs_d limit 10"
      val rs = stmt.executeQuery(sql)
      //STEP 5: Extract data from result set
      while (rs.next) { //Retrieve by column name
        val typeF = rs.getString("type")
        val first_page_pv3 = rs.getString("first_page_pv3")
        //Display values
        System.out.print("typeF: " + typeF)
        System.out.print(", first_page_pv3: " + first_page_pv3)
        println("")
      }
      //STEP 6: Clean-up environment
      rs.close()
      stmt.close
      conn.close
    } catch {
      case se: SQLException =>
        //Handle errors for JDBC
        se.printStackTrace()
      case e: Exception =>
        //Handle errors for Class.forName
        e.printStackTrace()
    } finally {
      //finally block used to close resources
      try if (stmt != null) stmt.close
      catch {
        case se2: SQLException =>

      } // nothing we can do

      try if (conn != null) conn.close
      catch {
        case se: SQLException =>
          se.printStackTrace()
      } //end finally try

    } //end try

    System.out.println("Goodbye!")
  }
}
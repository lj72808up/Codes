package utils

import java.sql.{Connection, DriverManager, ResultSet}

object MysqlUtil extends PrintLogger {

  /**
    * 测试数据库连接，能连接返回true，不能连接返回false，不抛出异常
    * 参数为HOST, PORT, USERNAME, PASSWORD, DBNAME
    **/
  def tryConnection(host: String, port: String, username: String, passwd: String, dbName: String): Boolean = {

    try {
      //      val conn = DriverManager.getConnection(s"jdbc:mysql://$host:$port/$dbName?user=$username&password=$passwd&useSSL=false")
      val conn = JdbcHelper.jonesDataSource.getConnection
      conn.close()
      true
    } catch {
      case ex: Exception =>
        ex.printStackTrace()
        false
    }
  }

  /**
    * 测试数据库连接，能连接返回true，不能连接返回false，不抛出异常
    * 参数为connection_str
    **/
  def tryConnection(connStr: String): Boolean = {
    try {
      //      val conn = DriverManager.getConnection(connStr)
      val conn = JdbcHelper.jonesDataSource.getConnection
      conn.close()
      true
    } catch {
      case ex: Exception =>
        ex.printStackTrace()
        false
    }
  }

  /**
    * 返回一个连接，失败抛出异常
    * 参数为HOST, PORT, USERNAME, PASSWORD, DBNAME
    **/
  def getConnection(host: String, port: String, username: String, passwd: String, dbName: String): Connection = {
    try {
      //      DriverManager.getConnection(s"jdbc:mysql://$host:$port/$dbName?user=$username&password=$passwd&useSSL=false")
      JdbcHelper.jonesDataSource.getConnection
    } catch {
      case ex: Exception => throw ex
    }
  }

  /**
    * 返回一个连接，失败抛出异常
    * 参数为connection_str
    **/
  def getConnection(connStr: String): Connection = {
    try {
      //      DriverManager.getConnection(connStr)
      JdbcHelper.jonesDataSource.getConnection
    } catch {
      case ex: Exception => throw ex
    }
  }

  /**
    * 执行query，传入sql语句，失败抛出异常
    **/
  def executeQuery(conn: Connection, sql: String): ResultSet = {

    try {
      val statement = conn.createStatement(ResultSet.TYPE_FORWARD_ONLY, ResultSet.CONCUR_READ_ONLY)
      statement.executeQuery(sql)
    } catch {
      case ex: Exception => throw ex
    }
  }


  /**
    * 执行query，传入prepareStatement语句，和对应的值，失败抛出异常
    **/
  def executeQuery(conn: Connection, sql: String, params: Any*): ResultSet = {
    try {
      val prepareStatement = conn.prepareStatement(sql)
      params.zipWithIndex.foreach {
        case (param, index) =>
          val x = param
          prepareStatement.setObject(index + 1, param)
      }
      //      info(prepareStatement.toString)
      prepareStatement.executeQuery()

    } catch {
      case ex: Exception => throw ex
    }
  }

  /**
    * 执行update，传入sql语句，失败抛出异常
    **/
  def executeUpdate(conn: Connection, sql: String): Int = {

    try {
      val statement = conn.createStatement()
      statement.executeUpdate(sql)
    } catch {
      case ex: Exception => throw ex
    }

  }

  /**
    * 执行update，传入prepareStatement语句，和对应的值，失败抛出异常
    **/
  def executeUpdate(conn: Connection, sql: String, params: Any*): Int = {

    info(conn.getMetaData.getConnection.getMetaData.getURL)
    info(s"sql:{${sql}}  ,  param:${params}")
    try {
      val prepareStatement = conn.prepareStatement(sql)
      params.zipWithIndex.foreach {
        case (param, index) =>
          val x = param
          prepareStatement.setObject(index + 1, param)
      }
      prepareStatement.executeUpdate()

    } catch {
      case ex: Exception => ex.printStackTrace(); throw ex
    }

  }

  /**
    * 执行update，传入prepareStatement语句，和对应的值，失败抛出异常
    **/
  def executeInsert(conn: Connection, sql: String, params: Any*): Int = {

    try {
      val prepareStatement = conn.prepareStatement(sql)
      params.zipWithIndex.foreach {
        case (param, index) =>
          val x = param
          prepareStatement.setObject(index + 1, param)
      }
      prepareStatement.executeUpdate()

    } catch {
      case ex: Exception => throw ex
    }

  }


  /**
    * 关闭一个连接
    **/
  def closeConnection(conn: Connection) {
    conn.close()
  }

  /**
    * 关闭一个结果集
    **/
  def closeResultSet(rs: ResultSet): Unit = {
    rs.close()
  }

  def getTaskNameById(queryId: String): String = {
    val conn = JdbcHelper.jonesDataSource.getConnection
    val sql = "SELECT query_name FROM query_submit_log WHERE query_id = ?"
    val pstmt = conn.prepareStatement(sql)
    pstmt.setString(1, queryId)
    val rs = pstmt.executeQuery()
    var queryName = ""
    if (rs.next()) {
      queryName = rs.getString(1)
    }
    rs.close()
    conn.close()
    println(s"查询名字: $sql => $queryId => 结果为: $queryName")
    queryName
  }

}

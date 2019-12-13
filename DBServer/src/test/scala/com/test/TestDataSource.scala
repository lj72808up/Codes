package com.test

import java.sql.ResultSetMetaData

import com.sogou.datasource.DateSourceBuilder

object TestDataSource {
  def main(args: Array[String]): Unit = {

    def getFiledRes(columnIndex:Integer,metaData: ResultSetMetaData): Unit ={
      val columnType = metaData.getColumnType(columnIndex)
      columnType match {
        case int =>
      }
    }

    val dataSource = new DateSourceBuilder().build()

    import java.sql.SQLException
    try { // 获得连接:
      val conn = dataSource.getConnection
      // 编写SQL：
      val sql = "select * from scheduler_job where task_id=?"
      val pstmt = conn.prepareStatement(sql)
      //索引从1开始
      pstmt.setLong(1, 1)
      // 执行sql:
      val rs = pstmt.executeQuery
      val metaData = rs.getMetaData
      val columnCount = metaData.getColumnCount

      for (i <- 1 to columnCount){
        print(s"${metaData.getColumnName(i)}:${metaData.getColumnType(i)}\t")
      }
      println()

      while (rs.next()){
        for (i <- 1 to columnCount)
          print(rs.getObject(i)+"\t")
      }

      pstmt.close()
      conn.close()
    } catch {
      case e: SQLException =>
        e.printStackTrace()
    }
  }
}

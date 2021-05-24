package com.test

import java.io.{BufferedOutputStream, FileOutputStream, PrintWriter}
import java.util

import com.alibaba.druid.filter.Filter
import com.alibaba.druid.filter.stat.StatFilter
import com.alibaba.druid.pool.DruidDataSource
import com.typesafe.config.Config
import javax.sql.DataSource
import models.config.{JdbcConf, MainConf}
import utils.{ConfigUtil, JdbcHelper, PrintLogger}

object TestClickhouse {
  def main(args: Array[String]): Unit = {
    testJdbc()
    /*val url = "jdbc:clickhouse://ck.adtech.sogou:8123/"
    val user = "default"
    val passwd = "6lYaUiFi"
    val ds = buildCKDS(url, user, passwd)
    val conn = ds.getConnection
    val pstmt = conn.prepareStatement("show databases")
    val rs = pstmt.executeQuery
    while (rs.next()) {
      println(rs.getString(1))
    }*/
  }

  def testJdbc(): Unit = {
    val url = "jdbc:mysql://10.139.36.81:3306/adtl_test?useSSL=false&useUnicode=true&characterEncoding=UTF-8"
    val user = "yuhancheng"
    val passwd = "yuhancheng"
    val ds = buildMysqlDS(url, user, passwd)
    val conn = ds.getConnection
    val pstmt = conn.prepareStatement(s"explain select * from connection_manager")
    val rs = pstmt.executeQuery
    /*while (rs.next()) {
      println(rs.getString(1))
    }*/
  }


  def buildMysqlDS(url: String, user: String, passwd: String): DataSource = {
    val dataSource = new DruidDataSource()

    dataSource.setUrl(url)
    dataSource.setUsername(user)
    dataSource.setPassword(passwd)
    dataSource.setMaxActive(1)
    dataSource.setInitialSize(1)
    dataSource.setMaxWait(1)
    dataSource.setMinIdle(1)

    dataSource
  }

  def buildCKDS(url: String, user: String, passwd: String): DataSource = {
    val dataSource = new DruidDataSource()

    dataSource.setUrl(url)
    dataSource.setUsername(user)
    dataSource.setPassword(passwd)
    dataSource.setMaxActive(1)
    dataSource.setInitialSize(1)
    dataSource.setMaxWait(1)
    dataSource.setMinIdle(1)

    dataSource
  }
}


object TestDruid extends PrintLogger {
  def initConf(): Unit = {
    val confPath = "application.conf"
    val conf: Config = ConfigUtil.confFromFile(confPath)
    MainConf.addConf("basicConf", MainConf(conf))
  }

  def main3(args: Array[String]): Unit = {
    initConf()

    //val helper = new JdbcHelper()

    /*    val params = new util.ArrayList[Any]()
        params.add("全部")
        params.add("20180926")*/
    //val writer = new PrintWriter(new BufferedOutputStream(new FileOutputStream("test.file")))
    //helper.execSql("mysql", 1, "select 1", writer.write, writer.write)
    //    helper.execSql("select * from small_flow_res where pw=? and dt=? ",params, print, print)
    //writer.close()
    info("finish")
  }


  /** 导库小工具  */
  def main11(args: Array[String]): Unit = {
    val encodingMap = new util.HashMap[String, String]()
    encodingMap.put("utf8", "UTF-8")
    encodingMap.put("latin1", "Cp1252")

    val srcDS = new DruidDataSource()
    srcDS.setUrl("jdbc:mysql://adtd.jones.rds.sogou:3306/adtd_jones?useSSL=false&useUnicode=true&characterEncoding=UTF-8")
    srcDS.setUsername("adtd")
    srcDS.setPassword("adtdJonesApi")
    srcDS.setMaxActive(2)
    srcDS.setInitialSize(1)
    srcDS.setMaxWait(60000)
    srcDS.setMinIdle(1)

    val distDS = new DruidDataSource()
    distDS.setUrl("jdbc:mysql://adtd.datacenter.rds.sogou:3306/datacenter?useSSL=false&characterEncoding=UTF-8")
    distDS.setUsername("adtd")
    distDS.setPassword("noSafeNoWork2016")
    distDS.setMaxActive(2)
    distDS.setInitialSize(1)
    distDS.setMaxWait(60000)
    distDS.setMinIdle(1)


    val srcConn = srcDS.getConnection
    val selectSql =
      """
        |select t1.table_a_id, t2.table_name as table_a_name,
        |	t1.table_b_id,t3.table_name as table_b_name
        |from table_relation t1
        |left join table_meta_data t2 on t1.table_a_id=t2.tid
        |left join table_meta_data t3 on t1.table_b_id=t3.tid
        |""".stripMargin
    val srcPS = srcConn.prepareStatement(selectSql)
    val rs = srcPS.executeQuery

    val distConn = distDS.getConnection

    while (rs.next()) {
      val tableAName = rs.getString(2)
      val tableBName = rs.getString(4)
      println((s"$tableAName :: $tableBName"))

      val distPS1 = distConn.prepareStatement(s"select tid from table_meta_data where table_name=?")
      distPS1.setString(1, tableAName)
      val rs1 = distPS1.executeQuery()

      val distPS2 = distConn.prepareStatement(s"select tid from table_meta_data where table_name=?")
      distPS2.setString(1, tableBName)
      val rs2 = distPS2.executeQuery()

      if (rs1.next() && rs2.next()) {
        val t1id = rs1.getInt(1)
        val t2id = rs2.getInt(1)

        val distPS3 = distConn.prepareStatement(
          s"""
             |insert into table_relation (
             |  table_a_id,table_b_id
             |) values
             |(?,?)
             |""".stripMargin
        )
        distPS3.setInt(1, t1id)
        distPS3.setInt(2, t2id)
        distPS3.execute()
      }
    }
  }

  /** 导库小工具: relation_meta_data  */
  def main4(args: Array[String]): Unit = {

    val srcDS = new DruidDataSource()
    srcDS.setUrl("jdbc:mysql://adtd.jones.rds.sogou:3306/adtd_jones?useSSL=false&useUnicode=true&characterEncoding=UTF-8")
    srcDS.setUsername("adtd")
    srcDS.setPassword("adtdJonesApi")
    srcDS.setMaxActive(2)
    srcDS.setInitialSize(1)
    srcDS.setMaxWait(60000)
    srcDS.setMinIdle(1)

    val distDS = new DruidDataSource()
    distDS.setUrl("jdbc:mysql://adtd.datacenter.rds.sogou:3306/datacenter?useSSL=false&characterEncoding=UTF-8")
    distDS.setUsername("adtd")
    distDS.setPassword("noSafeNoWork2016")
    distDS.setMaxActive(2)
    distDS.setInitialSize(1)
    distDS.setMaxWait(60000)
    distDS.setMinIdle(1)


    val srcConn = srcDS.getConnection
    val selectSql =
      """
        |select
        |	table_a_name,
        |	field_a_name,
        |	table_b_name,
        |	field_b_name,
        |	description
        |from (
        |select t1.table_a_id,
        |	t2.table_name as table_a_name,
        |	t1.field_a_id,
        |	t3.field_name as field_a_name,
        |	t1.table_b_id,
        |	p2.table_name as table_b_name,
        |	t1.field_b_id,
        |	p3.field_name as field_b_name,
        |	t1.description
        |from relation_meta_data t1
        |left join table_meta_data t2 on t1.table_a_id=t2.tid
        |left join field_meta_data t3 on t1.field_a_id=t3.fid
        |left join table_meta_data p2 on t1.table_b_id=p2.tid
        |left join field_meta_data p3 on t1.field_b_id=p3.fid
        |) a1
        |where table_a_name is not null
        |  and field_a_name is not NULL
        |  and table_b_name is not null
        |  and field_b_name is not NULL
        |""".stripMargin
    val srcPS = srcConn.prepareStatement(selectSql)
    val rs = srcPS.executeQuery

    val distConn = distDS.getConnection

    while (rs.next()) {
      val t1Name = rs.getString(1)
      val f1Name = rs.getString(2)
      val t2Name = rs.getString(3)
      val f2Name = rs.getString(4)
      val desc = rs.getString(5)

      val distPS1 = distConn.prepareStatement(s"select tid from table_meta_data where table_name=?")
      distPS1.setString(1, t1Name)
      val rs1 = distPS1.executeQuery()


      val distPS3 = distConn.prepareStatement(s"select tid from table_meta_data where table_name=?")
      distPS3.setString(1, t2Name)
      val rs3 = distPS3.executeQuery()


      if (rs1.next() && rs3.next()) {
        val t1id = rs1.getInt(1)
        val t2id = rs3.getInt(1)

        val distPS2 = distConn.prepareStatement(s"select fid,version from field_meta_data where field_name=? and table_id=?")
        distPS2.setString(1, f1Name)
        distPS2.setInt(2, t1id)
        val rs2 = distPS2.executeQuery()

        val distPS4 = distConn.prepareStatement(s"select fid,version from field_meta_data where field_name=? and table_id=?")
        distPS4.setString(1, f2Name)
        distPS4.setInt(2, t2id)
        val rs4 = distPS4.executeQuery()

        if (rs2.next() && rs4.next()) {
          val f1id = rs2.getInt(1)
          val f1version = rs2.getString(2)

          val f2id = rs4.getInt(1)
          val f2version = rs4.getString(2)

          println(s"$t1id :: $f1id ===>  $t2id :: $f2id")
          //        Thread.sleep(500)

          val distPS9 = distConn.prepareStatement(
            s"""
               |insert into relation_meta_data (
               |  table_a_id,field_a_id,version_a,
               |  table_b_id,field_b_id,version_b,
               |  description
               |) values
               |(?,?,?,?,?,?,?)
               |""".stripMargin)
          distPS9.setInt(1, t1id)
          distPS9.setInt(2, f1id)
          distPS9.setString(3, f1version)
          distPS9.setInt(4, t2id)
          distPS9.setInt(5, f2id)
          distPS9.setString(6, f2version)
          distPS9.setString(7, desc)
          distPS9.execute()
        }
      }
    }
  }
}
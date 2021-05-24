package utils

import java.sql.{Connection, DriverManager, ResultSet, ResultSetMetaData}
import java.util

import adtl.platform.Dataproto.QueryInfo
import com.alibaba.druid.filter.Filter
import com.alibaba.druid.filter.stat.StatFilter
import com.alibaba.druid.pool.DruidDataSource
import javax.sql.DataSource
import models.config.{JdbcConf, MainConf}
import play.api.libs.json.{Format, Json}

sealed trait DataSourceType

case class JonesDataSource() extends DataSourceType

case class WorkSheetDataSource() extends DataSourceType

case class DBDescriptor(user: String, passwd: String, url: String)

class JdbcHelper extends PrintLogger {

  def checkDataSource(dsType: String, dsId: Long): Boolean = {
    var isExist = false
    dsType match {
      case "mysql" =>
        isExist = JdbcHelper.mysqlMapping.get(dsId) != null
      case "clickhouse" =>
        isExist = JdbcHelper.ckMapping.get(dsId) != null
      case "hive" =>
        isExist = true
      case "kylin" =>
        if (dsId != 0){
          isExist = JdbcHelper.kylinMapping.get(dsId) != null
        } else {
          isExist = JdbcHelper.kylinMapping.entrySet().size()>0
        }

    }
    isExist
  }

  def execSql(dsType: String,
              queryInfo: QueryInfo,
              metaFun: (Seq[String]) => Unit,
              rsFun: (ResultSet, Int, Seq[String]) => Unit,
              params: util.ArrayList[Any] = new util.ArrayList[Any]()
             ): Unit = {
    if (!checkDataSource(dsType, queryInfo.dsid)) {
      if (!JdbcHelper.addDataSourceById(queryInfo.dsid)) {
        throw new RuntimeException(s"数据源id:${queryInfo.dsid}不存在")
      }
    }
    var conn: Connection = null
    dsType match {
      case "mysql" =>
        conn = JdbcHelper.mysqlMapping.get(queryInfo.dsid).getConnection
        info(s"[INFO] ${conn.getMetaData.getConnection.getMetaData.getURL}, ${queryInfo.sql.head}")
      case "clickhouse" =>
        conn = JdbcHelper.ckMapping.get(queryInfo.dsid).getConnection
        info(s"[INFO] ${conn.getMetaData.getConnection.getMetaData.getURL}, ${queryInfo.sql.head}")
      case "hive" =>
        Class.forName("org.apache.hive.jdbc.HiveDriver")
        conn = DriverManager.getConnection("jdbc:hive2://rsync.master09.saturn.hadoop.js.ted:10005/default;principal=hive/hive@SATURN.HADOOP.COM")
        val stmt = conn.createStatement()
        stmt.execute(s"Set mapreduce.job.queuename=hive")
        stmt.execute(s"Set mapreduce.job.name=${queryInfo.queryId}")
      case "kylin" => // kylin
        if (queryInfo.dsid != 0){
          conn = JdbcHelper.kylinMapping.get(queryInfo.dsid).getConnection
        } else {
          conn = JdbcHelper.kylinMapping.entrySet().iterator().next().getValue.getConnection
        }
        info(s"[INFO] ${conn.getMetaData.getConnection.getMetaData.getURL}, ${queryInfo.sql.head}")

      case _ => info(s"this type: [$dsType], cannot use JdbcHelp")
    }

    val pstmt = conn.prepareStatement(queryInfo.sql.head)

    //索引从1开始
    for (i <- 1 to params.size()) {
      pstmt.setObject(i, params.get(i - 1))
    }

    // 执行sql:
    val rs: ResultSet = pstmt.executeQuery
    val metaData: ResultSetMetaData = rs.getMetaData
    val columnCount = metaData.getColumnCount

    val header = if (queryInfo.dimensions.size == 0 && queryInfo.outputs.size == 0) {
      (1 to metaData.getColumnCount).map(
        idx => metaData.getColumnName(idx)
      )
    } else {
      queryInfo.dimensions.toList ::: queryInfo.outputs.toList
    }

    metaFun(header)
    rsFun(rs, columnCount, header)
    pstmt.close()
    conn.close()
  }
}

object JdbcHelper extends PrintLogger {
  private val jonesConf = MainConf("basicConf").jdbc

  lazy val jonesDataSource: DataSource = buildJonesDS(jonesConf)
  val mysqlMapping: util.HashMap[Long, DataSource] = initMysqlDS()
  val ckMapping: util.HashMap[Long, DataSource] = initClickHouseDS()
  val kylinMapping: util.HashMap[Long, DataSource] = initKylinDS() // 初始化麒麟的数据源


  def buildJonesDS(conf: JdbcConf): DataSource = {
    buildDS(conf, conf.url, conf.user, conf.passwd)
  }

  def buildDS(conf: JdbcConf, url: String, user: String, passwd: String): DataSource = {
    val dataSource = new DruidDataSource()
    val isEnableMonitor = false
    //setName的意义在于，如果存在多个数据源，监控的时候可以通过名字来区分开来。如果没有配置，将会生成一个名字
    //    dataSource.setName("testds")
    dataSource.setUrl(url)
    dataSource.setUsername(user)
    dataSource.setPassword(passwd)
    //    dataSource.setFilters("config")
    dataSource.setMaxActive(conf.maxActive)
    dataSource.setInitialSize(conf.initialSize)
    dataSource.setMaxWait(conf.maxWait)
    dataSource.setMinIdle(conf.minIdle)

    val filters = new util.ArrayList[Filter]()
    if (isEnableMonitor) {
      val statFilter = new StatFilter()
      statFilter.setLogSlowSql(true) // 记录慢查询
      statFilter.setMergeSql(true) // 监控显示合并sql
      statFilter.setSlowSqlMillis(2000)

      filters.add(statFilter)
      dataSource.setProxyFilters(filters)
    }
    dataSource
  }

  def initClickHouseDS(): util.HashMap[Long, DataSource] = {
    val conn = JdbcHelper.jonesDataSource.getConnection
    val checkSql = "select dsid,user,passwd,url from connection_manager where conn_type='clickhouse'"
    val pstmt = conn.prepareStatement(checkSql)
    val rs = pstmt.executeQuery()
    val dataSources = new util.HashMap[Long, DataSource]
    while (rs.next()) {
      val ds = buildDS(jonesConf, rs.getString(4), rs.getString(2), rs.getString(3))
      dataSources.put(rs.getLong(1), ds)
    }
    pstmt.close()
    conn.close()
    dataSources
  }

  def initMysqlDS(): util.HashMap[Long, DataSource] = {
    val conn = JdbcHelper.jonesDataSource.getConnection
    val checkSql = "select dsid,user,passwd,url from connection_manager where conn_type='mysql'"
    val pstmt = conn.prepareStatement(checkSql)
    val rs = pstmt.executeQuery()
    val dataSources = new util.HashMap[Long, DataSource]
    while (rs.next()) {
      val ds = buildDS(jonesConf, rs.getString(4), rs.getString(2), rs.getString(3))
      dataSources.put(rs.getLong(1), ds)
    }
    pstmt.close()
    conn.close()
    dataSources
  }

  def initKylinDS(): util.HashMap[Long, DataSource] = {
    val conn = JdbcHelper.jonesDataSource.getConnection
    val checkSql = "select dsid,user,host,port,`database`,passwd from connection_manager where conn_type='kylin'"
    val pstmt = conn.prepareStatement(checkSql)
    val rs = pstmt.executeQuery()
    val dataSources = new util.HashMap[Long, DataSource]
    while (rs.next()) {
      // url, user, passwd
      val jdbcUrl = s"jdbc:kylin://${rs.getString(3)}:${rs.getString(4)}/${rs.getString(5)}"
      val ds = buildDS(jonesConf, jdbcUrl, rs.getString(2), rs.getString(6))
      dataSources.put(rs.getLong(1), ds)
    }
    pstmt.close()
    conn.close()
    dataSources
  }

  def addDataSourceById(dsid: Long): Boolean = {
    info(s"add datasource: $dsid")
    val conn = JdbcHelper.jonesDataSource.getConnection
    val sql = "select conn_type,user,passwd,url from connection_manager where dsid = ?"
    val pstmt = conn.prepareStatement(sql)
    pstmt.setLong(1, dsid)
    val rs = pstmt.executeQuery()
    rs.next()
    val ds = buildDS(jonesConf, rs.getString(4), rs.getString(2), rs.getString(3))
    var flag = false
    rs.getString(1) match {
      case "mysql" =>
        mysqlMapping.put(dsid, ds)
        flag = true
      case "clickhouse" =>
        ckMapping.put(dsid, ds)
        flag = true
      case _@x =>
        info(s"can't add datasource ${dsid}, unknown connType ${x}")
        flag = false
    }
    flag
  }
}
package utils

import java.security.MessageDigest
import java.sql.ResultSet

import models.QueryInfoModel.conn_str
import org.apache.spark.sql.{DataFrame, SparkSession}
import play.api.libs.json.{Format, Json}

class LinageHelper extends PrintLogger {
  val spark: SparkSession = SparkUtil.getSparkSession()

  def parseLinage(sql: String, defaultDB: String, rolesStr: String): Array[Node] = {
    val upStream = SparkUtil.parseTable(sql)
    var resultTable = Array[Node]() // 过滤后真正上游的表
    var waitImport = Array[Node]() // 元数据中没有, 但是hive中有 ,需要被导入的表
    // 过滤不存在的表
    for (tbIdentifier <- upStream) {
      var db = defaultDB
      var tb = tbIdentifier
      val splits = tbIdentifier.split("\\.")
      if (splits.size >= 2) {
        db = splits(0)
        tb = splits(1)
      }
      val astarFlag = checkTableInDB(db, tb, "hive")
      if (!astarFlag) {
        val hiveFlag = checkTableInHive(tbIdentifier)
        if (hiveFlag) {
          info(s"$tbIdentifier 在元数据中不存在, 但是在hive中存在, 需要导入")
          waitImport = waitImport :+ Node("hive", db, tb)
          resultTable = resultTable :+ Node("hive", db, tb)
        } else {
          info(s"$tbIdentifier 在数据库和hive中都不存在, 是sql中的临时表")
        }
      } else {
        info(s"$tbIdentifier 在元数据中已存在, 作为血缘中的上游表")
        resultTable = resultTable :+ Node("hive", db, tb)
      }
    }

    // 插入还没在元数据库中的表
    if (waitImport != null && waitImport.length > 0) {
      insertMetaDataTable(waitImport, rolesStr)
    }

    info(s"需要导入的表为: ${waitImport.map(t => t.db + "." + t.tableName).mkString(", ")}")
    info(s"血缘解析的上游表为: ${resultTable.map(t => t.db + "." + t.tableName).mkString(", ")}")

    resultTable
  }

  private def checkTableInDB(db: String, tb: String, tbType: String): Boolean = {
    val conn = MysqlUtil.getConnection("dummy conn string")
    var flag = false
    tbType match {
      case "hive" => {
        var sql = ""
        var rs: ResultSet = null
        if ("".equals(db)) {
          sql = s"select * from table_meta_data where table_name=? and database_name=? and table_type=?"
          val ps = conn.prepareStatement(sql)
          ps.setString(1, tb)
          ps.setString(2, db)
          ps.setString(3, tbType)
          rs = ps.executeQuery()
        } else {
          sql = s"select * from table_meta_data where table_name=? and table_type=?"
          val ps = conn.prepareStatement(sql)
          ps.setString(1, tb)
          ps.setString(2, tbType)
          rs = ps.executeQuery()
        }
        flag = rs.next()
        rs.close()
      }
      case _ => info(s"暂不处理非hive的表: ${db}.${tb}")
    }
    MysqlUtil.closeConnection(conn)
    flag
  }

  private def checkTableInHive(tableIdentifier: String): Boolean = {
    val sql = s"desc $tableIdentifier"
    val spark = SparkUtil.getSparkSession()
    var flag = true
    try {
      spark.sql(sql)
    } catch {
      case e: Exception => flag = false
    }
    flag
  }

  private def insertMetaDataTable(nodes: Array[Node], rolesStr: String): Unit = {
    val hiveDsId = getHiveConnId()
    for (node <- nodes) {
      // 1. 每个 node 先检查数据库中是否存在
      if (checkTableInHive(s"${node.db}.${node.tableName}")) {
        val protoTable = new HiveMetaDataHelper(spark).getTable(node.db, node.tableName)

        val compress = protoTable.expendInfo("storageProperties")
        // todo 自动导入表后要插入权限
        val tableMetaData = TableMetaData(node.tableName, "hive", node.db, protoTable.expendInfo("owner"),
          protoTable.expendInfo("createdTime"), protoTable.expendInfo("location"), compress,
          getTableIdentifier(node.tableType, node.db, node.tableName), rolesStr, hiveDsId
        )

        insertMetaDataDB(tableMetaData, rolesStr)
      }

    }
  }

  private def insertMetaDataDB(node: TableMetaData, rolesStr: String): Unit = {
    val conn = MysqlUtil.getConnection("dummy conn string")
    info(s"${node.identifier}在数据库中还不存在, 开始导入. ")
    val sql =
      s"""INSERT INTO table_meta_data (
         |  table_name, table_type, database_name, location, creator, create_time, description, compress_mode, identifier, owner, dsid
         |)VALUES (
         |  ?, ?, ?, ?, ?, ?,"", ?, ?, ?, ?
         |)""".stripMargin

    val pstmt = conn.prepareStatement(sql)
    pstmt.setString(1, node.tableName)
    pstmt.setString(2, node.tableType)
    pstmt.setString(3, node.databaseName)
    pstmt.setString(4, node.location)
    pstmt.setString(5, node.creator)
    pstmt.setString(6, node.createTime)
    pstmt.setString(7, node.compressMode)
    pstmt.setString(8, node.identifier)
    pstmt.setString(9, node.owner)
    pstmt.setInt(10, node.dsid)

    pstmt.addBatch()
    pstmt.executeBatch()

    // 查询表id
    val checkSql = s"SELECT tid FROM table_meta_data WHERE identifier = ? "
    val pstmt2 = conn.prepareStatement(checkSql)
    pstmt2.setString(1, node.identifier)
    val checkRs = pstmt2.executeQuery()
    checkRs.next()
    val tid = checkRs.getInt(1)
    // 插入权限
    if (tid != 0) {
      val roles = rolesStr.split(",")
      val permitSql =
        s"""INSERT INTO group_metatable_mapping(
           |  group_name, table_id, valid, read_flag, exec_flag
           |) VALUES (
           |  ?,?,?,?,?
           |)
           |""".stripMargin
      val permitPstmt = conn.prepareStatement(permitSql)
      for (role <- roles) {
        permitPstmt.setString(1, role)
        permitPstmt.setString(2, tid + "")
        permitPstmt.setBoolean(3, true)
        permitPstmt.setBoolean(4, true)
        permitPstmt.setBoolean(5, true)
        permitPstmt.addBatch()
      }
      permitPstmt.executeBatch()
    } else {
      throw new Exception(s"sparkserver cannot find ${node.databaseName}.${node.tableName}")
    }
    conn.close()
  }


  private def getHiveConnId(): Int = {
    val conn = MysqlUtil.getConnection("dummy conn string")
    val ps = conn.prepareStatement(s"SELECT dsid FROM connection_manager WHERE conn_type='hive' LIMIT 1")
    val rs = ps.executeQuery()
    rs.next()
    val dsId = rs.getInt(1)
    rs.close()
    conn.close()
    dsId
  }


  private def getTableIdentifier(tableType: String, databaseName: String, tableName: String): String = {
    val str = s"$tableType-$databaseName-$tableName"
    Md5Util.getMd5(str)
  }
}


case class Node(tableType: String, db: String, tableName: String)

object Node {
  implicit val jsonFormat: Format[Node] = Json.format[Node]
}

case class TableMetaData(tableName: String,
                         tableType: String,
                         databaseName: String,
                         creator: String,
                         createTime: String,
                         location: String,
                         compressMode: String,
                         identifier: String,
                         owner: String,
                         dsid: Int
                        ) {}


object Md5Util {
  def getMd5(content: String): String = {
    val md5 = MessageDigest.getInstance("MD5")
    val encoded = md5.digest((content).getBytes)
    encoded.map("%02x".format(_)).mkString
  }

  def main(args: Array[String]): Unit = {
    val tableType = "hive"
    val databaseName = "datacenter"
    val tableName = "test_1214_02"
    val md5Content = getMd5(s"$tableType-$databaseName-$tableName")
    print(md5Content == "be2cac69457bd7544d04a3f351e7a5e2")
  }
}
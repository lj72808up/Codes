package utils

import java.util

import javax.sql.DataSource
import org.apache.spark.sql.SparkSession
import org.apache.spark.sql.catalyst.analysis.UnresolvedRelation

object SparkUtil extends PrintLogger {

  def getSparkSession(): SparkSession = {
//    return null
    val session = SparkSession.getActiveSession.orElse(SparkSession.getDefaultSession)
    if (session.isEmpty) {
      throw new IllegalStateException("Cannot find active or default SparkSession in the current " +
        "context")
    }

    session.get
  }

  def getClassLoader(): ClassLoader ={
    Option(Thread.currentThread().getContextClassLoader).getOrElse(getClass.getClassLoader)
  }

  def getMeta(sql: String): Array[String] = {
    val spark = SparkUtil.getSparkSession()
    val sqlDF = spark.sql(sql)
    val sqlFields = sqlDF.schema.fields
    val sqlTypes: Array[String] = sqlFields.map(f => f.dataType.sql)
    sqlTypes
  }

  def checkMeta(sql: String, outTable: String, outType: String): (Boolean, String) = {
    val spark = SparkUtil.getSparkSession()
    val sqlTypes = getMeta(sql)
    var outTypes = Array[String]()

    outType match {
      case "hive" => {
        val outTableDF = spark.sql(s"SELECT * FROM $outTable")
        val outFields = outTableDF.schema.fields
        outTypes = outFields.map(f => f.dataType.sql)
      }
      case t if (t == "mysql") || (t == "clickhouse") => {
        val (db, tb) = transTableName(outTable)
        outTypes = getSchemaFromAStarDB(tb, db, outType)
      }
    }

    print("sql产生的column: ")
    sqlTypes.foreach(x => print(s"${x}=> "))
    print("\ntable schema: ")
    outTypes.foreach(x => print(s"${x}=> "))
    println()

    outType match {
      case "hive" => {
        // 排除分区字段 dt 或 hour
        if (sqlTypes.length != outTypes.length - 1)
           return (false, "sql字段个数和输出表字段个数不同")
      }
      case "mysql" => {
        if (sqlTypes.length != outTypes.length - 1)
          return (false, "sql字段个数和输出表字段个数不同")
      }
      case "clickhouse" => {
        if (sqlTypes.length != outTypes.length)
          return (false, "sql字段个数和输出表字段个数不同")
      }
    }
    (true, "")
  }

  private def transTableName(tableName: String): (String, String) = {
    var db = "default"
    var tb = tableName

    val splits = tableName.split("\\.")
    if (splits.length > 1) {
      db = splits(0)
      tb = splits(1)
    }

    (db, tb)
  }

  private def getSchemaFromAStarDB(tb: String, db: String, tableType: String): Array[String] = {
    info(s"check table: db:$db, tb:$tb, tableType:$tableType")
    var fTypeArr = Array[String]()
    val conn = JdbcHelper.jonesDataSource.getConnection
    val sql =
      s"""
         |SELECT p2.field_name,p2.field_type FROM
         |(
         |  SELECT t1.tid, t2.`version` FROM
         |  (SELECT tid FROM table_meta_data WHERE table_type='$tableType' AND TABLE_NAME='$tb' AND database_name='$db') t1
         |  LEFT JOIN version_meta_data t2
         |  ON t1.tid=t2.table_id
         |  GROUP BY t1.tid, t2.version
         |  ORDER BY t2.version DESC
         |  LIMIT 1
         |) p1
         |LEFT JOIN
         |field_meta_data p2
         |ON
         |  p1.tid = p2.`table_id` AND
         |  p1.`version` = p2.`version`
         |""".stripMargin

    info(sql)
    val pStmt = conn.prepareStatement(sql)
    val rs = pStmt.executeQuery()
    while (rs.next()) {
      val fieldName = rs.getString(1)
      val filedType = rs.getString(2)

      print(s"$fieldName($filedType),")
      fTypeArr = fTypeArr :+ filedType
    }
    println()
    pStmt.close()
    conn.close()
    fTypeArr
  }

  @deprecated
  def parseTableOld(sql: String): Seq[String] = {
    val spark = getSparkSession()
    val logicalPlan = spark.sessionState.sqlParser.parsePlan(sql)
    import org.apache.spark.sql.catalyst.analysis.UnresolvedRelation
    val tables = logicalPlan.collect { case r: UnresolvedRelation => r.tableName }
    tables
  }

  def parseTable(sql: String): Seq[String] = {
    val spark = getSparkSession()
    val logicalPlan = spark.sessionState.sqlParser.parsePlan(sql)
    val len = logicalPlan.numberedTreeString.split("\n").length
    (0 until len).flatMap {
      idx =>
        logicalPlan.p(idx).collect { case r: UnresolvedRelation => r.tableName }
    }.toSet.toSeq
  }
}

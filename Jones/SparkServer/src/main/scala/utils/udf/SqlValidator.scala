package utils.udf

import utils.{JdbcHelper, PrintLogger, SparkUtil}

object SqlValidator extends PrintLogger{
  def validateSql(sql: String, sqlType: String, dsId: Long): String = {
    sqlType match {
      case "mysql" => assertMysql(sql, dsId)
      case "clickhouse" => assertClickHouse(sql, dsId)
      case "hive" => assertHive(sql)
      case _ => {
        info(s"校验的 sql 类型不在候选中, sqlType: ${sqlType}")
        "sql verify"
      }
    }
  }

  private def assertHive(sql: String): String = {
    val spark = SparkUtil.getSparkSession()
    spark.sql(sql)
    "sql verify"
  }

  private def assertClickHouse(sql: String, dsId: Long): String = {
    val validateSql = s"explain sql"
    val conn = JdbcHelper.ckMapping.get(dsId).getConnection
    val pstmt = conn.prepareStatement(validateSql)
    pstmt.executeQuery
    // 测试完毕, 测试不通过会抛异常
    pstmt.close()
    conn.close()
    "sql verify"
  }

  private def assertMysql(sql: String, dsId: Long): String = {
    val validateSql = s"explain $sql"
    val conn = JdbcHelper.mysqlMapping.get(dsId).getConnection
    val pstmt = conn.prepareStatement(validateSql)
    pstmt.executeQuery
    // 测试完毕, 测试不通过会抛异常
    pstmt.close()
    conn.close()
    "sql verify"
  }
}

package utils.dag

import java.sql.{Connection, DriverManager}

import airflow.DagTask
import org.apache.spark.sql.{DataFrame, SparkSession}
import utils.{PrintLogger, SparkUtil}

object MySqlParser extends PrintLogger {
  val spark: SparkSession = SparkUtil.getSparkSession()

  def createTableFromSql(selectSql: String,
                         jdbcUrl: String,
                         tableName: String,
                         staticFields: Map[String, String],
                         comment: String,
                         extraFields: Seq[String]): Unit = {
    val df = spark.sql(selectSql)
    val fields = df.schema.fields
    var fieldList = List[String]()
    fields.foreach(
      f => {
        val fType = "TEXT" //todo 字段类型硬编码 f.dataType.sql
        val fName = f.name
        val fComment = staticFields.getOrElse(fName, "")

        fieldList = fieldList :+ s"`$fName` $fType comment '$fComment'"
      }
    )
    extraFields.foreach(f => {
      fieldList = fieldList :+ s"`$f` Text " // todo 字段类型硬编码
    })

    var handleTableName = tableName
    val createSql =
      s"""CREATE TABLE IF NOT EXISTS $handleTableName (
         | ${fieldList.mkString("\t", ",\n\t", "")}
         |)
         |COMMENT '${comment}'
        """.stripMargin

    info(createSql)

    val conn = getConnection(jdbcUrl)
    val statement = conn.createStatement
    statement.executeUpdate(createSql)
    statement.close()
    conn.close()
  }

  def getConnection(jdbcURL: String): Connection = {
    Class.forName("com.mysql.jdbc.Driver")
    val conn = DriverManager.getConnection(jdbcURL)
    conn
  }

  def createTableFromTasks(spark: SparkSession, tasks: List[DagTask]): Unit = {
    val extraFields = Array[String]("date") // mysql表用于区分天和小时的字段
    for (i <- tasks.indices) {
      if (tasks(i).newOrOldOutTable == "newTable") {
        createTableFromSql(tasks(i).sql, tasks(i).lvpi.jdbcStr, tasks(i).tableName, tasks(i).sqlFields.getOrElse(Map()),
          tasks(i).tableDesc, extraFields)
      }
    }
  }
}


case class JdbcOption(url: String,
                      user: String,
                      password: String)
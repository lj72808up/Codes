package utils.dag

import airflow.DagTask
import org.apache.spark.sql.SparkSession
import utils.{PrintLogger, SparkUtil}

object NodeParser extends PrintLogger {
  val spark: SparkSession = SparkUtil.getSparkSession()

  def createTableFromTasks(tasks: List[DagTask]): Unit = {
    var sparkOuts = List[DagTask]()
    var ckOuts = List[DagTask]()
    var mysqlOuts = List[DagTask]()

    for (task <- tasks) {
      val outType = task.outTableType
      outType match {
        case t if t == "spark" || t == "hive" => sparkOuts = sparkOuts :+ task
        case "clickhouse" => ckOuts = ckOuts :+ task
        case "mysql" => mysqlOuts = mysqlOuts :+ task
      }
    }

    if (sparkOuts.nonEmpty) {
      SparkParser.createTableFromTasks(spark, sparkOuts)
    }

    if (ckOuts.nonEmpty) {
      CkParser.createTableFromTasks(spark, ckOuts)
    }

    if (mysqlOuts.nonEmpty) {
      MySqlParser.createTableFromTasks(spark, mysqlOuts)
    }
  }
}

package airflow

import play.api.libs.json.{Format, Json}

// 对应的dag任务是: sql语句, 输出表名, 任务类型
case class DagTask(outTableType: String,
                   sql: String,
                   tableName: String,
                   newOrOldOutTable: String,
                   tableDesc: String,
                   storageEngine: String,
                   taskType: String,
                   inputType: String,
                   taskTypeInfo: TaskTypeInfo,
                   sqlFields: Option[Map[String, String]],
                   params: Option[Map[String, String]],
                   lvpi: LvpiTask,
                   clickhouse: CkTask) // 比如输出类型是hive时, 选择文件类型

case class LvpiTask(jdbcStr: String, tableName: String)

case class CkTask(tableName: String, indices: Option[Array[String]], partitions: Option[Array[String]])

case class TaskTypeInfo(dsid: Int, hdfsLocation: String, splitter: String, validFieldCount: Int,
                        schema: List[Map[String, String]], extendSchema: List[Map[String, String]], hasExtend: Boolean,
                        hdfsDoneLocation: String, isCheckDone: Boolean)

object DagTask {
  implicit val jsonFormat: Format[DagTask] = Json.format[DagTask]
}

object TaskTypeInfo {
  implicit val jsonFormat: Format[TaskTypeInfo] = Json.format[TaskTypeInfo]
}

object LvpiTask {
  implicit val jsonFormat: Format[LvpiTask] = Json.format[LvpiTask]
}

object CkTask {
  implicit val jsonFormat: Format[CkTask] = Json.format[CkTask]
}

sealed class TaskType

case class HiveTask() extends TaskType

case class MysqlTask() extends TaskType

case class ClickHouseTask() extends TaskType

// hive存储格式的枚举
object HiveStorage extends Enumeration {
  type HiveStorage = Value
  val Tsv: airflow.HiveStorage.Value = Value("tsv")
  val Orc_Snappy: airflow.HiveStorage.Value = Value("orc/snappy")
}
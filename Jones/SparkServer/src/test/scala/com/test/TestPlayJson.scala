package com.test

import adtl.platform.Dataproto.{FilterOperation, HiveField, HiveTable}
import airflow.DagTask
import play.api.libs.json.{Format, JsArray, Json, Reads, Writes}
import routes.AStarRes
import utils.HiveMetaDataHelper

import scala.util.parsing.json.JSONArray

object TestPlayJson {

  case class DagTask(outTableType: String,
                     sql: String,
                     tableName: String,
                     newOrOldOutTable: String,
                     tableDesc: String,
                     storageEngine: String,
                     taskType: String,
                     sqlFields: Option[Map[String, String]],
                     params: Option[Map[String, String]],
                     lvpi: LvpiTask,
                     clickhouse: CkTask) // 比如输出类型是hive时, 选择文件类型

  case class LvpiTask(jdbcStr: String, tableName: String)

  case class CkTask(tableName: String, indices: Option[Array[String]], partitions: Option[Array[String]])

  object DagTask {
    implicit val jsonFormat: Format[DagTask] = Json.format[DagTask]
  }

  object LvpiTask {
    implicit val jsonFormat: Format[LvpiTask] = Json.format[LvpiTask]
  }

  object CkTask {
    implicit val jsonFormat: Format[CkTask] = Json.format[CkTask]
  }

  def main(args: Array[String]): Unit = {
    val maps = Map[String,Map[String, String]](
      "aa"-> Map("aa2"-> "bb2"),
      "bb"-> Map("aa2"-> "bb2")
    )
    println(maps.toString())
  }

  def main10(args: Array[String]): Unit = {
    val maps = Map[String,DagTask](
      "aa"-> DagTask("outTableType","sql","tableName","newOrOldOutTable","tableDesc","storageEngine","taskType",
        None,None,LvpiTask("xxx","yyyy"),CkTask("tn",None,None))
    )
    println(Json.toJson(maps).toString())
  }

  def main9(args: Array[String]): Unit = {
    val body =
      """
        |[
        |    {
        |        "outTableType":"",
        |        "tableName":"haha.hehe",
        |        "newOrOldOutTable":"newTable",
        |        "tableDesc":"的方式发生",
        |        "storageEngine":"",
        |        "sql":"",
        |        "taskType":"",
        |        "lvpi":{
        |            "dsid":1,
        |            "jdbcStr":"jdbc:mysql://adta.kl_history.rds.sogou:noSafeNoWork2016/kl??user=admin&password=noSafeNoWork2016&useSSL=false&useUnicode=true&characterEncoding=UTF-8",
        |            "tableName":"haha.hehe"
        |        },
        |        "clickhouse":{
        |            "tableName":""
        |        }
        |    }
        |]
        |""".stripMargin
    val filterObj = Json.parse(body).as[JsArray]
    filterObj.as[List[DagTask]]
  }

  def main8(args: Array[String]): Unit = {
    val value = Json.stringify(Json.toJson(Map(
      "res" -> true
    )))
    println(value)
  }

  def main0(args: Array[String]): Unit = {
    val str =
      """
        |{"name":"zhangsan","age":"24"}
        |""".stripMargin
    val parseObj = Json.parse(str)
    val map = parseObj.as[Map[String, String]]
    print(map)
  }

  def main1(args: Array[String]): Unit = {
    println("aaa")
    val str = "{\"name\":\"CXC\",\"values\":[\"百度\",\"哈哈\",\"aa\"],\"type\":\"QW_AM\"}"
    val filterObj = Json.parse(str)
    val qwPartition = (filterObj \ "values").as[JsArray]
    val arr = qwPartition.as[List[String]]
    println(arr)
  }

  def main2(args: Array[String]): Unit = {
    case class A(name: String)
    implicit val locationWrites = new Writes[A] {
      def writes(a: A) = Json.obj(
        "name" -> a.name
      )
    }
    val jsVal = Json.toJson(A("aaa"))
    println(Json.prettyPrint(jsVal))
  }

  def main3(args: Array[String]): Unit = {
    val table = new HiveTable(
      "test1", integratedField = List[HiveField](
        HiveField("name", "string", "姓名"),
        HiveField("age", "int", "年龄")
      )
    )
    import HiveMetaDataHelper._
    val jsVal = Json.toJson(table)
    val jsValStr = Json.prettyPrint(jsVal)
    println(jsValStr)
  }

  def main4(args: Array[String]): Unit = {
    val res = Array("aaa", "bbb")
    println(Json.stringify(Json.toJson(res)))
  }

  def main5(args: Array[String]): Unit = {
    val a = List(List(1, 2, 3))
    val jsonVal = Json.toJson(a)
    println(Json.stringify(jsonVal))
  }

  def main6(args: Array[String]): Unit = {
    val str =
      """
        |[{
        | "sql": "select 1",
        | "tableName": "out1",
        | "taskType": "day"
        |}]
        |""".stripMargin
    val filterObj = Json.parse(str).as[JsArray]
    implicit val jsonFormat: Format[DagTask] = Json.format[DagTask]
    println(filterObj.as[List[DagTask]])
  }

  def main7(args: Array[String]): Unit = {
    val res = AStarRes(flag = true, "good")
    println(Json.stringify(Json.toJson(res)))
  }
}

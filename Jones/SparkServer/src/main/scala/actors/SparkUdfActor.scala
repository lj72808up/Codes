package actors
import java.util
import java.util.function.Consumer

import akka.actor.Actor
import models.HiveMetaInfo
import org.apache.spark.sql.SparkSession
import utils.{DateTimeUtil, MysqlUtil, PrintLogger, SparkUtil}

case class Start(funcMap: util.HashMap[String, util.HashMap[String, String]])

class SparkUdfActor extends Actor with PrintLogger {
  val spark: SparkSession = SparkUtil.getSparkSession()

  override def receive: Receive = {
    case Start(funcMap) => {
      info(s"sparkUdfActor 收到参数: $funcMap")
      info("开始获取spark udf")
      val startTime = DateTimeUtil.getTsNow()
      val onlineFuncs = getAllUdByInvokeRemote(funcMap)
      val dbMarkFuncs = getDBMarkFuncs()
      markNewFuncs(diffNewFuncs(onlineFuncs, dbMarkFuncs), funcMap)
      val endTime = DateTimeUtil.getTsNow()
      info(s"获取spark udf成功, 耗时: ${(endTime - startTime) / 1000}")
    }
  }

  private def getAllUdfFromSpark(): Array[String] = {
    val udf = spark.udf
    val clazz = udf.getClass
    val field = clazz.getDeclaredField("functionRegistry")
    field.setAccessible(true)
    val functionRegistry = field.get(spark.udf).asInstanceOf[org.apache.spark.sql.catalyst.analysis.SimpleFunctionRegistry]
    functionRegistry.listFunction().map(identifier => identifier.funcName).toArray[String]
  }

  private def getAllUdByInvokeRemote(funcMap: util.HashMap[String,util.HashMap[String,String]]): Array[String] = {
    var arr = Array[String]()
    funcMap.keySet().forEach(new Consumer[String] {
      override def accept(t: String): Unit = {
        arr = arr :+ t
      }
    })
    arr
  }

  /**
    * 对比真实的funcs和db中记录的funcs, 真实的funcs增加的部分
    */
  private def diffNewFuncs(onlineFuncs: Array[String], dbMarkFuncs: Array[String]): Array[String] = {
    val onlineSet = new java.util.HashSet[String]()
    val dbSet = new java.util.HashSet[String]()

    onlineFuncs.foreach(fun => onlineSet.add(fun))
    dbMarkFuncs.foreach(fun => dbSet.add(fun))

    var newFuncs: Array[String] = Array[String]()
    onlineSet.removeAll(dbSet)
    onlineSet.forEach(new Consumer[String] {
      override def accept(t: String): Unit = {
        newFuncs = newFuncs :+ t
      }
    })
    info(s"发现${newFuncs.length}个新方法")
    newFuncs
  }

  private def getDBMarkFuncs(): Array[String] = {
    val conn = MysqlUtil.getConnection("dummy conn string")
    val rs = MysqlUtil.executeQuery(conn, s"select func_name from udf_manager")
    var dbFuncs = Array[String]()
    while (rs.next()) {
      val funcName = rs.getString(1)
      dbFuncs = dbFuncs :+ funcName
    }
    dbFuncs
  }

  private def markNewFuncs(newFuncs: Array[String], funcMap: util.HashMap[String,util.HashMap[String,String]]): Unit = {
    val conn = MysqlUtil.getConnection("dummy conn string")
    val sql = "INSERT INTO udf_manager(func_name, known, description, parameters, example) VALUES (?,?,?,?,?)"
    val preparedStatement = conn.prepareStatement(sql)

    newFuncs.foreach(func => {
      preparedStatement.setString(1, func)
      preparedStatement.setBoolean(2, false)
      preparedStatement.setString(3, funcMap.get(func).get("funcDesc"))
      val content =
        s"""${funcMap.get(func).get("funcDef")}
           |函数类型: ${funcMap.get(func).get("funcType")}
           |参数: ${funcMap.get(func).get("funcParams")}
           |返回值: ${funcMap.get(func).get("funcRet")}
           |""".stripMargin
      preparedStatement.setString(4, content)
      preparedStatement.setString(5, funcMap.get(func).get("funcUsage"))
      preparedStatement.addBatch()
    })

    preparedStatement.executeBatch()
    conn.close()
  }
}



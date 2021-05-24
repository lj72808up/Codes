package routes

import java.io.{File, FileInputStream, FileOutputStream}
import java.text.SimpleDateFormat
import java.util.{Date, UUID}

import adtl.platform.Dataproto._
import airflow.DagTask
import akka.http.scaladsl.marshalling.ToResponseMarshallable
import akka.http.scaladsl.model.{ContentTypes, HttpEntity}
import akka.http.scaladsl.server.Directives.{complete, entity, extractDataBytes, onComplete, path, withoutRequestTimeout, withoutSizeLimit, _}
import org.apache.log4j.{Level, Logger}
import org.apache.spark.sql.SparkSession
import akka.http.scaladsl.unmarshalling.{FromEntityUnmarshaller, PredefinedFromEntityUnmarshallers}
import akka.stream.ActorMaterializer
import akka.util.ByteString
import akka.http.scaladsl.model.{HttpResponse, Multipart, StatusCodes}
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.Route

import scala.concurrent.duration._
import models.{HiveMetaInfo, ServerInfo, ServerStatus}
import org.joda.time.DateTime
import utils.{DateTimeUtil, HDFSUtil, HiveMetaDataHelper, JedisHelper, JsonUtil, PrintLogger, SmallFlowTaskLib, SparkCalcLib, SparkUtil, TestUtil, XSSFUtil}
import models.QueryInfoModel
import org.apache.hadoop.fs.{FileSystem, Path}
import utils.dag.{NodeParser, SparkParser}
//import shadeio.poi.xssf.usermodel.{XSSFCell, XSSFWorkbook}
import play.api.libs.json.{JsArray, JsValue, Json}

import scala.concurrent.ExecutionContext
import scala.io.Source

class HttpRoute(spark: SparkSession, logger: Logger)(implicit executionContext: ExecutionContext, materializer: ActorMaterializer) extends PrintLogger {

  import spark.implicits._
  import implicits._

  object implicits {

    implicit val queryInfoUnmarshaller: FromEntityUnmarshaller[QueryInfo] =
      PredefinedFromEntityUnmarshallers.byteArrayUnmarshaller map { bytes =>
        QueryInfo.parseFrom(bytes)
      }
    implicit val SmallFlowUnmarshaller: FromEntityUnmarshaller[SmallFlow] =
      PredefinedFromEntityUnmarshallers.byteArrayUnmarshaller map { bytes =>
        SmallFlow.parseFrom(bytes)
      }
    implicit val MapStringStringUnmarshaller: FromEntityUnmarshaller[Map[String, String]] =
      PredefinedFromEntityUnmarshallers.byteArrayUnmarshaller map {
        bytes => {
          val body = new String(bytes)
          Json.parse(body).as[Map[String, String]]
        }
      }
    implicit val MapStringListMapUnmarshaller: FromEntityUnmarshaller[List[Map[String, String]]] =
      PredefinedFromEntityUnmarshallers.byteArrayUnmarshaller map {
        bytes => {
          val body = new String(bytes)
          Json.parse(body).as[List[Map[String, String]]]
        }
      }
    implicit val DagTasksUnmarshaller: FromEntityUnmarshaller[List[DagTask]] =
      PredefinedFromEntityUnmarshallers.byteArrayUnmarshaller map {
        bytes => {
          val body = new String(bytes)
          val filterObj = Json.parse(body).as[JsArray]
          filterObj.as[List[DagTask]]
        }
      }

  }


  def getTranStream(path: String): TranStream = {

    val respath = ServerInfo.Respath
    val originPath = s"${respath}/${path}"
    //    val file = Source.fromFile("/" + path)
    val file = Source.fromFile("/" + originPath)
    val iter = file.getLines()
    val header = iter.next().split("\t", -1)

    val lines = (0 until 5000).map {
      _ =>
        if (iter.hasNext) iter.next()
        else ""
    }.filterNot(_.isEmpty).map {
      x => Linedata().withField(x.split("\t", -1))
    }

    TranStream().withLine(lines).withSchema(header)
  }

  def getQueryIdFromQueryPath(path: String): String = {
    val splits = path.split("/")
    val fileName = splits(splits.size - 1)
    fileName.split("\\.")(0)
  }

  val route =
    path("sql" / "upload") {
      try {
        entity(as[Multipart.FormData]) { (formdata: Multipart.FormData) =>
          var fName = "" // 最终存储文件的全路径
          val fileNamesFuture = formdata.parts.mapAsync(1) { p =>
            info(s"Got part. name: ${p.name} filename: ${p.filename.get}") // 上传的文件名
            info("Writing file and Counting size...")
            @volatile var lastReport = System.currentTimeMillis()
            @volatile var lastSize = 0L
            val root = new File(s"/tmp/jones/sqlUpload/${new SimpleDateFormat("yyyyMMdd").format(new Date)}")
            if (!root.exists()) root.mkdirs()
            fName = s"${root}/${p.filename.get}"
            val distFile = new FileOutputStream(fName)

            def receiveChunk(counter: (Long, Long), chunk: ByteString): (Long, Long) = {
              val (oldSize, oldChunks) = counter
              val newSize = oldSize + chunk.size
              val newChunks = oldChunks + 1

              val now = System.currentTimeMillis()
              if (now > lastReport + 1000) {
                val lastedTotal = now - lastReport
                val bytesSinceLast = newSize - lastSize
                val speedMBPS = bytesSinceLast.toDouble / 1000000 /* bytes per MB */ / lastedTotal * 1000 /* millis per second */

                info(f"Already got $newChunks%7d chunks with total size $newSize%11d bytes avg chunksize ${newSize / newChunks}%7d bytes/chunk speed: $speedMBPS%6.2f MB/s")

                lastReport = now
                lastSize = newSize
              }
              // 存储文件
              distFile.write(chunk.toByteBuffer.array())
              distFile.flush()

              (newSize, newChunks)
            }

            p.entity.dataBytes.runFold((0L, 0L))(receiveChunk).map {
              case (size, numChunks) =>
                info(s"Size is $size, numChunks is $numChunks")
                info(s"store file: $fName")
                distFile.close()
                // 加载文件成为表
                //                SparkCalcLib.loadCSVtoTable(spark, fName, p.filename.get)
                (fName, p.filename, size)
            }
          }.runFold(Seq.empty[(String, Option[String], Long)])(_ :+ _).map(_.mkString(", "))
          complete {
            s"""{ "code":"0", "msg":"success" }"""
          }
        }
      } catch {
        case ex: Exception =>
          val res = s"""{ "code":"-1", "msg": ${ex.getMessage} }"""
          complete {
            res
          }
      }
    } ~
      path("query_sql") {
        withoutRequestTimeout {
          withoutSizeLimit {
            entity(as[QueryInfo]) {
              queryInfo =>
                try {
                  val savePathTmp = utils.SparkCalcLib.RunSQLSaveLocalCSV(spark, queryInfo, ServerInfo.Respath)
                  val tran = getTranStream("/" + savePathTmp)
                  complete(
                    HttpEntity {
                      ContentTypes.`text/plain(UTF-8)`
                      tran.toByteArray
                    }
                  )
                } catch {
                  case ex: Exception =>
                    import models.QueryInfoModel.implicits._
                    QueryInfoModel.updateQueryInfo(queryInfo
                      .withStatus(13)
                      .withLastUpdate(DateTimeUtil.getTsNow())
                      .withExceptionMsg(ex.getClass + " " + ex.getMessage)
                    )
                    logger.error(ex.getClass + "\t" + ex.getMessage)
                    val tran = TranStream()

                    complete(HttpEntity {
                      ContentTypes.`application/octet-stream`
                      tran.toByteArray
                    })
                }
            }
          }
        }
      } ~
      path("query_custom") {
        withoutRequestTimeout {
          withoutSizeLimit {
            entity(as[QueryInfo]) {
              queryInfo =>
                try {
                  val savePathTmp = utils.SparkCalcLib.RunCustomTask(spark, queryInfo, ServerInfo.Respath)
                  val tran = getTranStream("/" + savePathTmp)
                  QueryInfoModel.updateQueryInfo(queryInfo.copy(status = QueryStatus.CustomQueryFinished, dataPath = savePathTmp, lastUpdate = DateTimeUtil.getTsNow))
                  complete(
                    HttpEntity {
                      ContentTypes.`text/plain(UTF-8)`
                      tran.toByteArray
                    }
                  )
                } catch {
                  case ex: Exception =>
                    QueryInfoModel.updateQueryInfo(queryInfo.copy(status = QueryStatus.CustomQueryFailed, exceptionMsg = ex.getClass + " " + ex.getMessage, lastUpdate = DateTimeUtil.getTsNow))
                    logger.error(ex.getClass + "\t" + ex.getMessage)
                    val tran = TranStream()

                    complete(HttpEntity {
                      ContentTypes.`application/octet-stream`
                      tran.toByteArray
                    })
                }
            }
          }
        }
      } ~
      path("shutdown") {
        ServerStatus.isRunning.set(true)
        complete("")
      } ~
      path("kill") {
        withoutRequestTimeout {
          withoutSizeLimit {
            entity(as[QueryInfo]) {
              queryInfo =>
                spark.sparkContext.cancelJobGroup(queryInfo.queryId)
                logger.info(s"kill JobGroup: queryId : ${queryInfo.queryId}")
                // update mysql db status
                import models.QueryInfoModel.implicits._
                QueryInfoModel.updateQueryInfo(queryInfo
                  .withStatus(13)
                  .withLastUpdate(DateTimeUtil.getTsNow())
                  .withExceptionMsg("job killed")
                )

                complete("")
            }
          }
        }
      } ~
      path("killByJobId") {
        withoutRequestTimeout {
          withoutSizeLimit {
            entity(as[String]) {
              queryId =>
                var entity: AStarRes = null
                try {
                  spark.sparkContext.cancelJobGroup(queryId)
                  entity = AStarRes(flag = true, s"已杀死")
                  val res = HttpEntity {
                    ContentTypes.`text/plain(UTF-8)`
                    Json.stringify(Json.toJson(entity)).getBytes
                  }
                  complete(res)

                } catch {
                  case e =>
                    e.printStackTrace()
                    entity = AStarRes(flag = false, e.getMessage)
                    val res = HttpEntity {
                      ContentTypes.`text/plain(UTF-8)`
                      Json.stringify(Json.toJson(entity)).getBytes
                    }
                    complete(res)
                }
            }
          }
        }
      } ~
      path("get_singleMeta") {
        get {
          parameters('dbName, 'tbName) { (dbName, tbName) =>
            info(s"${dbName}.${tbName}")
            val table = new HiveMetaDataHelper(spark).getTable(dbName, tbName)
            complete(
              HttpEntity {
                ContentTypes.`application/octet-stream`
                table.toByteArray
              }
            )
          }
        }
      } ~
      path("get_meta") {
        get {
          complete(
            HttpEntity {
              ContentTypes.`application/octet-stream`
              HiveMetaInfo.protoHiveMeta.toByteArray
            }
          )
        }
      } ~
      path("get_sparkUDF") {
        val udfs = Json.stringify(Json.toJson(HiveMetaInfo.sparkUDFs))
        complete(
          HttpEntity {
            ContentTypes.`text/plain(UTF-8)`
            udfs
          }
        )
      } ~
      path("get_limit") {
        withoutRequestTimeout {
          withoutSizeLimit {
            entity(as[QueryInfo]) {
              queryInfo =>
                val queryId = getQueryIdFromQueryPath(queryInfo.dataPath)
                var tran = new JedisHelper().getTrans(queryId)
                if (tran == null) {
                  val respath = ServerInfo.Respath
                  val originPath = s"${respath}/${queryInfo.dataPath}"

                  val hadoopConf = spark.sparkContext.hadoopConfiguration
                  val hdfs = FileSystem.get(hadoopConf)
                  val fsPath = new Path(originPath) // hdfs
                  val fis = hdfs.open(fsPath)

                  info("get_limit\t" + queryInfo.dataPath)
                  tran = XSSFUtil.printExcelByStream(fis, 5000)
                  info(s"数据获取完毕. ${queryInfo.dataPath}")

                  fis.close()
                }
                complete(
                  HttpEntity {
                    ContentTypes.`text/plain(UTF-8)`
                    tran.toByteArray
                  }
                )
            }
          }
        }
      } ~
      path("small_flow") {
        withoutRequestTimeout {
          withoutSizeLimit {
            entity(as[SmallFlow]) {
              smallFlow =>
                smallFlow.sqls.foreach(println)

                SmallFlowTaskLib.saveTask(spark, smallFlow)

                complete("")
            }
          }
        }
      } ~
      path("quickly_query") {
        withoutRequestTimeout {
          withoutSizeLimit {
            post {
              entity(as[Map[String, String]]) { map =>
                val sql = map("sql")
                val res = SparkCalcLib.quicklyQuery(spark, sql)
                complete(
                  HttpEntity {
                    ContentTypes.`text/plain(UTF-8)`
                    Json.stringify(Json.toJson(res)).getBytes
                  })
              }
            }
          }
        }
      } ~
      path("dag_getFieldName") {
        withoutRequestTimeout {
          withoutSizeLimit {
            post {
              entity(as[Map[String, String]]) { map =>
                val sql = map("sql")
                val taskType = map("taskType")
                info(s"$taskType:$sql")
                var res: ToResponseMarshallable = null
                try {
                  val df = spark.sql(sql)
                  val fields = df.schema.fields.map(x => x.name)
                  val entity = AStarRes(flag = true, Json.stringify(Json.toJson(fields)))
                  res = HttpEntity {
                    ContentTypes.`text/plain(UTF-8)`
                    Json.stringify(Json.toJson(entity)).getBytes
                  }
                } catch {
                  case e: Exception => {
                    e.printStackTrace()
                    val entity = AStarRes(flag = false, e.getMessage)
                    res = HttpEntity {
                      ContentTypes.`text/plain(UTF-8)`
                      Json.stringify(Json.toJson(entity)).getBytes
                    }
                  }
                }
                complete(
                  res
                )
              }
            }
          }
        }
      } ~
      path("dag_createTable") {
        withoutRequestTimeout {
          withoutSizeLimit {
            post {
              entity(as[List[DagTask]]) { tasks =>
                var tableSet = HiveTableSet()
                var res: ToResponseMarshallable = null
                try {
                  // 1. 创建输出表 (spark的, ck的, mysql的)
                  NodeParser.createTableFromTasks(tasks)
                  tableSet = SparkParser.getTableInfo(spark, tasks)

                  // 2. hive的返回表的元数据信息, ck和mysql的需要Astar自动导入
                  res = HttpEntity {
                    ContentTypes.`application/octet-stream`
                    tableSet.toByteArray
                  }
                } catch {
                  case e: Exception => {
                    e.printStackTrace()
                    val entity = AStarRes(flag = false, e.getMessage)
                    res = HttpEntity {
                      ContentTypes.`text/plain(UTF-8)`
                      Json.stringify(Json.toJson(entity)).getBytes
                    }
                  }
                }
                complete(
                  res
                )
              }
            }
          }
        }
      } ~
      path("dag_checkOutTable") {
        withoutRequestTimeout {
          withoutSizeLimit {
            post {
              entity(as[Map[String, String]]) { map =>
                val sql = map("sql")
                val outTable = map("outTable")
                val outType = map("outType")
                val (checkResult, msg) = SparkUtil.checkMeta(sql, outTable, outType)
                complete(
                  HttpEntity {
                    ContentTypes.`text/plain(UTF-8)`
                    Json.stringify(Json.toJson(AStarRes(checkResult, msg))).getBytes
                  })
              }
            }
          }
        }
      } ~
      path("dag_getMeta") { //  获取 spark sql 对应的meta信息
        withoutRequestTimeout {
          withoutSizeLimit {
            post {
              entity(as[Map[String, String]]) { map =>
                val sql = map("sql")
                val sqlTypes: Array[String] = SparkUtil.getMeta(sql)
                complete(
                  HttpEntity {
                    ContentTypes.`text/plain(UTF-8)`
                    Json.stringify(Json.toJson(sqlTypes)).getBytes
                  })
              }
            }
          }
        }
      } ~
      path("dag_getUpstream") {
        withoutRequestTimeout {
          withoutSizeLimit {
            post {
              entity(as[List[Map[String, String]]]) { maps =>
                var resp = Array[Map[String, String]]()
                for (map <- maps) {
                  val sql = map("sql")
                  val outTable = map("out")
                  val outType = map("outType")
                  val taskType = map("taskType")
                  val rolesStr = map("roles")
                  val upTables = SparkParser.getUpStreamTables(sql, rolesStr).map(
                    node => s"${node.db}.${node.tableName}"
                  )
                  resp = resp :+ Map(
                    "outTable" -> outTable,
                    "upstream" -> upTables.mkString(","),
                    "outType" -> outType,
                    "taskType" -> taskType
                  )
                }

                complete(
                  HttpEntity {
                    ContentTypes.`text/plain(UTF-8)`
                    Json.stringify(Json.toJson(resp)).getBytes
                  })
              }
            }
          }
        }
      } ~
      HiveSqlRouter.router
}

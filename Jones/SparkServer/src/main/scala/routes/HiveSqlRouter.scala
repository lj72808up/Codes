package routes


import akka.http.scaladsl.model.{ContentTypes, HttpEntity}
import akka.http.scaladsl.server.Directives.{as, complete, entity, path, post, withoutRequestTimeout, withoutSizeLimit}
import akka.http.scaladsl.server.Route
import akka.http.scaladsl.unmarshalling.{FromEntityUnmarshaller, PredefinedFromEntityUnmarshallers}
import play.api.libs.json.Json
import utils.udf.SqlValidator
import utils.{PrintLogger, SparkUtil}

object HiveSqlRouter extends PrintLogger {

  implicit val MapStringStringUnmarshaller: FromEntityUnmarshaller[Map[String, String]] =
    PredefinedFromEntityUnmarshallers.byteArrayUnmarshaller map {
      bytes => {
        val body = new String(bytes)
        Json.parse(body).as[Map[String, String]]
      }
    }

  val router: Route = path("sql_verify") { //  获取 spark sql 对应的meta信息
    withoutRequestTimeout {
      withoutSizeLimit {
        post {
          entity(as[Map[String, String]]) { map =>
            def fun(): String = {
              val sql = map.getOrElse("sql", "")
              val sqlType = map.getOrElse("sqlType", "")
              val dsId = map.getOrElse("dsId", "").toLong

              info(s"检验 => sqlType:$sqlType, dsId:$dsId, sql:$sql")
              SqlValidator.validateSql(sql, sqlType, dsId)
            }

            val response = makeResponse(fun)
            complete(
              response
            )
          }
        }
      }
    }
  }

  def makeResponse(fun: () => String): HttpEntity.Strict = {
    var astarRes: AStarRes = null
    try {
      val returnJson = fun()
      astarRes = AStarRes(
        flag = true,
        msg = returnJson
      )
    } catch {
      case e: Exception => {
        e.printStackTrace()
        astarRes = AStarRes(
          flag = false,
          msg = e.getMessage
        )
      }
    }

    HttpEntity {
      ContentTypes.`text/plain(UTF-8)`
      Json.stringify(Json.toJson(astarRes)).getBytes
    }
  }

}

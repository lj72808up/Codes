package utils

//import org.org.apache.commons.codec.binary.Base64
import org.joda.time.DateTime
import play.api.libs.json.Json
import scalaj.http.{Base64, Http}

case class KylinCubeUtil(baseUrl:String, user:String, passwd:String) {

  var encoding = GetEncoding(user, passwd)
  Login(encoding)

  def GetEncoding(user:String, passwd:String): String = {
    val keyBytes = (user + ":" + passwd)
    Base64.encodeString(keyBytes)
  }

  def Login(encoding:String) = {
    val para = "/user/authentication"
    Http(baseUrl + para).postData(encoding).asString
  }

  def GetCube(cubeName:String) = {
    val para = "/cubes/"
    Http(baseUrl + para + cubeName).header("Authorization","Basic " + GetEncoding(user, passwd)).asString
  }

  def BuildSegmentForCube(cubeName:String, startTime:Long, endTime:Long) = {
    val para = s"/cubes/$cubeName/rebuild"
    val data = Json.toJson(Map(
      "startTime" -> Json.toJson(startTime),
      "endTime" -> Json.toJson(endTime),
      "buildType" -> Json.toJson("BUILD")))
    Http(baseUrl + para).header("Authorization","Basic " + GetEncoding(user, passwd)).header("Content-Type","application/json;charset=UTF-8").put(data.toString()).asString
  }

  def GetCubeLastSegment(cubeName:String) = {
    val res = GetCube(cubeName).body
    val parseObj = Json.parse(res)
    val segmentsObj = Json.parse((parseObj \ "segments").get(0).toString())
    val dateRangeEnd = (segmentsObj \ "date_range_end").as[Long]
    dateRangeEnd
  }
}

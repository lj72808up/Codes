package routes

import play.api.libs.json.{Format, Json}

case class AStarRes (flag:Boolean, msg:String)
object AStarRes{
  implicit val jsonFormat: Format[AStarRes] = Json.format[AStarRes]
}

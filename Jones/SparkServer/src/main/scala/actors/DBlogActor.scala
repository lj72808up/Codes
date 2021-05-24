package actors

import actors.DBlogActor.{QueryInfoGet, QueryInfoUpdate, QueryInfoUpdateStart}
import adtl.platform.Dataproto.QueryInfo
import akka.actor.{Actor, ActorLogging, Props}


object DBlogActor {

  sealed trait DBlogAction

  case class QueryInfoGet(query: QueryInfo) extends DBlogAction

  case class QueryInfoUpdate(query: QueryInfo) extends DBlogAction

  case class QueryInfoUpdateStart(query: QueryInfo) extends DBlogAction

}

class DBlogActor extends Actor with ActorLogging{

  import RabbitMQActor._
  import models.QueryInfoModel.{getQueryInfo, updateQueryInfo, updateQueryInfoCreateTime, updateQueryInfoStartTime}

  override def receive: Receive = {
    case infoGet: QueryInfoGet => {
      try {
        getQueryInfo(infoGet.query)
      } catch {
        case ex: Exception => printMysqlEx(ex.getClass + " " + ex.getMessage)
      }
    }
    case infoUpdate: QueryInfoUpdate => {
      try {
        updateQueryInfo(infoUpdate.query)
      } catch {
        case ex: Exception => printMysqlEx(ex.getClass + "\t" + ex.getMessage)
      }
    }
    case infoUpdate: QueryInfoUpdateStart => {
      try {
        updateQueryInfoStartTime(infoUpdate.query)
      } catch {
        case ex: Exception => printMysqlEx(ex.getClass + "\t" + ex.getMessage)
      }
    }
  }

  def printMysqlEx(ex: String): Unit = {
    log.error("filed to update mysql" + ex)
  }
}


package models.config
import com.typesafe.config.Config

case class CommonConf(mq:String, host:String, userName:String, password:String)

object CommonConf{
  def apply(conf:Config): CommonConf = {
    val mq = conf.getString("mq")
    val host = conf.getString("host")
    val userName = conf.getString("userName")
    val password = conf.getString("password")
    new CommonConf(mq, host, userName, password)
  }
}

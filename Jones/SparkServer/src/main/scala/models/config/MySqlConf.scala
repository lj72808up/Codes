package models.config

import com.typesafe.config.Config

case class MySqlConf(host: String, port: String, user: String, passwd: String, db: String)

object MySqlConf {

  def apply(config: Config): MySqlConf = {
    MySqlConf(
      config.getString("host"),
      config.getString("port"),
      config.getString("user"),
      config.getString("passwd"),
      config.getString("db")
    )
  }
}


package models.config

import com.typesafe.config.Config

case class RedisConf(host: String, port: Int,  passwd: String)

object RedisConf {

  def apply(config: Config): RedisConf = {
    RedisConf(
      config.getString("host"),
      config.getInt("port"),
      config.getString("passwd")
    )
  }
}


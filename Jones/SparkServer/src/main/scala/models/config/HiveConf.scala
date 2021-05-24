package models.config

import com.typesafe.config.Config

case class HiveConf(jdbcHost: String)

object HiveConf {
  def apply(config: Config): HiveConf = {
    HiveConf(
      config.getString("jdbc")
    )
  }
}




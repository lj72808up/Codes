package models.config

import com.typesafe.config.Config

case class KylinConf(apiHost: String, jdbcHost: String, user: String, passwd: String, jonesCube: String)

object KylinConf {
  def apply(config: Config): KylinConf = {
    KylinConf(
      config.getString("api"),
      config.getString("jdbc"),
      config.getString("user"),
      config.getString("passwd"),
      config.getString("jones_cube")
    )
  }
}




package models.config

import com.typesafe.config.Config

case class JdbcConf(url: String, user: String, passwd: String,
                    initialSize: Int, minIdle: Int, maxActive: Int,
                    maxWait: Int,
                    testWhileIdle: Boolean = true, testOnBorrow: Boolean = false, testOnReturn: Boolean = false,
                    poolPreparedStatements: Boolean = true, maxPoolPreparedStatementPerConnectionSize: Int = 20
                   )

object JdbcConf {
  def apply(config: Config): JdbcConf = {
    new JdbcConf(
      config.getString("url"),
      config.getString("user"),
      config.getString("passwd"),
      config.getInt("initialSize"),
      config.getInt("minIdle"),
      config.getInt("maxActive"),
      config.getInt("maxWait")
    )
  }
}



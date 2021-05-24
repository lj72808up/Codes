package models.config

import com.typesafe.config.{Config, ConfigFactory}

case class MainConf(app: AppConf,
                    mysql: MySqlConf,
                    redis: RedisConf,
                    kylin: KylinConf,
                    hive: HiveConf,
                    jdbc: JdbcConf,
                    common: CommonConf
                   )

object MainConf {

  var appConf: AppConf = _
  var mySqlConf: MySqlConf = _
  var redisConf: RedisConf = _
  var kylinConf: KylinConf = _
  var hiveConf: HiveConf = _
  var jdbcConf: JdbcConf = _
  var worksheetConf: JdbcConf = _
  var commonConf: CommonConf = _

  var confMap: Map[String, MainConf] = Map[String, MainConf]()

  def apply(conf: Config): MainConf = {
    new MainConf(
      AppConf(conf.getConfig("app")),
      MySqlConf(conf.getConfig("mysql")),
      RedisConf(conf.getConfig("redis")),
      KylinConf(conf.getConfig("kylin")),
      HiveConf(conf.getConfig("hive")),
      JdbcConf(conf.getConfig("jdbc")),
      CommonConf(conf.getConfig("common"))
    )
  }

  def apply(name: String): MainConf = {
    if (!confMap.keySet.contains("basicConf")) {
      val conf = ConfigFactory.load()
      addConf("basicConf", MainConf(conf))
    }
    confMap(name)
  }

  def addConf(name: String, conf: MainConf) {
    confMap = confMap + (name -> conf)
  }
}

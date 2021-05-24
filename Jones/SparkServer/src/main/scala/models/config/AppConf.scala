package models.config

import com.typesafe.config.Config

case class AppConf(sparkName:String, host:String, port:String, resPath:String, fuctionJars:String)

object AppConf{

  def apply(config:Config): AppConf ={
    AppConf(
      config.getString("spark_name"),
      config.getString("host"),
      config.getString("port"),
      config.getString("res_path"),
      config.getString("functionJars")
    )
  }
}

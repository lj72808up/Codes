package utils

import com.typesafe.config.{Config, ConfigFactory}

import scala.io.Source

object ConfigUtil {
  def confFromFile(path:String) = {
    val content = Source.fromFile(path).getLines().mkString("\n")
    ConfigFactory.parseString(content)
  }

}
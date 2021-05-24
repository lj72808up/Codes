package utils

import java.time.{LocalDate, LocalTime}

trait PrintLogger {
  def info(info :String):Unit = {
    println(s"[${LocalDate.now} ${LocalTime.now}] ${this.getClass.getName} $info")
  }
}
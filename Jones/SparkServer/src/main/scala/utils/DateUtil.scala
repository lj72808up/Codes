package utils

import org.joda.time.DateTime
import org.joda.time.format.DateTimeFormat

object DateUtil {
  def UnixToDateTime(time:Long) = {
    new DateTime(time)
  }

  def UnixNextDate(time:Long) = {
    ((new DateTime(time)).plusDays(1)).getMillis
  }

   def getYesterDayDate(): String = {
    DateTime.now().minusDays(1).toString("yyyyMMdd")

  }
}

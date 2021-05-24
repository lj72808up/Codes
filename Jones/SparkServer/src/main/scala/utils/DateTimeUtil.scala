package utils

import org.joda.time.DateTime

object DateTimeUtil {

  def getTsNow():Long = {
    DateTime.now().getMillis
  }

}

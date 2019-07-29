package com.test.loggers

import com.test.log.Logger

class PureLogger extends Logger {

  def pureinfo(info:String):Unit = {
    log.info(info)
  }
}

package com.test.log

import org.slf4j.LoggerFactory

trait Logger {
  protected val log = LoggerFactory.getLogger(this.getClass.getName)
}
package com.test.api

class GreetingsServiceImpl extends GreetingsService {
  override def sayHi(name: String): String = {
    "hi, " + name
  }

}

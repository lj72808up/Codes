package com.test

class AA{
  def say():Unit = {
    println(AA.aa)
  }
}

object AA{
  val aa:String = build()
  def build():String = {println("init"); return "haha"}

}

object TestAmqp2 {
  def main(args: Array[String]): Unit = {
    new AA().say()
    new AA().say()
  }
}


package com.test

import java.io.{File, FileOutputStream}

import com.sogou.adtl.ahocorasick.trie.{PayloadTrie, TrieConfig}
import utils.PrintLogger

object TestLogger extends PrintLogger {
  def main1(args: Array[String]): Unit = {

    val trie = new PayloadTrie[String](new TrieConfig)
    trie.addKeyword("aab", "aaa1")
    trie.addKeyword("bbb", "bbb1")
    trie.addKeyword("ccc", "ccc1")

    trie.constructFailureStates()

    println(trie.parseText("aabbbccc"))
  }

  def main(args: Array[String]): Unit = {
    val map1 = scala.collection.mutable.Map[String,String]()
    try{
      map1("name") = "zhangsan"
      throw new Exception("测试异常")
      map1("age") = "23"
    }catch {
      case e:Exception => e.printStackTrace()
    }
    println(map1)
  }
}

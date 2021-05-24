package util

import com.sogou.adtl.ahocorasick.trie.{PayloadTrie, TrieConfig}
import org.apache.spark.SparkFiles

import scala.io.Source

class MTrie(val fileName:String) {
  val trie = new PayloadTrie[String](new TrieConfig)
  Source.fromFile(SparkFiles.get(fileName)).getLines().foreach{
    line=>
      trie.addKeyword(line, line)
  }
  trie.constructFailureStates()

  println("Init:" + fileName)

  def containsMatch(word:String) = {
    if(trie.parseText(word).size() > 0) 1 else 0
  }

  def equalMatch(word:String) = {
    import scala.collection.JavaConverters._
    if(trie.parseText(word).asScala.count(_.getKeyword == word) > 0) 1 else 0

  }
}

object MTrieFun {
  def getFun(file:String):String = {
    val packageName = this.getClass.getPackage.getName
    s"""import $packageName._
       |lazy val trie = new MTrie("$file")
       |def apply(word:String, matchType:Int):Int={
       |  if(matchType == 0) {
       |    trie.equalMatch(word)
       |  } else {
       |    trie.containsMatch(word)
       |  }
       |}""".stripMargin
  }
}
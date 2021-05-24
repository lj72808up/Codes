import org.apache.spark.sql.SparkSession
import util.{DynamicUDF, MTrieFun}

object TestStart {
  def main(args: Array[String]): Unit = {
    val qwPartition = "testWord"
    val spark = SparkSession.builder().master("local").getOrCreate()
    import spark.implicits._
    val df = Seq[String]("百度aa", "葡萄", "搜狗狗").toDF("word")
    println(df.schema)
    df.createOrReplaceTempView("t1")

    DynamicUDF.register(spark, MTrieFun.getFun(qwPartition), qwPartition)
    println("注册完毕")

//    spark.sql("select word from t1").collect().map(row => row.getString(0)).foreach(println)

    spark.sql("select word from t1 where testWord(t1.word, 0) = 1")
      .collect().map(row => row.getString(0)).foreach(println)


  }
}

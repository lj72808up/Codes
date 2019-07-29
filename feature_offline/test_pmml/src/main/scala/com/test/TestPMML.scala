package com.test

import java.io.File
import org.apache.spark.sql.functions._
import org.apache.spark.sql.SparkSession
import org.jpmml.evaluator.spark.TransformerBuilder
import org.jpmml.evaluator.{DefaultVisitorBattery, LoadingModelEvaluatorBuilder}
//https://github.com/jpmml/jpmml-evaluator-spark
object TestPMML {
  def main(args: Array[String]): Unit = {

    val spark = SparkSession
      .builder()
      .master("local")
      .appName("testpmml")
      .getOrCreate()

    import spark.implicits._
    val in = Seq((1,2,3,1),(2,4,1,5),(7,8,3,6),(4,8,4,7),(2,5,6,9)).toDF("x1","x2","x3","x4")

    val f = "/Users/liujie32/IdeaProjects/baidu/xdata/trustfinger/demo/FeatureEngineering/feature_fingerPrint_offline/test_pmml/src/main/python/demo.pmml"
    val evaluatorBuilder = new LoadingModelEvaluatorBuilder()
      .setLocatable(false)
      .setVisitors(new DefaultVisitorBattery())
      .load(new File(f));

    val evaluator = evaluatorBuilder.build()

    // Performing a self-check (duplicates as a warm-up)
    evaluator.verify();

    val pmmlTransformerBuilder = new TransformerBuilder(evaluator)
      .withTargetCols()
      .withOutputCols()
      .exploded(false);

    val pmmlTransformer = pmmlTransformerBuilder.build()
    val os = pmmlTransformer.transform(in)
    println(os.columns.mkString(","))
    os.select(col("pmml")).collect().foreach(println)
  }


  def fun1() = {

  }
}

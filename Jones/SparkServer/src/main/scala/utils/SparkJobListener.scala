package utils

import org.apache.spark.scheduler.{SparkListener, SparkListenerJobEnd, SparkListenerJobStart}

class SparkJobListener extends SparkListener {
  override def onJobStart(jobStart: SparkListenerJobStart): Unit = {
    println(jobStart.jobId, jobStart.properties.getProperty("spark.jobGroup.id"))
  }

  override def onJobEnd(jobEnd: SparkListenerJobEnd): Unit = {
    println(jobEnd.jobResult, jobEnd.time)
  }
}

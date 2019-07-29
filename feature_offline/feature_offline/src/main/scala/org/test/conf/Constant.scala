package org.test.conf

import com.typesafe.config.ConfigFactory

object Constant {

  val config = ConfigFactory.load()

  lazy val srcTable = config.getString("table.sourceTable")
  lazy val distTable = config.getString("table.distTable")
  lazy val levelTable = config.getString("table.levelTable")
  lazy val metaTable = config.getString("table.metaTable")

  lazy val ccDistTable = config.getString("cc.distTable")
  lazy val ccSrcTable = config.getString("cc.sourceTable")
  lazy val threshold = config.getString("cc.threshold")
  lazy val iterationLimit = config.getInt("cc.iterationLimit")

  lazy val createSql =
    s"""create table if not exists ${distTable}(featureindex string,t1_primarykey string, t1_clusterid string,t1_magic string, t1_zid string,t1_burry string,
      |features string,
      |t2_primarykey string, t2_clusterid string, t2_magic string, t2_zid string,t2_burry string, status string)
      |partitioned by (save_time string)
      |STORED AS ORC
      |tblproperties ("orc.compress"="ZLIB")""".stripMargin

  lazy val createCCSql =
    s"""
       |create table if not exists ${ccDistTable}(
       |    sha                     string,
       |    finger_new              string
       |)
       |PARTITIONED BY (stat_date string)
       |STORED AS ORC
       |tblproperties ("orc.compress"="ZLIB");
     """.stripMargin
}


case class ExtendStructure(t1_primarykey:Option[String],
                           t1_clusterid: Option[String],
                           t1_status:String,
                           t1_magic:Option[String],
                           t1_zid:Option[String],
                           t1_burry:Option[String],
                           imei: Option[String],
                           mac: Option[String],
                           imsi: Option[String],
                           iccid: Option[String],
                           android: Option[String],
                           sdcid: Option[String],
                           serial: Option[String],
                           manufacturer: Option[String],
                           model: Option[String],
                           brand: Option[String],
                           mainboard: Option[String],
                           equipment: Option[String],
                           hardware: Option[String],
                           product: Option[String],
                           android_ver: Option[String],
                           api_level: Option[String],
                           version: Option[String],
                           picture_external_sd: Option[String],
                           picture_internal_sd: Option[String],
                           music_internal_sd: Option[String],
                           video_internal_sd: Option[String],
                           music_external_sd: Option[String],
                           video_external_sd: Option[String],
                           resolution: Option[String],
                           size: Option[String],
                           cpu_core: Option[String],
                           cpu_kernel: Option[String],
                           cpu_core_num: Option[String],
                           sys_fingerprint: Option[String],
                           cpu_fiber: Option[String],
                           decode_encode: Option[String],
                           t2_primarykey:Option[String],
                           t2_clusterid: Option[String],
                           t2_status:String,
                           t2_magic:Option[String],
                           t2_zid:Option[String],
                           t2_burry:Option[String])


case class Feature(imei: Option[String],
                   mac: Option[String],
                   imsi: Option[String],
                   iccid: Option[String],
                   android: Option[String],
                   sdcid: Option[String],
                   serial: Option[String],
                   manufacturer: Option[String],
                   model: Option[String],
                   brand: Option[String],
                   mainboard: Option[String],
                   equipment: Option[String],
                   hardware: Option[String],
                   product: Option[String],
                   android_ver: Option[String],
                   api_level: Option[String],
                   version: Option[String],
                   picture_external_sd: Option[String],
                   picture_internal_sd: Option[String],
                   music_internal_sd: Option[String],
                   video_internal_sd: Option[String],
                   music_external_sd: Option[String],
                   video_external_sd: Option[String],
                   resolution: Option[String],
                   size: Option[String],
                   cpu_core: Option[String],
                   cpu_kernel: Option[String],
                   cpu_core_num: Option[String],
                   sys_fingerprint: Option[String],
                   cpu_fiber: Option[String],
                   decode_encode: Option[String])


case class FeatureResult(featureKey:String,
                         t1_primarykey:Option[String],
                         t1_clusterid: Option[String],
                         t1_magic:Option[String],
                         t1_zid:Option[String],
                         t1_burry:Option[String],
                         features:String,
                         t2_primarykey:Option[String],
                         t2_clusterid: Option[String],
                         t2_magic:Option[String],
                         t2_zid:Option[String],
                         t2_burry:Option[String],
                         status:String
                        )

object FeatureEncoders{
  implicit val extendStructureEncoder = org.apache.spark.sql.Encoders.product[ExtendStructure]
  implicit val featureEncoder = org.apache.spark.sql.Encoders.product[Feature]
  implicit val featureResultEncoder = org.apache.spark.sql.Encoders.product[FeatureResult]

}
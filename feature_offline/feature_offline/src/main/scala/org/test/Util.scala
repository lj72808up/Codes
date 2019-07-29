package org.test

import java.security.MessageDigest
import java.text.SimpleDateFormat
import java.util.{Calendar, Random, TimeZone}

import com.fasterxml.jackson.databind.{DeserializationFeature, ObjectMapper}
import com.fasterxml.jackson.module.scala.DefaultScalaModule
import com.fasterxml.jackson.module.scala.experimental.ScalaObjectMapper
import org.test.conf.{ExtendStructure, Feature, FeatureResult}

object ConvertUtil {
  def convertFeature(obj: ExtendStructure): FeatureResult = {
    val features = Feature(obj.imei,
      obj.mac,
      obj.imsi,
      obj.iccid,
      obj.android,
      obj.sdcid,
      obj.serial,
      obj.manufacturer,
      obj.model,
      obj.brand,
      obj.mainboard,
      obj.equipment,
      obj.hardware,
      obj.product,
      obj.android_ver,
      obj.api_level,
      obj.version,
      obj.picture_external_sd,
      obj.picture_internal_sd,
      obj.music_internal_sd,
      obj.video_internal_sd,
      obj.music_external_sd,
      obj.video_external_sd,
      obj.resolution,
      obj.size,
      obj.cpu_core,
      obj.cpu_kernel,
      obj.cpu_core_num,
      obj.sys_fingerprint,
      obj.cpu_fiber,
      obj.decode_encode)
    val featureStr = JsonUtil.toJson(features)
    val status = if ("add".equals(obj.t1_status) || "add".equals(obj.t2_status)){
      "add"
    }else if("update".equals(obj.t1_status) || "update".equals(obj.t2_status)){
      "update"
    }else{
      "unchange"
    }
    /*val d1 = Map("t1_sdcid1"->obj.t1_sdcid1.getOrElse(""),
      "t1_sdcid2"->obj.t1_sdcid2.getOrElse(""),
      "t1_magic"->obj.t1_magic.getOrElse(""),
      "t1_burry"->obj.t1_burry.getOrElse(""),
      "t1_imei"->obj.t1_imei.getOrElse(""),
      "t1_imsi"->obj.t1_imsi.getOrElse(""),
      "t1_iccid"->obj.t1_iccid.getOrElse(""),
      "t1_android"->obj.t1_android.getOrElse(""),
      "t1_mac"->obj.t1_mac.getOrElse(""),
      "t1_status"->obj.t1_status
    )

    val d2 = Map(
      "t2_sdcid1"->obj.t2_sdcid1.getOrElse(""),
      "t2_sdcid2"->obj.t2_sdcid2.getOrElse(""),
      "t2_magic"->obj.t2_magic.getOrElse(""),
      "t2_burry"->obj.t2_burry.getOrElse(""),
      "t2_imei"->obj.t2_imei.getOrElse(""),
      "t2_imsi"->obj.t2_imsi.getOrElse(""),
      "t2_iccid"->obj.t2_iccid.getOrElse(""),
      "t2_android"->obj.t2_android.getOrElse(""),
      "t2_mac"->obj.t2_mac.getOrElse(""),
      "t2_status"->obj.t2_status
    )*/
    val featureIndexStr = s"${featureStr}_${obj.t1_primarykey}_${obj.t2_primarykey}"
    val featureIndex = MessageDigest.getInstance("MD5").digest(featureIndexStr.getBytes).map("%02x".format(_)).mkString

    return FeatureResult(featureIndex,
      obj.t1_primarykey,
      obj.t1_clusterid,
      obj.t1_magic,
      obj.t1_zid,
      obj.t1_burry,
      featureStr,
      obj.t2_primarykey,
      obj.t2_clusterid,
      obj.t2_magic,
      obj.t2_zid,
      obj.t2_burry,
      status
    )
  }


}

object JsonUtil {
  val mapper = new ObjectMapper() with ScalaObjectMapper
  mapper.setTimeZone(TimeZone.getTimeZone("Asia/Shanghai"))
  mapper.registerModule(DefaultScalaModule)
  mapper.configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false)
  //  mapper.registerModule(new JodaModule)

  def toJson(value: Map[Symbol, Any]): String = {
    toJson(value map { case (k, v) => k.name -> v })
  }

  def toJson(value: Any): String = {
    mapper.writeValueAsString(value)
  }

  def toMap[V](json: String)(implicit m: Manifest[V]) = fromJson[Map[String, V]](json)

  def fromJson[T](json: String)(implicit m: Manifest[T]): T = {
    mapper.readValue[T](json)
  }

  def main(args: Array[String]): Unit = {
    val str1 = """[[{u'OMX.dolby.ac4.decoder':u'[audio/ac4]'},{u'OMX.dolby.eac3.decoder':u'[audio/eac3]'},{u'OMX.dolby.ac3.decoder':u'[audio/ac3]'},{u'OMX.hisi.video.encoder.avc':u'[video/avc]'},{u'OMX.hisi.video.decoder.hevc':u'[video/hevc]'},{u'OMX.hisi.video.decoder.vp9':u'[video/x-vnd.on2.vp9]'},{u'OMX.hisi.video.decoder.avc':u'[video/avc]'},{u'OMX.hisi.video.decoder.vp8':u'[video/x-vnd.on2.vp8]'},{u'OMX.hisi.video.encoder.hevc':u'[video/hevc]'},{u'OMX.hisi.video.decoder.mpeg2':u'[video/mpeg2]'},{u'OMX.dolby.eac3_joc.decoder':u'[audio/eac3-joc]'},{u'OMX.hisi.video.decoder.mpeg4':u'[video/mp4v-es]'}]]"""
    val str2 = """[[{'OMX.dolby.ac4.decoder':'[audio/ac4]'},{'OMX.dolby.eac3.decoder':'[audio/eac3]'},{'OMX.dolby.ac3.decoder':'[audio/ac3]'},{'OMX.hisi.video.encoder.avc':'[video/avc]'}, {'OMX.hisi.video.decoder.hevc':'[video/hevc]'},{'OMX.hisi.video.decoder.vp9':'[video/x-vnd.on2.vp9]'},{'OMX.hisi.video.decoder.avc':'[video/avc]'},   {'OMX.hisi.video.decoder.vp8':'[video/x-vnd.on2.vp8]'}, {'OMX.hisi.video.encoder.hevc':'[video/hevc]'},{'OMX.hisi.video.decoder.mpeg2':'[video/mpeg2]'}, {'OMX.dolby.eac3_joc.decoder':'[audio/eac3-joc]'}, {'OMX.hisi.video.decoder.mpeg4':'[video/mp4v-es]'}]]"""
    //println (JsonUtil.toMap[Any](str1.replace("u'","'")))
    println(JsonUtil.fromJson[Any](str1.replace("u'","'").replace("'","\"")))
  }
}

object ParseUtil{

  def parseSaveTime(args:Array[String]):String={
    var saveTime = ""
    if (args.length==0){
      val calendar = Calendar.getInstance()
      calendar.add(Calendar.DATE,-1)
      saveTime = new SimpleDateFormat("yyyyMMdd").format(calendar.getTime)
    }else{
      saveTime = args(0)
    }
    return saveTime
  }
}
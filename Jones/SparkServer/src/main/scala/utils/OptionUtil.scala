package utils

import adtl.platform.Dataproto._

import org.apache.spark.sql.SparkSession

case class OptionUtil(spark: SparkSession) {


  def getOptions: FrontOptions = {

    val dt = utils.DateUtil.getYesterDayDate()

    // CXCHY
    println("CXCHY")
    val cxchySql = s"SELECT DISTINCT dw_dim_query_indus.indus_name FROM ADTL_BIZ.dw_dim_query_indus"
    val cxchy: Array[String] = spark.sql(cxchySql).collect().map(_.getAs[String](0))
    val cxchyOption = FrontOption().withOptions(cxchy).withName("cxchyOptions")
    cxchy.foreach(print)

    // WMQY
    println("WMQY")
    val wmqySql = s"SELECT DISTINCT dw_dim_city.city_name FROM ADTL_BIZ.dw_dim_city"
    val wmqy: Array[String] = spark.sql(wmqySql).collect().map(_.getAs[String](0))
    val wmqyOption = FrontOption().withOptions(wmqy).withName("wmqySqlOptions")
    wmqy.foreach(print)

    // YJKHHY
    println("YJKHHY")
    val yjkhhySql = s"SELECT DISTINCT dw_dim_account.first_indus_name FROM ADTL_BIZ.dw_dim_account"
    val yjkhhy: Array[String] = spark.sql(yjkhhySql).collect().map(_.getAs[String](0))
    val yjkhhyOption = FrontOption().withOptions(yjkhhy).withName("yjkhhyOptions")
    yjkhhy.foreach(print)

    // EJKHHY
    println("EJKHHY")
    val ejkhlxSql = s"SELECT DISTINCT dw_dim_account.second_indus_name FROM ADTL_BIZ.dw_dim_account"
    val ejkhlx: Array[String] = spark.sql(ejkhlxSql).collect().map(_.getAs[String](0))
    val ejkhlxOption = FrontOption().withOptions(ejkhlx).withName("ejkhlxOptions")
    ejkhlx.foreach(print)

    // EJYWLX
    println("EJYWLX")
    val ejywlxSql = s"SELECT DISTINCT dw_dim_service.second_name FROM ADTL_BIZ.dw_dim_service"
    val ejywlx: Array[String] = spark.sql(ejywlxSql).collect().map(_.getAs[String](0))
    val ejywlxOption = FrontOption().withOptions(ejywlx).withName("ejywlxOptions")
    ejywlx.foreach(print)

    // SJYWLX
    println("SJYWLX")
    val sjywlxSql = s"SELECT DISTINCT dw_dim_service.third_name FROM ADTL_BIZ.dw_dim_service"
    val sjywlx: Array[String] = spark.sql(sjywlxSql).collect().map(_.getAs[String](0))
    val sjywlxOption = FrontOption().withOptions(sjywlx).withName("sjywlxOptions")
    sjywlx.foreach(print)

    // KHLX
    println("KHLX")
    val khlxSql = s"SELECT DISTINCT dw_dim_account.perftype_id FROM ADTL_BIZ.dw_dim_account"
    val khlx: Array[String] = spark.sql(khlxSql).collect().map(_.getAs[String](0))
    val khlxOption = FrontOption().withOptions(khlx).withName("khlxOptions")
    khlx.foreach(print)

    // DQ
    println("DQ")
    val dqSql = s"SELECT DISTINCT dw_dim_account.biztype_id_ext FROM ADTL_BIZ.dw_dim_account"
    val dq: Array[String] = spark.sql(dqSql).collect().map(_.getAs[String](0))
    val dqOption = FrontOption().withOptions(dq).withName("dqOptions")
    dq.foreach(print)

    // YSLB
    println("YSLB")
    val yslbSql = s"SELECT DISTINCT YSLB(dw_ws_ad_pv_click.ext_reserve, dw_ws_ad_pv_click.style_reserve) FROM ADTL_BIZ.dw_ws_ad_pv_click where dt = $dt"
    val yslb: Array[String] = spark.sql(yslbSql).collect().map(_.getAs[String](0))
    val yslbOption = FrontOption().withOptions(yslb).withName("yslbOptions")
    yslb.foreach(print)

    // YJQD
    println("YJQD")
    val yjqdSql = s"SELECT DISTINCT dw_dim_pid.first_channel FROM ADTL_BIZ.dw_dim_pid"
    val yjqd: Array[String] = spark.sql(yjqdSql).collect().map(_.getAs[String](0))
    val yjqdOption = FrontOption().withOptions(yjqd).withName("yjqdOptions")
    yjqd.foreach(print)

    // EJQD
    println("EJQD")
    val ejqdSql = s"SELECT DISTINCT dw_dim_pid.second_channel FROM ADTL_BIZ.dw_dim_pid"
    val ejqd: Array[String] = spark.sql(ejqdSql).collect().map(_.getAs[String](0))
    val ejqdOption = FrontOption().withOptions(ejqd).withName("ejqdOptions")
    ejqd.foreach(print)

    // WZ
    println("WZ")
    val wzSql = s"SELECT DISTINCT WZ(dw_ws_ad_pv_click.reserved) FROM ADTL_BIZ.dw_ws_ad_pv_click where dt = $dt"
    val wz: Array[String] = spark.sql(wzSql).collect().map(_.getAs[String](0))
    val wzOption = FrontOption().withOptions(wz).withName("wzOptions")
    wz.foreach(print)

    // SD
    println("SD")
    val sdSql = s"SELECT DISTINCT dw_fact_ws_pagepv.ad_num FROM ADTL_BIZ.dw_fact_ws_pagepv where dt = $dt"
    val sd: Array[String] = spark.sql(sdSql).collect().map(_.getAs[String](0))
    val sdOption = FrontOption().withOptions(sd).withName("sdOptions")
    sd.foreach(print)

    // PW
    println("PW")
    val pwSql = s"SELECT DISTINCT PW(dw_ws_ad_pv_click.reserved) FROM ADTL_BIZ.dw_ws_ad_pv_click where dt = $dt"
    val pw: Array[String] = spark.sql(pwSql).collect().map(_.getAs[String](0))
    val pwOption = FrontOption().withOptions(pw).withName("pwOptions")
    pw.foreach(print)

    // CFLX
    println("CFLX")
    val cflxSql = s"SELECT DISTINCT CFLX(dw_ws_ad_pv_click.retrieve_flag) FROM ADTL_BIZ.dw_ws_ad_pv_click where dt = $dt"
    val cflx: Array[String] = spark.sql(cflxSql).collect().map(_.getAs[String](0))
    val cflxOption = FrontOption().withOptions(cflx).withName("cflxOptions")
    cflx.foreach(print)

    //LLPT
    println("LLPT")
    val llptSql = s"SELECT DISTINCT LLPT(dw_ws_ad_pv_click.flow_platform) FROM ADTL_BIZ.dw_ws_ad_pv_click where dt = $dt"
    val llpt: Array[String] = spark.sql(llptSql).collect().map(_.getAs[String](0))
    val llptOption = FrontOption().withOptions(llpt).withName("llptOptions")
    llpt.foreach(print)

    val ans: Array[FrontOption] = Array[FrontOption](
      cxchyOption,
      wmqyOption,
      yjkhhyOption,
      ejkhlxOption,
      ejywlxOption,
      sjywlxOption,
      khlxOption,
      dqOption,
      yslbOption,
      yjqdOption,
      ejqdOption,
      wzOption,
      sdOption,
      pwOption,
      cflxOption,
      llptOption
    )

    FrontOptions().withFrontOpts(ans)
  }
  
}

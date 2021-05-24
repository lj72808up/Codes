package com.test

import java.sql.Driver
import java.util.Properties
;

object TestKylin {
  def main(args: Array[String]): Unit = {
    val driver = Class.forName("org.apache.kylin.jdbc.Driver").newInstance().asInstanceOf[Driver]
    val info = new Properties()
    info.put("user", "ADMIN")
    info.put("password", "KYLIN")
    // http://10.140.85.161:7070/kylin/api/cubes/cube_wsbidding_smallflow_for_platform
    val projectName = "adtl_test"
    val conn = driver.connect("jdbc:kylin://10.140.85.161:7070/" + projectName, info)
    //    val sql = "select sum(pv2) from adtl_biz.VIEW_DWS_WSBIDDING_JONESPAGE_DI where DT='20210411'"
    val sql = "select sum(price), count(*) from adtl_biz.view_dws_wsbidding_jonespage_di where dt = 20210425"
    val state = conn.prepareStatement(sql)
    //    state.setInt(1, 10)noob
    val resultSet = state.executeQuery

    while (resultSet.next) {
      println(resultSet.getString(1))
    }
  }
}

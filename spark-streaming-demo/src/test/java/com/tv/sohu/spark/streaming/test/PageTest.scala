package com.tv.sohu.spark.streaming.test

import com.tv.sohu.spark.streaming.dm3.page.{AppLogParse, BCLogParse}

object PageTest {
  def main(args: Array[String]): Unit = {
    // 曝光日志
    val bcLog = "220.250.51.111 - z.m.tv.sohu.com [25/Mar/2021:13:59:57  0800] \"jsonlog={\"atype\":\"apps\",\"channelid\":\"6892\",\"cv\":\"8.2.3\",\"enterid\":\"18_1000010000_1004035583\",\"imei\":\"\",\"mfo\":\"TCL\",\"mfov\":\"TCL P308L\",\"mos\":\"android\",\"mosv\":\"4.4.4\",\"msg\":\"clk_vv\",\"mtype\":\"6\",\"passport\":\"\",\"pro\":\"1\",\"sim\":\"1\",\"startid\":\"1616651858719\",\"tkey\":\"410c9f41144f11f852e936e7c19815598ab05387\",\"uid\":\"8c99784ce38786b1b429321345859a72\",\"vids\":[{\"catecode\":\"101\",\"ctype\":\"video\",\"datatype\":1,\"idx\":\"0004\",\"mdu\":\"3809\",\"memo\":\"{\\\"mdx\\\":\\\"0006\\\",\\\"abmod\\\":\\\"35B_43A_71A_79B_75A_87B_9998A\\\"}\",\"pg\":\"10000\",\"playlistid\":9601609,\"scn\":\"01\",\"site\":1,\"time\":1616651997206,\"vid\":5772733}],\"webtype\":\"WiFi\"}\" 204 0 \"com.sohu.sohuvideo/8002003 (Linux; U; Android 4.4.4; zh_CN; TCL P308L; Build/KTU84P)\"";
    println(BCLogParse.parseLog(bcLog));
    // 曝光日志 - 访客
    val bcLog2 = "220.250.51.111 - z.m.tv.sohu.com [25/Mar/2021:13:59:57  0800] \"jsonlog={\"atype\":\"apps\",\"channelid\":\"6892\",\"cv\":\"8.2.3\",\"enterid\":\"18_1000010000_1004035583\",\"imei\":\"\",\"mfo\":\"TCL\",\"mfov\":\"TCL P308L\",\"mos\":\"android\",\"mosv\":\"4.4.4\",\"msg\":\"clk_vv\",\"mtype\":\"6\",\"passport\":\"\",\"pro\":\"1\",\"sim\":\"1\",\"startid\":\"1616651858719\",\"tkey\":\"410c9f41144f11f852e936e7c19815598ab05387\",\"uid\":\"guest8c99784ce38786b1b429321345859a72\",\"vids\":[{\"catecode\":\"101\",\"ctype\":\"video\",\"datatype\":1,\"idx\":\"0004\",\"mdu\":\"3809\",\"memo\":\"{\\\"mdx\\\":\\\"0006\\\",\\\"abmod\\\":\\\"35B_43A_71A_79B_75A_87B_9998A\\\"}\",\"pg\":\"10000\",\"playlistid\":9601609,\"scn\":\"01\",\"site\":1,\"time\":1616651997206,\"vid\":5772733}],\"webtype\":\"WiFi\"}\" 204 0 \"com.sohu.sohuvideo/8002003 (Linux; U; Android 4.4.4; zh_CN; TCL P308L; Build/KTU84P)\"";
    println(BCLogParse.parseLog(bcLog2));

    //行为日志
    val actionLog = "119.186.249.16 - mb.hd.sohu.com.cn [25/Mar/2021:13:59:55 +0800] \"GET /mc.gif?uid=5c043e60a85c4f6ad41ea5c6583e059e&tkey=2d6a0b099b14dc1218277acea07fc55a46326887&vid=&tvid=&mtype=6&mfov=PACT00&memo=%7B%22channeled%22%3A%221000029150%22%7D&passport=1248601464356564992%40sohu.com&ltype=&cv=8.8.0&mos=2&mosv=10&pro=1&mfo=OPPO&webtype=WiFi&time=1616651994999&channelid=281&sim=1&playlistid=&catecode=&preid=&newuser=0&enterid=0&startid=1616651978623&abmod=35B_43B_71B_79B_75B_87A_9998C&guid=&build=03221605&oaid=&url=5028&type=1&value=&androidId=0942553fe6b972e1&imei= HTTP/1.1\" 204 0 \"com.sohu.sohuvideo/8008000 (Linux; U; Android 10; zh_CN; PACT00; Build/QP1A.190711.020)\""
    println(AppLogParse.parseActionLog(actionLog))

    // 行为日志 - 访客
    val actionLog2 = "119.186.249.16 - mb.hd.sohu.com.cn [25/Mar/2021:13:59:55 +0800] \"GET /mc.gif?uid=guest5c043e60a85c4f6ad41ea5c6583e059e&tkey=2d6a0b099b14dc1218277acea07fc55a46326887&vid=&tvid=&mtype=6&mfov=PACT00&memo=%7B%22channeled%22%3A%221000029150%22%7D&passport=1248601464356564992%40sohu.com&ltype=&cv=8.8.0&mos=2&mosv=10&pro=1&mfo=OPPO&webtype=WiFi&time=1616651994999&channelid=281&sim=1&playlistid=&catecode=&preid=&newuser=0&enterid=0&startid=1616651978623&abmod=35B_43B_71B_79B_75B_87A_9998C&guid=&build=03221605&oaid=&url=5028&type=1&value=&androidId=0942553fe6b972e1&imei= HTTP/1.1\" 204 0 \"com.sohu.sohuvideo/8008000 (Linux; U; Android 10; zh_CN; PACT00; Build/QP1A.190711.020)\""
    println(AppLogParse.parseActionLog(actionLog2))

  }
}

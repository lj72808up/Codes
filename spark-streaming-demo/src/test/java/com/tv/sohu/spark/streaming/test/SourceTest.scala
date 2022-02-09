package com.tv.sohu.spark.streaming.test

import com.tv.sohu.spark.streaming.dm3.source.{AppLog, LogParse}

import java.util.regex.{Matcher, Pattern}

object SourceTest {
  def main(args: Array[String]): Unit = {
//    regexTest()
    // vv 正常日志
    val log = "110.184.20.233 - mb.hd.sohu.com.cn [25/Mar/2021:12:59:53 +0800] \"GET /mvv.gif?uid=a2b7a8ec3faab44487cd8d6f8ab46ef926&tkey=824a6c16038adb589833e8c738229c1896731128&vid=3097316&tvid=84207543&mtype=100&mfov=MX6&memo=%7B%22playstyle%22%3A1%2C%22page_transition%22%3A0%2C%22column_id%22%3A%220004%22%2C%22scn%22%3A%2202%22%2C%22pg%22%3A%2270000%22%2C%22isfee%22%3A0%2C%22isvr%22%3A0%7D&passport=&ltype=0&cv=8.5.3&mzcv=1.0.9&mos=2&mosv=7.1.1&pro=86&mfo=Meizu&webtype=WiFi&time=1616648392827&channelid=1006036133&sim=1&playlistid=9138015&catecode=101150%3B101157&preid=&newuser=0&enterid=0&startid=1616643152409&abmod=9998A&guid=ddd830a9df35d3c29e1d6d368729cec0&mnc=460000&build=03021820&oaid=a2b7a8ec3faab44487cd8d6f8ab46ef9&msg=caltime&type=vrs&playtime=2160&td=3031&version=1&codetype=2&cateid=2&company=&channeled=1000160005&language=&area=6&wtype=1&playid=1616646148486&playmode=1&screen=1&pid=&site=1&ugcode2=MCFOTk%3DqPbYJm86cdVtFtW_JrNUwe5dWwFNY2Id8cND7theLSexm1US9j6k68NcgPu39rg_&isp2p=0&androidId=79da637582e1a9ce&imei=863026034043175 HTTP/1.1\" 204 0 \"com.sohu.sohuvideo.meizu/1000009 (Linux; U; Android 7.1.1; zh_CN; MX6; Build/NMF26O)\"";
    val appLog = LogParse.formatLog(log);
    println(appLog);

    // 访客日志
    val log2 = "110.184.20.233 - mb.hd.sohu.com.cn [25/Mar/2021:12:59:53 +0800] \"GET /mvv.gif?uid=guest_a2b7a8ec3faab44487cd8d6f8ab46ef926&tkey=824a6c16038adb589833e8c738229c1896731128&vid=3097316&tvid=84207543&mtype=100&mfov=MX6&memo=%7B%22playstyle%22%3A1%2C%22page_transition%22%3A0%2C%22column_id%22%3A%220004%22%2C%22scn%22%3A%2202%22%2C%22pg%22%3A%2270000%22%2C%22isfee%22%3A0%2C%22isvr%22%3A0%7D&passport=&ltype=0&cv=8.5.3&mzcv=1.0.9&mos=2&mosv=7.1.1&pro=86&mfo=Meizu&webtype=WiFi&time=1616648392827&channelid=1006036133&sim=1&playlistid=9138015&catecode=101150%3B101157&preid=&newuser=0&enterid=0&startid=1616643152409&abmod=9998A&guid=ddd830a9df35d3c29e1d6d368729cec0&mnc=460000&build=03021820&oaid=a2b7a8ec3faab44487cd8d6f8ab46ef9&msg=caltime&type=vrs&playtime=2160&td=3031&version=1&codetype=2&cateid=2&company=&channeled=1000160005&language=&area=6&wtype=1&playid=1616646148486&playmode=1&screen=1&pid=&site=1&ugcode2=MCFOTk%3DqPbYJm86cdVtFtW_JrNUwe5dWwFNY2Id8cND7theLSexm1US9j6k68NcgPu39rg_&isp2p=0&androidId=79da637582e1a9ce&imei=863026034043175 HTTP/1.1\" 204 0 \"com.sohu.sohuvideo.meizu/1000009 (Linux; U; Android 7.1.1; zh_CN; MX6; Build/NMF26O)\"";
    val appLog2 = LogParse.formatLog(log2);
    println(appLog2);

    // uid 为空的情况
    val log3 = "223.104.24.8 - mb.hd.sohu.com.cn [25/Mar/2021:13:12:50 +0800] \"GET /mvv.gif?msg=playCount&uid=&vid=0&type=local&playtime=0&ltype=0&mtype=6&cv=&mos=2&mosv=&pro=1&mfo=&mfov=&td=0&webtype=Unknown&memo=&version=0&time=1615931395677&passport=&cateid=0&company=&channeled=1000120010&playlistid=0&language=&area=&wtype=3&channelid=873&sim=1&catecode=&preid=&newuser=1&enterid=0&playid=1615931395677&startid=1615931390271&playmode=2&screen=1&loc= HTTP/1.1\" 204 0 \"-\""
    val appLog3 = LogParse.formatLog(log3);
    println(appLog3);
  }

  //测试正则表达式
  def regexTest()={
     val patternApp = Pattern.compile("^([0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.*) - (.*) \\[(.*)\\] \"GET /.*.gif\\?(.*) HTTP.*\" [0-9]{3} [0-9]{1,5} \"(.*)\"$")
    val log = "110.184.20.233 - mb.hd.sohu.com.cn [25/Mar/2021:12:59:53 +0800] \"GET /mvv.gif?uid=a2b7a8ec3faab44487cd8d6f8ab46ef926&tkey=824a6c16038adb589833e8c738229c1896731128&vid=3097316&tvid=84207543&mtype=100&mfov=MX6&memo=%7B%22playstyle%22%3A1%2C%22page_transition%22%3A0%2C%22column_id%22%3A%220004%22%2C%22scn%22%3A%2202%22%2C%22pg%22%3A%2270000%22%2C%22isfee%22%3A0%2C%22isvr%22%3A0%7D&passport=&ltype=0&cv=8.5.3&mzcv=1.0.9&mos=2&mosv=7.1.1&pro=86&mfo=Meizu&webtype=WiFi&time=1616648392827&channelid=1006036133&sim=1&playlistid=9138015&catecode=101150%3B101157&preid=&newuser=0&enterid=0&startid=1616643152409&abmod=9998A&guid=ddd830a9df35d3c29e1d6d368729cec0&mnc=460000&build=03021820&oaid=a2b7a8ec3faab44487cd8d6f8ab46ef9&msg=caltime&type=vrs&playtime=2160&td=3031&version=1&codetype=2&cateid=2&company=&channeled=1000160005&language=&area=6&wtype=1&playid=1616646148486&playmode=1&screen=1&pid=&site=1&ugcode2=MCFOTk%3DqPbYJm86cdVtFtW_JrNUwe5dWwFNY2Id8cND7theLSexm1US9j6k68NcgPu39rg_&isp2p=0&androidId=79da637582e1a9ce&imei=863026034043175 HTTP/1.1\" 204 0 \"com.sohu.sohuvideo.meizu/1000009 (Linux; U; Android 7.1.1; zh_CN; MX6; Build/NMF26O)\"";
    val matcher = patternApp.matcher(log);
    print(matcher.matches().toString)
  }
}

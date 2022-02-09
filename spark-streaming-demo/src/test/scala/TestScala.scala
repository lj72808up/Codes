import java.text.SimpleDateFormat
import java.util.Calendar

object TestScala {

  def main(args: Array[String]): Unit = {
    val days = Array(1028)
    days.foreach(day => {
      println(s"hdfs dfs -rm /user/vetl/hive/warehouse/t_app_bd_stat/p_day=2021$day/*")
      println(s"hdfs dfs -cp /user/vetl/hive/warehouse/t_app_bd_stat/p_day=1999$day/* /user/vetl/hive/warehouse/t_app_bd_stat/p_day=2021$day")
    })

    days.foreach(
      day => {
        println(s"select count(1) from (select uid from t_app_bd_stat where p_day=2021$day and newuser=1 group by uid ) t1 limit 1 ")
      }
    )

    def getDaytimeTime(startDate: String, endDate: String): Array[String] = {
      val format = new SimpleDateFormat("yyyyMMdd");
      val sDate = format.parse(startDate);
      val eDate = format.parse(endDate);
      var start = sDate.getTime(); //获得毫秒数
      val end = eDate.getTime();
      if (start > end) {
        return null;
      }
      val calendar = Calendar.getInstance();
      calendar.setTime(sDate); //降序 setTime(eDate)
      var res = Array[String]();
      while (end >= start) {
//        println( format.format(calendar.getTime()))
        res = res :+ format.format(calendar.getTime())
        calendar.add(Calendar.DAY_OF_MONTH, 1); //降序 add(Calender.DAY_OF_MONTH, -1)
        start = calendar.getTimeInMillis(); //降序 end = calender.getTimeInMillis();
      }
      res
    }

    val days2 = getDaytimeTime("20210918", "20211017")
    days2.foreach(
      day => println(s"curl  http://10.19.117.143:8080/flushCache?prefix=DM4_&suffix=$day&t=dmxdev&urlUsername=zhanhongliu")
    )
  }
}

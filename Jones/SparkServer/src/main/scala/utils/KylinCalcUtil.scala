package utils

import java.io.{File, FileOutputStream, OutputStreamWriter, PrintWriter}
import java.nio.charset.StandardCharsets
import java.sql.{Driver, DriverManager}
import java.util.Properties

import adtl.platform.Dataproto.QueryInfo
import org.joda.time.DateTime

case class KylinCalcUtil(baseUrl:String, user:String, passwd:String) {
  Class.forName("org.apache.kylin.jdbc.Driver")
  val info = new Properties()
  info.put("user", user)
  info.put("password", passwd)
  def RunSQLSaveLocalCSV(query:QueryInfo, respath:String) = {
    val savePathTmp: String = s"${respath}/${DateTime.now.toString("yyyyMMdd")}/${query.queryId}"
    val saveFile = new File("/" + savePathTmp)
    if(!saveFile.getParentFile.exists()) saveFile.getParentFile.mkdirs()
    //val conn = driver.connect(baseUrl,info)
    if(saveFile.exists() && saveFile.isFile) saveFile.delete()
    val conn = DriverManager.getConnection(baseUrl, info)
    val state = conn.createStatement()
    val resultSet = state.executeQuery(query.sql(0))
    val resultMeta = resultSet.getMetaData
    val metaHeader = query.dimensions.toList ::: query.outputs.toList
    import java.io.OutputStreamWriter
    import java.io.PrintWriter
    val writer = new PrintWriter(new OutputStreamWriter(new FileOutputStream(saveFile), StandardCharsets.UTF_8), true)

    (0 until metaHeader.length).foreach{
      idx=>
        if(idx == 0)
          writer.print(metaHeader(idx))
        else
          writer.print("\t" + metaHeader(idx))
    }
    writer.println()

    while(resultSet.next()) {
      (1 to resultMeta.getColumnCount).foreach{
        idx=>
          if(idx == resultMeta.getColumnCount)
            writer.println(resultSet.getString(idx))
          else
            writer.print(resultSet.getString(idx) + "\t")
      }
    }
    writer.close()
    savePathTmp
  }
}

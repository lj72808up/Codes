package utils

import adtl.platform.Dataproto.QueryInfo

object TestUtil {

  def printQuery(query: QueryInfo) {
    println("------------ Query --------------")
    println("queryId" + query.queryId) // 1
    println("session" + query.session) // 2
    println("sql" + query.sql) // 3
    println("status" + query.status) // 4
    println("lastUpdate" + query.lastUpdate) // 5
    println("engine" + query.engine) // 6
    println("dimension" + query.dimensions) // 7
    println("dataPath" + query.dataPath) // 8
    println("exceptionMsg" + query.exceptionMsg) // 9
  }
}

package com.test

import org.apache.hadoop.hbase.{HBaseConfiguration, TableName}
import org.apache.hadoop.hbase.client.{ConnectionFactory, Get, Scan}

object TestScan {
  def main(args: Array[String]): Unit = {
    val conf1 = HBaseConfiguration.create
    val connection = ConnectionFactory.createConnection(conf1)
    val table = connection.getTable(TableName.valueOf("Contacts"))
    try {
      import org.apache.hadoop.hbase.client.ResultScanner
      import org.apache.hadoop.hbase.util.Bytes
      // Instantiating the Scan class// Instantiating the Scan class

      val scan = new Scan

      // Scanning the required columns
      scan.addColumn(Bytes.toBytes("info"), Bytes.toBytes("sha1"))

      // Getting the scan result
      val scanner = table.getScanner(scan)

      // Reading values from scan result
      var result = scanner.next
      while (result != null) {
        println("Found row : " + result)
        result = scanner.next
      }
      scanner.close()
    }
    finally {
      table.close
      connection.close()
    }
  }
}

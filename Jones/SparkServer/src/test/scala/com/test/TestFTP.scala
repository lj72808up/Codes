package com.test

import java.io.FileInputStream

import org.apache.commons.net.ftp.{FTP, FTPClient}

object TestFTP {


//  public static long copyLarge(InputStream input, OutputStream output, byte[] buffer)
//  throws IOException {
//    long count = 0;
//    int n = 0;
//    while (EOF != (n = input.read(buffer))) {
//      output.write(buffer, 0, n);
//      count += n;
//    }
//    return count;
//  }

  def main(args: Array[String]): Unit = {
    val ftp = new FTPClient()
    ftp.setControlEncoding("utf-8")
    ftp.connect("10.140.38.140") // ROC申请域名并使用域名

    println(ftp.login("ftp",""))
    val in = new FileInputStream("hs_err_pid2480.log")
    val in2 = new FileInputStream("hs_err_pid2480.log")
    val in3 = new FileInputStream("hs_err_pid2480.log")
    ftp.setFileType(FTP.BINARY_FILE_TYPE)

    val al = new java.util.ArrayList[FileInputStream]
    al.add(in)
    al.add(in2)
    al.add(in3)
    import scala.collection.JavaConversions._
    val out = ftp.storeFileStream("datacenter/20201201/中文.txt")
    val by =  new Array[Byte](1024 * 4)
    var n = 0
    al.foreach{
      ix =>
        while(n != -1) {
          n = ix.read(by)
          if (n != -1) {
            out.write(by, 0, n)
          }
        }
        ix.close()
    }
    out.close()
    ftp.logout()
  }
}

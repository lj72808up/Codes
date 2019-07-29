import com.test.utils.CommandUtil

object TestCommand {
  def main(args: Array[String]): Unit = {
    f2()
  }

  def f1(): Unit = {
    val wd = "dir1"
    val file = "start.sh"
    val c = new CommandUtil
    println(c.compactShell(file,wd))
  }

  def f2(): Unit = {
    val cmd = "bash start.sh 0613 0614"
    val wd = "dir1"

    val bu = new CommandUtil
    val is = bu.executeStrBash(cmd,workdir=wd)
    bu.cmdOutput(is,println)
  }
}

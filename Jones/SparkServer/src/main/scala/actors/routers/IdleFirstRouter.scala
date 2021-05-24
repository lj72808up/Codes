package actors.routers

import java.util.concurrent.ThreadLocalRandom

import actors.TaskProxyActor
import akka.actor.{ActorRef, ActorSystem, SupervisorStrategy}
import akka.dispatch.Dispatchers
import akka.routing.{ActorRefRoutee, NoRoutee, Pool, Resizer, Routee, Router, SmallestMailboxRoutingLogic}
import com.typesafe.config.Config
import utils.PrintLogger

import scala.annotation.tailrec
import scala.collection.immutable

final case class IdleFirstPool(val nrOfInstances: Int, override val resizer: Option[Resizer] = None,
                                 override val supervisorStrategy: SupervisorStrategy = Pool.defaultSupervisorStrategy,
                                 override val routerDispatcher: String = Dispatchers.DefaultDispatcherId,
                                 override val usePoolDispatcher: Boolean = false
                                ) extends Pool {

  def this(config: Config) =
    this(
      nrOfInstances = config.getInt("nr-of-instances"),
      resizer = Resizer.fromConfig(config),
      usePoolDispatcher = config.hasPath("pool-dispatcher"))

  def this(nr: Int) = this(nrOfInstances = nr)

  override def createRouter(system: ActorSystem): Router = new Router(IdleFirstRoutingLogic())

  override def nrOfInstances(sys: ActorSystem): Int = this.nrOfInstances

  /**
    * Setting the supervisor strategy to be used for the “head” Router actor.
    */
  def withSupervisorStrategy(strategy: SupervisorStrategy): IdleFirstPool = copy(supervisorStrategy = strategy)

  /**
    * Setting the resizer to be used.
    */
  def withResizer(resizer: Resizer): IdleFirstPool = copy(resizer = Some(resizer))

  /**
    * Setting the dispatcher to be used for the router head actor,  which handles
    * supervision, death watch and router management messages.
    */
  def withDispatcher(dispatcherId: String): IdleFirstPool = copy(routerDispatcher = dispatcherId)

//  /**
//    * Uses the resizer and/or the supervisor strategy of the given RouterConfig
//    * if this RouterConfig doesn't have one, i.e. the resizer defined in code is used if
//    * resizer was not defined in config.
//    */
//  override def withFallback(other: RouterConfig): RouterConfig = this.overrideUnsetConfig(other)

}


class IdleFirstRoutingLogic extends SmallestMailboxRoutingLogic with PrintLogger{
  override def select(message: Any, routees: immutable.IndexedSeq[Routee]): Routee = {
    if (routees.isEmpty) NoRoutee
    else this.selectNext(routees)
  }

  // Worst-case a 2-pass inspection with mailbox size checking done on second pass, and only until no one empty is found.
  // Lowest score wins, score 0 is autowin
  // If no actor with score 0 is found, it will return that, or if it is terminated, a random of the entire set.
  //   Why? Well, in case we had 0 viable actors and all we got was the default, which is the DeadLetters, anything else is better.
  // Order of interest, in ascending priority:
  // 1. The NoRoutee
  // 2. A Suspended ActorRef
  // 3. An ActorRef with unknown mailbox size but with one message being processed
  // 4. An ActorRef with unknown mailbox size that isn't processing anything
  // 5. An ActorRef with a known mailbox size
  // 6. An ActorRef without any messages
  @tailrec private def selectNext(targets: immutable.IndexedSeq[Routee],
                                  proposedTarget: Routee = NoRoutee,
                                  currentScore: Long = Long.MaxValue, // 得分最低的actor胜出
                                  at: Int = 0): Routee = {
    if (targets.isEmpty)
      NoRoutee
    else if (at >= targets.size) {
      // 如果 proposedTarget 已经 stop , 就随机选择一个; 否则选择 proposedTarget
      if (super.isTerminated(proposedTarget)) {
        printAndReturn(targets(ThreadLocalRandom.current.nextInt(targets.size)))
      } else {
        printAndReturn(proposedTarget)
      }
    } else {
      val target = targets(at) // 查看第 at 个 actor
      val newScore: Long =
        if (super.isSuspended(target)) { // 目标actor挂起
          Long.MaxValue - 1
        } else { //Just about better than the DeadLetters
          // 没有正在处里消息的actor得分为0, 正在处里消息的actor得分为已经处里的秒数
          val processingScore = if (super.isProcessingMessage(target)) {
            val path = this.getPath(target)
            val lastStartTime = TaskProxyActor.JobActorTimeMapping.getOrDefault(path, 0L)
            val during = (now - lastStartTime) / 1000
            info(s"${path} 正在执行任务, 持续了 ${during} 秒")
            during
          } else 0L

          val messageScore = if (!hasMessages(target)) {
            0L
          } else { //Race between hasMessages and numberOfMessages here, unfortunate the numberOfMessages returns 0 if unknown
            val mailBoxCnt = numberOfMessages(target)
            mailBoxCnt * predictTimePerMsg
          }

          processingScore + messageScore
        }

      if (newScore == 0) {
        printAndReturn(target)
      } else if (newScore < 0 || newScore >= currentScore) { // 接下来对比目标actor和位置在at的actor
        this.selectNext(targets, proposedTarget, currentScore, at + 1)
      } else {
        this.selectNext(targets, target, newScore, at + 1)
      }
    }
  }

  protected def getPath(a: Routee): String = a match {
    case ActorRefRoutee(x: ActorRef) => x.path.toString
    case _ => ""
  }

  protected def now(): Long = {
    System.currentTimeMillis()
  }

  protected def predictTimePerMsg: Long = {
    5 * 60 // 每个 mailbox 中的任务预计执行300秒
  }

  private def printAndReturn (a: Routee) :Routee = {
    info(s"执行任务的actor: ${getPath(a)}")
    a
  }
}

object IdleFirstRoutingLogic {
  def apply(): IdleFirstRoutingLogic = new IdleFirstRoutingLogic
}
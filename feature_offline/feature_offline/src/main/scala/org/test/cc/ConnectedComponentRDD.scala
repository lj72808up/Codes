package org.baidu.cc

/**
  * Copyright (c) 2017 Kwartile, Inc., http://www.kwartile.com
  * Permission is hereby granted, free of charge, to any person obtaining a copy
  * of this software and associated documentation files (the "Software"), to deal
  * in the Software without restriction, including without limitation the rights
  * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
  * copies of the Software, and to permit persons to whom the Software is
  * furnished to do so, subject to the following conditions:
  *
  * The above copyright notice and this permission notice shall be included in all
  * copies or substantial portions of the Software.
  *
  * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
  * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
  * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
  * SOFTWARE.
  */

/**
  * Map-reduce implementation of Connected Component
  * Given lists of subgraphs, returns all the nodes that are connected.
  */

import org.apache.spark.{RangePartitioner, SparkConf, SparkContext}
import org.apache.spark.rdd.RDD
import org.apache.spark.storage.StorageLevel

import scala.annotation.tailrec
import scala.collection.mutable


object ConnectedComponent extends Serializable {

  /**
    * Applies Small Star operation on RDD of nodePairs
    * @param nodePairs on which to apply Small Star operations
    * @return new nodePairs after the operation and conncectivy change count
    */
  private def smallStar(nodePairs: RDD[(Long, Long)]): (RDD[(Long, Long)], Int) = {

    /**
      * generate RDD of (self, List(neighbors)) where self > neighbors
      * E.g.: nodePairs (1, 4), (6, 1), (3, 2), (6, 5)
      * will result into (4, List(1)), (6, List(1)), (3, List(2)), (6, List(5))
      */
    val neighbors = nodePairs.map(x => {
      val (self, neighbor) = (x._1, x._2)  // 每个(u,v), 如果u>v,则发出(u,v); 如果u<=v, 则发出(v,u)
      if (self > neighbor)
        (self, neighbor)
      else
        (neighbor, self)
    })
    /*    neighbors:
        (1,1)
        (8,1)
        (7,7)
        (8,7)
        (8,1)
        (9,1)
        (9,8)
        (5,5)
        (8,5)*/

    /**
      * reduce on self to get list of all its neighbors.
      * E.g: (4, List(1)), (6, List(1)), (3, List(2)), (6, List(5))
      * will result into (4, List(1)), (6, List(1, 5)), (3, List(2))
      * Note:
      * (1) you may need to tweak number of partitions.
      * (2) also, watch out for data skew. In that case, consider using rangePartitioner
      */
    val empty = mutable.HashSet[Long]()
    val allNeighbors = neighbors.aggregateByKey(empty)(
      (lb, v) => lb += v,
      (lb1, lb2) => lb1 ++ lb2
    )
    /*    allNeighbors:
        (1,Set(1))
        (7,Set(7))
        (8,Set(1, 5, 7))
        (9,Set(1, 8))
        (5,Set(5))*/

    /**
      * Apply Small Star operation on (self, List(neighbor)) to get newNodePairs and count the change in connectivity
      */

    val newNodePairsWithChangeCount = allNeighbors.map(x => {
      val self = x._1
      val neighbors = x._2.toList
      val minNode = argMin(self :: neighbors)  // 找到包含自己在内的最小节点
      val newNodePairs = (self :: neighbors).map(neighbor => {  // 对包含自己在内的及其邻居, 发射出 (v,min)
        (neighbor, minNode)
      }).filter(x => {  // 只保留( 邻居<=自己 && 邻居!=min的 ) || (自己和邻居相等)
        val neighbor = x._1
        val minNode = x._2
        (neighbor <= self && neighbor != minNode) || (self == neighbor)
      })
      val uniqueNewNodePairs = newNodePairs.toSet.toList    // 对newNodePairs执行去重

      /**
        * We count the change by taking a diff of the new node pairs with the old node pairs
        */
      val connectivityChangeCount = (uniqueNewNodePairs diff neighbors.map((self, _))).length
      (uniqueNewNodePairs, connectivityChangeCount)
    }).persist(StorageLevel.MEMORY_AND_DISK_SER)
    /*
        newNodePairsWithChangeCount:
        (List((1,1)),0)
        (List((7,7)),0)
        (List((8,1), (5,1), (7,1)),2)
        (List((9,1), (8,1)),1)
        (List((5,5)),0)
    */


    /**
      * Sum all the changeCounts
      */
    val totalConnectivityCountChange = newNodePairsWithChangeCount.mapPartitions(iter => {
      val (v, l) = iter.toSeq.unzip
      val sum = l.foldLeft(0)(_ + _)
      Iterator(sum)
    }).sum.toInt

    val newNodePairs = newNodePairsWithChangeCount.map(x => x._1).flatMap(x => x)
    newNodePairsWithChangeCount.unpersist(false)
    (newNodePairs, totalConnectivityCountChange)

    /* newNodePairs:
      * (1,1)
        (7,7)
        (8,1)
        (5,1)
        (7,1)
        (9,1)
        (8,1)
        (5,5)
      *
      */
  }

  /**
    * Apply Large Star operation on a RDD of nodePairs
    * @param nodePairs on which to apply Large Star operations
    * @return new nodePairs after the operation and conncectivy change count
    */
  private def largeStar(nodePairs: RDD[(Long, Long)]): (RDD[(Long, Long)], Int) = {

    /**
      * generate RDD of (self, List(neighbors))
      * E.g.: nodePairs (1, 4), (6, 1), (3, 2), (6, 5)
      * will result into (4, List(1)), (1, List(4)), (6, List(1)), (1, List(6)), (3, List(2)), (2, List(3)), (6, List(5)), (5, List(6))
      */

    val neighbors = nodePairs.flatMap(x => {
      val (self, neighbor) = (x._1, x._2)
      if (self == neighbor) /* 对比自己和邻居, 相等发一条, 不相等发正反两条 */
        List((self, neighbor))
      else
        List((self, neighbor), (neighbor, self))
    })
    // neighbors:
    //    (8,1)
    //    (1,8)
    //    (5,8)
    //    (8,5)
    //    (8,7)
    //    (7,8)
    //    (9,8)
    //    (8,9)

    /**
      * reduce on self to get list of all its neighbors.
      * E.g: (4, List(1)), (1, List(4)), (6, List(1)), (1, List(6)), (3, List(2)), (2, List(3)), (6, List(5)), (5, List(6))
      * will result into (4, List(1)), (1, List(4, 6)), (6, List(1, 5)), (3, List(2)), (2, List(3)), (5, List(6))
      * Note:
      * (1) you may need to tweak number of partitions.
      * (2) also, watch out for data skew. In that case, consider using rangePartitioner
      */

    val localAdd = (s: mutable.HashSet[Long], v: Long) => s += v
    val partitionAdd = (s1: mutable.HashSet[Long], s2: mutable.HashSet[Long]) => s1 ++= s2
    val allNeighbors = neighbors.aggregateByKey(mutable.HashSet.empty[Long]/*, rangePartitioner*/)(localAdd, partitionAdd)
    /* allNeighbors:
        (1,Set(8))
        (7,Set(8))
        (8,Set(9, 1, 5, 7))
        (9,Set(8))
        (5,Set(8))*/
    /**
      * Apply Large Star operation on (self, List(neighbor)) to get newNodePairs and count the change in connectivity
      *
      * scala> val neighbors = List((8,1),(1,8),(5,8),(8,5),(8,7),(7,8),(9,8),(8,9))
      * scala> val allNeighbors = List((1,Set(8)),(7,Set(8)),(8,Set(9, 1, 5, 7)),(9,Set(8)),(5,Set(8)))
      */

    val newNodePairsWithChangeCount = allNeighbors.map(x => {
      val self = x._1   // root节点
      val neighbors = x._2.toList   // 所有邻居节点形成的列表
      val minNode = argMin(self :: neighbors)    // 包含自己在内的最小节点
      val newNodePairs = (self :: neighbors).map(neighbor => {   // 包括自己在内的及其所有邻居, 发射(v,min)
        (neighbor, minNode)
      }).filter(x => {   // 只留v>=min的对儿
        val neighbor = x._1
        val minNode = x._2
        neighbor >= self
      })

      val uniqueNewNodePairs = newNodePairs.toSet.toList  // 对newNodePairs去重
      val connectivityChangeCount = (uniqueNewNodePairs diff neighbors.map((self, _))).length  // neighbors.map((self, _))的结果: 形成(self,每个neibor)的对儿
      (uniqueNewNodePairs, connectivityChangeCount)
    }).persist(StorageLevel.MEMORY_AND_DISK_SER)

    /*  newNodePairsWithChangeCount:
        (List((1,1), (8,1)),2)
        (List((7,7), (8,7)),2)
        (List((8,1), (9,1)),1)
        (List((9,8)),0)
        (List((5,5), (8,5)),2)
    */

    val totalConnectivityCountChange = newNodePairsWithChangeCount.mapPartitions(iter => {
      val (v, l) = iter.toSeq.unzip
      val sum = l.foldLeft(0)(_ + _)
      Iterator(sum)
    }).sum.toInt

    /**
      * Sum all the changeCounts
      */
    val newNodePairs = newNodePairsWithChangeCount.map(x => x._1).flatMap(x => x)  // 把newNodePairsWithChangeCount的结果打平
    /*    newNodePairs:
          (1,1)
          (8,1)
          (7,7)
          (8,7)
          (8,1)
          (9,1)
          (9,8)
          (5,5)
          (8,5)*/

    newNodePairsWithChangeCount.unpersist(false)
    (newNodePairs, totalConnectivityCountChange)
  }

  private def argMin(nodes: List[Long]): Long = {
    nodes.min(Ordering.by((node: Long) => node))
  }

  /**
    * Build nodePairs given a list of nodes.  A list of nodes represents a subgraph.
    * @param nodes that are part of a subgraph
    * @return nodePairs for a subgraph
    */
  private def buildPairs(nodes:List[Long]) : List[(Long, Long)] = {
    buildPairs(nodes.head, nodes.tail, null.asInstanceOf[List[(Long, Long)]])
  }

  @tailrec
  private def buildPairs(node: Long, neighbors:List[Long], partialPairs: List[(Long, Long)]) : List[(Long, Long)] = {
    if (neighbors.isEmpty) {
      if (partialPairs != null)
        List((node, node)) ::: partialPairs
      else
        List((node, node))
    } else if (neighbors.length == 1) {
      val neighbor = neighbors(0)
      if (node > neighbor)
        if (partialPairs != null) List((node, neighbor)) ::: partialPairs else List((node, neighbor))
      else
      if (partialPairs != null) List((neighbor, node)) ::: partialPairs else List((neighbor, node))
    } else {
      val newPartialPairs = neighbors.map(neighbor => {
        if (node > neighbor)
          List((node, neighbor))
        else
          List((neighbor, node))
      }).flatMap(x=>x)

      if (partialPairs != null)
        buildPairs(neighbors.head, neighbors.tail, newPartialPairs ::: partialPairs)
      else
        buildPairs(neighbors.head, neighbors.tail, newPartialPairs)
    }
  }

  /**
    * Implements alternatingAlgo.  Converges when the changeCount is either 0 or does not change from the previous iteration
    * @param nodePairs for a graph
    * @param largeStarConnectivityChangeCount change count that resulted from the previous iteration
    * @param smallStarConnectivityChangeCount change count that resulted from the previous iteration
    * @param didConverge flag to indicate the alorigth converged
    * @param currIterationCount counter to capture number of iterations
    * @param maxIterationCount maximum number iterations to try before giving up
    * @return RDD of nodePairs
    */

  @tailrec
  private def alternatingAlgo(nodePairs: RDD[(Long, Long)],
                              largeStarConnectivityChangeCount: Int, smallStarConnectivityChangeCount: Int, didConverge: Boolean,
                              currIterationCount: Int, maxIterationCount: Int): (RDD[(Long, Long)], Boolean, Int) = {

    val iterationCount = currIterationCount + 1
    if (didConverge)
      (nodePairs, true, currIterationCount)
    else if (currIterationCount >= maxIterationCount) {
      (nodePairs, false, currIterationCount)
    }
    else {

      val (nodePairsLargeStar, currLargeStarConnectivityChangeCount) = largeStar(nodePairs)

      val (nodePairsSmallStar, currSmallStarConnectivityChangeCount) = smallStar(nodePairsLargeStar)

      if ((currLargeStarConnectivityChangeCount == largeStarConnectivityChangeCount &&
        currSmallStarConnectivityChangeCount == smallStarConnectivityChangeCount) ||
        (currSmallStarConnectivityChangeCount == 0 && currLargeStarConnectivityChangeCount == 0)) {
        alternatingAlgo(nodePairsSmallStar, currLargeStarConnectivityChangeCount,
          currSmallStarConnectivityChangeCount, true, iterationCount, maxIterationCount)
      }
      else {
        alternatingAlgo(nodePairsSmallStar, currLargeStarConnectivityChangeCount,
          currSmallStarConnectivityChangeCount, false, iterationCount, maxIterationCount)
      }
    }
  }

  /**
    * Driver function
    * @param cliques list of nodes representing subgraphs (or cliques)
    * @param maxIterationCount maximum number iterations to try before giving up
    * @return Connected Components as nodePairs where second member of the nodePair is the minimum node in the component
    */
  def run(cliques:RDD[(Long, Long)], maxIterationCount: Int): (RDD[(Long, Long)], Boolean, Int) = {

    val nodePairs = cliques/*cliques.map(aClique => {
      buildPairs(aClique)
    }).flatMap(x=>x)*/

    val (cc, didConverge, iterCount) = alternatingAlgo(nodePairs, 9999999, 9999999, false, 0, maxIterationCount)

    if (didConverge) {
      (cc, didConverge, iterCount)
    } else {
      (null.asInstanceOf[RDD[(Long, Long)]], didConverge, iterCount)
    }
  }


  def main(args: Array[String]) = {

    val sparkConf = new SparkConf()
      .setAppName("ConnectedComponent")
      .setMaster("local[1]")

    val sc = new SparkContext(sparkConf)
    /*    val cliqueFile = args(0)
        val cliquesRec = sc.textFile(args(0))
        val cliques = cliquesRec.map(x => {
          val nodes = x.split("\\s+").map(y => y.toLong).toList
          nodes
        })*/

    val cliques = sc.parallelize(List((8l,1l),(5l,8l),(8l,7l),(9l,8l)))


    val (cc, didConverge, iterCount) = ConnectedComponent.run(cliques, 20)

    if (didConverge) {
      println("Converged in " + iterCount + " iterations")
      val cc2 = cc.map(x => {
        (x._2, List(x._1))
      })

      /**
        * Get all the nodes in the connected component/
        * We are using a rangePartitioner because the CliquesGenerator produces data skew.
        */
      /*val rangePartitioner = new RangePartitioner(cc2.getNumPartitions, cc2)
      val connectedComponents =  cc2.reduceByKey(rangePartitioner, (a, b) => {b ::: a})*/

      //connectedComponents.mapPartitionsWithIndex((index, iter) => {
      //  iter.toList.map(x => (index, x._1, x._2.size)).iterator
      //  }).collect.foreach(println)

      cc.collect().foreach(println)
      println("connected components")
      //connectedComponents.map(x => (x._2.length).toString + " " + x._1 + " " + x._2.sorted.mkString(" ")).saveAsTextFile(cliqueFile + "_cc_out")
    }
    else {
      println("Max iteration reached.  Could not converge")
    }
  }
}


package com.test

import org.apache.spark.sql.SparkSession
import org.apache.spark.sql.catalyst.analysis.UnresolvedRelation
import org.apache.spark.sql.catalyst.plans.logical._
import org.apache.spark.sql.execution.datasources.{CreateTable, HadoopFsRelation, LogicalRelation}
import java._


import org.apache.spark.sql.catalyst.catalog.HiveTableRelation
import utils.LinageHelper

object TestSparkSql {
  /**
    * GlobalLimit 10
    * +- LocalLimit 10
    * +- Aggregate [indus_first_name#16719], [indus_first_name#16719, sum(max_price#16673) AS sum(max_price)#16729, sum(price#16674) AS sum(price)#16730, sum(aa#16671) AS sum(aa)#16731, sum(clk#16672L) AS sum(clk)#16732L]
    * +- Join LeftOuter, (query_classify#16685 = query_classify#16718)
    * :- Join Inner, (pid#16706 = pid#16713)
    * :  :- SubqueryAlias a
    * :  :  +- Aggregate [pid#16706, query_classify#16685], [pid#16706, query_classify#16685, sum((cast(price#16703 as double) / cast(max_price#16704 as double))) AS aa#16671, count(click_id#16687) AS clk#16672L, sum(cast(max_price#16704 as double)) AS max_price#16673, sum(cast(price#16703 as double)) AS price#16674]
    * :  :     +- Filter ((dt#16712 >= 20190805) && (dt#16712 <= 20190811))
    * :  :        +- SubqueryAlias dwd_wsbidding_chargecln_di
    * :  :           +- Relation[pvid#16675,pid_src#16676,suid_yyid#16677,group_id#16678,ad_id#16679,account_id#16680,flag#16681,reserved_list#16682,keyword_list#16683,search_keyword#16684,query_classify#16685,pv_refer#16686,click_id#16687,creative_id#16688,plan_id#16689,acid_index#16690,kid_indus#16691,flow_flag#16692,click_refer#16693,click_region#16694,show_ip#16695,show_region#16696,domain#16697,user_agent#16698,... 14 more fields] orc
    * :  +- SubqueryAlias p
    * :     +- Distinct
    * :        +- Project [pid#16713]
    * :           +- Filter (second_channel#16715 = 手机QQ浏览器（优质）)
    * :              +- SubqueryAlias dim_pid_da
    * :                 +- HiveTableRelation `adtl_biz`.`dim_pid_da`, org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe, [pid#16713, first_channel#16714, second_channel#16715, typename#16716, account#16717]
    * +- SubqueryAlias c
    * +- Distinct
    * +- Project [query_classify#16718, indus_first_name#16719, indus_second_name#16720]
    * +- SubqueryAlias dw_dim_queryclassify
    * +- HiveTableRelation `adtl_biz`.`dw_dim_queryclassify`, org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe, [query_classify#16718, indus_first_name#16719, indus_second_name#16720]
    */

  val sql: String =
    """select
      |    c.indus_first_name, sum(a.max_price), sum(a.price), sum(a.aa), sum(a.clk)
      |from
      |  (
      |    select
      |      pid,
      |      query_classify, sum(price / max_price) as aa, count(click_id) as clk, sum(max_price) as max_price, sum(price) as price
      |    FROM
      |      adtl_biz.dwd_wsbidding_chargecln_di
      |    where
      |      dt between '20190805'
      |      and '20190811'
      |    group by
      |      pid,
      |      query_classify
      |  ) a
      |  inner join (
      |    SELECT
      |       t1.pid
      |    FROM
      |      adtl_biz.dim_pid_da t1
      |    LEFT JOIN
      |      default.query_word t2
      |    WHERE
      |      t1.second_channel = t2.part
      |  ) p on a.pid = p.pid
      |  left join (
      |    SELECT
      |      distinct query_classify, indus_first_name, indus_second_name
      |    FROM
      |      adtl_biz.dw_dim_queryclassify
      |  ) c on a.query_classify = c.query_classify
      |group by
      |  c.indus_first_name
      |limit
      |  10""".stripMargin


  def main2(args: Array[String]): Unit = {
    val spark = SparkSession.builder().master("local")
      .getOrCreate()

    val unresolvedPlan = spark.sessionState.sqlParser.parsePlan(sql)

    //      val resilvedPlan = df.queryExecution.analyzed   //  org.apache.spark.sql.catalyst.plans.logical.LogicalPlan --> resolved
    // // println(df.queryExecution.logical)
    val plan = unresolvedPlan
    val defaultDB = "hahaDB"
    val inputTable = new util.HashSet[DcTable]()
    val outputTable = new util.HashSet[DcTable]()

    val tableList = new util.ArrayList[String]()
    resolveLogic(plan, defaultDB, inputTable, outputTable, 0, tableList)

    println("end:" + tableList)
    // println(inputTable)
    // println(outputTable)

  }

  case class DcTable(db: String, tb: String)

  def resolveLogic(plan: LogicalPlan, defaultDB: String, inputTables: util.Set[DcTable], outputTables: util.Set[DcTable], level: Int, tableList: util.ArrayList[String]): Unit = {
    // println(s"<<<<<<<<<<<<<<<<<< $level >>>>>>>>>>>>>>>>>>")
    // println(plan)
    // println(s">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
    plan match {

      case plan: Project =>
        val project = plan.asInstanceOf[Project]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: Union => // 多个分支
        val project = plan.asInstanceOf[Union]
        for (child <- project.children) {
          resolveLogic(child, defaultDB, inputTables, outputTables, level + 1, tableList)
        }

      case plan: Join => // 两个分支
        val project = plan.asInstanceOf[Join]
        resolveLogic(project.left, defaultDB, inputTables, outputTables, level + 1, tableList)
        resolveLogic(project.right, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: Aggregate =>
        val project = plan.asInstanceOf[Aggregate]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: Filter =>
        val project = plan.asInstanceOf[Filter]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: Generate =>
        val project = plan.asInstanceOf[Generate]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: RepartitionByExpression =>
        val project = plan.asInstanceOf[RepartitionByExpression]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: SerializeFromObject =>
        val project = plan.asInstanceOf[SerializeFromObject]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: MapPartitions =>
        val project = plan.asInstanceOf[MapPartitions]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: DeserializeToObject =>
        val project = plan.asInstanceOf[DeserializeToObject]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: Repartition =>
        val project = plan.asInstanceOf[Repartition]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: Deduplicate =>
        val project = plan.asInstanceOf[Deduplicate]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: Window =>
        val project = plan.asInstanceOf[Window]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: MapElements =>
        val project = plan.asInstanceOf[MapElements]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: TypedFilter =>
        val project = plan.asInstanceOf[TypedFilter]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: Distinct =>
        val project = plan.asInstanceOf[Distinct]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: SubqueryAlias =>
        val project = plan.asInstanceOf[SubqueryAlias]
        val childInputTables = new util.HashSet[DcTable]()
        val childOutputTables = new util.HashSet[DcTable]()

        // 增加节点数组
        tableList.add(project.alias)
        resolveLogic(project.child, defaultDB, childInputTables, childOutputTables, level + 1, tableList)
        // 删除节点数组
        tableList.remove(tableList.size() - 1)

        if (childInputTables.size() > 0) {
          val iterator = childInputTables.iterator()
          while (iterator.hasNext()) {
            inputTables.add(iterator.next())
          }
        } else {
          inputTables.add(DcTable(defaultDB, project.alias))
        }

//      case plan: UnresolvedRelation => // UnresolvedRelation 是AST 树的叶子节点
//        val project = plan.asInstanceOf[UnresolvedRelation]
//        // 增加节点
//        tableList.add(project.tableName)
//        println(tableList)
//        val dcTable = DcTable(project.tableIdentifier.database.getOrElse(defaultDB), project.tableIdentifier.table)
//        inputTables.add(dcTable)
//        // 删除节点
//        tableList.remove(tableList.size() - 1)


        // 当 LogicPlan 被 resolve 后, 之前的 UnresolvedRelation 变为 HiveTableRelation
      case plan: HiveTableRelation =>
        val project = plan.asInstanceOf[HiveTableRelation]
        // 增加节点
        val tb = project.tableMeta.identifier
        tableList.add(s"${tb.database}.${tb.table}")
        println(tableList)
        val dcTable = DcTable(tb.database.getOrElse(defaultDB), tb.table)
        inputTables.add(dcTable)
        // 删除节点
        tableList.remove(tableList.size() - 1)


      // spark sql 涉及的第一个表, 作为 datasource 读入, 是BaseRelation子类中的一种 (hive表使用HadoopFsRelation)
      case plan: LogicalRelation=>
        val project = plan.asInstanceOf[LogicalRelation]
        println(s"source name : ${project.catalogTable.get.identifier}")
        // org.apache.spark.sql.execution.datasources.HadoopFsRelation
        println(s"实际的读入类: ${project.relation.getClass.getCanonicalName}")  // org.apache.spark.sql.execution.datasources.HadoopFsRelation


      case plan: GlobalLimit =>
        val project = plan.asInstanceOf[GlobalLimit]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case plan: LocalLimit =>
        val project = plan.asInstanceOf[LocalLimit]
        resolveLogic(project.child, defaultDB, inputTables, outputTables, level + 1, tableList)

      case `plan` => plan.getClass.getCanonicalName// println("******child plan******:\n" + plan)

    }
  }


  def main1(args: Array[String]): Unit = {
    // 页面在这里启动, 但是你没有job
    val spark = SparkSession.builder().master("local")
      .getOrCreate()
    // spark.sql 内部:
    // Dataset.ofRows(self, sessionState.sqlParser.parsePlan(sqlText))
    // (1) sessionState.sqlParser.parsePlan(sqlText) => antlr4 解析 AST 树, 生成 UnResolved LogicalPlan转
    // (2) Dataset.ofRows(...) 使用 Analyser 绑定sql中的表的元数据到 AST 的 TreeNode 上, 形成 Resolved LogicalPlan

    val df = spark.sql("select 1")
    //    df.queryExecution.logical.resolved
    // 有action动作时, 触发 optimizer (RBO, CBO), 生成 optimized LogicalPlan
    // SparkPlan 将 LogicalPlan 转换成 PhysicalPlan;
    // prepareForExecution()将 PhysicalPlan 转换成可执行物理计划;
    // 使用 execute()执行可执行物理计划
    df.collect()

    Thread.sleep(60 * 60 * 1000)
  }
}

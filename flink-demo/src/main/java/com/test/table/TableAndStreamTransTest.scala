package com.test.table

import org.apache.flink.api.scala.typeutils.Types
import org.apache.flink.streaming.api.scala.StreamExecutionEnvironment
import org.apache.flink.table.api.{Schema, Table}
import org.apache.flink.table.api.bridge.scala.StreamTableEnvironment
import org.apache.flink.table.connector.ChangelogMode
import org.apache.flink.types.{Row, RowKind}

import java.time.Instant

object TableAndStreamTransTest {
  def main(args: Array[String]): Unit = {
    //    test1()
    //    test2()
    //    test3()
    test4()
  }

  /**
   * （1）起始框架：Table 和 DataStream （insert-only,changlog） 互转
   */
  def test1(): Unit = {
    // create environments of both APIs（table API 和 DataStream API）
    val env = StreamExecutionEnvironment.getExecutionEnvironment
    // 如果有设置 checkpoint 的操作，在 TableEnv 声明之前设置
    //    env.getCheckpointConfig.setCheckpointingMode(CheckpointingMode.EXACTLY_ONCE)

    val tableEnv: StreamTableEnvironment = StreamTableEnvironment.create(env)

    // 创建 DataStream
    //    val dataStream = env.fromElements("Alice", "Bob", "John")
    val dataStream = env.fromElements(
      Row.of("Alice", Int.box(12)),
      Row.of("Bob", Int.box(10)),
      Row.of("Alice", Int.box(100))
    )(Types.ROW(Types.STRING, Types.INT))

    // DataStream 转成 insert-only Table
    //    val inputTable = tableEnv.fromDataStream(dataStream)
    val inputTable = tableEnv.fromDataStream(dataStream).as("name", "score")

    // 注册 table 为 view 才能查询
    tableEnv.createTemporaryView("InputTable", inputTable)
    //    val resultTable = tableEnv.sqlQuery("SELECT name FROM InputTable")
    val resultTable = tableEnv.sqlQuery("SELECT name, SUM(score) FROM InputTable GROUP BY name")


    // insert-only Table 转成 DataStream
    //    val resultStream = tableEnv.toDataStream(resultTable)
    // 因为table有聚合操作，使得table和dataStream的转换更像是 cdc 推送数据。例如把聚合后的 update 操作推给 mysql sink
    val resultStream = tableEnv.toChangelogStream(resultTable)
    resultStream.print()

    env.execute() // DataStream 启动 source->sink 要用 StreamExecutionEnvironment。execute
  }

  /**
   * （2） 什么时候 Table Query 开始执行？
   *  - TableEnvironment.executeSql() is called. This method is used for executing a given statement, and the sql query is translated immediately once this method is called.
   *  - TablePipeline.execute() is called. This method is used for executing a source-to-sink pipeline, and the Table API program is translated immediately once this method is called.
   *  - Table.execute() is called. This method is used for collecting the table content to the local client, and the Table API is translated immediately once this method is called.
   *  - StatementSet.execute() is called. A TablePipeline (emitted to a sink through StatementSet.add()) or an INSERT statement (specified through StatementSet.addInsertSql()) will be buffered in StatementSet first. They are transformed once StatementSet.execute() is called. All sinks will be optimized into one DAG.
   *  - A Table is translated when it is converted into a DataStream (see Integration with DataStream). Once translated, it’s a regular DataStream program and is executed when StreamExecutionEnvironment.execute() is called.
   */
  def test2(): Unit = {

    val env = StreamExecutionEnvironment.getExecutionEnvironment
    val tableEnv: StreamTableEnvironment = StreamTableEnvironment.create(env)
    val dataStream = env.fromElements(
      Row.of("Alice", Int.box(12)),
      Row.of("Bob", Int.box(10)),
      Row.of("Alice", Int.box(100))
    )(Types.ROW(Types.STRING, Types.INT))

    // 创建输入输出
    val inputTable: Table = tableEnv.fromDataStream(dataStream).as("name", "score")
    tableEnv.createTemporaryView("InputTable", inputTable)
    tableEnv.executeSql("CREATE TABLE OutputTable(name STRING,score INT) WITH ('connector' = 'print')")
    tableEnv.executeSql("CREATE TABLE OutputTable2(name STRING,score INT) WITH ('connector' = 'print')")

    // tableEnv 没有像 streamEnv 一样的 execute()，而是如下几种直接启动 source->sink 的方法：
    //    tableEnv.from("InputTable").executeInsert("OutputTable")                   // （1）
    //    tableEnv.executeSql("INSERT INTO OutputTable SELECT * FROM InputTable")    // （2）

    /*tableEnv.createStatementSet()                                              // （3）多 sink
      .addInsert("OutputTable", tableEnv.from("InputTable"))
      .addInsert("OutputTable2", tableEnv.from("InputTable"))
      .execute()                                       */

    /*tableEnv.createStatementSet()                                              // （4）多 sink
      .addInsertSql("INSERT INTO OutputTable SELECT * FROM InputTable")
      .addInsertSql("INSERT INTO OutputTable2 SELECT * FROM InputTable")
      .execute()                                       */

    // 收集到本地 sink
    //    tableEnv.from("InputTable").execute().print()                              // 本地（1）
    //    tableEnv.executeSql("SELECT * FROM InputTable").print()                    // 本地（2）

    // (1)
    // 增加一个带 printing sink 的 Table 分支到 StreamExecutionEnvironment
    // 这个任务要用后面的 env.execute() 启动
    tableEnv.toDataStream(inputTable).print()

    // (2)
    // 执行一个 Table API： end-to-end pipeline 作为 Flink job, 然后收集到本地打印
    // thus (1) has still not been executed
    inputTable.execute().print()

    // executes the DataStream API pipeline with the sink defined in (1) as a Flink job, (2) was already running before
    env.execute()
  }

  /**
   * (3) 各种形式的 fromDataStream
   * 从 Table API 的视角， from / to DataStream 像是从一个 virtual table connector 读或写到一个 virtual table connector
   * 这个 virtual DataStream table connector 会给每一行暴露一个 rowtime 元数据属性
   */
  def test3(): Unit = {
    val env = StreamExecutionEnvironment.getExecutionEnvironment
    val tableEnv: StreamTableEnvironment = StreamTableEnvironment.create(env)

    import org.apache.flink.api.scala._ // for implicit conversions
    // create a DataStream
    val dataStream = env.fromElements(
      User("Alice", 4, Instant.ofEpochMilli(1000)),
      User("Bob", 6, Instant.ofEpochMilli(1001)),
      User("Alice", 10, Instant.ofEpochMilli(1002)))

    // === EXAMPLE 1 ===
    // fromDataStream 自动推断所有列名及类型，列名是 case class 的属性名
    val table1 = tableEnv.fromDataStream(dataStream)
    table1.printSchema()
    // prints:
    // (
    //  `name` STRING,
    //  `score` INT,
    //  `event_time` TIMESTAMP_LTZ(9)
    // )

    // === EXAMPLE 2 ===
    // 推断物理列名，并加了一个 proc_time 属性（处理时间）
    val table2 = tableEnv.fromDataStream(
      dataStream,
      Schema.newBuilder()
        .columnByExpression("proc_time", "PROCTIME()")
        .build())
    table2.printSchema()
    // prints:
    // (
    //  `name` STRING,
    //  `score` INT NOT NULL,
    //  `event_time` TIMESTAMP_LTZ(9),
    //  `proc_time` TIMESTAMP_LTZ(3) NOT NULL *PROCTIME* AS PROCTIME()
    //)

    // === EXAMPLE 3 ===
    // 增加了一个 rowtime 属性列，和一个自定义水印策略
    val table3 =
      tableEnv.fromDataStream(
        dataStream,
        Schema.newBuilder()
          .columnByExpression("rowtime", "CAST(event_time AS TIMESTAMP_LTZ(3))")
          .watermark("rowtime", "rowtime - INTERVAL '10' SECOND")
          .build())
    table3.printSchema()
    // prints:
    // (
    //  `name` STRING,
    //  `score` INT,
    //  `event_time` TIMESTAMP_LTZ(9),
    //  `rowtime` TIMESTAMP_LTZ(3) *ROWTIME* AS CAST(event_time AS TIMESTAMP_LTZ(3)),
    //  WATERMARK FOR `rowtime`: TIMESTAMP_LTZ(3) AS rowtime - INTERVAL '10' SECOND
    // )

    // === EXAMPLE 4 ===
    // derive all physical columns automatically
    // but access the stream record's timestamp for creating a rowtime attribute column
    // also rely on the watermarks generated in the DataStream API

    // we assume that a watermark strategy has been defined for `dataStream` before
    // (not part of this example)
    val table4 =
      tableEnv.fromDataStream(
        dataStream,
        Schema.newBuilder()
          .columnByMetadata("rowtime", "TIMESTAMP_LTZ(3)")
          .watermark("rowtime", "SOURCE_WATERMARK()") // 使用 datastream 上的水印策略
          .build())
    table4.printSchema()
    // prints:
    // (
    //  `name` STRING,
    //  `score` INT,
    //  `event_time` TIMESTAMP_LTZ(9),
    //  `rowtime` TIMESTAMP_LTZ(3) *ROWTIME* METADATA,
    //  WATERMARK FOR `rowtime`: TIMESTAMP_LTZ(3) AS SOURCE_WATERMARK()
    // )

    // === EXAMPLE 5 ===

    // define physical columns manually
    // in this example,
    //   - we can reduce the default precision of timestamps from 9 to 3
    //   - we also project the columns and put `event_time` to the beginning
    val table5 =
      tableEnv.fromDataStream(
        dataStream,
        Schema.newBuilder()
          .column("event_time", "TIMESTAMP_LTZ(3)") // 自己指定列对应的类型，自动推断是 TIMESTAMP_LTZ(9)
          .column("name", "STRING")
          .column("score", "INT")
          .watermark("event_time", "SOURCE_WATERMARK()")
          .build());
    table5.printSchema();
    // prints:
    // (
    //  `event_time` TIMESTAMP_LTZ(3) *ROWTIME*,
    //  `name` VARCHAR(200),
    //  `score` INT
    // )
    // note: the watermark strategy is not shown due to the inserted column reordering projection
  }

  /**
   * 4. 各种形式的 fromChangelogStream
   * - flink table 的运行时是一个 changelog 的处理器
   * - fromChangelogStream 的 DataStream 只能是 Row[RowKind] 类型。RowKind：insert ， upsert
   */
  def test4():Unit = {
    val env = StreamExecutionEnvironment.getExecutionEnvironment
    val tableEnv: StreamTableEnvironment = StreamTableEnvironment.create(env)
    // === EXAMPLE 1 ===
    // interpret the stream as a retract stream
    // create a changelog DataStream
    val dataStream = env.fromElements(
      Row.ofKind(RowKind.INSERT, "Alice", Int.box(12)),
      Row.ofKind(RowKind.INSERT, "Bob", Int.box(5)),
      Row.ofKind(RowKind.UPDATE_BEFORE, "Alice", Int.box(12)),
      Row.ofKind(RowKind.UPDATE_AFTER, "Alice", Int.box(100))
    )(Types.ROW(Types.STRING, Types.INT))

    // interpret the DataStream as a Table
    val table = tableEnv.fromChangelogStream(dataStream)

    // register the table under a name and perform an aggregation
    tableEnv.createTemporaryView("InputTable", table)
    tableEnv
      .executeSql("SELECT f0 AS name, SUM(f1) AS score FROM InputTable GROUP BY f0")
      .print()

    // prints:
    // +----+--------------------------------+-------------+
    // | op |                           name |       score |
    // +----+--------------------------------+-------------+
    // | +I |                            Bob |           5 |
    // | +I |                          Alice |          12 |
    // | -D |                          Alice |          12 |
    // | +I |                          Alice |         100 |
    // +----+--------------------------------+-------------+


    // === EXAMPLE 2 ===
    // interpret the stream as an upsert stream (without a need for UPDATE_BEFORE)
    // create a changelog DataStream
    val dataStream2 = env.fromElements(
      Row.ofKind(RowKind.INSERT, "Alice", Int.box(12)),
      Row.ofKind(RowKind.INSERT, "Bob", Int.box(5)),
      Row.ofKind(RowKind.UPDATE_AFTER, "Alice", Int.box(100))
    )(Types.ROW(Types.STRING, Types.INT))

    // interpret the DataStream as a Table
    // - fromChangelogStream(DataStream, Schema, ChangelogMode): Gives full control about how to interpret a stream as a changelog.
    // --               The passed ChangelogMode helps the planner to distinguish between insert-only, upsert, or retract behavior
    val table2 =
      tableEnv.fromChangelogStream(
        dataStream2,
        Schema.newBuilder().primaryKey("f0").build(),
        ChangelogMode.upsert())

    // register the table under a name and perform an aggregation
    tableEnv.createTemporaryView("InputTable2", table2)
    tableEnv
      .executeSql("SELECT f0 AS name, SUM(f1) AS score FROM InputTable2 GROUP BY f0")
      .print()

    // prints:
    // +----+--------------------------------+-------------+
    // | op |                           name |       score |
    // +----+--------------------------------+-------------+
    // | +I |                            Bob |           5 |
    // | +I |                          Alice |          12 |
    // | -U |                          Alice |          12 |
    // | +U |                          Alice |         100 |
    // +----+--------------------------------+-------------+

  }
}


case class User(name: String, score: java.lang.Integer, event_time: java.time.Instant)

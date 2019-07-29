name := "feature_fingerPrint_offline"

val spark_version = "2.4.2"

val sparkCore = "org.apache.spark" %% "spark-core" % spark_version
val sparkStreaming = "org.apache.spark" %% "spark-streaming" % spark_version
val sparkSql = "org.apache.spark" %% "spark-sql" % spark_version
val sparkHive = "org.apache.spark" %% "spark-hive" % spark_version
val sparkMLib = "org.apache.spark" %% "spark-mllib" % spark_version
val config = "com.typesafe" % "config" % "1.3.2"
val jpmml = "org.jpmml" % "jpmml-evaluator-spark" % "1.2.2"

val hadoopCommon= "org.apache.hadoop" % "hadoop-common" % "2.7.5"
val hdfs = "org.apache.hadoop" % "hadoop-hdfs" % "2.7.5"
val akka= "com.typesafe.akka" %% "akka-actor" % "2.5.10"
val akkaRemote= "com.typesafe.akka" %% "akka-remote" % "2.5.10"
val akkaContrib= "com.typesafe.akka" %% "akka-contrib" % "2.5.10"
val akkaTestKit= "com.typesafe.akka" %% "akka-testkit" % "2.5.10"
val akkaSLF4j= "com.typesafe.akka" %% "akka-slf4j" % "2.5.10"
val akkaQuartzScheduler = "com.enragedginger" %% "akka-quartz-scheduler" % "1.7.1-akka-2.5.x"
val commonMail= "org.apache.commons" % "commons-email" % "1.5"
val sunMail = "com.sun.mail" % "javax.mail" % "1.6.2"
val hiveJdbc = "org.apache.hive" % "hive-jdbc" % "2.1.1"
val play= "com.typesafe.play" %% "play" % "2.6.12"
val logbackClassic= "ch.qos.logback" % "logback-classic" % "1.2.2"
val logbackCore= "ch.qos.logback" % "logback-core" % "1.2.2"
val log4jtoSLF4j = "org.slf4j" % "log4j-over-slf4j" % "1.7.25"
val slf4j= "org.slf4j" % "slf4j-api" % "1.7.25"

val jackson = "com.fasterxml.jackson.module" %% "jackson-module-scala" % "2.6.5"
val elasticJob = "com.dangdang" % "elastic-job-lite-core" % "2.1.5"
val aeroClient = "com.aerospike" % "aerospike-client" % "4.2.3"

val mysqlConnector = "mysql" % "mysql-connector-java" % "5.1.47"
// https://mvnrepository.com/artifact/com.h2database/h2
val h2db = "com.h2database" % "h2" % "1.4.193"



lazy val sparkDependencies = Seq(sparkCore,sparkStreaming,sparkSql,sparkHive,sparkMLib)
lazy val commonDependencies = Seq(config)
lazy val commonSettings = Seq(
  version := "0.1",
  scalaVersion := "2.11.8"
)

lazy val assemblySettings = Seq(
  assemblyMergeStrategy in assembly := {
    case PathList("META-INF", xs @ _*) => MergeStrategy.discard
    case PathList("pom.properties", xs @ _*) => MergeStrategy.discard
    case PathList("pom.xml", xs @ _*) => MergeStrategy.discard
    case PathList("reference.conf") => MergeStrategy.concat
    case "*.conf" => MergeStrategy.concat
    case PathList(ps @ _*) => MergeStrategy.first
    case x =>
      val oldStrategy = (assemblyMergeStrategy in assembly).value
      oldStrategy(x)
  }
)

lazy val feature_offine = (project in file("./feature_offline"))
  .settings(Seq(
    version := "0.1",
    scalaVersion := "2.11.8"
  ):_*)
  .settings(
    name := "feature_offline",
    assemblySettings,
    assemblyOption in assembly := (assemblyOption in assembly).value.copy(includeScala = false),
    libraryDependencies ++= sparkDependencies.map(_ % "provided" withSources()),
    libraryDependencies  ++= Seq(config)
  ) //.dependsOn(sub2)

lazy val test_pmml = (project in file("./test_pmml"))
  .settings(commonSettings:_*)
  .settings(
    name := "test_pmml",
    assemblySettings,
    assemblyOption in assembly := (assemblyOption in assembly).value.copy(includeScala = false),
    libraryDependencies ++= sparkDependencies.map(_  withSources()),
    libraryDependencies  ++= Seq(jpmml)
  ) //.dependsOn(sub2)

lazy val free_marker=(project in file("./freemarker_component"))
  .settings(commonSettings:_*)
  .settings(
    name := "freemarker_component",
    libraryDependencies ++= Seq("org.freemarker" % "freemarker" % "2.3.28")
  )

lazy val test_hbase=(project in file("./test_hbase"))
  .settings(commonSettings:_*)
  .settings(
    name := "test_hbase",
    resolvers += "spark-hbase" at "http://repo.hortonworks.com/content/groups/public/",
    libraryDependencies ++= Seq(
      "com.hortonworks" % "shc-core" % "1.1.1-2.1-s_2.11"),
    libraryDependencies ++= Seq(
      "org.apache.hbase" % "hbase-client" % "2.1.4",
      "org.apache.hbase" % "hbase-common" % "2.1.4",
      "org.apache.hbase" % "hbase-server" % "2.1.4",
      "org.apache.hbase" % "hbase" % "2.1.4",
      "org.apache.hbase" % "hbase-mapreduce" % "2.1.4"
    ),
    assemblySettings,
    assemblyOption in assembly := (assemblyOption in assembly).value.copy(includeScala = true),
    libraryDependencies ++= sparkDependencies.map(_ % "provided" withSources()),
    libraryDependencies  ++= Seq(config,aeroClient)
  )

lazy val test_flink=(project in file("./test_flink"))
  .settings(commonSettings:_*)
  .settings(
    name := "test_flink",
    libraryDependencies ++= Seq(
      "org.apache.flink" % "flink-scala_2.11" % "1.8.0",
      "org.apache.flink" % "flink-streaming-scala_2.11" % "1.8.0"
    )
  )

lazy val feature_core = (project in file("./feature_core"))
  .settings(Seq(
    version := "0.1",
    scalaVersion := "2.11.8"
  ):_*)
  .settings(
    name := "feature_core",
    libraryDependencies ++= Seq(jackson,config,logbackClassic,logbackCore),
    libraryDependencies ++= Seq(hadoopCommon,hdfs).map(x=>x exclude("org.slf4j","slf4j-log4j12"))
  )

lazy val feature_pipeline = (project in file("./feature_pipeline"))
  .settings(Seq(
    version := "0.1",
    scalaVersion := "2.11.8"
  ):_*)
  .settings(
    name := "feature_pipeline",
    assemblySettings,
    assemblyOption in assembly := (assemblyOption in assembly).value.copy(includeScala = true),
    libraryDependencies ++= Seq(akkaQuartzScheduler,akka,akkaSLF4j,h2db)
  ).dependsOn(feature_core)

lazy val root = (project in file("."))
  .settings(commonSettings:_*)
  .settings(
    name := "root"
  ).aggregate(feature_offine,test_pmml,test_hbase,feature_pipeline)
//  .aggregate(feature_offine,feature_pipeline)
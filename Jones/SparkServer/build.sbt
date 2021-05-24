name := "adtl_spark_server"

version := "0.1"

scalaVersion := "2.12.5"
val akkaVersion = "2.6.8"
val akkaHttpVersion = "10.2.3"

resolvers += "Nexus" at "http://tool.adtd.sogou:8082/repository/maven-releases/"

libraryDependencies ++= Seq(
  // 要将spark lib下的 "commons-compress 换成这个新的版本, 否则导出超大excel时, 由于压缩缓存会产生报错
  "org.apache.commons" % "commons-compress" % "1.20",
  "com.rabbitmq" % "amqp-client" % "3.4.3",
  "com.thesamet.scalapb" %% "scalapb-runtime" % scalapb.compiler.Version.scalapbVersion % "protobuf",
  "com.alibaba" % "druid" % "1.1.20",
  "sogou.adtl" % "ahocorasick_2.11" % "1.0.0",
  "com.typesafe.akka" %% "akka-http" % akkaHttpVersion,
  "com.typesafe.akka" %% "akka-stream" % akkaVersion,
  "com.typesafe.akka" %% "akka-actor" % akkaVersion,
  "com.typesafe.akka" %% "akka-stream" % akkaVersion,
  "com.typesafe" % "config" % "1.2.1",
  "org.scalaj" %% "scalaj-http" % "2.3.0",
  "com.typesafe.play" %% "play-json" % "2.7.4",
  "org.joda" % "joda-convert" % "1.8.1",
  "joda-time" % "joda-time" % "2.9",
  "org.apache.spark" %% "spark-core" % "3.0.1", // % "provided",
  "org.apache.spark" %% "spark-sql" % "3.0.1", //% "provided",
  "commons-net" % "commons-net" % "3.7",
  "sogou.adtl" %% "spark-excel" % "1.1",
  "com.monitorjbl" % "xlsx-streamer" % "2.1.0",
  "sogou.adtl" % "clickhouse-tools_2.11" % "0.1.3",
  "redis.clients" % "jedis" % "3.2.0",
  "org.scala-lang" % "scala-compiler" % scalaVersion.value,
  "org.apache.kylin" % "kylin-jdbc" % "4.0.0-alpha",
  "mysql" % "mysql-connector-java" % "5.1.6"
)

dependencyOverrides ++= {
  Seq(
    "com.fasterxml.jackson.module" %% "jackson-module-scala" % "2.6.7.1",
    "com.fasterxml.jackson.core" % "jackson-databind" % "2.6.7",
    "com.fasterxml.jackson.core" % "jackson-core" % "2.6.7"
  )
}

PB.targets in Compile := Seq(
  scalapb.gen() -> (sourceManaged in Compile).value
)

assemblyMergeStrategy in assembly := {
  case PathList("META-INF", xs@_*) => MergeStrategy.discard
  case PathList("pom.properties", xs@_*) => MergeStrategy.discard
  case PathList("pom.xml", xs@_*) => MergeStrategy.discard
  case PathList("reference.conf") => MergeStrategy.concat
  case "*.conf" => MergeStrategy.concat
  case PathList(ps@_*) => MergeStrategy.first
  case x =>
    val oldStrategy = (assemblyMergeStrategy in assembly).value
    oldStrategy(x)
}

assemblyOption in assembly := (assemblyOption in assembly).value.copy(includeScala = true)
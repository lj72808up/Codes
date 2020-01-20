name := "AmqpClient"

version := "0.1"

scalaVersion := "2.12.10"


libraryDependencies ++= Seq("com.typesafe" % "config" % "1.3.2",
  "com.alibaba" % "druid" % "1.1.20",
  "mysql" % "mysql-connector-java" % "5.1.6",
  "com.rabbitmq" % "amqp-client" % "5.7.3",
  // https://mvnrepository.com/artifact/junit/junit
  "junit" % "junit" % "4.13" % Test,
  // https://mvnrepository.com/artifact/org.slf4j/slf4j-api
  "org.slf4j" % "slf4j-api" % "1.7.30"

)

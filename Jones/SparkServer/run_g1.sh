/search/odin/spark-3.0.1-bin-sogou/bin/spark-submit --master=yarn --class WebServer --deploy-mode=client \
--num-executors 10 \
--driver-memory 12G \
--executor-memory 12G \
--executor-cores 4 --conf \
spark.sql.shuffle.partitions=1000 --conf \
spark.dynamicAllocation.enabled=true --conf \
spark.dynamicAllocation.maxExecutors=140 --conf \
spark.shuffle.service.enabled=true --conf \
spark.sql.planner.skewJoin.threshold=true --conf \
spark.sql.autoBroadcastJoinThreshold=104857600 --conf \
spark.default.parallelism=1000  \
--conf "spark.driver.extraJavaOptions=-XX:+UseG1GC -XX:+PrintTenuringDistribution -XX:+PrintHeapAtGC -XX:+PrintReferenceGC -XX:+PrintGCApplicationStoppedTime -Xloggc:/root/test/spark3.0_server/gc/gc-%t.log -XX:+UseGCLogFileRotation -XX:NumberOfGCLogFiles=14 -XX:GCLogFileSize=100M -XX:NativeMemoryTracking=summary" \
--jars sogou-adtl-dw-tool-udf-spark3-assembly-0.1.jar adtl_spark_server-assembly-0.1_test.jar

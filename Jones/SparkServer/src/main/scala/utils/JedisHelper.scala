package utils

import java.io.InputStream

import adtl.platform.Dataproto.{Linedata, TranStream}
import models.config.{MainConf, RedisConf}
import redis.clients.jedis.{Jedis, JedisPool, JedisPoolConfig}

class JedisHelper extends PrintLogger {

  val pool: JedisPool = JedisHelper.pool

  def saveRecorderFromIs(in: InputStream, queryId: String): Unit = {
    val tran = XSSFUtil.printExcelByStream(in, 5000)
    val protobufBytes = tran.toByteArray
    val keyTimeOut = 7 * 24 * 3600
    this.putKey(queryId, protobufBytes, keyTimeOut)
  }


  /**
    * @param key   : queryId
    * @param value : TranStream
    */
  def putKey(key: String, value: Array[Byte], keyTimeOut: Int = 0): Unit = {
    var jedis: Jedis = null
    try {
      jedis = pool.getResource
      jedis.set(key.getBytes(), value)
      if (keyTimeOut != 0) {
        jedis.expire(key, keyTimeOut)
      }
    } finally {
      if (jedis != null) {
        jedis.close()
      }
    }
  }

  def getKey(key: String): Array[Byte] = {
    var jedis: Jedis = null
    try {
      jedis = pool.getResource
      jedis.get(key.getBytes)
    } finally {
      if (jedis != null) {
        jedis.close()
      }
    }
  }

  def getTrans(queryId: String): TranStream = {
    val bytes = this.getKey(queryId)
    if (bytes != null) {
      info(s"${queryId} receive data from redis finish")
      TranStream.parseFrom(bytes)
    } else {
      info(s"${queryId} receive empty data from redis !!")
//      TranStream().withLine(Array[Linedata]()).withSchema(Array[String]())
      null
    }
  }
}


object JedisHelper {
  var pool: JedisPool = _

  def initPool(): Unit = {
    val redisConf: RedisConf = MainConf("basicConf").redis
    val host: String = redisConf.host
    val port: Int = redisConf.port
    val auth: String = redisConf.passwd

    val timeout = 10000
    val jedisPool = new JedisPool(new JedisPoolConfig, host, port, timeout, auth)
    pool = jedisPool
  }

  def getSchemaKey(queryId: String): String = {
    s"${queryId}_schema"
  }

  def close(): Unit = {
    pool.close()
  }
}
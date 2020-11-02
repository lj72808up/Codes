import redis.clients.jedis.{Jedis, JedisPool, JedisPoolConfig}

object TestJedis {
  def main(args: Array[String]): Unit = {
    // (1) Jedis 实例不是线程安全的, 不能再线程中共享; 如果想共享, 可以共享 JedisPool  (依赖 Commons Pool 2 实现)
    val host = "test.astar.ms.redis.sogou"
    val port = 1927
    val auth = "noSafeNoWork2020"
    val timeout = 10000
    val pool = new JedisPool(new JedisPoolConfig, host, port, timeout, auth)
    var jedis:Jedis = null
    try {
      jedis = pool.getResource

      val key = "test-1027"
      //      jedis.set(key,"xxx");
      //      println(jedis.get(key))
      for (i <- 0 to 10) {
        jedis.rpush(key,  s"$i")
      }
      val res = jedis.lrange(key, 0, -1)
      // expire "test-1027" 30  超时时间30秒
      // ttl "test-1027"

      val keyTimeOut = 7 * 24 * 3600
      jedis.expire(key, keyTimeOut)
      println(res)
      println(jedis.ttl(key))


    } finally {
      if (jedis != null) {
        jedis.close();
      }
    }

    pool.close()
  }

  def testClassLoader(args: Array[String]): Unit = {
    val a = TestJedis.getClass.getClassLoader
    val b = classOf[Jedis].getClassLoader

    println(a)
    println(b)
    println(Thread.currentThread().getContextClassLoader)
  }
}

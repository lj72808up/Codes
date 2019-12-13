package com.sogou.datasource

import java.util

import com.alibaba.druid.filter.Filter
import com.alibaba.druid.filter.stat.StatFilter
import com.alibaba.druid.pool.DruidDataSource
import javax.sql.DataSource

class DateSourceBuilder {
  def build(): DataSource = {

    val dataSource = new DruidDataSource()
    val isEnableMonitor = true
    dataSource.setName("testds")
    dataSource.setUrl("jdbc:mysql://10.139.36.81:3306/adtl_test")
    dataSource.setUsername("yuhancheng")
    dataSource.setPassword("yuhancheng")
    //    连接池加密
    //    dataSource.setConnectionProperties(
    //      "config.decrypt=trueconfig.decrypt.key=" + druidDataSourceProperties.getPwdPublicKey())
    dataSource.setFilters("config")
    dataSource.setMaxActive(100)
    dataSource.setInitialSize(5)
    dataSource.setMaxWait(50 * 1000)
    dataSource.setMinIdle(5)
    //    dataSource.setTimeBetweenEvictionRunsMillis(druidDataSourceProperties.getTimeBetweenEvictionRunsMillis())
    //    dataSource.setMinEvictableIdleTimeMillis(druidDataSourceProperties.getMinEvictableIdleTimeMillis())
    //    dataSource.setValidationQuery(druidDataSourceProperties.getValidationQuery())
    //    dataSource.setTestWhileIdle(druidDataSourceProperties.isTestWhileIdle())
    //    dataSource.setTestOnBorrow(druidDataSourceProperties.isTestOnBorrow())
    //    dataSource.setTestOnReturn(druidDataSourceProperties.isTestOnReturn())
    //    dataSource.setPoolPreparedStatements(druidDataSourceProperties.isPoolPreparedStatements())
    //    dataSource.setMaxOpenPreparedStatements(druidDataSourceProperties.getMaxOpenPreparedStatements())
    val filters = new util.ArrayList[Filter]()
    if (isEnableMonitor) {
      val statFilter = new StatFilter()
      statFilter.setLogSlowSql(true)  // 记录慢查询
      statFilter.setMergeSql(true) // 监控显示合并sql
      statFilter.setSlowSqlMillis(2000)
      
      filters.add(statFilter)
      dataSource.setProxyFilters(filters)
    }
    dataSource
  }

}

'use strict'
const merge = require('webpack-merge')
const prodEnv = require('./prod.env')

module.exports = merge(prodEnv, {
  NODE_ENV: '"production"',
  BASE_API_URL: '"http://pre.astar.adtech.sogou/jones"'
})

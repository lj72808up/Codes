<template>
  <div id="app">
    <router-view v-if="validatePass"/>
  </div>
</template>

<script>
import {delCookie, getCookie, setCookie} from './conf/utils'

const axios = require('axios')
const RSA = require('node-rsa')

function getQueryString(name) {
  let reg = `(^|&)${name}=([^&]*)(&|$)`
  let r = window.location.search.substr(1).match(reg)
  if (r != null) return unescape(r[2])
  return null
}

export default {
  name: 'App',
  data() {
    return {
      validatePass: false
    }
  },
  methods: {
    isNotExpire(ts1, ts2) {
      let interval = 2 * 60 * 1000
      return -interval <= ts1 - ts2 <= interval
    },
    async parseToken(pToken) {
      // todo 解析ptoken改为axios访问接口解析astar接口
      console.log(pToken)
      // alert(pToken)
      let param = {
        'pToken': pToken
      }
      let uid = -1
      let userName = ''
      let now = new Date().getTime()

      let url = axios.defaults.baseURL.replace(/jones/g, '') + 'xiaop/jone/parseToken'
      let tokenObj = {
        'uid': -99999,
        'name': "dummyUser",
        'ts': -99999
      }

      await axios.post(url, param, {responseType: 'json'})
        .then(function (response) {
          let data = response.data
          if (data.Res) {  // 验证通过
            let returnObj = JSON.parse(data.Info)
            uid = returnObj.uid
            userName = returnObj.name
          } else {
            // alert('token解析失败, 请重新登录')
          }
          tokenObj = {
            'uid': uid,
            'name': userName,
            'ts': now
          }
        })
        .catch(function (error) {
          console.log(error)
          // alert(error)
          tokenObj = {
            'uid': -1,
            'name': '',
            'ts': now
          }
        })
      return tokenObj
    },
    makeJumpUrl() {
      let jumpURL = 'http://datacenter.adtech.sogou/'
      let curURL = process.env.BASE_API_URL // window.location.href
      if ((curURL + '').search('pre') != -1) {
        jumpURL = 'http://pre.datacenter.adtech.sogou/'
      }
      let routePath = decodeURI(getCookie('_adtech_path'))
      if (routePath != null && routePath != '' && routePath != undefined && routePath != 'null') {
        jumpURL = jumpURL + '#' + routePath
      }
      return jumpURL
    },
    getUserFromCookieAndValidDate(now) {
      // 校验用户名是否存在和是否过期,
      let lastLoginTime = getCookie('_adtech_last_ts')
      let userId = getCookie('_adtech_user')
      if (userId !== null && userId !== undefined && userId !== "" && userId != "-1") {
        // alert ("cookie中存在 userId:"+userId)
        if (this.isNotExpire(lastLoginTime, now)) {
          // alert ("userId没有过期")
          return userId
        }
      }
      return null
    },
    postUserMapping(uid, username) {
      let param = {
        "uid": uid,
        "username": username,
      }
      axios.post('/auth/user/postUser', param, {responseType: 'json'})
        .then(function (response) {
        }).catch(function (error) {
        console.error(error)
      })
    },
    async checkLogin() {
      let jumpURL = this.makeJumpUrl()
      let now = new Date().getTime()

      let cookieUserId = this.getUserFromCookieAndValidDate(now)
      // alert("最终, cookie中俄的userId:"+cookieUserId)

      if (cookieUserId === null || cookieUserId === undefined || cookieUserId === "" || cookieUserId == -1) {
        let pToken = getQueryString('ptoken')
        if (pToken === null || pToken === '' || pToken === undefined) {
          // alert("cache中没有userId, URL中也没有Token, 跳转登录页面")
          window.location.href = 'https://login.sogou-inc.com/?appid=1812&sso_redirect=' + jumpURL
        } else {
          // alert('URL参数中存在token: ' + pToken)
          let tokenObj = await this.parseToken(pToken)
          console.log("解析pToken后")
          // alert("解析URL参数中的token" + JSON.stringify(tokenObj))
          let userId = tokenObj.uid
          let userName = tokenObj.name
          let ts = tokenObj.ts
          // alert('解析token后, userId:' + userId + 'userName:' + userName + 'ts:' + ts)

          setCookie('_adtech_user', userId, 300)
          setCookie('_adtech_user_name', encodeURI(userName), 300)
          setCookie('_adtech_last_ts', ts, 1)

          let param = {
            "uid": userId,
            "username": userName,
          }

          let url = axios.defaults.baseURL.replace(/jones/g, '') + 'xiaop/auth/user/postUser'
          axios.post(url, param, {responseType: 'json'})
            .then(function (response) {
              window.location.href = jumpURL // 'http://datacenter.adtech.sogou/#/'
            }).catch(function (error) {
            console.error(error)
          })
        }
      } else {
        // alert("cookie token 没过期, 正常放行, 跳转到:" + jumpURL)
        this.validatePass = true
        this.refreshCache()
      }
    },
    refreshCache() {
      /** 每次进入系统,刷新一下cache*/
      let now = new Date().getTime()
      setCookie('_adtech_last_ts', now, 1)
    }
  },
  component: {},
  activated: function () {
    this.getCase()
  },
  /*created: function () {
    axios.get('/auth/routes')
  },*/
  created() {
    this.checkLogin()
    // setCookie('_adtech_user', "liujie02", 1)
    // this.validatePass = true
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

.el-dialog {
  display: flex;
  flex-direction: column;
  margin: 0 !important;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}
</style>



//获取cookie、
export function getCookie (name) {
  var arr, reg = new RegExp('(^| )' + name + '=([^;]*)(;|$)')
  if (arr = document.cookie.match(reg)) {
    return (arr[2])
  } else {
    return null
  }
}

//设置cookie,增加到vue实例方便全局调用
export function setCookie (c_name, value, expiredays) {
  var exdate = new Date()
  exdate.setDate(exdate.getDate() + expiredays)
  document.cookie = c_name + '=' + value + ((expiredays == null) ? '' : ';expires=' + exdate.toGMTString())
}

//设置根域名cookie,
export function setRootCookie (c_name, value, expiredays) {
  var exdate = new Date()
  exdate.setDate(exdate.getDate() + expiredays)
  document.cookie = c_name + '=' + value + ((expiredays == null) ? '' : ';expires=' + exdate.toGMTString()) + ';Path=/;domain=adtech.sogou;'
}

//删除cookie
export function delCookie (name) {
  var exp = new Date()
  exp.setTime(exp.getTime() - 1)
  var cval = getCookie(name)
  if (cval != null) {
    document.cookie = name + '=' + cval + ';expires=' + exp.toGMTString()
  }
}

export function makeLvpiURL(lvpiURL) {
  let completeUrl = axios.defaults.baseURL.replace(/jones/g, '') + 'lvpi' + lvpiURL // lvpiURL: 'xiaop/auth/user/postUser'
  return completeUrl
}

// format date yyyyMMdd
export function formatDate (date) {
  var year = date.getFullYear().toString()
  var month = (date.getMonth() + 101).toString().substring(1)
  var day = (date.getDate() + 100).toString().substring(1)
  return year + month + day
}

// format date yyyy-MM-dd
export function formatDateLine (date) {
  var year = date.getFullYear().toString()
  var month = (date.getMonth() + 101).toString().substring(1)
  var day = (date.getDate() + 100).toString().substring(1)
  return year + '-' + month + '-' + day
}

//file to json
const axios = require('axios')

export function getFileContent (file) {
  var formData = new FormData()
  formData.append('file', file)
  axios.post('task/file', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  }).then(function (response) {
    return response.data
  }).catch(function (error) {
    return null
  })
}

export function isEmpty (v) {
  switch (typeof v) {
    case 'undefined':
      return true
    case 'string':
      if (v.replace(/(^[ \t\n\r]*)|([ \t\n\r]*$)/g, '').length === 0) return true
      break
    case 'boolean':
      if (!v) return true
      break
    case 'number':
      if (0 === v || isNaN(v)) return true
      break
    case 'object':
      if (null === v || v.length === 0) return true
  }
  return false
}

// 判断任务状态
export function isRunFinsh (statusCode) {
  return statusCode === 3 || statusCode === 6 || statusCode === 12
}

export function isRunFinshPagani (statusCode) {
  return statusCode === 22
}

// 判断任务类型
export function type2route (type) {
  if (type === '无线_query') {
    return 'wap'
  } else if (type === 'PC_query') {
    return 'pc'
  } else if (type === '无线搜索查询') {
    return 'wapTotal'
  } else if (type === 'PC搜索查询') {
    return 'pcTotal'
  } else if (type === 'yh') {
    return 'gala'
  } else if (type === 'sql查询') {
    return 'sql'
  } else if (type === '工单') {
    return 'jdbc'
  } else {
    return false
  }
}

// is user himself
export function isUserHimself (user) {
  return user === getCookie('_adtech_user')
}

// 解析URL 获取指定param的值
export function getQueryVariable (query, variable) {
  var vars = query.split('&')
  for (var i = 0; i < vars.length; i++) {
    var pair = vars[i].split('=')
    if (pair[0] == variable) {
      return pair[1]
    }
  }
  return (false)
}

// 判断是否包含中文
export function checkChinese (content) {
  return /.*[\u4e00-\u9fa5]+.*/.test(content);
}


export function refreshCache(){
  /** 每次进入系统,刷新一下cache*/
  let cookieToken = getCookie('_adtech_local_token')
  if (!(cookieToken === null || cookieToken === '' || cookieToken === undefined)) {
    setCookie('_adtech_local_token', cookieToken, 1)
  }

  let now = new Date().getTime()
  setCookie('_adtech_last_ts', now, 1)
}

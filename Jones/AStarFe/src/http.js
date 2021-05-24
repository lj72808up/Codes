/*
 * http axios conf
 */

import axios from 'axios'
import router from './router'
import {getCookie,refreshCache} from './conf/utils'
import { Message, Loading } from 'element-ui'

// axios全局配置
// online server
// axios.defaults.baseURL = 'http://10.134.101.120:8080/jones'; // old
// axios.defaults.baseURL = 'http://10.140.38.140:8080/jones';  // new
// test server
// axios.defaults.baseURL = 'http://astar.adtech.sogou/jones';
axios.defaults.baseURL = process.env.BASE_API_URL;
// axios.defaults.baseURL = 'http://localhost:9090/jones'


axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';

// http request 拦截器
axios.interceptors.request.use(
  (config) => {
    var username = getCookie('_adtech_user');
    var password = getCookie('_adtech_local_token');
    refreshCache()
    if(username) {
      config.headers.Authorization = 'Basic ' + btoa(username +':' + password);
    }
    return config;
  },
  (err) => {
    return Promise.reject(err);
  }
);

// http response 拦截器
axios.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if(error.response) {
      switch (error.response.status) {
        case 403:
          router.replace({
            path: '/403',
            query: {redirect: router.currentRoute.fullPath}
          });
        // case 404:
        //   router.replace({
        //     path: '/404',
        //     query: {redirect: router.currentRoute.fullPath}
        //   })
      }
    }
    return Promise.reject(error.response.data);
  }
);
export default axios;

<template>
  <el-container>
    <el-aside v-show="!isMobile" width="auto">
      <el-scrollbar style="height: 100%;">
        <el-row>
          <div style="float: right; margin-right: 25px; margin-top: 20px; cursor: move" id="drag-div">
            <i class="el-icon-s-unfold" style="font-size: 25px;color: white"
               @click="toggleCollapse()"/>
          </div>
        </el-row>

        <!--class="el-menu-vertical-demo"-->
        <el-menu :default-active="activeIndex"
                 :router=false
                 :collapse="isCollapse"
                 :collapse-transition="true"
                 class="el-menu-vertical-demo"
                 background-color="#304156"
                 text-color="#fff"
                 active-text-color="#ffd04b"
                 :unique-opened="true"
                 @select="onSelect"
                 id="myLeft">
          <template v-for="(item,idx) in urls">
            <el-submenu :index="item.Index" v-if="item.Children != null && item.Children.length>0 ">
              <template slot="title">
                <i :class="item.Class"/>
                <span slot="title">{{item.Title}}</span>
              </template>
              <el-menu-item-group>
                <template v-for="ele in item.Children">
                  <el-menu-item :index="ele.Index">{{ele.Title}}</el-menu-item>
                </template>
              </el-menu-item-group>
            </el-submenu>

            <el-menu-item :index="item.Index" v-if="item.Children === null || item.Children.length===0 ">
              <i :class="item.Class"/>
              <span>{{item.Title}}</span>
            </el-menu-item>
          </template>
        </el-menu>
      </el-scrollbar>

    </el-aside>


    <el-container>
      <!--      <div class="sideDiv">-->
      <el-header v-show="!isMobile" height="60px">
        <!--        <div style="display:inline; float: left; margin-left: -10px; ">-->
        <!--          <i class="el-icon-caret-right" style="font-size:25px; color: red"/>-->
        <!--        </div>-->
        <!--<div style="float: left; margin-right: 10px; margin-top: 2px; display: inline-block; cursor: move"
             id="drag-div" >
          <i class="el-icon-s-unfold" style="font-size: 25px;color: #304156"/>
        </div>-->
        <div class="name">
          <img src="../assets/logo.png" height="40px" align="left">
        </div>

        <div class="logoutClass">
          <el-dropdown class="avatar-container" trigger="click">
            <div class="avatar-wrapper">
              <el-image class="avatar-inner"
                        style="width: 35px; height: 35px; border-radius: 50%; "
                        :src="userImg"
                        fit="full">
              </el-image>
              <span class="avatar-inner">{{echoUser}}</span>
              <i class="fa fa-sort-desc avatar-inner"/>
            </div>
            <el-dropdown-menu slot="dropdown" class="user-dropdown">
              <router-link to="/doc">
                <el-dropdown-item>
                  帮助
                </el-dropdown-item>
              </router-link>

              <!--              <a href="javascript:loginOut()">-->
              <span @click="loginOut">
                <el-dropdown-item>登出</el-dropdown-item>
              </span>
              <!--              </a>-->

            </el-dropdown-menu>
          </el-dropdown>
        </div>
      </el-header>

      <el-main>
        <div v-if="this.$route.path==='/'" style="margin: -18px">

          <!--  :height="imgHeight+'px'"  -->
          <!--          <img src="static/img/homePage1.png" width="100%">-->
          <el-carousel indicator-position="outside" height="580px">
            <el-carousel-item v-for="item in imgs" :key="item">
              <img :src="item" width="100%" ref="image"/>
            </el-carousel-item>
          </el-carousel>
          <!--文字卡片1-->
          <el-row>
            <el-col :span="10" :offset="3">
              <h2>数据查询</h2>
              <span>支持 spark, mysql, clickhouse 等查询引擎 </span> <br/>
              <span>与元数据联通, 输入 sql 时进行全方位的提示</span><br/>
              <span>支持自定义参数配置</span>
            </el-col>
            <el-col :span="11">
              <img src="static/img/fine1.jpeg" width="60%"/>
            </el-col>
          </el-row>
          <el-divider/>
          <!--文字卡片2-->
          <el-row>
            <el-col :span="10" :offset="3">
              <img src="static/img/fine2.jpeg" width="60%"/>
            </el-col>
            <el-col :span="11">
              <h2>数据图表</h2>
              <span>丰富的数据可视化操作, 让数据直观, 明确的展现给分析人员 </span> <br/>
              <span>支持结果对比高亮, 自定义数据更新周期等</span><br/>
            </el-col>
          </el-row>
          <el-divider/>
          <!--文字卡片3-->
          <el-row>
            <el-col :span="10" :offset="3">
              <h2>自动数仓</h2>
              <span>一站式自建数仓方便用户创建自己的分析流程</span><br/>
              <span>拖拽式工作流让创建更简单, 高效</span><br/>
              <span>支持自定义邮件发送结果</span><br/>
            </el-col>
            <el-col :span="11">
              <img src="static/img/fine3.jpeg" width="60%"/>
            </el-col>
          </el-row>
          <el-divider/>
          <!--尾部介绍-->
          <el-row style="background-color:#1b1c1f; ">
            <div style="margin-top: 30px; background-color:#1b1c1f; ">
              <div style="background-color:#1b1c1f;">
                <el-col :span="4" :offset="8" style="background-color:#1b1c1f; ">
                  <el-row align="middle" style="background-color:#1b1c1f; ">
                    <span style="color: white">关于我们</span>
                  </el-row>
                  <ul style="list-style:none; margin-left: -50px; color: white; font-size: 11px">
                    <li>商业数据产品</li>
                    <li>商业数据开发</li>
                  </ul>

                </el-col>
                <el-col :span="6" style="background-color:#1b1c1f; ">
                  <el-row style="background-color:#1b1c1f; ">
                    <span style="color: white">联系方式</span>
                  </el-row>
                  <img src="../../static/img/xiaopScan.jpeg" width="20%" style="margin-left: -5px"/>
                </el-col>
              </div>
              <div style="background-color:#1b1c1f;">
                <el-col :offset="10" :span="8" style="background-color:#1b1c1f; ">
                  <span style="font-size: 11px; color: #555">ADTECH-数据平台</span>
                </el-col>
              </div>
            </div>
          </el-row>
        </div>
        <router-view v-else/>
      </el-main>
    </el-container>
  </el-container>
</template>

<style>

  .box-card {
    width: 100%;
  }

  .clearfix:before,
  .clearfix:after {
    display: table;
    content: "";
  }

  .clearfix:after {
    clear: both
  }

  .el-carousel__item h3 {
    color: #475669;
    font-size: 18px;
    opacity: 0.75;
    line-height: 300px;
    margin: 0;
  }

  .el-header {
    /*background-color: #545c64;*/
    color: #333;
    line-height: 50px;
    border: 1px solid #eee;
    padding: 0 0;
  }

  .el-aside {
    color: #545c64;
    background-color: #304156;
  }

  .el-main {
    margin-top: 0 !important;
    /*margin-left: 0;*/
    padding: 5px 20px 5px 20px !important;
  }

  .logoutClass {
    float: right;
    height: 100%;
    line-height: 50px;
  }

  .user-avatar {
    cursor: pointer;
    width: 40px;
    height: 40px;
    border-radius: 10px;
  }

  /* 菜单栏(非伸缩) */
  .el-menu-vertical-demo:not(.el-menu--collapse) {
    width: 200px;
    border-right: 0px;
    /*min-height: 400px;*/
  }

  /* 菜单栏(伸缩) */
  .el-menu--collapse {
    border-right: 0px;
    /*min-height: 400px;*/
  }

  .el-menu {
    border-right: 0px;
  }

  .avatar-inner {
    margin: 5px;
  }

  html, body, #app, .el-container {
    padding: 0px;
    margin: 0px;
    height: 100%;
  }

  .avatar-wrapper {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .menu-expanded {
    flex: 0 0 180px;
    width: 180px;
  }

  .el-menu {
    width: 100%;
  }

  .el-submenu .el-menu-item {
    min-width: 0px;
  }

  .el-menu--vertical > .el-menu--popup {
    max-height: 100vh;
    overflow-x: auto;
    overflow-y: auto;
    border: 1px solid;
  }

</style>

<script>
  import Home from './Layout.vue'
  import SqlQuery from './SqlQuery'
  import SecondComponent from './Queue'
  import NewJohn from './WapJohn'
  import {getCookie, setRootCookie, delCookie, setCookie} from '../conf/utils'
  import axios from 'axios'

  export default {
    name: '',
    data () {
      return {
        activeIndex: 'doc',
        isLimitWorkSheet: false,
        urls: [],
        isCollapse: false,
        userName: '',
        echoUser: '',
        userImg: 'http://puboa.sogou-inc.com/moa/sylla/mapi/img?&s=1&w=128&id=' + getCookie('_adtech_user'),
        cacheMap: {},
        timer: null,
        imgHeight: '',
        imgs: ['static/img/homePage1.png', 'static/img/homePage2.png']//['static/img/num0.jpeg', 'static/img/num1.jpeg']
      }
    },
    methods: {
      cancelDrag () {
        let svgResize = document.getElementById('drag-div')
        svgResize.onmousedown = null
      },
      drag () {
        let svgResize = document.getElementById('drag-div')
        svgResize.onmousedown = function (e) {
          let startX = e.clientX
          // svgResize.top = svgResize.offsetTop
          let curWidth = document.getElementById('myLeft').offsetWidth
          document.onmousemove = function (e) {
            let endX = e.clientX
            let moveLen = (endX - startX)
            /*let moveLen = svgResize.top + (endY - startY)
            if (moveLen < 30) moveLen = 30
            if (moveLen > maxT - 30) moveLen = maxT - 30*/
            // svgResize.style.top = moveLen
            // console.log(moveLen)
            // console.log('之前:' + document.getElementById('myLeft').offsetWidth)
            if (curWidth + moveLen > 200 && curWidth + moveLen < 500) {   // 最窄这么宽
              document.getElementById('myLeft').style.width = curWidth + moveLen + 'px'
            }
            // console.log('之后:' + document.getElementById('myLeft').offsetWidth)
          }
          document.onmouseup = function (evt) {
            document.onmousemove = null
            document.onmouseup = null
            svgResize.releaseCapture && svgResize.releaseCapture()
          }
          svgResize.setCapture && svgResize.setCapture()
          return false
        }
      },
      _isMobile () {
        let flag = navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i)
        return flag
      },

      setTestTime () {
        let self = this
        axios.get('/redash/timestamp')
          .then(function (response) {
            var res = []
            self.timestamp = response.data
            self.sec186 = (self.timestamp + 7482952) * 3
            setRootCookie('sec186', self.sec186.toString(), 1)

          }).catch(function (error) {
        })
        console.log('set timer sec186')
      },

      onSelect (index, indexPath) {
        console.log(index + ':::' + indexPath)
        if (index.startsWith('redashshow')) {
          console.log('是报表')

          let indexKey = index.split('/')[1]
          let cacheStr = this.cacheMap[indexKey]
          sessionStorage.setItem(indexKey, cacheStr)

          //this.$router.push({path: 'redashinsert'})
          //this.$router.push({path: 'redashshow'})
          //this.$router.replace({path: index})
          //this.$router.push({path: index})
          console.log(indexKey)
          this.$router.push({name: 'redashshow', params: {redashname: indexKey}})

        } else if (index.startsWith('dashboardshow')) {
          console.log('是dashboard')

          let indexKey = index.split('/')[1]
          let cacheStr = this.cacheMap[indexKey]
          sessionStorage.setItem(indexKey, cacheStr)
          console.log(indexKey)
          this.$router.push({name: 'dashboardshow', params: {dashboardname: indexKey}})

        } else if (index.startsWith('redashqueryshow')) {
          console.log('--------------------jump to option web----------')
          //this.$router.push({name: 'redashqueryoption'})
        } else {
          this.$router.push({name: index})    // 非报表
        }
      },
      toggleCollapse () {
        this.isCollapse = !this.isCollapse
        if (this.isCollapse) {
          this.cancelDrag()
        } else {
          this.drag()
        }
      },
      loginOut () {
        delCookie('SG_SSO')
        this.$confirm('您确定要退出吗?', '退出管理平台', {
          confirmButtonText: '确定',
          cancelButtonText: '取消'
        }).then(() => {
          delCookie('_adtech_local_token')
          delCookie('_adtech_user')
          let jumpURL = 'http://datacenter.adtech.sogou/'
          let curURL = process.env.BASE_API_URL // window.location.href
          if ((curURL + '').search('pre') != -1) {
            jumpURL = 'http://pre.datacenter.adtech.sogou/'
          }
          window.location.href = 'https://login.sogou-inc.com/logout.jsp?appid=1812&sso_redirect=' + jumpURL
        }).catch(function (error) {
          console.log(error)
        })
      }
    },
    component: {
      Home, SqlQuery, SecondComponent, NewJohn
    },
    computed: {
      isMobile: function () {
        if (this._isMobile()) {
          return true
        } else {
          return false
        }
      }
    },
    mounted () {
      console.log(this.$route.path)
      let user = decodeURI(getCookie('_adtech_user_name'))
      if (user === null || user === undefined) {
        user = getCookie('_adtech_user')
      }
      this.echoUser = user
      this.userName = getCookie('_adtech_user')
      this.userImg = 'http://puboa.sogou-inc.com/moa/sylla/mapi/img?&s=1&w=128&id=' + this.userName
      let username = this.userName
      let self = this

      this.timer = setInterval(this.setTestTime, 180000)

      function parseInfo (info) {
        // console.log(info)
        let infoObj = JSON.parse(info)
        let infoMap = new Map(Object.entries(infoObj))
        for (let [k, v] of infoMap) {
          self.cacheMap[k] = v
        }
      }

      axios.get('/auth/labels/checked/' + username)
        .then(function (response) {
          self.urls = response.data
          for (let i = 0, len = response.data.length; i < len; i++) {
            let firstLabel = response.data[i]
            // 首层
            if (firstLabel.Info !== null && firstLabel.Info !== '') {
              parseInfo(firstLabel.Info)
            }
            // 二层
            if (firstLabel.Children !== null && firstLabel.Children.length !== 0) {
              for (let j = 0, len2 = firstLabel.Children.length; j < len2; j++) {
                let secondLabel = firstLabel.Children[j]
                if (secondLabel.Info !== null && secondLabel.Info !== '') {
                  parseInfo(secondLabel.Info)
                }
              }
            }
          }
          // console.log(self.cacheMap)
        })
        .catch(function (error) {
          console.log(error)
        })
      this.drag()
    },
    beforeDestroy () {
      clearInterval(this.timer)
    }
  }
</script>

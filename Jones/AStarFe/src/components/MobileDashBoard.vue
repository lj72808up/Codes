<template>
  <div>

    <!--<mt-navbar v-model="selected">
      <mt-tab-item id="1">外卖</mt-tab-item>
      <mt-tab-item id="2">订单</mt-tab-item>
      <mt-tab-item id="3">发现</mt-tab-item>
    </mt-navbar>-->


    <div v-if="!showChildOptions">
      <!--显示iframe的地方 -->
      <div v-bind:style="{'height':myHeighgt+'px'}" class="outDiv">
        <iframe id="show-iframe" :src='url_src'
                frameborder="no"
                marginheight="0" marginwidth="0"
                scrolling="yes"/>

        <h3 style="display:none">
          This is redashshow
        </h3>
      </div>
    </div>

    <div v-if="showChildOptions">
      <mt-tab-container v-model="selected">
        <mt-tab-container-item id="1">
          <mt-cell v-for="opt in childOptions"
                   :title="opt.Label"
                   @click.native="showIframe(opt)"
                   is-link
                   value=""/>
        </mt-tab-container-item>
      </mt-tab-container>
    </div>

    <h3 style="margin-bottom: 50px;display:none">
      This is redashshow
    </h3>

    <!--固定到底部-->
    <mt-tabbar v-model="active">
      <mt-tab-item :id="bar.Mid" @click.native="onChooseContainer(bar)" v-for="bar in firstBars">
        <i :class="bar.Img" style="display: block; font-size: 22px; margin-bottom: 5px"/>
        {{bar.Label}}
      </mt-tab-item>
    </mt-tabbar>
  </div>
</template>

<script>
  import {getCookie, setRootCookie} from '../conf/utils'

  const axios = require('axios')

  export default {
    data() {
      return {
        active: 1,
        selected: '1',
        showChildOptions: true,
        url_src: '',
        timer: null,
        myHeighgt: '',
        timestamp: 0,
        sec186: 0,
        firstBars: [],
        childOptions: [],
      }
    },
    methods: {
      onChooseContainer(bar) {
        console.log(bar)
        let url = bar.Url
        if (url != '') {
          this.url_src = url
          this.showChildOptions = false
        } else {
          this.url_src = ''   // 有子菜单, 显示列表页
          this.showChildOptions = true
          // 查找子菜单
          let mid = bar.Mid
          let self = this
          axios.get('/bodong/getChildOptions/' + mid)
            .then(function (response) {
              self.childOptions = response.data
            })
            .catch(function (error) {
              console.error(error)
            })
        }
      },
      getFirstLevel() {
        let self = this
        axios.get('/bodong/getFirst')
          .then(function (response) {
            self.firstBars = response.data
            self.url_src = self.firstBars[0].Url
          })
          .catch(function (error) {
            console.error(error)
          })
      },
      showIframe(bar) {
        console.log(bar.Url)
        this.url_src = bar.Url
        this.showChildOptions = false
      },
      setTestTime() {
        let self = this
        /*        axios.get('/redash/timestamp')
                  .then(function (response) {
                    self.timestamp = response.data
                    self.sec186 = (self.timestamp + 7482952) * 3
                    setRootCookie('sec186', self.sec186.toString(), 1)

                  }).catch(function (error) {
                })
                console.log('set timer sec186')*/

        axios.get('/redash/timestamp')
          .then(function (response) {
            var res = []
            self.timestamp = response.data
            self.sec186 = (self.timestamp + 7482952) * 3
            // 打印sec186
            // alert(self.sec186.toString())
            setRootCookie('sec186', self.sec186.toString(), 1)

            self.isShowFrame = true

          }).catch(function (error) {
          console.log(error)
        })
      },
    },
    created() {
      this.setTestTime()
      this.timer = setInterval(this.setTestTime, 180000)
    },
    mounted() {
      let uid = getCookie('_adtech_user')
      let username = getCookie('_adtech_user_name')

      console.log('uid:' + uid + ',username:' + username + ',')

      let height = document.body.clientHeight
      console.log(height)
      this.myHeighgt = height - 15
      this.getFirstLevel()
      this.showChildOptions = false
    },
    beforeDestroy() {
      clearInterval(this.timer)
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  #show-iframe {
    width: 100%;
    height: 100% !important;
    background-color: transparent;
    overflow-x: scroll;
    overflow-y: hidden;
  }

  .outDiv {
    margin-left: -20px;
    margin-right: -20px;
    margin-top: -5px;
  }
</style>

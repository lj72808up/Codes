<template>
  <div v-bind:style="{'height':myHeighgt+'px'}" class="outDiv">
    <iframe id="show-iframe" :src='url_src'
            frameborder="no"
            marginheight="0" marginwidth="0"
            scrolling="yes"/>

    <!--    style="display:none">-->
    <h3 style="display:none">
      This is redashshow/{{$route.params.redashname}}
    </h3>
  </div>

</template>

<script>
  import {delCookie, getCookie, setCookie, setRootCookie} from '../conf/utils'

  const axios = require('axios')
  export default {
    name: 'MobileDashBoardShow',
    data () {
      return {
        url_src: '',
        myHeighgt: ''
      }
    },
    methods: {
      refreshTimeStamp () {
        // 访问一下timestamp
        axios.get('/redash/timestamp')
          .then(function (response) {
            let timestamp = response.data
            let sec186 = (timestamp + 7482952) * 3
            setRootCookie('sec186', sec186.toString(), 1)

          }).catch(function (error) {
          console.log(error)
        })
      },

      getRedashUrl (redashName) {
        let self = this
        axios.get('/mobileDashboard/getDashboardUrl/' + redashName)
          .then(function (resp) {
            let redashUrl = resp.data
            self.url_src = redashUrl
          })
          .catch(function (error) {
            console.log(error)
          })
      }
    },
    created () {
      let routePath = this.$route.path
      console.log(routePath)
      setCookie('_adtech_path', encodeURI(routePath), 1)
    },
    beforeDestroy () {
      delCookie('_adtech_path')
    },
    mounted () {
      let uid = getCookie('_adtech_user')
      let username = getCookie('_adtech_user_name')

      console.log('uid:' + uid + ',username:' + username + ',')

      let height = document.body.clientHeight
      console.log(height)
      this.myHeighgt = height - 15

      this.refreshTimeStamp()

      let redashName = this.$route.params.redashname
      this.getRedashUrl(redashName)

      console.log(this.$route)
    },
  }

</script>

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

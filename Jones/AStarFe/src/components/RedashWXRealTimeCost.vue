<template>
  <div style="height:100%; margin-left: -20px;margin-right: -20px; margin-top: -5px">
    <iframe id="show-iframe" :src='url_src'
            frameborder="no"
            marginheight="0" marginwidth="0"
            scrolling="yes"/>
  </div>
</template>

<style scoped>
  #show-iframe {
    width: 100%;
    height: 100% !important;
    background-color: transparent;
    overflow-x: scroll;
    overflow-y: hidden;
  }
</style>

<script>
  const axios = require('axios')
  import {getCookie, setCookie, setRootCookie, isEmpty} from '../conf/utils'

  export default {
    data () {
      return {
        url_src: '',
        timestamp: 0,
        sec186: 0,
        iframe_loaded: false,
        rolesData: [],
      }
    },
    methods: {
      initRedashURL () {
        let self = this
        axios.get('/auth/roles')
          .then(function (response) {
            var res = []
            response.data.forEach(ele => {
              res.push({'roleName': ele})
            })
            self.rolesData = res
          }).catch(function (error) {
          console.log(error)
        })
      },
      setInframeHeight () {
        let self = this

        var oIframe = document.getElementById('show-iframe')
        if (oIframe) {
          var iframeWin = oIframe.contentWindow || oIframe.contentDocument.parentWindow

          if (iframeWin.document.body) {
            oIframe.style.height = (Number(iframeWin.document.documentElement.scrollHeight || iframeWin.document.body.scrollHeight)) + 'px'
            console.log(oIframe.height)
          }
        }
        self.iframe_loaded = true
      },
    },
    created: function () {
      let self = this
      console.log('created')

      axios.get('/redash/timestamp')
        .then(function (response) {
          var res = []
          self.timestamp = response.data
          self.sec186 = (self.timestamp + 7482952) * 3
          setRootCookie('sec186', self.sec186.toString(), 1)
          // self.url_src = 'http://redash.adtech.sogou/public/dashboards/vSCVP1S0jSkDJmNHftcPLGdbcC5ivRbanqtlT9zA?org_slug=default&refresh=120'

        }).catch(function (error) {
        console.log(error)
      })
    },
    mounted () {
      let redashUrl = sessionStorage.getItem('redash_url')
      console.log('redashUrl:' + redashUrl)
      this.url_src = redashUrl
      const oIframe = document.getElementById('show-iframe')
      const deviceWidth = document.documentElement.clientWidth
      const deviceHeight = document.documentElement.clientHeight
      oIframe.style.height = (Number(deviceHeight) - 200) + 'px' //数字是页面布局高度差
    }
  }
</script>

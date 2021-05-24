<template>
  <div v-bind:style="{'height':myHeighgt+'px'}" class="outDiv">
<!--    style="height:780px !important; margin-left: -20px;margin-right: -20px; margin-top: -5px">-->

    <iframe id="show-iframe" v-if=isShowFrame :src='url_src'
            frameborder="no"
            marginheight="0" marginwidth="0"
            scrolling="yes"
    style="height: 100%"/>
    <h3 style="display:none">
        This is redashshow/{{$route.params.redashname}}
    </h3>
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
  .outDiv{
    margin-left: -20px;
    margin-right: -20px;
    margin-top: -5px;
  }
</style>

<script>
  const axios = require('axios')
  import {getCookie,setCookie,setRootCookie, isEmpty} from '../conf/utils'

  export default {
    data () {
      return {
        url_src: '',
        timestamp: 0,
        sec186: 0,
        isShowFrame: false,
        rolesData: [],
        myHeighgt:'400'
      }
    },
    methods: {
      updateCookie () {
        let self = this
        axios.get('/redash/timestamp')
          .then(function (response) {
            var res = []
            self.timestamp = response.data
            self.sec186 = (self.timestamp+7482952)*3
            setRootCookie('sec186', self.sec186.toString(), 1);

            self.isShowFrame = true

          }).catch(function (error) {
          console.log(error)
        })
        console.log('++++++++++++++++++')
      },
    },
    created: function () {
      console.log('created')

      let self = this
      self.isShowFrame = false
      this.url_src = sessionStorage.getItem(self.$route.params.redashname)
      console.log('redashUrl:'+this.url_src)

      axios.get('/redash/timestamp')
          .then(function (response) {
            var res = []
            self.timestamp = response.data
            self.sec186 = (self.timestamp+7482952)*3
            setRootCookie('sec186', self.sec186.toString(), 1);

            self.isShowFrame = true

          }).catch(function (error) {
          console.log(error)
        })
    },

    mounted () {
      console.log('mounted');
      let height = document.body.clientHeight
      console.log(height)
      this.myHeighgt = height - 80
    },

    beforeUpdate: function () {
      console.log('before update')
      this.url_src = sessionStorage.getItem(this.$route.params.redashname)
    },

    updated () {
      console.log("updated")
    }
  }
</script>

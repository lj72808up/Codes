<template>
  <div>
    <div v-bind:style="{'height':myHeighgt+'px'}" class="outDiv">
      <iframe id="show-iframe" src='http://redash.adtech.sogou/public/dashboards/8xfxxyHcF3AQNfbxmPZJ1PhGzxNh9FBhJOSUM16y?org_slug=default&refresh=1200'
              frameborder="no"
              marginheight="0" marginwidth="0"
              scrolling="yes"
              style="height: 100%"/>
    </div>
  </div>
</template>

<script>

  export default {
    name: 'DatabasesContent',
    data () {
      return {
        myHeighgt:""
      }
    },
    methods: {

    },
    mounted () {
      let height = document.body.clientHeight
      console.log(height)
      this.myHeighgt = height - 15
    },
  }
</script>

<style>
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

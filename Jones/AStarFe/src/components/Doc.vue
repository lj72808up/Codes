<template>
    <div>
      <el-row>
        <h3 class="inline-title">帮助与说明 </h3>
        <el-button
          type="text"
          icon="el-icon-edit"
          @click="redirectClick"
          circle></el-button>
      </el-row>

      <div v-html="content"></div>
    </div>
</template>

<script>
  import router from '../router';
  const axios = require('axios');

  export default {
    data() {
      return {
        content: ""
      }
    },
    methods:{
      handleRefreshingContent: function() {
        let self = this;
        axios.get('doc/gethtml')
          .then(function (response) {
            self.content = response.data
          })
          .catch(function (err) {
            self.$notify.error({
              title: '错误',
              message: err
            })
          });
      },
      redirectClick(){
        router.push({path: 'docedit'})
      },
    },
    created: function () {
      this.handleRefreshingContent();
    }
  }
</script>

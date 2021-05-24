<template>
  <div>
    <el-row>
      <h3 class="inline-title">文档编辑</h3>
      <span></span>
      <el-button
        type="text"
        icon="el-icon-view"
        @click="redirectClick"
        circle></el-button>
    </el-row>

    <!--
    <div>
      <el-row>
        <vue-editor
          v-model="docForm.content"
          :editorToolbar="customToolbar">
        </vue-editor>
      </el-row>
    </div>
    -->

    <div id="editor">
        <mavon-editor 
          style="height: 100%"
          v-model = 'docForm.content'
          :ishljs="true"
          ref=md @imgAdd="$imgAdd" @imgDel="$imgDel" @save="$saveDoc"          
        />
    </div>

    <!--
    <div class="formCommit">
      <el-row type="flex" justify="end">
        <el-col :span="6">
          <el-button type="primary" @click="handleSavingContent">保存</el-button>
          <el-button @click="handleRefreshingContent">刷新</el-button>
        </el-col>
      </el-row>
    </div>
    -->

  </div>
</template>

<script>
  import {getCookie } from "../conf/utils";
  import router from '../router';

  import { mavonEditor } from 'mavon-editor'
  import 'mavon-editor/dist/css/index.css'
  import qs from 'qs'

  const axios = require('axios');

  export default {
    name: 'editor',
    components: {
        mavonEditor,
        // or 'mavon-editor': mavonEditor
    },
    data(){
      return{
        docForm:{
          content: "",    // md 原始文本
          html: "",       //html 解析后的文本
          submitter: "ADMIN",
        },
        //editorContent:'',
        img_file:[],
        code_style:'github',

      }
    },
    methods:{
      rndStr(len) {
        let text = ""
        let chars = "abcdefghijklmnopqrstuvwxyz"
      
        for( let i=0; i < len; i++ ) {
          text += chars.charAt(Math.floor(Math.random() * chars.length))
        }
        return text
      },      
      // 绑定@imgAdd event
      $imgAdd(pos, $file) {
          // 第一步.将图片上传到服务器.
          let self = this
          this.img_file[pos] = $file;
          let param = "sign_encode1=" + this.rndStr(16) + "&" + "encode1=" + encodeURIComponent($file['miniurl'])
          
          /* 因为axios 设置了 axios.defaults.baseURL，所以只能用http全路径访问， 以免请求到 Astar中
          axios.post('http://datasearch.adtech.sogou:8080/yuntu/http_upload?appid=201061', param )
            .then(function (response) {
                let _res = response.data[0];
                console.log("=================111=")
                console.log(_res)
                console.log("==================111")
                // 第二步.将返回的url替换到文本原位置![...](0) -> ![...](url)
                self.$refs.md.$img2Url(pos, _res.url);                

            }).catch(function (error) {
            console.log(error)
          })
          */
          fetch('/yuntu/http_upload?appid=201061', {
              method: "POST",
              //mode: 'cors',
              headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
              body:param,
            }
          ).then(response => {
            return response.json()
          })
          .then((data) => {
              let _res = data[0];
              // 第二步.将返回的url替换到文本原位置![...](0) -> ![...](url)
              self.$refs.md.$img2Url(pos, _res.url);
          })
          .catch((error) => {
            console.log(error)
          })
          
      },
      $imgDel(pos) {
          delete this.img_file[pos];
      },

      $saveDoc(value, render ) {
          let self = this;
          self.docForm['content'] = value;
          self.docForm['html'] = render;

          axios({
            url: 'doc/submit',
            method: 'post',
            data: this.docForm
          }).then((response) => {
              self.$notify.success({
                title: '提示',
                message: response.data
              });
            })
            .catch((err) => {
              self.$notify.error({
                title: '错误',
                message: err
              })
            })          
      },


      handleRefreshingContent: function() {
        let self = this;
        axios.get('doc/get')
          .then(function (response) {
            self.docForm.content = response.data
          })
          .catch(function (err) {
            self.$notify.error({
              title: '错误',
              message: err
            })
          });
      },
      redirectClick(){
        router.push({path: 'doc'})
      },
    },
    mounted: function () {
      this.handleRefreshingContent();
    //  read author from cookie
      this.docForm.submitter = getCookie('_adtech_user')
    },

  }
</script>

<style>
#editor {
    margin: auto;
    width: 100%;
    height: 580px;
}
</style>
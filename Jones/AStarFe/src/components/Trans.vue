<template>

  <div>
    <h3>一键转化工具</h3>
    <el-row>
      <span>转化前原始文件上传或填写。上传文件使用换行分隔，直接填写使用逗号分隔。</span>
      <span>注意imei的前后不要有空格！</span>
    </el-row>
    <el-row>
      <el-col :span="14">
        <el-input v-model="userInfo"></el-input>
      </el-col>
      <el-col :span="2" style="text-align: center">
        <el-upload
          class="upload-demo"
          action="123"
          :before-upload="beforeUpload"
          show-file-list=false>
          <i class="el-icon-upload"></i>
        </el-upload>
      </el-col>
      <el-col :span="4" :offset="1">
        <el-button type="success" @click="onTrans">转换</el-button>
      </el-col>
    </el-row>
    <div v-show="isShowRes">
      <span>您的提交的转化queryId: </span><strong>{{ transQueryId }}</strong>
      <h3>查询结果</h3>
      <span><i class="el-icon-bell"></i> 同步任务查询耗时较长。请耐心等待或转到任务队列查看。</span>
      <el-table
        v-loading="transLoading"
        :data="transData"
        style="width: 100%">
        <el-table-column
          v-for="sqlDataCol in transDataCols"
          :key="sqlDataCol.prop"
          :prop="sqlDataCol.prop"
          :label="sqlDataCol.label"
          :sortable="sqlDataCol.sortable">
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script>
  const axios = require('axios');
  export default {
    data() {
      return {
        isShowRes: false,
        userInfo: '',
        transQueryId:'',
        transLoading: true,
        transDataCols: [],
        transData : [],
      }
    },
    methods: {
      beforeUpload(file) {
        let self = this;
        var formData = new FormData();
        formData.append("file", file);
        axios.post('task/file', formData,  {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
          .then(function (response) {
            self.userInfo = response.data
          }).catch(function (error) {
          console.log(error)
        })
        return false ; // 返回false不会自动上传
      },
      onTrans () {
        let self = this;
        this.$message({
          type: 'success',
          message: '提交成功!'
        });
        var params = new URLSearchParams();
        params.append('sql', this.userInfo);
        axios({
          method: 'post',
          url: '/dw/sql/id',
          data: params,
          responseType: 'string'
        }).then(function (response) {
            self.transQueryId = response.data;
            self.isShowRes = true;
            self.$options.methods.onTransTask.bind(self)();
          })
          .catch(function (error) {
            this.$message({
              type: 'danger',
              message: '任务失败!'
            });
            console.log(error);
          });
      },
      onTransTask(){
        var params = new URLSearchParams();
        params.append('sql', this.userInfo);
        params.append('queryId', this.transQueryId)
        axios.post('/dw/sql/custom', params, {responseType: 'json'})
          .then(function (response) {
            self.transDataCols = response.data.cols;
            self.transData = response.data.data;
            self.transLoading = false;
          })
          .catch(function (error) {
            this.$message({
              type: 'danger',
              message: '任务失败!'
            });
            self.isShowRes = false;
            console.log(error);
          });
      }
    },

  }
</script>



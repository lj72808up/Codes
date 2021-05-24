<template>
  <div>
    <h3>数据查询</h3>
    <el-form ref="dataForm" :model="dataForm" label-width="80px">
      <el-row>
        <el-col :span="14">
          <el-form-item label="项目ID" prop="name">
            <el-input v-model="dataForm.task_id" placeholder="请输入" clearable></el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-form-item label="起止日期">
          <el-date-picker
          v-model="dataForm.date_range"
          type="daterange"
          align="right"
          unlink-panels
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="yyyyMMdd">
          </el-date-picker>
        </el-form-item>
      </el-row>

      <el-row>
        <el-form-item label="上传提示">
          <div><span>请上传UTF8编码的纯文本文件，注意使用换行分隔。上传成功后内容会显示在输入框中。</span></div>
        </el-form-item>
      </el-row>

      <el-row>
        <el-col :span="11">
          <el-form-item label="客户ID" >
            <el-input v-model="dataForm.custom_id" type="array" placeholder='请注意使用","分隔' clearable></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="1" style="text-align: center">
          <el-upload
            class="upload-demo"
            action="123"
            :before-upload="beforeUpload"
            :show-file-list=false>
            <i class="el-icon-upload2" ></i>
          </el-upload>
        </el-col>
      </el-row>

      <el-row>
        <el-form-item label="结果筛选">
          <el-col :span="5">
            <el-switch
              v-model="dataForm.is_mobil_qq"
              active-text="手Q渠道结果">
            </el-switch>
          </el-col>
          <el-col :span="4">
            <el-switch
              v-model="dataForm.is_same"
              active-text="同质结果">
            </el-switch>
          </el-col>
          <el-col :span="4">
            <el-switch
              v-model="dataForm.is_avg"
              active-text="分日结果">
            </el-switch>
          </el-col>
          <!--<el-col :span="4">-->
            <!--<el-switch-->
              <!--v-model="dataForm.is_custom"-->
              <!--active-text="客户维度">-->
            <!--</el-switch>-->
          <!--</el-col>-->
        </el-form-item>
      </el-row>
      <el-row>
        <el-form-item label="查询筛选">
          <el-radio-group v-model="dataForm.query" size="small">
            <el-radio-button label="整体页面"></el-radio-button>
            <el-radio-button label="实验样式-样式"></el-radio-button>
            <el-radio-button label="实验样式-页面"></el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-row>
      <el-row>
        <el-form-item label="下载筛选">
          <el-radio-group v-model="dataForm.download" size="small" @change="changeSame(dataForm.download)">
            <el-radio-button label="实验样式页面"></el-radio-button>
            <el-radio-button label="实验样式"></el-radio-button>
            <el-radio-button label="实验样式同质页面"></el-radio-button>
            <el-radio-button label="同质样式"></el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-row>
      <el-row type="flex" justify="end">
        <el-col :span="6">
          <el-button type="primary" @click="submitDataForm">立即查询</el-button>
          <el-button @click="downloadDataForm">提交下载</el-button>
        </el-col>
      </el-row>
    </el-form>

    <h3>数据结果</h3>
    <el-row>
      <span>点击按钮即可下载数据表。</span>
      <el-button type="success" icon="el-icon-download" circle size="mini" @click="exportExcel"></el-button>
    </el-row>
    <el-table id="abTable"
      v-loading="resLoading"
      :data="resData"
      stripe
      style="width: 100%">
      <el-table-column
        v-for="resCol in resCols"
        v-if="!resCol.hidden"
        :key="resCol.prop"
        :prop="resCol.prop"
        :label="resCol.label"
        :sortable="resCol.sortable">
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
  import {formatDate, isEmpty} from '../conf/utils'
  import router from '../router';
  import XLSX from 'xlsx';
  import FileSaver from 'file-saver';
  import global_ from './Global';

  const axios = require('axios');

  export default {
    data(){
      return{
        dataForm:{
          task_id:'',
          date_range:global_.dateRange,
          custom_id: [],
          is_avg: false,
          is_mobil_qq: false,
          is_same: false,
          is_custom: false,
          query: "整体页面",
          download: "实验样式页面",
          type: 0
        },
        resCols:[],
        resData: [],
        resLoading: false,
      }
    },
    methods:{
      changeSame (con) {
        let self = this;
        self.dataForm.is_same = con === '实验样式同质页面' || con === '同质样式';
      },
      beforeUpload (file) {
        let self = this;
        var formData = new FormData();
        formData.append("file", file);
        axios.post('task/file', formData,  {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
          .then(function (response) {
            self.dataForm.custom_id = response.data
          }).catch(function (error) {
          console.log(error)
        })
        return false ;
      },
      redirectClick(){
        router.push({path: 'que'})
      },
      msgNotify(queryId) {
        this.$notify({
          title: queryId,
          message: '服务器接收到您的请求并分配查询ID。您可以转到任务队列页面查看结果！',
          type: 'success',
          duration: 2000,
          onClick: this.redirectClick,
        });
      },
      errorNotify(msg) {
        this.$notify.error({
          title: '错误',
          message: msg
        });
      },
      submitDataForm(){
        let self = this;
        // empty case
        if ( this.dataForm.custom_id === ""){
          this.dataForm.custom_id = [];
        }
        // custom_id not empty string to array
        if (typeof this.dataForm.custom_id === "string") {
          const tmp = this.dataForm.custom_id;
          this.dataForm.custom_id = tmp.split(",")
        }
        self.resLoading = true;
        axios({
          method:'post',
          url:'/small/query',
          data: JSON.stringify(self.dataForm)
        }).then(function (response) {
             self.resLoading = false;
             self.resCols = response.data.cols;
             self.resData = response.data.data;
          })
          .catch(function (error) {
            self.errorNotify(error);
          });
      },
      downloadDataForm(){
        let self = this;
        // empty case
        if ( this.dataForm.custom_id === ""){
          this.dataForm.custom_id = [];
        }
        // custom_id not empty string to array
        if (typeof this.dataForm.custom_id === "string") {
          const tmp = this.dataForm.custom_id;
          this.dataForm.custom_id = tmp.split(",")
        }
        axios({
          method:'post',
          url:'/small/download',
          data: JSON.stringify(self.dataForm)
        })
          .then(function (response) {
          self.$options.methods.msgNotify.bind(self)(response.data);
        })
          .catch(function (error) {
            self.errorNotify("出错了！请联系管理员。");
            console.log(error)
          });
      },
      exportExcel () {
        /* generate workbook object from table */
        let wb = XLSX.utils.table_to_book(document.querySelector("#abTable"));
        /* get binary string as output */
        let wbout = XLSX.write(wb, { bookType: 'xlsx', bookSST: true, type: 'array' });
        try {
          FileSaver.saveAs(new Blob([wbout], { type: 'application/octet-stream' }), "ABtestRes.xlsx");
        } catch (e)
        {
          if (typeof console !== 'undefined')
            console.log(e, wbout)
        }
        return wbout
      }
    },
    created: function () {
      const abtestInfo = JSON.parse(sessionStorage.getItem("abtest_info"));
      if (abtestInfo){
        this.dataForm.task_id = abtestInfo.task_id;
      }
    }
  }
</script>

<template>
  <div>
    <h3>银河数据查询</h3>
    <el-form ref="galaxyForm" :model="galaxyForm" :rules="galaxyRules" size="medium" label-width="100px">
      <el-form-item prop="dimensions" label-width="0px">
        <el-checkbox-group v-model="galaxyForm.dimensions">
          <el-row>
            <el-col :span="2">
              <el-tag size="medium">基本</el-tag>
            </el-col>
            <el-col :span="22">
              <el-checkbox v-for="(value, key) in galaxyBase" :label="value" :key="value">{{key}}</el-checkbox>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="2">
              <el-tag size="medium">广告</el-tag>
            </el-col>
            <el-col :span="22">
              <el-checkbox v-for="(value, key) in galaxyGG" :label="value" :key="value">{{key}}</el-checkbox>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="2">
              <el-tag size="medium">客户</el-tag>
            </el-col>
            <el-col :span="22">
              <el-checkbox v-for="(value, key) in galaxyKH" :label="value" :key="value">{{key}}</el-checkbox>
            </el-col>
          </el-row>
        </el-checkbox-group>
      </el-form-item>
      <hr>

      <h3>筛选条件</h3>
      <el-row>
        <el-form-item label="起止日期" prop="date">
          <el-date-picker
            v-model="galaxyForm.date"
            type="daterange"
            align="right"
            unlink-panels
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="yyyyMMdd"
            :picker-options="galaxyDateOptions">
          </el-date-picker>
        </el-form-item>
      </el-row>

      <!--may have file-->
      <el-row>
        <el-form-item label="上传提示">
          <div><span>请上传UTF8编码的纯文本文件，注意使用换行分隔。上传成功后内容会显示在输入框中。</span></div>
        </el-form-item>
      </el-row>

      <el-row>
        <el-col :span="7">
          <el-form-item label="代理商ID" prop="filters.DLSID.values">
            <el-input v-model="galaxyForm.filters.DLSID.values"></el-input>
          </el-form-item>
        </el-col>
        <el-col class='type' :span="4">
          <el-form-item label="类型" label-width="50px" prop="filters.DLSID.type">
            <el-select v-model="galaxyForm.filters.DLSID.type" placeholder="">
              <el-option label="黑名单" value="NOTIN"></el-option>
              <el-option label="白名单" value="IN"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="1" style="text-align: center">
          <el-upload
            class="upload-demo"
            action="123"
            :before-upload="beforeUploadDLSID"
            :show-file-list=false>
            <i class="el-icon-upload2" ></i>
          </el-upload>
        </el-col>
        <el-col :span="7">
          <el-form-item label="客户ID" prop="filters.KHID.values">
            <el-input v-model="galaxyForm.filters.KHID.values"></el-input>
          </el-form-item>
        </el-col>
        <el-col class='type' :span="4">
          <el-form-item label="类型" label-width="45px" prop="filters.KHID.type">
            <el-select v-model="galaxyForm.filters.KHID.type" placeholder="">
              <el-option label="黑名单" value="NOTIN"></el-option>
              <el-option label="白名单" value="IN"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="1" style="text-align: center">
          <el-upload
            class="upload-demo"
            action="123"
            :before-upload="beforeUploadKHID"
            :show-file-list=false>
            <i class="el-icon-upload2" ></i>
          </el-upload>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="7">
          <el-form-item label="查询词" prop="filters.CXC.values">
            <el-input v-model="galaxyForm.filters.CXC.values"></el-input>
          </el-form-item>
        </el-col>
        <el-col class='type' :span="4" prop="filters.CXC.type">
          <el-form-item label="类型" label-width="50px">
            <el-select v-model="galaxyForm.filters.CXC.type" placeholder="">
              <el-option label="黑名单" value="NOTIN"></el-option>
              <el-option label="白名单" value="IN"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="1" style="text-align: center">
          <el-upload
            class="upload-demo"
            action="123"
            :before-upload="beforeUploadCXC"
            :show-file-list=false>
            <i class="el-icon-upload2"></i>
          </el-upload>
        </el-col>
        <el-col :span="7">
          <el-form-item label="计划ID" prop="filters.JHID.values">
            <el-input v-model="galaxyForm.filters.JHID.values"></el-input>
          </el-form-item>
        </el-col>
        <el-col class='type' :span="4" prop="filters.JHID.type">
          <el-form-item label="类型" label-width="45px">
            <el-select v-model="galaxyForm.filters.JHID.type">
              <el-option label="黑名单" value="NOTIN"></el-option>
              <el-option label="白名单" value="IN"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="1" style="text-align: center">
          <el-upload
            class="upload-demo"
            action="123"
            :before-upload="beforeUploadJHID"
            :show-file-list=false>
            <i class="el-icon-upload2" ></i>
          </el-upload>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="7">
          <el-form-item label="组ID" prop="filters.ZID.values">
            <el-input v-model="galaxyForm.filters.ZID.values"></el-input>
          </el-form-item>
        </el-col>
        <el-col class='type' :span="4">
          <el-form-item label="类型" label-width="50px" prop="filters.ZID.type">
            <el-select v-model="galaxyForm.filters.ZID.type" placeholder="">
              <el-option label="黑名单" value="NOTIN"></el-option>
              <el-option label="白名单" value="IN"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="1" style="text-align: center">
          <el-upload
            class="upload-demo"
            action="123"
            :before-upload="beforeUploadZID"
            :show-file-list=false>
            <i class="el-icon-upload2" ></i>
          </el-upload>
        </el-col>
        <el-col :span="7">
          <el-form-item label="PID" prop="filters.PID.values">
            <el-input v-model="galaxyForm.filters.PID.values"></el-input>
          </el-form-item>
        </el-col>
        <el-col class='type' :span="4">
          <el-form-item label="类型" label-width="45px" prop="filters.PID.type">
            <el-select v-model="galaxyForm.filters.PID.type" placeholder="">
              <el-option label="黑名单" value="NOTIN"></el-option>
              <el-option label="白名单" value="IN"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="1" style="text-align: center">
          <el-upload
            class="upload-demo"
            action="123"
            :before-upload="beforeUploadPID"
            :show-file-list=false>
            <i class="el-icon-upload2"></i>
          </el-upload>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="7">
          <el-form-item label="关键词" prop="filters.GJC.values">
            <el-input v-model="galaxyForm.filters.GJC.values"></el-input>
          </el-form-item>
        </el-col>
        <el-col class='type' :span="4">
          <el-form-item label="类型" label-width="50px" prop="filters.GJC.type">
            <el-select v-model="galaxyForm.filters.GJC.type">
              <el-option label="黑名单" value="NOTIN"></el-option>
              <el-option label="白名单" value="IN"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="1" style="text-align: center">
          <el-upload
            class="upload-demo"
            action="123"
            :before-upload="beforeUploadGJC"
            :show-file-list=false>
            <i class="el-icon-upload2"></i>
          </el-upload>
        </el-col>
        <el-col :span="7">
          <el-form-item label="关键词ID" prop="filters.GJCID.values">
            <el-input v-model="galaxyForm.filters.GJCID.values"></el-input>
          </el-form-item>
        </el-col>
        <el-col class='type' :span="4">
          <el-form-item label="类型" label-width="45px" prop="filters.GJCID.type">
            <el-select v-model="galaxyForm.filters.GJCID.type" placeholder="">
              <el-option label="黑名单" value="NOTIN"></el-option>
              <el-option label="白名单" value="IN"></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="1" style="text-align: center">
          <el-upload
            class="upload-demo"
            action="123"
            :before-upload="beforeUploadGJCID"
            :show-file-list=false>
            <i class="el-icon-upload2"></i>
          </el-upload>
        </el-col>
      </el-row>

      <!--other filters-->
      <el-row>
        <el-form-item label="勾选提示">
          <div><span>选项太多找不到？直接搜索试试。</span></div>
        </el-form-item>
      </el-row>

      <el-row>
        <el-col :span="11">
          <el-form-item label="一级客户行业" prop="filters.YJKHHY.values">
            <el-select v-model="galaxyForm.filters.YJKHHY.values" multiple filterable placeholder="请选择">
              <el-option
                v-for="value in yjkhhyOptions"
                :key="value"
                :label="value"
                :value="value">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="11" :offset="1">
          <el-form-item label="二级客户行业" prop="filters.EJKHHY.values">
            <el-select v-model="galaxyForm.filters.EJKHHY.values" multiple filterable placeholder="请选择">
              <el-option
                v-for="value in ejkhhyOptions"
                :key="value"
                :label="value"
                :value="value">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="11">
          <el-form-item label="网民区域" prop="filters.WMQY.values">
            <el-select v-model="galaxyForm.filters.WMQY.values" multiple filterable placeholder="请选择">
              <el-option-group
                v-for="group in wmqyOptions"
                :key="group.label"
                :label="group.label">
                <el-option
                  v-for="value in group.children"
                  :key="value"
                  :label="value"
                  :value="value">
                </el-option>
              </el-option-group>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="11" :offset="1">
          <el-form-item label="大区">
            <el-select v-model="galaxyForm.filters.DQ.values" multiple filterable placeholder="请选择">
              <el-option
                v-for="value in dqOptions"
                :label="value"
                :key="value"
                :value="value">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>

      </el-row>

        <el-row>
          <el-col :span="11">
            <el-form-item label="二级业务类型" prop="filters.EJYWLX.values">
              <el-select v-model="galaxyForm.filters.EJYWLX.values" multiple filterable placeholder="请选择">
                <el-option
                  v-for="value in ejywlxOptions"
                  :label="value"
                  :key="value"
                  :value="value">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="11" :offset="1">
            <el-form-item label="三级业务类型" prop="filters.SJYWLX.values">
              <el-select v-model="galaxyForm.filters.SJYWLX.values" multiple filterable placeholder="请选择">
                <el-option
                  v-for="value in sjywlxOptions"
                  :key="value"
                  :label="value"
                  :value="value">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row>
          <el-col >
            <el-form-item label="客户类型" prop="filters.KHLX.values">
              <el-checkbox-group v-model="galaxyForm.filters.KHLX.values">
                <el-checkbox v-for="value in khlxOptions" :label="value" :key="value">{{value}}</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
          </el-col>
        </el-row>

        <!--<el-row>-->
          <!--<el-col>-->
            <!--<el-form-item label="样式类别">-->
              <!--<el-checkbox-group v-model="galaxyForm.filters.YSLB.values">-->
                <!--<el-checkbox v-for="value in yslbOptions" :label="value" :key="value">{{value}}</el-checkbox>-->
              <!--</el-checkbox-group>-->
            <!--</el-form-item>-->
          <!--</el-col>-->
        <!--</el-row>-->

        <!--<el-row>-->
          <!--<el-col>-->
            <!--<el-form-item label="排位">-->
              <!--<el-checkbox-group v-model="galaxyForm.filters.PW.values">-->
                <!--<el-checkbox v-for="value in pwOptions" :label="value" :key="value">{{value}}</el-checkbox>-->
              <!--</el-checkbox-group>-->
            <!--</el-form-item>-->
          <!--</el-col>-->
        <!--</el-row>-->

        <!--<el-row>-->
          <!--<el-col>-->
            <!--<el-form-item label="触发类型">-->
              <!--<el-checkbox-group v-model="galaxyForm.filters.CFLX.values">-->
                <!--<el-checkbox v-for="value in cflxOptions" :label="value" :key="value">{{value}}</el-checkbox>-->
              <!--</el-checkbox-group>-->
            <!--</el-form-item>-->
          <!--</el-col>-->
        <!--</el-row>-->

        <hr>

        <h3>输出指标</h3>
        <el-row>
          <el-form-item label="基础指标" prop="outputs">
              <el-checkbox-group v-model="galaxyForm.outputs">
                <el-checkbox label="DJ">点击</el-checkbox>
                <el-checkbox label="XH">消耗</el-checkbox>
              </el-checkbox-group>
          </el-form-item>
        </el-row>

        <hr>
        <h3>任务信息</h3>
        <el-row>
          <el-col>
            <el-form-item label="任务名称" prop="name">
              <el-col :span="10">
                <el-input v-model="galaxyForm.name" placeholder="请输入任务名称"></el-input>
              </el-col>
            </el-form-item>
          </el-col>
        </el-row>

    </el-form>

    <div class="formCommit">
      <el-row type="flex" justify="end">
        <el-col :span="6">
          <el-button type="primary" @click="onSubmit('galaxyForm')">立即提交</el-button>
          <el-button @click="onReset('galaxyForm')">重置</el-button>
        </el-col>
      </el-row>
    </div>

  </div>
</template>

<script>
  const axios = require('axios');
  import {getFileContent, isEmpty} from '../conf/utils'
  import global_ from './Global'
  export default {
    data(){
      return{
        galaxyBase: global_.galaxyBase,
        galaxyGG: global_.galaxyGG,
        galaxyKH: global_.galaxyKH,

        galaxyRules: global_.strictRules,
        galaxyDateOptions: global_.dateOptions,

        // options
        wmqyOptions: global_.provincesChina,
        yjkhhyOptions:[],
        ejkhhyOptions:[],
        sjkhhyOptions:[],
        khlxOptions:[],
        dqOptions:[],
        ejywlxOptions:[],
        sjywlxOptions:[],
        // yslbOptions:[],
        // pwOptions:[],
        // cflxOptions:[],

        // submit form data
        galaxyForm: {
          date: global_.dateRange,
          dimensions: [],
          filters:{
            DLSID:{
              name:'DLSID',
              values:[],
              type:'IN'
            },
            KHID:{
              name:'KHID',
              values:[],
              type:'IN'
            },
            CXC:{
              name:'CXC',
              values:[],
              type:'IN'
            },
            JHID:{
              name:'JHID',
              values:[],
              type:'IN'
            },
            ZID:{
              name:'ZID',
              values:[],
              type:'IN'
            },
            PID:{
              name:'PID',
              values:[],
              type:'IN'
            },
            GJC:{
              name:'GJC',
              values:[],
              type:'IN'
            },
            GJCID:{
              name:'GJCID',
              values:[],
              type:'IN'
            },
            QDMC:{
              name:'QDMC',
              values:[],
              type:'IN'
            },
            WMQY:{
              name:'WMQY',
              values:[],
              type:'IN'
            },
            YJKHHY:{
              name:'YJKHHY',
              values:[],
              type:'IN'
            },
            EJKHHY:{
              name:'EJKHHY',
              values:[],
              type:'IN'
            },
            SJKHHY:{
              name:'SJKHHY',
              values:[],
              type:'IN'
            },
            KHLX:{
              name:'KHLX',
              values:[],
              type:'IN'
            },
            DQ:{
              name:'DQ',
              values:[],
              type:'IN'
            },
            EJYWLX:{
              name:'EJYWLX',
              values:[],
              type:'IN'
            },
            SJYWLX:{
              name:'SJYWLX',
              values:[],
              type:'IN'
            },
            YSLB:{
              name:'YSLB',
              values:[],
              type:'IN'
            },
            PW:{
              name:'PW',
              values:[],
              type:'IN'
            },
            CFLX:{
              name:'CFLX',
              values:[],
              type:'IN'
            }
          },
          outputs: [],
          name: '',
        },

        queryId: '',
        errorMsg:''

      }
    },
    methods:{
      onSubmit(formName){
        let self = this;
        this.$refs[formName].validate((valid)=> {
          if (valid) {
            // build my filter json object
            var postFilters = [];
            for (var item in this.galaxyForm.filters) {
              if (! isEmpty(this.galaxyForm.filters[item].values)) {
                if (typeof this.galaxyForm.filters[item].values === "string") {
                  const tmp = this.galaxyForm.filters[item].values;
                  this.galaxyForm.filters[item].values = tmp.split(",")
                }
                postFilters.push(this.galaxyForm.filters[item])
              }
            }
            this.galaxyForm.filters = postFilters;
            axios({
              method: 'post',
              url: '/jone/combination_yh',
              data: JSON.stringify(this.galaxyForm)
            }).then(function (response) {
              self.queryId = response.data;
              self.$options.methods.msgNotify.bind(self)();
            })
            .catch(function (error) {
              self.errorMsg = error;
              self.$options.methods.errorNotify.bind(self)();
            });
          }
        });
      },
      onReset(formName){
        this.$refs[formName].resetFields();
      },
      msgNotify() {
        this.$notify.info({
          title: this.queryId,
          message: '服务器接收到您的请求并分配查询ID。您可以转到任务队列页面查看结果！',
          type: 'success'
        });
      },
      errorNotify() {
        this.$notify.error({
          title: '错误',
          message: this.errorMsg
        });
      },
      beforeUploadDLSID (file) {
        let self = this;
        var formData = new FormData();
        formData.append("file", file);
        axios.post('task/file', formData,  {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        .then(function (response) {
          self.galaxyForm.filters.DLSID.values = response.data
        }).catch(function (error) {
          self.errorMsg = error;
          self.$options.methods.errorNotify.bind(self)();
        })
        return false ; // 返回false不会自动上传
      },
      beforeUploadKHID (file) {
        let self = this;
        var formData = new FormData();
        formData.append("file", file);
        axios.post('task/file', formData,  {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        .then(function (response) {
          self.galaxyForm.filters.KHID.values = response.data
        }).catch(function (error) {
          self.errorMsg = error;
          self.$options.methods.errorNotify.bind(self)();
        })
        return false ; // 返回false不会自动上传
      },
      beforeUploadCXC (file) {
        let self = this;
        var formData = new FormData();
        formData.append("file", file);
        axios.post('task/file', formData,  {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        .then(function (response) {
          self.galaxyForm.filters.CXC.values = response.data
        }).catch(function (error) {
          self.errorMsg = error;
          self.$options.methods.errorNotify.bind(self)();
        })
        return false ;
      },
      beforeUploadJHID (file) {
        let self = this;
        var formData = new FormData();
        formData.append("file", file);
        axios.post('task/file', formData,  {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        .then(function (response) {
          self.galaxyForm.filters.JHID.values = response.data
        }).catch(function (error) {
          self.errorMsg = error;
          self.$options.methods.errorNotify.bind(self)();
        })
        return false ;
      },
      beforeUploadZID (file) {
        let self = this;
        var formData = new FormData();
        formData.append("file", file);
        axios.post('task/file', formData,  {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        .then(function (response) {
          self.galaxyForm.filters.ZID.values = response.data
        }).catch(function (error) {
          self.errorMsg = error;
          self.$options.methods.errorNotify.bind(self)();
        })
        return false ;
      },
      beforeUploadPID (file) {
        let self = this;
        var formData = new FormData();
        formData.append("file", file);
        axios.post('task/file', formData,  {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        .then(function (response) {
          self.galaxyForm.filters.PID.values = response.data
        }).catch(function (error) {
          self.errorMsg = error;
          self.$options.methods.errorNotify.bind(self)();
        })
        return false ;
      },
      beforeUploadGJC (file) {
        let self = this;
        var formData = new FormData();
        formData.append("file", file);
        axios.post('task/file', formData,  {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        .then(function (response) {
          self.galaxyForm.filters.GJC.values = response.data
        }).catch(function (error) {
          self.errorMsg = error;
          self.$options.methods.errorNotify.bind(self)();
        })
        return false ;
      },
      beforeUploadGJCID (file) {
        let self = this;
        var formData = new FormData();
        formData.append("file", file);
        axios.post('task/file', formData,  {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        .then(function (response) {
          self.galaxyForm.filters.GJCID.values = response.data
        }).catch(function (error) {
          self.errorMsg = error;
          self.$options.methods.errorNotify.bind(self)();
        })
        return false ;
      },
      getOptions (){
        let self = this;
        axios({
          method :'post',
          url: 'dw/sql/options_yh/date',
          data: this.galaxyForm.date
        }).then(function (response) {
          self.yjkhhyOptions = response.data.yjkhhy_options;
          self.ejkhhyOptions = response.data.ejkhhy_options;
          self.sjkhhyOptions = response.data.sjkhhy_options;
          self.ejywlxOptions = response.data.ejywlx_options;
          self.sjywlxOptions = response.data.sjywlx_options;
          self.qdmcOptions = response.data.qdmc_options;
          self.khlxOptions = response.data.khlx_options;
          self.dqOptions = response.data.dq_options;
          // self.pwOptions = response.data.pw_options;
          // self.yslbOptions = response.data.yslb_options;
          // self.cflxOptions = response.data.cflx_options;
        }).catch(function (error) {
          console.log(error)
        })
      },
      loadHistory (){
        let self = this;
        const hist = JSON.parse(sessionStorage.getItem('sub_history'));
        sessionStorage.removeItem('sub_history');
        if (hist) {
          self.galaxyForm.dimensions = hist['dimensions'];
          self.galaxyForm.date = hist['date'];
          self.galaxyForm.outputs = hist['outputs'];
          const historyForm = hist['filters'];
          historyForm.forEach(function (ele) {
            self.galaxyForm.filters[ele.name].values = ele.values;
            self.galaxyForm.filters[ele.name].type = ele.type;
          })
        }
      }
    },
    created: function () {
      this.getOptions();
    },
    mounted: function () {
      this.loadHistory();
    },
  }
</script>

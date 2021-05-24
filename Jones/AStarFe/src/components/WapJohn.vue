<template>
  <div class="wap">
      <el-form ref="wapForm" :model="wapForm" size="medium" :rules="wapRules" label-width="110px">
      <h3>无线查询维度</h3>
      <el-form-item prop="dimensions" label-width="0px">
        <el-checkbox-group v-model="wapForm.dimensions">
          <el-row>
            <el-col :span="2">
              <el-tag size="medium">基本</el-tag>
            </el-col>
            <el-col :span="22">
              <el-checkbox v-for="(value, key) in optionsBase" :label="value" :key="value">{{key}}</el-checkbox>
            </el-col>
          </el-row>
          <!--<el-row>
            <el-col :span="2">
              <el-tag size="medium">广告</el-tag>
            </el-col>
            <el-col :span="22">
              <el-checkbox v-for="(value, key) in optionsGG" :label="value" :key="key">{{key}}</el-checkbox>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="2">
              <el-tag size="medium">客户</el-tag>
            </el-col>
            <el-col :span="22">
              <el-checkbox v-for="(value, key) in optionsKH" :label="value" :key="value">{{key}}</el-checkbox>
            </el-col>
          </el-row>-->
        </el-checkbox-group>
      </el-form-item>
      <hr>

      <h3>筛选条件</h3>
      <div class="filters">
        <el-row>
          <el-form-item label="起止日期" prop="date">
            <el-date-picker
              v-model="wapForm.date"
              type="daterange"
              align="right"
              unlink-panels
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="yyyyMMdd"
              :picker-options="dateOptions"
              @change="getOptions">
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
            <el-form-item label="查询词" prop="filters.CXC.values">
              <el-input v-model="wapForm.filters.CXC.values"></el-input>
            </el-form-item>
          </el-col>
          <el-col class='type' :span="4">
            <el-form-item label="类型" label-width="50px" prop="filters.CXC.type">
              <el-select v-model="wapForm.filters.CXC.type" placeholder="">
                <el-option label="精确匹配" value="QW_AM"></el-option>
                <el-option label="模糊匹配" value="QW_FM"></el-option>
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
        </el-row>
        <!--other filters-->
        <el-row>
          <el-form-item label="勾选提示">
            <div><span>选项太多找不到？直接搜索试试。</span></div>
          </el-form-item>
        </el-row>

        <el-row>
          <el-col :span="8">
            <el-form-item label="一级查询词行业" prop="filters.CXCHY_YIJI.values">
              <el-select v-model="wapForm.filters.CXCHY_YIJI.values" multiple filterable placeholder="请选择">
                <el-option
                  v-for="value in filterOptions.cxchy_yiji_options"
                  :key="value"
                  :label="value"
                  :value="value">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="二级查询词行业" prop="filters.CXCHY_ERJI.values">
              <el-select v-model="wapForm.filters.CXCHY_ERJI.values" multiple filterable placeholder="请选择">
                <el-option
                  v-for="value in filterOptions.cxchy_erji_options"
                  :key="value"
                  :label="value"
                  :value="value">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>

          <el-col :span="8">
            <el-form-item label="一级渠道" prop="filters.YJQD.values">
              <el-select v-model="wapForm.filters.YJQD.values" multiple filterable placeholder="请选择">
                <el-option
                  v-for="value in filterOptions.yjqd_options"
                  :label="value"
                  :key="value"
                  :value="value">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>

          <el-col :span="8">
            <el-form-item label="二级渠道" prop="filters.EJQD.values">
              <el-select v-model="wapForm.filters.EJQD.values" multiple filterable placeholder="请选择">
                <el-option
                  v-for="value in filterOptions.ejqd_options"
                  :key="value"
                  :label="value"
                  :value="value">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="8" >
            <el-form-item label="网民区域" prop="filters.WMQY.values">
              <el-select v-model="wapForm.filters.WMQY.values" multiple filterable placeholder="请选择">
                <el-option-group
                  v-for="group in filterOptions.wmqy_options"
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
        </el-row>
        <hr>

        <h3>输出指标</h3>
        <el-row>
          <el-form-item label="基础指标" prop="outputs">
            <el-col>
              <el-checkbox-group v-model="wapForm.outputs" @change="changeOrderColumns">
                <el-checkbox v-for="(key, value) in optionsOutputs" :label="key" :key="key">{{value}}</el-checkbox>
              </el-checkbox-group>
            </el-col>
          </el-form-item>
        </el-row>
        <el-row>
          <el-form-item label="输出条数限制" prop="limitNumber">
            <el-col :span="8">
              <el-input v-model="wapForm.limitNumber" placeholder="不填表示没有限制"/>
            </el-col>
          </el-form-item>
        </el-row>
        <el-row>
          <el-form-item label="指标排序" prop="orderColumns">
            <el-select v-model="wapForm.orderColumns" multiple filterable placeholder="请选择">
              <el-option
                v-for="item in orderOptions"
                :label="item.key"
                :key="item.key"
                :value="item.value">
              </el-option>
            </el-select>
            <span style="color:#BDBDBD;font-size:14px;">( 点击的顺序表示排序的指标顺序 )</span>
          </el-form-item>
        </el-row>

        <hr>
        <h3>任务信息</h3>
        <el-row>
          <el-col>
            <el-form-item label="任务名称" prop="name">
              <el-col :span="10">
                <el-input v-model="wapForm.name" placeholder="请输入任务名称"></el-input>
              </el-col>
            </el-form-item>
          </el-col>
        </el-row>

      </div>
    </el-form>
    <div class="formCommit">
      <el-row type="flex" justify="end">
        <el-col :span="6">
<!--          <el-button type="primary" @click="onSubmit('wapForm')">立即提交</el-button>-->
          <el-button type="primary" @click="onSubmit('wapForm')"
                     :icon="buttonIcon" :disabled="onProcessing">{{buttonText}}</el-button>
          <el-button @click="onReset('wapForm')">重置</el-button>
        </el-col>
      </el-row>

    </div>
  </div>

</template>

<script>
  const axios = require('axios')

  import router from '../router'
  import {setCookie, getCookie, delCookie, isEmpty} from '../conf/utils'
  import global_ from './Global'

  export default {
    data () {
      return {
        isOnSubmit: false,
        wapRules: global_.strictRules,
        dateOptions: global_.dateOptions,

        optionsBase: global_.wapBase,
        optionsGG: global_.wapGG,
        optionsKH: global_.wapKH,
        optionsOutputs: global_.wapOutputs,

        // options
        filterOptions: {
          cxchy_yiji_options: [],
          cxchy_erji_options: [],
          wmqy_options: global_.provincesChina,
          yjkhhy_options: [],
          ejkhhy_options: [],
          sjkhhy_options: [],
          khlx_options: [],
          dq_options: [],
          ejywlx_options: [],
          sjywlx_options: [],
          yslb_options: [],
          yjqd_options: [],
          ejqd_options: [],
          wz_options: [],
          pw_options: [],
          sd_options: [],
          gglx_options: [],
          llpt_options: [],
          sjly_options: [],
        },

        // submit form data
        wapForm: {
          date: global_.dateRange,
          dimensions: [],
          filters: {
            DLSID: {
              name: 'DLSID',
              values: [],
              type: 'IN'
            },
            KHID: {
              name: 'KHID',
              values: [],
              type: 'IN'
            },
            CXC: {
              name: 'CXC',
              values: [],
              type: 'QW_AM'
            },
            JHID: {
              name: 'JHID',
              values: [],
              type: 'IN'
            },
            ZID: {
              name: 'ZID',
              values: [],
              type: 'IN'
            },
            PID: {
              name: 'PID',
              values: [],
              type: 'IN'
            },
            GJC: {
              name: 'GJC',
              values: [],
              type: 'IN'
            },
            GJCID: {
              name: 'GJCID',
              values: [],
              type: 'IN'
            },
            CXCHY_YIJI: {
              name: 'CXCHY_YIJI',
              values: [],
              type: 'IN'
            },
            CXCHY_ERJI: {
              name: 'CXCHY_ERJI',
              values: [],
              type: 'IN'
            },
            WMQY: {
              name: 'WMQY',
              values: [],
              type: 'IN'
            },
            YJKHHY: {
              name: 'YJKHHY',
              values: [],
              type: 'IN'
            },
            EJKHHY: {
              name: 'EJKHHY',
              values: [],
              type: 'IN'
            },
            SJKHHY: {
              name: 'SJKHHY',
              values: [],
              type: 'IN'
            },
            KHLX: {
              name: 'KHLX',
              values: [],
              type: 'IN'
            },
            DQ: {
              name: 'DQ',
              values: [],
              type: 'IN'
            },
            EJYWLX: {
              name: 'EJYWLX',
              values: [],
              type: 'IN'
            },
            SJYWLX: {
              name: 'SJYWLX',
              values: [],
              type: 'IN'
            },
            YSLB: {
              name: 'YSLB',
              values: [],
              type: 'IN'
            },
            YJQD: {
              name: 'YJQD',
              values: [],
              type: 'IN'
            },
            EJQD: {
              name: 'EJQD',
              values: [],
              type: 'IN'
            },
            WZ: {
              name: 'WZ',
              values: [],
              type: 'IN'
            },
            PW: {
              name: 'PW',
              values: [],
              type: 'IN'
            },
            SD: {
              name: 'SD',
              values: [],
              type: 'IN'
            },
            GGLX: {
              name: 'GGLX',
              values: [],
              type: 'IN'
            }
          },
          outputs: [],
          limitNumber: '',
          orderColumns:[],
          name: '',
        },

        // response msg
        queryId: '',
        errorMsg: '',
        // 上传的查询词文件
        queryWordFile: '',
        orderOptions: [],   // 排序的列
      }
    },
    methods: {
      isProcessing() {
        return this.isOnSubmit
      },
      changeOrderColumns () {
        let checkOutputs = this.wapForm.outputs
        let outputOptions = this.optionsOutputs

        let options = []
        checkOutputs.forEach(out => {
          for (let key in outputOptions) {
            if (outputOptions[key] === out) {
              options.push({
                'key': key,
                'value': out,
              })
            }
          }
        })
        this.orderOptions = options
      },
      onSubmit (formName) {
        //todo 构造参数,将查询词文件也传上去
        let self = this
        this.$refs[formName].validate((valid) => {
          if (valid) {
            this.$confirm('确认选择了所需的全部条件', '提示', {
              distinguishCancelAndClose: true,
              confirmButtonText: '确定',
              cancelButtonText: '取消',
              type: 'warning',
              center: true
            }).then(() => {
              // build my filter json object
              let postFilters = []
              for (let item in this.wapForm.filters) {
                if (!isEmpty(this.wapForm.filters[item].values)) {
                  if (typeof this.wapForm.filters[item].values === 'string') {
                    const tmp = this.wapForm.filters[item].values
                    this.wapForm.filters[item].values = tmp.split(',')
                  }
                  postFilters.push(this.wapForm.filters[item])
                }
              }

              let formData = new FormData()
              formData.append('typeName','无线_query')
              formData.append('qwFile', this.wapForm.filters.CXC.values)

              let optionOutputsValue = []
              for(let key in this.optionsOutputs) {
                optionOutputsValue.push(this.optionsOutputs[key])
              }
              console.log(optionOutputsValue)
              let outputs = this.wapForm.outputs.sort(function (a,b) {
                return optionOutputsValue.indexOf(a) - optionOutputsValue.indexOf(b)
              })

              formData.append('queryJSON', JSON.stringify({
                dimensions: this.wapForm.dimensions,
                filters: postFilters,
                outputs: outputs,
                date: this.wapForm.date,
                name: this.wapForm.name,
                orderColumns: this.wapForm.orderColumns,
                limitNum: this.wapForm.limitNumber,
              }))

              console.log(formData)
              self.isOnSubmit = true
              axios.post('/jone/combination'
                , formData
                , {headers: {'Content-Type': 'multipart/form-data'}}
              ).then(function (response) {
                self.queryId = response.data
                self.$options.methods.msgNotify.bind(self)()
                self.isOnSubmit = false
              }).catch(function (error) {
                self.errorMsg = error
                self.$options.methods.errorNotify.bind(self)()
                self.isOnSubmit = false
              })
            })
          }
        })
      },
      onReset (formName) {
        this.$refs[formName].resetFields()
      },
      redirectClick () {
        router.push({path: 'que'})
      },
      msgNotify () {
        this.$notify({
          title: this.queryId,
          message: '服务器接收到您的请求并分配查询ID。您可以转到任务队列页面查看结果！',
          type: 'success',
          duration: 2000,
          onClick: this.redirectClick,
        })
      },
      errorNotify () {
        this.$notify.error({
          title: '错误',
          message: this.errorMsg
        })
      },
      beforeUploadDLSID (file) {
        let self = this
        var formData = new FormData()
        formData.append('file', file)
        axios.post('task/file', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
          .then(function (response) {
            self.wapForm.filters.DLSID.values = response.data
          }).catch(function (error) {
          self.errorMsg = error
          self.$options.methods.errorNotify.bind(self)()
        })
        return false // 返回false不会自动上传
      },
      beforeUploadKHID (file) {
        let self = this
        var formData = new FormData()
        formData.append('file', file)
        axios.post('task/file', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
          .then(function (response) {
            self.wapForm.filters.KHID.values = response.data
          }).catch(function (error) {
          self.errorMsg = error
          self.$options.methods.errorNotify.bind(self)()
        })
        return false // 返回false不会自动上传
      },
      beforeUploadCXC (file) {
        let self = this
        this.queryWordFile = file
        var formData = new FormData()
        formData.append('file', file)
        axios.post('task/file', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
          .then(function (response) {
            self.wapForm.filters.CXC.values = response.data
          }).catch(function (error) {
          self.errorMsg = error
          self.$options.methods.errorNotify.bind(self)()
        })
        return false
      },
      beforeUploadJHID (file) {
        let self = this
        var formData = new FormData()
        formData.append('file', file)
        axios.post('task/file', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
          .then(function (response) {
            self.wapForm.filters.JHID.values = response.data
          }).catch(function (error) {
          self.errorMsg = error
          self.$options.methods.errorNotify.bind(self)()
        })
        return false
      },
      beforeUploadZID (file) {
        let self = this
        var formData = new FormData()
        formData.append('file', file)
        axios.post('task/file', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
          .then(function (response) {
            self.wapForm.filters.ZID.values = response.data
          }).catch(function (error) {
          self.errorMsg = error
          self.$options.methods.errorNotify.bind(self)()
        })
        return false
      },
      beforeUploadPID (file) {
        let self = this
        var formData = new FormData()
        formData.append('file', file)
        axios.post('task/file', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
          .then(function (response) {
            self.wapForm.filters.PID.values = response.data
          }).catch(function (error) {
          self.errorMsg = error
          self.$options.methods.errorNotify.bind(self)()
        })
        return false
      },
      beforeUploadGJC (file) {
        let self = this
        var formData = new FormData()
        formData.append('file', file)
        axios.post('task/file', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
          .then(function (response) {
            self.wapForm.filters.GJC.values = response.data
          }).catch(function (error) {
          self.errorMsg = error
          self.$options.methods.errorNotify.bind(self)()
        })
        return false
      },
      beforeUploadGJCID (file) {
        let self = this
        var formData = new FormData()
        formData.append('file', file)
        axios.post('task/file', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
          .then(function (response) {
            self.wapForm.filters.GJCID.values = response.data
          }).catch(function (error) {
          self.errorMsg = error
          self.$options.methods.errorNotify.bind(self)()
        })
        return false
      },
      getOptions () {
        let self = this
        axios({
          method: 'post',
          url: 'dw/sql/options_wap/date',
          data: this.wapForm.date
        }).then(function (response) {
          const options = response.data
          for (var ele in options) {
            self.filterOptions[ele] = options[ele]
          }
        }).catch(function (error) {
          console.log(error)
        })
      },
      loadHistory () {
        let self = this
        const hist = JSON.parse(sessionStorage.getItem('sub_history'))
        sessionStorage.removeItem('sub_history')
        if (hist) {
          try {
            self.wapForm.dimensions = hist['dimensions']
            self.wapForm.date = hist['date']
            self.wapForm.outputs = hist['outputs']
            self.wapForm.limitNumber = hist['limitNum']
            const historyForm = hist['filters']
            self.changeOrderColumns()
            self.wapForm.orderColumns = hist['orderColumns']
            historyForm.forEach(function (ele) {
              self.wapForm.filters[ele.name].values = ele.values
              self.wapForm.filters[ele.name].type = ele.type
            })
          } catch (e) {
            this.onReset('wapForm')
          }
        }
      },
    },
    created: function () {
      this.getOptions()
    },
    computed:{
      onProcessing: function () {
        return this.isProcessing()
      },
      buttonText: function () {
        if (this.isProcessing()) {
          return '提交中'
        } else {
          return '确定'
        }
      },
      buttonIcon: function () {
        if (this.isProcessing()) {
          return 'el-icon-loading'
        } else {
          return ''
        }
      },
    },
    mounted: function () {
      this.loadHistory()
    },
  }
</script>

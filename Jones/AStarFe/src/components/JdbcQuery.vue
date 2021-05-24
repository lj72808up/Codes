<template>
  <div>
    <h3>SQL工单</h3>
    <el-col :span="4">
      <div id="treeDiv" style="height: 700px;overflow:auto">
        <span>库中的模板</span>
        <el-input
          placeholder="输入关键字过滤"
          v-model="filterText">
        </el-input>
        <!--同步加载树-->
        <el-tree :data="treeData" :props="defaultProps" :filter-node-method="filterNode" ref="tree"
                 id="el-tree"
                 @node-click="queryInfo"/>

        <!--右键菜单-->
        <div v-show="canModify">
          <context-menu class="right-menu"
                        :target="contextMenuTarget"
                        :show="menuVisible"
                        @update:show="(show) => menuVisible = show">
            <a href="javascript:" @click="()=>{treeAdd();menuVisible=false}">该节点下创建工单</a>
            <a href="javascript:" @click="onCreateRootDir">根节点下创建新目录</a>
            <a href="javascript:" @click="onDelJdbc">删除节点</a>
            <a href="javascript:" @click="onChangeDir">更改工单所属目录</a>
          </context-menu>
        </div>
      </div>
    </el-col>

    <el-col :span="14" :offset="1">
      <el-form ref="form" :model="form" :rules="rules">
        <el-row>
          <el-form-item prop="name">
            <el-col :span="10">
              <el-input v-model="form.name" placeholder="请输入工单名称"/>
            </el-col>

            <el-col :span="4" :offset="3">
              <!--              <el-select v-model="dsTypeSelectValue" placeholder="请选择数据源类型"-->
              <!--                         v-show="isShow  && (!onlyRead) &&(!idEcho)"-->
              <!--                         @change="onChangeDSType" clearable filterable>-->
              <el-select v-model="dsTypeSelectValue" placeholder="请选择数据源类型"
                         v-show="isShow  && (!onlyRead)"
                         @change="onChangeDSType" clearable filterable>
                <el-option
                  v-for="item in dsTypeOptions"
                  :key="item.name"
                  :label="item.name"
                  :value="item.name">
                </el-option>
              </el-select>
            </el-col>
            <el-col :span="6">
              <!--              <el-select v-model="selectValue" placeholder="请选择工单所属的数据源"-->
              <!--                         v-show="isShow  && (!onlyRead) &&(!idEcho)"-->
              <!--                         @change="onChangeDataSource">-->
              <el-select v-model="selectValue" placeholder="请选择工单所属的数据源"
                         v-show="isShow  && (!onlyRead)"
                         @change="onChangeDataSource">
                <el-option
                  v-for="item in options"
                  :key="item.dsid"
                  :label="item.DatasourceName"
                  :value="item.Dsid">
                </el-option>
              </el-select>
            </el-col>
            <el-col :span="1">
              <div style="margin-left: 5px">
                <span style="color:#BDBDBD" v-show="true">{{ selectValue }}</span>
              </div>
            </el-col>
          </el-form-item>
        </el-row>

        <el-row v-show="(!isShow) && canModify && (!onlyRead)">
          <span style="color:#BDBDBD;font-size:12px;">若工单名称或工单sql已经修改, 请点击旁边的按钮保存</span>
          <el-button @click="onChangeSql('form')" type="Dashed" size="small">保存</el-button>
        </el-row>

        <el-row>
          <el-input v-model="description"
                    v-show="true"
                    placeholder="工单描述"
                    :rows="1"
                    :max-rows="3"
                    :readonly="onlyRead">
            <template slot="prepend">工单备注:</template>
          </el-input>
        </el-row>

        <el-form-item prop="sql">
          <codemirror v-model="form.sql" :options="cmOptions" ref="cmEditor"/>
        </el-form-item>

        <el-form-item>
          <el-col span="16">
            <el-row type="flex" justify="start">
              <el-button @click="onParam('form')">生成参数列表</el-button>
            </el-row>
          </el-col>
          <el-col span="8">
            <el-row type="flex" justify="end">
              <el-checkbox v-show="(isShow && (!onlyRead))" v-model="query_checked">同步创建Redash_Query</el-checkbox>
              <el-button type="primary" @click="onCreateSqlView('form')"
                         v-show="((!isShow) && createFlag && viewAccess)">创建可视化Query
              </el-button>
              <el-button type="success" @click="onJumpSqlView('form')"
                         v-show="((!isShow) && (!createFlag && viewAccess))">跳转可视化
              </el-button>
              <el-button type="primary" @click="onSubmit('form')" v-show="!isShow" :disabled="isOnSubmit">提交执行
              </el-button>
              <el-button @click="onReset('form')">重置</el-button>
              <el-button @click="onCreate('form')" v-show="(isShow && (!onlyRead))" type="primary">新建工单
              </el-button>
            </el-row>
          </el-col>

        </el-form-item>
      </el-form>

    </el-col>

    <el-col :span="4" :offset="1">

      <el-row>
        <label>参数配置</label>
        <el-tooltip class="item" effect="dark" content="参数用于提交的sql中, 不会保存在模板内" placement="top-start">
          <i class="el-icon-question" @click="makeParam()"/>
        </el-tooltip>
        <el-button @click="makeParam()" v-show="false">显示参数</el-button>
      </el-row>

      <template v-for="(item,index) in paramNames">
        <el-row>
          <span>参数名: {{ item }}</span>
        </el-row>
        <el-row>
          <el-select v-model="paramTypes[index]" placeholder="请选择参数类型" @change="onChangeParamType(index)"
                     size="mini" style="width: 100%">
            <el-option
              v-for="item in getParamOptions(item)"
              :key="item.value"
              :label="item.label"
              :value="item.value">
            </el-option>
          </el-select>
        </el-row>
        <el-row>
          <el-input v-model="paramValues[index]" placeholder="请输入"
                    v-show="(paramTypes[index]==='string')||(paramTypes[index]==='timeInterval')" size="mini"/>

          <el-row>
            <el-col span="22">
              <el-input size="mini" v-show="paramTypes[index]==='qw_partition'"
                        placeholder="上传查询词"
                        v-model="cxcValues"/>
            </el-col>
            <el-col span="1" offset="1">
              <el-upload
                class="upload-demo"
                action="123"
                :before-upload="beforeUploadCXC"
                :show-file-list=false
                v-show="paramTypes[index]==='qw_partition'">
                <i class="el-icon-upload2"/>
              </el-upload>
            </el-col>
          </el-row>
          <el-row v-show="paramTypes[index]==='qw_partition'">
            <span style="color:#BDBDBD;font-size:12px;">0:精确匹配,1:模糊匹配</span>
          </el-row>


          <el-date-picker
            size="mini"
            v-model="paramValues[index]"
            v-show="paramTypes[index]==='date'"
            type="date"
            placeholder="选择任务日期"
            format="yyyyMMdd"
            value-format="yyyyMMdd"
            style="width: 100%"
          />
        </el-row>
      </template>


      <el-row>
        <el-divider content-position="center"/>
        <label style="font-size: 15px">账号ID格式化工具</label>
        <el-input type="textarea"
                  v-model="beforeTransId"
                  placeholder="请输入内容 (使用回车分隔)"
                  :rows="8"
                  :max-rows="8">
        </el-input>
        <el-button @click="transId(false)" size="mini">转换</el-button>
        <el-button @click="transId(true)" size="mini">转换字符串</el-button>
        <el-input type="textarea"
                  v-model="afterTransId"
                  :rows="3"
                  :max-rows="3">
        </el-input>
      </el-row>

      <el-dialog title="选择新目录" :visible.sync="newFolderVisible">
        <div id="folderTreeDiv" style="height:500px;;overflow:auto">
          <el-scrollbar style="height: 100%;">
            <el-tree :data="treeData" :props="defaultProps" :filter-node-method="filterNode" ref="folderTree"
                     id="folderTree"
                     @node-click="selectDistFolder"/>
          </el-scrollbar>
        </div>
        <el-row type="flex" justify="end">
          <el-button size="small" @click="submitChangeDir">确认目录更改</el-button>
        </el-row>
      </el-dialog>
    </el-col>
  </div>

</template>

<script>
import router from '../router'

const axios = require('axios')

import {getQueryVariable} from '../conf/utils'
import {getCookie} from '../conf/utils'
import 'codemirror/mode/sql/sql.js'  // 语言模式
import 'codemirror/lib/codemirror.js'
import 'codemirror/mode/clike/clike.js'
import 'codemirror/addon/display/autorefresh.js'
import 'codemirror/addon/edit/matchbrackets.js'
import 'codemirror/addon/selection/active-line.js'
import 'codemirror/addon/display/fullscreen.js'
import 'codemirror/addon/hint/show-hint.js'
import 'codemirror/addon/hint/sql-hint.js'
// theme css
import 'codemirror/addon/hint/show-hint.css'
import 'codemirror/addon/display/fullscreen.css'
import sqlFormatter from 'sql-formatter'

export default {
  data() {
    var validateSql = (rule, value, callback) => {
      let dropExp = /drop/i
      let delExp = /delete/i
      let trunExp = /truncate/i
      let insrtExp = /insert/i
      if (dropExp.test(value) === true
        || delExp.test(value) === true
        || trunExp.test(value) === true
        || insrtExp.test(value) === true
      ) {
        callback(new Error('SQL中存在危险字符!'))
      } else {
        callback()
      }
    }
    return {
      distFolder: 0,
      newFolderVisible: false,
      isOnSubmit: false,
      contextMenuTarget: null,
      cmOptions: {
        // codemirror options
        tabSize: 4,
        mode: 'text/x-mysql',   // https://codemirror.net/mode/index.html
        theme: 'default',
        lineNumbers: true,
        dragDrop: false,
        line: true,
        hintOptions: {
          tables: {},
          databases: {}
        },
        extraKeys: {
          'Tab': 'autocomplete',
          'F12': function (cm) {
            console.log(cm.getOption('fullScreen'))
            cm.setOption('fullScreen', !cm.getOption('fullScreen'))
          },
          'Esc': function (cm) {
            if (cm.getOption('fullScreen')) cm.setOption('fullScreen', false)
          }
        },
        readOnly: 'nocursor',
      },
      cxcValues: [''],
      form: {
        name: '',
        sql: '',
      },
      options: [],
      selectValue: '',
      canDirModify: false,
      menuVisible: false,
      rules: {
        name: [
          {required: true, message: '请输入任务名称', trigger: 'blur'},
          {min: 3, message: '名称长度至少 3 个字符', trigger: 'blur'}
        ],
        sql: [
          {require: true, message: 'SQL语句不能为空', trigger: 'blur'},
          {validator: validateSql, trigger: 'blur'}
        ]
      },
      // isShowRes: false,
      sqlLoading: true,
      sqlQueryId: '',
      sqlData: [],
      sqlDataCols: [],
      treeData: '',
      defaultProps: {
        children: 'children',
        label: 'label',
        templateId: 'templateId',
        redashid: 'redashid',
      },
      filterText: '',
      sqlText: '',
      isShow: true,       // 切换'新建工单'和'点击工单树显示'
      canModify: false,   // 是否可以保存对点击选中的工单的修改 (包括描述)
      onlyRead: false,    // 顶栏标签的配置中, 是否配置了可写权限 (是否有新建工单的权限)
      curTmplateId: '',
      curPid: '',
      description: '',
      paramNames: [],
      paramTypes: [],
      paramValues: [],
      paramOptions: [{
        label: '一般类型',
        value: 'string',
        defaultValue: ''
      }, {
        label: '日期类型',
        value: 'date',
        defaultValue: '19990101'
      }, {
        label: '查询词文件分区',
        value: 'qw_partition',
        defaultValue: ''
      }],
      jsonParams: '',
      beforeTransId: '',
      afterTransId: '',
      queryWordsPartition: '', //'amy_1586947661',
      file: '',
      dataSourceIdNameMapping: {},

      dsTypeOptions: [{'name': 'mysql'}, {'name': 'hive'}, {'name': 'clickhouse'},],
      dsTypeSelectValue: '',
      query_checked: false,
      createFlag: true,
      curRedashQueryId: 0,
      viewAccess: false,
      jumpAccess: false,
      newFolderData: null,
    }
  },
  watch: {
    filterText(val) {
      this.$refs.tree.filter(val)
    },
  },
  computed: {
    codemirror() {
      console.log('计算了')
      return this.$refs.cmEditor.codemirror
    },
  },
  methods: {
    submitChangeDir() {
      let newFolder = this.newFolderData
      let template = this.nodeData
      let params = {
        'templateId': template.templateId + '',
        'folderId': newFolder.templateId + '',
      }
      let self = this
      axios.post('/jdbc/template/ChangeTemplateAffiliation', params, {responseType: 'json'})
        .then(function (response) {
          if (response.data.Res) {
            self.$notify({
              message: '更改成功',
              type: 'success'
            })
            self.newFolderVisible = false
            self.getTreeData()
          } else {
            self.$notify({
              message: response.data.Info,
              type: 'error'
            })
          }
        })
        .catch(function (error) {
          console.log(error)
        })
    },
    selectDistFolder(data) {
      let templateId = this.nodeData.templateId   // 选择的事那个工单的目录
      if (data.pid != 0) {
        this.$notify({
          message: '只能选择一级节点为父目录',
          type: 'error',
          duration: 2000
        })
        this.newFolderData = null
      } else {
        this.$notify({
          message: '已选择 "' + data.label + '" 作为上层目录',
          type: 'success',
          duration: 2000
        })
        this.newFolderData = data
      }
    },
    onChangeDir() {
      let pid = this.nodeData.pid
      console.log('点击了:' + JSON.stringify(this.nodeData))
      if (pid === '0') {
        this.$alert('不能更改目录所属, 只能更改工单所属的目录')
      } else {
        this.menuVisible = false
        this.newFolderVisible = true
      }
    },
    onDelJdbc() {
      let templateId = this.nodeData.templateId
      let self = this
      let params = {
        'templateId': templateId
      }
      this.$confirm('确认永久删除工单吗', {
        distinguishCancelAndClose: true,
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
        center: true
      }).then(() => {
        axios.post('/jdbc/template/deleteTemplateById', params, {responseType: 'json'})
          .then(function (response) {
            if (response.data.Res) {
              self.$notify({
                message: '删除成功',
                type: 'success'
              })
            } else {
              self.$notify({
                message: response.data.Info,
                type: 'error'
              })
            }
            self.getTreeData()
          })
          .catch(function (error) {
            console.log(error)
          })
      })
      this.menuVisible = false
    },
    onCreateRootDir() {
      this.onReset('form')   // 重置 form
      this.menuVisible = false     // 点击后隐藏菜单
      this.$message({   // 提示
        type: 'success',
        message: '已选择根节点作为上级'
      })
    },
    getParamOptions(paramName) {
      let intervalOptions = [
        {
          label: '一般类型.',
          value: 'timeInterval',
          defaultValue: ''
        }
      ]
      // 参数名包含点
      if (paramName.search('\\.') !== -1) {
        return intervalOptions
      } else {
        return this.paramOptions
      }
    },
    onChangeDSType() {
      let self = this
      switch (this.dsTypeSelectValue) {
        /*case 'hive':
          this.selectValue = ''
          this.options = [{'Dsid': 'hive', 'DatasourceName': 'hive'}]
          break*/
        default:
          axios.get('/jdbc/template/getAllDataSources?connType=' + this.dsTypeSelectValue)
            .then(function (response) {
              self.selectValue = ''
              self.options = response.data
            })
            .catch(function (error) {
              self.selectValue = ''
              self.options = []
              console.log(error)
            })
          break
      }
    },
    formatCode() {
      let editor = this.codemirror
      console.log(editor)
      let sqlContent = editor.getValue()
      editor.setValue(sqlFormatter.format(sqlContent))
    },
    beforeUploadCXC(file) {
      let self = this
      var formData = new FormData()
      formData.append('file', file)
      axios.post('task/file', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
        .then(function (response) {
          self.cxcValues = response.data
        }).catch(function (error) {
        self.errorMsg = error
        self.$options.methods.errorNotify.bind(self)()
      })
      return false
    },
    onChangeDataSource() {
      let curDsid = this.selectValue
      let curDsName = ''
      for (let i = 0; i < this.options.length; i++) {
        if (this.options[i]['Dsid'] === curDsid) {
          curDsName = this.options[i]['DatasourceName']
          break
        }
      }
      // console.log('当前数据源名称:' + curDsName)
      // console.log('当前数据源名称:' + curDsName.toLowerCase())
      if ((curDsName.toLowerCase().search('spark') !== -1)) {
        console.log('当前选择的数据源中包含spark')
        // 遍历当前的数据类型中是否有查询词类型,没有则添加
        let flag = false
        for (let i = 0; i < this.paramOptions.length; i++) {
          if (this.paramOptions[i].value === 'qw_partition') {
            flag = true
          }
        }
        if (!flag) {
          this.paramOptions.push({
            label: '查询词文件分区',
            value: 'qw_partition',
            defaultValue: ''
          })
        }
      } else {
        // 遍历当前的数据类型汇总是否含有查询词类型,有则删除
        for (let i = this.paramOptions.length - 1; i >= 0; i--) {
          if (this.paramOptions[i].value === 'qw_partition') {
            this.paramOptions.splice(i, 1)
          }
        }
      }
      // console.log(this.paramOptions)
    },
    handleFileUpload() {
      let index = 0
      for (let i = 0; i < this.paramTypes.length; i++) {
        if (this.paramTypes[i] === 'qw_partition') {
          index = i
        }
      }
      console.log(this.$refs.file[index].files[0])
      this.file = this.$refs.file[index].files[0]
    },
    submitFile() {
      let self = this
      console.log(this.file)
      if (this.file) {
        this.$confirm('请确保提交的文件中只包含查询词一列; 表生成完毕会产生提示, 分区名会在底部显示', '提示', {
          distinguishCancelAndClose: true,
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
          center: true
        }).then(() => {
          let formData = new FormData()
          let username = getCookie('_adtech_user')
          formData.append('file', this.file)
          formData.append('username', username)

          axios.post('/dw/upload/file', formData, {headers: {'Content-Type': 'multipart/form-data'}}
          ).then(function (response) {
            let resData = response.data
            resData = resData.match(/part='(\S*)'/)[1]
            console.log(resData)
            self.queryWordsPartition = resData
            self.$message({
              type: 'success',
              message: '提交成功,分区名为:' + response.data
            })
          })
            .catch(function (error) {
              console.log(error)
              self.$message({
                type: 'success',
                message: '提交失败' + error
              })
            })
        })
      } else {
        // alert('先选择文件')
        this.$confirm('请先选择带查询词的csv文件, 文件中只包含查询词一列;', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
          center: true
        }).then(() => {
        })
      }
    },

    makeParam() {
      let params = []
      for (let i = 0; i < this.paramNames.length; i++) {
        params.push({
          'name': this.paramNames[i],
          'type': this.paramTypes[i],
          'value': this.paramValues[i]
        })
      }
      console.log('paramsName:' + this.paramNames + '\nparamTypes:' + this.paramTypes + '\nparamValues:' + this.paramValues)
      this.jsonParams = JSON.stringify(params)
      console.log('jsonParam:' + this.jsonParams)
    },
    redirectClick() {
      router.push({path: 'que'})
    },
    onCreateSqlView(formName) {
      let self = this
      self.curRedashQueryId = 0

      let params = {
        'templateId': self.curTmplateId,
      }
      axios.post('/jdbc/template/createQueryView', params, {responseType: 'json'})
        .then(function (response) {
          if (response.data.msg.indexOf('成功') >= 0) {
            self.curRedashQueryId = response.data.redash_query_id
            self.createFlag = false
            self.isShow = false
            self.$message({
              type: 'success',
              message: response.data.msg
            })
          } else {
            self.$message({
              type: 'error',
              message: response.data.msg
            })
          }

        })
        .catch(function (error) {
          self.$options.methods.saveUpdateNotify.bind(self)()
          console.log(error)
        })

    },
    onJumpSqlView(formName) {
      let self = this
      console.log('----------------------onJumpSqlView')
      let params = {
        'templateId': self.curTmplateId,
      }
      axios.post('/jdbc/template/getjumpaccess', params, {responseType: 'json'})
        .then(function (response) {
          if (response.data.jumpaccess == true) {
            let queryUrl = 'http://inner-redash.adtech.sogou/queries/' + self.curRedashQueryId + '/source'
            self.$router.push({name: 'redashqueryshow', params: {queryid: self.curRedashQueryId, queryUrl: queryUrl}})
          } else {
            self.$message({
              type: 'error',
              message: '可视化非本人创建,请在自助可视化中跳转!'
            })
          }
          console.info(response.data)
        })
        .catch(function (error) {
          self.$options.methods.saveUpdateNotify.bind(self)()
          console.log(error)
        })

    },
    onChangeSql(formName) {
      let self = this
      let params = {
        // 'templateId': '',
        'title': self.form.name,
        'template': self.form.sql,
        'templateId': self.curTmplateId,
        'description': self.description,
      }
      axios.post('/jdbc/template/updateJdbcTemplate', params, {responseType: 'json'})
        .then(function (response) {
          self.$message({
            type: 'success',
            message: '模板更新成功!'//+self.form.name+':'+this.form.sql+':'+this.curTmplateId
          })
          self.getTreeData()
        })
        .catch(function (error) {
          self.$options.methods.saveUpdateNotify.bind(self)()
          console.log(error)
        })
    },
    onCreate(formName) {
      let self = this
      var username = getCookie('_adtech_user')
      this.$refs[formName].validate((valid) => {
        if (valid) {
          console.log(username)
          console.log(self.query_checked)
          console.log(self.dsTypeSelectValue)
          if (self.query_checked == true && self.dsTypeSelectValue != 'mysql' && self.dsTypeSelectValue != 'clickhouse') {
            this.$alert('当前引擎不支持创建可视化！')
            return false
          }

          this.$confirm('请确认sql语句正确, 是否创建模板?', '提示', {
            distinguishCancelAndClose: true,
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
            center: true
          }).then(() => {
            this.$message({
              type: 'success',
              message: '请求已提交!'
            })
            var params = {
              'creator': username,
              'title': this.form.name,
              'template': this.form.sql,
              'dsid': this.selectValue,
              'pid': parseInt(this.curPid),
              'description': this.description,
              'createredash': this.query_checked,
            }
            console.log(this.query_checked)
            axios({
              method: 'put',
              url: '/jdbc/template/put',
              data: params,
              responseType: 'string'
            }).then(function (response) {
              console.log(response.data)
              // location.reload()
              if (response.data.indexOf('fail') >= 0) {
                self.$message({
                  type: 'error',
                  message: response.data,
                  duration: 5000
                })
              } else {
                self.$message({
                  type: 'success',
                  message: '创建成功!'
                })
              }

              axios.get('/jdbc/template/getAllTemplate')
                .then(function (response) {
                  self.treeData = response.data
                })
                .catch(function (error) {
                  console.log(error)
                })
              self.onReset('form')

            }).catch(function (error) {
              console.log(error)
            })
          })
        } else {
          return false
        }
      })
    },

    onSubmit(formName) {
      let self = this
      //把类型为qw_partition类型的paramValue改成this.queryWordsPartition
      for (let i = 0; i < this.paramNames.length; i++) {
        if (this.paramTypes[i] === 'qw_partition') {
          console.log(i)
          this.paramValues[i] = this.queryWordsPartition
        }
      }
      this.makeParam()
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$confirm('不正确的SQL查询会造成资源浪费，请反复校验查询语句。 是否提交?', '提示', {
            distinguishCancelAndClose: true,
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
            center: true
          }).then(() => {

            console.info('dsid:' + this.selectValue)
            let qws = ''
            if (typeof this.cxcValues == 'string') {
              qws = this.cxcValues
            } else {
              qws = this.cxcValues.join(',')
            }

            let params = {
              'templateId': this.curTmplateId + '',
              'sql': this.form.sql,
              'param': this.jsonParams,
              'dsid': this.selectValue + '',
              'queryName': this.form.name,
              //'queryWordsPartition': this.queryWordsPartition,
              'queryWords': qws
            }
            console.log(params)
            self.isOnSubmit = true
            axios.post('/jdbc/template/run', params, {responseType: 'json'})
              .then(function (response) {
                console.log(response.data)
                self.sqlQueryId = response.data
                self.show = true

                self.onReset('form')
                // self.$options.methods.msgNotify.bind(self)()

                if (response.data.Res) {
                  self.$message({
                    type: 'success',
                    message: '提交成功! ' + '(' + response.data.Info + ')',
                  })
                  self.$options.methods.msgNotify.bind(self)()
                } else {
                  self.$message({
                    type: 'error',
                    message: '提交失败: ' + response.data.Info,
                    duration: 5000
                  })
                }
                self.isOnSubmit = false
              })
              .catch(function (error) {
                self.isOnSubmit = false
                self.$options.methods.errorNotify.bind(self)()
                // self.isShowRes = false
                console.log(error)
                self.$message({
                  type: 'error',
                  message: '提交失败!'
                })
              })
          })
        } else {
          return false
        }
      })
    },
    msgNotify() {
      this.$notify({
        message: '服务器接收到您的请求并分配查询ID。您可以转到任务队列页面查看结果！',
        type: 'success',
        duration: 2000,
        onClick: this.redirectClick,
      })
    },
    onChangeParamType(index) {
      this.paramValues[index] = ''  // 赋一个默认值
      console.log(this.paramTypes)
    },
    resetParam() {
      this.paramTypes = []
      this.paramNames = []
      this.paramValues = []
      this.jsonParams = ''
      this.createFlag = true
      this.curRedashQueryId = 0
    },
    onParam(formName) {
      this.resetParam()
      let sql = this.form.sql
      let expr = /\{\{(.+?)\}\}/g   // 匹配双括号中的参数名
      let originParams = sql.match(expr)
      console.log('originParams:' + originParams)
      originParams.forEach(p => {
        let paramName = p.substring(2, p.length - 2).trim()
        if (this.paramNames.indexOf(paramName) === -1) {  // 如果不存在再插入
          this.paramNames.push(paramName)
        }
      })
      console.log('paramNames:' + this.paramNames)
      for (let step = 0; step < this.paramNames.length; step++) {
        this.paramValues.push('')
        // 参数默认类型的选择
        // (1) 带time的默认 "日期类型"
        // (2) 其它参数默认 "普通类型"
        if (this.paramNames[step].toLowerCase().search('\\.') !== -1) {
          this.paramTypes.push('timeInterval')
        } else if (this.paramNames[step].toLowerCase().search('time') !== -1) {
          this.paramTypes.push('date')
        } else {
          this.paramTypes.push('string')
        }
      }
    },
    onReset(formName) {

      /* Reset our form values */
      this.resetParam()
      this.$refs[formName].resetFields()
      /* Trick to reset/clear native browser form validation state */
      this.$nextTick(() => {
        // this.isShowRes = false
      })
      this.isShow = true
      this.selectValue = ''
      this.curPid = ''
      this.description = ''
    },
    /*onTask () {
      var params = new URLSearchParams()
      params.append('name', this.form.name)
      params.append('sql', this.form.sql)
      params.append('queryId', this.sqlQueryId)
      axios.post('/dw/sql/task', params, {responseType: 'json'})
        .then(function (response) {
          self.sqlDataCols = response.data.cols
          self.sqlData = response.data.data
          self.show = true
        })
        .catch(function (error) {
          self.$options.methods.errorNotify.bind(self)()
          self.isShowRes = false
          console.log(error)
        })
    },*/
    filterNode(value, data) {
      if (!value) return true
      let lowerValue = value.toLowerCase()
      let lowerData = data.label.toLowerCase()
      return lowerData.indexOf(lowerValue) !== -1
    },
    errorNotify() {
      this.$notify.error({
        title: '错误',
        message: 'SQL填写有误。请重新填写后查询！'
      })
    },
    saveUpdateNotify() {
      this.$notify.error({
        title: '错误',
        message: '更新失败'
      })
    },
    loadHistory() {
      let self = this
      const hist = sessionStorage.getItem('jdbc_history')
      sessionStorage.removeItem('jdbc_history')
      if (hist) {
        try {
          let json2map = JSON.parse(hist)
          self.form.sql = json2map['SqlStr']
          self.form.name = json2map['QueryName']
          let params = JSON.parse(json2map['QueryJson'])
          self.selectValue = params['dsid']
          self.isShow = false

          self.jsonParams = params['param']
          let arr = JSON.parse(self.jsonParams)
          for (let i = 0; i < arr.length; i++) {
            self.paramNames.push(arr[i]['name'])
            self.paramTypes.push(arr[i]['type'])
            self.paramValues.push(arr[i]['value'])
          }
          console.log('paramsName:' + this.paramNames + '\nparamTypes:' + this.paramTypes + '\nparamValues:' + this.paramValues)
          console.log('jsonParam:' + this.jsonParams)
        } catch (e) {
          console.log(e)
          this.onReset('form')
        }
      }
    },
    queryInfo(data) {
      console.log(data)
      this.nodeData = data
      let self = this
      if (data.pid !== '0') {
        this.resetParam()
        if (data.redashid != 0) {
          self.curRedashQueryId = data.redashid
          self.createFlag = false
        }

        self.isShow = false
        self.curTmplateId = data.templateId
        console.log('templateId:' + data.templateId)
        console.log('label:' + data.label)
        axios.get('/jdbc/template/getById/' + data.templateId)
          .then(function (response) {
            // console.log(response.data.Template)
            self.form.name = response.data.Title
            self.form.sql = response.data.Template
            self.selectValue = response.data.Dsid
            // self.canModify =
            // console.log('当前数据源id:' + self.selectValue)
            self.onChangeDataSource()
            self.description = response.data.Description
          }).catch(function (error) {
          console.log(error)
        })

        // 获取是否有写权限  (权限都是限制在父级的工单目录上的)
        axios.get('/jdbc/template/canModify/' + data.pid)
          .then(function (response) {
            console.log('当前工单可修改:' + response.data)
            self.canModify = response.data
          }).catch(function (error) {
          console.log(error)
        })
      } else {
        this.resetParam()
        if (data.redashid != 0) {
          self.curRedashQueryId = data.redashid
          self.createFlag = false
        }

        self.isShow = false
        self.curTmplateId = data.templateId
        console.log('templateId:' + data.templateId)
        console.log('label:' + data.label)

        self.form.name = data.label   // 显示目录名
        self.form.sql = ''            // sql置空
        self.description = ''

        // 获取是否有写权限  (权限都是限制在父级的工单目录上的)
        axios.get('/jdbc/template/canModify/' + data.templateId)
          .then(function (response) {
            console.log('当前工单可修改:' + response.data)
            self.canModify = response.data
          }).catch(function (error) {
          console.log(error)
        })

        // 获取是否有写权限  (权限都是限制在父级的工单目录上的)
        axios.get('/jdbc/template/canModify/' + data.templateId)
          .then(function (response) {
            console.log('当前工单目录可修改:' + response.data)
            self.canDirModify = response.data
          }).catch(function (error) {
          console.log(error)
        })
      }

    },

    rightClick(data) {
      console.log(data)
      // this.menuVisible = true
      // this.nodeData = object
      // this.menuVisible = true
      // let menu = document.querySelector('#menu')
      // menu.style.cssText = 'position: fixed; left: ' + (MouseEvent.clientX - 10) + 'px' + '; top: ' + (MouseEvent.clientY - 25) + 'px; z-index: 999; cursor:pointer;'
    },
    treeAdd() {
      if ((this.nodeData.pid + '') != '0') {  // 二级节点
        this.$message({
          type: 'error',
          message: '只能在一级节点下创建模板'
        })
      } else {
        this.onReset('form')
        this.curPid = this.nodeData.templateId
        this.$message({
          type: 'success',
          message: '已选择:"' + this.nodeData.label + '"(' + this.curPid + ')作为上级'
        })
      }
    },

    // 异步树叶子节点懒加载逻辑
    loadNode(node, resolve) {
      // console.log(node, resolve)
      // 一级节点处理
      if (node.level === 0) {
        this.requestTree(resolve)
      }
      // 其余节点处理
      if (node.level >= 1) {
        // 注意！把resolve传到你自己的异步中去
        this.getIndex(node, resolve)
      }
    },
    // 异步加载叶子节点数据函数
    getIndex(node, resolve) {
      axios.get('/jdbc/template/getAllTemplate/' + node.data.templateId)
        .then(function (response) {
          let data = response.data
          resolve(data)
        })
        .catch(function (error) {
          console.log(error)
        })
    },
    // 首次加载一级节点数据函数
    requestTree(resolve) {
      axios.get('/jdbc/template/getAllTemplate/0')
        .then(function (response) {
          let data = response.data
          resolve(data)
        })
        .catch(function (error) {
          console.log(error)
        })
    },
    transId(isString) {
      let before = this.beforeTransId.replace(/(^\s*)|(\s*$)/g, '')  // 删除首尾空格
      let after = ''
      if (isString) {
        after = before.replace(/\r/g, '').replace(/\n/g, '\',\'')
        after = '\'' + after + '\''
      } else {
        after = before.replace(/\r/g, '').replace(/\n/g, ',')
      }
      this.afterTransId = after
    },
    constructContextMenu() {
      // 等 axios 树渲染完毕再构建这个树
      this.$nextTick(() => {
        // vue-context-menu 需要传入一个触发右键事件的元素，等页面 dom 渲染完毕后才可获取
        this.contextMenuTarget = document.querySelector('#el-tree')
        // console.log(this.contextMenuTarget)
        // 获取所有的 treeitem，循环监听右键事件
        const tree = document.querySelectorAll('#el-tree [role="treeitem"]')
        // console.log(tree)
        tree.forEach(i => {
          i.addEventListener('contextmenu', event => {
            // 如果右键了，则模拟点击这个treeitem
            event.target.click()
          })
        })
      })
    },
    getTreeData() {
      let self = this
      axios.get('/jdbc/template/getAllTemplate')
        .then(function (response) {
          self.treeData = response.data
          self.constructContextMenu()
        })
        .catch(function (error) {
          console.log(error)
        })
    }
  },
  mounted() {
    $(document).ready(function () {
      let h = document.body.clientHeight
      let codeMirrorHeight = h - 370 + 'px'
      console.log('codemirror:' + codeMirrorHeight)
      $('.CodeMirror').css('height', codeMirrorHeight)
      let treeDiv = document.getElementById('treeDiv')
      treeDiv.style.height = (h - 150) + 'px'

      /*let folderTreeDiv = document.getElementById('folderTreeDiv')
      folderTreeDiv.style.height = (h - 150) + 'px'*/
    })
    this.loadHistory()
    let self = this
    this.getTreeData()

    axios.get('/jdbc/template/getAllDataSources')
      .then(function (response) {
        self.options = response.data
      })
      .catch(function (error) {
        console.log(error)
      })

    axios.get('/jdbc/template/getEdit')
      .then(function (response) {
        self.onlyRead = response.data.readOnly   // 顶栏标签的配置中, 是否配置了可写权限
        self.cmOptions.readOnly = self.onlyRead ? 'nocursor' : false
      })
      .catch(function (error) {
        console.log(error)
      })

    axios.get('/jdbc/template/getviewaccess')
      .then(function (response) {
        self.viewAccess = response.data.viewaccess   // 顶栏标签的配置中, 是否配置了 可视化权限
        console.log('Access')
        console.log(self.viewAccess)
      })
      .catch(function (error) {
        console.log(error)
      })

  }

}
</script>

<style type='text/css'>
.CodeMirror {
  border: 1px solid #eee;
  /*height: 500px;*/
  font-size: 13px;
  line-height: 120%;
}

html,
body {
  height: 100%;
}

#el-tree {
  user-select: none;
}

.right-menu {
  font-size: 13px;
  position: fixed;
  color: #606266;
  background: #fff;
  border: solid 1px rgba(0, 0, 0, .2);
  border-radius: 3px;
  z-index: 999;
  display: none;
}

.right-menu a {
  width: 150px;
  height: 28px;
  line-height: 28px;
  text-align: center;
  display: block;
  color: #1a1a1a;
}

.right-menu a:hover {
  background: #eee;
  color: #fff;
}

.right-menu {
  border: 1px solid #eee;
  box-shadow: 0 0.5em 1em 0 rgba(0, 0, 0, .1);
  border-radius: 1px;
  height: 130px;
}

a {
  text-decoration: none;
}

.right-menu a {
  padding: 2px;
}

.right-menu a:hover {
  background: #99A9BF;
}

.el-scrollbar__bar {
  right: 0px;
  bottom: 0px;
}
</style>



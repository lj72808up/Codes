<template>
  <div>
    <h3>SQL查询API</h3>
    <el-col :span=19 style="padding-right: 5px">
      <div id="box">
        <div id="top">
          <el-form ref="form" :model="form" :rules="rules">
            <el-row style="margin-bottom: -10px">
              <el-form-item prop="name">
                <el-col :span=6>
                  <el-input v-model="form.name" size="small" placeholder="请输入任务名称"/>
                </el-col>
                <el-col :span=4>
                  <div style="margin-left: 5px">
                    <el-select v-model="datasource.typeSelectValue" placeholder="请选择数据源类型" size="small"
                               @change="onChangeDSType" clearable filterable>
                      <el-option
                        v-for="item in datasource.typeOptions"
                        :key="item.name"
                        :label="item.name"
                        :value="item.name">
                      </el-option>
                    </el-select>
                  </div>
                </el-col>
                <el-col :span=7 :offset="7">
                  <el-select v-model="datasource.selectValue" placeholder="请选择数据源" size="small"
                             @change="onChangeDataSource" clearable filterable>
                    <el-option
                      v-for="item in datasource.options"
                      :key="item.dsid"
                      :label="item.DatasourceName"
                      :value="item.Dsid">
                    </el-option>
                  </el-select>
                  <el-button type="primary" size="small" @click="onClickSave()">保存</el-button>

                </el-col>
              </el-form-item>
            </el-row>
            <el-form-item prop="sql">
              <codemirror v-model="form.sql" :options="cmOptions" ref="cmEditor"/>
            </el-form-item>
          </el-form>

          <!--文件上传-->
          <el-row :gutter="20" style="margin-top: -7px;">
            <el-col :span=14>
              <el-row type="flex" justify="start">
                <div>
                  <el-button size="small" v-show="datasource.typeSelectValue!=='clickhouse'" @click="verifySql"
                             :disabled="onVerify" style="margin-right: 13px">
                    {{verifyButtonText}}
                  </el-button>
                </div>

                <el-upload
                  class="upload-demo"
                  ref="upload"
                  action="dummyUrl"
                  :limit="1"
                  :file-list="fileList"
                  :auto-upload="false"
                  :before-upload="beforeUpload"
                  :on-exceed="onExceed"
                >
                  <el-button slot="trigger" size="small" type="primary" v-if="visibility">选取查询词文件</el-button>
                  <el-button style="margin-left: 10px;" size="small" type="success" @click="submitFile"
                             v-if="visibility">上传到服务器
                  </el-button>
                  <span style="color:#BDBDBD;font-size:12px;" v-show="true">{{uploadTableName}}</span>
                </el-upload>
              </el-row>
            </el-col>
            <el-col :span=10>
              <el-button size="small" @click="formatCode()">格式化</el-button>
              <el-button size="small" @click="onReset('form')">重置</el-button>
              <el-button size="small" type="primary" @click="onSubmit('form')" :disabled="isOnSubmit">提交</el-button>
              <el-button size="small"  v-show="(viewAccess)" type="primary" @click="onCreateSqlView()" :disabled="disableRedashCreteBtn"  >
                生成可视化
              </el-button>

            </el-col>
          </el-row>
        </div> <!--end top-->

        <div id="resize"> <!--class="line"-->
          <el-divider><i class="el-icon-d-caret"/></el-divider>
        </div>

        <div id="down">
          <span style="display: block; font-weight:bold;">表预览 {{explorerTable.tableName}}</span>
          <el-tabs tab-position="top" @tab-click="onClickTab" type="border-card">
            <el-tab-pane label="字段">
              <!--              <div style="max-height: 500px; overflow: auto">-->
              <!--字段标签页-->
              <el-table
                v-loading="fieldLoading"
                :data="fieldList"
                stripe
                style="width: 80%"
                :default-sort="{prop:'', order:'descending'}"
                @cell-dblclick="onRightClickFields">
                <el-table-column key="fieldName" prop="fieldName" label="字段名称" sortable="true"/>
                <el-table-column key="fieldType" prop="fieldType" label="字段类型" sortable="true"/>
                <el-table-column key="description" prop="description" label="描述" sortable="true"/>
              </el-table>

              <div v-show="partitionDataList.length>0">
                <h4>分区字段</h4>
                <el-col span=15>
                  <el-table
                    :data="partitionDataList"
                    stripe
                    style="width: 100%"
                    :default-sort="{prop:'', order:'descending'}">
                    <el-table-column key="fieldName" prop="fieldName" label="字段名称" sortable="true"/>
                    <el-table-column key="fieldType" prop="fieldType" label="字段类型" sortable="true"/>
                  </el-table>
                </el-col>
              </div>
              <!--              </div>-->

            </el-tab-pane>

            <el-tab-pane label="表信息">
              <el-form ref="tableForm" :model="tableForm" label-width="120px">
                <el-row>
                  <el-col :span=7>
                    <el-form-item label="表名称 : " size="small">
                      <el-input v-model="tableForm.tableName" :readonly="true"/>
                      <!--                    <span>{{tableForm.tableName}}</span>-->
                    </el-form-item>
                  </el-col>

                  <el-col :span=7>
                    <el-form-item label="所属库 : " size="small">
                      <el-input v-model="tableForm.database" :readonly="true"/>
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-row>
                  <el-col span=7>
                    <el-form-item label="创建者 : " size="small">
                      <el-input v-model="tableForm.creator" :readonly="true"/>
                    </el-form-item>
                  </el-col>

                  <el-col span=7>
                    <el-form-item label="创建时间 : " size="small">
                      <el-input v-model="tableForm.createTime" :readonly="true"/>
                    </el-form-item>
                  </el-col>

                  <el-col :span=7>
                    <el-form-item label="表类型 : " size="small">
                      <el-input v-model="tableForm.tableType" :readonly="true"/>
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-row>
                  <el-col span=14>
                    <el-form-item label="文件存储位置 : ">
                      <el-input v-model="tableForm.location" :readonly="true"/>
                    </el-form-item>
                  </el-col>
                  <el-col span=7>
                    <el-form-item label="压缩方式 : ">
                      <el-input v-model="tableForm.compressMode" :readonly="true"/>
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-row>
                  <el-col span=21>
                    <el-form-item label="描述 : ">
                      <el-input v-model="tableForm.description" :readonly="true"/>
                    </el-form-item>
                  </el-col>
                </el-row>

              </el-form>
              <!--              </div>-->
            </el-tab-pane>
            <el-tab-pane label="数据预览">
              <el-table
                v-loading="exploreData.loading"
                element-loading-text="拼命加载中"
                element-loading-spinner="el-icon-loading"
                :data="exploreData.dataList"
                stripe
                style="width: 100%"
                :default-sort="{prop:'', order:'descending'}">
                <el-table-column
                  v-for="col in exploreData.tableColumn"
                  :key="col.field"
                  :prop="col.field"
                  :label="col.title"
                  :sortable="true">
                </el-table-column>
              </el-table>
            </el-tab-pane>


            <el-tab-pane label="历史任务">
              <el-table
                v-loading="taskLoading"
                :data="taskData"
                stripe
                style="width: 100%"
                :default-sort="{prop:'last_update', order:'descending'}">
                <el-table-column
                  v-for="taskCol in taskCols"
                  v-if="!taskCol.hidden"
                  :key="taskCol.prop"
                  :prop="taskCol.prop"
                  :label="taskCol.label"
                  :sortable="taskCol.sortable">
                </el-table-column>
                <el-table-column
                  fixed="right"
                  label="异常">
                  <template slot-scope="scope">
                    <el-button @click="checkHistoryMessage(scope.row)" type="text" size="mini" v-show="scope.row.status_code==4">查看异常
                    </el-button>
                  </template>
                </el-table-column>
                <el-table-column
                  fixed="right"
                  label="操作"
                  width="210">
                  <template slot-scope="scope">
                    <el-button @click="pathClick(scope.row)" type="text" size="mini">下载结果</el-button>
                    <el-button @click="redirectHistoryTaskClick(scope.row)" type="text" size="mini">查看任务SQL</el-button>
                  </template>
                </el-table-column>
              </el-table>

            </el-tab-pane>


            <!-- 提交过的任务结果查看 -->
            <el-tab-pane
              v-for="(task, index) in exploreResult.tasks"
              :key="task.id"
              :label="task.name"
              :name="task.id">

              <el-form ref="curTaskRes" :model="curTaskRes">
                <el-form-item label="任务id:" style="margin-bottom: 1px;">
                  {{task.id}}
                  <el-button type="primary" size="small" @click="removeTab()" style="margin-left: 20px">关闭</el-button>
                </el-form-item>
                <el-form-item label="状态:" style="margin-bottom: 1px;">
                  {{curTaskRes.status}}

                  <el-button v-show="curTaskRes.status==='任务取消'"
                             style="margin-left: 30px" size="mini" type="danger"
                             @click="checkMessage()">查看错误
                  </el-button>
                  <el-button v-show="curTaskRes.status==='任务完成'"
                             style="margin-left: 30px" size="small" type="primary"
                             @click="downloadTaskResult()">下载结果
                  </el-button>
                </el-form-item>
              </el-form>
              <el-progress :text-inside="true" :stroke-width="24" :percentage="curTaskRes.process" status="success">
              </el-progress>
              <!--任务结果预览5000条-->
              <pl-table ref="plTable"
                        selectTrClass="selectTr"
                        :max-height="500"
                        :use-virtual="true"
                        header-drag-style
                        fixedColumnsRoll
                        inverse-current-row
                        v-loading="curTaskRes.dataLoading">
                <template slot="empty">
                  没有查询到符合条件的记录
                </template>
                <pl-table-column
                  v-for="item in curTaskRes.contentCols"
                  :key="item.id"
                  :resizable="true"
                  :show-overflow-tooltip="true"
                  :prop="item.prop"
                  :label="item.label"
                  :fixed="item.fixed"/>
              </pl-table>

            </el-tab-pane>
          </el-tabs>
        </div>

      </div>
    </el-col>
    <el-col :span=5>
      <div id="treeDiv">
        <el-tabs v-model="activeRightTabs">
          <el-tab-pane label="库表结构" name="first">
            库表结构
            <el-input
              placeholder="输入关键字过滤"
              size="small"
              v-model="filterText"/>
            <el-tree :data="treeData" :props="defaultProps" :filter-node-method="filterNode" ref="tree">
              <span class="el-tree-node__content" slot-scope="{ node, data }">
                <span @click="onClickTree(data)" @dblclick="rightClickTree(data)">{{ node.label }}</span>
              </span>
            </el-tree>
          </el-tab-pane>
          <el-tab-pane label="字段索引" name="second">字段索引
            <el-row>
              <el-col :span=19 style="padding-right: -10px">
                <el-input placeholder="查询包含该字段的表" size="small" v-model="searchFieldName"
                          @keydown.enter.native="searchField"/>
              </el-col>
              <el-col :span=5>
                <el-button size="small" type="primary" @click="searchField">查询</el-button>
              </el-col>
            </el-row>
            <el-tree :data="fieldIndexData" :props="defaultProps" :filter-node-method="filterNode" ref="indexFieldTree">
              <span class="el-tree-node__content" slot-scope="{ node, data }">
                <span @click="onClickSearchFieldTree(data)" @dblclick="rightClickTree(data)">{{ node.label }}</span>
              </span>
            </el-tree>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-col>

    <el-dialog title="保存至工单" :visible.sync="saveOption.showDialog" @close="closeDialog()">
      <span style="color:red;font-size:12px;">{{nodeData.label + '  -  ' + saveOption.wsName}}</span>
      <el-row style="margin-bottom: 10px">
        <el-col :span=10>
          <el-input placeholder="输入存储名称" v-model="saveOption.wsName" size="small"/>
        </el-col>
        <el-col :span=2 :offset="1">
          <el-button size="mini" type="primary" @click="submitNewWs()">提交</el-button>
        </el-col>
      </el-row>
      <div style="height:500px; overflow: auto">
        <el-tree :data="wsTree" :props="defaultProps" :filter-node-method="filterNode" ref="wsTree" accordion
                 @node-click="onClickWSTree"/>
        <!--  @node-contextmenu="ondbClickWSTree"   -->

        <!--增加右键事件 @mouseleave="menuVisible=!menuVisible"-->
        <div v-show="menuVisible" @mouseleave="menuVisible=!menuVisible">
          <button id="menu" class="el-button el-button--small" @click="treeAdd">在该节点下创建模板</button>
          <!--<ul id="menu2">
            <li tabindex="-1" style="list-style-type:none" @click="treeAdd">
              <el-link icon="el-icon-circle-plus-outline" :underline="false">在该节点下创建模板</el-link>
            </li>
          </ul>-->
        </div>
      </div>
    </el-dialog>

    <el-dialog
      title="异常信息"
      :visible.sync="curTaskRes.errDialogVisible"
      width="50%">
      <el-input type="textarea"
                v-model="curTaskRes.errMsg"
                placeholder="无异常内容显示"
                :rows="20"
                :max-rows="30"
                :readonly="true"/>
    </el-dialog>

    <el-dialog
      title="异常信息"
      :visible.sync="historyDialogVisible"
      width="50%">
      <el-input type="textarea"
                v-model="historyErrMsg"
                placeholder="无异常内容显示"
                :rows="20"
                :max-rows="30"
                :readonly="true"></el-input>
    </el-dialog>
  </div>

</template>

<script>
  import sqlFormatter from 'sql-formatter'
  import 'codemirror/mode/sql/sql.js' // 语言模式
  import 'codemirror/lib/codemirror.js'
  import 'codemirror/mode/clike/clike.js'
  import 'codemirror/addon/display/autorefresh.js'
  import 'codemirror/addon/edit/matchbrackets.js'
  import 'codemirror/addon/selection/active-line.js'
  import 'codemirror/addon/display/fullscreen.js'
  import 'codemirror/addon/hint/show-hint.js'
  import 'codemirror/addon/hint/sql-hint.js'
  // theme css
  import 'codemirror/theme/base16-dark.css'
  import 'codemirror/addon/hint/show-hint.css'
  import 'codemirror/addon/display/fullscreen.css'
  import router from '../router'
  import {getCookie, setCookie, formatDateLine, isRunFinsh, isRunFinshPagani, type2route} from '../conf/utils'


  const axios = require('axios')
  let time = 50
  let timeOut = null

  export default {
    data () {
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
        onVerify: false,
        searchFieldDataSource: {
          dsType: '',
          dsId: '',
        },
        dataPathRoot: '',
        dataPathRootPaga: '',
        curTaskRes: {
          taskId: '',
          status: '',
          statusCode: '',
          errMsg: '',
          errDialogVisible: false,
          dataPath: '',
          filePath: '',
          dataLoading: false,
          contentCols: [],
          contentData: [],
          process: 0,
          tableIdx: 0,
        },
        // 结果标签预览
        exploreResult: {
          maxRunningTask: 3,   // 最大任务执行数量
          runningTaskNum: 0,   // 正在执行的任务个数
          taskTotalCnt: 0,     // 全部提交过的任务数
          tasks: [],         // 全部提交过的任务id集合  [{id: 'c1fac271227806877913ccf196e69eb7', name: '再来一次'}]
        },
        isOnSubmit: false,
        visibility: false,
        nodeData: {
          label: ''
        },
        datasource: {
          selectValue: '',
          options: [],
          mysqlOptions: [],
          typeOptions: [{'name': 'mysql'}, {'name': 'hive'}, {'name': 'clickhouse'},{'name': 'kylin'}],
          typeSelectValue: '',
        },
        wsTree: '',
        activeRightTabs: 'first',
        menuVisible: false,
        saveOption: {
          showDialog: false,
          wsShowName: '',
          wsName: '',
        },
        tableForm: {
          tableName: '',
          tableType: '',
          database: '',
          location: '',
          compressMode: '',
          description: '',
          tid: '',
          creator: '',
          createTime: ''
        },
        exploreData: {
          dataList: [],
          tableColumn: [],
          loading: false
        },
        searchFieldName: '',
        form: {
          name: '',
          sql: '',
        },
        fieldLoading: false,  //字段预览加载
        fieldList: [],
        tableColumn: [],
        partitionDataList: [],
        curTab: '',
        explorerTable: {
          tableId: '',
          tableName: '',
        },
        uploadTableName: '',
        rules: {
          sql: [
            {require: true, message: 'SQL语句不能为空', trigger: 'blur'},
            {validator: validateSql, trigger: 'blur'}
          ]
        },
        isShowRes: false,
        sqlLoading: true,
        sqlQueryId: '',
        sqlData: [],
        sqlDataCols: [],
        treeData: '',
        fieldIndexData: '',
        defaultProps: {
          children: 'children',
          label: 'label'
        },
        curPid: '',
        filterText: '',
        wsFilterText: '',
        fileList: [],
        queryWordsPartition: '',
        exploreTableName: '',

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
          }
        },

        timer: null,
        curRedashQueryId: "",
        disableRedashCreteBtn: false,
        viewAccess: false,

        taskLoading: false,
        taskCols: [],
        taskData: [],
        taskForm: {
          session: '',
          date: formatDateLine(new Date()),
          startDate: formatDateLine(new Date(new Date().getTime() - 7 * 24 * 60 * 60 * 1000)),
          status: ''
        },

        historyErrMsg: '',
        historyDialogVisible: false,
      }
    },
    created () {
      let self = this
      axios.get('/task/log/rsync')
        .then(function (response) {
          self.dataPathRoot = response.data[0]
          self.dataPathRootPaga = response.data[1]
        }).catch(function (error) {
          console.log(error)
        }
      )
    },
    watch: {
      filterText (val) {
        this.$refs.tree.filter(val)
      },
      wsFilterText (val) {
        this.$refs.wsTree.filter(val)
      },
    },
    methods: {
      searchField () {
        let params = {
          'fieldName': this.searchFieldName
        }
        let self = this
        axios.post('/dw/searchField', params, {responseType: 'json'})
          .then(function (response) {
            self.fieldIndexData = response.data.children
          }).catch(function (error) {
          console.error(error)
        })
      },
      verifySql () {
        let sql = this.form.sql
        let sqlType = this.datasource.typeSelectValue
        let dsId = this.datasource.selectValue
        let self = this
        let params = {
          'sql': sql,
          'sqlType': sqlType,
          'dsId': dsId + '',
        }
        this.onVerify = true
        axios.post('/dw/sqlVerify', params, {responseType: 'json'})
          .then(function (response) {
            console.log(response.data)
            let data = response.data
            self.onVerify = false
            if (data.Res) {
              self.$message({
                message: 'sql校验通过',
                type: 'success'
              })
            } else {
              self.$alert(data.Info, 'sql校验失败', {
                confirmButtonText: '确定'
              })
            }
          }).catch(function (error) {
          self.onVerify = false
          console.error(error)
        })

      },
      removeTab () {
        let curTaskId = this.curTaskRes.taskId
        let tasks = this.exploreResult.tasks
        let idx = -1
        for (let i = 0; i < tasks.length; i++) {
          if (tasks[i].id === curTaskId) {
            idx = i
          }
        }
        console.log('关闭:' + curTaskId + '; id:' + idx)
        tasks.splice(idx, 1)
        this.exploreResult.tasks = tasks
        this.curTaskRes.taskId = '';
      },
      checkHistoryMessage (row) {
        this.historyErrMsg = ''
        this.historyDialogVisible = true
        let self = this
        let params = {
          'queryId': row.query_id
        }
        axios.post('/task/checkException', params, {responseType: 'json'})
          .then(function (response) {
            self.historyErrMsg = response.data.res
          }).catch(function (error) {
          console.error(error)
        })
      },
      checkMessage () {
        this.curTaskRes.errMsg = ''
        this.curTaskRes.errDialogVisible = true
        let self = this
        let params = {
          'queryId': this.curTaskRes.taskId
        }
        axios.post('/task/checkException', params, {responseType: 'json'})
          .then(function (response) {
            self.curTaskRes.errMsg = response.data.res
          }).catch(function (error) {
          console.error(error)
        })
      },
      ProcessBarIncrease () {
        this.percentage += 10
        if (this.percentage > 100) {
          this.percentage = 100
        }
      },
      downloadTaskResult () {
        window.open(this.curTaskRes.filePath + '?op=OPEN')
      },
      onRightClickFields (row, column, cell, event) {
        let checkWord = cell.children[0].innerHTML  // 双击选中的文字
        this.insertCodeMirror(checkWord)
      },
      onChangeDSType () {
        let self = this
        if (this.datasource.typeSelectValue === 'hive') {
          this.visibility = true
        } else {
          this.visibility = false
        }
        this.treeData = ''
        axios.get('/dw/GetConnByPrivilege/' + this.datasource.typeSelectValue)
          .then(function (response) {
            self.datasource.selectValue = ''
            self.datasource.options = response.data
            // 如下两行,默认选择数据源的第一个选项,并更新树
            self.datasource.selectValue = self.datasource.options[0].Dsid
            self.onChangeDataSource()
          })
          .catch(function (error) {
            self.datasource.selectValue = ''
            self.datasource.options = []
            console.log(error)
          })
      },
      closeDialog () {
        // 重置表中的数据
        this.saveOption.wsName = '' // 输入名称
        this.curPid = '0'// pid
        this.nodeData.label = ''// 父节点名字
      },
      onChangeDataSource () {
        let curdsType = this.datasource.typeSelectValue
        let curdsId = this.datasource.selectValue
        let self = this
        console.log(curdsType + ':' + curdsId)
        let params = {
          'dsType': curdsType,
          'dsId': curdsId + '',
        }
        // 1. 元数据树
        axios.post('/dw/meta/get', params, {responseType: 'json'})
          .then(function (response) {
            self.treeData = response.data.children
          })
          .catch(function (error) {
            console.log(error)
          })

        this.cmOptions.hintOptions.tables = {}
        // 2. sql表和字段的输入提示
        axios.post('/dw/hint/getSqlHint', params, {responseType: 'json'})
          .then(function (response) {
            self.cmOptions.hintOptions.tables = response.data
            console.log(self.cmOptions.hintOptions.tables)
            console.log('this is current codemirror object', self.codemirror)
            self.result = 'finish'
          })
          .catch(function (error) {
            console.log(error)
          })

        // 3. udf函数提示  (hive模式下)
        if (curdsType.indexOf('hive') !== -1) {
          axios.get('/dw/hint/getSqlUDF')
            .then(function (response) {
              let funcs = response.data
              if (funcs != null && funcs.length > 0) {
                funcs.forEach(ele => {
                  self.cmOptions.hintOptions.tables[ele] = []
                })
              }
            })
            .catch(function (error) {
              console.log(error)
            })
        }
        console.log(self.cmOptions.hintOptions)
      },
      submitNewWs () {
        let params = {
          'creator': getCookie('_adtech_user'),
          'title': this.saveOption.wsName,
          'template': this.form.sql,
          'dsid': this.datasource.selectValue,
          'pid': parseInt(this.curPid),
          'description': '',
        }
        let self = this
        axios({
          method: 'put',
          url: '/jdbc/template/put',
          data: params,
          responseType: 'string'
        }).then(function (response) {
          self.saveOption.showDialog = false
          self.$message({
            type: 'success',
            message: '模板保存成功!'
          })
          console.log(response.data)
          // location.reload()
        }).catch(function (error) {
          console.log(error)
          self.$message({
            type: 'error',
            message: '模板保存失败!' + error
          })
        })
      },
      treeAdd (event) {
        console.log(this.nodeData)
        if ((this.nodeData.pid + '') !== '0') {  // 二级节点
          this.$message({
            type: 'error',
            message: '只能在一级节点下创建模板'
          })
        } else {   // 一级节点
          this.curPid = this.nodeData.templateId
          this.$message({
            type: 'success',
            message: '已选择:"' + this.nodeData.label + '"(' + this.curPid + ')作为上级'
          })
          this.saveOption.wsShowName = this.nodeData.label + ' - ' + this.saveOption.wsName
        }
      },
      onClickWSTree (object, node, component) {
        // pid != 0的是二级节点, 不能选中
        if (object.pid === '0') {
          this.idEcho = false
          this.nodeData = object
          this.curPid = this.nodeData.templateId
          this.saveOption.wsShowName = this.nodeData.label + ' - ' + this.saveOption.wsName
        }

        console.log(object)
      },
      ondbClickWSTree (MouseEvent, object, Node, element) {
        this.idEcho = false
        this.nodeData = object
        this.menuVisible = true
        let menu = document.querySelector('#menu')
        console.log('鼠标位置:' + MouseEvent.clientX + ',' + MouseEvent.clientY)
        menu.style.cssText = 'position: fixed; left: ' + (MouseEvent.clientX - 280) + 'px' + '; top: ' + (MouseEvent.clientY - 105) + 'px; z-index: 999; cursor:pointer;'
      },
      onClickSave () {
        this.saveOption.showDialog = true
        let self = this
        axios.get('/jdbc/template/getAllTemplate')
          .then(function (response) {
            self.wsTree = response.data
          })
          .catch(function (error) {
            console.log(error)
          })
      },
      onExploreData (tableName, dsType, dsId) {
        console.log('explore data的参数: dsType:' + dsType + ',dsId:' + dsId)
        let self = this
        this.exploreData.dataList = []
        let sendRequest = async () => {
          try {
            console.log('开始获取' + tableName + '的数据')
            self.exploreData.loading = true
            let params = {
              'tableName': tableName,
              'dsType': (dsType != null && dsType !== '') ? dsType : this.datasource.typeSelectValue,
              'dsId': (dsId != null && dsId !== '') ? dsId : this.datasource.selectValue + ''
            }
            const resp = await axios.post('/dw/explore', params, {responseType: 'json'})
            self.exploreData.loading = false
            self.exploreData.dataList = resp.data.dataList
            self.exploreData.tableColumn = resp.data.tableColumn
          } catch (err) {
            console.error(err)
          }
        }
        sendRequest()
      },
      insertCodeMirror (data) {
        let editor = this.$refs.cmEditor.codemirror
        let doc = editor.getDoc()
        let cursor = doc.getCursor()// 获取输入光标
        let line = cursor.line  // 第几行
        let ch = cursor.ch        // 第几个字读
        console.log(cursor)
        let pos = {
          line: line,
          ch: ch
        }
        doc.replaceRange(data, pos)
        editor.focus()
      },
      rightClickTree (data) {
        clearTimeout(timeOut)
        console.log('双击事件')
        this.insertCodeMirror(data.tableName)
      },
      onClickSearchFieldTree (node) {
        let self = this
        console.log(node)
        if (node.tableName !== '') {
          let dsType = node.tableType
          let dsId = node.dsid + ''
          self.searchFieldDataSource.dsType = dsType
          self.searchFieldDataSource.dsId = dsId
          self.explorerTable.tableName = node.label
          self.explorerTable.tableId = node.tid

          // 不要联动数据源类型和数据源, 这两个是任务的数据源
          //数据预览
          if (self.curTab === '数据预览') {
            self.onExploreData(node.tableName, dsType, dsId)
          }
          self.exploreTableName = node.tableName
          // 请求字段信息
          let params = {
            'tableId': self.explorerTable.tableId + '',
            'version': 'latest',
          }
          axios.post('/metaData/fields/getByVersion', params, {responseType: 'json'})
            .then(function (response) {
              self.tableColumn = response.data.tableColumn
              self.fieldList = response.data.dataList
              // 分区信息
              axios.post('/metaData/fields/getPartitionByVersion', params, {responseType: 'json'})
                .then(function (response) {
                  self.partitionDataList = response.data
                })
                .catch(function (error) {
                  console.log(error)
                })
            })
            .catch(function (error) {
              console.log(error)
            })

          // 表信息
          let tid = node.tid
          axios.get('/metaData/getTableById/' + tid, {responseType: 'json'})
            .then(function (response) {
              console.log(response.data)
              self.tableForm.tableName = response.data.TableName
              self.tableForm.tableType = response.data.TableType
              self.tableForm.database = response.data.DatabaseName
              self.tableForm.location = response.data.Location
              self.tableForm.compressMode = response.data.CompressMode
              self.tableForm.description = response.data.Description
              self.tableForm.creator = response.data.Creator
              self.tableForm.createTime = response.data.CreateTime
            })
            .catch(function (error) {
              console.log(error)
            })
        }
      },
      onClickTree (node) {
        // clearTimeout(timeOut) // 清除第一个单击事件
        let self = this
        // timeOut = setTimeout(function () {
        console.log('单击事件')
        console.log(node)
        if (node.tableName !== '') {
          self.explorerTable.tableName = node.label
          self.explorerTable.tableId = node.tid
          //数据预览
          if (self.curTab === '数据预览') {
            self.onExploreData(node.tableName)
          }
          self.exploreTableName = node.tableName
          // 请求字段信息
          let params = {
            'tableId': self.explorerTable.tableId + '',
            'version': 'latest',
          }
          axios.post('/metaData/fields/getByVersion', params, {responseType: 'json'})
            .then(function (response) {
              self.tableColumn = response.data.tableColumn
              self.fieldList = response.data.dataList
              // 分区信息
              axios.post('/metaData/fields/getPartitionByVersion', params, {responseType: 'json'})
                .then(function (response) {
                  self.partitionDataList = response.data
                })
                .catch(function (error) {
                  console.log(error)
                })
            })
            .catch(function (error) {
              console.log(error)
            })

          // 表信息
          let tid = node.tid
          axios.get('/metaData/getTableById/' + tid, {responseType: 'json'})
            .then(function (response) {
              console.log(response.data)
              self.tableForm.tableName = response.data.TableName
              self.tableForm.tableType = response.data.TableType
              self.tableForm.database = response.data.DatabaseName
              self.tableForm.location = response.data.Location
              self.tableForm.compressMode = response.data.CompressMode
              self.tableForm.description = response.data.Description
              self.tableForm.creator = response.data.Creator
              self.tableForm.createTime = response.data.CreateTime
            })
            .catch(function (error) {
              console.log(error)
            })
        }
      },
      getExploreTaskData (tableIndex) {
        let self = this
        console.log("getExploreTaskData")
        console.log(tableIndex)
        console.log(this.curTaskRes.status)

        if(self.curTaskRes.tableIdx !== tableIndex){
          return
        }
        if (this.curTaskRes.status === '任务完成') {
          self.curTaskRes.dataLoading = true
          let params = new URLSearchParams()
          params.append('data_path', this.curTaskRes.dataPath)
          params.append('status', this.curTaskRes.statusCode)
          console.log(self.$refs.plTable[tableIndex])
          if (self.$refs.plTable[tableIndex] !== undefined) {   // 第一次弹出时为undefined
            self.$refs.plTable[tableIndex].reloadData([])
          }
          console.log(params)
          axios({
            method: 'post',
            url: '/task/data/path',
            data: params
          }).then(function (response) {
            let res = response.data
            if (res.Res) {
              let table = JSON.parse(res.Info)
              self.curTaskRes.contentCols = table.cols
              self.curTaskRes.contentData = table.data
              self.curTaskRes.dataLoading = false
              if ((table.data !== null) && (table.data.length !== undefined)) {
                self.$refs.plTable[tableIndex].reloadData(self.curTaskRes.contentData)
              }
            } else {
              self.curTaskRes.dataLoading = false
              self.$message({
                message: res.Info,
                type: 'error'
              })
            }
          }).catch(function (error) {
            self.curTaskRes.dataLoading = false
            self.$message({
              message: error,
              type: 'error'
            })
            console.log(error)
          })
        }
      },
      onClickTab (tabObj) {
        // console.log(tabObj)
        let self = this
        this.curTab = tabObj.label

        if (tabObj.index < 4) {
          self.curTaskRes.taskId = '';
        } else {
          self.curTaskRes.taskId = tabObj.name
        }

        self.curTaskRes.status = ''
        self.curTaskRes.status_code = 0
        self.curTaskRes.dataPath = ''
        self.curTaskRes.filePath = ''
        self.curTaskRes.dataLoading = false
        self.curTaskRes.errMsg = ''
        self.curTaskRes.contentCols = []
        self.curTaskRes.contentData = []

        if (tabObj.label === '数据预览') {
          if (this.activeRightTabs === 'second') {
            let dsType = this.searchFieldDataSource.dsType
            let dsId = this.searchFieldDataSource.dsId
            if (dsId !== '' && dsId != null && dsId != undefined) {
              this.onExploreData(this.exploreTableName, dsType, dsId)
            }
          } else {
            this.onExploreData(this.exploreTableName)
          }
        } else if (tabObj.label === '历史任务') {
          this.loadTaskData();
        } else if (tabObj.index >= 4) {   // 第5个标签及以后是提交的任务
          let taskName = tabObj.label     // 绑定的任务名
          //let taskId = self.curTaskRes.taskId        // 绑定的任务id

          console.log(tabObj.index)
          self.curTaskRes.tableIdx = tabObj.index - 4
          self.curTaskRes.process = self.exploreResult.tasks[self.curTaskRes.tableIdx]['process']

          axios.get('/task_god/query/' + self.curTaskRes.taskId)
            .then(function (response) {
              let data = response.data
              self.curTaskRes.status = data.status
              self.curTaskRes.statusCode = data.status_code
              // self.curTaskRes.errMsg = data.ExceptionMsg // 弹框
              // (2) 数据文件地址
              self.curTaskRes.dataPath = data.data_path
              self.curTaskRes.filePath = self.dataPathRoot + 'datacenter/' + self.curTaskRes.dataPath
              // (3) 预览5000条的接口
              //self.curTaskRes.taskId = taskId
              console.log(self.curTaskRes)

              if (self.curTaskRes.status === '任务完成' || self.curTaskRes.status === '任务取消') {
                self.curTaskRes.process = 100
                self.exploreResult.tasks[self.curTaskRes.tableIdx]['process'] = self.curTaskRes.process
              }

              self.getExploreTaskData(tabObj.index - 4)    // 计算是第几个plTable
            }).catch(function (error) {
            console.log(error)
          })
        }
      }
      ,
      beforeUpload (file) {
        let self = this
        let formData = new FormData()
        let username = getCookie('_adtech_user')
        formData.append('file', file)
        formData.append('username', username)

        axios.post('/dw/upload/file', formData, {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          }
        ).then(function (response) {
          console.log(response.data)
          let resData = response.data
          resData = resData.match(/part='(\S*)'/)[1]
          console.log(resData)
          self.queryWordsPartition = resData

          self.uploadTableName = response.data
          self.$message({
            type: 'success',
            message: '提交成功,表名为:' + response.data
          })
        }).catch(function (error) {
          console.log('失败')
          console.log(error)
          self.$message({
            type: 'success',
            message: '提交失败' + error
          })
        })
        return false   // 阻断action中的传输, 用于自己实现axios传输文件, 返回true会访问action中的url
      }
      ,
      onExceed (files, fileList) {
        this.$set(fileList[0], 'raw', files[0])
        this.$set(fileList[0], 'name', files[0].name)
        this.$refs['upload'].clearFiles()//清除文件
        this.$refs['upload'].handleStart(files[0])//选择文件后的赋值方法
      }
      ,
      formatCode () {
        let editor = this.codemirror
        console.log(editor)
        let sqlContent = editor.getValue()
        editor.setValue(sqlFormatter.format(sqlContent))
      }
      ,
      handleFileUpload () {
        this.file = this.$refs.file.files[0]
        console.log(this.file)
      }
      ,
      submitFile () {
        let self = this
        this.$confirm('请确保提交的文件中只包含查询词一列; 表生成完毕会产生提示, 分区名会在底部显示', '提示', {
          distinguishCancelAndClose: true,
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
          center: true
        }).then(() => {
          self.$refs.upload.submit()
        })

      }
      ,

      onSubmit (formName) {
        let self = this
        this.$refs[formName].validate((valid) => {
          if (valid) {
            this.$confirm('不正确的SQL查询会造成资源浪费，请反复校验查询语句。 是否提交?', '提示', {
              distinguishCancelAndClose: true,
              confirmButtonText: '确定',
              cancelButtonText: '取消',
              type: 'warning',
              center: true
            }).then(() => {
              let params = new URLSearchParams()
              params.append('sql', this.form.sql)
              self.isOnSubmit = true
              axios({
                method: 'post',
                url: '/dw/sql/id',
                data: params,
                responseType: 'string'
              })
                .then(function (response) {
                  self.sqlQueryId = response.data
                  self.isShowRes = true
                  self.$options.methods.onTask.bind(self)()
                  self.isOnSubmit = false
                })
                .catch(function (error) {
                  self.isOnSubmit = false
                  this.$message({
                    type: 'error',
                    message: error
                  })
                  console.log(error)
                })
            })
          } else {
            return false
          }
        })
      }
      ,
      onReset (formName) {
        /* Reset our form values */
        this.$refs[formName].resetFields()
        /* Trick to reset/clear native browser form validation state */
        this.$nextTick(() => {
          this.isShowRes = false
        })
      }
      ,
      msgNotify () {
        this.$notify({
          message: '服务器接收到您的请求并分配查询ID。您可以转到任务队列页面查看结果！',
          type: 'success',
          duration: 2000,
          onClick: this.redirectClick,
        })
      },

      onTask () {
        let self = this
        let params = new URLSearchParams()
        let taskName = this.form.name
        if (taskName === '') {
          taskName = this.sqlQueryId
        }
        params.append('name', taskName)
        params.append('sql', this.form.sql)
        params.append('queryId', this.sqlQueryId)
        params.append('queryWordsPartition', this.queryWordsPartition)
        params.append('sqlType', this.datasource.typeSelectValue)
        params.append('dsId', this.datasource.selectValue)

        axios.post('/dw/sql/task', params, {responseType: 'json'})
          .then(function (response) {
            self.show = true
            if (response.data.Res) {
              self.$options.methods.msgNotify.bind(self)()
              self.exploreResult.tasks.push({
                'id': self.sqlQueryId,
                'name': taskName,
                'process': 0,
              })    // 标签中增加任务名
            } else {
              self.$message({
                type: 'error',
                message: response.data.Info
              })
            }
          })
          .catch(function (error) {
            self.$options.methods.errorNotify.bind(self)()
            self.isShowRes = false
            console.log(error)
            self.$message({
              type: 'error',
              message: error
            })
          })
      }
      ,
      onCreateSqlView () {
        if(this.disableRedashCreteBtn === true){
          return;
        }
        this.disableRedashCreteBtn = true;

        let self = this
        let params = new URLSearchParams()
        let taskName = this.form.name
        if (taskName === '') {
          this.$alert('请输入查询任务名称!')
          this.disableRedashCreteBtn = false
          return
        }

        let sqlType = this.datasource.typeSelectValue;
        if (sqlType != 'mysql' && sqlType != 'clickhouse') {
          this.$alert('当前引擎不支持创建可视化！')
          this.disableRedashCreteBtn = false
          return
        }

        params.append('name', taskName)
        params.append('sql', this.form.sql)
        params.append('queryId', this.sqlQueryId)
        params.append('queryWordsPartition', this.queryWordsPartition)
        params.append('sqlType', this.datasource.typeSelectValue)
        params.append('dsId', this.datasource.selectValue)


        self.$message({
          type: 'success',
          message: "可视化SQL创建中..."
        })
        axios.post('/jdbc/sql/createQueryView', params, {responseType: 'json'})
          .then(function (response) {
            if (response.data.msg.indexOf('成功') >= 0) {
              self.curRedashQueryId = response.data.redash_query_id
              let queryUrl = 'http://inner-redash.adtech.sogou/queries/' + self.curRedashQueryId + '/source'
              self.$router.push({name: 'redashqueryshow', params: {queryid: self.curRedashQueryId, queryUrl: queryUrl}})
            } else {
              self.$message({
                type: 'error',
                message: response.data.msg
              });
              self.disableRedashCreteBtn = false;
            }

          })
          .catch(function (error) {
            console.log(error)
            self.disableRedashCreteBtn = false;
          })

      },
      filterNode (value, data) {
        if (!value) return true
        return data.label.indexOf(value) !== -1
      }
      ,
      errorNotify () {
        this.$notify.error({
          title: '错误',
          message: 'SQL填写有误。请重新填写后查询！'
        })
      }
      ,
      loadHistory () {
        let self = this
        const hist = sessionStorage.getItem('sub_history')
        sessionStorage.removeItem('sub_history')
        if (hist) {
          let json2map = JSON.parse(hist)
          try {
            self.form.sql = json2map['SqlStr']
            self.form.name = json2map['QueryName']
            let queryJson = JSON.parse(json2map['QueryJson'])
            self.datasource.typeSelectValue = queryJson['sqlType']
            // 元素可见性
            if (self.datasource.typeSelectValue === 'hive') {
              self.visibility = true
            } else {
              self.visibility = false
            }
            self.treeData = ''
            axios.get('/dw/GetConnByPrivilege/' + this.datasource.typeSelectValue)
              .then(function (response) {
                self.datasource.options = response.data
                self.datasource.selectValue = queryJson['dsid']
                self.onChangeDataSource()
              })
              .catch(function (error) {
                self.datasource.selectValue = ''
                self.datasource.options = []
                console.log(error)
              })

            // self.form.sql = decodeURIComponent(getQueryVariable(hist, 'sql').replace(/\+/g, ' '))
            // if (getQueryVariable(hist, 'name')) {
            //   self.form.name = decodeURIComponent(getQueryVariable(hist, 'name').replace(/\+/g, ' '))
            // }
          } catch (e) {
            this.onReset('form')
          }
        }
      },
      setProcessBar () {
        let self = this

        let curTaskId = this.curTaskRes.taskId
        let curTaskIdx = 0
        for (let i = 0; i < this.exploreResult.tasks.length; i++) {
          if (this.exploreResult.tasks[i]['id'] === curTaskId) {
            curTaskIdx = i
            self.curTaskRes.process = this.exploreResult.tasks[i]['process']

            continue
          }

          if (this.exploreResult.tasks[i]['process'] < 95) {
            this.exploreResult.tasks[i]['process'] += Math.floor(((Math.random() * 10) + 0) % 6 / 5)
          }
        }

        if (self.curTaskRes.process === 100 || self.curTaskRes.taskId === '') {
          return
        }

        axios.get('/task_god/query/' + self.curTaskRes.taskId)
          .then(function (response) {
            let data = response.data

            if(self.curTaskRes.taskId !== curTaskId){
              return
            }
            self.curTaskRes.status = data.status
            self.curTaskRes.statusCode = data.status_code
            // self.curTaskRes.errMsg = data.ExceptionMsg // 弹框
            // (2) 数据文件地址
            self.curTaskRes.dataPath = data.data_path
            self.curTaskRes.filePath = self.dataPathRoot + 'datacenter/' + self.curTaskRes.dataPath
            // (3) 预览5000条的接口
            console.log(self.curTaskRes)
            self.getExploreTaskData(self.curTaskRes.tableIdx)    // 计算是第几个plTable

            if (self.curTaskRes.status === '任务完成' || self.curTaskRes.status === '任务取消') {
              self.curTaskRes.process = 100
              self.exploreResult.tasks[self.curTaskRes.tableIdx]['process'] = self.curTaskRes.process
            } else {
              if (self.curTaskRes.process < 95) {
                self.curTaskRes.process += Math.floor(((Math.random() * 10) + 0) % 6 / 5)
                self.exploreResult.tasks[self.curTaskRes.tableIdx]['process'] = self.curTaskRes.process
              }
            }
            console.log(self.curTaskRes.process)
          }).catch(function (error) {
          console.log(error)
        })
      },


      pathClick (row) {
        let self = this
        const statusCode = row.status_code
        let dataPath = row.data_path
        if (dataPath.indexOf("xlsx") === -1) {
          dataPath = dataPath + ".csv"
        }
        console.log(dataPath)
        if (isRunFinsh(statusCode) || isRunFinshPagani(statusCode)) {
          var path = this.dataPathRoot + 'datacenter/' + dataPath + '?op=OPEN'
          var pathShow = this.dataPathRoot + '\n' + 'datacenter/' + dataPath + '?op=OPEN'
          if (isRunFinshPagani(statusCode)) {
            path = this.dataPathRootPaga + dataPath + '?op=OPEN'
            pathShow = this.dataPathRootPaga + '\n' + dataPath + '?op=OPEN'
          }
          window.open(path)
        } else {
          self.$message({
            message: '任务没有正常完成，不支持查看！',
            type: 'error'
          })
        }
      },

      redirectClick () {
        router.push({path: 'que'})
      },
      redirectHistoryTaskClick (row) {
        let self = this
        const queryId = row.query_id
        const taskName = row.query_name

        axios.get('/task/queryLog/' + queryId)
          .then(function (response) {
            console.log(response.data)
            sessionStorage.setItem('sub_history', JSON.stringify(response.data))

            let tasks = self.exploreResult.tasks
            let idx = -1
            for (let i = 0; i < tasks.length; i++) {
              if (tasks[i].id === queryId) {
                idx = i
              }
            }

            if(idx === -1){
              self.exploreResult.tasks.push({
                  'id': queryId,
                  'name': taskName,
                  'process': 0,
                })    // 标签中增加任务名
            }

            self.loadHistory();
          })
          .catch(function (error) {
            console.log(error)
          })
        //const route = type2route(row.type)
      },
      filterHistoryTask (task) {
        const newtask = [];
        for(let i=0;i<task.length;i++){
          console.log(task[i]);
          if (task[i]["type"] === "sql查询"){
            newtask.push(task[i])
          }
        }
        return newtask;
      },
      loadTaskData () {
        let self = this;
        self.taskLoading = true;
        self.taskForm.session = getCookie('_adtech_user');

        const userURL = '/task/log/config?session=' + this.taskForm.session
          + '&startDate=' + this.taskForm.startDate + '&date=' + this.taskForm.date + '&status=' + this.taskForm.status;
        axios.get(userURL)
          .then(function (response) {
            self.taskLoading = false;
            self.taskCols = response.data.cols;

            //self.taskData = response.data.data;
            self.taskData = self.filterHistoryTask(response.data.data);

          }).catch(function (error) {
          console.log(error);
        })
      },

    },
    computed: {
      codemirror () {
        console.log('计算了')
        return this.$refs.cmEditor.codemirror
      },
      verifyButtonText: function () {
        if (this.onVerify) {
          return '校验中'
        } else {
          return '校验sql'
        }
      }
    }
    ,
    mounted () {
      let self = this;
      axios.get('/jdbc/template/getviewaccess')
        .then(function (response) {
          // 顶栏标签的配置中, 是否配置了 可视化权限
          if(response.data.viewaccess===true)   {
            self.viewAccess = true
          }
          console.log('Access')
          console.log(self.viewAccess)
        })
        .catch(function (error) {
          console.log(error)
        });


      $(document).ready(function () {
        let h = document.body.clientHeight
        console.log(h)

        let svgResize = document.getElementById('resize')
        let svgTop = document.getElementById('top')
        let svgDown = document.getElementById('down')
        let svgBox = document.getElementById('box')
        let treeDiv = document.getElementById('treeDiv')

        svgBox.style.height = (h - 120) + 'px'
        svgTop.style.height = (h - 120) * 0.5 + 'px'
        svgDown.style.height = (h - 120) * 0.5 + 'px'
        treeDiv.style.height = (h - 120) + 'px'

        let codeMirrorHeight = svgTop.offsetHeight - 120 + 'px'
        console.log('codemirror:' + codeMirrorHeight)
        console.log('top:' + svgTop.style.height)
        $('.CodeMirror').css('height', codeMirrorHeight)

        svgResize.onmousedown = function (e) {
          let startY = e.clientY
          svgResize.top = svgResize.offsetTop
          document.onmousemove = function (e) {
            let endY = e.clientY
            let moveLen = svgResize.top + (endY - startY)
            let maxT = svgBox.clientHeight - svgResize.offsetHeight
            if (moveLen < 30) moveLen = 30
            if (moveLen > maxT - 30) moveLen = maxT - 30
            svgResize.style.top = moveLen
            svgTop.style.height = moveLen + 'px'
            svgDown.style.height = (svgBox.clientHeight - moveLen - 5) + 'px'

            let oldCMHeight = svgTop.offsetHeight - 120 + 'px'
            $('.CodeMirror').css('height', oldCMHeight)
          }
          document.onmouseup = function (evt) {
            document.onmousemove = null
            document.onmouseup = null
            svgResize.releaseCapture && svgResize.releaseCapture()
          }
          svgResize.setCapture && svgResize.setCapture()
          return false
        }

      });


      this.loadHistory()




      this.timer = setInterval(this.setProcessBar, 4000)
    },
    beforeDestroy () {
      clearInterval(this.timer)
    }
  }
</script>

<style type='text/css' scoped>

  .CodeMirror {
    border: 1px solid #eee;
    /*height: 400px;*/
    font-size: 13px;
    line-height: 120%;
  }

  /*  .line {
      width: 100%;
      height: 1px;
      background: #e0e7eb;
      position: relative;
      text-align: center;
      margin-top: 40px;
      cursor: s-resize;
    }*/

  h3 {
    font-size: large;
    margin-block-start: 0px;
    margin-block-end: 0px;
  }

  #box {
    width: 100%;
    height: 100%;
    position: relative;
    overflow: hidden;
    /*border: 1px solid blue;*/
  }

  #top {
    /*height:calc(64%);*/
    width: 100%;
    float: left;
    overflow: auto;
  }

  #resize {
    margin-top: -20px;
    margin-bottom: 10px;
    position: relative;
    /*background: #e0e7eb;*/
    height: 20px;
    width: 100%;
    cursor: pointer;
    float: left;
    text-align: center;
  }

  #down {
    height: 100%;
    width: 100%;
    float: left;
    overflow: auto;
  }

  #treeDiv {
    height: 100%;
    width: 100%;
    overflow: auto;
  }


</style>



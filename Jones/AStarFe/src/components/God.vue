/* eslint-disable */
<template>
  <div>
    <h3>任务队列</h3>
    <!--筛选框-->
    <el-row>

      <el-col :span="8">
        <el-date-picker
          size="mini"
          v-model="taskForm.date"
          type="date"
          placeholder="选择任务日期"
          format="yyyy 年 MM 月 dd 日"
          value-format="yyyy-MM-dd"
          @change="loadTaskData">
        </el-date-picker>
      </el-col>

      <el-col :span="8">
        <el-select size="mini" v-model="taskForm.status" placeholder="选择任务状态" @change="loadTaskData">
          <el-option label="全部" value=""></el-option>
          <el-option label="任务提交中" value="0"></el-option>
          <el-option label="任务执行中" value="1"></el-option>
          <el-option label="任务正在取消" value="2"></el-option>
          <el-option label="任务完成" value="3"></el-option>
          <el-option label="任务取消" value="4"></el-option>
        </el-select>
      </el-col>

      <el-col :span="5">
        <el-input v-model="taskForm.session"
                  clearable
                  size="mini"
                  placeholder="输入用户名"
                  @change="loadTaskData">
        </el-input>
      </el-col>

    </el-row>

    <!--任务队列表-->
    <el-table
      v-loading="taskLoading"
      :data="taskData"
      stripe
      style="width: 100%"
      :default-sort='{prop:"create_time", order:"descending"}'>
      <el-table-column
        v-for="taskCol in taskCols"
        :key="taskCol.prop"
        :prop="taskCol.prop"
        :label="taskCol.label"
        :sortable="taskCol.sortable">
      </el-table-column>
      <el-table-column
        fixed="right"
        label="异常">
        <template slot-scope="scope">
          <el-button @click="checkMessage(scope.row)" type="text" size="mini" v-show="scope.row.status_code==4">查看异常
          </el-button>
        </template>
      </el-table-column>
      <el-table-column
        fixed="right"
        label="操作"
        width="140">
        <template slot-scope="scope">
          <el-button @click="pathClick(scope.row)" type="text" size="mini">路径</el-button>
          <el-button @click="dataClick(scope.row)" type="text" size="mini">数据</el-button>
          <el-button @click="redirectClick(scope.row)" type="text" size="mini">配置</el-button>
          <el-button @click="assignUser(scope.row)" type="text" size="mini">指派</el-button>
          <el-button @click="deleteTask(scope.row)" type="text" size="mini">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      title="异常信息"
      :visible.sync="dialogVisible"
      width="50%">
      <el-input type="textarea"
                v-model="errMsg"
                placeholder="无异常内容显示"
                :rows="20"
                :max-rows="30"
                :readonly="true"></el-input>
    </el-dialog>

    <!--数据结果表-->
    <el-dialog title="查询结果" :visible.sync="dialogTableVisible" width="70%">
      <el-row>
        <span> 该表为样例数据，最多展现5000条。</span>
        <!--        <el-button type="success" icon="el-icon-download" circle size="mini" @click="exportExcel"></el-button>-->
      </el-row>
      <!--<el-row>
        <span>提示：下载数据为表格展示数据，获取全量数据请查看路径使用RSYNC下载。</span>
      </el-row>-->
      <pl-table ref="plTable"
                selectTrClass="selectTr"
                :max-height="500"
                :use-virtual="true"
                header-drag-style
                fixedColumnsRoll
                inverse-current-row
                v-loading="dataLoading">
        <template slot="empty">
          没有查询到符合条件的记录
        </template>
        <pl-table-column
          v-for="item in contentCols"
          :key="item.id"
          :resizable="true"
          :show-overflow-tooltip="true"
          :prop="item.prop"
          :label="item.label"
          :fixed="item.fixed"/>
      </pl-table>
    </el-dialog>

    <el-dialog title="分发给用户" :visible.sync="assignUserVisible">
      <el-form>
        <el-form-item>
          <el-input v-model="assignUserName" placeholder="请输入用户名"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onClickAssign()">确定</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>

<style>
  .el-table .warning-row {
    background: oldlace;
  }

  .el-table .success-row {
    background: #f0f9eb;
  }
</style>

<script>
  const axios = require('axios')
  import XLSX from 'xlsx'
  import FileSaver from 'file-saver'
  import router from '../router'
  import {getCookie, setCookie, formatDateLine, isRunFinsh, isRunFinshPagani, type2route} from '../conf/utils'

  export default {
    data () {
      return {
        taskCols: [],
        taskData: [],
        taskForm: {
          session: '',
          date: formatDateLine(new Date()),
          status: '',
        },
        contentCols: [],
        contentData: [],
        contentFiled: '',
        dialogTableVisible: false,
        taskLoading: false,
        dataLoading: false,
        dataPathRoot: '',
        dataPathRootPaga: '',
        dialogVisible: false,
        errMsg: '',
        assignUserName: '',
        assignTaskId: '',
        assignUserVisible: false,
        maxHeight: '',
      }
    },
    methods: {
      deleteTask (row) {
        let queryId = row.query_id
        let self = this
        axios({
          method: 'delete',
          url: '/task/deleteTaskRecoder/' + queryId
        }).then(function (response) {
          if (response.data.Res) {
            self.$message({
              message: '删除成功',
              type: 'success'
            })
            self.loadTaskData()
          } else {
            self.$message({
              message: '删除失败:' + response.data.Info,
              type: 'error'
            })
          }
        }).catch(function (error) {
          console.log(error)
        })
      },
      onClickAssign () {
        let param = {
          'assignUserName': this.assignUserName,
          'assignTaskId': this.assignTaskId
        }
        console.log(param)
        let self = this
        axios.post('/task/changeTaskUser', param, {responseType: 'json'})
          .then(function (response) {
            let res = response.data
            if (res.Res) {
              self.$message({
                message: '指派成功',
                type: 'success'
              })
              self.loadTaskData()
            } else {
              self.$message({
                message: res.Info,
                type: 'error'
              })
            }
          }).catch(function (error) {
          console.error(error)
        })
        this.assignUserVisible = false
      },
      checkMessage (row) {
        this.errMsg = ''
        this.dialogVisible = true
        let self = this
        let params = {
          'queryId': row.query_id
        }
        axios.post('/task/checkException', params, {responseType: 'json'})
          .then(function (response) {
            self.errMsg = response.data.res
          }).catch(function (error) {
          console.error(error)
        })
      },
      thisType2route (type) {
        let res = type2route(type)
        console.log(res)
        return res
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
          this.$confirm(pathShow, 'RSYNC路径', {
            confirmButtonText: '点我下载',
            cancelButtonText: '取消'
          }).then(() => {
            this.$copyText(path).then(function () {
              window.open(path)
            }, function () {
              self.$message({
                message: '复制失败',
                type: 'error'
              })
            })
          })
        } else {
          self.$message({
            message: '任务没有正常完成，不支持查看！',
            type: 'error'
          })
        }
      },
      dataClick (row) {
        let self = this
        const statusCode = row.status_code
        const dataPath = row.data_path
        if (isRunFinsh(statusCode) || isRunFinshPagani(statusCode)) {
          self.dialogTableVisible = true
          self.dataLoading = true
          var params = new URLSearchParams()
          params.append('data_path', dataPath)
          params.append('status', statusCode)
          if (self.$refs.plTable !== undefined) {   // 第一次弹出时为undefined
            self.$refs.plTable.reloadData([])
          }

          axios({
            method: 'post',
            url: '/task/data/path',
            data: params
          }).then(function (response) {
            let res = response.data
            if (res.Res) {
              let table = JSON.parse(res.Info)
              self.contentCols = table.cols
              self.contentData = table.data
              self.dataLoading = false
              if ((table.data !== null) && (table.data.length !== undefined)) {
                self.$refs.plTable.reloadData(self.contentData)
              }
            } else {
              self.dataLoading = false
              self.$message({
                message: res.Info,
                type: 'error'
              })
            }
          }).catch(function (error) {
            self.dataLoading = false
            self.$message({
              message: error,
              type: 'error'
            })
            console.log(error)
          })
        } else {
          self.$message({
            message: '任务没有正常完成，不支持查看！',
            type: 'error'
          })
        }
      },
      assignUser (row) {
        console.log(row)
        this.assignUserVisible = true
        this.assignTaskId = row.query_id
      },
      redirectClick (row) {
        let self = this
        const queryId = row.query_id
        const route = type2route(row.type)
        if (route) {
          switch (route) {
            case 'jdbc':
              axios.get('/task/queryLog/' + queryId)
                .then(function (response) {
                  console.log(response.data)
                  sessionStorage.setItem('jdbc_history', JSON.stringify(response.data))
                  router.push({path: route})
                })
                .catch(function (error) {
                  console.log(error)
                })
              break
            case 'sql':
              axios.get('/task/queryLog/' + queryId)
                .then(function (response) {
                  console.log(response.data)
                  sessionStorage.setItem('sub_history', JSON.stringify(response.data))
                  router.push({path: route})
                })
                .catch(function (error) {
                  console.log(error)
                })
              break
            default:
              axios.get('/task/log/query/' + queryId)
                .then(function (response) {
                  sessionStorage.setItem('sub_history', JSON.stringify(response.data))
                  // if (row.type === '自定义查询') {
                  //   sessionStorage.setItem('sub_history', response.data)
                  // }
                  router.push({path: route})
                })
                .catch(function (error) {
                  console.log(error)
                })
          }
        } else {
          // 待修改
          self.$message({
            message: '任务类型不支持查看配置！',
            type: 'error'
          })
        }
      },
      reRun (row) {
        let self = this
        const queryId = row.query_id
        console.log('queryId:' + queryId + ', taskType:' + row.type)
        let params = {
          'queryId': queryId
        }

        axios.post('/task/reRun', params, {responseType: 'json'})
          .then(function (response) {
            console.info(response.data.queryId)
            self.$confirm('重跑成功, 点击确定刷新页面', '提示', {
              distinguishCancelAndClose: true,
              confirmButtonText: '确定',
              cancelButtonText: '取消',
              type: 'warning',
              center: true
            }).then(() => {
              location.reload()
            })
          })
          .catch(function (error) {
            self.$message({
              message: '重跑失败！' + error.toString(),
              type: 'error'
            })
          })

      },
      exportExcel () {
        /* generate workbook object from table */
        let wb = XLSX.utils.table_to_book(document.querySelector('#resTable'))
        /* get binary string as output */
        let wbout = XLSX.write(wb, {bookType: 'xlsx', bookSST: true, type: 'array'})
        try {
          FileSaver.saveAs(new Blob([wbout], {type: 'application/octet-stream'}), '查询结果表.xlsx')
        } catch (e) {
          if (typeof console !== 'undefined') {
            console.log(e, wbout)
          }
        }
        return wbout
      },
      loadTaskData () {
        let self = this
        self.taskLoading = true
        const userURL = '/task_god/log/config?session=' + this.taskForm.session
          + '&date=' + this.taskForm.date + '&status=' + this.taskForm.status
        axios.get(userURL)
          .then(function (response) {
            self.taskLoading = false
            self.taskCols = response.data.cols
            self.taskData = response.data.data
          }).catch(function (error) {
          console.log(error)
        })
      }
    },
    mounted: function () {
      this.loadTaskData()
      this.maxHeight = (document.body.clientHeight) * 0.8 - 100 + 'px'
    },
    created: function () {
      let self = this
      axios.get('/task/log/rsync')
        .then(function (response) {
          self.dataPathRoot = response.data[0]
          self.dataPathRootPaga = response.data[1]
        }).catch(function (error) {
          console.log(error)
        }
      )
    }
  }
</script>

<style>
  .el-table tbody tr:hover > td {
    background-color: #ecf5ff !important
  }
</style>

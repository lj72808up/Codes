<template>
  <div>
    <el-row align="middle">
      <el-col :span="1">
        <h3>DAGs</h3>
      </el-col>
      <el-col :offset="20" :span="1">
        <el-button type="primary" style="margin-top: 15px"
                   @click="onClickNew">新建DAG
        </el-button>
      </el-col>
    </el-row>
    <el-row>
      <div>
        <el-table :data="dagData" ref="dagTable" row-key="name" stripe>
          <el-table-column label="开关">
            <template slot-scope="dag">
              <el-switch v-if="dag.row.status === 2" v-model="dag.row.is_open" active-color="#13ce66"
                         inactive-color="#ff4949"
                         :active-value="1" :inactive-value="0" @change="changeSwitchFlag(dag.row)">
              </el-switch>
              <span v-else-if="dag.row.status === 1">Dag创建中</span>
              <span v-else-if="dag.row.status === 8">Dag删除中</span>
            </template>
          </el-table-column>

          <el-table-column property="name" label="DAG Id"/>
          <el-table-column property="exhibitName" label="DAG Name"/>
          <el-table-column prop="schedule" label="Schedule"/>
          <el-table-column prop="user" label="user"/>
          <el-table-column fixed="right" label="操作" width="280">
            <template slot-scope="dag">
              <el-button icon="el-icon-delete" type="text" size="small" @click="DeleteDag(dag.$index,dag.row.name)">[删除]
              </el-button>
              <el-button icon="el-icon-setting" type="text" size="small" @click="DagTaskShow(dag.row)">[调度]</el-button>
              <el-button icon="el-icon-setting" type="text" size="small" @click="DagTaskConfig(dag.row)">[配置]
              </el-button>
              <el-button icon="el-icon-setting" type="text" size="small" @click="manualTrigger(dag.row.name)">[手动触发]
              </el-button>
              <el-button icon="el-icon-thumb" type="text" size="small" @click="dispatch(dag.row.name)">[指派]</el-button>
            </template>
          </el-table-column>
          <!-- dag 状态暂时放到 task里
          <el-table-column  label="调度状态">
            <template slot-scope="dag">
              <el-button icon="el-icon-star-off" type="text" size="small" @click="MarkDagSheduleStatus(dag.$index,dag.row,'notrun')">清除</el-button>
              <el-button icon="el-icon-check" type="text" size="small" @click="MarkDagSheduleStatus(dag.$index,dag.row,'suc')">成功</el-button>
              <el-button icon="el-icon-close" type="text" size="small" @click="MarkDagSheduleStatus(dag.$index,dag.row,'fail')">失败</el-button>
            </template>
          </el-table-column>
          -->

        </el-table>

        <el-dialog title="分发给用户" :visible.sync="assignUser.visible">
          <el-form>
            <el-form-item>
              <el-input v-model="assignUser.name" placeholder="请输入用户名"/>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="onClickAssign()">确定</el-button>
            </el-form-item>
          </el-form>
        </el-dialog>

        <el-dialog title="请选择调度的起始时间" :visible.sync="triggerVisible">
          <el-form style="margin-top: 8px; ">
            <el-row>
              <el-col :span="11">
                <el-form-item label="开始日期: ">
                  <el-date-picker
                    v-model="triggerForm.startTime"
                    format="yyyy年MM月dd日"
                    value-format="yyyyMMdd"
                    placeholder="开始日期"
                    style="width: calc(100% - 100px)"
                    type="date"
                    size="small"/>
                </el-form-item>
              </el-col>
              <el-col :span="11">
                <el-form-item label="结束日期: ">
                  <el-date-picker
                    v-model="triggerForm.endTime"
                    format="yyyy年MM月dd日"
                    value-format="yyyyMMdd"
                    placeholder="结束日期"
                    style="width: calc(100% - 100px)"
                    type="date"
                    size="small"/>
                </el-form-item>
              </el-col>
              <el-col :span="1">
                <el-form-item>
                  <el-button type="primary" @click="onTrigger()" size="small" style="float:right">确定</el-button>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
          <!--          <el-row>-->
          <div style="max-height: 500px ; overflow-y:scroll; "> <!--style="max-height: 500px ; overflow-y:scroll; "-->
            <el-scrollbar style="height: 100%;">
              <el-table :data="readyData" stripe style="width: 100%">
                <el-table-column v-for="readyCol in readyCols" v-if="!readyCol.hidden"
                                 :key="readyCol.prop" :prop="readyCol.prop" :label="readyCol.label"
                                 :sortable="readyCol.sortable"/>
                <el-table-column key="operate" prop="operate" label="操作" width="100px" align="center">
                  <template slot-scope="scope">
                    <el-button icon="el-icon-info" type="text" @click="onClickLog(scope.row)" size="mini">日志
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-scrollbar>
          </div>
          <!--          </el-row>-->
        </el-dialog>

        <el-dialog title="查看日志" :visible.sync="logVisible">
          <el-tabs tab-position="left" style="height: 600px;" @tab-click="checkModelLog" v-model="defaultModel">
            <el-tab-pane v-for="modelId in triggerModelIds" :label="modelId">
              <div style="max-height: 600px; overflow-y: auto">
                <span style="white-space:pre-wrap">{{ modelLog }}</span>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-dialog>
      </div>
    </el-row>
  </div>
</template>

<script>
const axios = require('axios')
import router from '../router'
import {getCookie, isEmpty} from '../conf/utils'

export default {
  data() {
    return {
      defaultModel: "",
      modelLog: "",
      triggerModelIds: [],
      logVisible: false,
      readyData: [],
      readyCols: [],
      triggerForm: {
        startTime: "",
        endTime: "",
      },
      dagData: [],
      assignUser: {
        visible: false,
        name: ''
      },
      curDagId: '',
      triggerVisible: false,
      triggerUser: '',
    }
  },
  methods: {
    checkModelLog(tab) {
      let labelSplits = tab.label.split("_")
      let modelId = labelSplits[labelSplits.length - 1]
      if (labelSplits.length === 1) {
        modelId = labelSplits[0]
      }
      let self = this
      let param = {
        "remoteId": modelId,
        "userName": this.triggerUser,
      }
      axios.post('/airflow/getTriggerLog', param, {responseType: 'json'})
        .then(function (response) {
          self.modelLog = response.data
        })
        .catch(function (error) {
          console.log(error)
        })
    },
    onClickLog(row) {
      this.triggerUser = row.username
      this.triggerModelIds = []
      let remoteIds = []
      let remoteArr = JSON.parse(row.remote_ids)
      remoteArr.forEach(ele => {
        let id = ele.modelName + "_" + ele.remoteId
        remoteIds.push(id)
      })
      // this.$alert(remoteIds)
      // let status =  row.status
      if (remoteIds !== "") {
        this.triggerModelIds = remoteIds
        this.logVisible = true
        this.defaultModel = this.triggerModelIds[0]
        console.log(this.triggerModelIds)
        console.log("defaultModel:" + this.defaultModel)
      }
    },
    getTriggerHistory() {
      let self = this
      let param = {
        "dagName": this.curDagId
      }
      axios.post('/airflow/getManualTriggerHistory', param, {responseType: 'json'})
        .then(function (response) {
          self.readyData = response.data.data
          self.readyCols = response.data.cols
        })
        .catch(function (error) {
          console.log(error)
          self.$message({
            type: 'error',
            message: error,
          })
        })
    },
    onTrigger() {
      console.log(this.triggerForm)
      console.log(this.curDagId)
      if (this.triggerForm.startTime > this.triggerForm.endTime) {
        this.$alert("开始时间应该小于结束时间")
      } else {
        let self = this
        let param = {
          "startTime": this.triggerForm.startTime,
          "endTime": this.triggerForm.endTime,
          "dagName": this.curDagId,
        }
        axios.post('/airflow/manualTrigger', param, {responseType: 'json'})
          .then(function (response) {
            self.$message({
              type: 'success',
              message: '手动触发成功',
            })
            self.triggerVisible = false
            self.triggerForm = {
              startTime: "",
              endTime: "",
            }
          })
          .catch(function (error) {
            console.log(error)
            self.$message({
              type: 'error',
              message: error,
            })
          })
      }
    },
    onClickAssign() {
      console.log(this.assignUser.name)
      let param = {
        "user": this.assignUser.name,
        "dagId": this.curDagId,
      }
      let self = this
      axios.post('/airflow/assignDag', param, {responseType: 'json'})
        .then(function (response) {
          self.$message({
            type: 'success',
            message: '指派成功',
          })
          self.assignUser.visible = false
          self.getDagList()
        })
        .catch(function (error) {
          console.log(error)
          self.$message({
            type: 'error',
            message: error,
          })
        })
      this.assignUser.name = ""
    },
    dispatch(dagId) {
      console.log(dagId)
      this.curDagId = dagId
      this.assignUser.visible = true
    },
    manualTrigger(dagId) {
      console.log(dagId)
      this.curDagId = dagId
      this.triggerVisible = true
      this.getTriggerHistory()
    },
    changeSwitchFlag(daginfo) {
      let self = this
      let params = {
        'name': daginfo.name,
        'isopen': daginfo.is_open.toString(),
      }

      axios.post('/airflow/switchdag', params, {responseType: 'json'})
        .then(function (response) {
          let res = response.data

          if (res['flag'] == false) {
            self.$message({
              type: 'error',
              message: '切换失败',
            })
          }
        })
        .catch(function (error) {
          console.log(error)
        })

    },
    DeleteDag(rowIdx, dagname) {
      let self = this
      let params = {
        'name': dagname,
      }
      console.log('delete index' + rowIdx.toString())

      this.$confirm('此操作将永久删除该任务, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        axios.post('/airflow/deletedag', params, {responseType: 'json'})
          .then(function (response) {
            let res = response.data

            if (res['flag'] === false) {
              self.$message({
                type: 'error',
                message: '删除失败',
              })
            } else {
              // self.dagData.splice(rowIdx, 1)
              self.$message({
                type: 'success',
                message: '操作成功, 请等待系统删除',
              })
              self.getDagList()
            }
          })
          .catch(function (error) {
            console.log(error)
          })
      }).catch(function (error) {
        console.log(error)
      })
    },

    MarkDagSheduleStatus(index, daginfo, newstatus) {
      let self = this
      let params = {
        'name': daginfo.name,
        'run_time': daginfo.lastrun,
        'new_status': newstatus,
      }
      let idx = index

      axios.post('/airflow/markdag', params, {responseType: 'json'})
        .then(function (response) {
          let res = response.data

          if (res['flag'] == false) {
            self.$message({
              type: 'error',
              message: '更新失败',
            })
          } else {
            this.dagData[index]['lastrunstatus'] = newstatus
            self.$message({
              type: 'success',
              message: '更新成功！',
            })
          }
        })
        .catch(function (error) {
          console.log(error)
        })

    },

    DagTaskConfig(daginfo) {
      sessionStorage.setItem('jones_airflow_dagId', daginfo.name)
      router.push({path: 'dag'})
    },
    onClickNew() {
      router.push({path: 'dag'})
    },

    DagTaskShow(daginfo) {
      console.log(daginfo.name)
      // 放入cache
      sessionStorage.setItem('jones_airflow_dagId_show', daginfo.name)
      //let info = {
      //  'tableId': row.tid
      //}
      //sessionStorage.setItem('metaData_table_id', JSON.stringify(info))
      // 跳转
      router.push({path: 'dagTaskTrace'})
    },

    getDagList() {
      let self = this
      axios.get('/airflow/airflowall')
        .then(function (response) {
          let adata = []
          let checkedIdx = []
          let index = 0
          self.dagData = []
          response.data.forEach(ele => {
            console.log('==================++')
            console.log(ele.IsOpen)
            self.dagData.push({
              'name': ele.Name,
              'exhibitName': ele.ExhibitName,
              'is_open': ele.IsOpen,
              'schedule': ele.Schedule,
              'user': ele.User,
              'status': ele.Status,
              'state_suc': ele.StateSuccess,
              'state_failed': ele.StateFailed,
              'state_running': ele.StateRunning,
              'task_suc': ele.TaskStateSuc,
              'task_running': ele.TaskStateRunning,
              'task_failed': ele.TaskStateFailed,
              'task_up_failed': ele.TaskStateUpFailed,
              'task_skipped': ele.TaskStateSkipped,
              'task_up_for_retry': ele.TaskStateUpForRetry,
              'task_up_for_reshedule': ele.TaskStateUpForReshedule,
              'task_queued': ele.TaskStateQueued,
              'task_null': ele.TaskStateNull,
              'task_sheduled': ele.TaskStateSheduled,
            })
            console.log(self.dagData)
          })

        }).catch(function (error) {
        console.log(error)
      })
    }
  },
  created: function () {
    this.getDagList()
  },
  mounted() {
    console.log('mounted')
  }
}
</script>

<style lang="scss" scoped>
.el-dialog-div {
  height: 60vh;
  overflow: auto;
}

</style>

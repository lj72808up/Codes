<template>
  <div>
    <h3> {{ dagId }} 任务详情:</h3>
    <el-row style="margin-top: 10px">
      <el-form inline>
        <el-form-item label="基准日期:">
          <el-date-picker
            v-model="baseDate"
            format="yyyy年MM月dd日"
            value-format="yyyy-MM-dd"
            placeholder="基准日期"
            type="date"
            @change="loadTaskData"/>
        </el-form-item>
        <el-form-item label="调度次数:">
          <el-select v-model="cnt" placeholder="选择调度次数" @change="loadTaskData">
            <el-option label="5" value="5"/>
            <el-option label="25" value="25"/>
            <el-option label="50" value="50"/>
            <el-option label="100" value="100"/>
          </el-select>
        </el-form-item>
        <el-form-item label="任务名称:">
          <el-input v-model="searchWord" @input="onChangeSearch" clearable/>
        </el-form-item>
      </el-form>
    </el-row>
    <el-table
      :data="tableForm.data"
      :row-class-name="getTableRowClassName"
      style="width: 100%">
      <el-table-column
        v-for="col in tableForm.cols"
        v-if="!col.hidden"
        :key="col.prop"
        :prop="col.prop"
        :label="col.label"
        :sortable="col.sortable">
      </el-table-column>

      <el-table-column
        fixed="right"
        label="状态"
        width="140">
        <template slot-scope="scope">
          <el-button v-if="scope.row.state==='success'" type="success" size="small">{{ stateOption.success }}
          </el-button>
          <el-button v-else-if="scope.row.state==='failed'" type="danger" size="small">{{ stateOption.danger }}
          </el-button>
          <el-button v-else-if="scope.row.state==='notrun'" size="small">{{ stateOption.notrun }}</el-button>
          <el-button v-else type="info" size="small">{{ getTaskState(scope.row.state) }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column label="状态操作">
        <template slot-scope="scope">
          <el-button icon="el-icon-star-off" size="small" @click="TaskStateOpt(scope.$index,scope.row)">修改Task状态
          </el-button>
        </template>
      </el-table-column>
      <el-table-column label="状态操作">
        <template slot-scope="scope">
          <el-button icon="el-icon-star-off" size="small" @click="DagStateOpt(scope.$index,scope.row)">修改所属Dag状态
          </el-button>
        </template>
      </el-table-column>
      <el-table-column label="日志" width="100px" align="center">
        <template slot-scope="scope">
          <el-button icon="el-icon-edit" type="text" @click="checkLog(scope.row)"/>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog title="Task状态修改" :visible.sync="taskStateOptVisible">
      <h3>{{ curTask.dagName }} -- {{ curTask.taskId }} -- {{ curTask.execution_date }} 状态修改:</h3>
      <el-button round size="small" @click="MarkTaskSheduleStatus('notrun')">重跑</el-button>
      <el-button type="success" round size="small" @click="MarkTaskSheduleStatus('success')">成功</el-button>
      <el-button type="danger" round size="small" @click="MarkTaskSheduleStatus('failed')">失败</el-button>
    </el-dialog>

    <el-dialog title="Dag状态修改" :visible.sync="dagStateOptVisible">
      <h3>{{ curTask.dagName }} -- {{ curTask.execution_date }} 状态修改:</h3>
      <el-button round size="small" @click="MarkDagSheduleStatus('notrun')">重跑</el-button>
      <!--          <el-button  type="success" round size="small" @click="MarkDagSheduleStatus('success')">成功</el-button>-->
      <!--          <el-button  type="danger" round size="small" @click="MarkDagSheduleStatus('failed')">失败</el-button>-->
    </el-dialog>

    <el-dialog title="查看日志" :visible.sync="logVisible" width="80%">
      <el-row>
        <el-col :span="11">
<!--          <span>{{logUrl1}}</span>-->
          <span>域名:api.airflow.adtech.sogou</span><br/>
          <iframe id="show-iframe1" :src='logUrl1'
                  frameborder="no"
                  marginheight="0" marginwidth="0"
                  scrolling="yes"
                  style="height: 660px;width: 100%"/>
        </el-col>

        <el-col :span="11" :offset="1">
<!--          <span>{{logUrl2}}</span>-->
          <span>域名: airflow.adtech.sogou</span><br/>
          <iframe id="show-iframe2" :src='logUrl2'
                  frameborder="no"
                  marginheight="0" marginwidth="0"
                  scrolling="yes"
                  style="height: 660px;width: 100%"/>
        </el-col>
      </el-row>
    </el-dialog>
  </div>
</template>

<style>
.el-table .odd-row {
  background: #fff;
}

.el-table .even-row {
  background: #ebeef5;
}
</style>
<script>
const axios = require('axios')
export default {
  data() {
    return {
      stateOption: {
        success: '成功',
        danger: '失败',
        notrun: '等待执行',
      },
      taskCnt: 1,
      baseDate: "",
      tableForm: {
        loading: true,
        data: [],
        cols: []
      },
      dagId: '',
      cnt: 0,
      searchWord: '',
      tableFormData: [],
      taskStateOptVisible: false,
      dagStateOptVisible: false,
      curTask: {
        index: -1,
        taskId: '',
        dagName: '',
        execution_date: '',
      },
      logUrl1: '',
      logUrl2: '',
      logVisible: false,
    }
  },
  methods: {
    checkLog(row) {
      let executionDate = row.execution_date
      let taskId = row.task_id
      let url1 = "http://api.airflow.adtech.sogou:8282/" + this.dagId + "/" + taskId + "/" + encodeURIComponent(executionDate)
      let url2 = "http://airflow.adtech.sogou:8282/" + this.dagId + "/" + taskId + "/" + encodeURIComponent(executionDate)
      this.logUrl1 = url1
      this.logUrl2 = url2
      this.logVisible = true
    },
    onChangeSearch() {
      if (this.searchWord === "") {
        this.tableForm.data = this.tableFormData
        return
      }
      this.tableForm.data = this.tableFormData.filter(
        x => x.task_id === this.searchWord
      )
      // console.log(this.tableForm.data[0].task_id)
    },
    loadTaskData() {
      let self = this
      this.searchWord = ""
      // 获取1个dag中的任务个数
      /*axios.get('airflow/getTaskCntInOneDag/' + self.dagId)
        .then(function (response) {
          self.taskCnt = response.data.taskCount

        }).catch(function (error) {
        console.log(error)
      })*/
      // 获取task列表
      let param = {
        "dagId": this.dagId,
        "baseDate": this.baseDate,
        "cnt": this.cnt + "",
        "taskId": "",
      }
      axios.post('airflow/getDetailByDagId/', param, {responseType: 'json'})
        .then(function (response) {
          self.tableForm.cols = response.data.cols
          self.tableForm.data = response.data.data
          self.tableFormData = response.data.data

          self.tableForm.loading = false
        }).catch(function (error) {
        console.log(error)
      })
      // console.log()
    },
    getTableRowClassName({row, rowIndex}) {
      // let groupCnt = Math.floor(rowIndex / this.taskCnt)
      if (row.epoch % 2 === 1) {  // 判断奇偶
        return 'odd-row'
      } else {
        return 'even-row'
      }
    },
    getTaskState(state) {
      if (state === "") {
        return "running"
      } else {
        return state
      }
    },
    MarkTaskSheduleStatus(newstatus) {
      let self = this
      let params = {
        'dag_name': self.curTask.dagName,
        'task_id': self.curTask.taskId,
        'execution_date': self.curTask.execution_date,
        'new_status': newstatus,
      }
      let idx = self.curTask.index

      axios.post('/airflow/markTaskState', params, {responseType: 'json'})
        .then(function (response) {
          let res = response.data

          if (res["flag"] == false) {
            self.$message({
              type: 'error',
              message: '更新失败',
            })
          } else {
            self.tableForm.data[idx]['state'] = newstatus
            self.$message({
              type: 'success',
              message: '更新成功！',
            })
          }
          self.taskStateOptVisible = false
        })
        .catch(function (error) {
          console.log(error)
          self.taskStateOptVisible = false
        })

    },
    MarkDagSheduleStatus(newstatus) {
      let self = this
      let params = {
        'name': self.curTask.dagName,
        'execution_date': self.curTask.execution_date,
        'new_status': newstatus,
      }

      let idx = self.curTask.index
      axios.post('/airflow/markdagstate', params, {responseType: 'json'})
        .then(function (response) {
          let res = response.data
          if (res["flag"] == false) {
            self.$message({
              type: 'error',
              message: '更新失败',
            })
          } else {
            self.$message({
              type: 'success',
              message: '更新成功！',
            })
            for (var i = 0; i < self.tableForm.data.length; i++) {
              if (self.tableForm.data[i]['execution_date'] === self.curTask.execution_date) {
                self.tableForm.data[i]['state'] = ""
              }

            }

          }
          self.dagStateOptVisible = false
        })
        .catch(function (error) {
          console.log(error)
          self.dagStateOptVisible = false
        })
    },

    TaskStateOpt(index, taskinfo) {
      let self = this
      self.curTask.taskId = taskinfo.task_id
      self.curTask.dagName = self.dagId
      self.curTask.execution_date = taskinfo.execution_date
      self.curTask.index = index

      self.taskStateOptVisible = true
    },
    DagStateOpt(index, taskinfo) {
      let self = this
      self.curTask.taskId = taskinfo.task_id
      self.curTask.dagName = self.dagId
      self.curTask.execution_date = taskinfo.execution_date
      self.curTask.index = index

      self.dagStateOptVisible = true
    },
  },
  mounted() {
    // this.dagId = 'experiment'
    let d = new Date()
    this.baseDate = d.getFullYear() + '-' + (d.getMonth() + 1) + '-' + d.getDate()
    // console.log("当前的基准日期: " + this.baseDate)

    this.dagId = sessionStorage.getItem('jones_airflow_dagId_show')
    // console.log("当前dagId:" + this.dagId)
    this.cnt = 5
    this.loadTaskData()
    // sessionStorage.removeItem('jones_airflow_dagId_show');
  }
}
</script>

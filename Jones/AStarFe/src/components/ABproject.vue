<template>
  <div>
    <h3>待上线项目</h3>
    <el-row>
      <el-table
        v-loading="readyLoading"
        :data="readyData"
        stripe
        style="width: 100%"
        :default-sort = "{prop:'', order:'descending'}">
        <el-table-column
          v-for="readyCol in readyCols"
          v-if="!readyCol.hidden"
          :key="readyCol.prop"
          :prop="readyCol.prop"
          :label="readyCol.label"
          :sortable="readyCol.sortable">
        </el-table-column>
        <el-table-column
          fixed="right"
          label="操作"
          width="200">
          <template slot-scope="scope">
            <el-button @click="projectOnline(scope.row)" type="text" size="mini">上线</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-row>

    <h3>上线项目</h3>
    <el-row>
      <el-table
        v-loading="projectLoading"
        :data="projectData"
        stripe
        style="width: 100%"
        :default-sort = "{prop:'lastUpdate', order:'descending'}">
        <el-table-column
          v-for="projectCol in projectCols"
          v-if="!projectCol.hidden"
          :key="projectCol.prop"
          :prop="projectCol.prop"
          :label="projectCol.label"
          :sortable="projectCol.sortable">
        </el-table-column>
        <el-table-column
          fixed="right"
          label="操作"
          width="200">
          <template slot-scope="scope">
            <el-button @click="projectClick(scope.row)" type="text" size="mini">查看数据</el-button>
            <el-button @click="projectResubmit(scope.row)" type="text" size="mini">查看配置</el-button>
            <el-button @click="projectOffline(scope.row)" type="text" size="mini">下线</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-row>

    <h3>下线项目</h3>
    <el-row>
      <el-table
        v-loading="offlineLoading"
        :data="offlineData"
        stripe
        style="width: 100%"
        :default-sort = "{prop:'', order:'descending'}">
        <el-table-column
          v-for="offlineCol in offlineCols"
          v-if="!offlineCol.hidden"
          :key="offlineCol.prop"
          :prop="offlineCol.prop"
          :label="offlineCol.label"
          :sortable="offlineCol.sortable">
        </el-table-column>
        <el-table-column
          fixed="right"
          label="操作"
          width="200">
          <template slot-scope="scope">
            <el-button @click="projectOnline(scope.row)" type="text" size="mini">重新上线</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-row>
  </div>

</template>

<script>
  const axios = require('axios');
  import router from '../router';
  import {isUserHimself} from '../conf/utils';

  export default {
    data(){
      return{
        readyLoading: false,
        readyData: [],
        readyCols: [],
        projectLoading: false,
        projectData: [],
        projectCols: [],
        offlineLoading: false,
        offlineData: [],
        offlineCols: []
      }
    },
    methods:{
      projectClick(row){
        sessionStorage.setItem('abtest_info', JSON.stringify(row))
        router.push({path: 'abdata'})
      },
      projectResubmit(row) {
        let self = this;
        const taskId = row.task_id;
        const taskType = row.task_type;
        axios.get('/small/log/query/' + taskId)
          .then(function (response) {
            sessionStorage.setItem('abtest_resubmit', JSON.stringify(response.data));
            if (taskType !== "无线小流量任务"){
              router.push({path: 'abpc'});
            } else {
              router.push({path: 'abwap'})
            }
          })
            .catch(function (error) {
            self.errorNotify(error)
          })

      },
      loadReadyData (){
        let self = this;
        self.readyLoading = true;
        axios.get('/small/log/ready')
          .then(function(response) {
            self.readyLoading = false;
            self.readyCols = response.data.cols;
            self.readyData = response.data.data;
          });
      },
      loadProjectData (){
        let self = this;
        self.projectLoading = true;
        axios.get('/small/log/enable')
          .then(function(response) {
            self.projectLoading = false;
            self.projectCols = response.data.cols;
            self.projectData = response.data.data;
          });
      },
      loadOfflineData (){
        let self = this;
        self.offlineLoading = true;
        axios.get('/small/log/disable')
          .then(function(response) {
            self.offlineLoading = false;
            self.offlineCols = response.data.cols;
            self.offlineData = response.data.data;
          });
      },
      loadData(){
        this.loadReadyData();
        this.loadProjectData();
        this.loadOfflineData();
      },
      errorNotify(msg) {
        this.$notify.error({
          title: '错误',
          message: msg
        });
      },
      projectOnline(row) {
        let self = this;
        const taskId = row.task_id;
        const session = row.session;
        if (isUserHimself(session)) {
          axios.get('/small/task/online/' + taskId)
            .then(function () {
              self.loadData();
            })
            .catch(function (error) {
              self.errorNotify(error)
            })
        } else{
          this.errorNotify("您不是项目所有者，无权操作！")
        }

      },
      projectOffline(row) {
        let self = this;
        const taskId = row.task_id;
        const session = row.session;
        if (isUserHimself(session)) {
          axios.get('/small/task/offline/' + taskId)
            .then(function () {
              self.loadData();
            })
            .catch(function (error) {
              self.errorNotify(error);
            })
        } else {
          this.errorNotify("您不是项目所有者，无权操作！")
        }
      }
    },
    mounted: function () {
      this.loadData();
    },
  }
</script>

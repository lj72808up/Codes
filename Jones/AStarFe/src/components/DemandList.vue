<template>
  <div>
    <el-button type="primary" @click="createDemand" size="small">创建需求</el-button>
    <el-table
      :data="demandList"
      stripe
      style="width: 100%">
      <el-table-column
        v-for="col in cols"
        v-if="!col.hidden"
        :key="col.prop"
        :prop="col.prop"
        :label="col.label"
        :sortable="true">
        <template slot-scope="scope">
          <span v-if="col.prop==='status'">{{getDevStatus(scope.row.status)}}</span> <!--转换状态-->
          <span v-else>{{scope.row[col.prop]}}</span>
        </template>
      </el-table-column>
      <el-table-column key="operate" prop="operate" label="操作" width="100px" align="center">
        <template slot-scope="scope">
          <el-button icon="el-icon-setting" type="text" @click="showDemandInfo(scope.row.rid)" size="mini">查看
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
  import router from "../router";

  const axios = require('axios');
  export default {
    name: "DemandList",
    data() {
      return {
        demandList: [],
        cols: [],
        devStatusMapping: {},
      }
    },
    methods: {
      createDemand() {
        this.showDemandInfo("new")
      },
      getAllDemand() {
        let self = this;
        axios.get('/requirement/getAllDemand')
          .then(function (response) {
            self.demandList = response.data.data;
            self.cols = response.data.cols
          }).catch(function (error) {
          console.log(error)
        })
      },
      getAllStatusMapping() {
        let self = this;
        axios.get('/requirement/GetDevStatusMapping')
          .then(function (response) {
            self.devStatusMapping = response.data
          }).catch(function (error) {
          console.log(error)
        })
      },
      getDevStatus(status) {
        return this.devStatusMapping[status]
      },
      showDemandInfo(rid) {
        sessionStorage.setItem('demandId', rid + "")
        router.push({name: 'demand', params: {demandId: rid+""}})
      }
    },
    mounted() {
      this.getAllDemand();
      this.getAllStatusMapping()
    }
  }
</script>

<style scoped>

</style>

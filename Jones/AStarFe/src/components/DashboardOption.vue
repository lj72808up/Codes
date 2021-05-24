<template>
  <div>
    <h3>DashBoard 新增</h3>
    <el-form inline="true" label-width="120px">
      <el-form-item label="DashBoard名称 ">
        <el-input  v-model.trim="dashboardName" size="small"/>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onInsertDashboard()" size="small">添加</el-button>
      </el-form-item>
    </el-form>

    <el-form v-show="false" inline="true" label-width="120px">
      <el-form-item label="DashboardURL">
        <el-input  v-model="dashboardURL" size="big" style="width:800px"/>
      </el-form-item>
    </el-form>

    <h3>DashBoard 管理</h3>
    <el-row>
      <el-table :data="redashData" stripe style="width: 100%">
        <el-table-column prop="name" label="DashBoard名称" min-width="50%">
          <template slot-scope="scope">
            <el-input v-if="(modifyStatus && scope.$index == modifyIdx)" v-model="modifyNewName"></el-input>
            <span v-else>{{scope.row.name}}</span>  
          </template>
        </el-table-column>      
        <el-table-column v-if=showURl prop="url" label="DashBoard_URL" min-width="50%">
        </el-table-column>        
        <el-table-column prop="index" label="DashBoard索引" min-width="50%">
        </el-table-column> 
        <el-table-column prop="creator" label="Creator" min-width="50%">
        </el-table-column> 
        <el-table-column label="操作">
          <template slot-scope="rdata">
            <el-button @click="assignUser(rdata.row)" type="info" size="mini">指派</el-button>
            <el-button v-show="(rdata.$index != modifyIdx)" @click="ModifyDashBoard(rdata.$index,rdata.row.index,rdata.row.name)" type="danger" size="mini">修改</el-button>
            <el-button v-show="(rdata.$index == modifyIdx)" @click="UpdateDashBoard(rdata.$index,rdata.row.index)" type="success" icon="el-icon-check" circle></el-button>
            <el-button v-show="(rdata.$index == modifyIdx)" @click="CancelUpdate()" type="info" icon="el-icon-close" circle></el-button>
            <el-button @click="DeleteDashBoard(rdata.$index,rdata.row.index)" type="danger" size="mini">删除</el-button>
            <el-button @click="JumpDashBoard(rdata.$index,rdata.row.index,rdata.row.url)" type="success" size="mini">跳转</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-row>

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


<script>
  const axios = require('axios')
  import {getCookie,setCookie,setRootCookie, isEmpty} from '../conf/utils'


  export default {
    data () {
      return {
        redashData: [],
        redashURL: '',
        dashboardIndex: '',

        dashboardName: '',
        dashboardURL: '',
        dashboardRefresh: 120,

        newLid: -1,
        privilegeRoles:[],
        roleOptions: [],
        parentNode : '数据图表',

        modifyNewName : '',
        modifyStatus : false,
        modifyIdx : -1,
        showURl: false,

        assignUserVisible : false,
        assignIndex : "",
        assignUserName : "",
      }
    },
    methods: {
        onInsertDashboard:async function () {
          let self = this
          let mergeUrl = ""

          if(this.dashboardName == ""){
            this.$alert("名称不可为空")
            return 
          }
          let params = {
            'dashboardname': this.dashboardName,
            'usrid' : getCookie('_adtech_user')
          }

          //添加新dashboard表
          await axios.post('/redash/dashboardinsert', params, {responseType: 'json'})
            .then( (response) => {
              self.newLid = response.data

              if(self.newLid != -1){
                self.$message({
                  message: '添加成功！',
                  type: 'success'
                })
                sessionStorage.setItem(params['dashboardIndex'], params['dashboardURL'])
                this.GetAllDashboard()
              }
              else{
                  self.$message({
                    message: '添加失败！',
                    type: 'error'
                  })
              }
            })
            .catch( (error) => {
              self.$message({
                message: '添加失败！' + error.toString(),
                type: 'error'
              })
          })
        },   
        JumpDashBoard (rowIdx,dashboardIndex,dashboardUrl) {
          //router.push({path: 'dag'})
          this.$router.push({name: 'dashboardshow', params: {dashboardname: dashboardIndex,dashboardUrl:dashboardUrl}})
        },

        GetAllDashboard () {
          let self = this
          axios.get('/redash/dashboardall')
            .then(function (response) {
              var res = []
              response.data.forEach(ele => {
                console.log(ele);
                res.push(ele)
              })
              self.redashData = res
              console.log("------------------")
              console.log(self.redashData)
            }).catch(function (error) {
            console.log(error)
          })
        },
        ModifyDashBoard (rowIdx,dashboardIndex,dashboardName) {
          let self = this
          if (this.modifyStatus == true) {
              self.$message({
                message: '正在编辑...',
                type: 'error'
              })
              return 
          }
          this.modifyNewName = dashboardName
          this.modifyStatus = true
          this.modifyIdx = rowIdx 
        },
        UpdateDashBoard(rowIdx,dashboardIndex) {
          let self = this

          let params = {
            'dashboardIndex': dashboardIndex,
            'dashboardName' : this.modifyNewName,
          }
          console.log("begin update this dashboard : " + dashboardIndex)
          axios.post('/redash/dashboardupdate', params, {responseType: 'json'})
            .then( (response) => {
              self.$message({
                message: '更新成功！',
                type: 'success'
              })
              self.CancelUpdate()
              self.GetAllDashboard()
            })
            .catch( (error) => {
              self.$message({
                message: '更新失败！' + error.toString(),
                type: 'error'
              })
            })           
        }, 
        CancelUpdate() {
          let self = this
          this.modifyNewName = ""
          this.modifyStatus = false
          this.modifyIdx = -1 
        },                     
        DeleteDashBoard (rowIdx,dashboardIndex) {
          let self = this
          let cancelFlag = false

          let params = {
            'dashboardIndex': dashboardIndex
          }

          if(confirm("确定删除吗?")==false){
            return 
          }

          console.log("begin delete this redash : " + dashboardIndex)
          axios.post('/redash/dashboarddel', params, {responseType: 'json'})
            .then( (response) => {
              self.$message({
                message: '删除成功！',
                type: 'success'
              })
              self.GetAllDashboard()
              self.redashData.splice(rowIdx,1)
              console.log(self.redashData)
            })
            .catch( (error) => {
              self.$message({
                message: '删除失败！' + error.toString(),
                type: 'error'
              })
            })  
        },
        assignUser (row) {
          console.log(row)
          this.assignUserVisible = true
          this.assignIndex = row.index
        },
        onClickAssign () {
          let param = {
            'assignUserName': this.assignUserName,
            'assignIndex': this.assignIndex,
            'assignType': "dashboard",
          }
          console.log(param)
          let self = this
          axios.post('/redash/user_change', param, {responseType: 'json'})
            .then(function (response) {
              let res = response.data
              if (res.Res) {
                self.$message({
                  message: '指派成功',
                  type: 'success'
                })
                self.GetAllDashboard()
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
    },

    created: function () {
      let self = this
      axios.get('/redash/timestamp')
        .then(function (response) {
          self.timestamp = response.data
          self.sec186 = (self.timestamp + 7482952) * 3
          setRootCookie('sec186', self.sec186.toString(), 1)
        }).catch(function (error) {
        console.log(error)
      })

      axios.get('/redash/write_access')
          .then(function (response) {
            self.showURl = response.data['dashboard_write']
          }).catch(function (error) {
          console.log(error)
        })
      this.GetAllDashboard()
    },
    mounted () {
      console.log('mounted');

    }
  }
</script>

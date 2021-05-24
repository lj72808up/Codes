<template>
  <div>
    <h3>新增绿皮</h3>

    <el-form ref="form">
      <el-row style="margin-bottom: -10px">
        <el-form-item prop="name">
          <el-col :span="4">
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
          <el-col :span="6" :offset="1">
            <el-select v-model="datasource.selectValue" placeholder="请选择数据源" size="small"
                        @change="onChangeDataSource" clearable filterable>
              <el-option
                v-for="item in datasource.options"
                :key="item.dsid"
                :label="item.DatasourceName"
                :value="item.Dsid">
              </el-option>
            </el-select>
          </el-col>
          <el-col :span="6">
            <span>编码</span>
            <el-select v-model="encode_value" placeholder="编码">
              <el-option
                v-for="item in encode_options"
                :key="item.value"
                :label="item.label"
                :value="item.value">
              </el-option>
            </el-select>            
          </el-col>
          <el-col :span="4">
            <el-checkbox v-model="date_pianyi_1">日期偏移+1</el-checkbox>
          </el-col>
        </el-form-item>
      </el-row>


      <el-row style="margin-bottom: 50px">
        <el-tree :data="treeData" :props="defaultProps"  ref="tree">
          <span class="el-tree-node__content" slot-scope="{ node, data }">
            <span @click="onClickTree(data)">{{ node.label }}</span>
          </span>
        </el-tree>
      </el-row>
      <el-row style="margin-left: 40%">
        <span>绿皮名称</span>
        <el-input  v-model="lvpi_name" size="big" style="width:300px"/>
        <el-button type="primary" size="small" @click="onCreateSqlView()">创建绿皮</el-button>
      </el-row>      
    </el-form>
  
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

        treeData: '',
        datasource: {
          selectValue: '',
          options: [],
          mysqlOptions: [],
          typeOptions: [{'name': 'mysql'},  {'name': 'clickhouse'},],
          typeSelectValue: '',
        },
        defaultProps: {
          children: 'children',
          label: 'label'
        },   
        explorerTable: {
          tableId: '',
          tableName: '',
        }, 
        exploreTableName : '',
        date_pianyi_1 : false,

        encode_options: [{
          value: 'UTF8',
          label: 'UTF8'
        }, {
          value: 'LATIN1',
          label: 'LATIN1'
        }],
        encode_value: 'LATIN1',
        lvpi_name: '',           
      }
    },
    methods: {  
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
        onChangeDSType () {
          let self = this
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
        },
        onClickTree (node) {
          let self = this
          if (node.tableName !== '') {
            self.explorerTable.tableName = node.label
            self.explorerTable.tableId = node.tid
            self.exploreTableName = node.tableName

          }
        }, 
        onCreateSqlView () {
          let self = this
          console.log(self.exploreTableName)
          if (self.exploreTableName === "") {
            self.$message({
              type: 'error',
              message: "请点击表名文字选中表!"
            })
            return          
          }
          if (self.lvpi_name.trim() === ""){
            self.$message({
              type: 'error',
              message: "请输入绿皮名称!"
            })
            return  
          }

          let params = new URLSearchParams()
          let ds_name = ""
          for (let idx=0; idx<self.datasource.options.length; idx++) {
            var ds = self.datasource.options[idx]
            if (self.datasource.selectValue == ds["Dsid"]) {
              ds_name = ds["DatasourceName"].split("(")[0]
            }
          }
          console.log(ds_name)
          console.log(self.explorerTable)
          console.log("1232222222222222222222222")

          params.append('tablename', self.exploreTableName)
          params.append('tableid',self.explorerTable.tableId)
          params.append('encode', this.encode_value)
          params.append('date_pianyi',this.date_pianyi_1)
          params.append('sqlType', this.datasource.typeSelectValue)
          params.append('dsId', this.datasource.selectValue)
          //params.append('name','绿皮_'+ds_name+"_"+ self.exploreTableName.replace(".","_"))
          params.append('name',self.lvpi_name)


          self.$message({
            type: 'info',
            message: "绿皮创建中..."
          })        
          axios.post('/jdbc/sql/createLvpiView', params, {responseType: 'json'})
            .then(function (response) {
              if (response.data.msg.indexOf('成功') >= 0) {
                self.curRedashQueryId = response.data.redash_query_id
                //let queryUrl = 'http://inner-redash.adtech.sogou/queries/' + self.curRedashQueryId + '/source'
                let queryUrl = 'http://inner-redash.adtech.sogou/queries/' + self.curRedashQueryId + '/'
                self.$router.push({name: 'redashqueryshow', params: {queryid: self.curRedashQueryId, queryUrl: queryUrl}})
              } else {
                self.$message({
                  type: 'error',
                  message: response.data.msg
                });
              }

            })
            .catch(function (error) {
              console.log(error)
            })
            
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

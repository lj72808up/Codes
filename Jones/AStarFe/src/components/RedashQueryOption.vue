<template>
  <div>
    <h3>Redash Query 管理</h3>
    <el-row>
      <el-table :data="redashData" stripe style="width: 100%">
        <el-table-column prop="name" label="RedashQuery名称">
        </el-table-column>
        <el-table-column v-if=showURl prop="url" label="RedashQuery_URL">
        </el-table-column>        
        <el-table-column prop="index" label="RedashQuery索引">
        </el-table-column> 
        <el-table-column prop="creator" label="Creator">
        </el-table-column> 
        <el-table-column label="操作">
          <template slot-scope="rdata">
            <el-button @click="assignUser(rdata.row)" type="info" size="mini">指派</el-button>
            <el-button @click="DeleteRedashQuery(rdata.$index,rdata.row.index)" type="danger" size="mini">删除</el-button>
            <el-button @click="JumpQueryView(rdata.$index,rdata.row.index,rdata.row.url)" type="success" size="mini">跳转</el-button>
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
        queryIndex: '',
        showURl: false,

        assignUserVisible : false,
        assignIndex : "",
        assignUserName : "",
      }
    },
    methods: {
        DeleteRedashQuery (rowIdx,queryIndex) {
          let self = this
          let cancelFlag = false

          let params = {
            'queryIndex': queryIndex
          }

          if(confirm("确定删除吗?")==false){
            return 
          }

          console.log("begin delete this redash : " + queryIndex)
          axios.post('/redash/redashquerydel', params, {responseType: 'json'})
            .then( (response) => {
              self.$message({
                message: '删除成功！',
                type: 'success'
              })
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
        JumpQueryView (rowIdx,queryIndex,queryUrl) {
          console.log(queryIndex)
          let queryid = queryIndex.split('/')[1]
          console.log(queryid)
          //router.push({path: 'dag'})
          this.$router.push({name: 'redashqueryshow', params: {queryid: queryid,queryUrl:queryUrl}})
        },
        GetAllQuery() {
          let self = this
          axios.get('/redash/redashqueryall')
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
        assignUser (row) {
          console.log(row)
          this.assignUserVisible = true
          this.assignIndex = row.index
        },
        onClickAssign () {
          let param = {
            'assignUserName': this.assignUserName,
            'assignIndex': this.assignIndex,
            'assignType': "query",
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
                self.GetAllQuery()
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
          self.showURl = response.data['redashquery_write']
        }).catch(function (error) {
        console.log(error)
      })  
      self.GetAllQuery()
    },
    mounted () {
      console.log('mounted');

    }
  }
</script>

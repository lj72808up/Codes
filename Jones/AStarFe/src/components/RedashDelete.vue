<template>
  <div>
    <h3>删除图表</h3>
    <el-row>
      <el-table :data="redashData" stripe style="width: 100%">
        <el-table-column prop="title" label="图表名称">
        </el-table-column>
        <el-table-column prop="info" label="图表URL">
        </el-table-column>        
        <el-table-column prop="index" label="图标索引">
        </el-table-column> 
        <el-table-column label="操作">
          <template slot-scope="rdata">
            <el-button @click="DeleteRedash(rdata.$index,rdata.row.index)" type="text" size="mini">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-row>

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
        redashIndex: '',
      }
    },
    methods: {
        DeleteRedash (rowIdx,redashIndex) {
          let self = this
          let cancelFlag = false

          let params = {
            'redashIndex': redashIndex
          }

          if(confirm("确定删除吗?")==false){
            return 
          }

          console.log("begin delete this redash : " + redashIndex)
          axios.post('/redash/redashdel', params, {responseType: 'json'})
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
    },

    created: function () {
      let self = this
      axios.get('/redash/redashall')
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
    mounted () {
      console.log('mounted');

    }
  }
</script>

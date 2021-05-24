<template>
  <div>
    <h3>新增标签</h3>
    <el-form  label-width="120px">
      <el-form-item label="请输入标签名">
        <el-col :span="8">
          <el-input v-model="newLabelTitle" size="small"/>
        </el-col>
      </el-form-item>
      <el-form-item label="请输入URL索引">
        <el-col :span="8">        
          <el-input v-model="newLabelIndex" size="small"/>
        </el-col>
      </el-form-item>    
      <el-form-item inline="true" label="请输入命名空间">
        <el-col :span="8">        
          <el-input v-model="newLabelNS" size="small"/>
        </el-col>
        <el-button type="primary" @click="onAddLabel()" size="small">添加</el-button>
      </el-form-item>          

    </el-form>

    <h3>配置标签关联</h3>
    <el-row>
      <div style="max-height: 600px ; overflow:auto">
        <el-table
          :data="labelData"
          ref="labelTable"
          row-key="labelId"
          stripe
        >
          <el-table-column property="labelId" label="标签id"/>
          <el-table-column prop="labelTitle2" label="标签名称"/>
          <el-table-column prop="labelIndex2" label="标签索引">
            <template slot-scope="scope">
              <el-input v-if="scope.row.status" v-model="scope.row.labelIndex2"></el-input>
              <span v-else>{{scope.row.labelIndex2}}</span>
            </template>
          </el-table-column>          
          <el-table-column prop="labelNameSpace2" label="标签命名空间">
            <template slot-scope="scope">
              <el-input v-if="scope.row.status" v-model="scope.row.labelNameSpace2"></el-input>
              <span v-else>{{scope.row.labelNameSpace2}}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="rdata">
              <el-button @click="UpdateLabel(rdata.row)" type="text" size="mini">修改</el-button>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="rdata">
              <el-button @click="CancelUpdateLabel(rdata.row)" type="text" size="mini">放弃修改</el-button>
            </template>
          </el-table-column>          
          <el-table-column label="操作">
            <template slot-scope="rdata">
              <el-button @click="SaveLabel(rdata.row)" type="text" size="mini">保存</el-button>
            </template>
          </el-table-column>                    
        </el-table>
      </div>
    </el-row>
  </div>
</template>

<script>
  const axios = require('axios')
  import {getCookie, isEmpty} from '../conf/utils'

  export default {
    data () {
      return {
          labelData : [],
          newLabelTitle: '',
          newLabelIndex: '',
          newLabelNS: '',
      }
    },
    methods: {
      onAddLabel() {
        let self = this
        let params = {
          'title': self.newLabelTitle,
          'index': self.newLabelIndex,
          'namespace':self.newLabelNS,
        }

        axios.post('/auth/labels/rdadd', params, {responseType: 'json'})
          .then(function (response) {
            let newLid = response.data
            console.log(newLid)
            console.log("==================")
            
            self.labelData.push({
              'labelId': newLid,
              'labelTitle': self.newLabelTitle,
              'labelIndex': self.newLabelIndex,
              'labelNameSpace': self.newLabelNS,
              'labelTitle2': self.newLabelTitle,
              'labelIndex2': self.newLabelIndex,
              'labelNameSpace2': self.newLabelNS,
              'status': false,
            })

            self.$message({
              type: 'success',
              message: '添加成功,标签id' + newLid.toString() +'!'
            })

            self.newLabelTitle = ''
            self.newLabelIndex = ''
            self.newLabelNS = ''
            console.info(response.data)
          })
          .catch(function (error) {
            console.log(error)
          })        
      },
      UpdateLabel(curLabel){
        curLabel.status = true
      },
      CancelUpdateLabel(curLabel){
        curLabel.labelTitle2 = curLabel.labelTitle,
        curLabel.labelIndex2 = curLabel.labelIndex,
        curLabel.labelNameSpace2 = curLabel.labelNameSpace,
        curLabel.status = false
      },
      SaveLabel(curLabel){
        let self = this
        if(curLabel.labelIndex2 == curLabel.labelIndex && curLabel.labelNameSpace2 == curLabel.labelNameSpace )
        {
            self.$message({
              type: 'success',
              message: '修改前后未变化!'
            })
            return
        }

        let params = {
          'lid': curLabel.labelId.toString(),
          'index':curLabel.labelIndex2,
          'namespace':curLabel.labelNameSpace2,
        }

        axios.post('/auth/labels/rdupdate', params, {responseType: 'json'})
          .then(function (response) {
            curLabel.labelIndex = curLabel.labelIndex2,
            curLabel.labelNameSpace = curLabel.labelNameSpace2,
            curLabel.status = false,          

            self.$message({
              type: 'success',
              message: '更新成功!'
            })
            console.info(response.data)
          })
          .catch(function (error) {
            console.log(error)
          })
      },      
    },
    created: function () {
      let self = this
      axios.get('/auth/labels/rdlabels')
        .then(function (response) {
          let nameRes = []
          let checkedIdx = []
          let index = 0
          response.data.forEach(ele => {
            nameRes.push({
              'labelId': ele.Lid,
              'labelTitle': ele.Title,
              'labelIndex': ele.Index,
              'labelNameSpace': ele.NameSpace,
              'labelTitle2': ele.Title,
              'labelIndex2': ele.Index,
              'labelNameSpace2': ele.NameSpace,
              'status': false,
            })
          })
          self.labelData = nameRes

        }).catch(function (error) {
        console.log(error)
      })
    },
    mounted () {
      console.log('mounted')
    }
  }
</script>

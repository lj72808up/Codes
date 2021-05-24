<template>
  <div>
    <h3>新增dim</h3>
    <el-form label-width="150px">
      <el-form-item label="广告样式ID">
        <el-col :span="8">
          <el-input v-model="newAdStyleId" size="small"/>
        </el-col>
      </el-form-item>
      <el-form-item label="一级广告样式名称">
        <el-col :span="8">
          <el-input v-model="newFirstStyleName" size="small"/>
        </el-col>
      </el-form-item>
      <el-form-item inline="true" label="二级广告样式名称">
        <el-col :span="8">
          <el-input v-model="newSecondStyleName" size="small"/>
        </el-col>
      </el-form-item>
      <el-form-item inline="true" label="三级广告样式名称">
        <el-col :span="8">
          <el-input v-model="newThirdStyleName" size="small"/>
        </el-col>
        <el-button type="primary" @click="onAddLabel()" size="small">添加</el-button>
      </el-form-item>
    </el-form>

    <h3>配置dim</h3>
    <el-row>
      <div style="max-height: 600px ; overflow:auto">
        <el-table
          :data="labelData"
          ref="labelTable"
          row-key="id"
          stripe
        >
          <el-table-column property="Id" label="Id"/>
          <el-table-column prop="AdStyleId2" label="一级广告样式名称">
            <template slot-scope="scope">
              <el-input v-if="scope.row.status" v-model="scope.row.AdStyleId2"></el-input>
              <span v-else>{{ scope.row.AdStyleId2 }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="FirstStyleName2" label="一级广告样式名称">
            <template slot-scope="scope">
              <el-input v-if="scope.row.status" v-model="scope.row.FirstStyleName2"></el-input>
              <span v-else>{{ scope.row.FirstStyleName2 }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="SecondStyleName2" label="二级广告样式名称">
            <template slot-scope="scope">
              <el-input v-if="scope.row.status" v-model="scope.row.SecondStyleName2"></el-input>
              <span v-else>{{ scope.row.SecondStyleName2 }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="ThirdStyleName2" label="三级广告样式名称">
            <template slot-scope="scope">
              <el-input v-if="scope.row.status" v-model="scope.row.ThirdStyleName2"></el-input>
              <span v-else>{{ scope.row.ThirdStyleName2 }}</span>
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
              <el-button @click="DelLabel(rdata.row)" type="text" size="mini">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-row>
  </div>
</template>

<script>
const axios = require('axios')

export default {
  data() {
    return {
      labelData: [],
      newAdStyleId: '',
      newFirstStyleName: '',
      newSecondStyleName: '',
      newThirdStyleName: ''
    }
  },
  methods: {
    onAddLabel() {
      let self = this
      let params = {
        "AdStyleId": parseInt(this.newAdStyleId),
        "ThirdStyleName": this.newThirdStyleName,
        "SecondStyleName": this.newSecondStyleName,
        "FirstStyleName": this.newFirstStyleName
      }
      axios.post('/single/addDim', params, {responseType: 'json'})
        .then(function (response) {
          self.$message({
            type: 'success',
            message: '添加成功',
          })
          self.initTable()
          self.newAdStyleId = ''
          self.newFirstStyleName = ''
          self.newSecondStyleName = ''
          self.newThirdStyleName = ''
        })
        .catch(function (error) {
          console.log(error)
        })
    },
    UpdateLabel(curLabel) {
      console.log("更改状态" + curLabel.status)
      curLabel.status = true
    },
    CancelUpdateLabel(curLabel) {
      curLabel.AdStyleId2 = curLabel.AdStyleId,
        curLabel.FirstStyleName2 = curLabel.FirstStyleName,
        curLabel.SecondStyleName2 = curLabel.SecondStyleName,
        curLabel.ThirdStyleName2 = curLabel.ThirdStyleName,
        curLabel.status = false
    },
    DelLabel(curLabel) {
      let self = this
      axios.get('/single/delDim/'+curLabel.Id, {responseType: 'json'})
        .then(function (response) {
          curLabel.status = false
          self.$message({
            type: 'success',
            message: '删除成功!'
          })
          console.info(response.data)
          self.initTable()
        })
        .catch(function (error) {
          console.log(error)
        })
    },
    SaveLabel(curLabel) {
      let self = this
      let params = {
        'Id': curLabel.Id,
        'AdStyleId': parseInt(curLabel.AdStyleId2),
        'FirstStyleName': curLabel.FirstStyleName2,
        'SecondStyleName': curLabel.SecondStyleName2,
        'ThirdStyleName': curLabel.ThirdStyleName2
      }

      axios.post('/single/updateDim', params, {responseType: 'json'})
        .then(function (response) {
            curLabel.status = false
            self.$message({
              type: 'success',
              message: '更新成功!'
            })
          console.info(response.data)
          self.initTable()
        })
        .catch(function (error) {
          console.log(error)
        })
    },
    initTable() {
      let self = this
      axios.get('/single/getDims')
        .then(function (response) {
          self.labelData = response.data
          let nameRes = []
          response.data.forEach(ele => {
            nameRes.push({
              'Id': ele.Id,
              'AdStyleId': ele.AdStyleId,
              'FirstStyleName': ele.FirstStyleName,
              'SecondStyleName': ele.SecondStyleName,
              'ThirdStyleName': ele.ThirdStyleName,
              'AdStyleId2': ele.AdStyleId,
              'FirstStyleName2': ele.FirstStyleName,
              'SecondStyleName2': ele.SecondStyleName,
              'ThirdStyleName2': ele.ThirdStyleName,
              'status': false,
            })
          })
          self.labelData = nameRes

        }).catch(function (error) {
        console.log(error)
      })
    }
  },
  created: function () {
    this.initTable()
  },
  mounted() {
    console.log('mounted')
  }
}
</script>

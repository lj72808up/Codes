<template>
  <div>
    <h3>新增图表展示</h3>
    <el-form inline="true" label-width="120px">
      <el-form-item label="图表名称">
        <el-input v-model="redashName" size="small"/>
      </el-form-item>
      <el-form-item label="上级域名">
        <el-select v-model="superiorDomain" filterable placeholder="请选择" size="small">
          <el-option
            v-for="domain in superiorDomainOptions"
            :key="domain.Did"
            :label="domain.Name"
            :value="domain.Did">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onInsertRedash()" size="small">添加</el-button>
      </el-form-item>
    </el-form>
    <el-form inline="true" label-width="120px">
      <el-form-item label="刷新时间(秒) ">
        <el-input placeholder="请输入秒:如120" v-model="redashRefresh" number type="number" maxlength="5" size="small"/>
      </el-form-item>
      <el-form-item label-width="200px" label="(置为空,表示不自动刷新)">
      </el-form-item>
    </el-form>
    <el-form inline="true" label-width="120px">
      <el-form-item label="图表来源URL">
        <el-input v-model="redashURL" size="big" style="width:800px"/>
      </el-form-item>
    </el-form>

    <el-form inline="true" label-width="120px">
      <el-form-item v-show="true" label="有权使用的角色:">
        <el-select v-model="privilegeRoles" multiple filterable placeholder="请选择" size="small">
          <el-option
            v-for="value in roleOptions"
            :key="value"
            :label="value"
            :value="value">
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>

    <el-form inline="true" label-width="120px">
      <el-form-item label="图表类型">
        <el-radio-group v-model="parentNode">
          <el-radio label="数据图表" border>数据图表</el-radio>
          <el-radio label="波动图表" border>波动图表</el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>


  </div>
</template>


<script>
const axios = require('axios')
import {getCookie, setCookie, setRootCookie, isEmpty} from '../conf/utils'


export default {
  data() {
    return {
      redashName: '',
      redashURL: '',
      redashIndex: '',
      redashRefresh: 3600,
      superiorDomain: '',  // 上级域名
      superiorDomainOptions: [],
      newLid: -1,
      privilegeRoles: [],
      roleOptions: [],
      parentNode: '数据图表',
    }
  },
  methods: {
    getRoleOptions() {
      let self = this
      axios.get('/auth/roles')
        .then(function (response) {
          self.roleOptions = response.data
        })
        .catch(function (error) {
          console.log(error)
        })
    },

    rndStr(len) {
      let text = " "
      let chars = "abcdefghijklmnopqrstuvwxyz"

      for (let i = 0; i < len; i++) {
        text += chars.charAt(Math.floor(Math.random() * chars.length))
      }
      return text
    },
    getSuperDomains(){
      let self = this
      let url = axios.defaults.baseURL.replace(/jones/g, '') + 'lvpi/getAllDomain'
      axios.get(url)
        .then(function (response) {
          self.superiorDomainOptions = response.data
        })
        .catch(function (error) {
          console.error(error)
        })
    },
    onInsertRedash: async function () {
      let self = this
      //todo  执行插入lvpi的接口

      let mergeUrl = ""
      let cancelFlag = false
      self.redashIndex = this.rndStr(8)

      if (this.redashRefresh == '') {
        mergeUrl = this.redashURL
      } else {
        mergeUrl = this.redashURL + "&refresh=" + this.redashRefresh
      }

      //添加此URL
      let params = {
        'redashname': this.redashName,
        'redashURL': mergeUrl,
        'redashIndex': this.redashIndex,
        'privilegeRoles': this.privilegeRoles.join(","),
        'type': this.parentNode,
        'superiorDomain': this.superiorDomain+"",
      }

      console.log("begin redash check -------------------")
      //查询是否URL重复
      await axios.post('/redash/redashcheck', params, {responseType: 'json'})
        .then((response) => {

          if (response.data['code'] != '0000') {
            let msg = "URL与[" + response.data['title'] + "]重复，确认添加吗?"
            if (confirm(msg) == false) {
              this.cancelFlag = true
            } else {
              this.cancelFlag = false
            }
          }
        })
        .catch((error) => {
          self.$message({
            message: '查询失败！' + error.toString(),
            type: 'error'
          })
        })

      if (this.cancelFlag == true) {
        return
      }

      console.log("begin redash insert -------------------")
      //添加新redash表
      await axios.post('/redash/redashinsert', params, {responseType: 'json'})
        .then((response) => {
          self.newLid = response.data

          if (self.newLid != -1) {
            self.$message({
              message: '添加成功！',
              type: 'success'
            })
            sessionStorage.setItem(params['redashIndex'], params['redashURL'])
          } else {
            self.$message({
              message: '添加失败！',
              type: 'error'
            })
          }
        })
        .catch((error) => {
          self.$message({
            message: '添加失败！' + error.toString(),
            type: 'error'
          })
        })
    },
  },

  created: function () {
    console.log('created')
  },
  mounted() {
    console.log('mounted')
    this.getRoleOptions()
    this.getSuperDomains()
  }
}
</script>

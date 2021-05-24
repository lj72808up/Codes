<template>
  <div>
    <el-row>
      <h3>连接管理</h3>
    </el-row>

    <el-row>
      <el-col :span="6">
        <el-select v-model="connSelectType" placeholder="请选择数据库类型"
                   @change="onChangeConnType" size="small" clearable>
          <el-option label="hive" value="hive"/>
          <el-option label="mysql" value="mysql"/>
          <el-option label="clickhouse" value="clickhouse"/>
          <el-option label="kylin" value="kylin"/>
        </el-select>
        <el-button size="small" type="primary" @click="triggerNew" v-show="!isReadOnly">新增</el-button>

        <div style="margin-top: 10px;">
          <el-tree :data="treeData"
                   :props="defaultProps"
                   :filter-node-method="filterNode" ref="tree"
                   @node-click="onClickTree"
                   accordion/>
        </div>
      </el-col>
      <el-col :span="1">
        <el-divider direction="vertical"/>
      </el-col>
      <el-col :span="8">
        <el-form ref="connForm" :model="connForm" :inline="false" label-width="100px">
          <el-form-item label="连接类型: " prop="connType">
            <el-select v-model="connForm.ConnType" placeholder="请选择连接类型"
                       :readonly="isReadOnly" size="small">
              <el-option label="hive" value="hive"/>
              <el-option label="mysql" value="mysql"/>
              <el-option label="clickhouse" value="clickhouse"/>
              <el-option label="kylin" value="kylin"/>
            </el-select>
          </el-form-item>

          <el-form-item label="连接名: " prop="DatasourceName">
            <el-input v-model="connForm.DatasourceName" size="small" :readonly="isReadOnly"/>
          </el-form-item>

          <el-form-item label="主机: " prop="Host">
            <el-input v-model="connForm.Host" size="small" :readonly="isReadOnly"/>
          </el-form-item>

          <el-form-item label="用户名: " prop="User">
            <el-input v-model="connForm.User" size="small" :readonly="isReadOnly"/>
          </el-form-item>

          <el-form-item label="密码: " prop="Passwd">
            <el-input v-model="connForm.Passwd" size="small" :readonly="isReadOnly"/>
          </el-form-item>

          <el-form-item label="数据库: " prop="Database">
            <el-input v-model="connForm.Database" size="small" :readonly="isReadOnly"/>
          </el-form-item>

          <el-form-item label="端口: " prop="Port">
            <el-input v-model="connForm.Port" size="small" :readonly="isReadOnly"/>
          </el-form-item>

          <el-form-item label="编码: " prop="Encode">
            <el-input v-model="connForm.Encode" size="small" :readonly="isReadOnly"/>
          </el-form-item>
        </el-form>

        <el-row type="flex" justify="start">
          <el-col :span="6" :offset="2">
            <el-button size="small" @click="onTest">测试连接</el-button>
            <i class="el-icon-loading" v-show="testLoading"/>
          </el-col>
          <el-col :span="2" :offset="14">
            <el-button size="small" @click="onSaveUpdate" type="primary" v-show="(!isReadOnly) && (connSelectId!=='')">
              保存修改
            </el-button>
            <el-button size="small" @click="onSaveNew" type="primary" v-show="(!isReadOnly) && (connSelectId==='')">
              提交新增
            </el-button>
          </el-col>
        </el-row>
        <el-row type="flex" justify="start">
          <el-col :offset="2">
            <div style="line-height: 20px; margin-top: 2px">
              <el-alert v-show="testShow && testSuccess"
                        title="测试成功"
                        type="success"
                        effect="light"
                        :closable="false"
                        show-icon/>
              <el-alert v-show="testShow && !testSuccess"
                        :title="testErrMsg"
                        type="error"
                        effect="light"
                        :closable="false"
                        show-icon/>
            </div>
          </el-col>
        </el-row>
      </el-col>

      <el-col :span="6" :offset="3">
        <span style="font-size: 14px;color: #606266;">有权使用的其它角色:</span>
        <el-checkbox-group v-model="privilegeSelect" size="mini">
          <el-checkbox v-for="role in roleOptions" :label="role" :key="role" border>{{role}}</el-checkbox>
        </el-checkbox-group>
      </el-col>
    </el-row>
  </div>

</template>

<script>
  const axios = require('axios')

  export default {
    data () {
      return {
        defaultProps: {
          children: 'children',
          label: 'DatasourceName'
        },
        connSelectType: '',  // 默认选中mysql
        connSelectId: '',         // 当前选中的数据源id
        treeData: [],
        connForm: {
          ConnType: '',
          DatasourceName: '',
          Host: '',
          User: '',
          Passwd: '',
          Database: '',
          Port: '',
          Encode: '',
          PrivilegeRoles: [],
        },
        testLoading: false,
        testShow: false,
        testSuccess: false,
        testErrMsg: '',
        isReadOnly: false,
        roleOptions: [],
        privilegeSelect: [],
      }
    },
    methods: {
      onSaveUpdate () {
        let self = this
        let param = {
          Dsid: this.connSelectId,
          ConnType: this.connForm.ConnType,
          DatasourceName: this.connForm.DatasourceName,
          Host: this.connForm.Host,
          User: this.connForm.User,
          Passwd: this.connForm.Passwd,
          Database: this.connForm.Database,
          Port: this.connForm.Port,
          Encode: this.connForm.Encode,
          PrivilegeRoles: this.privilegeSelect.join(','),
        }
        axios.post('/connMeta/updateConn', param, {responseType: 'json'})
          .then(function (response) {
            if (response.data.Res) {
              self.$message({
                message: '修改成功！',
                type: 'success'
              })
              self.privilegeSelect = []
            } else {
              self.$message({
                message: response.data.Info,
                type: 'error'
              })
            }
            self.getConnTree()
            self.$refs['connForm'].resetFields()
          }).catch(function (error) {
          console.error(error)
        })
      },
      onSaveNew () {
        let self = this
        console.log(this.connForm.PrivilegeRoles)
        this.connForm.PrivilegeRoles = this.privilegeSelect.join(',')
        let param = this.connForm
        axios.post('/connMeta/addConn', param, {responseType: 'json'})
          .then(function (response) {
            if (response.data.Res) {
              self.$message({
                message: '创建成功！',
                type: 'success'
              })
            } else {
              self.$message({
                message: response.data.Info,
                type: 'error'
              })
            }
            self.getConnTree()
          }).catch(function (error) {
          console.error(error)
        })
      },
      onTest () {
        this.testLoading = true
        this.testShow = false
        let self = this
        let param = {
          connType: this.connSelectType,
          user: this.connForm.User,
          passwd: this.connForm.Passwd,
          database: this.connForm.Database,
          host: this.connForm.Host,
          port: this.connForm.Port
        }
        axios.post('/connMeta/testConn', param, {responseType: 'json'})
          .then(function (response) {
            self.testLoading = false
            let isSuccess = response.data.Res
            self.testSuccess = isSuccess
            self.testErrMsg = response.data.Info
            self.testShow = true
          }).catch(function (error) {
          console.error(error)
        })
      },
      triggerNew () {
        // 点击新增按钮
        this.connSelectId = ''
        this.connSelectType = ''
        this.privilegeSelect = []
        this.$refs['connForm'].resetFields()
      },
      onChangeConnType () {
        console.log(this.connSelectType)
        console.log(this.$refs['connForm'])
        this.$refs['connForm'].resetFields()
        this.testShow = false
        this.getConnTree()
      },
      getConnTree () {
        let self = this
        axios.get('/connMeta/getByOwner?connType=' + this.connSelectType)
          .then(function (response) {
            self.treeData = response.data
          })
          .catch(function (error) {
            console.log(error)
          })
      },
      onClickTree (node) {
        this.testShow = false
        this.connSelectId = node.Dsid
        this.connSelectType = node.ConnType
        this.testLoading = false
        this.privilegeSelect = []
        console.log('当前数据源头id: ' + this.connSelectType + '-' + this.connSelectId)
        let self = this
        axios.get('/connMeta/getById/' + node.Dsid)
          .then(function (response) {
            // console.log(response.data)
            self.connForm.DatasourceName = response.data.DatasourceName
            self.connForm.Host = response.data.Host
            self.connForm.User = response.data.User
            self.connForm.Passwd = response.data.Passwd
            self.connForm.Database = response.data.Database
            self.connForm.Port = response.data.Port
            self.connForm.Encode = response.data.Encode
            self.connForm.ConnType = response.data.ConnType
          })
          .catch(function (error) {
            console.log(error)
          })
        axios.get('/connMeta/getPrivilegeRolesById/' + node.Dsid)
          .then(function (response) {
            let res = eval(response.data)
            self.privilegeSelect = res
            console.log(res)

          })
          .catch(function (error) {
            console.log(error)
          })
      },
      checkReadOnly () {
        /*let self = this
        axios.get('/connMeta/getEdit')
          .then(function (response) {
            self.isReadOnly = response.data.readOnly
            console.log('conn readOnly: ' + self.isReadOnly)
          })
          .catch(function (error) {
            console.log(error)
          })*/
      },
      getRoleOptions () {
        let self = this
        axios.get('/auth/roles')
          .then(function (response) {
            self.roleOptions = response.data
          })
          .catch(function (error) {
            console.log(error)
          })
      }
    },
    mounted () {
      this.checkReadOnly()
      this.getConnTree()
      this.getRoleOptions()
    }
  }
</script>

<style scoped>

</style>

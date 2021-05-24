<template>
  <div>
    <h3>增加新角色</h3>
    <el-form inline="true" label-width="120px">
      <el-form-item label="请输入角色名">
        <el-input v-model="newRoleName" size="small"/>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onCreateRole()" size="small">创建</el-button>
      </el-form-item>
    </el-form>

    <h3>角色与工单的增删改查</h3>
    <el-row>
      <el-table
        :data="rolesData"
        stripe
        style="width: 100%">
        <el-table-column
          prop="roleName"
          label="角色名称">
        </el-table-column>
        <el-table-column
          label="工单权限">
          <template slot-scope="role_worksheet">
            <el-button @click="getCheckedWorkSheet(role_worksheet.row.roleName)" type="text" size="mini">查看工单权限
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-row>

    <el-dialog title="配置角色能看到的工单" :visible.sync="worksheetVisible">
      <el-row>
        <div style="max-height: 500px ; overflow:auto">
          <el-table
            :data="worksheetData"
            ref="worksheetTable"
            row-key="TemplateId"
            @selection-change="handleSelectionChange">
            <el-table-column property="worksheetName" label="工单目录"/>
            <el-table-column property="worksheetName" label="修改">
              <template slot-scope="scope">
                <el-switch
                  v-model="scope.row.canModify"
                  :active-value="true"
                  :inactive-value="false"
                  active-color="#02538C"
                  inactive-color="#B9B9B9"
                  @change="changeWorkSheetReadOnly(scope.row)"/>
              </template>
            </el-table-column>
            <el-table-column
              type="selection"
              label="请选择">
            </el-table-column>
          </el-table>
        </div>
      </el-row>

      <el-row>
        <el-col>
          <el-button type="primary" @click="updateRoleWorksheet()" size="small">确定</el-button>
        </el-col>
      </el-row>
    </el-dialog>

    <div style="margin-bottom: 20px">
      <h3>批量添加</h3>
      <label style="font-size: 14px">请选择角色: </label>
      <el-select v-model="batch.roles" multiple filterable placeholder="请选择" size="small" style="width: 80%">
        <el-option
          v-for="role in rolesData"
          :key="role.roleName"
          :label="role.roleName"
          :value="role.roleName">
        </el-option>
      </el-select>

      <br/>
      <label style="font-size: 14px">请选择工单: </label>
      <el-select v-model="batch.workSheetIds" multiple filterable placeholder="请选择" size="small" style="width: 80%">
        <el-option
          v-for="ws in batchWorksheetData"
          :key="ws.worksheetId"
          :label="ws.worksheetName"
          :value="ws.worksheetId">
        </el-option>
      </el-select>
      <el-button @click="onBatchInsert" size="small" plain type="primary">添加</el-button>
    </div>
  </div>
</template>

<script>
  const axios = require('axios')
  import {getCookie, isEmpty} from '../conf/utils'

  export default {
    data () {
      return {
        batch:{
          roles:[],
          workSheetIds:[],
        },
        checkDuplicateRole: '',
        newRoleName: '',
        rolesData: [],
        curWorksheetGroup: '',

        rolesUsersTableVisible: false,
        rolesUsersData: [],
        addRolesUsersName: '',
        batchWorksheetData: [],

        rolesRoutesTableVisible: false,
        rolesRoutesData: [],
        addrolesRoutesName: '',

        worksheetVisible: false,
        worksheetData: [],
        checkedWorksheet: [5],

        tmpRole: '',

        urlVisible: false,
        urlData: [],
        checkedURL: [],
        curURLGroup: '',
        writableURL: [],
      }
    },
    methods: {
      onBatchInsert(){
        console.log(this.batch)
        let self = this
        let params = {
          "roles": this.batch.roles.join(","),
          "worksheetIds": this.batch.workSheetIds.join(","),
        }
        axios.post('/auth/jdbc/batchChangeWorkSheetMapping', params, {responseType: 'json'})
          .then(function (response) {
              self.$message({
                message: '创建成功！',
                type: 'success'
              })
          })
          .catch(function (error) {
            self.$message({
              message: '创建失败！' + error.toString(),
              type: 'error'
            })
          })

      },
      changeWorkSheetReadOnly (row) {
        console.log(row)
      },
      onCreateRole () {
        let roleName = this.newRoleName
        let self = this
        let params = {
          'role': this.newRoleName
        }
        axios.post('/auth/roles/addRole', params, {responseType: 'json'})
          .then(function (response) {
            if (!response.data.Res) {
              self.$message({
                message: '角色名不能为空',
                type: 'error'
              })
            } else {
              self.$message({
                message: '创建成功！',
                type: 'success'
              })
              self.initGetRoles()
            }
          })
          .catch(function (error) {
            self.$message({
              message: '创建失败！' + error.toString(),
              type: 'error'
            })
          })
      },
      toggleSelection (rows) {
        if (rows) {
          rows.forEach(row => {
            this.$refs.worksheetTable.toggleRowSelection(row)
          })
        } else {
          this.$refs.worksheetTable.clearSelection()
        }
      },
      updateRoleWorksheet () {
        let self = this
        let role = this.curWorksheetGroup
        console.log(role)
        let params = this.checkedWorksheet
        // console.log(this.checkedWorksheet)
        axios.post('/auth/jdbc/changeWorkSheetMapping/' + role, params, {responseType: 'json'})
          .then(function (response) {
            self.worksheetVisible = false
            self.$message({
              type: 'success',
              message: '模板权限更新成功!'
            })
            console.info(response.data)
          })
          .catch(function (error) {
            self.$options.methods.saveUpdateNotify.bind(self)()
            console.log(error)
          })
      },
      handleSelectionChange (selectedItems) {
        // 每次点击表格中的选择框, 都会触发这个方法.
        // selectedItems是选中的表格对象
        this.checkedWorksheet = selectedItems
      },
      getCheckedWorkSheet (name) {
        console.log(name)
        let self = this
        self.curWorksheetGroup = name
        axios.get('/auth/template/getAllWorksheetCategory/' + name)
          .then(function (response) {
            let nameRes = []
            let checkedIdx = []
            let index = 0
            response.data.forEach(ele => {
              // 不放HasPrivilege属性, 因为表格中是哦福选中的状态不是通过绑定属性实现的
              // 表格的选中是通过控件的toggleRowSelection()方法
              nameRes.push({
                'worksheetName': ele.worksheetName,
                'worksheetId': ele.worksheetId,
                'canModify': ele.canModify,
              })
              // 找出有权限的worksheet的索引
              if (ele.hasPrivilege) {
                checkedIdx.push(index)
              }
              index++
            })
            self.worksheetData = nameRes
            self.worksheetVisible = true

            let checkedRes = []
            checkedIdx.forEach(ele => {
              checkedRes.push(self.worksheetData[ele])
            })
            self.$nextTick(function () {
              self.toggleSelection(checkedRes)
              // self.toggleSelection([self.worksheetData[1], self.worksheetData[2]]);//每次更新了数据，触发这个函数即可。
            })
          }).catch(function (error) {
          console.log(error)
        })

      },
      initGetRoles () {
        let self = this
        axios.get('/auth/roles')
          .then(function (response) {
            var res = []
            response.data.forEach(ele => {
              // console.log(ele);
              res.push({'roleName': ele})
            })
            self.rolesData = res
          }).catch(function (error) {
          console.log(error)
        })
      },
      getAllWorkSheet (role) {
        let self = this
        axios.get('/auth/template/getAllWorksheetCategory/' + role)
          .then(function (response) {
            let nameRes = []
            response.data.forEach(ele => {
              // 不放HasPrivilege属性, 因为表格中是哦福选中的状态不是通过绑定属性实现的
              // 表格的选中是通过控件的toggleRowSelection()方法
              nameRes.push({
                'worksheetName': ele.worksheetName,
                'worksheetId': ele.worksheetId,
                'canModify': ele.canModify,
              })
            })
            self.batchWorksheetData = nameRes
            console.log(self.batchWorksheetData)
          }).catch(function (error) {
          console.log(error)
        })

      },
    },
    created: function () {
      let self = this
      axios.get('/auth/roles')
        .then(function (response) {
          let res = []
          response.data.forEach(ele => {
            // console.log(ele);
            res.push({'roleName': ele})
          })
          self.rolesData = res
          self.getAllWorkSheet(res[0])
        }).catch(function (error) {
        console.log(error)
      })
    },
    mounted () {
      console.log('mounted')
    }
  }
</script>

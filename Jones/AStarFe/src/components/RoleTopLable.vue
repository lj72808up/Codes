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


    <h3>角色与顶栏标签的增删改查</h3>
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
          label="顶栏标签权限">
          <template slot-scope="role_url">
            <el-button @click="getCheckedURL(role_url.row.roleName)" type="text" size="mini">查看顶栏标签</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-row>

    <el-dialog :title="''+ curURLGroup+' 的顶栏标签配置'" :visible.sync="urlVisible">
      <el-row>
        <div style="max-height: 600px ; overflow:auto">
          <el-table
            :data="urlData"
            ref="urlTable"
            row-key="Lid"
            @selection-change="handleURLSelectionChange">
            <el-table-column property="urlName" label="URL目录"/>
            <el-table-column property="urlIndex" label="URL地址"/>
            <el-table-column property="nameSpace" label="命名空间"/>
            <el-table-column label="写权限">
              <template slot-scope="scope">
                <el-switch
                  v-model="scope.row.writable"
                  :active-value="1"
                  :inactive-value="2"
                  active-color="#02538C"
                  inactive-color="#B9B9B9"
                  @change="changeReadOnly(scope.row)"/>
              </template>
            </el-table-column>

            <el-table-column
              type="selection"
              width="55"
              label="请选择">
            </el-table-column>
          </el-table>
        </div>
      </el-row>

      <el-row>
        <el-col>
          <el-button type="primary" @click="updateRoleURL()" size="small">确定</el-button>
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
      <label style="font-size: 14px">请选择标签: </label>
      <el-select v-model="batch.labels" multiple filterable placeholder="请选择" size="small" style="width: 80%">
        <el-option
          v-for="urlData in batchUrlData"
          :key="urlData.urlId"
          :label="urlData.urlName"
          :value="urlData.urlName">
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
  data() {
    return {
      batch: {
        roles: [],
        labels: []
      },
      checkDuplicateRole: '',
      newRoleName: '',
      rolesData: [],
      curWorksheetGroup: '',

      rolesUsersTableVisible: false,
      rolesUsersData: [],
      addRolesUsersName: '',

      rolesRoutesTableVisible: false,
      rolesRoutesData: [],
      addrolesRoutesName: '',

      worksheetVisible: false,
      worksheetData: [],
      checkedWorksheet: [5],

      tmpRole: '',

      urlVisible: false,
      urlData: [],
      batchUrlData: [],
      checkedURL: [],
      curURLGroup: '',
      writableURL: [],
    }
  },
  methods: {
    onBatchInsert(){
      console.log(this.batch)
      let self = this
      let labelName = []
      let delimiter = " - "
      this.batch.labels.forEach(label=>{
        if (label.search(delimiter) !== -1){
          labelName.push(label.split(delimiter)[1])
        } else {
          labelName.push(label)
        }
      })
      let params = {
        "roles": this.batch.roles.join(","),
        "labelNames": labelName.join(","),
      }
      axios.post('/auth/labels/addNewLabels', params, {responseType: 'json'})
        .then(function (response) {
          console.info(response.data)
          self.$message({
            message: '创建成功！',
            type: 'success'
          })
          self.initGetRoles()
        })
        .catch(function (error) {
          self.$message({
            message: '创建失败！' + error.toString(),
            type: 'error'
          })
        })
    },
    // 改变是否只读的属性
    changeReadOnly(row) {
      // console.log(row)
      if (row.writable === 1) {  // 打开状态
        this.writableURL.push(row.urlId + '')
      } else {  //关闭状态
        // 列表中删除指定元素
        for (let i = 0; i < this.writableURL.length; i++) {
          if (this.writableURL[i] === row.urlId + '') {
            this.writableURL.splice(i, 1)
            break
          }
        }
      }
      console.log(this.writableURL)
    },
    onCreateRole() {
      let roleName = this.newRoleName
      let self = this
      let params = {
        'role': this.newRoleName
      }
      axios.post('/auth/roles/addRole', params, {responseType: 'json'})
        .then(function (response) {
          console.info(response.data)
          self.$message({
            message: '创建成功！',
            type: 'success'
          })
          self.initGetRoles()
        })
        .catch(function (error) {
          self.$message({
            message: '创建失败！' + error.toString(),
            type: 'error'
          })
        })
    },
    getCheckedURL(roleName) {
      console.log(roleName)
      let self = this
      self.curURLGroup = roleName
      axios.get('/auth/labels/' + roleName)
        .then(function (response) {
          let nameRes = []
          let checkedIdx = []
          let index = 0
          self.writableURL = []
          response.data.forEach(ele => {
            nameRes.push({
              'urlName': ele.Title,
              'urlId': ele.Lid,
              'urlIndex': ele.Index,
              'nameSpace': ele.NameSpace,
              'writable': ele.IsWritable
            })
            if (ele.IsWritable === 1) {
              self.writableURL.push(ele.Lid + '')
            }
            if (ele.HasPrivilege) {
              checkedIdx.push(index)
            }
            index++
          })
          self.urlData = nameRes
          self.urlVisible = true

          let checkedRes = []
          checkedIdx.forEach(ele => {
            checkedRes.push(self.urlData[ele])
          })
          self.$nextTick(function () {
            checkedRes.forEach(row => {
              self.$refs.urlTable.toggleRowSelection(row)
            })
            //self.$refs.urlTable.toggleRowSelection(self.urlData[1]) //每次更新了数据，触发这个函数即可。
          })
          console.log('writableURL:' + self.writableURL)
        }).catch(function (error) {
        console.log(error)
      })
    },
    toggleSelection(rows) {
      if (rows) {
        rows.forEach(row => {
          this.$refs.worksheetTable.toggleRowSelection(row)
        })
      } else {
        this.$refs.worksheetTable.clearSelection()
      }
    },
    updateRoleURL() {
      let self = this
      let role = this.curURLGroup
      let checkedIds = this.checkedURL
      console.log(role + ':' + checkedIds)
      let params = {
        'role': [role],
        'ids': checkedIds,
        'writableUrls': this.writableURL
      }
      axios.post('/auth/labels/update', params, {responseType: 'json'})
        .then(function (response) {
          self.urlVisible = false
          self.$message({
            type: 'success',
            message: 'URL权限更新成功!'
          })
          console.info(response.data)
        })
        .catch(function (error) {
          self.$options.methods.saveUpdateNotify.bind(self)()
          console.log(error)
        })
    },
    handleURLSelectionChange(val) {
      let checkedIds = []
      val.forEach(ele => {
        checkedIds.push(ele.urlId + '')
      })
      this.checkedURL = checkedIds
    },
    initGetRoles() {
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
    getLabels(roleName){
      let self = this
      axios.get('/auth/labels/' + roleName)
        .then(function (response) {
          let nameRes = []
          response.data.forEach(ele => {
            nameRes.push({
              'urlName': ele.Title,
              'urlId': ele.Lid,
              'urlIndex': ele.Index,
              'nameSpace': ele.NameSpace,
            })
          })
          self.batchUrlData = nameRes
          console.log(self.batchUrlData)
        }).catch(function (error) {
        console.log(error)
      })
    }
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
        console.log(self.rolesData)
        self.getLabels(self.rolesData[0])
      }).catch(function (error) {
      console.log(error)
    })
  },
  mounted() {
    console.log('mounted')
  }
}
</script>

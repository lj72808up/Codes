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
    <h3>角色路由及用户的增删改查</h3>
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
          label="路由操作">
          <template slot-scope="scope">
            <el-button @click="getRolesRoutes(scope.row.roleName)" type="text" size="mini">查看全部路由</el-button>
          </template>
        </el-table-column>
        <el-table-column
          label="用户操作">
          <template slot-scope="role">
            <el-button @click="getRolesUsers(role.row.roleName)" type="text" size="mini">查看所有用户</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-row>

    <el-dialog title="角色所有路由" :visible.sync="rolesRoutesTableVisible">
      <el-row>
        <el-col>
          <span>删除、添加路由后不能立即更新用户列表！请关闭对话框后再次查看。</span>
        </el-col>
        <el-col :span="16">
          <el-input v-model="addRolesRoutesName" placeholder="请输入添加的路由" size="mini"/>
        </el-col>
        <el-col :span="2" :offset="1">
          <el-button type="success" size="mini" @click="addRolesRoutes()">添加</el-button>
        </el-col>
      </el-row>

      <el-table :data="rolesRoutesData">
        <el-table-column property="routeName" label="路由"/>
        <el-table-column label="操作">
          <template slot-scope="route">
            <el-button @click="delRolesRoutes(route.row.routeName)" type="danger" size="mini">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <el-dialog title="角色所有用户" :visible.sync="rolesUsersTableVisible">
      <el-row>
        <el-col :span="17">
          <span>删除、添加用户后不能立即更新用户列表！请关闭对话框后再次查看。</span>
        </el-col>
        <el-col :span="6" :offset="1">
          <el-checkbox v-model="checkDuplicateRole" size="mini">(去重用户组)</el-checkbox>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="16">
          <el-input v-model="addRolesUsersName" placeholder="请输入添加的用户名" size="mini"/>
        </el-col>
        <el-col :span="2" :offset="1">
          <el-button type="success" size="mini" @click="addRolesUsers()">添加</el-button>
        </el-col>
      </el-row>

      <el-table :data="rolesUsersData">
        <el-table-column property="userName" label="姓名"/>
        <el-table-column label="操作">
          <template slot-scope="user">
            <el-button @click="delRolesUsers(user.row.userName)" type="danger" size="mini">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

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
        <el-table
          :data="worksheetData"
          ref="worksheetTable"
          row-key="TemplateId"
          @selection-change="handleSelectionChange">
          <el-table-column property="worksheetName" label="工单目录"/>
          <el-table-column
            type="selection"
            width="55"
            label="请选择">
          </el-table-column>
        </el-table>
      </el-row>

      <el-row>
        <el-col>
          <el-button type="primary" @click="updateRoleWorksheet()" size="small">确定</el-button>
        </el-col>
      </el-row>
    </el-dialog>


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

    <el-dialog title="配置角色能看到的顶栏标签" :visible.sync="urlVisible">
      <el-row>
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
      </el-row>

      <el-row>
        <el-col>
          <el-button type="primary" @click="updateRoleURL()" size="small">确定</el-button>
        </el-col>
      </el-row>
    </el-dialog>
  </div>
</template>

<script>
  const axios = require('axios')
  import {getCookie, isEmpty} from '../conf/utils'

  export default {
    data () {
      return {
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
        checkedURL: [],
        curURLGroup: '',
        writableURL: [],
      }
    },
    methods: {
      // 改变是否只读的属性
      changeReadOnly (row) {
        // console.log(row)
        if(row.writable===1){  // 打开状态
          this.writableURL.push(row.urlId+'')
        }else{  //关闭状态
          // 列表中删除指定元素
          for (let i = 0; i < this.writableURL.length; i++) {
            if (this.writableURL[i] === row.urlId+'') {
              this.writableURL.splice(i, 1);
              break;
            }
          }
        }
        console.log(this.writableURL)
      },
      onCreateRole () {
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
      getCheckedURL (roleName) {
        console.log(roleName)
        let self = this
        self.curURLGroup = roleName
        axios.get('/auth/labels/' + roleName)
          .then(function (response) {
            let nameRes = []
            let checkedIdx = []
            let index = 0
            response.data.forEach(ele => {
              console.log(ele)
              nameRes.push({
                'urlName': ele.Title,
                'urlId': ele.Lid,
                'urlIndex': ele.Index,
                'nameSpace': ele.NameSpace,
                'writable': ele.IsWritable
              })
              console.log(nameRes)
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
          }).catch(function (error) {
          console.log(error)
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
      updateRoleURL () {
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
      updateRoleWorksheet () {
        let self = this
        let role = this.curWorksheetGroup
        let checkedIds = this.checkedWorksheet
        console.log(role)
        let params = {
          'role': [role],
          'ids': checkedIds,
        }
        axios.post('/auth/jdbc/changeWorkSheetMapping', params, {responseType: 'json'})
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
      handleSelectionChange (val) {
        let checkedIds = []
        val.forEach(ele => {
          checkedIds.push(ele.worksheetId + '')
        })
        this.checkedWorksheet = checkedIds
        console.log(this.checkedWorksheet)
      },
      handleURLSelectionChange (val) {
        let checkedIds = []
        val.forEach(ele => {
          checkedIds.push(ele.urlId + '')
        })
        this.checkedURL = checkedIds
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
              console.log(ele)
              nameRes.push({
                'worksheetName': ele.Title,
                'worksheetId': ele.TemplateId
              })
              if (ele.HasPrivilege) {
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
            console.log('checkedRes:---')
            console.log(checkedRes)
            self.$nextTick(function () {
              self.toggleSelection(checkedRes)
              // self.toggleSelection([self.worksheetData[1], self.worksheetData[2]]);//每次更新了数据，触发这个函数即可。
            })
          }).catch(function (error) {
          console.log(error)
        })

      },
      getRolesRoutes (name) {
        let self = this
        const comURL = '/auth/routes/role/' + name
        axios.get(comURL)
          .then(function (response) {
            var res = []
            response.data.forEach(ele => {
              res.push({'routeName': ele})
            })
            self.rolesRoutesData = res
            console.log(res)
            self.rolesRoutesTableVisible = true
            self.tmpRole = name
          }).catch(function (error) {
          console.log(error)
        })
      },
      addRolesRoutes () {
        let self = this
        if (!isEmpty(this.addRolesRoutesName)) {
          const comURL = '/auth/roles/route/' + this.tmpRole
          var params = new URLSearchParams()
          params.append('rolename', this.tmpRole)
          params.append('route', this.addRolesRoutesName)
          axios({
            method: 'post',
            url: comURL,
            data: params
          }).then(function () {
            self.addRolesRoutesName = ''
            self.rolesRoutesTableVisible = false
            self.$message('添加成功！')
          }).catch(function (error) {
            console.log(error)
          })
        }
      },
      delRolesRoutes (name) {
        let self = this
        const comURL = '/auth/roles/route/' + this.tmpRole + '?route=' + name
        axios({
          method: 'delete',
          url: comURL
        }).then(function () {
          self.rolesRoutesTableVisible = false
          self.$message('删除成功！')
        }).catch(function (error) {
          console.log(error)
        })
      },
      getRolesUsers (name) {
        let self = this
        const comURL = '/auth/users/role/' + name
        axios.get(comURL)
          .then(function (response) {
            var res = []
            response.data.forEach(ele => {
              res.push({'userName': ele})
            })
            self.rolesUsersData = res
            self.rolesUsersTableVisible = true
            self.tmpRole = name
          }).catch(function (error) {
          console.log(error)
        })
      },
      addRolesUsers () {
        let self = this
        if (!isEmpty(this.addRolesUsersName)) {
          const comURL = '/auth/roles/user/' + this.tmpRole
          var params = new URLSearchParams()
          params.append('rolename', this.tmpRole)
          params.append('username', this.addRolesUsersName)
          params.append('checkDuplicate',this.checkDuplicateRole)
          axios({
            method: 'post',
            url: comURL,
            data: params
          }).then(function (response) {
            if (response.data.Res) {
              self.$message({
                message: response.data.Info,
                type: 'success'
              })
              self.addRolesUsersName = ''
              self.rolesUsersTableVisible = false
            }else{
              self.$message({
                message: "该用户已存在角色 : " + response.data.Info,
                type: 'error'
              })
            }
          }).catch(function (error) {
            console.log(error)
          })
        }
      },
      delRolesUsers (name) {
        let self = this
        const comURL = '/auth/roles/user/' + this.tmpRole + '?username=' + name
        axios({
          method: 'delete',
          url: comURL
        }).then(function () {
          self.rolesUsersTableVisible = false
          self.$message('删除成功！')
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
      }
    },
    created: function () {
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
    mounted () {
      console.log('mounted')
    }
  }
</script>

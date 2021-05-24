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
<!--        <el-table-column-->
<!--          label="路由操作">-->
<!--          <template slot-scope="scope">-->
<!--            <el-button @click="getRolesRoutes(scope.row.roleName)" type="text" size="mini">查看全部路由</el-button>-->
<!--          </template>-->
<!--        </el-table-column>-->
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

      <div style="max-height: 600px ; overflow:auto">
        <el-table :data="rolesUsersData">
          <el-table-column property="userName" label="姓名"/>
          <el-table-column label="操作">
            <template slot-scope="user">
              <el-button @click="delRolesUsers(user.row.userName)" type="danger" size="mini">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
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
          params.append('checkDuplicate', this.checkDuplicateRole)
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
            } else {
              self.$message({
                message: '该用户已存在角色 : ' + response.data.Info,
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

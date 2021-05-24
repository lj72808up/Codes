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

    <h3>角色与数据源的增删改查</h3>
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
          label="数据源权限">
          <template slot-scope="role_metatable">
            <el-button @click="getCheckedMetaTable(role_metatable.row.roleName,'read')" type="text" size="mini">查询权限
            </el-button>
          </template>
        </el-table-column>
<!--        <el-table-column
          label="数据源权限">
          <template slot-scope="role_metatable">
            <el-button @click="getCheckedMetaTable(role_metatable.row.roleName,'exec')" type="text" size="mini">执行权限
            </el-button>
          </template>
        </el-table-column>   -->
      </el-table>
    </el-row>


    <el-dialog title="表级查询权限配置" :visible.sync="metatableReadVisible">
      <el-row>
        <div style="max-height: 500px ; overflow:auto">
        <el-tree
          :data="metatableData"
          ref="ReadTree"
          show-checkbox
          node-key="Id"
          :default-expanded-keys="expandTreeId"
          :default-checked-keys="curSelectTableIdList"
          :props="treeDefaultProps">
        </el-tree>
        </div>
      </el-row>

      <el-row>
        <el-col>
          <el-button type="primary" @click="updateRoleMetatable('read')" size="small">确定</el-button>
        </el-col>
      </el-row>
    </el-dialog>

    <el-dialog title="表级执行权限配置" :visible.sync="metatableExecVisible">
      <el-row>
        <div style="max-height: 500px ; overflow:auto">
        <el-tree
          :data="metatableData"
          ref="ExecTree"
          show-checkbox
          node-key="Id"
          :default-expanded-keys="expandTreeId"
          :default-checked-keys="curSelectTableIdList"
          :props="treeDefaultProps">
        </el-tree>
        </div>
      </el-row>

      <el-row>
        <el-col>
          <el-button type="primary" @click="updateRoleMetatable('exec')" size="small">确定</el-button>
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
        newRoleName: '',
        rolesData: [],
        curMetatableGroup: '',

        metatableReadVisible: false,
        metatableExecVisible: false,
        metatableData: [],
        curSelectTableIdList: [],

        expandTreeId : Array.from({ length: 1000 }, (_, i) => 1000001 + i),
        //操作权限类型  读取 read ， 执行 exec
        type : 'read',

        treeDefaultProps: {
          children: 'NodeChildren',
          label: 'NodeName'
        },
      }
    },
    methods: {
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

      updateRoleMetatable (type) {
        let self = this

        if(type == "read"){
          this.curSelectTableIdList = this.$refs.ReadTree.getCheckedKeys()
        }
        else{
          this.curSelectTableIdList = this.$refs.ExecTree.getCheckedKeys()
        }

        let role = this.curMetatableGroup
        let checkedIds = []
        for(var  idx in this.curSelectTableIdList) {
          if(this.curSelectTableIdList[idx] >= 1000001) {
            continue
          } else {
            checkedIds.push(this.curSelectTableIdList[idx]+'')
          }
        }
        console.log(checkedIds)

        let params = {
          'role': [role],
          'ids': checkedIds,
          'type': [type],
        }
        axios.post('/auth/metadata/changeMetaTableMapping', params, {responseType: 'json'})
          .then(function (response) {
            type == "read" ? self.metatableReadVisible = false : self.metatableExecVisible = false

            self.$message({
              type: 'success',
              message: '数据源权限更新成功!'
            })
            console.info(response.data)
          })
          .catch(function (error) {
            self.$options.methods.saveUpdateNotify.bind(self)()
            console.log(error)
          })
      },

      //从后台获取当前role的 metaTableId 读取标志位，执行标志位
      getCheckedMetaTable:async function (name,curType){
        let self = this
        self.type = curType
        self.curSelectTableIdList = []
        self.curMetatableGroup = name

        await axios.get('/auth/metadata/GetMetaTableAccess/' + name)
          .then( (response) => {
            response.data.forEach(ele => {
              if(curType == 'read' && ele.ReadAccess == 1 ){
                self.curSelectTableIdList.push(ele.Tid)
              }
              if(curType == 'exec' && ele.ExecAccess == 1 ){
                self.curSelectTableIdList.push(ele.Tid)
              }
            })
          }).catch( (error) => {
          console.log(error)
        })

        await axios.get('/auth/metadata/GetAllMetaTableTreeAccess/' + name)
          .then( (response) => {
            self.metatableData = response.data
            curType == 'read' ? self.metatableReadVisible = true : self.metatableExecVisible = true

          }).catch( (error) => {
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

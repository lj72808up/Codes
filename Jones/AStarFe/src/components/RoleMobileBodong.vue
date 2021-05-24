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

    <h3>波动图标权限配置</h3>
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
          <template slot-scope="scope">
            <el-button @click="getCheckedBodongMenu(scope.row.roleName,'read')" type="text" size="mini">查询权限
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-row>


    <el-dialog title="表级查询权限配置" :visible.sync="bodongTreeVisible">
      <el-row>
        <div style="height: 500px; max-height: 500px ; overflow:auto">
          <el-tree
            :data="demoData"
            ref="bodongTree"
            show-checkbox
            node-key="menuId"
            default-expand-all
            @close="closeTree"
            :check-strictly="true"
            :default-expanded-keys="expandTreeId"
            :default-checked-keys="curCheckedIds"
            :props="treeDefaultProps"/>
        </div>
        <!--@check-change="handleFirstLevelCheck"-->
      </el-row>

      <el-row type="flex" justify="end">
        <el-button type="primary" @click="updateBodongMapping()" size="small">确定</el-button>
      </el-row>
    </el-dialog>
  </div>
</template>

<script>
  const axios = require('axios')
  import {getCookie, isEmpty} from '../conf/utils'

  export default {
    data() {
      return {
        newRoleName: '',
        rolesData: [],
        curRole: '',

        bodongTreeVisible: false,
        curCheckedIds: [],

        expandTreeId: Array.from({length: 1000}, (_, i) => 1000001 + i),
        //操作权限类型  读取 read ， 执行 exec
        type: 'read',

        /*treeDefaultProps: {
          children: 'NodeChildren',
          label: 'NodeName'
        },*/
        treeDefaultProps: {
          children: 'children',
          label: 'label'
        },
        demoData: [],
      }
    },
    methods: {
      handleFirstLevelCheck(targetNode, isCheckedSelf, isCheckedInChildren) {
        // 处理第一层的节点选中状态
        if (targetNode.pid === 0) {
          let childKeys = []
          if (isCheckedSelf) {  // 首层节点选中
            if (targetNode.children !== undefined){
              targetNode.children.forEach(child => {
                childKeys.push(child.menuId)
              })
            }
            let checkedSet = new Set(childKeys)
            checkedSet.add(targetNode.menuId)
            this.curCheckedIds.forEach(id => {
              checkedSet.add(id)
            })
            this.$refs.bodongTree.setCheckedKeys(Array.from(checkedSet))
          } else { //首层节点反选

          }
        }
        console.log(targetNode)
        console.log(isCheckedSelf)
        console.log(isCheckedInChildren)
      },
      closeTree() {
        this.curCheckedIds = []
      },
      onCreateRole() {
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

      updateBodongMapping() {
        let checkedNodes = this.$refs.bodongTree.getCheckedNodes()
        let role = this.curRole
        let checkedIdSet = new Set()
        checkedNodes.forEach(node => {
          checkedIdSet.add(node.menuId)
          if (node.pid !== 0) {
            checkedIdSet.add(node.pid)
          }
        })
        this.curCheckedIds = Array.from(checkedIdSet)
        console.log(checkedIdSet)

        let self = this
        let params = {
          'role': role,
          'ids': this.curCheckedIds.join(","),
        }
        console.log(params)
        axios.post('/auth/bodong/UpdateCheckedBodongNode', params, {responseType: 'json'})
          .then(function (response) {
            if (response.data.Res) {
              self.$message({
                type: 'success',
                message: '波动图表权限更新成功!'
              })
              self.bodongTreeVisible = false
            } else {
              self.$message({
                type: 'error',
                message: response.data.Info
              })
            }
          })
          .catch(function (error) {
            console.log(error)
          })
      },

      //从后台获取当前role的 metaTableId 读取标志位，执行标志位
      getCheckedBodongMenu(role) {
        let self = this
        this.curRole = role
        this.bodongTreeVisible = true
        axios.get('/auth/bodong/GetAllBodongTree')
          .then((response) => {
            self.demoData = response.data
          }).catch((error) => {
          console.log(error)
        })
        // 选中保存的关系
        this.setCheckedKeys()
      },
      setCheckedKeys() {
        let self = this
        axios.get('/auth/bodong/GetCheckedBodongNode?role=' + this.curRole)
          .then(function (response) {
            self.$refs.bodongTree.setCheckedKeys(response.data)
          }).catch(function (error) {
          console.log(error)
        })
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
    mounted() {
      console.log('mounted')
    }
  }
</script>

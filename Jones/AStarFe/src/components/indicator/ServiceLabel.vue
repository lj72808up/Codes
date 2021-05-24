<template>
  <div>
    <el-row>
      <h3 style="margin-bottom: 0px; margin-top: 5px; display:inline;">业务标签</h3>
      <el-tooltip content="将指标按照业务逻辑进行分类管理，分类会展示在指标配置中" placement="top">
        <i class="el-icon-warning-outline"/>
      </el-tooltip>
    </el-row>
    <el-button type="primary" size="medium" style="margin-bottom: 10px" @click="newTagDialog=true" v-show="!topLabelReadOnly">创建业务标签</el-button>

    <el-table :data="dataList" stripe style="width: 100%" :border="false"
              :header-cell-style="{background:'#eee',color:'#606266'}">
      <el-table-column key="tag_name" prop="tag_name" label="业务标签名"/>
      <el-table-column key="description" prop="description" label="描述"/>
      <el-table-column key="create_time" prop="create_time" label="创建时间" :sortable="true"/>
      <el-table-column key="modify_time" prop="modify_time" label="最后修改时间" :sortable="true"/>
      <el-table-column key="operate" prop="operate" label="操作" width="100px" align="center">
        <template slot-scope="scope">
          <el-button icon="el-icon-edit" type="text" @click="onClickEditServiceTag(scope.row)"/>
        </template>
      </el-table-column>
    </el-table>


    <el-dialog title="新建业务标签" :visible.sync="newTagDialog" width="600px" center="true" :show-close="false"
               :before-close="closeTagDialog">
      <el-form ref="tagForm" :model="tagForm" label-width="80px"
               :rules="tagFormRules">
        <el-form-item label="标签名称" prop="name">
          <el-input v-model="tagForm.name" size="mini"/>
        </el-form-item>
        <el-form-item label="描述(选填)" prop="desc">
          <el-input v-model="tagForm.desc" size="mini" type="textarea" :rows="10"/>
        </el-form-item>
      </el-form>

      <el-row type="flex" justify="end">
        <el-button size="small" type="primary" @click="makeNewServiceTag">创建</el-button>
        <el-button size="small" @click="closeTagDialog">取消</el-button>
      </el-row>

    </el-dialog>


    <el-dialog title="编辑业务标签" :visible.sync="editTagDialog" width="700px" center="true" :show-close="false"
               :before-close="closeTagDialog">
      <el-form ref="tagForm" :model="tagForm" label-width="80px"
               :rules="tagFormRules">
        <el-form-item label="标签名称" prop="name">
          <el-input v-model="tagForm.name" size="mini"/>
        </el-form-item>
        <el-form-item label="描述(选填)" prop="desc">
          <el-input v-model="tagForm.desc" size="mini" type="textarea" :rows="10"/>
        </el-form-item>
      </el-form>

      <el-row type="flex" justify="end">
        <el-button size="small" type="primary" @click="makeEditServiceTag" v-show="!topLabelReadOnly">提交</el-button>
        <el-button size="small" @click="closeTagDialog">取消</el-button>
      </el-row>

    </el-dialog>
  </div>

</template>

<script>
  const axios = require('axios')
  export default {
    data () {
      return {
        dataList: [],
        tableColumn: [],
        newTagDialog: false,
        editTagDialog: false,
        tagForm: {
          name: '',
          desc: ''
        },
        topLabelReadOnly: false,
        editTagId: 0,
        tagFormRules: {
          'name': [{required: true, message: '请输入标签名', trigger: 'blur'}]
        },
      }
    },
    methods: {
      getEdit(){
        let self = this
        axios.get('/indicator/getEdit')
          .then(function (response) {
            self.topLabelReadOnly = response.data.readOnly
            console.log('顶栏标签配置的serviceTag只读权限:'+self.topLabelReadOnly)
          })
          .catch(function (error) {
            console.log(error)
          })
      },
      onClickEditServiceTag (row) {
        this.editTagDialog = true
        console.log(row)
        this.$nextTick(() => {
          this.tagForm.name = row.tag_name
          this.tagForm.desc = row.description
          this.editTagId = row.tid
        })
      },
      closeTagDialog () {
        console.log(this.$refs['tagForm'])
        this.$refs['tagForm'].resetFields()
        this.$nextTick(() => {
          // alert(JSON.stringify(this.tagForm))
          this.newTagDialog = false
          this.editTagDialog = false
        })
      },
      getColumnWidth (field) {
        if (field === 'description') {
          return '60%'
        }
        if (field === 'tid') {
          return '20%'
        }
      },
      getAllServiceTag () {
        let self = this
        axios.get('/indicator/serviceTag/getAll')
          .then(function (response) {
            self.tableColumn = response.data.cols
            self.dataList = response.data.data
          }).catch(function (error) {
          console.log(error)
        })
      },
      makeEditServiceTag () {
        let self = this
        this.$refs['tagForm'].validate((valid) => {
          if (valid) {
            // 表单校验通过
            let params = {
              'tagName': this.tagForm.name,
              'description': this.tagForm.desc,
              'tid': this.editTagId
            }
            axios.post('/indicator/serviceTag/editServiceTag', params, {responseType: 'json'})
              .then(function (response) {
                if (response.data.Res) {
                  self.$message({
                    type: 'success',
                    message: '编辑成功'
                  })
                  self.getAllServiceTag()
                } else {
                  self.$message({
                    type: 'error',
                    message: response.data.Info
                  })
                }
                self.editTagDialog = false
              }).catch(function (error) {
              console.log(error)
            })
          }
        })
      },
      makeNewServiceTag () {
        let self = this
        this.$refs['tagForm'].validate((valid) => {
          if (valid) {
            // 表单校验通过
            let params = {
              'name': this.tagForm.name,
              'desc': this.tagForm.desc
            }
            axios.post('/indicator/serviceTag/putServiceTag', params, {responseType: 'json'})
              .then(function (response) {
                if (response.data.Res) {
                  self.$message({
                    type: 'success',
                    message: '创建成功'
                  })
                  self.getAllServiceTag()
                } else {
                  self.$message({
                    type: 'error',
                    message: response.data.Info
                  })
                }
                self.newTagDialog = false
              }).catch(function (error) {
              console.log(error)
            })
          }
        })
      }
    },
    mounted: function () {
      this.getAllServiceTag()
      this.getEdit()
    }
  }
</script>

<style scoped>

</style>

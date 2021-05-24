<template>
  <div>
    <div>
      <el-row>
        <h3 style="margin-bottom: 0px; margin-top: 5px; display:inline;">原子指标</h3>
      </el-row>
      <el-button type="primary" size="medium" style="margin-bottom: 10px"
                 @click="onClickNew">新建原子指标</el-button>

      <el-table :data="dataList" stripe style="width: 100%" :border="false"
                :header-cell-style="{background:'#eee',color:'#606266'}">
        <el-table-column key="name" prop="name" label="指标名" min-width="12%"/>
        <el-table-column key="meaning" prop="meaning" label="指标含义" min-width="52%"/>
        <el-table-column key="aggregation" prop="aggregation" label="聚合方式" min-width="12%"/>
        <el-table-column key="modify_time" prop="modify_time" label="最后修改时间" :sortable="true" min-width="12%"/>
        <el-table-column key="business_label" prop="business_label" label="业务标签" min-width="12%">
          <template slot-scope="scope">
            <el-tag style="margin-right: 5px" v-for="(businessLabel, index) in scope.row.business_label.split(',')">
              {{serviceTagMap[businessLabel]}}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column key="operate" prop="operate" label="操作" width="100px" align="center">
          <template slot-scope="scope">
            <el-button icon="el-icon-edit" type="text" @click="onClickEditIndicator(scope.row)"/>
            <el-button icon="el-icon-delete" type="text" @click="onDeleteIndicator(scope.row)"/>
          </template>
        </el-table-column>
      </el-table>
    </div>

  </div>

</template>

<script>
  import router from '../../router'

  const axios = require('axios')
  export default {
    data () {
      return {
        editId: 0,
        dataList: [],
        tableColumn: [],
        indicatorForm: {
          name: '',
          aggregation: '',
          meaning: '',
          serviceTagIds: [],
        },
        indicatorFormRules: {
          'name': [{required: true, message: '请输入标签名', trigger: 'blur'}],
          'serviceTagIds': [{required: true, message: '请选择业务类型', trigger: ['blur', 'change']}]
        },
        serviceTagList: [],
        serviceTagMap: {},
        factRelationList: [],
        facTableForm: {
          tableType: '',
          factTable: '',
          foreignKey: ''
        },
        dsTypeOptions: [{'name': 'mysql'}, {'name': 'hive'}, {'name': 'clickhouse'}],
        factTableOptions: [],
        foreignKeyOptions: [],
        topLabelReadOnly: false,
        fromDictionary: false,
      }
    },
    methods: {
      onClickNew() {
        sessionStorage.setItem('preIndicatorId', '');
        router.push({path: 'preDefineIndicatorCreate'})
      },
      postFactRelation () {
        let self = this
        let param = {
          'indicatorId': this.editId,
          'factId': this.facTableForm.factTable,
          'associatedField': this.facTableForm.foreignKey,
        }
        axios.post('/indicator/predefine/addIndicatorFactRelation', param)
          .then(function (response) {
            if (response.data.Res) {
              self.$notify({
                message: '创建成功',
                type: 'success'
              })
              self.getFactRelation()
            } else {
              self.$notify({
                message: response.data.Info,
                type: 'error'
              })
            }
          })
          .catch(function (error) {
            console.log(error)
          })
      },
      onDeleteFactRelation (row) {
        let self = this
        let ifid = row.ifid
        axios.post('/indicator/predefine/delIndicatorFactRelation/' + ifid)
          .then(function (response) {
            if (response.data.Res) {
              self.$notify({
                message: '删除成功',
                type: 'success'
              })
              self.getFactRelation()
            } else {
              self.$notify({
                message: response.data.Info,
                type: 'error'
              })
            }
          })
          .catch(function (error) {
            console.log(error)
          })
      },
      genTableIdentifier (database, tableName) {
        return database + '.' + tableName
      },
      onChangeFactTable () {
        let self = this
        let params = {
          'tableId': self.facTableForm.factTable + '',
          'version': 'latest',
        }
        axios.post('/metaData/fields/getByVersion', params, {responseType: 'json'})
          .then(function (response) {
            self.foreignKeyOptions = response.data.dataList
          })
          .catch(function (error) {
            console.log(error)
          })
      },
      onChangeFactType () {
        let tableType = this.facTableForm.tableType
        this.facTableForm.factTable = ''
        this.foreignKeyOptions = []
        let self = this
        let param = {
          'curPage': '1',
          'pageSize': '9999',
          'table_type': tableType
        }
        axios.post('metaData/getTableByConditions', param, {responseType: 'json'})
          .then(function (response) {
            self.factTableOptions = response.data.data
          }).catch(function (error) {
          console.log(error)
        })
      },
      getFactRelation () {
        let self = this
        axios.get('indicator/getIndicatorFactRelation/' + this.editId)
          .then(function (response) {
            self.factRelationList = response.data.data
          })
          .catch(function (error) {
            console.log(error)
          })
      },
      onDeleteIndicator (row) {
        this.$confirm('确定删除吗', {
          confirmButtonText: '删除',
          cancelButtonText: '取消'
        }).then(() => {
          this.editId = row.id
          let self = this
          axios.post('/indicator/delIndicator/' + this.editId)
            .then(function (response) {
              if (response.data.Res) {
                self.$message({
                  type: 'success',
                  message: '删除成功'
                })
                self.getAllPredefineIndicator()
              } else {
                self.$message({
                  type: 'error',
                  message: response.data.Info
                })
              }
            }).catch(function (error) {
            console.log(error)
          })
        })
      },
      onClickEditIndicator (row) {
        sessionStorage.setItem('preIndicatorId', row.id);
        router.push({path: 'preDefineIndicatorCreate'})
      },
      closeIndicatorDialog () {
        this.$refs['indicatorForm'].resetFields()
      },
      browseTableInfo (tid) {
        let info = {
          'tableId': tid
        }
        sessionStorage.setItem('metaData_table_id', JSON.stringify(info))
        // 跳转
        router.push({path: 'tableInfo'})
      },
      getColumnWidth (field) {
        if (field === 'meaning') {
          return '50%'
        }
        if (field === 'tid') {
          return '20%'
        }
      },
      getAllPredefineIndicator () {
        let self = this
        axios.get('/indicator/getAllIndicator')
          .then(function (response) {
            self.tableColumn = response.data.cols
            self.dataList = response.data.data
          }).catch(function (error) {
          console.log(error)
        })
      },
      makeEditPredefineIndicator () {
        // todo 更新接口
        let self = this
        let params = {
          'name': this.indicatorForm.name,
          'meaning': this.indicatorForm.meaning,
          'BusinessLabels': this.indicatorForm.serviceTagIds,
          'aggregation': this.indicatorForm.aggregation,
          'id': this.editId,
        }
        axios.post('/indicator/editPreIndicator', params, {responseType: 'json'})
          .then(function (response) {
            if (response.data.Res) {
              self.$message({
                type: 'success',
                message: '修改成功'
              })
              self.getAllPredefineIndicator()
            } else {
              self.$message({
                type: 'error',
                message: response.data.Info
              })
            }
          }).catch(function (error) {
          console.log(error)
        })
      },
      makeNewPredefineIndicator () {
        if (this.editId > 0) {
          this.makeEditPredefineIndicator()
        } else {
          let self = this
          this.$refs['indicatorForm'].validate((valid) => {
            if (valid) {
              // 表单校验通过
              let params = {
                'name': this.indicatorForm.name,
                'meaning': this.indicatorForm.meaning,
                'BusinessLabels': this.indicatorForm.serviceTagIds,
                'aggregation': this.indicatorForm.aggregation
              }
              axios.post('/indicator/postPreIndicator', params, {responseType: 'json'})
                .then(function (response) {
                  if (response.data.Res) {
                    self.$message({
                      type: 'success',
                      message: '创建成功'
                    })
                    self.getAllPredefineIndicator()
                  } else {
                    self.$message({
                      type: 'error',
                      message: response.data.Info
                    })
                  }
                }).catch(function (error) {
                console.log(error)
              })
            }
          })
        }

      },
      getAllServiceTag (func) {
        let self = this
        axios.get('/indicator/serviceTag/getAll')
          .then(function (response) {
            let dataList = response.data.data
            let serviceTags = []
            dataList.forEach(ele => {
              serviceTags.push({
                'tagName': ele.tag_name,
                'tid': ele.tid + ''
              })
              self.serviceTagMap[ele.tid] = ele.tag_name
            })
            self.serviceTagList = serviceTags
            console.log(self.serviceTagMap)
            if ((func !== null) && (func !== undefined) && (func !== "")) {
              func() // 执行回调函数
            }
          }).catch(function (error) {
          console.log(error)
        })
      },
      getEdit () {
        let self = this
        axios.get('/indicator/getEdit')
          .then(function (response) {
            self.topLabelReadOnly = response.data.readOnly
            console.log('顶栏标签配置的indicator只读权限:' + self.topLabelReadOnly)
          })
          .catch(function (error) {
            console.log(error)
          })
      },
    },
    /*created () {
      let originId = sessionStorage.getItem('originId')
      console.log(originId)
      if ((originId !== null) && (originId !== undefined) && (originId !== "")) {
        sessionStorage.removeItem('originId')
        this.topLabelReadOnly = true   // 跳转过来是查看的
        this.getAllServiceTag()
        let self = this
        let params = {
          "originId": originId+'',
          "assetType": "predefine"
        }
        this.fromDictionary = true
        axios.post('/indicator/dictionary/searchById', params, {responseType: 'json'})
          .then(function (response) {
            let row = {   // 接口返回的是json格式, 转换成TableData的形式
              "id":response.data.id,
              "name":response.data.name,
              "business_label": response.data.businessLabel,
              "aggregation":response.data.aggregation,
              "meaning":response.data.meaning,
            }
            self.onClickEditIndicator (row)
          }).catch(function (error) {
          console.log(error)
        })
      } else {
        this.getAllServiceTag(this.getAllPredefineIndicator)
      }
    },*/
    mounted: function () {
      this.getAllServiceTag(this.getAllPredefineIndicator)
    }
  }
</script>

<style scoped>

</style>

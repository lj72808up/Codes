<template>
  <div>
    <div>
      <h3 style="margin-bottom: 0px; margin-top: 5px;">维度管理</h3>
      <el-select v-model="isDerivative" placeholder="请选择" size="small" @change="getAllDimensions">
        <el-option label="主维度" :value="false"/>
        <el-option label="衍生维度" :value="true"/>
      </el-select>
      <el-popconfirm
        confirm-button-text='衍生维度'
        cancel-button-text='主维度'
        cancel-button-type='success'
        confirm-button-type='info'
        @confirm="onClickCreate(true)"
        @cancel="onClickCreate(false)"
        title="请选择要创建的维度类型?">
        <!--<el-button type="primary" size="medium" style="margin-bottom: 10px"
                   @click="onCreateDimension">创建维度
        </el-button>-->
        <el-button type="primary" size="medium" style="margin-bottom: 10px" slot="reference">创建维度</el-button>
      </el-popconfirm>
      <!--主维度列表-->
      <el-table :data="dataList" stripe style="width: 100%" :border="false"
                :header-cell-style="{background:'#eee',color:'#606266'}">
        <el-table-column key="name" prop="name" label="中文名称"/>
        <el-table-column key="alias_name" prop="alias_name" label="英文名称"/>
        <el-table-column v-if="isDerivative===false" key="pk_field" prop="pk_field" label="维度主键"/>
        <el-table-column v-if="isDerivative===false" key="primary_table_data_source" prop="primary_table_data_source" label="主表数据源"
                         :sortable="true"/>
        <el-table-column v-if="isDerivative===false" key="primary_table_name" prop="primary_table_name" label="维度主表">
          <template slot-scope="scope">
          <span @click="browseTableInfo(scope.row.primary_table_id)"
                style="cursor:pointer;text-decoration: underline">
            {{scope.row.primary_database_name}}.{{scope.row.primary_table_name}}
          </span>
          </template>
        </el-table-column>
        <!--        <el-table-column key="primary_database_name" prop="primary_database_name" label="主表数据库"/>-->
        <el-table-column key="business_label" prop="business_label" label="业务标签">
          <template slot-scope="scope">
            <el-tag style="margin-right: 5px" v-for="(businessLabel, index) in scope.row.business_label.split(',')"
                    v-show="scope.row.business_label.split(',').length>0 && scope.row.business_label.split(',')[0]!==''">
              {{serviceTagMap[businessLabel]}}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column key="operate" prop="operate" label="操作" width="100" align="center">
          <template slot-scope="scope">
            <el-button icon="el-icon-edit" type="text" @click="onEditDimension(scope.row)"/>
          </template>
        </el-table-column>
      </el-table>

      <!--&lt;!&ndash;衍生维度&ndash;&gt;
      <el-table :data="dataList" stripe style="width: 100%" :border="false"
                :header-cell-style="{background:'#eee',color:'#606266'}" v-if="isDerivative===true">
        <el-table-column key="name" prop="name" label="中文名称"/>
        <el-table-column key="alias_name" prop="alias_name" label="英文名称"/>
      </el-table>-->
    </div>

  </div>
</template>

<script>
  import router from '../../router'

  const axios = require('axios')
  export default {
    data() {
      return {
        isDerivative: false,
        name: 'Dimension',
        dataList: [],
        createDimension: false,
        serviceTagList: [],
        serviceTagMap: {},
        dimensionForm: {
          id: 0,
          name: '',
          dataSourceType: '',
          dataSource: '',
          primaryTableId: '',
          primaryTableName: '',
          pkField: '',
          serviceTagIds: [],
          businessLabel: '',
        },
        dsTypeOptions: [{'name': 'mysql'}, {'name': 'hive'}, {'name': 'clickhouse'}],
        dsOptions: [],
        tableOptions: [],
        pkOptions: [],
        rules: {
          name: [{required: true, message: '请输入维度名', trigger: 'blur'}],
          dataSource: [{required: true, message: '请选择数据源', trigger: 'blur'}],
          primaryTableId: [{required: true, message: '请学则维度主表', trigger: 'blur'}],
          pkField: [{required: true, message: '请选择维度主键', trigger: 'blur'}],
        },
        facTableForm: {
          tableType: '',
          factTable: '',
          foreignKey: ''
        },
        factTableOptions: [],
        foreignKeyOptions: [],
        factRelationList: [],
        topLabelReadOnly: false,
        fromDictionary: false,
      }
    },
    methods: {
      onClickCreate(isDerivative) {
        if (isDerivative){  // 跳转衍生维度
          sessionStorage.setItem('dimensionId', '')
          router.push({path: 'dimensionCreateDerivative'})
        } else {
          sessionStorage.setItem('dimensionId', '')
          router.push({path: 'dimensionCreate'})
        }
      },
      onDeleteFactRelation(row) {
        let dfId = row.dfid
        let self = this
        let params = {
          'dfid': dfId
        }
        this.$confirm('确定删除这个事实表关系吗', {
          confirmButtonText: '删除',
          cancelButtonText: '取消'
        }).then(() => {
          axios.post('/indicator/dimension/postDimFactRelation', params, {responseType: 'json'})
            .then(function (response) {
              if (response.data.Res) {
                self.$notify({
                  message: '删除成功',
                  type: 'success'
                })
                self.$refs['facTableForm'].resetFields()
                self.getFactRelationByDimId()
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
        })
      }
      ,
      getFactRelationByDimId() {
        let url = '/indicator/dimension/getDimFactRelations/' + this.dimensionForm.id
        let self = this
        axios.get(url)
          .then(function (response) {
            self.$nextTick(() => {
              self.factRelationList = response.data.data
            })
          })
          .catch(function (error) {
            console.log(error)
          })
      }
      ,
      postFactRelation() {
        let self = this
        let params = {
          'dimId': this.dimensionForm.id,
          'factId': this.facTableForm.factTable,
          'foreignKey': this.facTableForm.foreignKey
        }
        axios.post('/indicator/dimension/postDimFactRelation', params, {responseType: 'json'})
          .then(function (response) {
            if (response.data.Res) {
              self.$notify({
                message: '更新成功',
                type: 'success'
              })
              self.$refs['facTableForm'].resetFields()
              self.getFactRelationByDimId()
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
      }
      ,
      onChangeFactTable() {
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
      }
      ,
      genTableIdentifier(database, tableName) {
        return database + '.' + tableName
      }
      ,
      onChangeFactType() {
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

      }
      ,
      getAllDimensions() {
        let self = this
        axios.get('/indicator/dimension/getAll?isDerivative=' + this.isDerivative)
          .then(function (response) {
            self.dataList = response.data.data
          }).catch(function (error) {
          console.log(error)
        })
      }
      ,
      onEditDimension(row) {
        if (row.is_derivative){ // 衍生维度
          sessionStorage.setItem('dimensionId', row.id)
          router.push({path: 'dimensionCreateDerivative'})
        } else {   // 主维度
          sessionStorage.setItem('dimensionId', row.id)
          router.push({path: 'dimensionCreate'})
        }
      }
      ,
     /* onCreateDimension() {
        // sessionStorage 存入
        sessionStorage.setItem('dimensionId', '')
        router.push({path: 'dimensionCreate'})
      }
      ,*/
      closeCreate() {
        this.$refs['dimensionForm'].resetFields()
        this.$refs['facTableForm'].resetFields()
        this.dimensionForm.id = 0  // 没出现在表单中的属性, 手动置0
        // alert (JSON.stringify(this.indicatorForm))
        this.$nextTick(() => {
          this.createDimension = false
        })

      }
      ,
      browseTableInfo(tid) {
        let info = {
          'tableId': tid
        }
        sessionStorage.setItem('metaData_table_id', JSON.stringify(info))
        // 跳转
        router.push({path: 'tableInfo'})
      }
      ,
      onChangeDSType() {
        let self = this
        axios.get('/dw/GetConnByPrivilege/' + this.dimensionForm.dataSourceType)
          .then(function (response) {
            self.dsOptions = response.data
            // 如下两行,默认选择数据源的第一个选项,并更新树
            self.dimensionForm.dataSource = self.dsOptions[0].Dsid
            self.onChangeDataSource()
          })
          .catch(function (error) {
            console.log(error)
          })
      }
      ,
      onChangeDataSource() {
        let self = this
        console.log(this.dimensionForm.dataSource)
        let curdsType = this.dimensionForm.dataSourceType
        let curdsId = this.dimensionForm.dataSource
        let params = {
          'dsType': curdsType,
          'dsId': curdsId + '',
        }

        axios.post('/dw/meta/get', params, {responseType: 'json'})
          .then(function (response) {
            let treeData = response.data.children
            let tbList = []
            treeData.forEach(ele => {
              try {
                ele.children.forEach(cld => {
                  tbList.push({
                    'tableId': cld.tid,
                    'tableName': cld.tableName,
                  })
                })
              } catch (err) {
                console.log(err)
              }
            })
            self.tableOptions = tbList
          })
          .catch(function (error) {
            console.log(error)
          })
      }
      ,
      onChangeTableName() {
        let self = this
        this.findTableNameInOptions()
        let params = {
          'tableId': self.dimensionForm.primaryTableId + '',
          'version': 'latest',
        }
        axios.post('/metaData/fields/getByVersion', params, {responseType: 'json'})
          .then(function (response) {
            self.pkOptions = response.data.dataList
            self.foreignKeyOptions = response.data.dataList
          })
          .catch(function (error) {
            console.log(error)
          })
      }
      ,
      findTableNameInOptions() {
        let tableId = this.dimensionForm.primaryTableId
        let options = this.tableOptions
        options.forEach(ele => {
          if (ele.tableId == tableId) {
            this.dimensionForm.primaryTableName = ele.tableName
          }
        })
      }
      ,
      getAllServiceTag(func) {
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
            func()  // 执行回调函数
          }).catch(function (error) {
          console.log(error)
        })
      }
    }
    ,
    mounted: function () {
      this.getAllServiceTag(this.getAllDimensions)
    },
  }
</script>

<style scoped>
</style>

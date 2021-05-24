<template>
  <div>
    <!--维度创建-->
    <div>
      <h3 style="margin-bottom: 0px; margin-top: 5px;">维度主表配置</h3>
      <el-row style="margin-top: 10px; width: 100%">
        <el-form ref="dimensionForm" :model="dimensionForm" label-width="110px" :rules="rules">
          <el-row>
            <el-form-item label="中文名称:" prop="name" style="display: inline-block; width: 30%">
              <el-input v-model="dimensionForm.name" size="small"/>
            </el-form-item>
            <el-form-item label="英文名称:" prop="aliasName" style="display: inline-block; width: 30%">
              <el-input v-model="dimensionForm.aliasName" size="small"/>
            </el-form-item>
          </el-row>
          <!--          <span>维度主表配置</span>-->
          <el-form-item label="主表数据源:" prop="dataSource" style="width: 60%">
            <el-row>
              <el-col :span="11">
                <el-select v-model="dimensionForm.dataSourceType" filterable placeholder="请选择数据源类型:"
                           style="width: 100%" @change="onChangeDSType" size="small">
                  <el-option
                    v-for="ele in dsTypeOptions"
                    :key="ele.name"
                    :label="ele.name"
                    :value="ele.name">
                  </el-option>
                </el-select>
              </el-col>
              <el-col style="text-align: center;" :span="1">-</el-col>
              <el-col :span="12">
                <el-select v-model="dimensionForm.dataSource" filterable placeholder="请选择数据源:"
                           style="width: 100%" @change="onChangeDataSource" size="small">
                  <el-option
                    v-for="ele in dsOptions"
                    :key="ele.Dsid"
                    :label="ele.DatasourceName"
                    :value="ele.Dsid">
                  </el-option>
                </el-select>
              </el-col>
            </el-row>
          </el-form-item>

          <el-form-item label="维度表表名:" prop="primaryTableId" style="display: inline-block; width: 30%">
            <el-select v-model="dimensionForm.primaryTableId" filterable clearable placeholder="请选择表名:"
                       @change="onChangeTableName" style="width: 100%" size="small">
              <el-option
                v-for="ele in tableOptions"
                :key="ele.tableId"
                :label="ele.tableName"
                :value="ele.tableId">
              </el-option>
            </el-select>
          </el-form-item>


          <el-form-item label="主键(ID):" prop="pkField" style="display: inline-block;width: 30%">
            <el-select v-model="dimensionForm.pkField" filterable placeholder="维度主键:"
                       style="width: 100%" size="small">
              <el-option
                v-for="ele in pkOptions"
                :key="ele.fieldName"
                :label="ele.fieldName"
                :value="ele.fieldName">
              </el-option>
            </el-select>
          </el-form-item>

          <!--<el-form-item label="值键(名称):" prop="valField" style="display: inline-block;width: 30%">
            <el-select v-model="dimensionForm.valField" filterable placeholder="维度值键"
                       style="width: 100%" size="small">
              <el-option
                v-for="ele in pkOptions"
                :key="ele.fieldName"
                :label="ele.fieldName"
                :value="ele.fieldName">
              </el-option>
            </el-select>
          </el-form-item>-->

          <!--<el-row>
            <el-form-item label="属性字段:" prop="assistInfo" style="display: inline-block;width: 90%">
              <el-select v-model="dimensionForm.assistInfoArr" filterable placeholder="辅助信息:"
                         style="width: 100%" size="small" multiple="true">
                <el-option
                  v-for="ele in pkOptions"
                  :key="ele.fieldName"
                  :label="ele.fieldName"
                  :value="ele.fieldName">
                </el-option>
              </el-select>
            </el-form-item>
          </el-row>
-->
          <el-form-item label="维度描述:" prop="description" style="width: 90%">
            <el-input type="textarea" :autosize="{ minRows: 2, maxRows: 3}" style="width: 100%"
                      v-model="dimensionForm.description"/>
          </el-form-item>

          <el-form-item label="业务标签:" prop="serviceTagIds" style="width: 90%">
            <el-select v-model="dimensionForm.serviceTagIds" multiple filterable placeholder="请选择"
                       style="width: 100%" size="mini">
              <el-option
                v-for="ele in serviceTagList"
                :key="ele.tid"
                :label="ele.tagName"
                :value="ele.tid">
              </el-option>
            </el-select>
          </el-form-item>

        </el-form>
      </el-row>
      <el-row type="flex" justify="end" style="width: 90%">
        <el-button size="small" type="primary" @click="onCreateDimension" v-show="!isCheck">提交</el-button>
        <el-button size="small" @click="$router.back(-1)">返回</el-button>
      </el-row>

      <h3 style="margin-bottom: 10px; margin-top: 5px; display: inline-block">关联事实表配置</h3>
      <!--关联事实表配置-->
      <el-tooltip class="item" effect="dark" content="事实表外键与维度表主键相同则符合规范" placement="top-start">
        <el-button icon="el-icon-question" type="text" style="color: #E6A23C; font-size: 16px"/>
      </el-tooltip>
      <el-form ref="facTableForm" :model="facTableForm" :inline="true" label-width="100px">
        <el-form-item label="事实表:">
          <el-select v-model="facTableForm.tableType" filterable placeholder="请选择事实表类型"
                     size="small" @change="onChangeFactType">
            <el-option
              v-for="ele in dsTypeOptions"
              :key="ele.name"
              :label="ele.name"
              :value="ele.name">
            </el-option>
          </el-select>
          <i class="el-icon-right"/>
          <el-select v-model="facTableForm.factTable" placeholder="请选择事实表" filterable
                     size="small" @change="onChangeFactTable" clearable>
            <el-option
              v-for="ele in factTableOptions"
              :key="ele.tid"
              :label="genTableIdentifier(ele.database_name,ele.table_name)"
              :value="ele.tid">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="关联外键:">
          <el-select v-model="facTableForm.foreignKey" filterable placeholder="请选择事实表的关联外键"
                     size="small" clearable>
            <el-option
              v-for="ele in foreignKeyOptions"
              :key="ele.fieldName"
              :label="ele.fieldName"
              :value="ele.fieldName">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button size="small" type="success" @click="postFactRelation" v-show="!isCheck">提交</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="factRelationList" stripe style="width: 85%; margin-left:30px" :border="true">
        <el-table-column key="fact_database" prop="fact_database" label="事实表库"/>
        <el-table-column key="fac_table_name" prop="fac_table_name" label="事实表名称">
          <template slot-scope="scope">
          <span @click="browseTableInfo(scope.row.fact_id)"
                style="cursor:pointer;text-decoration: underline">{{scope.row.fac_table_name}}</span>
          </template>
        </el-table-column>
        <el-table-column key="foreign_key" prop="foreign_key" label="关联外键">
          <template slot-scope="scope">
            <span>{{scope.row.foreign_key}}</span>
          </template>
        </el-table-column>
        <el-table-column key="foreign_key" prop="foreign_key_valid" label="规范符合" width="200" align="center">
          <template slot-scope="scope">
            <el-button v-if="scope.row.foreign_key===dimensionForm.pkField" icon="el-icon-s-flag"
                       style="color: green; font-size: 22px" type="text"/>
            <el-button v-else icon="el-icon-warning" style="color: red; font-size: 22px" type="text"/>
          </template>
        </el-table-column>
        <el-table-column key="operate" prop="operate" label="操作" width="100px" align="center">
          <template slot-scope="scope">
            <el-button icon="el-icon-delete" type="text" @click="onDeleteFactRelation(scope.row)"/>
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
    data() {
      return {
        name: 'Dimension',
        isCheck: false,
        dataList: [],
        serviceTagList: [],
        serviceTagMap: {},
        dimensionForm: {
          id: 0,
          name: '',
          aliasName: '',
          dataSourceType: '',
          dataSource: '',
          primaryTableId: '',
          primaryTableName: '',
          description: '',
          pkField: '',
          valField: '',
          assistInfo: '',
          assistInfoArr: [],
          serviceTagIds: [],
          businessLabel: '',
        },
        dsTypeOptions: [{'name': 'mysql'}, {'name': 'hive'}, {'name': 'clickhouse'}],
        dsOptions: [],
        tableOptions: [],
        pkOptions: [],
        foreignKeyOptions: [],
        rules: {
          name: [{required: true, message: '请输入维度名', trigger: 'blur'}],
          aliasName: [{required: true, message: '请输入英文别名', trigger: 'blur'}],
          dataSource: [{required: true, message: '请选择数据源', trigger: 'blur'}],
          primaryTableId: [{required: true, message: '请学则维度主表', trigger: 'blur'}],
          pkField: [{required: true, message: '请选择维度主键', trigger: 'blur'}],
          valField: [{required: true, message: '请选择维度值字段', trigger: 'blur'}]
        },
        facTableForm: {
          tableType: '',
          factTable: '',
          foreignKey: ''
        },
        factTableOptions: [],
        // pkOptions: [],
        factRelationList: [],
        fromDictionary: false,
      }
    },
    methods: {
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
        console.log("id为" + this.dimensionForm.id)
        if (this.dimensionForm.id === "" || this.dimensionForm.id === undefined || this.dimensionForm.id === null || this.dimensionForm.id == 0) {
          this.$alert("请先提交维度主表配置")
          return
        }
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
            let fieldNameArr = self.foreignKeyOptions.map(ele => ele.fieldName)
            if (fieldNameArr.indexOf(self.dimensionForm.pkField) > -1) {
              self.facTableForm.foreignKey = self.dimensionForm.pkField
            } else {
              self.facTableForm.foreignKey = ''
            }
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
      },
      onEditDimension(row) {
        this.$nextTick(() => {
          console.log(row)
          this.dimensionForm.id = row.id
          this.dimensionForm.name = row.name
          this.dimensionForm.aliasName = row.aliasName
          this.dimensionForm.dataSourceType = row.data_source_type
          this.dimensionForm.valField = row.val_field
          this.getFactRelationByDimId()
          this.onChangeDSType()
          this.dimensionForm.dataSource = row.primary_table_data_source
          this.onChangeDataSource()
          this.dimensionForm.primaryTableId = row.primary_table_id
          this.dimensionForm.primaryTableName = row.primary_table_name
          this.onChangeTableName()
          this.dimensionForm.pkField = row.pk_field
          this.dimensionForm.description = row.description
          if (row.assist_info !== '') {
            this.dimensionForm.assistInfoArr = row.assist_info.split(',')
          }
          if (row.business_label !== '') {
            this.dimensionForm.serviceTagIds = row.business_label.split(',')
          }
        })
      }
      ,
      onCreateDimension() {
        console.log(this.dimensionForm)
        let self = this
        this.$refs['dimensionForm'].validate((valid) => {
          if (valid) {
            this.dimensionForm.businessLabel = this.dimensionForm.serviceTagIds.join(',')
            this.dimensionForm.assistInfo = this.dimensionForm.assistInfoArr.join(',')
            axios.post('/indicator/dimension/postNew', this.dimensionForm, {responseType: 'json'})
              .then(function (response) {
                if (response.data.Res) {
                  self.$message({
                    type: 'success',
                    message: '创建成功',
                  })
                  // todo 返回dimensionId
                  self.dimensionForm.id = parseInt(response.data.Info)
                } else {
                  self.$message({
                    type: 'error',
                    message: response.data.Info,
                  })
                }
              })
              .catch(function (error) {
                console.log(error)
              })
          }
        })
      },
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
            self.pkOptions = response.data.dataList
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
      ,
      getDimensionById(dimensionId) {
        let self = this
        let params = {
          'originId': dimensionId + '',
          'assetType': 'dimension'
        }
        this.fromDictionary = true
        axios.post('/indicator/dictionary/searchById', params, {responseType: 'json'})
          .then(function (response) {
            let row = {   // 接口返回的是json格式, 转换成TableData的形式
              'id': response.data.Id,
              'name': response.data.Name,
              'aliasName': response.data.AliasName,
              'data_source_type': response.data.DataSourceType,
              'primary_table_data_source': response.data.PrimaryTableDataSource,
              'primary_table_id': response.data.PrimaryTableId,
              'primary_table_name': response.data.PrimaryTableName,
              'pk_field': response.data.PkField,
              'val_field': response.data.ValField,
              'assist_info': response.data.AssistInfo,
              'business_label': response.data.BusinessLabel,
              'description': response.data.Description,
            }
            self.onEditDimension(row)
          }).catch(function (error) {
          console.log(error)
        })
      }
    }
    ,
    mounted() {
      let dimensionId = sessionStorage.getItem('dimensionId')
      this.isCheck = (sessionStorage.getItem("dimensionCreate") === 'true')
      sessionStorage.removeItem("dimensionCreate")
      console.log('收到的dimensionId: ' + dimensionId)
      this.getAllServiceTag()
      // 判断是新建还是修改
      if ((dimensionId !== null) && (dimensionId !== undefined) && (dimensionId !== '')) {
        this.getDimensionById(dimensionId)
      }
    },
  }
</script>

<style scoped>
</style>

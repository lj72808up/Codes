<template>
  <div>
    <!--新建指标的dialog-->
    <div style="width: 80%">
      <h3 style="margin-bottom: 0px; margin-top: 5px;">原子指标</h3>
      <el-form ref="indicatorForm" :model="indicatorForm" label-width="80px"
               :rules="indicatorFormRules">
        <el-row>
          <el-col :span="11">
            <el-form-item label="中文名称" prop="name">
              <el-input v-model="indicatorForm.name" size="small"/>
            </el-form-item>
          </el-col>
          <el-col :span="11" :offset="1">
            <el-form-item label="英文名称" prop="fieldName">
              <el-input v-model="indicatorForm.fieldName" size="small"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="11">
            <el-form-item label="聚合方式" prop="aggregate">
              <el-select size="small" v-model="indicatorForm.aggregation" filterable placeholder="请选择" style="width: 100%">
                <el-option label="sum" value="sum"/>
                <el-option label="min" value="min"/>
                <el-option label="max" value="max"/>
                <el-option label="avg" value="avg"/>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="11" :offset="1">
            <el-form-item label="业务标签" prop="serviceTagIds">
              <el-select v-model="indicatorForm.serviceTagIds" multiple filterable placeholder="请选择"
                         style="width: 100%" size="small">
                <el-option
                  v-for="ele in serviceTagList"
                  :key="ele.tid"
                  :label="ele.tagName"
                  :value="ele.tid">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <!-- 如果有 -->
        <el-form-item label="计算公式">
          <el-button class="button-new-expression" size="mini" type="primary" plain @click="expressionVisible=true">+
            计算公式
          </el-button>
          <span style="margin-left: 10px; color: #b3b3b3">{{omitExpression}}</span>
          <!--              <el-input v-model="indicatorForm.expression" type="textarea" :min-rows="1" :max-rows="3"/>-->
        </el-form-item>

        <!--标签选择框-->
        <el-form-item label="指标含义" prop="meaning">
          <el-input v-model="indicatorForm.meaning" size="small" type="textarea" :rows="3"/>
        </el-form-item>

      </el-form>
      <el-row type="flex" justify="end">
        <el-button size="small" type="primary" @click="submitIndicator" v-show="!isCheck">提交</el-button>
        <el-button size="small" @click="$router.back(-1)">返回</el-button>
      </el-row>

      <h3 style="display: inline-block">关联表配置</h3>
      <el-tooltip class="item" effect="dark" content="事实表的'关联字段名称'和原子指标的'指标字段名'相同则符合规范" placement="top-start">
        <el-button icon="el-icon-question" type="text" style="color: #E6A23C; font-size: 16px"/>
      </el-tooltip>
      <!--关联事实表配置-->
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
        <el-form-item label="关联值字段:">
          <el-select v-model="facTableForm.foreignKey" filterable placeholder="请选择事实表的关联值字段"
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
          <el-button size="small" type="primary" @click="postFactRelation" v-show="!isCheck">提交</el-button>
        </el-form-item>
      </el-form>
      <el-table :data="factRelationList" stripe style="width: 100%; margin-left:30px" :border="true">
        <el-table-column key="fact_name" prop="fact_name" label="事实表名称">
          <template slot-scope="scope">
          <span @click="browseTableInfo(scope.row.fact_id)"
                style="cursor:pointer;text-decoration: underline">{{scope.row.fact_name}}</span>
          </template>
        </el-table-column>
        <el-table-column key="database_name" prop="database_name" label="事实表库名"/>
        <el-table-column key="associated_field" prop="associated_field" label="关联值字段"/>
        <el-table-column key="fact_type" prop="fact_type" width="150" align="center" label="事实表类型"/>
        <el-table-column key="associated_field_valid" prop="associated_field_valid" label="规范符合" width="150"
                         align="center">
          <template slot-scope="scope">
            <el-button v-if="scope.row.associated_field===indicatorForm.fieldName" icon="el-icon-s-flag"
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

    <!--todo -->
    <el-dialog title="配置计算公式" :visible.sync="expressionVisible" width="60%">
      <div style="margin-right: 10px">
        <div style="border:1px solid #DCDFE6; height: 400px;">
          <el-scrollbar style="height: 100%;">
            <div id="添加的指标">
              <el-col :span="12">
                <el-select v-model="indicatorForm.predictItems" multiple filterable placeholder="请选择需要组合的原子指标"
                           style="margin-left: 5px; width: 80%" size="mini">
                  <el-option
                    v-for="ele in preIndicatorOptions"
                    :key="ele.id"
                    :label="ele.name"
                    :value="ele.id">
                  </el-option>
                </el-select>
                <el-button size="mini" style="margin-left: 2px" type="primary" plain @click="addDependencies">添加
                </el-button>
                <el-row v-for="(indicate, idx) in indexDependencies">
                  <el-col :span="22">
                    <el-row
                      style="margin-left: 15px; border:1px solid #DCDFE6; margin-bottom: 5px; background-color:#FAFAFA">
                      <el-col :span="4">
                        <div class="indexIcon">{{idx+1}}</div>
                      </el-col>
                      <el-col :span="20">
                        <span>{{indicate.name}}</span>
                      </el-col>
                    </el-row>
                  </el-col>
                </el-row>
                <span class="expressionPlaceHolder" v-show="indexDependencies.length===0">添加指标到这里将进行组合</span>
              </el-col>

              <el-col :span="12" style="margin-top: 17px;">
                <el-input id="expressionTextarea" type="textarea" v-model="indicatorForm.expression"
                          :rows="17" placeholder="请输入计算表达式" size="medium" style="width: 97%"/>
              </el-col>
            </div>
          </el-scrollbar>
        </div>
      </div>
    </el-dialog>
  </div>

</template>

<script>
  import router from '../../router'

  const axios = require('axios')
  export default {
    data() {
      return {
        editId: 0,
        dataList: [],
        isCheck: false,
        tableColumn: [],
        indicatorForm: {
          name: '',
          fieldName: '',
          aggregation: '',
          meaning: '',
          serviceTagIds: [],
          predictItems: [],
          expression: '',
        },
        expressionVisible: false,
        preIndicatorOptions: [],
        indexDependencies: [],
        indicatorFormRules: {
          'name': [{required: true, message: '请输入指标名', trigger: 'blur'}],
          'fieldName': [{required: true, message: '请输入指标字段', trigger: 'blur'}],
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
        fromDictionary: false,
      }
    },
    computed: {
      omitExpression() {
        if (this.indicatorForm.expression.length > 40) {
          return this.indicatorForm.expression.substr(0,40) + '...'
        }else {
          return this.indicatorForm.expression
        }
      }
    },
    methods: {
      postFactRelation() {
        if (this.arbitrateNewOrUpdate() === 'new') {
          this.$alert("请先保存上方的基本信息")
          return
        }
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
      deleteDependencies(idx) {
        this.indexDependencies.splice(idx, 1)
      },
      addDependencies() {
        let selectIds = this.indicatorForm.predictItems
        this.transStrToIndexDependencies(selectIds)

        this.indexDependencies.forEach(ele => {
          this.indicatorForm.expression = this.indicatorForm.expression + '\n' + ele.name
        })
      },
      transStrToIndexDependencies(selectIds) {
        // 从preIndicatorOptions中找出selectIds的集合, 放入indexDependencies
        this.indexDependencies = []
        this.preIndicatorOptions.forEach(indicator => {
          if (selectIds.indexOf(indicator.id) !== -1) {
            this.indexDependencies.push(indicator)
          }
        })
        this.indicatorForm.predictItems = selectIds
      },
      onDeleteFactRelation(row) {
        this.$confirm('确定删除吗', {
          confirmButtonText: '删除',
          cancelButtonText: '取消'
        }).then(() => {
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
        })
      },
      genTableIdentifier(database, tableName) {
        return database + '.' + tableName
      },
      onChangeFactTable() {
        let self = this
        let params = {
          'tableId': self.facTableForm.factTable + '',
          'version': 'latest',
        }
        axios.post('/metaData/fields/getByVersion', params, {responseType: 'json'})
          .then(function (response) {
            //todo 筛选维度主键
            self.foreignKeyOptions = response.data.dataList
            let fieldNameArr = self.foreignKeyOptions.map(ele => ele.fieldName)
            if (fieldNameArr.indexOf(self.indicatorForm.fieldName) > -1) {
              self.facTableForm.foreignKey = self.indicatorForm.fieldName
            } else {
              self.facTableForm.foreignKey = ''
            }
          })
          .catch(function (error) {
            console.log(error)
          })
      },
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
      getFactRelation() {
        let self = this
        axios.get('indicator/getIndicatorFactRelation/' + this.editId)
          .then(function (response) {
            self.factRelationList = response.data.data
          })
          .catch(function (error) {
            console.log(error)
          })
      },
      onDeleteIndicator(row) {
        this.editId = row.id
        let self = this
        axios.post('/indicator/delIndicator/' + this.editId)
          .then(function (response) {
            if (response.data.Res) {
              self.$message({
                type: 'success',
                message: '删除成功'
              })
              self.getAllIndicator()
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
      getAllPreDefineIndex() {
        let self = this
        axios.get('/indicator/getAllIndicator')
          .then(function (response) {
            self.preIndicatorOptions = response.data.data
          }).catch(function (error) {
          console.log(error)
        })
      },
      closeIndicatorDialog() {
        router.push({path: 'preIndicator'})
      },
      browseTableInfo(tid) {
        let info = {
          'tableId': tid
        }
        sessionStorage.setItem('metaData_table_id', JSON.stringify(info))
        // 跳转
        router.push({path: 'tableInfo'})
      },
      getColumnWidth(field) {
        if (field === 'meaning') {
          return '50%'
        }
        if (field === 'tid') {
          return '20%'
        }
      },
      makeEditPredefineIndicator() {
        let self = this
        let dependencyIndexIds = []
        this.indexDependencies.forEach(indicator => {
          dependencyIndexIds.push(indicator.id)
        })
        let params = {
          'name': this.indicatorForm.name,
          'FieldName': this.indicatorForm.fieldName,
          'meaning': this.indicatorForm.meaning,
          'BusinessLabels': this.indicatorForm.serviceTagIds,
          'aggregation': this.indicatorForm.aggregation,
          'id': this.editId,
          'expression': this.indicatorForm.expression,  // 表达式
          'dependencies': dependencyIndexIds,
        }
        axios.post('/indicator/editPreIndicator', params, {responseType: 'json'})
          .then(function (response) {
            if (response.data.Res) {
              self.$message({
                type: 'success',
                message: '修改成功'
              })
              // self.getAllIndicator()
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
      submitIndicator() {
        this.checkDuplicateIndicator(this.makeNewPredefineIndicator)
      },
      checkDuplicateIndicator(fun) {
        let params = {
          'name': this.indicatorForm.name,
          'fieldName': this.indicatorForm.fieldName,
          'id': this.editId + ""
        }
        let self = this
        axios.post('/indicator/dictionary/checkIndicatorDuplicate', params, {responseType: 'json'})
          .then(function (response) {
            let indicators = response.data
            if (indicators.length > 0) {
              // 存在同名的
              self.$message({
                type: 'error',
                message: '指标重复: ' + indicators[0].name + ' 指标字段: ' + indicators[0].fieldName
              })
            } else {
              // 不存在同名的
              fun()
            }
          })
          .catch(function (error) {
            console.log(error)
          })
      },
      makeNewPredefineIndicator() {
        this.$refs['indicatorForm'].validate((valid) => {
          if (valid) {
            if (this.arbitrateNewOrUpdate() === 'update') {
              this.makeEditPredefineIndicator()
            } else {
              let self = this
              this.$refs['indicatorForm'].validate((valid) => {
                if (valid) {
                  let dependencyIndexIds = []
                  this.indexDependencies.forEach(indicator => {
                    dependencyIndexIds.push(indicator.id)
                  })
                  // 表单校验通过
                  let params = {
                    'name': this.indicatorForm.name,
                    'meaning': this.indicatorForm.meaning,
                    'BusinessLabels': this.indicatorForm.serviceTagIds,
                    'aggregation': this.indicatorForm.aggregation,
                    'fieldName': this.indicatorForm.fieldName,
                    'expression': this.indicatorForm.expression,  // 表达式
                    'dependencies': dependencyIndexIds,
                  }
                  axios.post('/indicator/postPreIndicator', params, {responseType: 'json'})
                    .then(function (response) {
                      if (response.data.Res) {
                        self.$message({
                          type: 'success',
                          message: '创建成功, id:' + response.data.Info
                        })
                        self.editId = parseInt(response.data.Info)
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
          }
        })
      },
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
            if ((func !== null) && (func !== undefined) && (func !== "")) {
              func() // 执行回调函数
            }
          }).catch(function (error) {
          console.log(error)
        })
      },
      arbitrateNewOrUpdate() {
        let editId = this.editId
        if ((editId === undefined) || (editId === null) || (editId == '') || (editId == '0')) {
          return 'new'
        } else {
          return 'update'
        }
      },
      getDependencies() {
        let self = this
        axios.get('/indicator/getIndicatorDependenciesById/' + this.editId)
          .then(function (response) {
            self.transStrToIndexDependencies(response.data)
          }).catch(function (error) {
          console.log(error)
        })
      },
      getPreIndicatorById(preIndicatorId) {
        let self = this
        let params = {
          'originId': preIndicatorId + '',
          'assetType': 'predefine'
        }
        this.fromDictionary = true
        axios.post('/indicator/dictionary/searchById', params, {responseType: 'json'})
          .then(function (response) {
            self.editId = parseInt(response.data.id)
            self.getDependencies()
            self.indicatorForm.name = response.data.name
            self.indicatorForm.meaning = response.data.meaning
            self.indicatorForm.aggregation = response.data.aggregation
            self.indicatorForm.fieldName = response.data.fieldName
            self.indicatorForm.expression = response.data.expression
            self.getFactRelation()
            if (response.data.businessLabel !== '') {
              self.indicatorForm.serviceTagIds = response.data.businessLabel.split(',')
            }
          }).catch(function (error) {
          console.log(error)
        })
      },
    },
    mounted: function () {
      let preIndicatorId = sessionStorage.getItem("preIndicatorId")
      this.isCheck = (sessionStorage.getItem("preIndicatorCheck") === 'true')
      sessionStorage.removeItem("preIndicatorCheck")
      this.getAllServiceTag()
      if ((preIndicatorId !== null) && (preIndicatorId !== undefined) && (preIndicatorId !== '')) {
        this.getPreIndicatorById(preIndicatorId)
      }
      this.getAllPreDefineIndex()
    }
  }
</script>

<style scoped>
  .indexIcon {
    height: 24px;
    width: 24px;
    margin-left: 15px;
    background-color: #f2f4fc;
    color: #222f73;
    border-radius: 50%;
    line-height: 24px;
    text-align: center;
    margin-top: 8px;
  }

  .expressionPlaceHolder {
    color: #ccc;
    font-family: PingFang SC;
    font-size: 12px;
    font-weight: 500;
    margin-left: calc(50% - 80px);
    display: block;
    margin-top: 56px;
  }

  .button-new-expression {
    height: 32px;
    line-height: 30px;
    padding-top: 0;
    padding-bottom: 0;
  }
</style>

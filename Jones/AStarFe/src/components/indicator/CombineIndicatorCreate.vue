<template>
  <div>

    <!--新建指标的dialog-->
    <div>
      <el-row>
        <h3 style="margin-bottom: 0px; margin-top: 5px; display:inline;">派生指标</h3>
      </el-row>
      <el-row>
        <el-form ref="indicatorForm" :model="indicatorForm" label-width="100px"
                 :rules="indicatorFormRules">
          <el-row>
            <el-col :span="8">
              <el-form-item label="中文名称:" prop="name">
                <el-input v-model="indicatorForm.name" size="mini"/>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item label="英文名称:" prop="fieldName">
                <el-input v-model="indicatorForm.fieldName" size="mini"/>
              </el-form-item>
            </el-col>
          </el-row>
          <!--标签选择框-->
          <el-form-item label="时间修饰词">
            <el-select v-model="indicatorForm.timeQualifier" filterable placeholder="请选择" size="mini">
              <el-option
                v-for="ele in timeQualifierOption"
                :label="ele.CName"
                :key="ele.EName"
                :value="ele.Qid">
              </el-option>
            </el-select>

            <el-button type="primary" plain size="small" @click="qualifierVisible = true">+ 新增时间修饰词</el-button>
          </el-form-item>

          <el-form-item label="其它修饰词">
            <el-tag
              :key="tag"
              v-for="tag in indicatorForm.qualifiers"
              closable
              :disable-transitions="false"
              @close="handleClose(tag)">
              {{tag}}
            </el-tag>
            <el-input
              class="input-new-tag"
              v-if="inputVisible"
              v-model="inputValue"
              ref="saveTagInput"
              size="small"
              @keyup.enter.native="handleInputConfirm"
              @blur="handleInputConfirm"/>
            <el-button v-else class="button-new-tag" size="small" @click="showInput">+ 新修饰词</el-button>

          </el-form-item>

          <el-form-item label="业务标签" prop="serviceTagIds">
            <el-select v-model="indicatorForm.serviceTagIds" multiple filterable placeholder="请选择"
                       style="width: 40%" size="mini">
              <el-option
                v-for="ele in serviceTagList"
                :key="ele.tid"
                :label="ele.tagName"
                :value="ele.tid">
              </el-option>
            </el-select>
          </el-form-item>

          <el-row>
            <el-col :span="19">
              <el-form-item label="指标含义" prop="meaning">
                <el-input v-model="indicatorForm.meaning" size="small" type="textarea" :rows="4"/>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>

      </el-row>
      <el-row type="flex" justify="end" style="width: 80%">
        <el-button size="small" type="primary" @click="makeNewCombineIndicator" v-show="!isCheck">提交
        </el-button>
        <el-button size="small" @click="$router.back(-1)">返回</el-button>
      </el-row>

      <div>

      </div>
    </div>


    <el-dialog title="添加时间修饰词" :visible.sync="qualifierVisible" width="400px">
      <el-form ref="timeQualifierForm" :model="timeQualifierForm" label-width="100px"
               :rules="timeQualifierFormRules">
        <el-row>
          <el-col :span="22">
            <el-form-item label="中文名称" prop="cName">
              <el-input v-model="timeQualifierForm.cName" size="mini"/>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="22">
            <el-form-item label="英文名称" prop="eName">
              <el-input v-model="timeQualifierForm.eName" size="mini"/>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row type="flex" justify="end">
          <el-col :span="4">
            <el-button type="primary" plain size="mini" @click="onSubmitQualifier">提交</el-button>
          </el-col>
        </el-row>
      </el-form>
    </el-dialog>

  </div>

</template>

<script>
  import {setRootCookie} from "../../conf/utils"

  const axios = require('axios')
  export default {
    data() {
      return {
        editId: 0,
        dataList: [],
        isCheck: false,
        tableColumn: [],
        editCombineIndicator: false,
        indicatorForm: {
          name: '',
          fieldName: '',
          meaning: '',
          expression: '',
          serviceTagIds: [],
          predictItems: [],
          qualifiers: [],
          timeQualifier: '',
        },
        indicatorFormRules: {
          'name': [{required: true, message: '请输入中文名', trigger: 'blur'}],
          'serviceTagIds': [{required: true, message: '请选择业务类型', trigger: ['blur', 'change']}],
          'fieldName': [{required: true, message: '请输入英文名', trigger: 'blur'}],
        },
        timeQualifierForm: {
          cName: '',
          eName: '',
          qualifierType: '',
        },
        timeQualifierFormRules: {
          'cName': [{required: true, message: '请输入中文名', trigger: 'blur'}],
          'eName': [{required: true, message: '请输入英文名', trigger: ['blur', 'change']}],
          // 'qualifierType': [{required: true, message: '请选择修饰词类型', trigger: 'blur'}],
        },
        qualifierVisible: false,
        serviceTagList: [],
        serviceTagMap: {},
        preIndicatorOptions: [],
        timeQualifierOption: [],
        indexDependencies: [],
        inputVisible: false,
        inputValue: '',
        fromDictionary: false,
      }
    },
    methods: {
      handleClose(tag) {
        this.indicatorForm.qualifiers.splice(this.indicatorForm.qualifiers.indexOf(tag), 1)
      },

      showInput() {
        this.inputVisible = true
        this.$nextTick(_ => {
          this.$refs['saveTagInput'].$refs['input'].focus()
        })
      },

      handleInputConfirm() {
        let inputValue = this.inputValue
        if (inputValue) {
          this.indicatorForm.qualifiers.push(inputValue)
        }
        this.inputVisible = false
        this.inputValue = ''
      },
      insertToExpression(content) {
        console.log()
        this.indicatorForm.expression += content
      },
      async addDependencies() {
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
      onSubmitQualifier() {
        this.$refs['timeQualifierForm'].validate((valid) => {
          if (valid) {
            let self = this
            let param = this.timeQualifierForm
            param.qualifierType = '时间周期'
            axios.post('/indicator/qualifier/add', param)
              .then(function (response) {
                if (response.data.Res) {
                  self.$message({
                    type: 'success',
                    message: '修饰词添加成功'
                  })
                  self.getAllIndicator()
                  self.timeQualifierForm = {
                    cName: '',
                    eName: '',
                    qualifierType: ''
                  }
                  self.qualifierVisible = false
                } else {
                  self.$message({
                    type: 'error',
                    message: response.data.Info
                  })
                }
              }).catch(function (error) {
              console.error(error)
            })
          }
        })
      },
      deleteDependencies(idx) {
        this.indexDependencies.splice(idx, 1)
      },
      onDeleteIndicator(row) {
        this.editId = row.id
        let self = this
        this.$confirm('确认删除该UDF记录吗', '注意', {
          distinguishCancelAndClose: true,
          confirmButtonText: '确定',
          cancelButtonText: '取消'
        }).then(() => {
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
        })
      },
      onClickEditIndicator(row) {
        this.indexDependencies = []
        // 正在编辑的 id
        this.editId = row.id
        this.indicatorForm.name = row.name
        this.indicatorForm.meaning = row.meaning
        this.indicatorForm.expression = row.expression
        this.indicatorForm.serviceTagIds = row.business_label.split(',')
        this.indicatorForm.qualifiers = row.qualifiers === '' ? [] : row.qualifiers.split(',')
        console.log(row)
        let self = this
        axios.get('/indicator/getIndicatorDependenciesById/' + this.editId)
          .then(function (response) {
            self.transStrToIndexDependencies(response.data)
          }).catch(function (error) {
          console.log(error)
        })
      },
      closeIndicatorDialog() {
        this.$router.back(-1)
      },
      getColumnWidth(field) {
        if (field === 'meaning') {
          return '50%'
        }
        if (field === 'tid') {
          return '20%'
        }
      },
      getAllIndicator() {
        let self = this
        axios.get('/indicator/getAllCombineIndicator')
          .then(function (response) {
            self.tableColumn = response.data.cols
            self.dataList = response.data.data
          }).catch(function (error) {
          console.log(error)
        })
      },
      onClickNewIndex() {
        this.indexDependencies = []
        this.indicatorForm.name = ''
        this.indicatorForm.meaning = ''
        this.indicatorForm.expression = ''
        this.indicatorForm.serviceTagIds = []
        this.indicatorForm.predictItems = []
        this.editId = 0
      },
      makeNewCombineIndicator() {
        let self = this
        this.$refs['indicatorForm'].validate((valid) => {
          if (valid) {
            let dependencyIndexIds = []
            this.indexDependencies.forEach(indicator => {
              dependencyIndexIds.push(indicator.id)
            })
            let params = {
              'name': this.indicatorForm.name,
              'meaning': this.indicatorForm.meaning,
              'expression': this.indicatorForm.expression,
              'BusinessLabels': this.indicatorForm.serviceTagIds,
              'dependencies': dependencyIndexIds,
              'qualifiers': this.indicatorForm.qualifiers,
              'fieldName': this.indicatorForm.fieldName,
              'timeQualifier': this.indicatorForm.timeQualifier,
            }
            if (this.arbitrateNewOrUpdate() === 'update') { // 更新
              params.id = this.editId
              axios.post('/indicator/postCombineIndicator', params, {responseType: 'json'})
                .then(function (response) {
                  if (response.data.Res) {
                    self.$message({
                      type: 'success',
                      message: '更新成功'
                    })
                  } else {
                    self.$message({
                      type: 'error',
                      message: response.data.Info
                    })
                  }
                }).catch(function (error) {
                console.log(error)
              })
            } else { // 新建
              axios.post('/indicator/postCombineIndicator', params, {responseType: 'json'})
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
              func()  // 执行回调函数
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
      getCombineIndicatorById(combineIndicatorId) {
        let self = this
        let params = {
          'originId': combineIndicatorId + '',
          'assetType': 'combine'
        }
        this.fromDictionary = true
        axios.post('/indicator/dictionary/searchById', params, {responseType: 'json'})
          .then(function (response) {
            self.editId = parseInt(response.data.id)
            self.getDependencies()
            self.indicatorForm.name = response.data.name
            self.indicatorForm.meaning = response.data.meaning
            self.indicatorForm.aggregation = response.data.aggregation
            self.indicatorForm.expression = response.data.expression
            self.indicatorForm.fieldName = response.data.fieldName
            self.fieldName = response.data.fieldName
            self.indicatorForm.timeQualifier = response.data.timeQualifier
            if (response.data.qualifiers !== '') {
              self.indicatorForm.qualifiers = response.data.qualifiers.split(',')
            }
            if (response.data.businessLabel !== '') {
              self.indicatorForm.serviceTagIds = response.data.businessLabel.split(',')
            }
          }).catch(function (error) {
          console.log(error)
        })
      },
      getTimeQualifierOption() {
        let self = this
        // todo  get
        axios.get('/indicator/qualifier/time/get')
          .then(function (response) {
            self.timeQualifierOption = response.data
          }).catch(function (err) {
            console.error(err)
        })
      },

      onMounted() {
        let combineIndicatorId = sessionStorage.getItem("combineIndicatorId")
        this.isCheck = (sessionStorage.getItem("combineIndicatorCheck") === 'true')
        sessionStorage.removeItem("combineIndicatorCheck")
        if ((combineIndicatorId !== null) && (combineIndicatorId !== undefined) && (combineIndicatorId !== '')) {
          this.getCombineIndicatorById(combineIndicatorId)
        }
      }
    },
    mounted: function () {
      this.getAllPreDefineIndex()
      this.getAllServiceTag(this.onMounted)
      this.getTimeQualifierOption()
    }
  }
</script>

<style scoped>

  .el-tag + .el-tag {
    margin-left: 10px;
  }

  .button-new-tag {
    margin-left: 10px;
    height: 32px;
    line-height: 30px;
    padding-top: 0;
    padding-bottom: 0;
  }

  .input-new-tag {
    width: 90px;
    margin-left: 10px;
    vertical-align: bottom;
  }
</style>

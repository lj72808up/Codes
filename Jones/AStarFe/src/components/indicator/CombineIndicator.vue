<template>
  <div>
    <div>
      <el-row>
        <h3 style="margin-bottom: 0px; margin-top: 5px; display:inline;">计算指标</h3>
      </el-row>
      <el-button type="primary" size="medium" style="margin-bottom: 10px" @click="onClickNewIndex()">新建计算指标
      </el-button>

      <el-table :data="dataList" stripe style="width: 100%" :border="false"
                :header-cell-style="{background:'#eee',color:'#606266'}">
        <el-table-column key="name" prop="name" label="指标名" min-width="12%"/>
        <el-table-column key="meaning" prop="meaning" label="指标含义" min-width="26%"/>
        <el-table-column key="expression" prop="expression" label="指标计算方式" min-width="26%"/>
        <el-table-column key="create_time" prop="create_time" label="创建时间" :sortable="true" min-width="12%"/>
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
  import router from "../../router"

  const axios = require('axios')
  export default {
    data() {
      return {
        editId: -1,
        dataList: [],
        tableColumn: [],
        newCombineIndicator: false,
        editCombineIndicator: false,
        topLabelReadOnly: false,
        indicatorForm: {
          name: '',
          meaning: '',
          expression: '',
          serviceTagIds: [],
          predictItems: [],
          qualifiers: [],
        },
        indicatorFormRules: {
          'name': [{required: true, message: '请输入标签名', trigger: 'blur'}],
          'serviceTagIds': [{required: true, message: '请选择业务类型', trigger: ['blur', 'change']}],
        },
        serviceTagList: [],
        serviceTagMap: {},
        preIndicators: [],
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
      addDependencies() {
        let selectIds = this.indicatorForm.predictItems
        this.transStrToIndexDependencies(selectIds)
        // 添加后还原select选择框的结果
        this.indicatorForm.predictItems = []
      },
      transStrToIndexDependencies(selectIds) {
        // 从preIndicators中找出selectIds的集合, 放入indexDependencies
        this.preIndicators.forEach(indicator => {
          if (selectIds.indexOf(indicator.id) !== -1) {
            this.indexDependencies.push(indicator)
          }
        })
      },
      deleteDependencies(idx) {
        this.indexDependencies.splice(idx, 1)
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
              self.getAllCombineIndicator()
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
      onClickEditIndicator(row) {
        let combineIndicatorId = row.id
        sessionStorage.setItem('combineIndicatorId', combineIndicatorId)
        router.push({path: 'combineIndicatorCreate'})
      },
      closeIndicatorDialog() {
        this.$refs['indicatorForm'].resetFields()
        this.$nextTick(() => {
          this.newCombineIndicator = false
          this.editCombineIndicator = false
        })
      },
      getColumnWidth(field) {
        if (field === 'meaning') {
          return '50%'
        }
        if (field === 'tid') {
          return '20%'
        }
      },
      getAllCombineIndicator() {
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
        sessionStorage.setItem('combineIndicatorId', '')
        this.$router.push({path: 'combineIndicatorCreate'})
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
    },
    mounted: function () {
      this.getAllServiceTag(this.getAllCombineIndicator)
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

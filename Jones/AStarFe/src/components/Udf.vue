<template>
  <div>
    <h4>UDF函数预览</h4>
    <vxe-table style="margin-left: 200px;margin-right: 200px"
               border
               show-overflow
               keep-source
               ref="existTable"
               :data="dataList"
               :edit-config="{trigger: 'dblclick', mode: 'cell', showStatus: true, icon: 'fa fa-pencil-square-o'}">

      <vxe-table-column align="center" field="funcName" title="函数名" :visible="true" width="330">
        <template style="text-align:center" slot-scope="oneFunc">
          <el-button @click="ShowUDFDetail(oneFunc.row)" type="text">{{oneFunc.row.funcName}}</el-button>
        </template>
      </vxe-table-column>
      <vxe-table-column align="center" field="description" title="功能描述" :visible="true"/>
      <vxe-table-column align="center" field="operation" title="操作" :visible="true" width="100">
        <template slot-scope="scope">
          <el-button @click="deleteById(scope.row.funcId)" type="text" size="mini" icon="el-icon-delete"/>
        </template>
      </vxe-table-column>
    </vxe-table>

    <h4>本次启动新增的UDF</h4>
    <vxe-table style="margin-left: 150px;margin-right: 150px"
               border
               show-overflow
               keep-source
               ref="newTable"
               :data="newDataList"
               :edit-config="{trigger: 'click', mode: 'cell', showStatus: true, icon: 'fa fa-pencil-square-o'}">

      <vxe-table-column align="center" field="funcName" title="函数名" :visible="true">
        <template style="text-align:center" slot-scope="oneFunc">
          <el-button @click="ShowUDFDetail(oneFunc.row)" type="text">{{oneFunc.row.funcName}}</el-button>
        </template>
      </vxe-table-column>
      <vxe-table-column align="center" field="description" title="功能描述" :visible="true"/>
    </vxe-table>

    <el-dialog height="600px" width="70%" title="函数详情" :visible.sync="funcVisible" @opened='openedDialog'
               @close='closeDialog'>
      <el-container>
        <el-header height="120px">
          <p style="font-size:30px;lineHeight:5px;margin-left: 30px">{{curUDF.funcName}}</p>
          <el-input v-if="curUDFModifyFlag" v-model="curUDF.description"></el-input>
          <span style="margin-top:10px;margin-left:50px" v-else>{{curUDF.description}}</span>
        </el-header>

        <el-container>
          <el-main>
            <h4>参数说明</h4>
            <el-input type="textarea" :autosize="{ minRows: 7, maxRows: 7}" :readOnly="!curUDFModifyFlag"
                      v-model="curUDF.parameters">
            </el-input>
          </el-main>
          <el-main>
            <h4>Hive Usage</h4>
            <codemirror v-model="curUDF.example" :options="cmOptions" ref="cmEditor"/>
          </el-main>
        </el-container>

        <el-footer>
          <el-row type="flex" justify="end" v-if="!topLabelReadOnly">
            <el-button size="small" style="margin-top: 10px" @click="modifyUdfDesc()">修改</el-button>
            <el-button size="small" style="margin-top: 10px" @click="updateUdfDesc()">提交</el-button>
          </el-row>
        </el-footer>
      </el-container>

    </el-dialog>


  </div>
</template>


<script>
  const axios = require('axios')

  export default {
    name: 'Udf',
    data () {
      return {
        curUDFModifyFlag: false,
        cmOptions: {
          // codemirror options
          tabSize: 4,
          mode: 'text/x-mysql',   // https://codemirror.net/mode/index.html
          theme: 'default',
          lineNumbers: true,
          dragDrop: false,
          readOnly: true,
          lineWrapping: true,
          line: true,
          hintOptions: {
            tables: {},
            databases: {}
          },
          extraKeys: {
            'Tab': 'autocomplete'
          }
        },
        // 已配置的udf
        dataList: [],
        tableColumn: [],
        // 本次启动新增的udf
        newDataList: [],
        newTableColumn: [],
        curUDF: {},
        funcVisible: false,
        showcode: '',
        topLabelReadOnly: false    // 顶栏标签权限是否限制了只读
      }
    },
    methods: {
      deleteById (funcId) {
        let self = this
        this.$confirm('确认删除该UDF记录吗', '注意', {
          distinguishCancelAndClose: true,
          confirmButtonText: '确定',
          cancelButtonText: '取消'
        }).then(() => {
          axios.get('/udf/delete/' + funcId)
            .then(function (response) {
              self.getUdfList()
              self.$message({
                type: 'success',
                message: '保存修改'
              })
            })
            .catch(function (error) {
              console.log(error)
            })
        }).catch(action => {
          console.log(action)
        })
      },
      modifyUdfDesc () {
        this.curUDFModifyFlag = true
        this.cmOptions.readOnly = false
      },
      closeDialog () {
        this.curUDFModifyFlag = false
        this.cmOptions.readOnly = true
        this.curUDF = {}
        this.showcode = ''
      },
      openedDialog () {
        this.showcode = Prism.highlight(this.curUDF.example, Prism.languages.sql, 'sql')
        //Prism.highlightAll();
        //var precode = document.getElementById('sqldata')
        //Prism.highlightElement(precode);
      },
      updateUdfDesc () {
        let updateRecords = [this.curUDF]

        let self = this
        axios.post('/udf/update/description', updateRecords, {responseType: 'json'})
          .then(function (response) {
            if (response.data.Res) {
              self.$message({
                message: response.data.Info,
                type: 'success'
              })
            } else {
              self.$message.error(response.data.Info)
            }
            self.getUdfList()
            self.getNewUdfList()
            self.funcVisible = false
          })
          .catch(function (error) {
            console.log(error)
          })
      },
      getNewUdfList () {
        let self = this
        let param = {
          'known': false
        }
        axios.post('/udf/getList', param, {responseType: 'json'})
          .then(function (response) {
            self.newTableColumn = response.data.tableColumn
            self.newDataList = response.data.dataList
          })
          .catch(function (error) {
            console.log(error)
          })
      },
      getUdfList () {
        let self = this
        let param = {
          'known': true
        }
        axios.post('/udf/getList', param, {responseType: 'json'})
          .then(function (response) {
            self.tableColumn = response.data.tableColumn
            self.dataList = response.data.dataList
          })
          .catch(function (error) {
            console.log(error)
          })
      },
      ShowUDFDetail (oneFunc) {
        let self = this
        //import lodash from 'lodash'
        //var deep = lodash.cloneDeep(objects);
        self.curUDF = Object.assign({}, oneFunc)
        console.log(self.curUDF)
        self.funcVisible = true
        return
      },
      getEdit () {
        let self = this
        axios.get('/udf/getEdit')
          .then(function (response) {
            self.topLabelReadOnly = response.data.readOnly
            console.log('顶栏标签配置的udf只读权限:' + self.topLabelReadOnly)
          })
          .catch(function (error) {
            console.log(error)
          })
      }
    },
    computed: {
      descriptionWidth: function (title) {
        if (title === 'description') {
          return '60%'
        } else {
          return ''
        }
      }
    },
    mounted () {
      this.getUdfList()
      this.getNewUdfList()
      this.getEdit()
    },
  }
</script>

<style scoped>
  .el-main >>> .CodeMirror {
    border: 1px solid #eee;
    height: auto;
    min-height: 150px;
    max-height: 150px;
  }

  .el-main >>> .CodeMirror-scroll {
    height: auto;
    min-height: 150px;
    max-height: 150px;
  }
</style>

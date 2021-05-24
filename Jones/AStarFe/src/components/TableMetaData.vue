/* eslint-disable */
<template>
  <div>
    <h3>元数据管理</h3>

    <!--    <el-form ref="tableForm" :model="tableForm" :inline="true">-->
    <el-row>
      <el-col span="9">
        类型:
        <el-select placeholder="请选择类型" @change="onChangeTableType()" size="mini" filterable="true"
                   v-model="conditionSelect.tableType"
                   clearable="true">
          <el-option
            v-for="v in conditionOptions.tableTypes"
            :key="v"
            :label="v"
            :value="v">
          </el-option>
        </el-select>
        <span>数据库:</span>
        <el-select placeholder="请选择数据库" @change="onChangeDataBase()" size="mini" filterable="true"
                   v-model="conditionSelect.database" clearable="true">
          <el-option
            v-for="v in conditionOptions.databases"
            :key="v"
            :label="v"
            :value="v">
          </el-option>
        </el-select>
      </el-col>

      <el-col :span="6" offset="7">
        模糊查询:
        <el-input size="mini" v-model="conditionSelect.tableNameLike" style="width:180px"
                  clearable="true"
                  @input="reInitPageTable()"/>
      </el-col>

    </el-row>


    <el-form ref="tableForm" :model="tableForm" :inline="true">
      <el-table
        v-loading="tableForm.loading"
        :data="tableForm.data"
        stripe
        style="width: 100%">
        <el-table-column
          v-for="col in tableForm.cols"
          v-if="!col.hidden"
          :key="col.prop"
          :prop="col.prop"
          :label="col.label"
          :sortable="col.sortable">
        </el-table-column>

        <el-table-column
          fixed="right"
          label="操作"
          width="140">
          <template slot-scope="scope">
            <el-row>
              <!--ADD 20210324 只有是owner的, 且有编辑权限的, 才能删除-->
              <el-button icon="el-icon-delete" type="text" size="small" @click="deleteTable(scope.row)"
                         v-show="(!readOnly) && (scope.row.is_owner)"/>
              <el-button icon="el-icon-setting" type="text" size="small" @click="editTable(scope.row)">详情</el-button>
              <el-button icon="el-icon-printer" type="text" size="small" @click="graphTable(scope.row)"/>
<!--              <span>{{scope.row.is_owner}}</span>-->
            </el-row>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="pageParam.currentPage"
        :page-sizes="[10, 50, 100]"
        :page-size="pageParam.pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 20px"
        :total="pageParam.totalItems">
      </el-pagination>
    </el-form>

    <div style="margin-top:80px" v-show="!readOnly">
      <h3>新增元数据</h3>
      <el-form ref="importForm" :model="importForm" :inline="true" label-width="120px">
        <el-row>
          <el-form-item label="类型:">
            <el-select placeholder="请选择类型" size="mini" filterable="true"
                       v-model="importForm.tableType"
                       @change="onChangeImportType()"
                       clearable="true">
              <el-option
                v-for="v in conditionOptions.tableTypes"
                :key="v"
                :label="v"
                :value="v">
              </el-option>
            </el-select>
          </el-form-item>

          <el-row>
            <el-form-item v-show="importForm.tableType !=='hdfs'" label="数据源:">
              <el-select placeholder="请选择数据源" size="mini" filterable="true"
                         v-model="importForm.importId"
                         @change="onChangeImportId"
                         clearable="true">
                <el-option
                  v-for="item in importForm.mysqlOptions"
                  :key="item.ImportName"
                  :label="item.ImportName"
                  :value="item.ImportId">
                </el-option>
              </el-select>
              <span style="color: rgb(189, 189, 189);">{{showUrl}}</span>
            </el-form-item>
          </el-row>

          <el-form-item v-show="importForm.tableType!=='' && importForm.tableType!=='hdfs'" label="输入数据库:">
            <el-input size="mini" v-model="importForm.databaseName"/>
          </el-form-item>

          <el-form-item v-show="importForm.tableType!=='' && importForm.tableType!=='hdfs'" label="输入表:">
            <el-input size="mini" v-model="importForm.tableName"/>
          </el-form-item>

          <el-form-item v-show="importForm.tableType==='hdfs'" label="文件名:">
            <el-input size="mini" v-model="importForm.tableName"/>
          </el-form-item>

          <el-form-item v-show="importForm.tableType==='hdfs'" label="文件存储位置:">
            <el-input size="mini" v-model="importForm.location"/>
          </el-form-item>

          <el-form-item v-show="true" label="有权使用的角色:">
            <el-select v-model="importForm.privilegeRoles" multiple filterable placeholder="请选择" size="small">
              <el-option
                v-for="value in roleOptions"
                :key="value"
                :label="value"
                :value="value">
              </el-option>
            </el-select>
          </el-form-item>

        </el-row>
        <el-row v-show="importForm.tableType==='hdfs'">
          <el-col span="18">
            <el-input type="textarea"
                      :rows="5"
                      :placeholder="importForm.hdfsPlaceHolder"
                      size="mini"
                      v-model="importForm.hdfsImportStr"/>
          </el-col>
        </el-row>

        <div style="margin-top: 20px" v-show="!readOnly">
          <el-row>
            <el-form-item>
              <el-button size="small" type="primary" @click="onImport()">提交导入</el-button>
              <el-button size="small" type="primary" @click="onSubmitNewTables()"
                         v-show="importForm.tableType!=='hdfs'">提交新增
              </el-button>
            </el-form-item>
          </el-row>
        </div>


      </el-form>

    </div>

  </div>

</template>

<style>
  .el-table .info-row {
    background: oldlace;
  }

  .el-table .success-row {
    background: #f0f9eb;
  }
</style>

<script>
  const axios = require('axios')
  import XLSX from 'xlsx'
  import FileSaver from 'file-saver'
  import router from '../router'
  import {getCookie, setCookie, formatDateLine, isRunFinsh, isRunFinshPagani, type2route} from '../conf/utils'

  export default {
    data () {
      return {
        readOnly: true,
        importSelection: ['hive', 'clickhouse', 'mysql'],
        conditionSelect: {
          tableType: '',
          database: '',
          tableNameLike: '',
        },

        importForm: {
          tableType: '',
          databaseName: '',
          tableName: '',
          importId: '',
          location: '',
          mysqlOptions: [],
          hdfsPlaceHolder: '请输入字段信息, 用逗号分隔: 字段名称,字段类型,业务信息描述',
          hdfsImportStr: '',
          privilegeRoles: [],
        },

        conditionOptions: {
          tableTypes: '',
          databases: ''
        },

        tableForm: {
          loading: true,
          data: [],
          cols: []
        },

        newTableForm: [{
          tableName: '',
          tableType: '',
          database: '',
          location: '',
          compressMode: '',
          description: ''
        }],

        tableSelect: {
          tid: 0,
          tableName: '',
          tableType: '',
          database: '',
          location: '',
          compressMode: '',
          description: '',
          creator: '',
          createTime: '',
        },
        dialogVisible: false,
        curTid: '',

        pageParam: {
          currentPage: 1,
          pageSize: 10,
          totalItems: 100
        },
        graphVisible: false,
        tableTypes: [],
        showUrl: '',
        roleOptions: [],
      }
    },
    methods: {
      getRoleOptions () {
        let self = this
        axios.get('/auth/roles')
          .then(function (response) {
            self.roleOptions = response.data
          })
          .catch(function (error) {
            console.log(error)
          })
      },

      onChangeImportId (checkValue) {
        console.log(checkValue)
        for (let opt of this.importForm.mysqlOptions) {
          if (opt.ImportId === checkValue) {
            this.showUrl = opt.ConnectionUrl
            break
          }
        }
      },

      onChangeImportType () {
        this.importForm.importId = ''
        this.showUrl = ''
        let tableType = this.importForm.tableType
        // if (tableType === 'mysql' || tableType === 'clickhouse') {
        let self = this
        axios.get('/metaData/tableImport/getConnections?connType=' + tableType).then(
          function (response) {
            self.importForm.mysqlOptions = response.data
          }
        ).catch(function (error) {
          console.log(error)
        })
        // }
      },

      onImport () {
        console.log(this.importForm.tableType)
        let self = this
        let params = {
          'tableType': this.importForm.tableType,
          'databaseName': this.importForm.databaseName,
          'tableName': this.importForm.tableName,
          'importId': this.importForm.importId + '',
          'hdfsImportStr': this.importForm.hdfsImportStr,
          'location': this.importForm.location,
          'privilegeRoles': this.importForm.privilegeRoles.join(','),
        }
        console.log(params)
        axios.post('metaData/import/single', params, {responseType: 'json'})
          .then(function (response) {
            let types = {
              'success': '导入成功',
              'error': '导入失败'
            }
            console.log(response.data)
            let resFlag = 'error'
            if (response.data.Res) {
              resFlag = 'success'
            }
            self.$message({
              type: resFlag,
              message: self.importForm.tableType + types[resFlag] + '; ' + response.data.Info
            })
            self.reInitPageTable()
          }).catch(function (error) {
          console.log(error)
        })
      },

      checkField (row) {
        let info = {
          'tableId': row.tid
        }
        sessionStorage.setItem('metaData_table_id', JSON.stringify(info))
        router.push({path: 'fieldMetaData'})
      },

      graphTable (row) {
        console.log(row)
        let info = {
          'tableId': row.tid,
          'tableName': row.table_name
        }
        sessionStorage.setItem('metaData_graph_tid', JSON.stringify(info))
        router.push({path: 'tableGraph'})
      },

      editTable (row) {
        console.log(row)
        this.curTid = row.tid
        this.tableSelect.tid = this.curTid
        this.tableSelect.tableName = row.table_name
        this.tableSelect.tableType = row.table_type
        this.tableSelect.database = row.database_name
        this.tableSelect.location = row.location
        this.tableSelect.compressMode = row.compress_mode
        this.tableSelect.description = row.description
        this.tableSelect.creator = row.creator
        this.tableSelect.createTime = row.create_time

        // 放入cache
        // sessionStorage.setItem('metaData_tableEdit', JSON.stringify(this.tableSelect))
        let info = {
          'tableId': row.tid
        }
        sessionStorage.setItem('metaData_table_id', JSON.stringify(info))
        // 跳转
        router.push({path: 'tableInfo'})
      },

      deleteTable (row) {
        console.log('tid: ' + row.tid)
        let self = this
        this.$confirm('确定删除该表吗?', '提示', {
          distinguishCancelAndClose: true,
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'info',
          center: true
        }).then(() => {
          axios.get('/metaData/table/delete/' + row.tid).then(function (response) {
            self.getTableTypes()

          }).catch(function (error) {
            console.log(error)
          })
        })
      },

      clearData () {
        this.newTableForm = [{
          tableName: '',
          tableType: '',
          database: '',
          location: '',
          compressMode: '',
          description: ''
        }]

        this.tableSelect = {
          tid: 0,
          tableName: '',
          tableType: '',
          database: '',
          location: '',
          compressMode: '',
          description: '',
          creator: '',
          create_time: '',
        }
      },

      checkTableName () {
        let flag = true
        flag = (this.importForm.tableName !== '') && (this.importForm.tableType !== '')
        // this.newTableForm.forEach(item => {
        //   if (item.tableName === '') {
        //     flag = false
        //   }
        // })
        return flag
      },
      onSubmitNewTables () {
        if (this.checkTableName()) {
          console.log(this.importForm)
          let params = [{
            tableName: this.importForm.tableName,
            tableType: this.importForm.tableType,
            database: this.importForm.databaseName,
            location: this.importForm.location,
            compressMode: '',
            description: '',
            importId: this.importForm.importId + '',
            privilegeRoles: this.importForm.privilegeRoles.join(',')
          }]
          let self = this
          axios.post('metaData/postTables', params, {responseType: 'json'})
            .then(function (response) {
              let resFlag = 'error'
              if (response.data.Res) {
                resFlag = 'success'
              }
              self.$message({
                type: resFlag,
                message: response.data.Info
              })
              if (response.data.Res) {
                self.getTableTypes()
                self.clearData()
              }
            }).catch(function (error) {
            console.log(error)
          })
        } else {
          this.$message({
            type: 'error',
            message: '请确保表名不为空'
          })
        }
      },
      removeNewTable (item) {
        let index = this.newTableForm.indexOf(item)
        if (index !== -1) {
          this.newTableForm.splice(index, 1)
        }
      },

      addNewTable () {
        this.newTableForm.push({
          tableName: '',
          tableType: '',
          database: ''
        })
      },
      onChangeTableType () {
        let self = this
        this.pageParam.currentPage = 1
        if (self.conditionSelect.tableType === '') {
          // 清空数据库选项
          self.conditionOptions.databases = ''
          self.conditionSelect.database = ''
        } else {
          axios.get('metaData/getDataBaseByType/' + this.conditionSelect.tableType)
            .then(function (response) {
              self.conditionOptions.databases = response.data
              self.conditionSelect.database = ''

            }).catch(function (error) {
            console.log(error)
          })
        }
        this.reInitPageTable()
      },

      onChangeDataBase () {
        this.pageParam.currentPage = 1
        this.reInitPageTable()
      },
      getTableTypes () {
        // 查询出有多少个类型
        let self = this
        axios.get('metaData/tableTypes/get')
          .then(function (response) {
            self.conditionOptions.tableTypes = response.data
            self.tableTypes = response.data
          }).catch(function (error) {
          console.log(error)
        })

        // 查出表格数据
        let param = {
          'curPage': this.pageParam.currentPage + '',
          'pageSize': this.pageParam.pageSize + '',
        }
        axios.post('metaData/getTableByConditions', param, {responseType: 'json'})
          .then(function (response) {
            self.tableForm.cols = response.data.cols
            self.tableForm.data = response.data.data

            self.tableForm.loading = false
          }).catch(function (error) {
          console.log(error)
        })

      },

      // 每次更该筛选条件后, 都可以调用reInitPage, 刷新总条数和表格数据
      reInitPageTable () {
        let self = this
        this.pageParam.currentPage = 1
        this.tableForm.loading = true
        let params = {
          'table_type': self.conditionSelect.tableType,
          'database_name': self.conditionSelect.database,
          'curPage': this.pageParam.currentPage + '',
          'pageSize': this.pageParam.pageSize + '',
          'tableNameLike': this.conditionSelect.tableNameLike,
        }
        axios.post('metaData/table/getNums', params, {responseType: 'json'})
          .then(function (response) {
            console.log('count:' + response.data.count)
            self.pageParam.totalItems = response.data.count
            self.updateTableData()
            this.tableForm.loading = false
          }).catch(function (error) {
          this.tableForm.loading = false
          console.log(error)
        })
      },

      initPageTable () {
        let self = this
        let params = {
          'table_type': self.conditionSelect.tableType,
          'database_name': self.conditionSelect.database,
          'curPage': this.pageParam.currentPage + '',
          'pageSize': this.pageParam.pageSize + '',
        }
        axios.post('metaData/table/getNums', params, {responseType: 'json'})
          .then(function (response) {
            console.log('count:' + response.data.count)
            self.pageParam.totalItems = response.data.count
            self.getTableTypes()
          }).catch(function (error) {
          console.log(error)
        })
      },

      handleSizeChange (pageSize) {
        console.log(`每页 ${pageSize} 条`)
        this.pageParam.pageSize = pageSize
        this.handleCurrentChange(1)
      },

      handleCurrentChange (curPage) {
        console.log(`当前页: ${curPage}`)
        this.pageParam.currentPage = curPage
        this.tableForm.loading = true
        this.updateTableData()
      },// 前台控件的page号从1开始

      updateTableData () {
        // 查出表格数据
        let param = {
          'curPage': this.pageParam.currentPage + '',
          'pageSize': this.pageParam.pageSize + '',
          'table_type': this.conditionSelect.tableType,
          'database_name': this.conditionSelect.database,
          'tableNameLike': this.conditionSelect.tableNameLike,
        }
        let self = this
        axios.post('metaData/getTableByConditions', param, {responseType: 'json'})
          .then(function (response) {
            self.tableForm.cols = response.data.cols
            self.tableForm.data = response.data.data

            self.tableForm.loading = false
          }).catch(function (error) {
          console.log(error)
        })
      }
    }
    ,
    // methods end
    mounted: function () {
      this.initPageTable()
      this.getRoleOptions()
      let self = this
      axios.get('/metaData/getEdit')
        .then(function (response) {
          self.readOnly = response.data.readOnly
          console.log(self.readOnly)
        })
        .catch(function (error) {
          console.log(error)
        })
    }

  }
</script>

<style scoped>
  .line {
    float: right;
    width: 100%;
    height: 1px;
    margin-top: -0.5em;
    margin-bottom: 0.5em;
    background: #d4c4c4;
    position: relative;
    text-align: center;
  }
</style>

<template>
  <div class="app-container">
    <el-form ref="form" :model="form" label-width="120px">
      <h3>基础信息</h3>

      <el-row>
        <el-col :span="6">
          <el-form-item label="表名称" size="small">
            <el-input :readonly="readOnly" v-model="form.tableName"/>
          </el-form-item>
        </el-col>

        <el-col :span="6">
          <el-form-item label="所属库" size="small">
            <el-input :readonly="readOnly" v-model="form.database"/>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col span="6">
          <el-form-item label="创建者" size="small">
            <el-input :readonly="readOnly" type="" v-model="form.creator"/>
          </el-form-item>
        </el-col>

        <el-col span="6">
          <el-form-item label="创建时间" size="small">
            <el-input :readonly="readOnly" type="" v-model="form.createTime"/>
          </el-form-item>
        </el-col>

        <el-col :span="6">
          <el-form-item label="表类型" size="small">
            <el-select v-model="form.tableType" placeholder="please select your zone">
              <el-option label="hive" value="hive"/>
              <el-option label="mysql" value="mysql"/>
              <el-option label="hdfs" value="hdfs"/>
              <el-option label="clickhouse" value="clickhouse"/>
              <el-option label="kylin" value="kylin"/>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col span="12">
          <el-form-item label="文件存储位置">
            <el-input :readonly="readOnly" type="" v-model="form.location" size="small"/>
          </el-form-item>
        </el-col>
        <el-col span="6">
          <el-form-item label="压缩方式">
            <el-input :readonly="readOnly" type="" v-model="form.compressMode" size="small"/>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col span="18">
          <el-form-item label="描述">
            <el-input :readonly="readOnly" type="textarea" rows=3 v-model="form.description" size="small"/>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="7">
          <el-form-item label="有权使用的角色" size="small">
            <el-select v-model="form.privilegeRoles" multiple filterable placeholder="请选择" size="small"
                       style="width: 100%" :readonly="readOnly">
              <el-option
                v-for="value in roleOptions"
                :key="value"
                :label="value"
                :value="value">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>

    <el-row type="flex" justify="end" v-show="!readOnly" style="width: 75%; ">
      <el-button @click="submitUpdate()" size="small" type="primary" plain style="margin-top: -30px">提交修改</el-button>
    </el-row>

    <!-- 字段信息 -->
    <div style="margin-top: 20px">
      <el-row >
        <el-col span="2">
          <h3 v-if="(!totallyNew)&&(versions !== null)&&(!readOnly)">编辑元数据</h3>
          <h3 v-else>查看元数据</h3>
        </el-col>
        <el-col :span="3">
          <el-button @click="onClickNewVersion()" icon="el-icon-circle-plus" size="small"
                     v-show="(!totallyNew)&&(versions !== null)&&(!readOnly)" style="margin-top: 20px"> 创建新版本
          </el-button>
        </el-col>
      </el-row>

      <el-form :inline="true" v-show="!totallyNew" label-width="120px">
        <el-form-item label="已存在版本">
          <el-select v-model="curVersion" placeholder="请选择版本" @change="onChangeVersion(curVersion)" size="mini">
            <el-option
              v-for="v in versions"
              :key="v"
              :label="v"
              :value="v">
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item v-show="!readOnly">
          <el-button size="small" :icon="buttonIcon" type="primary" plain @click="updateCurVersionField"
                     :disabled="isUpdateFieldNow">{{buttonText}}
          </el-button>
        </el-form-item>

        <el-row>
          <el-form-item label="版本更新描述">
            <el-input :readonly="readOnly" v-model="curVersionDesc" size="small" style="width: 480px"/>
            <el-button @click="updateVersionDesc()" v-show="!readOnly" size="small">更新描述</el-button>
          </el-form-item>
        </el-row>

        <!--        <h4>字段配置</h4>-->
        <div style="margin-left: 20px">
          <vxe-table
            border="inner"
            show-overflow
            keep-source
            ref="xTable"
            stripe="true"
            :data="dataList"
            :edit-config="{trigger: 'dblclick', mode: 'cell', showStatus: true, icon: 'fa fa-pencil-square-o'}">
            <!--显示字段-->
            <template v-for="(col,idx) in tableColumn">
              <vxe-table-column :field="col.field" :title="col.title" :visible="col.visible"
                                :edit-render="col.editRender" :width="col.field==='description'?'40%':''"/>
            </template>
            <!--固定列-->
            <vxe-table-column v-if="!readOnly" title="操作" width="100">
              <template v-slot="{ row, rowIndex }">
                <el-button icon="el-icon-delete" type="text" size="mini" @click="deleteField(row)" v-show="!readOnly"/>
                <el-button icon="el-icon-setting" type="text" size="mini" @click="editRelation(row)">配置</el-button>
              </template>
            </vxe-table-column>
          </vxe-table>

          <el-row v-show="!readOnly">
            <el-col>
              <vxe-toolbar>
                <template v-slot:buttons>
                  <vxe-button @click="getUpdateEvent">提交字段修改</vxe-button>
                </template>
              </vxe-toolbar>
            </el-col>
          </el-row>
        </div>
      </el-form>

      <el-row>
        <h4>分区字段</h4>
        <el-col span="15" style="margin-left: 20px">
          <vxe-table
            border="inner"
            show-overflow
            keep-source
            ref="partitionTable"
            stripe="true"
            :data="partitionDataList"
            :edit-config="{trigger: 'dblclick', mode: 'cell', showStatus: true, icon: 'fa fa-pencil-square-o'}"
            @edit-closed="updatePartitionField">
            <!--固定列-->
            <vxe-table-column field="fid" title="字段id" :visible="false"/>
            <vxe-table-column field="fieldName" title="字段名称"/>
            <vxe-table-column field="fieldType" title="字段类型"/>
            <vxe-table-column field="description" title="描述" :edit-render="{name: 'input'}"/>
          </vxe-table>
        </el-col>
      </el-row>

      <el-form :inline="true" v-show="(!totallyNew)&&(versions !== null)&&(!readOnly)" label-width="120px"
               :rules="newFieldRules" ref="newFieldForm" :model="newField">
        <el-row>
          <div style="margin-bottom: 100px;margin-top: 30px">
            <el-col :span="4">
              <h4>当前版本下添加字段</h4>
            </el-col>
          </div>
        </el-row>

        <el-row>
          <el-form-item label="字段名称" prop="fieldName">
            <el-input :readonly="readOnly" v-model="newField.fieldName" size="small"/>
          </el-form-item>
          <el-form-item label="字段类型" prop="fieldType">
            <el-input :readonly="readOnly" v-model="newField.fieldType" size="small"/>
          </el-form-item>
        </el-row>

        <el-row>
          <el-form-item label="业务信息描述">
            <el-input :readonly="readOnly" v-model="newField.description" type="textarea" style="width:500px"/>
          </el-form-item>
          <el-button @click="addFieldInCurrentVersion()" size="small" v-show="!readOnly">添加</el-button>
        </el-row>
      </el-form>

      <!--全新版本-->

      <el-form ref="newVersionForm" :model="newVersionForm" :inline="true" v-show="totallyNew">
        <el-row>
          <el-form-item label="数据物理表:">
            <el-input :readonly="readOnly" size="mini" disabled v-model="tableName"/>
          </el-form-item>
          <el-form-item label="输入版本号:">
            <el-input :readonly="readOnly" v-model="newVersionForm.Version" size="mini"/>
          </el-form-item>
        </el-row>
        <el-row>
          <el-form-item
            v-for="(totallyNewField, index) in newVersionForm.Fields"
            :label="'字段' + index"
            :rules="{
              required: true, message: '域名不能为空', trigger: 'blur'
            }">
            <!--  :key="mean.key"
              :prop="'domains.' + index + '.value'"
              :rules="{
                required: true, message: '域名不能为空', trigger: 'blur'
              }"
            -->
            <el-row>
              <el-col span="4">
                <el-button @click="addTotallyNewField()" type="text" icon="el-icon-circle-plus-outline"
                           style=" font-size:16px; color: blue"/>
                <el-button @click.prevent="removeTotallyNewField(totallyNewField)"
                           icon="el-icon-remove-outline" type="text" style=" font-size:16px; color: red"/>
              </el-col>
            </el-row>

            <el-row>
              <el-form-item label="字段名称">
                <el-input :readonly="readOnly" v-model="totallyNewField.fieldName" size="small"/>
              </el-form-item>

              <el-form-item label="字段类型">
                <el-input :readonly="readOnly" v-model="totallyNewField.fieldType" size="small"/>
              </el-form-item>

            </el-row>

            <el-row>
              <el-form-item label="业务信息描述">
                <el-input :readonly="readOnly" v-model="totallyNewField.description" type="textarea"
                          style="width:500px"/>
              </el-form-item>
            </el-row>
          </el-form-item>
        </el-row>

        <el-row>
          <el-form-item>
            <!--            <el-button size="small" @click="addTotallyNewField()">添加字段</el-button>-->
            <!--<el-button size="small" @click="resetTotallyNewField('newVersionForm')">重置</el-button>-->
            <el-button size="small" type="primary" @click="onSubmitTotallyNewField()">提交</el-button>
            <el-button size="small" type="text" @click="totallyNew = !totallyNew">返回</el-button>
          </el-form-item>
        </el-row>
      </el-form>

      <el-dialog
        title="配置字段的关联关系"
        width="50%"
        :visible.sync="dialogVisible"
      >

        <template>
          <div class="container" v-show="!readOnly">
            <div class="large-12 medium-12 small-12 cell">
              <input type="file" id="file" ref="file" v-on:change="handleFileUpload()"/>
              <button v-on:click="submitFile()">上传关系文件</button>
              <span style="color:#BDBDBD;font-size:12px;">{{info}}</span>
            </div>
          </div>
        </template>

        <h2></h2>
        <el-form ref="relationSelect"
                 :model="relationSelect"
                 :rules="relationRules"
                 :inline="true"
                 label-width="80px"
                 v-show="!readOnly">
          <el-form-item label="类型:">
            <el-select placeholder="请选择类型" @change="onChangeRelationTableType()" size="mini" filterable="true"
                       v-model="relationSelect.typeSelect">
              <el-option
                v-for="v in relationOptions.tableTypes"
                :key="v"
                :label="v"
                :value="v">
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="数据库:">
            <el-select placeholder="请选择数据库" @change="onChangeRelationDataBase()" size="mini" filterable="true"
                       v-model="relationSelect.dataBaseSelect">
              <el-option
                v-for="v in relationOptions.databases"
                :key="v"
                :label="v"
                :value="v">
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="表:">
            <el-select placeholder="请选择表" @change="onChangeRelationTable()" size="mini" filterable="true"
                       v-model="relationSelect.tableSelect">
              <el-option
                v-for="v in relationOptions.tables"
                :key="v.Tid"
                :label="v.TableName"
                :value="v.Tid">
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="版本:">
            <el-select placeholder="请选择版本" @change="onChangeRelationVersion()" size="mini" filterable="true"
                       v-model="relationSelect.versionSelect" clearable>
              <el-option
                v-for="v in relationOptions.versions"
                :key="v.Version"
                :label="v.Version"
                :value="v.Version">
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="字段:">
            <el-select placeholder="请选择字段" @change="" size="mini" filterable="true"
                       v-model="relationSelect.fieldSelect">
              <el-option
                v-for="v in relationOptions.fields"
                :key="v.fid"
                :label="v.fieldName"
                :value="v.fid">
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="关系描述:">
            <el-input :readonly="readOnly" v-model="relationSelect.relationDesc"/>
          </el-form-item>

          <el-row>
            <el-col span="5" offset="19">
              <el-button round type="primary" size="mini"
                         @click="addRelation('relationSelect')"
                         v-show="!readOnly">添加
              </el-button>
            </el-col>
          </el-row>

        </el-form>

        <el-table
          v-loading="relationTable.loading"
          :data="relationTable.data"
          stripe
          style="width: 100%">
          <el-table-column
            v-for="col in relationTable.cols"
            v-if="!col.hidden"
            :key="col.prop"
            :prop="col.prop"
            :label="col.label"
            :sortable="col.sortable">
          </el-table-column>

          <el-table-column
            fixed="right"
            label="操作">
            <template slot-scope="scope">
              <el-button icon="el-icon-delete" type="info" size="mini" round @click="deleteRelation(scope.row)"/>
            </template>
          </el-table-column>

        </el-table>
      </el-dialog>

    </div>


  </div>

</template>

<script>
  import {formatDateLine, getCookie} from '../conf/utils'

  const axios = require('axios')
  export default {
    data () {
      return {
        newFieldRules: {
          fieldName: [
            {required: true, message: '请输入字段名称'},
            // {min: 3, max: 5, message: '长度在 3 到 5 个字符'}
          ],
          fieldType: [
            {required: true, message: '请输入字段类型'}
          ]
        },
        readOnly: true,
        roleOptions: [],
        partitionDataList: '',
        form: {
          tableName: '',
          tableType: '',
          database: '',
          location: '',
          compressMode: '',
          description: '',
          tid: '',
          creator: '',
          createTime: '',
          privilegeRoles: [],
        },

        curTid: '',

        info: 'csv内容形如: mysql,default_database,table1,field1,v0,我是描述',
        taskSwitch: false,
        tableColumn: [],
        dataList: [],
        taskForm: {
          session: '',
          date: formatDateLine(new Date()),
          status: ''
        },
        contentCols: [],
        contentData: [],
        contentFiled: '',
        dialogTableVisible: false,
        fieldLoading: false,
        dataLoading: false,
        dataPathRoot: '',
        dataPathRootPaga: '',
        dialogVisible: false,
        errMsg: '',

        // start
        versions: [],
        curVersion: '',
        curVersionDesc: '',
        versionDescList: {},

        userName: '',
        tableId: 2,
        tableName: '',

        newField: {
          fieldName: '',
          fieldType: '',
          description: '',
        },

        totallyNew: false,
        newVersion: '',
        newVersionForm: {
          Version: '',
          Fields: [{
            fieldName: '',
            fieldType: '',
            description: ''
          }]
        },
        isUpdateFieldNow: false,

        relationTable: {
          loading: false,
          data: [],
          cols: []
        },

        relationOptions: {
          tableTypes: [],
          databases: [],
          tables: [],
          versions: [],
          fields: []
        },
        relationSelect: {
          typeSelect: '',
          dataBaseSelect: '',
          tableSelect: '',
          versionSelect: '',
          fieldSelect: '',
          relationDesc: ''
        },

        curFieldId: '',
        relationRules: {
          label: [
            {required: true, message: '请输入关系描述', trigger: 'blur'},
            //{min: 3, message: '名称长度至少 3 个字符', trigger: 'blur'}
          ],
          // sql: [
          //   {require: true, message: 'SQL语句不能为空', trigger: 'blur'},
          //   {validator: validateSql, trigger: 'blur'}
          // ]
        },
        file: '',
      }
    },
    computed: {
      buttonIcon: function () {
        if (this.isUpdateFieldNow) {
          return 'el-icon-loading'
        } else {
          return 'el-icon-upload'
        }
      },
      buttonText: function () {
        if (this.isUpdateFieldNow) {
          return '更新中'
        } else {
          return '更新当前版本'
        }
      }
    },
    methods: {
      updateCurVersionField () {
        let version = this.curVersion
        let tid = this.tableId

        this.isUpdateFieldNow = true

        let param = {
          'version': version,
          'tableId': tid + '',
        }

        let self = this
        axios.post('/metaData/updateTableByVersion', param, {responseType: 'json'})
          .then(function (response) {
              if (response.data.Res) {
                self.$message({
                  message: '元数据更新成功',
                  type: 'success'
                })
                self.loadFieldMetaData()
              } else {
                self.$message({
                  message: response.data.Info,
                  type: 'error'
                })
              }
              self.isUpdateFieldNow = false
            }
          )
          .catch(function (error) {
            self.isUpdateFieldNow = false
            self.$options.methods.saveUpdateNotify.bind(self)()
            console.log(error)
          })
      },
      updatePartitionField ({row, column}) {
        let partitionTable = this.$refs.partitionTable
        // console.log(row)
        let field = column.property  // 属性名
        let self = this
        let cellValue = row[field]   // 属性值
        // 判断单元格值是否被修改
        if (partitionTable.isUpdateByRow(row, field)) {
          axios.post('/metaData/field/updatePartition', row, {responseType: 'json'})
            .then(function (response) {
              if (response.data.Res) {
                self.$message({
                  message: response.data.Info,
                  type: 'success'
                })
                self.$refs.partitionTable.reloadRow(row, null, field)
              } else {
                self.$message({
                  message: response.data.Info,
                  type: 'error'
                })
              }
            })
            .catch(function (error) {
              self.$options.methods.saveUpdateNotify.bind(self)()
              console.log(error)
            })
        }

        // setTimeout(() => {
        //   this.$message({
        //     message: '修改成功!'+ field+":"+cellValue,
        //     type: 'success'
        //   })
        //   // 局部更新单元格为已保存状态
        //   this.$refs.partitionTable.reloadRow(row, null, field)
        // }, 300)
      },
      submitUpdate () {
        console.log(this.form.tid)
        let params = {
          'tid': parseInt(this.form.tid),
          'tableName': this.form.tableName,
          'tableType': this.form.tableType,
          'databaseName': this.form.database,
          'description': this.form.description,
          'location': this.form.location,
          'compressMode': this.form.compressMode,
          'creator': this.form.creator,
          'createTime': this.form.createTime,
          'privilegeRoles': this.form.privilegeRoles.join(','),
        }
        let self = this
        axios.post('metaData/table/update', params, {responseType: 'json'})
          .then(function (response) {
            self.$message({
              message: '修改成功!',
              type: 'success'
            })
          }).catch(function (error) {
          console.log(error)
        })
      },

      updateVersionDesc () {
        console.log(this.tableId + ':' + this.curVersion)
        let param = {
          'Description': this.curVersionDesc,
          'Version': this.curVersion,
          'TableId': this.tableId
        }
        let self = this
        axios.post('/metaData/version/update', param, {responseType: 'json'})
          .then(function (response) {
            self.$message({
              message: '更新成功!',
              type: 'success'
            })
          })
          .catch(function (error) {
            console.log(error)
          })
      },

      onClickNewVersion () {
        let self = this
        axios.get('metaData/fields/version/getNewVersion/' + this.tableId)
          .then(function (response) {
            self.newVersionForm.Version = response.data
            // self.totallyNew=!self.totallyNew

            console.log(self.versions[0])
            let params = {
              'tableId': self.tableId + '',
              'version': self.versions[0]
            }
            axios.post('/metaData/fields/getListByVersion', params, {responseType: 'json'})
              .then(function (response) {
                console.log(response.data)
                if (response.data !== null) {
                  self.newVersionForm.Fields = response.data
                }
              })
              .catch(function (error) {
                console.log(error)
              })

            self.totallyNew = !self.totallyNew

          }).catch(function (error) {
          console.log(error)
        })
      },
      handleFileUpload () {
        this.file = this.$refs.file.files[0]
        console.log(this.file)
      },
      submitFile () {
        let self = this
        if (this.file) {
          this.$confirm('请确保提交的文件中包含6列:关联的表类型,关联的数据库, 关联的表名,关联的字段名,关联的版本号,关系描述', '提示', {
            distinguishCancelAndClose: true,
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
            center: true
          }).then(() => {
            let formData = new FormData()
            formData.append('file', this.file)
            formData.append('tableAId', this.tableId)
            formData.append('fieldAId', this.curFieldId)
            formData.append('tableAVersion', this.curVersion)

            axios.post('/metaData/upload/relation/file', formData,
              {
                headers: {
                  'Content-Type': 'multipart/form-data'
                }
              }
            ).then(function (response) {
              console.log(response.data)
              // self.uploadTableName = response.data
              self.$message({
                type: 'success',
                message: response.data
              })
              self.getRelation()
            })
              .catch(function (error) {
                console.log('失败')
                console.log(error)
                self.$message({
                  type: 'success',
                  message: '提交失败' + error
                })
              })
          })

        } else {
          // alert('先选择文件')
          this.$confirm('请先选择带查询词的csv文件, 第一行为表头, 字段用半角逗号分隔', '提示', {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning',
            center: true
          }).then(() => {
          })
        }
      },

      deleteRelation (row) {
        console.log(row.rid)
        let self = this
        let params = {
          'rid': row.rid
        }
        axios.post('metaData/fields/relation/delete', params, {responseType: 'json'})
          .then(function (response) {
            self.getRelation()
          }).catch(function (error) {
          console.log(error)
        })
      },

      addRelation (formName) {
        let self = this
        // this.$refs[formName].validate((valid) => {
        //   if (valid) {
        //调用添加关系的接口
        let params = {
          'tableAId': this.tableId + '',
          'fieldAId': this.curFieldId + '',
          'versionA': this.curVersion,
          'tableBId': this.relationSelect.tableSelect + '',
          'fieldBId': this.relationSelect.fieldSelect + '',
          'versionB': this.relationSelect.versionSelect,
          'description': this.relationSelect.relationDesc
        }
        console.log(JSON.stringify(params))
        axios.post('metaData/fields/relation/update', params, {responseType: 'json'})
          .then(function (response) {
            self.getRelation()
          }).catch(function (error) {
          console.log(error)
        })
        return true
        //   } else {
        //     //console.log("校验未通过")
        //     return false
        //   }
        // })
      },

      onChangeRelationVersion () {
        let self = this
        let tableId = this.relationSelect.tableSelect
        let version = this.relationSelect.versionSelect

        let params = {
          'tableId': tableId + '',
          'version': version
        }
        axios.post('metaData/fields/getByVersion', params, {responseType: 'json'})
          .then(function (response) {

            self.relationOptions.fields = response.data.dataList
          }).catch(function (error) {
          console.log(error)
        })

      },
      onChangeRelationTable () {
        let self = this
        let tableId = this.relationSelect.tableSelect
        axios.get('metaData/fields/version/get/' + tableId)
          .then(function (response) {
            self.relationOptions.versions = response.data.Versions
          }).catch(function (error) {
          console.log(error)
        })
      },

      onChangeRelationDataBase () {
        let self = this
        let params = {
          'tableType': self.relationSelect.typeSelect,
          'database': self.relationSelect.dataBaseSelect
        }
        axios.post('metaData/getTablesByDataBase', params, {responseType: 'json'})
          .then(function (response) {
            /*let tables = []
            for (let i=0; i<response.data.length; i++){
              tables.push(response.data[i].TableName)
            }*/

            self.relationOptions.tables = response.data
          }).catch(function (error) {
          console.log(error)
        })
      },

      onChangeRelationTableType () {
        // TODO
        let self = this
        axios.get('metaData/getDataBaseByType/' + this.relationSelect.typeSelect)
          .then(function (response) {
            self.relationOptions.databases = response.data
            self.relationSelect.dataBaseSelect = ''
            self.relationSelect.tableSelect = ''
            self.relationSelect.fieldSelect = ''

          }).catch(function (error) {
          console.log(error)
        })
      },
      getRelation () {
        let self = this
        // 查询关联关系
        axios.get('metaData/fields/relation/get/' + self.curFieldId)
          .then(function (response) {
            self.relationTable.data = response.data.data
            self.relationTable.cols = response.data.cols
            self.relationTable.loading = false
          }).catch(function (error) {
          console.log(error)
        })
      },
      editRelation (row) {
        this.dialogVisible = true
        this.relationTable.loading = true
        let self = this
        // console.log(row)
        this.curFieldId = row.fid + ''
        // 查询出有多少个类型
        axios.get('metaData/tableTypes/get')
          .then(function (response) {
            self.relationOptions.tableTypes = response.data
          }).catch(function (error) {
          console.log(error)
        })
        self.getRelation()
      },

      onSubmitTotallyNewField () {
        this.newVersionForm.TableId = this.tableId + ''
        this.newVersionForm.Fields = JSON.stringify(this.newVersionForm.Fields) + ''
        this.newVersionForm.TotallyNew = 'true'
        // this.$message(JSON.stringify(this.newVersionForm))
        let params = JSON.stringify(this.newVersionForm)
        console.log(params)
        let self = this
        axios.post('/metaData/fields/newVersion/put', params, {responseType: 'json'})
          .then(function (response) {
            // console.info(response.data)
            // self.$message({
            //   type: 'success',
            //   message: '元数据添加成功!'
            // })
            self.getAllVersion()
            self.clearNewField()
            self.totallyNew = false
          })
          .catch(function (error) {
            self.$options.methods.saveUpdateNotify.bind(self)()
            console.log(error)
          })
        //
      },
      resetTotallyNewField (formName) {
        this.$refs[formName].resetFields()
        this.newVersionForm = {
          Version: '',
          Fields: [{
            fieldName: '',
            fieldType: '',
            description: ''
          }]
        }
      },
      removeTotallyNewField (item) {
        let index = this.newVersionForm.Fields.indexOf(item)
        if (index !== -1) {
          this.newVersionForm.Fields.splice(index, 1)
        }
      },
      addTotallyNewField () {
        this.newVersionForm.Fields.push({
          fieldName: '',  // 词性
          fieldType: '', // 中文意义
          description: '' // 例句
        })
      },

      onChangeVersion () {
        console.log('curVersion:' + this.curVersion + ', versions:' + this.versions)
        this.curVersionDesc = this.versionDescList[this.curVersion]
        this.loadFieldMetaData()
      },

      deleteField (row) {
        let fid = row.fid
        let self = this
        console.log('fid:' + fid)

        this.$confirm('确定删除该字段吗?', '提示', {
          distinguishCancelAndClose: true,
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'info',
          center: true
        }).then(() => {
          let params = {
            'fieldId': fid + '',
          }
          console.log(params)
          axios.post('/metaData/fields/delete', params, {responseType: 'json'})
            .then(function (response) {
              console.info(response.data)
              self.$message({
                type: 'success',
                message: '字段已删除!'
              })
              self.loadFieldMetaData()
            })
            .catch(function (error) {
              self.$options.methods.saveUpdateNotify.bind(self)()
              console.log(error)
            })
        })
      },
      getAllVersion () {
        let self = this
        let tableId = self.tableId
        // 基本信息
        axios.get('metaData/getTableById/' + tableId)
          .then(function (response) {
            self.form.tableName = response.data.TableName
            self.form.tableType = response.data.TableType
            self.form.database = response.data.DatabaseName
            self.form.location = response.data.Location
            self.form.compressMode = response.data.CompressMode
            self.form.description = response.data.Description
            self.form.tid = response.data.Tid
            self.form.creator = response.data.Creator
            self.form.createTime = response.data.CreateTime
          }).catch(function (error) {
          console.log(error)
        })
        // 字段信息
        axios.get('metaData/fields/version/get/' + tableId)
          .then(function (response) {
            let descLists = []
            let descMap = {}
            let localversions = []
            response.data.Versions.forEach(function (item) {
              localversions.push(item.Version)
              descMap[item.Version] = item.Description
              descLists.push(item.Description)
            })
            console.log(self.versions)
            self.curVersionDesc = descLists[0]
            self.versions = localversions
            self.versionDescList = descMap
            self.tableName = response.data.TableName
            self.curVersion = self.versions[0]
            console.log('versions:' + self.versions + ';curVersion:' + self.curVersion)

            self.userName = getCookie('_adtech_user')
            self.fieldLoading = true
            let params = {
              'tableId': self.tableId + '',
              'version': self.curVersion,
            }
            console.log(params)
            axios.post('/metaData/fields/getByVersion', params, {responseType: 'json'})
              .then(function (response) {
                self.fieldLoading = false
                self.tableColumn = response.data.tableColumn
                self.dataList = response.data.dataList
              })
              .catch(function (error) {
                console.log(error)
              })
            // 分区信息
            axios.post('/metaData/fields/getPartitionByVersion', params, {responseType: 'json'})
              .then(function (response) {
                self.partitionDataList = response.data
              })
              .catch(function (error) {
                console.log(error)
              })

          }).catch(function (error) {
          console.log(error)
        })

      },

      loadFieldMetaData () {

        let self = this
        self.fieldLoading = true
        let params = {
          'tableId': self.tableId + '',
          'version': self.curVersion,
        }
        console.log(params)
        axios.post('/metaData/fields/getByVersion', params, {responseType: 'json'})
          .then(function (response) {
            self.fieldLoading = false
            self.tableColumn = response.data.tableColumn
            self.dataList = response.data.dataList
          })
          .catch(function (error) {
            console.log(error)
          })

        // 分区信息
        axios.post('/metaData/fields/getPartitionByVersion', params, {responseType: 'json'})
          .then(function (response) {
            self.partitionDataList = response.data
          })
          .catch(function (error) {
            console.log(error)
          })

      },

      addFieldInCurrentVersion () {
        console.log(this.newField)
        let self = this
        let params = {
          'TableId': self.tableId + '',
          'Version': self.curVersion,
          'TotallyNew': 'false',
          'Fields': '[{"fieldName":"' + self.newField.fieldName + '","fieldType":"' + self.newField.fieldType + '","description":"' + self.newField.description + '"}]'
        }
        // console.log(params)
        console.log(this.$refs['newFieldForm'])
        this.$refs['newFieldForm'].validate((valid) => {
          if (valid) {
            axios.post('/metaData/fields/newVersion/put', params, {responseType: 'json'})
              .then(function (response) {
                console.info(response.data)
                self.$message({
                  type: 'success',
                  message: '元数据添加成功!'//+self.form.name+':'+this.form.sql+':'+this.curTmplateId
                })
                self.clearNewField()
                self.loadFieldMetaData()
              })
              .catch(function (error) {
                self.$options.methods.saveUpdateNotify.bind(self)()
                console.log(error)
              })
          } else {
            console.log('没验证通过啊')
            return false
          }
        })

      },

      clearNewField () {
        this.newField.fieldName = ''
        this.newField.fieldType = ''
        this.newField.description = ''
        this.resetTotallyNewField('newVersionForm')
      },

      getUpdateEvent () {
        let updateRecords = this.$refs.xTable.getUpdateRecords()
        console.log(JSON.stringify(updateRecords))
        let params = JSON.stringify(updateRecords)
        let self = this
        this.$confirm('确定提交修改吗? (提交过程较长,请等待提示后再退出页面)', '提交修改', {
          confirmButtonText: '确定',
          cancelButtonText: '取消'
        }).then(() => {
            axios.post('/metaData/field/update', params, {responseType: 'json'})
              .then(function (response) {
                if (response.data.Res) {
                  self.$message({
                    message: response.data.Info,
                    type: 'success'
                  })
                  self.loadFieldMetaData()
                } else {
                  self.$message({
                    message: response.data.Info,
                    type: 'error'
                  })
                }
              })
              .catch(function (error) {
                self.$options.methods.saveUpdateNotify.bind(self)()
                console.log(error)
              })
          }
        ).catch(() => {
        })
      },
      getRoleOptions () {
        let self = this
        axios.get('/auth/roles')
          .then(function (response) {
            self.roleOptions = response.data
          })
          .catch(function (error) {
            console.log(error)
          })
        axios.get('/metaData/getTablePrivilegeByTid/' + this.tableId)
          .then(function (response) {
            let res = eval(response.data)
            self.form.privilegeRoles = res
            console.log(res)
          })
          .catch(function (error) {
            console.log(error)
          })
      }
    },

    mounted () {
      // const cacheInfo = sessionStorage.getItem('metaData_tableEdit')
      // sessionStorage.removeItem('metaData_tableEdit')
      // this.form = JSON.parse(cacheInfo)
      // console.log(this.form)

      let param = JSON.parse(sessionStorage.getItem('metaData_table_id'))
      // sessionStorage.removeItem('metaData_table_id')
      console.log('tableId:' + param.tableId + ' versions:' + this.versions)
      this.tableId = param.tableId
      let self = this
      axios.get('/metaData/getOwnerEdit?tid='+this.tableId)
        .then(function (response) {
          self.readOnly = response.data.readOnly
          console.log(self.readOnly)
        })
        .catch(function (error) {
          console.log(error)
        })

      this.getAllVersion()
      this.getRoleOptions()
    }

  }
</script>

<style scoped>
  /*.line {*/
  /*  text-align: center;*/
  /*}*/

  .line {
    float: right;
    width: 100%;
    height: 1px;
    margin-top: 1em;
    margin-bottom: 0.5em;
    background: #d4c4c4;
    position: relative;
    text-align: center;
  }

  /*  .titleBig {
      font-size: 1.17em;
      margin-block-start: 1em;
      margin-block-end: 1em;
      margin-inline-start: 0px;
      margin-inline-end: 0px;
      font-weight: bold;
    }*/

</style>

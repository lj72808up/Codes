<template>
  <div>
    <h3>增加一级域</h3>
    <el-form inline="true" label-width="120px">
      <el-form-item label="请输入域名">
        <el-input v-model="newDomainName" size="small"/>
      </el-form-item> <br/>
      <el-form-item label="管理员(逗号分隔)" style="margin-top: -20px">
        <el-input v-model="newDomainAdministrator" size="small"/>
<!--        <el-button type="primary" @click="onCreateDomain()" size="small">创建</el-button>-->
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onCreateDomain()" size="small">创建</el-button>
      </el-form-item>
    </el-form>
    <h3>配置域名管理员</h3>
    <el-row>
      <el-table
        :data="domainData"
        stripe
        style="width: 100%">
        <el-table-column
          prop="Name"
          label="域名">
        </el-table-column>
        <el-table-column
          label="用户操作">
          <template slot-scope="role">
            <el-button @click="renameAdministrator(role.row)" type="text" size="mini">修改管理员</el-button>
          </template>
        </el-table-column>
        <el-table-column
          label="绿皮配置">
          <template slot-scope="scope">
            <el-button @click="getCheckedURL(scope.row.Did,scope.row.Name)" type="text" size="mini">查看包含绿皮</el-button>
          </template>
        </el-table-column>

      </el-table>
    </el-row>

    <el-dialog :title="''+ curURLGroup+' 域包含的绿皮配置'" :visible.sync="lvpiOptionVisable">
          <el-row>
            <div style="max-height: 600px ; overflow:auto">
              <el-table
                :data="lvpiOptions"
                ref="lvpiOptionTable"
                row-key="Lid"
                @selection-change="handleLvpiSelectionChange">
                <el-table-column property="title" label="图表名称"/>
                <el-table-column property="pid" label="图表所属域"/>
                <el-table-column prop="userCanWatch" label="请输入有权查看的用户" style="width: 80%">
                  <template slot-scope="scope">
                    <el-input v-model="scope.row.userCanWatch" size="small" @input="onChangeUserWatch(scope.row)"/>
                  </template>
                </el-table-column>
                <el-table-column
                  type="selection"
                  width="85"
                  label="请选择">
                </el-table-column>
              </el-table>
            </div>
          </el-row>

          <el-row>
            <el-col>
              <el-button type="primary" @click="updateLvpiChecked()" size="small">提交</el-button>
            </el-col>
          </el-row>
      </el-dialog>


    <el-dialog title="修改管理员" :visible.sync="renameAdministratorVisible">
      <el-row>
        <span>{{curDomain.Did}}-{{curDomain.Name}}的管理员</span> <br/>
      </el-row>
      <el-row>
        <el-col :span="16">
          <el-input v-model="onRenameAdministratorName" placeholder="请输入添加的用户名" size="mini"/>
          <span style="margin-top:-20px; color: red">管理员用逗号分隔</span>
        </el-col>
        <el-col :span="2" :offset="1">
          <el-button type="success" size="mini" @click="onRenameAdministrator()">修改</el-button>
        </el-col>
      </el-row>

    </el-dialog>
  </div>
</template>

<script>
const axios = require('axios')
import {getCookie, isEmpty, makeLvpiURL} from '../conf/utils'

export default {
  data() {
    return {
      newDomainAdministrator: '',
      newDomainName: '',
      curDomain: '',
      domainData: [],
      labelPTitleToPidPM: {},      
      renameAdministratorVisible: false,
      onRenameAdministratorName: '',
      urlData: [],
      checkedURL: [],
      curURLGroup: '',
      writableURL: [],
      lvpiOptionVisable: false,
      curDomainId: '',
    }
  },
  methods: {
    onCreateDomain() {
      let self = this
      let params = {
        'name': this.newDomainName,
        'administrator': this.newDomainAdministrator
      }
      axios.post('/auth/lvpiDomain/postFirstDomain', params, {responseType: 'json'})
        .then(function (response) {
          self.$message({
            message: '创建成功！',
            type: 'success'
          })
          self.initGetDomains()
        })
        .catch(function (error) {
          self.$message({
            message: '创建失败！' + error.toString(),
            type: 'error'
          })
        })
    },                                                                                                                

    onChangeUserWatch (row) {
      console.log('输入用户中')
      this.$refs.lvpiOptionTable.toggleRowSelection(row, false)
      this.$refs.lvpiOptionTable.toggleRowSelection(row, true)
    },

    
    renameAdministrator(row) {
      let self = this
      this.curDomain = row
      this.onRenameAdministratorName = row.Administrator
      self.renameAdministratorVisible = true
    },
    onRenameAdministrator() {
      let self = this
      let param = {
        'administrator':this.onRenameAdministratorName,
        'did': this.curDomain.Did,
      }
      axios.post('/auth/lvpiDomain/postAdministrator', param, {responseType: 'json'})
        .then(function (response) {
          self.$message({
            message: '修改成功！',
            type: 'success'
          })
          self.initGetDomains()
          self.renameAdministratorVisible = false
        })
        .catch(function (error) {
          self.$message({
            message: '修改失败！' + error.toString(),
            type: 'error'
          })
        })

    },

    getCheckedURL (lvpiId, lvpiName) {
      let self = this
      self.curURLGroup = lvpiName + '(' + lvpiId + ')'
      self.curDomainId = lvpiId
      self.lvpiOptionVisable = true
      try {
        self.$refs.lvpiOptionTable.clearSelection()
      } catch(error) {
        console.log(error)
      }
      // todo 查看原来选中了哪几个绿皮, 并在空间上选中 (toggleSelection)
      axios.get(makeLvpiURL('/getLvpis/'))
        .then(function (response) {
          ///根据id选择row
          let checkedRow = []
          const pidDict = {}
          response.data.forEach(ele=>{
            pidDict[ele.Did] = ele; 
          })
          
          for(let key in pidDict) {
            const ele = pidDict[key]
            for (let i=0; i<self.lvpiOptions.length; i++){
              if(self.lvpiOptions[i].lid == ele.Lid) {
                self.lvpiOptions[i].userCanWatch = ele.UserCanWatch
                self.lvpiOptions[i].pid = pidDict[ele.Pid].Name
                self.lvpiOptions[i].valid = ele.Valid
                if(ele.Pid == lvpiId) {
                  self.$refs.lvpiOptionTable.toggleRowSelection(self.lvpiOptions[i])
                }
              }
            }
          }
          console.log(self.lvpiOptions)
          //self.$refs.lvpiOptionTable.toggleRowSelection(row)
        }).catch(function (error) {
        console.error(error)
      })
    },
    toggleSelection (checkedRes) {
      this.$nextTick(function () {
        if (checkedRes) {
          checkedRes.forEach(row => {
            this.$refs.lvpiOptionTable.toggleRowSelection(row)
          })
        } else {
          this.$refs.lvpiOptionTable.clearSelection()
        }
      })
    },
    updateLvpiChecked () {
      console.log('已选中lvpi:')
      let self = this
      let params = this.checkedURL
      console.log(this.checkedURL)
      axios.post(makeLvpiURL('/postSecondLvpi'), params)
        .then(function (response) {
          self.lvpiOptionVisable = false
          self.$message('修改成功')
          self.toggleSelection()  // 取消选择
        }).catch(function (error) {
        console.error(error)
      })
    },
    handleLvpiSelectionChange (val) {
      let checkedIds = []
      val.forEach(ele => {
        let item = {
          'lid': ele.lid,
          'pid': this.curDomainId,
          'name': ele.title,
          'index': ele.index,
          'info': ele.info,
          'userCanWatch': ele.userCanWatch,
          'valid': ele.valid,
        }
        checkedIds.push(item)
      })
      this.checkedURL = checkedIds
      // console.log("已选中:"+this.checkedURL)
    },
    initGetDomains () {
      // todo 改成获取域名管理员的域
      let self = this
      axios.get(makeLvpiURL('/getAllDomain'))
        .then(function (response) {
          self.domainData = response.data
          console.log(self.domainData)
        }).catch(function (error) {
        console.log(error)
      })

    },
    ///////////// 绿皮层级
    changeCountPM (row) {
      console.log(row)
      let labelId = row.labelId
      let sortId = row.sortId
      let param = {
        'labelId': labelId + '',
        'sortId': sortId + ''
      }
      let self = this
      axios.post('/auth/label/updateSortId', param, {responseType: 'json'})
        .then(function (response) {
          self.$message({
            message: '排序ID更新成功！',
            type: 'success'
          })
        }).catch(function (error) {
        console.log(error)
      })
    },
    UpdateLabelPM (curLabel) {
      curLabel.status = true
    },
    CancelUpdateLabelPM (curLabel) {
      curLabel.labelTitle2 = curLabel.labelTitle,
        curLabel.labelClass2 = curLabel.labelClass,
        curLabel.labelPid2 = curLabel.labelPid,
        curLabel.labelPTitle2 = curLabel.labelPTitle,
        curLabel.status = false
    },
    choosePidTitlePM (curLabel) {
      curLabel.labelPid2 = this.labelPTitleToPidPM[curLabel.labelPTitle2]
    },
    SaveLabelPM (curLabel) {
      console.log(curLabel)
      let labelId = curLabel.labelId
      let self = this
      if (curLabel.labelTitle2 == curLabel.labelTitle && curLabel.labelClass2 == curLabel.labelClass && curLabel.labelPid2 == curLabel.labelPid) {
        self.$message({
          type: 'success',
          message: '修改前后未变化!'
        })
        return
      }

      let params = {
        'lid': curLabel.labelId.toString(),
        'title': curLabel.labelTitle2,
        'class': curLabel.labelClass2,
        'pid': curLabel.labelPid2.toString(),
      }
      console.log(params)

      axios.post('/auth/labels/pmupdate', params, {responseType: 'json'})
        .then(function (response) {
          curLabel.labelTitle = curLabel.labelTitle2,
            curLabel.labelClass = curLabel.labelClass2,
            curLabel.labelPid = curLabel.labelPid2,
            curLabel.labelPTitle = curLabel.labelPTitle2,
            curLabel.status = false,

            self.$message({
              type: 'success',
              message: '更新成功!'
            })
          console.info(response.data)
        })
        .catch(function (error) {
          console.log(error)
        })
    },

    getAllLvpi: function () {
      let self = this
      axios.get('/redash/redashall')
        .then(function (response) {
          let items = []
          response.data.forEach(ele => {
            items.push({
              'lid': ele.lid,
              'pid': '',
              'title': ele.title,
              'index': ele.index,
              'info': ele.info,
              'userCanWatch': '',
              'valid': '',
            })
          })
          self.lvpiOptions = items
        }).catch(function (error) {
        console.log(error)
      })
    }
  },
  
  created: function () {
    console.log('mounted')
    this.getAllLvpi()
    this.initGetDomains()
  },
  mounted() {
    console.log('mounted')
  }
}
</script>

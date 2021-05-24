<template xmlns="http://www.w3.org/1999/html">
  <div>
    <div v-if="!newDemand">
      <el-button type="primary" size="small" @click="drawer=true">新建</el-button>
    </div>
    <!--新建-->
    <div v-if="newDemand">
      <el-card style="margin-bottom: 2px">
        <div slot="header" class="clearfix">
          <span>基本信息</span>
        </div>
        <el-form ref="basicInfoForm" :model="basicInfoForm" label-width="120px"
                 :rules="basicInfoFormRules">
          <el-form-item label="需求名称" prop="projectName" style="width: 80%;">
            <el-input v-model="basicInfoForm.projectName" size="small"/>
          </el-form-item>

          <el-form-item label="来源部门" prop="apartment" style="width: 45%; display: inline-block">
            <el-select v-model="basicInfoForm.apartment" filterable placeholder="请选择" style="width: 60%" size="small">
              <el-option v-for="name in apartmentOptions" :key="name" :label="name" :value="name"/>
            </el-select>
          </el-form-item>

          <el-form-item label="优先级" prop="priority" style="width: 45%; display: inline-block">
            <el-radio-group v-model="basicInfoForm.priority" size="small" style="width: 100%">
              <el-radio-button label="P0"/>
              <el-radio-button label="P1"/>
              <el-radio-button label="P2"/>
              <el-radio-button label="P3"/>
            </el-radio-group>
          </el-form-item>

          <el-form-item label="处理人" prop="developer" style="width: 45%;display: inline-block" size="small">
            <el-input v-model="basicInfoForm.developer" style="width: 60%"/>
          </el-form-item>

          <el-form-item label="需求阶段" prop="status" style="width: 45%;display: inline-block" size="small">
            <el-select v-model="basicInfoForm.status" filterable placeholder="请选择" style="width: 60%">
              <el-option
                v-for="(textValue,numKey) in devStatusMapping"
                :key="textValue"
                :label="textValue"
                :value="numKey">
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="期望时间" prop="expectTime" style="width: 45%; display: inline-block">
            <el-date-picker
              v-model="basicInfoForm.expectTime"
              type="date"
              value-format="yyyy-MM-dd"
              style="width: 60%"
              placeholder="选择日期" size="small"/>
          </el-form-item>

          <el-form-item label="预计开始时间" prop="startTime" style="width: 45%; display: inline-block" size="small">
            <el-date-picker
              v-model="basicInfoForm.startTime"
              type="date"
              value-format="yyyy-MM-dd"
              style="width: 60%"
              placeholder="选择日期" size="small"/>
          </el-form-item>

          <el-form-item label="预计结束时间" prop="endTime" style="width: 45%; display: inline-block" size="small">
            <el-date-picker
              v-model="basicInfoForm.endTime"
              type="date"
              style="width: 60%"
              value-format="yyyy-MM-dd"
              placeholder="选择日期" size="small"/>
          </el-form-item>
          <!--<el-form-item label="" prop="status" style="width: 50%;display: inline-block" size="small">
            <el-row type="flex" justify="end">
              <el-button type="primary" size="small" @click="uploadDataRequest('saveOnly')">保存</el-button>
            </el-row>
          </el-form-item>-->
        </el-form>

      </el-card>

      <el-card>
        <!--        <div slot="header" class="clearfix">-->
        <!--          <span>需求信息</span>-->
        <!--        </div>-->
        <el-tabs v-model="activeTab" type="card">
          <el-tab-pane label="需求文档" name="first">
            <!--标题在宽度非常窄的时候, 会形成双层, 导致高度增加-->
            <div style="min-height: 460px; max-height: 598px; overflow-y: auto">
              <quill-editor
                id="editor"
                :disabled="isdisabled"
                ref="myQuillEditor"
                v-model="basicInfoForm.document"
                :options="editorOption"/>
            </div>

            <div style="position:absolute; top:70px; right:10px; z-index:1000;">
              <el-button @click="isdisabled=false" type="primary" plain v-show="canEditDoc">进入编辑</el-button>
            </div>
          </el-tab-pane>
          <el-tab-pane label="指标说明" name="second">
            <div style="margin-left: 10px; min-height: 400px">

              <el-button type="primary" size="small" @click="onImportDataset">导入指标</el-button>
              <el-button type="success" size="small" @click="addIndicatorVisible=true">添加指标</el-button>
              <el-dialog title="添加指标" :visible.sync="addIndicatorVisible">
                <el-form ref="userDefineIndicatorForm" :model="userDefineIndicatorForm" label-width="80px">
                  <el-form-item label="资产名称">
                    <el-input v-model="userDefineIndicatorForm.name" size="small" style="width: 60%"/>
                  </el-form-item>

                  <el-form-item label="资产描述" style="width: 100%">
                    <el-input v-model="userDefineIndicatorForm.description" size="mini"
                              type="textarea" :rows="5"/>
                  </el-form-item>

                  <el-form-item>
                    <el-button type="success" size="small" @click="onAddUdIndicator" style="float: right">添加描述
                    </el-button>
                  </el-form-item>
                </el-form>
              </el-dialog>
              <!--这个是不显示的-->
              <!--              <el-table-->
              <!--                :data="multipleSelection"-->
              <!--                stripe-->
              <!--                style="width:99%; margin-bottom: 10px; display:none"-->
              <!--                :show-header="true">-->
              <!--                <el-table-column key="type_name" prop="type_name" label="资产类型"/>-->
              <!--                <el-table-column key="item_name" prop="item_name" label="资产名"/>-->
              <!--                <el-table-column key="description" prop="description" label="描述"/>-->
              <!--              </el-table>-->
              <!--只显示合并后的指标-->
              <el-table
                :data="unionIndicator"
                stripe
                style="width:99%; margin-bottom: 10px;"
                :show-header="true">
                <el-table-column key="type_name" prop="type_name" label="资产类型"/>
                <el-table-column key="item_name" prop="item_name" label="资产名"/>
                <el-table-column key="description" prop="description" label="描述"/>
              </el-table>
            </div>
          </el-tab-pane>
          <el-tab-pane label="需求附件" name="third">
            <div style="margin-left: 10px; min-height: 400px">
              <el-col :span="12">
                <el-upload
                  class="upload-demo"
                  ref="upload"
                  action="dummyUrl"
                  :limit="1"
                  :file-list="fileList"
                  :auto-upload="false"
                  :before-upload="beforeUpload"
                  :on-exceed="onExceed"
                >
                  <el-button slot="trigger" size="small" type="primary">选择需求文档</el-button>
                  <el-button style="margin-left: 10px;" size="small" type="success" @click="submitFile">上传</el-button>
                  <span style="color:#BDBDBD;font-size:12px;">{{ uploadTableName }}</span>
                </el-upload>

                <div style="height: 360px ; overflow:scroll; ">
                  <span>评论区:</span>
                  <el-button size="mini" type="primary" plain @click="commentVisible=true">添加评论</el-button>
                  <el-timeline style="margin-left: -30px;margin-top: 5px">
                    <el-timeline-item v-for="item in criticalItems" :timestamp="item.Time+' ('+item.Commentator+')'"
                                      placement="top">
                      <el-card>
                        <!--                        <h4>评论人:{{item.Commentator}}</h4>-->
                        <p>{{ item.Comment }}</p>
                      </el-card>
                    </el-timeline-item>
                  </el-timeline>
                  <!--                  <div v-for="item in criticalItems" style="margin-bottom: 20px; border: #00bb00">
                                      <span style="font-size: 14px">评论人: {{ item.Commentator }}</span> <br/>
                                      <span style="font-size: 14px">评论内容: {{ item.Comment }}</span>
                                    </div>-->
                  <!--                  </el-scrollbar>-->
                </div>
              </el-col>

              <el-col :span="12">
                <el-table
                  :data="documentList"
                  stripe
                  style="width:99%; margin-bottom: 10px"
                  :show-header="false">
                  <el-table-column key="filename" prop="filename" label="文件名称">
                    <template slot-scope="scope">
                      <sapn @click="dowloadAttachment(scope.row)" style="cursor: pointer">{{ scope.row.Filename }}
                      </sapn>
                    </template>
                  </el-table-column>
                  <el-table-column key="operate" prop="operate" label="操作" width="100px">
                    <template slot-scope="scope">
                      <el-button @click="delAttachment(scope.row)" type="text" size="small" icon="el-icon-close"/>
                    </template>
                  </el-table-column>
                </el-table>
              </el-col>

            </div>
          </el-tab-pane>
          <el-tab-pane label="历史记录" name="fourth">
            <el-table :data="historyData" stripe style="width: 100%" ref="historyTable"
                      highlight-current-row>
              <el-table-column v-for="historyCol in historyCols" v-if="!historyCol.hidden"
                               :key="historyCol.prop" :prop="historyCol.prop" :label="historyCol.label"
                               :sortable="historyCol.sortable"/>
              <el-table-column fixed="right" label="操作" width="110">
                <template slot-scope="scope">
                  <el-button @click="checkHistory(scope.row)" type="text" size="mini">查看数据</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
        </el-tabs>

      </el-card>

      <el-row type="flex" justify="end" style="margin-top: 10px">
        <el-button type="primary" size="small" @click="uploadDataRequest">提交</el-button>
        <el-button size="small" type="danger" @click="$router.push({name: 'demandList'})">取消</el-button>
      </el-row>
    </div>

    <!-- 选择依赖的数据资产 -->
    <el-dialog title="选择数据资产" :visible.sync="dialogTableVisible" width="80%">
      <span>请在第一列勾选</span>
      <div style="max-height: 500px ; overflow:scroll; ">
        <el-scrollbar style="height: 100%;">
          <el-table :data="assetList" stripe style="width: 100%" :border="false" ref="assetTable"
                    @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55"/>
            <el-table-column key="type_name" prop="type_name" label="资产类型" min-width="18%"/>
            <el-table-column key="item_name" prop="item_name" label="资产名" min-width="18%">
              <template slot-scope="scope">
                <span>{{ scope.row.item_name }}</span>
              </template>
            </el-table-column>
            <el-table-column key="description" prop="description" label="描述" min-width="46%"/>
            <el-table-column key="business_label" prop="business_label" label="业务标签" min-width="18%">
              <template slot-scope="scope">
                <el-tag style="margin-right: 5px" v-for="(businessLabel, index) in scope.row.business_label.split(',')"
                        v-show="scope.row.business_label.split(',').length>0 && scope.row.business_label.split(',')[0]!==''">
                  {{ serviceTagMap[businessLabel] }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-scrollbar>
      </div>
      <el-row type="flex" justify="end" style="margin-top: 10px; margin-bottom: -5px">
        <el-button type="primary" size="small" @click="dialogTableVisible=false">确定</el-button>
      </el-row>
    </el-dialog>


    <el-dialog title="上传评论" :visible.sync="commentVisible" width="600px">
      <span>请在下面书写评论</span>
      <el-input v-model="commentStr" size="mini" type="textarea" :rows="10"/>

      <el-row type="flex" justify="end" style="margin-top: 10px; margin-bottom: -5px">
        <el-button type="primary" size="small" @click="uploadComment">提交</el-button>
      </el-row>
    </el-dialog>
  </div>
</template>

<script>
import 'mavon-editor/dist/css/index.css'
import router from "../router"

const axios = require('axios')

export default {
  name: 'DemandRequirement',
  data() {
    return {
      isdisabled: true,
      commentStr: "",
      commentVisible: false,
      editorOption: {
        placeholder: "此处编辑需求说明",
        theme: "snow",
      },
      drawer: true,
      criticalItems: [{"userName": "zhangsan", "content": "嘿嘿嘿"}, {
        "userName": "lisi",
        "content": "后水水水水是是是"
      }, {"userName": "lisi", "content": "后水水水水是是是"}, {"userName": "lisi", "content": "后水水水水是是是"}, {
        "userName": "lisi",
        "content": "后水水水水是是是"
      },],
      activeTab: 'first',
      dialogTableVisible: false,
      activeStep: 2,
      newDemand: 'rtl',
      serviceTagList: [],
      serviceTagMap: {},
      apartmentOptions: ["营销事业部/商业产品-搜索", "营销事业部/商业产品-投放平台", "营销事业部/商业产品-产品运营",
        "营销事业部/商业产品-品牌-银河", "营销事业部/商业产品-网盟-信息流", "营销事业部/商业产品-创新", "营销事业部/商业产品-阅读",
        "营销事业部/商业产品-反作弊", "营销事业部/商业产品-数据", "营销事业部/战略合作部", "营销事业部/商业分析中心",
        "营销事业部/广告营销中心", "营销事业部/风控", "营销事业部/ADM", "搜索事业部", "输入法事业部", "商业平台事业部"],
      assetList: [],
      demandId: "",
      userDefineIndicatorForm: {
        "name": "",
        "description": "",
      },
      userDefineIndicatorSet: [],
      basicInfoForm: {
        'projectName': '',
        'expectTime': '',
        'startTime': '',
        'endTime': '',
        'priority': '',
        'developer': '',
        'document': "<h2><span style=\"color: rgb(230, 0, 0);\">*</span>1.背景（或项目链接）。</h2><p><span style=\"color: rgb(230, 0, 0);\"> 必填</span></p><p><br></p><h2><span style=\"color: rgb(230, 0, 0);\">*</span>2.价值。</h2><p><span style=\"color: rgb(230, 0, 0);\"> 必填</span></p><p><br></p><h2><span style=\"color: rgb(230, 0, 0);\">*</span>3.需求详情（指标、维度、口径等）。&nbsp;</h2><p><span style=\"color: rgb(230, 0, 0);\"> 必填</span></p><p><br></p><h2><span style=\"color: rgb(230, 0, 0);\">*</span>4.实现形式（如接口、报表、固定SQL、工单、数据表等）。&nbsp;</h2><p><span style=\"color: rgb(230, 0, 0);\"> 必填</span></p><p><br></p><h2>5.需求PRD（链接、附件等）。&nbsp;</h2><p>选填</p><p><br></p><h2>6.历史数据回溯时间。&nbsp;</h2><p>选填</p><p><br></p><h2>7.埋点文档。</h2><p>选填</p>",
        'status': "1",
        'apartment': "",
      },
      fileList: [],    // 需求文件
      basicInfoFormRules: {
        'projectName': [{required: true, message: '请输入需求名', trigger: 'blur'}],
        'expectTime': [{required: true, message: '请输入期望时间', trigger: 'blur'}],
        'priority': [{required: true, message: '请选择优先级', trigger: 'blur'}],
        'apartment': [{required: true, message: '请选来源部门', trigger: 'blur'}],
      },
      developer: '',  // 开发负责人
      documentList: [],
      multipleSelection: [],   // 选择的表格项
      devStatusMapping: {},
      addIndicatorVisible: false,
      historyData: [],
      historyCols: [],
    }
  },
  methods: {
    checkHistory(row) {
      // this.$alert(row)
      this.basicInfoForm.projectName = row.project_name
      this.basicInfoForm.expectTime = row.expect_time
      this.basicInfoForm.startTime = row.start_time
      this.basicInfoForm.endTime = row.end_time
      this.basicInfoForm.priority = row.priority
      this.basicInfoForm.developer = row.developer
      this.basicInfoForm.document = row.demand_text
      this.basicInfoForm.status = row.status+""
      this.basicInfoForm.apartment = row.apartment
      this.$refs.historyTable.setCurrentRow(row)
      this.$message({
        type: 'success',
        message: '显示历史'+ row.hid
      })
    },
    getComments() {
      let self = this
      axios.get('/requirement/indicator/getDataRequirement?demandId=' + self.demandId)
        .then(function (response) {
          self.criticalItems = response.data
        }).catch(function (error) {
        console.log(error)
      })
    },
    uploadComment() {
      let self = this
      if (this.demandId === "" || this.demandId === 0) {
        this.$alert("请先提交基本信息")
        return
      }
      let param = {
        demandId: parseInt(this.demandId),
        comment: this.commentStr,
      }
      axios.post('/requirement/indicator/addDataRequirement', param, {responseType: 'json'})
        .then((response) => {
          if (response.data.Res) {
            self.$message({
              message: '评论成功',
              type: 'success'
            })
            this.commentVisible = false
            self.getComments()
          } else {
            self.$message({
              message: '评论失败:' + response.data.Info,
              type: 'xxx',
            })
          }
        })
        .catch((error) => {
          console.log(error)
        })
    },
    onDeleteUdI(row) {
      let index = this.userDefineIndicatorSet.indexOf(row)
      if (index > -1) {
        this.userDefineIndicatorSet.splice(index, 1)
      }
    },
    onAddUdIndicator() {
      this.userDefineIndicatorSet.push({
        "item_name": this.userDefineIndicatorForm.name,
        "description": this.userDefineIndicatorForm.description,
      })
      this.userDefineIndicatorForm = {
        "name": "",
        "description": "",
      }
      this.addIndicatorVisible = false
    },
    axiosDownload(data, fileName) {
      if (!data) {
        return
      }
      let url = window.URL.createObjectURL(new Blob([data]))
      let link = document.createElement('a')
      link.style.display = 'none'
      link.href = url
      link.setAttribute('download', fileName)

      document.body.appendChild(link)
      link.click()
    },
    delAttachment(row) {
      let aid = row.Aid
      let self = this
      this.$confirm('确定删除附件吗', '提示', {
        distinguishCancelAndClose: true,
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
        center: true
      }).then(() => {
        axios.get('/requirement/attachment/delete/' + aid)
          .then(function (response) {
            self.getUploadFileList()
          }).catch(function (error) {
          console.log(error)
        })
      })
    },
    dowloadAttachment(row) {
      let aid = row.Aid
      let filename = row.Filename
      let self = this

      axios({
        method: 'post',
        url: '/requirement//attachment/download/' + aid,
        responseType: 'blob'
      }).then(response => {
        console.log("下载请求完毕")
        self.axiosDownload(response.data, filename)
      }).catch((error) => {
        console.log(error)
      })
    },
    getAllStatusMapping() {
      let self = this
      axios.get('/requirement/GetDevStatusMapping')
        .then(function (response) {
          self.devStatusMapping = response.data
          // self.$alert(self.devStatusMapping)
        }).catch(function (error) {
        console.log(error)
      })
    },
    getDevStatus(status) {
      return this.devStatusMapping[status]
    },
    onImportDataset() {
      this.dialogTableVisible = true
      let self = this
      // 触发反选表格, 延时200ms, 等待 dom 树建立完毕
      setTimeout(function () {
        self.multipleSelection.forEach(row => {
          self.$refs.assetTable.toggleRowSelection(row, true)
        })
      }, 200)
    },
    handleClose() {
      this.drawer = false
    },
    saveDoc(value, render) {
      let docForm = {
        'content': value,
        'html': render,
      }
      console.log(docForm)
    },
    onExceed(files, fileList) {
      this.$set(fileList[0], 'raw', files[0])
      this.$set(fileList[0], 'name', files[0].name)
      this.$refs['upload'].clearFiles() //清除文件
      this.$refs['upload'].handleStart(files[0])//选择文件后的赋值方法
    },
    submitFile() {
      let self = this
      let param = {
        "projectName": this.basicInfoForm.projectName,
        "expectTime": this.basicInfoForm.expectTime,
        "proposer": "",
        "developer": this.basicInfoForm.developer,
        "developTime": "",
        "status": parseInt(this.basicInfoForm.status, 10),
        "demandText": this.basicInfoForm.document,
        "priority": this.basicInfoForm.priority
      }
      this.$confirm('确定上传文件吗', '提示', {
        distinguishCancelAndClose: true,
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
        center: true
      }).then(() => {
        // 新增需求时, 上传文件先提交基本信息
        if (this.arbitrateUpdateOrInsert() === 'new') {
          // 执行插入操作
          axios.post('/requirement/post', param, {responseType: 'json'})
            .then((response) => {
              if (response.data.Res) {
                self.demandId = response.data.Info
                self.postUserDefineIndicator(self.demandId)
                self.uploadDemandAssetRelation()
                self.$notify({
                  message: "提交成功, 需求id:" + response.data.Info,
                  type: 'success',
                  duration: 2000
                })
                self.demandId = response.data.Info
                self.$refs.upload.submit()   // 会触发 beforeUpload 方法
              } else {
                self.$notify({
                  message: "提交失败: " + response.data.Info,
                  type: 'error',
                  duration: 2000
                })
              }
            })
            .catch((error) => {
              console.log(error)
            })
        } else {
          self.$refs.upload.submit()   // 会触发 beforeUpload 方法
        }
      })

    },
    getUploadFileList() {
      let rid = this.demandId  // 需求Id
      let self = this
      axios.get('/requirement/attachment/get/' + rid)
        .then(function (response) {
          self.documentList = response.data
        }).catch(function (error) {
        console.log(error)
      })
    },
    beforeUpload(file) {
      let self = this
      let formData = new FormData()
      formData.append('file', file)

      axios.post('/requirement/attachment/upload/' + this.demandId, formData, {
          headers: {'Content-Type': 'multipart/form-data'}
        }
      ).then(function (response) {
        let data = response.data
        if (data.Res) {
          self.$notify({
            message: "附件上传成功",
            type: 'success',
            duration: 2000
          })
          self.getUploadFileList()
        } else {
          self.$notify({
            title: "附件上传失败",
            message: data.Info,
            type: 'error',
            duration: 3000
          })
        }
      }).catch(function (error) {
        console.log(error)
      })
      return false   // 阻断action中的传输, 用于自己实现axios传输文件, 返回true会访问action中的url
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
    getAllDictionary(func) {
      let self = this
      axios.get('/indicator/dictionary/getAll')
        .then(function (response) {
          self.assetList = response.data.data
          if ((func !== null) && (func !== undefined) && (func !== "")) {
            func()  // 执行回调函数
          }
        }).catch(function (error) {
        console.log(error)
      })
    },
    handleSelectionChange(val) {
      this.multipleSelection = val
    },
    uploadDemandAssetRelation() {
      let self = this
      let assets = []
      try {
        this.multipleSelection.forEach(ele => {
          let originId = ele.origin_id + ""
          let typeName = ele.type_name
          assets.push({
            "origin_id": originId,
            "type_name": typeName,
          })
        })
        axios.post('/requirement/postRequireAssetRelation/' + this.demandId, assets, {responseType: 'json'})
          .then((response) => {
            if (response.data.Res) {
              self.$notify({
                message: "数据资产绑定成功",
                type: 'success',
                duration: 2000
              })
            } else {
              self.$notify({
                title: "数据资产关系绑定失败",
                message: response.data.Info,
                type: 'error',
                duration: 2000
              })
            }
          })
          .catch((error) => {
            self.$message({
              message: '数据资产关系绑定失败: ' + error.toString(),
              type: 'error'
            })
          })
      } catch (e) {
        console.log(e)
      }
    },
    arbitrateUpdateOrInsert() {
      if (this.demandId === "" ||
        this.demandId === "0" ||
        this.demandId === 0 ||
        this.demandId === undefined ||
        this.demandId === null) {
        return "new"
      } else {
        return "update"
      }
    },
    getHistory() {
      let self = this
      axios.get('/requirement/getHistory/' + this.demandId)
        .then(function (response) {
          self.historyData = response.data.data
          self.historyCols = response.data.cols
        })
        .catch(function (error) {
          console.log(error)
        })
    },
    uploadDataRequest() {
      let flag = this.arbitrateUpdateOrInsert()
      let self = this
      let param = {
        "projectName": this.basicInfoForm.projectName,
        "expectTime": this.basicInfoForm.expectTime,
        "proposer": "",
        "developer": this.basicInfoForm.developer,
        "developTime": "",
        "status": parseInt(this.basicInfoForm.status, 10),
        "demandText": this.basicInfoForm.document,
        "priority": this.basicInfoForm.priority,
        "startTime": this.basicInfoForm.startTime,
        "endTime": this.basicInfoForm.endTime,
        "apartment": this.basicInfoForm.apartment,

      }
      if (flag === 'update') {
        // 执行更新操作
        self.postUserDefineIndicator(this.demandId)
        axios.post('/requirement/update/' + this.demandId, param, {responseType: 'json'})
          .then((response) => {
            if (response.data.Res) {
              self.uploadDemandAssetRelation()
              self.$notify({
                message: "更新成功",
                type: 'success',
                duration: 2000
              })
              self.isdisabled = true
              router.push({path: 'demandList'})  // 更新后跳转
            } else {
              self.$notify({
                message: "提交失败: " + response.data.Info,
                type: 'error',
                duration: 2000
              })
            }
          })
          .catch((error) => {
            console.log(error)
          })
      } else {
        // 执行插入操作
        axios.post('/requirement/post', param, {responseType: 'json'})
          .then((response) => {
            if (response.data.Res) {
              self.demandId = response.data.Info
              self.postUserDefineIndicator(self.demandId)
              self.uploadDemandAssetRelation()
              self.$notify({
                message: "提交成功, 需求id:" + response.data.Info,
                type: 'success',
                duration: 2000
              })
              self.demandId = response.data.Info
              router.push({path: 'demandList'})
            } else {
              self.$notify({
                message: "提交失败: " + response.data.Info,
                type: 'error',
                duration: 2000
              })
            }
          })
          .catch((error) => {
            console.log(error)
          })
      }
      this.getHistory()
    },
    postUserDefineIndicator(rid) {
      let self = this
      let param = []
      this.userDefineIndicatorSet.forEach(ele => {
        param.push({
          "Rid": parseInt(rid + ""),
          "AssetName": ele.item_name,
          "Description": ele.description
        })
      })
      axios.post('/requirement/indicator/userDefine/post/' + rid, param, {responseType: 'json'})
        .then((response) => {
        })
        .catch((error) => {
          self.$message({
            message: '提交失败！' + error.toString(),
            type: 'error'
          })
        })

    },
    getDemandById() {
      let self = this
      axios.get('/requirement/getDemandById/' + this.demandId)
        .then(function (response) {
          let data = response.data
          // self.$alert(data);
          self.basicInfoForm.projectName = data.ProjectName
          self.basicInfoForm.expectTime = data.ExpectTime
          self.basicInfoForm.priority = data.Priority
          self.basicInfoForm.developer = data.Developer
          self.basicInfoForm.document = data.DemandText
          self.basicInfoForm.status = data.Status + ""
          self.basicInfoForm.startTime = data.StartTime
          self.basicInfoForm.endTime = data.EndTime
          self.basicInfoForm.apartment = data.Apartment
        }).catch(function (error) {
        console.log(error)
      })
      this.getUploadFileList()
    },
    getUserDefineIndicators() {
      let self = this
      axios.get('/requirement/indicator/userDefine/get/' + this.demandId)
        .then(function (response) {
          let data = response.data
          data.forEach(ele => {
            self.userDefineIndicatorSet.push({
              "item_name": ele.AssetName,
              "description": ele.Description,
            })
          })
        }).catch(function (error) {
        console.log(error)
      })
    },
    getDemandRelationById() {
      let self = this
      axios.get('/requirement/getAssetRelationById/' + this.demandId)
        .then(function (response) {
          let relations = response.data
          let eleInAssetList = []
          relations.forEach(r => {
            let target = self.findAssetInList(r)
            eleInAssetList.push(target)
          })
          self.multipleSelection = eleInAssetList

        }).catch(function (error) {
        console.log(error)
      })
    },
    findAssetInList(asset) {
      let target = null
      let targetType = ''
      if (asset.Atype === 'predefine') {
        targetType = '原子指标'
      } else if (asset.Atype === 'combine') {
        targetType = '计算指标'
      } else if (asset.Atype === 'dimension') {
        targetType = '维度'
      }
      this.assetList.forEach(ele => {
        if ((ele.origin_id === asset.Assetid) && (ele.type_name === targetType)) {
          target = ele
        }
      })
      return target
    }

  },
  computed: {
    canEditDoc: function () {
      let status = this.devStatusMapping[this.basicInfoForm.status]
      let flag = (status.indexOf("开发") === -1) && (status.indexOf("完成") === -1) && (status.indexOf("验收") === -1)
      // console.log("触发计算:" + flag)
      // console.log("status:" + status)
      return flag
    },
    // 合并导入的指标和添加后的指标
    unionIndicator: function () {
      try {
        let importIndicators = this.multipleSelection
        let addIndicators = this.userDefineIndicatorSet
        let unionSet = []
        console.log(importIndicators)
        console.log(addIndicators)
        // 合并添加的指标
        if (addIndicators != null) {
          addIndicators.forEach(ele => {
            unionSet.push({
              "description": ele.description,
              "item_name": ele.item_name,
              "type_name": '手动添加'
            })
          })
        }
        // 合并导入的指标
        if (importIndicators != null) {
          importIndicators.forEach(ele => {
            unionSet.push({
              "description": ele.description == null ? "" : ele.description,
              "item_name": ele.item_name,
              "type_name": ele.type_name,
            })
          })
        }
        console.log("合并指标完成")
        return unionSet
      } catch (err) {
        return []
      }
    }
  },
  mounted() {
    this.getAllStatusMapping()
    this.demandId = this.$route.params.demandId //sessionStorage.getItem("demandId");
    console.log("获取的需求Id:" + this.demandId)
    if ((this.demandId != null) && (this.demandId !== "") && (this.demandId !== undefined) && (this.demandId !== 'new')) {
      this.getDemandById()
      this.getUserDefineIndicators()
    }
    // 因为要反选表格, 所以在获得所有资产列表后再获取需求和资产的绑定关系
    if (this.demandId === 'new') {
      this.demandId = ""
    }
    this.getAllServiceTag(this.getAllDictionary(this.getDemandRelationById))
    this.getComments()
    this.getHistory()
  }
}
</script>

<style scoped>
.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}

.clearfix:after {
  clear: both
}

.el-card /deep/ .el-card__header {
  padding: 0 0 20px 20px;
  font-size: 14px;
}

#editor {
  margin: auto;
  width: 100%;
  height: 400px;
}

.el-table >>> .el-table__body tr.current-row > td {
  background-color: #E1F3D8 !important;
}
</style>

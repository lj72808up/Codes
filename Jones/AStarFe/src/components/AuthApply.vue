<template>
  <div style="margin-left: 10px">
    <h3 v-show="!isFinish">权限申请</h3>
    <el-form ref="applyForm" :model="applyForm" label-width="90px" v-show="!isFinish">
      <el-row>
        <el-col :span="totalWidth">
          <el-form-item label="姓名">
            <el-input v-model="applyForm.name"/>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row v-show="false">
        <el-col :span="totalWidth">
          <el-form-item label="分机">
            <el-input v-model="applyForm.phone"/>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="totalWidth">
          <el-form-item label="直接主管">
            <el-input v-model="applyForm.directManager"/>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row v-show="false">
        <el-col :span="totalWidth">
          <el-form-item label="Sogou邮箱">
            <el-input v-model="applyForm.email"/>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="2 * totalWidth / 3">
          <el-form-item label="部门">
            <el-input v-model="applyForm.department"/>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="2 * totalWidth / 3">
          <el-form-item label="职位">
            <el-input v-model="applyForm.position"/>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="totalWidth">
          <el-form-item label="申请原因">
            <el-input type="textarea" v-model="applyForm.demand"/>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-button @click="onApply()" type="primary">申请</el-button>
      </el-row>
    </el-form>

    <h3 v-show="isFinish">{{msg}}</h3>
  </div>
</template>

<script>
  const axios = require('axios')

  export default {
    data () {
      return {
        applyForm: {
          name: '',
          phone: '',
          directManager: '',
          email: '',
          department: '',
          position: '',
          demand: '',
        },
        totalWidth: 24,
        uid: '',
        isFinish: false,

        applyProcessing: false,  // 申请还在进行中
        applySuccess: false, // 申请已经审核通过
        msg: ""
      }
    },
    methods: {
      _isMobile () {
        let flag = navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i)
        return flag
      },
      onApply () {
        let self = this
        let params = {
          boss: this.applyForm.directManager,
          apartment: this.applyForm.department,
          position: this.applyForm.position,
          name: this.applyForm.name,
          reason: this.applyForm.demand,
          uid: this.uid
        }
        this.$confirm('是否提交申请?', '提示', {
          distinguishCancelAndClose: true,
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
          center: true
        }).then(() => {
          let url = axios.defaults.baseURL.replace(/jones/g, '') + 'xiaop/apply'
          axios({
            method: 'post',
            url: url,
            data: params,
            responseType: 'string'
          }).then(function (response) {
            self.$message({
              type: 'success',
              message: '提交成功.'
            })
            self.applyForm = {
              name: '',
              phone: '',
              directManager: '',
              email: '',
              department: '',
              position: '',
              demand: '',
            }
            self.isFinish = true
          }).catch(function (error) {
            self.$message({
              type: 'error',
              message: error
            })
            console.log(error)
          })
        })
      }
    },
    mounted () {
      let uid = this.$route.query.uid
      this.uid = uid
      console.log('uid:' + uid)
      if (this._isMobile()) {
        this.totalWidth = 24
      } else {
        this.totalWidth = 9
      }
      // 查看该uid的审批状态
      let url = axios.defaults.baseURL.replace(/jones/g, '') + 'xiaop/applyStatus/'+uid
      let self = this
      axios.get(url)
        .then(function (response) {
          let res = response.data
          if (res.uid === ""){// 还未申请

          }else {  // 已经申请
            if (res.status === 0) {
              // 还在审批中
              self.applyProcessing = true
              self.isFinish = true
              self.msg = "等待审核选中.."
            } else {
              // 审核通过
              self.applySuccess = true
              self.isFinish = true
              self.msg = "审核完成, "+res.detail
            }
          }
        })
        .catch(function (error) {
          console.log(error)
        })
    }
  }
</script>

<style scoped>

</style>

<template>
  <div style="margin-left: 10px">
    <h3 v-show="!isFinish">审批</h3>
    <el-form v-show="!isFinish" ref="applyForm" :model="applyForm" label-width="90px">
      <el-row>
        <el-col :span="totalWidth">
          <el-form-item label="姓名" prop="name">
            <el-input v-model="applyForm.name" :readonly="true"/>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row v-show="false">
        <el-col :span="totalWidth">
          <el-form-item label="分机" prop="phone">
            <el-input v-model="applyForm.phone" :readonly="true"/>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="totalWidth">
          <el-form-item label="直接主管" prop="directManager">
            <el-input v-model="applyForm.directManager" :readonly="true"/>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row v-show="false">
        <el-col :span="totalWidth">
          <el-form-item label="Sogou邮箱" prop="email">
            <el-input v-model="applyForm.email" :readonly="true"/>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="2 * totalWidth / 3">
          <el-form-item label="部门" prop="department">
            <el-input v-model="applyForm.department" :readonly="true"/>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="2 * totalWidth / 3" prop="position">
          <el-form-item label="职位">
            <el-input v-model="applyForm.position" :readonly="true"/>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col :span="totalWidth" prop="demand">
          <el-form-item label="申请原因">
            <el-input type="textarea" v-model="applyForm.demand" :readonly="true"/>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col>
          <el-form-item label="分配角色" prop="applyRole">
            <el-select v-model="applyForm.applyRole">
              <el-option
                v-for="item in roleOptions"
                :key="item"
                :label="item"
                :value="item">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row>
        <el-col>
          <el-button @click="onSubmitApproval" type="primary">确定提交</el-button>
        </el-col>
      </el-row>

    </el-form>

    <h4 v-show="isFinish">
      {{msg}}
    </h4>
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
          applyRole: '',
        },
        totalWidth: 24,
        uid: '',
        roleOptions: [],
        isFinish: false,
        msg: '',
      }
    },
    methods: {
      _isMobile () {
        return navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i)
      },
      onSubmitApproval () {
        let self = this
        let params = {
          uid: this.uid,
          role: this.applyForm.applyRole
        }
        console.log(params)
        let url = axios.defaults.baseURL.replace(/jones/g, '') + 'xiaop/role/add'
        axios.post(url, params, {responseType: 'json'})
          .then(function (response) {
            location.reload()
            /*let flag = 'success'
            if (!response.data.Res) {
              flag = 'error'
            }
            self.$message({
              type: flag,
              message: '申请用户:' + self.uid + ', ' + response.data.Info
            })*/
          }).catch(function (error) {
          console.error(error)
        })
      }
    },
    mounted () {
      // 判断手机和电脑
      if (this._isMobile()) {
        this.totalWidth = 24
      } else {
        this.totalWidth = 9
      }

      this.uid = this.$route.query.uid
      this.applyForm.name = this.$route.query.name
      this.applyForm.directManager = this.$route.query.boss
      this.applyForm.department = this.$route.query.apartment
      this.applyForm.position = this.$route.query.position
      this.applyForm.demand = this.$route.query.reason

      let self = this
      let url = axios.defaults.baseURL.replace(/jones/g, '') + 'xiaop/applyStatus/' + this.uid
      axios.get(url)
        .then(function (response) {
          let res = response.data
          console.log(res)
          if (res.status !== 0) {
            // 审核完成
            self.isFinish = true
            self.msg = res.detail
          }
        })
        .catch(function (error) {
          console.log(error)
        })

      let rolesUrl = axios.defaults.baseURL.replace(/jones/g, '') + 'xiaop/roles'
      // 1. 查看有哪些角色
      axios.get(rolesUrl).then(function (response) {
        console.log(response)
        self.roleOptions = response.data
      }).catch(function (error) {
        console.log(error)
      })

    }
  }
</script>

<style scoped>

</style>

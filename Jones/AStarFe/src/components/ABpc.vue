<template>
  <div>
    <h3>PC小流量监测工具申请</h3>
    <el-form ref="ABform" :model="ABform" :rules="ABrules" label-width="80px">
         <h4>项目配置</h4>
         <el-row>
           <el-col :span="18">
             <el-form-item label="项目名称" prop="name">
               <el-input v-model="ABform.name" clearable></el-input>
             </el-form-item>
           </el-col>
         </el-row>
         <!--<el-row>-->
           <!--<el-form-item label="起止日期">-->
             <!--<el-date-picker-->
               <!--v-model="ABform.dateRange"-->
               <!--type="daterange"-->
               <!--align="right"-->
               <!--unlink-panels-->
               <!--range-separator="至"-->
               <!--start-placeholder="开始日期"-->
               <!--end-placeholder="结束日期"-->
               <!--value-format="yyyyMMdd"-->
               <!--:picker-options="dateOptions">-->
             <!--</el-date-picker>-->
           <!--</el-form-item>-->
         <!--</el-row>-->
         <el-row>
             <el-form-item label="全局位置">
               <el-form-item label="" prop="ad_pos_restrict_up">
                 <el-col :span="3">
                   <el-checkbox :indeterminate="isIndeterminatePosUp" v-model="checkAllPosUp" @change="handleCheckAllChangePosUp">全选上方</el-checkbox>
                 </el-col>
                 <el-col :span="20">
                   <el-checkbox-group v-model="ABform.ad_pos_restrict_up" @change="handleCheckedChangePosUp">
                     <el-checkbox v-for="pos in ad_pos_up" :label="pos" :key="pos">{{pos}}</el-checkbox>
                   </el-checkbox-group>
                 </el-col>
               </el-form-item>

               <el-form-item label="" prop="ad_pos_restrict_down">
                 <el-col :span="3">
                   <el-checkbox :indeterminate="isIndeterminatePosDown" v-model="checkAllPosDown" @change="handleCheckAllChangePosDown">全选下方</el-checkbox>
                 </el-col>
                 <el-col :span="20">
                   <el-checkbox-group v-model="ABform.ad_pos_restrict_down" @change="handleCheckedChangePosDown">
                     <el-checkbox v-for="pos in ad_pos_down" :label="pos" :key="pos">{{pos}}</el-checkbox>
                   </el-checkbox-group>
                 </el-col>
               </el-form-item>
             </el-form-item>
         </el-row>

         <el-row>
             <el-form-item label="全局IP" prop="ip_restrict">
               <el-col :span="22">
                 <el-checkbox-group v-model="ABform.ip_restrict">
                   <el-checkbox
                     v-for="value in ip"
                     :key="value"
                     :label="value"
                     :value="value">
                   </el-checkbox>
                 </el-checkbox-group>
               </el-col>
             </el-form-item>
         </el-row>

         <h4>实验路样式判断</h4>
         <el-row>
           <el-form-item label="new_extend_reserve">
             <el-col :span="4" :offset="2">
               <span>PV日志$61字段 rshift:</span>
             </el-col>
             <el-col :span="2">
               <el-form-item prop="expr_extend_reserve.rshift" required>
                 <el-input v-model="ABform.expr_extend_reserve.rshift" size="mini" clearable></el-input>
               </el-form-item>
             </el-col>
             <el-col :span="1" :offset="1">
               <span>and:</span>
             </el-col>
             <el-col :span="2">
               <el-form-item prop="expr_extend_reserve.and" required>
                 <el-input v-model="ABform.expr_extend_reserve.and" size="mini" clearable></el-input>
               </el-form-item>
             </el-col>
             <el-col :span="1" :offset="1">
               <span>value:</span>
             </el-col>
             <el-col :span="2">
               <el-form-item prop="expr_extend_reserve.value" required>
                 <el-input v-model="ABform.expr_extend_reserve.value" size="mini" clearable></el-input>
               </el-form-item>
             </el-col>
             <el-col :span="2" :offset="1">
               <span>逗号分隔第</span>
             </el-col>
             <el-col :span="2">
               <el-form-item prop="expr_extend_reserve.pos" required>
                 <el-input v-model="ABform.expr_extend_reserve.pos" size="mini" clearable></el-input>
               </el-form-item>
             </el-col>
             <el-col :span="1">
               <span>位</span>
             </el-col>
           </el-form-item>
         </el-row>

         <hr>
         <h4>对照路样式判断</h4>
         <el-row>
           <el-form-item label="new_extend_reserve">
              <el-col :span="4" :offset="2">
                <span>PV日志$61字段 rshift:</span>
              </el-col>
              <el-col :span="2">
                <el-form-item prop="ctrl_extend_reserve.rshift" required>
                  <el-input v-model="ABform.ctrl_extend_reserve.rshift" size="mini" clearable></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="1" :offset="1">
                <span>and:</span>
              </el-col>
              <el-col :span="2">
                <el-form-item prop="ctrl_extend_reserve.and" required>
                  <el-input v-model="ABform.ctrl_extend_reserve.and" size="mini" clearable></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="1" :offset="1">
                <span>value:</span>
              </el-col>
              <el-col :span="2">
                <el-form-item prop="ctrl_extend_reserve.value" required>
                  <el-input v-model="ABform.ctrl_extend_reserve.value" size="mini" clearable></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="2" :offset="1">
                <span>逗号分隔第</span>
              </el-col>
              <el-col :span="2">
                <el-form-item prop="ctrl_extend_reserve.pos" required>
                  <el-input v-model="ABform.ctrl_extend_reserve.pos" size="mini" clearable></el-input>
                </el-form-item>
              </el-col>
              <el-col :span="1">
                <span>位</span>
              </el-col>
            </el-form-item>
         </el-row>


       </el-form>
    <div class="formCommit">
      <el-row type="flex" justify="end">
        <el-col :span="6">
          <el-button type="primary" @click="onSubmit('ABform')">立即提交</el-button>
          <el-button @click="onReset('ABform')">重置</el-button>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
  const axios = require('axios');
  require('promise/polyfill-done');

  import router from '../router';
  import global_ from './Global';

  export default {
    data(){
      return {
        dateOptions: global_.dateOptions,

        ad_pos_up : global_.abtestAdPosUpPc,
        ad_pos_down : global_.abtestAdPosDownPc,
        isIndeterminatePosUp: true,
        isIndeterminatePosMid: true,
        isIndeterminatePosDown: true,
        checkAllPosUp: false,
        checkAllPosMid: false,
        checkAllPosDown: false,

        ip: global_.abtestIP.pc,

        ABrules: global_.nameRule,

        ABform: {
          name: '',
          date_range: global_.dateRange,
          type: 'PC小流量任务',
          ad_pos_restrict_up: [],
          ad_pos_restrict_down: [],
          ad_pos_restrict: [],
          ip_restrict: [],
          expr_extend_reserve:{
            rshift:'',
            and:'',
            value:'',
            pos:''
          },
          ctrl_extend_reserve: {
            rshift: '',
            and: '',
            value: '',
            pos: ''
          }
        }
      }
    },
    methods:{
      handleCheckAllChangePosUp(val) {
        this.ABform.ad_pos_restrict_up = val ? this.ad_pos_up : [];
        this.isIndeterminatePosUp = false;
      },
      handleCheckedChangePosUp(value) {
        let checkedCount = value.length;
        this.checkAllPosUp = checkedCount === this.ad_pos_up.length;
        this.isIndeterminatePosUp = checkedCount > 0 && checkedCount < this.ad_pos_up.length;
      },
      handleCheckAllChangePosMid(val) {
        this.ABform.ad_pos_restrict_mid = val ? this.ad_pos_mid : [];
        this.isIndeterminatePosMid = false;
      },
      handleCheckedChangePosMid(value) {
        let checkedCount = value.length;
        this.checkAllPosMid = checkedCount === this.ad_pos_mid.length;
        this.isIndeterminatePosMid = checkedCount > 0 && checkedCount < this.ad_pos_mid.length;
      },
      handleCheckAllChangePosDown(val) {
        this.ABform.ad_pos_restrict_down = val ? this.ad_pos_down : [];
        this.isIndeterminatePosDown = false;
      },
      handleCheckedChangePosDown(value) {
        let checkedCount = value.length;
        this.checkAllPosDown = checkedCount === this.ad_pos_down.length;
        this.isIndeterminatePosDown = checkedCount > 0 && checkedCount < this.ad_pos_down.length;
      },
      redirectClick (){
        router.push({path: 'abproject'})
      },
      errorNotify(msg) {
        this.$notify.error({
          title: '错误',
          message: msg
        });
      },
      onSubmit(formName){
        let self= this;
        this.ABform.ad_pos_restrict = this.ABform.ad_pos_restrict_up.concat(this.ABform.ad_pos_restrict_down);
        this.$refs[formName].validate((valid)=>{
            if (valid){
              axios({
                method:'post',
                url:'/small/submit2pc',
                data: JSON.stringify(self.ABform)
              }).then(function (response) {
                self.$notify({
                  title: response.data,
                  message: '服务器接收到您的请求并分配查询ID。您可以转到查询页面查看项目！',
                  type: 'success',
                  duration: 2000,
                  onClick: self.redirectClick,
                });
              }).catch(function (error) {
                self.$options.methods.errorNotify(error);
              });
            } else {
              self.$options.methods.errorNotify("请补充内容后重提交！");
            }
          }
        )
      },
      onReset(formName){
        this.$refs[formName].resetFields();
      },
      getIp(){
        let self = this;
        axios.get('small/ip/pc')
          .then((response) => {
            var arrIp = []
            response.data.forEach(function (item) {
              arrIp = arrIp.concat(item['ips'])
            })
            self.ip = arrIp
        })
      },
      loadHistory(){
        let self = this;
        const hist = JSON.parse(sessionStorage.getItem('abtest_resubmit'));
        sessionStorage.removeItem('abtest_resubmit');
        if (hist) {
          for (var ele in hist){
            self.ABform[ele] = hist[ele];
          }
        }
      },
    },
    mounted: function () {
      this.loadHistory();
    },
    created: function () {
      this.getIp();
    }
  }
</script>

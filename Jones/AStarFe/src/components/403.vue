<template>
  <div>
    <div class="page-container">
      <p>403 Forbidden</p>
      <p>您目前没有该系统查看权限！请先申请权限</p>
      <h6>若已申请, 请转至<a href="http://datacenter.adtech.sogou">http://datacenter.adtech.sogou</a></h6>
      <div align="left">
        <el-col :span="15" offset="9">
          <h6>一 使用xiaoP扫描如下账号关注</h6>
          <img src="../assets/xiaopScan.jpeg" width="40%"/>
          <h6>二 输入"申请权限"</h6>
          <img src="../assets/xiaopApply.png" width="40%"/>
          <h6>三 填写申请单, 等待审核通过</h6>
          <img src="../assets/xiaopForm.png" width="40%"/>
        </el-col>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    name: ''
  }
</script>

<style scoped>
  .page-container {
    margin: 40px;
    font-size: 30px;
    text-align: center;
    color: rgb(192, 204, 218);
  }
</style>

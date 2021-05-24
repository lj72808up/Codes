<template>
  <div>
    <h3 style="margin-bottom: 0px; margin-top: 5px;">数据资产</h3>
    <el-row>
      <el-col :span="8">
        <span class="el-form-item__label">模糊查询:</span>
        <el-input v-model="searchWord"
                  size="mini"
                  clearable
                  style="width: 50%; margin-bottom: 5px"
                  placeholder="输入关键词"
                  @input="searchItems"/>
      </el-col>
      <el-col :span="8">
        <span class="el-form-item__label">资产类型:</span>
        <el-select v-model="assetType" filterable placeholder="请选择" clearable
                   style="width: 50%" size="mini" @change="searchItems">
          <el-option
            v-for="ele in assetList"
            :key="ele.id"
            :label="ele.name"
            :value="ele.id">
          </el-option>
        </el-select>
      </el-col>
      <el-col :span="8">
        <span class="el-form-item__label">业务标签:</span>
        <el-select v-model="businessLabel" filterable placeholder="请选择" clearable
                   style="width: 50%" size="mini" @change="searchItems">
          <el-option
            v-for="ele in serviceTagList"
            :key="ele.tid"
            :label="ele.tagName"
            :value="ele.tid">
          </el-option>
        </el-select>
      </el-col>
    </el-row>
    <el-table :data="dataList" stripe style="width: 100%" :border="false"
              :header-cell-style="{background:'#eee',color:'#606266'}">
      <el-table-column key="type_name" prop="type_name" label="资产类型" min-width="18%"/>
      <el-table-column key="item_name" prop="item_name" label="资产名" min-width="18%">
        <template slot-scope="scope">
          <span @click="browseTableInfo(scope.row.origin_id, scope.row.type_name)"
                style="cursor:pointer;text-decoration: underline">{{scope.row.item_name}}</span>
        </template>
      </el-table-column>
      <el-table-column key="description" prop="description" label="描述" min-width="46%"/>
      <el-table-column key="business_label" prop="business_label" label="业务标签" min-width="18%">
        <template slot-scope="scope">
          <el-tag style="margin-right: 5px" v-for="(businessLabel, index) in scope.row.business_label.split(',')"
                  v-show="scope.row.business_label.split(',').length>0 && scope.row.business_label.split(',')[0]!==''">
            {{serviceTagMap[businessLabel]}}
          </el-tag>
        </template>
      </el-table-column>
    </el-table>
  </div>

</template>

<script>
  import router from '../../router'

  const axios = require('axios')
  export default {
    data() {
      return {
        dataList: [],
        tableColumn: [],
        serviceTagList: [],
        serviceTagMap: {},
        searchWord: '',
        businessLabel: '',
        assetType: '',
        assetList: [{'id': 'predefine', 'name': '原子指标'},
          {'id': 'combine', 'name': '计算指标'},
          {'id': 'dimension', 'name': '维度'}],
      }
    },
    methods: {
      browseTableInfo(id, accessType) {
        console.log(accessType + "::" + id);
        switch (accessType) {
          case "原子指标":
            sessionStorage.setItem('preIndicatorId', id);
            sessionStorage.setItem('preIndicatorCheck', 'true');
            router.push({path: 'preDefineIndicatorCreate'});
            break;
          case "计算指标":
            sessionStorage.setItem('combineIndicatorId', id);
            sessionStorage.setItem('combineIndicatorCheck', 'true');
            router.push({path: 'combineIndicatorCreate'});
            break;
          case  "维度":
            sessionStorage.setItem('dimensionId', id);
            sessionStorage.setItem('dimensionCreate', 'true');
            router.push({path: 'dimensionCreate'});
            break;
          default:
        }
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
      searchItems() {
        let self = this
        let params = {
          'searchWord': this.searchWord,
          'businessLabel': this.businessLabel,
          'assetType': this.assetType,
        }
        axios.post('/indicator/dictionary/search', params, {responseType: 'json'})
          .then(function (response) {
            self.dataList = response.data.data
          })
          .catch(function (error) {
            console.log(error)
          })
      },
      getAllDictionary() {
        let self = this
        axios.get('/indicator/dictionary/getAll')
          .then(function (response) {
            self.dataList = response.data.data
          }).catch(function (error) {
          console.log(error)
        })
      }
    },
    mounted: function () {
      this.getAllServiceTag(this.getAllDictionary())
    }
  }
</script>

<style>
  .el-table .cell {
    white-space: pre-line !important;
  }

</style>

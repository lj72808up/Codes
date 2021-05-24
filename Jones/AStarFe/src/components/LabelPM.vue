<template>
  <div>
    <h3>配置标签关联</h3>
    <el-row>
      <div>
        <el-table
          :data="labelData"
          ref="labelTable"
          row-key="labelId"
          stripe
          @selection-change="handleURLSelectionChange">
          <!--          <el-table-column property="labelId" label="标签id"/>-->
          <el-table-column prop="labelTitle2" label="标签名称" sortable>
            <template slot-scope="scope">
              <el-input v-if="scope.row.status" v-model="scope.row.labelTitle2"></el-input>
              <span v-else>{{scope.row.labelTitle2}}</span>
            </template>
          </el-table-column>
          <el-table-column property="labelClass2" label="标签图标">
            <template slot-scope="scope">
              <e-icon-picker v-if="scope.row.status" v-model="scope.row.labelClass2"/>
              <i v-else :class="scope.row.labelClass2"/>
            </template>
          </el-table-column>
          <!--          <el-table-column property="labelPid2" label="上级标签id"/>-->
          <el-table-column label="上级标签名称" sortable>
            <template slot-scope="scope">
              <el-select v-if="scope.row.status" v-model=scope.row.labelPTitle2 @change="choosePidTitle(scope.row)">
                <el-option
                  v-for="item in labelFatherList"
                  :key="item.labelFatherPid"
                  :label="item.labelFatherTitle"
                  :value="item.labelFatherTitle">
                </el-option>
              </el-select>
              <span v-else>{{scope.row.labelPTitle}}</span>
            </template>
          </el-table-column>
          <el-table-column property="sortId" label="排序ID" sortable>
            <template slot-scope="scope">
              <el-input-number v-model="scope.row.sortId" @change="changeCount(scope.row)"
                               :min="1"
                               :max="10000"
                               size="small"/>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="rdata">
              <el-button @click="UpdateLabel(rdata.row)" type="text" size="mini">修改</el-button>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="rdata">
              <el-button @click="CancelUpdateLabel(rdata.row)" type="text" size="mini">放弃修改</el-button>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="rdata">
              <el-button @click="SaveLabel(rdata.row)" type="text" size="mini">保存</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-row>
  </div>
</template>

<script>
  const axios = require('axios')
  import {getCookie, isEmpty} from '../conf/utils'

  export default {
    data () {
      return {
        labelData: [],
        labelPTitleToPid: {},
        labelFatherList: [],

      }
    },
    methods: {
      changeCount (row) {
        console.log(row)
        let labelId = row.labelId
        let sortId = row.sortId
        let param = {
          'labelId': labelId+'',
          'sortId': sortId+''
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
      UpdateLabel (curLabel) {
        curLabel.status = true
      },
      CancelUpdateLabel (curLabel) {
        curLabel.labelTitle2 = curLabel.labelTitle,
          curLabel.labelClass2 = curLabel.labelClass,
          curLabel.labelPid2 = curLabel.labelPid,
          curLabel.labelPTitle2 = curLabel.labelPTitle,
          curLabel.status = false
      },
      choosePidTitle (curLabel) {
        curLabel.labelPid2 = this.labelPTitleToPid[curLabel.labelPTitle2]
      },
      SaveLabel (curLabel) {
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

    },
    created: function () {
      let self = this
      axios.get('/auth/labels/pmlabels')
        .then(function (response) {
          let nameRes = []
          let checkedIdx = []
          let index = 0
          response.data.forEach(ele => {
            nameRes.push({
              'labelId': ele.Lid,
              'labelTitle': ele.Title,
              'labelClass': ele.Class,
              'labelPid': ele.Pid,
              'labelPTitle': ele.Ptitle,
              'labelTitle2': ele.Title,
              'labelClass2': ele.Class,
              'labelPid2': ele.Pid,
              'labelPTitle2': ele.Ptitle,
              'status': false,
              'sortId': ele.SortId,
            })

            //将子节点的父亲节点加入map
            self.labelPTitleToPid[ele.Ptitle] = ele.Pid
            //将根节点的所有子节点加入map
            if (ele.Pid == 0) {
              self.labelPTitleToPid[ele.Title] = ele.Lid
            }
          })
          self.labelData = nameRes

          //获取所有父节点信息
          for (var title in self.labelPTitleToPid) {
            self.labelFatherList.push({
              labelFatherPid: self.labelPTitleToPid[title],
              labelFatherTitle: title
            })
          }

        }).catch(function (error) {
        console.log(error)
      })
    },
    mounted () {
      console.log('mounted')
    }
  }
</script>

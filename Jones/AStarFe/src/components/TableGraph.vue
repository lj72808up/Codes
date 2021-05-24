/* eslint-disable */
<template>
  <div>
    <div>
      <el-row>
        <h3>
          {{title}}:
        </h3>
      </el-row>

      <el-row>
        <div class="chart-wrapper">
          <v-chart :options="myOptions" @click="onClick"/>
        </div>
      </el-row>


      <el-row>
        <el-button type="text" size="small" @click="handleClose()">回退</el-button>
      </el-row>


    </div>


    <el-dialog :title="source_target" :visible.sync="dialogVisible" width="80%">
      <vxe-table
        border
        show-overflow
        keep-source
        ref="xTable"
        :data="dataList"
        :edit-config="{trigger: 'click', mode: 'cell', showStatus: true, icon: 'fa fa-pencil-square-o'}">
        <!--显示字段-->
        <template v-for="(col,idx) in tableColumn">
          <vxe-table-column :field="col.field" :title="col.title" :visible="col.visible" :edit-render="col.editRender"/>
        </template>
      </vxe-table>
    </el-dialog>

  </div>
</template>

<style scoped>
  .chart-wrapper {
    width: 100%;
    height: 600px;
  }

  .echarts {
    width: 100%;
    height: 90%;
  }
</style>

<script>
  import router from '../router'

  import ECharts from 'vue-echarts'
  import 'echarts/lib/chart/line'
  import 'echarts/lib/component/polar'
  import 'echarts/lib/chart/graph'

  const axios = require('axios')

  export default {
    data () {
      return {
        info: 'csv内容形如: mysql,default_database,table1,field1,v0,我是描述',
        dialogVisible: false,
        dataList:[],
        tableColumn:[],
        title: '',
        source_target: '',
        myOptions: {
          autoresize: true,
          tooltip: {},
          animationDurationUpdate: 1500,
          animationEasingUpdate: 'quinticInOut',
          series: [
            {
              type: 'graph',
              layout: 'none',
              symbol: 'rect',
              symbolSize: [80, 40],  // 大小
              roam: true,
              label: {
                show: true
              },
              edgeSymbol: ['', 'arrow'],
              edgeSymbolSize: [10, 20],
              edgeLabel: {
                fontSize: 20
              },

              data: [],
              // links: [],
              links: [],
              lineStyle: {
                opacity: 0.9,
                width: 2,
                curveness: 0.2
              },
            }
          ]
        }
      }
    },
    methods: {
      handleClose () {
        sessionStorage.removeItem('metaData_graph_tid')
        router.push({path: 'tableMetaData'})
      },
      onClick (params) {
        console.log(params)
        if (params.dataType === 'edge') {
          // 点击到了 graph 的 edge（边）上。
          console.log('点到了边上')
          let param = {
            'source': params.data.source,
            'target': params.data.target
          }

          this.source_target = params.data.source+ " 和 " + params.data.target

          this.dialogVisible = true
          let self = this
          axios.post('/metaData/graph/relation/check', param, {responseType: 'json'})
            .then(function (response) {
              console.log(response)
              self.tableColumn = response.data.tableColumn
              self.dataList = response.data.dataList
            }).catch(function (error) {
            console.log(error)
          })

        } else {
          let clickTableId = params.data.tableId
          let info = {
            'tableId': clickTableId
          }
          window.sessionStorage.setItem('metaData_table_id', JSON.stringify(info))
          window.location.href = '#/tableInfo'
        }
      }
    }
    ,
    mounted: function () {
      let self = this

      let param = JSON.parse(sessionStorage.getItem('metaData_graph_tid'))
      console.log('tableId:' + param.tableId)
      // let tableId = 15
      let tableId = param.tableId
      this.title = param.tableName + '的拓扑关系'

      axios.get('metaData/graph/' + tableId)
        .then(function (response) {
          self.myOptions.series[0].data = response.data.data
          self.myOptions.series[0].links = response.data.links
        }).catch(function (error) {
        console.log(error)
      })

    }
  }
</script>

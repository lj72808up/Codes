<script>
  import {formatDate} from '../conf/utils'

  const yesterday = new Date()
  yesterday.setTime(yesterday.getTime() - 3600 * 1000 * 24)

  const ad_pos_up = ['0', '1', '2', '3', '4', '5', '6', '7']
  const ad_pos_mid = ['8', '9']
  const ad_pos_down = ['10', '11', '12', '13', '14', '15']
  const ad_pos_up_pc = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '10']
  const ad_pos_down_pc = ['11', '12', '13', '14', '15']
  const ip_north = ['10.134.115.72', '10.134.115.73', '10.134.115.74', '10.134.89.114',
    '10.139.17.51', '10.139.17.52', '10.139.17.53', '10.139.17.54', '10.139.20.72']
  const ip_south = ['10.135.8.67', '10.135.8.68', '10.135.8.69',
    '10.135.8.70', '10.135.8.71', '10.135.8.72']
  const ip_jiangsu = ['10.140.19.110', '10.140.19.123', '10.140.20.22']
  const ip_pc = ['10.134.89.114', '10.139.20.72',
    '10.139.20.59', '10.134.89.115', '10.134.105.25', '10.139.36.56', '10.139.21.59', '10.134.81.56',
    '10.135.73.91', '10.135.73.92', '10.135.73.93', '10.140.11.39', '10.140.19.67', '10.140.11.22']
  const ip_wap = ['10.139.17.51', '10.139.20.71', '10.134.73.35',
    '10.139.35.58', '10.134.89.67', '10.139.20.52', '10.134.104.119',
    '10.134.105.24', '10.139.36.101',
    '10.135.8.67', '10.135.66.32', '10.135.66.41', '10.135.73.41', '10.135.73.90',
    '10.140.19.110', '10.140.11.23', '10.140.26.80', '10.140.26.81']

  const wap_base_total = {
    '日期': 'RQ',
    '月份': 'RQ_YUEFEN',
    '季度': 'RQ_JIDU',
    '样式类别': 'YSLB',
    'PID': 'PID',
    '排位': 'PW',
    '深度': 'SD',
    '位置': 'WZ',
    '查询词': 'CXC',
    '查询词行业ID':'CXCID',
    '一级查询词行业': 'CXCHY_YIJI',
    '二级查询词行业': 'CXCHY_ERJI',
    '一级渠道': 'YJQD',
    '二级渠道': 'EJQD',
    '二级业务类型': 'EJYWLX',
    '三级业务类型': 'SJYWLX',
    '业务类型ID': 'YWLXID',
    '大区': 'DQ',
    '网民区域': 'WMQY'
  }

  const pc_base_total = {
    '日期': 'RQ',
    '月份': 'RQ_YUEFEN',
    '季度': 'RQ_JIDU',
    '一级渠道': 'YJQD',
    '二级渠道': 'EJQD',
    'PID': 'PID',
    '查询词': 'CXC',
    '查询词行业ID': 'CXCID',
    '一级查询词行业': 'CXCHY_YIJI',
    '二级查询词行业': 'CXCHY_ERJI',
    '二级业务类型': 'EJYWLX',
    '三级业务类型': 'SJYWLX',
    '业务类型ID': 'YWLXID',
    '大区': 'DQ',
    '网民区域': 'WMQY',
    '深度': 'SD'
  }

  const wap_base = {
    '日期': 'RQ',
    '月份': 'RQ_YUEFEN',
    '季度': 'RQ_JIDU',
    '查询词': 'CXC',
    '一级查询词行业': 'CXCHY_YIJI',
    '二级查询词行业': 'CXCHY_ERJI',
    '一级渠道': 'YJQD',
    '二级渠道': 'EJQD',
    '网民区域': 'WMQY'
  }
  const wap_GG = {
    '关键词ID': 'GJCID',
    '关键词': 'GJC',
    '创意ID': 'CYID',
    '客户ID': 'KHID',
    '计划ID': 'JHID',
    '计划名称': 'JHMC',
    '组ID': 'ZID',
    '组名称': 'ZMC',
    '广告类型': 'GGLX'
  }
  const wap_KH = {
    '客户ID': 'KHID',
    '客户名称': 'KHMC',
    '客户属性': 'KHSX',
    '一级客户行业': 'YJKHHY',
    '二级客户行业': 'EJKHHY',
    '三级客户行业': 'SJKHHY',
    '客户类型': 'KHLX',
    '代理商ID': 'DLSID',
    '代理商名称': 'DLSMC'
  }
  const wap_outputs = {'PV1': 'PV1', 'PV2': 'PV2', 'PV3': 'PV3', '点击': 'DJ', '消耗': 'XH'}

  const pc_base = {
    '日期': 'RQ',
    '月份': 'RQ_YUEFEN',
    '季度': 'RQ_JIDU',
    '查询词': 'CXC',
    '一级查询词行业': 'CXCHY_YIJI',
    '二级查询词行业': 'CXCHY_ERJI',
    '一级渠道': 'YJQD',
    '二级渠道': 'EJQD',
    '网民区域': 'WMQY'
  }
  const pc_GG = {
    '关键词ID': 'GJCID',
    '关键词': 'GJC',
    '创意ID': 'CYID',
    '计划ID': 'JHID',
    '计划名称': 'JHMC',
    '组ID': 'ZID',
    '组名称': 'ZMC',
    '广告类型': 'GGLX',
    '样式类别': 'YSLB',
    '排位': 'PW'
  }
  const pc_KH = {
    '客户ID': 'KHID',
    '客户名称': 'KHMC',
    '客户属性': 'KHSX',
    '客户类型': 'KHLX',
    '代理商ID': 'DLSID',
    '代理商名称': 'DLSMC',
    '一级客户行业': 'YJKHHY',
    '二级客户行业': 'EJKHHY',
    '三级客户行业': 'SJKHHY'
  }
  const pc_outputs = {'PV1': 'PV1', 'PV2': 'PV2', 'PV3': 'PV3', '点击': 'DJ', '消耗': 'XH'}

  const galaxy_base = {
    '日期': 'RQ', 'PID': 'PID', '查询词': 'CXC',
    '二级业务类型': 'EJYWLX', '三级业务类型': 'SJYWLX', '大区': 'DQ', '网民区域': 'WMQY'
  }
  const galaxy_GG = {
    '关键词ID': 'GJCID',
    '关键词': 'GJC',
    '创意ID': 'CYID',
    '计划ID': 'JHID',
    '计划名称': 'JHMC',
    '组ID': 'ZID',
    '组名称': 'ZMC'
  }
  const galaxy_KH = {
    '客户ID': 'KHID',
    '客户名称': 'KHMC',
    '客户类型': 'KHLX',
    '一级客户行业': 'YJKHHY',
    '二级客户行业': 'EJKHHY',
    '代理商ID': 'DLSID',
    '代理商名称': 'DLSMC'
  }

  export default {
    yesterday: yesterday,
    dateOptions: {
      shortcuts: [{
        text: '最近三天',
        onClick (picker) {
          const end = new Date()
          const start = new Date()
          end.setTime(end.getTime() - 3600 * 1000 * 24)
          start.setTime(start.getTime() - 3600 * 1000 * 24 * 3)
          picker.$emit('pick', [start, end])
        }
      },
        {
          text: '最近一周',
          onClick (picker) {
            const end = new Date()
            const start = new Date()
            end.setTime(end.getTime() - 3600 * 1000 * 24)
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            picker.$emit('pick', [start, end])
          }
        }, {
          text: '最近两周',
          onClick (picker) {
            const end = new Date()
            const start = new Date()
            end.setTime(end.getTime() - 3600 * 1000 * 24)
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 14)
            picker.$emit('pick', [start, end])
          }
        }]
    },
    dateRange: [formatDate(yesterday), formatDate(yesterday)],
    provincesChina: [
      {
        label: '华北地区',
        children: ['北京', '天津', '河北', '内蒙古', '山西']
      },
      {
        label: '东北地区',
        children: ['黑龙江', '吉林', '辽宁']
      },
      {
        label: '华东地区',
        children: ['上海', '福建', '江苏', '江西', '山东', '浙江']
      },
      {
        label: '华中地区',
        children: ['河南', '湖北', '湖南']
      },
      {
        label: '华南地区',
        children: ['广东', '广西', '海南']
      },
      {
        label: '西南地区',
        children: ['重庆', '贵州', '四川', '西藏', '云南']
      },
      {
        label: '西北地区',
        children: ['甘肃', '宁夏', '青海', '陕西', '新疆']
      },
      {
        label: '其他地区',
        children: ['港澳台及海外', '无此地区']
      },
      {
        label: '全国',
        children: ['全国']
      }
    ],
    nameRule: {
      name: [
        {required: true, message: '请输入名称', trigger: 'blur'},
        {min: 3, message: '长度至少在 3 个字符', trigger: 'blur'}
      ],
    },
    strictRules: {
      name: [
        {required: true, message: '请输入名称', trigger: 'blur'},
        {min: 3, message: '名称长度至少 3 个字符', trigger: 'blur'}
      ],
      dimensions: [
        {type: 'array', required: true, message: '请选择查询维度', trigger: 'change'},
      ],
      outputs: [
        {type: 'array', required: true, message: '请选择输出指标', trigger: 'change'},
      ]
    },

    //  abtest dimensions
    abtestAdPosUp: ad_pos_up,
    abtestAdPosMid: ad_pos_mid,
    abtestAdPosDown: ad_pos_down,
    abtestAdPosUpPc: ad_pos_up_pc,
    abtestAdPosDownPc: ad_pos_down_pc,
    abtestIP: {
      north: ip_north,
      south: ip_south,
      jiangsu: ip_jiangsu,
      pc: ip_pc,
      wap: ip_wap
    },

    //   wap dimensions
    wapBase: wap_base,
    wapBaseTotal: wap_base_total,
    wapGG: wap_GG,
    wapKH: wap_KH,
    wapOutputs: wap_outputs,

    //   pc dimensions
    pcBase: pc_base,
    pcBaseTotal: pc_base_total,
    pcGG: pc_GG,
    pcKH: pc_KH,
    pcOutputs: pc_outputs,

    //   galaxy dimensions
    galaxyBase: galaxy_base,
    galaxyGG: galaxy_GG,
    galaxyKH: galaxy_KH
  }
</script>

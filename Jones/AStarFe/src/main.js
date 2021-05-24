// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import ECharts from 'vue-echarts'

// 手动引入 ECharts 各模块来减小打包体积
// 引入提示框和标题组件
import 'echarts/lib/component/tooltip'
import 'echarts/lib/component/title'
// 注册组件后即可使用
Vue.component('v-chart', ECharts)

import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';


import BootstrapVue from 'bootstrap-vue'
Vue.use(BootstrapVue);
// import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

Vue.use(ElementUI);
Vue.config.productionTip = false


import plTable from 'pl-table'
import 'pl-table/themes/index.css'
import 'pl-table/themes/plTableStyle.css'
Vue.use(plTable);

import MintUI from 'mint-ui'
import 'mint-ui/lib/style.css'
Vue.use(MintUI)

//引入配置的http文件
import axios from './http'

import VueClipboard from 'vue-clipboard2'
Vue.use(VueClipboard)

import DAGBoard from './dag-diagram/index'
Vue.use(DAGBoard)

// import VueEditor from "vue2-editor";
// Vue.use(VueEditor)

import VueJsonp from 'vue-jsonp'
Vue.use(VueJsonp)

// codemirror
import VueCodeMirror from 'vue-codemirror'
import 'codemirror/lib/codemirror.css'
Vue.use(VueCodeMirror)


import './assets/global.css'

import VueContextMenu from '@xunlei/vue-context-menu'
Vue.use(VueContextMenu)

//图标选择器
import iconPicker from 'e-icon-picker';
import 'e-icon-picker/dist/index.css';//基础样式
import 'e-icon-picker/dist/main.css'; //fontAwesome 图标库样式
Vue.use(iconPicker);

// vxe-table
import 'xe-utils'
import VXETable from 'vxe-table'
import 'vxe-table/lib/index.css'
Vue.use(VXETable)
// 给 vue 实例挂载全局窗口对象
Vue.prototype.$XModal = VXETable.modal


import VueQuillEditor from 'vue-quill-editor'
import 'quill/dist/quill.core.css' // import styles
import 'quill/dist/quill.snow.css' // for snow theme
import 'quill/dist/quill.bubble.css' // for bubble theme
Vue.use(VueQuillEditor)


new Vue({
  el: '#app',
  router,
  render: h => h(App),
  components: { App },
  template: '<App/>'
})




import SVGContent from './dagContent'
import NodeBus from './nodeBus'
import {getCookie,setCookie} from '../conf/utils'

const DAGBoard = {
  install(Vue, options) {
    Vue.component(SVGContent.name, SVGContent)
    Vue.component(NodeBus.name, NodeBus)
  }
}

if (typeof window !== 'undefined' && window.Vue) {
  window.Vue.use(DAGBoard)
}

export default DAGBoard

/*let ptoken = getCookie('ptoken')
alert("ptoken: "+ ptoken)
if (ptoken==="" || ptoken===null){
  alert("跳转了")
  window.location.href="https://login.sogou-inc.com/?appid=1755&sso_redirect=http://datacenter.adtech.sogou/"
} else {
  alert("不用跳转")
}*/


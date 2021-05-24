import Vue from 'vue'
import Router from 'vue-router'
import Forbidden from '../components/403'
import NotFound from '../components/404'
import Home from '../components/Layout'
import SqlQuery from '../components/SqlQuery'
import Queue from '../components/Queue'
import WapJohn from '../components/WapJohn'
import WapJohnTotal from '../components/WapJohnTotal'
import PCJohn from '../components/PCJohn'
import PCJohnTotal from '../components/PCJohnTotal'
import Galaxy from '../components/GalaxyJohn'
import Role from '../components/Role'
// import User from '../components/User'
import Trans from '../components/Trans'
import God from '../components/God'
import ABwap from '../components/ABwap'
import ABpc from '../components/ABpc'
import ABproject from '../components/ABproject'
import ABdata from '../components/ABdata'
import Doc from '../components/Doc'
import DocEdit from '../components/DocEdit'
import JdbcQuery from '../components/JdbcQuery'
import TableMetaData from '../components/TableMetaData'
import TableGraph from '../components/TableGraph'
import TableInfo from '../components/TableInfo'
import Test from '../components/Test'
import RoleUser from '../components/RoleUser'
import RoleWorkSheet from '../components/RoleWorkSheet'
import RoleTopLable from '../components/RoleTopLable'
import RedashWXRealTimeCost from '../components/RedashWXRealTimeCost'
import RedashInsert from '../components/RedashInsert'
import RedashShow from '../components/RedashShow'
import RedashDelete from '../components/RedashDelete'
import LabelPM from '../components/LabelPM'
import LabelRD from '../components/LabelRD'
import Diagram from '../components/Diagram'
import ConnMetaData from '../components/ConnMetaData'
import RoleMetaTable from '../components/RoleMetaTable'
import Airflow from '../components/Airflow'
import DagTaskTrace from '../components/DagTaskTrace'
import DashboardShow from '../components/DashboardShow'
import DashboardOption from '../components/DashboardOption'
import RedashQueryShow from '../components/RedashQueryShow'
import RedashQueryOption from '../components/RedashQueryOption'
import MobileDashBoard from '../components/MobileDashBoard'
import MobileDashBoardShow from '../components/MobileDashBoardShow'
import Udf from '../components/Udf'
import ServiceLabel from '../components/indicator/ServiceLabel'
import PreDefineIndicator from '../components/indicator/PreDefineIndicator'
import CombineIndicator from '../components/indicator/CombineIndicator'
import Dimension from '../components/indicator/Dimension'
import Dictionary from '../components/indicator/Dictionary'
import DemandRequirement from '../components/DemandRequirement'
import DemandList from "../components/DemandList";
import DimensionCreate from "../components/indicator/DimensionCreate";
import PreDefineIndicatorCreate from "../components/indicator/PreDefineIndicatorCreate";
import CombineIndicatorCreate from "../components/indicator/CombineIndicatorCreate";
import DimensionCreateDerivative from "../components/indicator/DimensionCreateDerivative"
import RedashLvpi from '../components/RedashLvpi'
import RoleMobileBodong from "../components/RoleMobileBodong"
import RoleDomain from "../components/RoleDomain"
import DimStyleClassify from "../components/singleTable/DimStyleClassify"

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: '',
      component: Home,
      children: [
        {
          path: 'jdbc',
          name: 'jdbc',
          component: JdbcQuery
        },
        {
          path: 'sql',
          name: 'sql',
          component: SqlQuery
        },
        {
          path: 'wap',
          name: 'wap',
          component: WapJohn
        },
        {
          path: 'pc',
          name: 'pc',
          component: PCJohn
        },
        {
          path: 'wapTotal',
          name: 'wapTotal',
          component: WapJohnTotal
        },
        {
          path: 'pcTotal',
          name: 'pcTotal',
          component: PCJohnTotal
        },
        {
          path: 'gala',
          name: 'gala',
          component: Galaxy
        },
        {
          path: 'que',
          name: 'que',
          component: Queue
        },
        {
          path: 'role',
          name: 'role',
          component: Role
        },
        /*{
          path: 'user',
          name: 'user',
          component: User
        },*/
        {
          path: 'roleuser',
          name: 'roleuser',
          component: RoleUser
        },
        {
          path: 'roleworksheet',
          name: 'roleworksheet',
          component: RoleWorkSheet
        },
        {
          path: 'rolemetatable',
          name: 'rolemetatable',
          component: RoleMetaTable
        },
        {
          path: 'roletoplable',
          name: 'roletoplable',
          component: RoleTopLable
        },
        {
          path: 'labelpm',
          name: 'labelpm',
          component: LabelPM,
        },
        {
          path: 'labelrd',
          name: 'labelrd',
          component: LabelRD,
        },
        {
          path: 'airflow',
          name: 'airflow',
          component: Airflow,
        },
        {
          path: 'redashwxrealtimecost',
          name: 'redashwxrealtimecost',
          component: RedashWXRealTimeCost
        },
        {
          path: 'redashinsert',
          name: 'redashinsert',
          component: RedashInsert
        },
        {
          path: 'redashdelete',
          name: 'redashdelete',
          component: RedashDelete,
        },
        {
          path: 'bodonginsert',
          name: 'bodonginsert',
          component: RedashInsert
        },
        {
          path: 'bodongdelete',
          name: 'bodongdelete',
          component: RedashDelete,
        },
        {
          path: 'redashshow/:redashname',
          name: 'redashshow',
          component: RedashShow,
        },
        {
          path: 'dashboardoption',
          name: 'dashboardoption',
          component: DashboardOption,
        },
        {
          path: 'dashboardshow/:dashboardname',
          name: 'dashboardshow',
          component: DashboardShow,
        },
        {
          path: 'redashlvpi',
          name: 'redashlvpi',
          component: RedashLvpi,
        },
        {
          path: 'redashqueryoption',
          name: 'redashqueryoption',
          component: RedashQueryOption,
        },
        {
          path: 'redashqueryshow',
          name: 'redashqueryshow',
          component: RedashQueryShow,
        },
        {
          path: 'trans',
          name: 'trans',
          component: Trans
        },
        {
          path: 'abwap',
          name: 'abwap',
          component: ABwap
        },
        {
          path: 'abpc',
          name: 'abpc',
          component: ABpc
        },
        {
          path: 'abproject',
          name: 'abproject',
          component: ABproject
        },
        {
          path: 'abdata',
          name: 'abdata',
          component: ABdata
        },
        {
          path: 'god',
          name: 'god',
          component: God
        },
        {
          path: 'docedit',
          name: 'docedit',
          component: DocEdit
        },
        {
          path: 'doc',
          name: 'doc',
          component: Doc
        },
        {
          path: 'tableMetaData',
          name: 'tableMetaData',
          component: TableMetaData
        },
        {
          path: 'tableGraph',
          name: 'tableGraph',
          component: TableGraph
        }, {
          path: 'tableInfo',
          name: 'tableInfo',
          component: TableInfo
        }, {
          path: 'test',
          name: 'test',
          component: Test
        }, {
          path: 'dag',
          name: 'dag',
          component: Diagram
        }, {
          path: 'connMetaData',
          name: 'connMetaData',
          component: ConnMetaData
        }, {
          path: 'dagTaskTrace',
          name: 'dagTaskTrace',
          component: DagTaskTrace
        }, {
          path: 'mobileDashBoard',
          name: 'mobileDashBoard',
          component: MobileDashBoard
        }, {
          path: 'mobileDashBoardShow/:redashname',
          name: 'mobileDashBoardShow',
          component: MobileDashBoardShow,
        }, {
          path: 'udf',
          name: 'udf',
          component: Udf,
        }, {
          path: 'serviceLabel',
          name: 'serviceLabel',
          component:ServiceLabel
        }, {
          path: 'preIndicator',
          name: 'preIndicator',
          component:PreDefineIndicator
        }, {
          path: 'preDefineIndicatorCreate',
          name: 'preDefineIndicatorCreate',
          component:PreDefineIndicatorCreate
        }, {
          path: 'combineIndicator',
          name: 'combineIndicator',
          component:CombineIndicator
        }, {
          path: 'combineIndicatorCreate',
          name: 'combineIndicatorCreate',
          component:CombineIndicatorCreate
        }, {
          path: 'dimension',
          name: 'dimension',
          component:Dimension
        }, {
          path: 'dimensionCreate',
          name: 'dimensionCreate',
          component:DimensionCreate
        },{
          path: 'dimensionCreateDerivative',
          name: 'dimensionCreateDerivative',
          component:DimensionCreateDerivative
        }, {
          path: 'dictionary',
          name: 'dictionary',
          component:Dictionary
        }, {
          path: 'demand/:demandId',
          name: 'demand',
          component:DemandRequirement
        }, {
          path: 'demandList',
          name: 'demandList',
          component:DemandList
        }, {
          path: 'roleMobileBodong',
          name: 'roleMobileBodong',
          component:RoleMobileBodong
        }, {
          path: 'roleDomain',
          name: 'roleDomain',
          component: RoleDomain
        }, {
          path: 'dimsc',
          name: 'dimsc',
          component: DimStyleClassify
        }
      ]
    },
    {
      path: '/403',
      component: Forbidden,
      name: '',
      hidden: false
    },
    {
      path: '/404',
      component: NotFound,
      name: '',
      hidden: false
    }
  ]
})



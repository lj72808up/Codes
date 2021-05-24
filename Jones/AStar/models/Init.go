package models

import (
	"AStar/models/indicator"
	"AStar/utils"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/authz"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/casbin/beego-orm-adapter"
	"github.com/casbin/casbin"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var Enforcer *casbin.Enforcer

var DBConfig DBConfigModel

var RedisConfig RedisConfigModel

var TriggerDagConfig TriggerDagModel

var ServerInfo ServerInfoModel

var KylinInfo KylinInfoModel

var WsSQLModel SqlModel

var PcSQLModel SqlModel

var YhSQLModel SqlModel

var AirflowInfo AirflowInfoModel

var WsFactProiOrdList []SqlFactTableModel
var PcFactProiOrdList []SqlFactTableModel

var Logger *logs.BeeLogger

var ManualTriggerDagChannel = make(chan DagTriggerTask, 9999) // 生成一个长度为 9999 的缓存 channel

func Init() {

	// init logger
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/beego.log"}`)
	Logger = logs.NewLogger()
	_ = Logger.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/AStar.log","separate":["error", "info", "debug"]}`)
	_ = Logger.SetLogger(logs.AdapterConsole)
	Logger.Async(1e3)
	Logger.EnableFuncCallDepth(true)

	utils.Logger = Logger
	// CORS是一个W3C标准，全称是"跨域资源共享"（Cross-origin resource sharing）
	// Beego中支持跨域
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

	//init DB, ServerInfo, WsSQLModel, KylinInfo,RedisInfo
	utils.LoadModelFromConfig(&DBConfig)
	utils.LoadModelFromConfig(&RedisConfig)
	utils.LoadModelFromConfig(&ServerInfo)
	utils.LoadModelFromConfig(&KylinInfo)
	utils.LoadModelFromConfig(&AirflowInfo)
	utils.LoadModelFromConfig(&TriggerDagConfig)

	WsSQLModel, _ = SqlModel{}.InitSQL("conf/wsSqlConfig.json")
	PcSQLModel, _ = SqlModel{}.InitSQL("conf/pcSqlConfig.json")
	YhSQLModel, _ = SqlModel{}.InitSQL("conf/yhSqlConfig.json")

	WsFactProiOrdList = WsSQLModel.InitFactsPriorityOrderedList()
	PcFactProiOrdList = PcSQLModel.InitFactsPriorityOrderedList()

	//init orm QuerySubmitLog
	orm.RegisterModel(new(QuerySubmitLog))
	orm.RegisterModel(new(SmallFlowTask))
	orm.RegisterModel(new(FrontFilterModel))
	orm.RegisterModel(new(PlatformDoc))
	orm.RegisterModel(new(SchedulerJob))
	orm.RegisterModel(new(SqlTemplate))
	orm.RegisterModel(new(GroupWorksheetMapping))
	orm.RegisterModel(new(TopLabel))
	orm.RegisterModel(new(GroupLabelMapping))
	orm.RegisterModel(new(TableMetaData))
	orm.RegisterModel(new(FieldMetaData))
	orm.RegisterModel(new(VersionMetaData))
	orm.RegisterModel(new(RelationMetaData))
	orm.RegisterModel(new(TableRelation))
	orm.RegisterModel(new(TableTypes))
	orm.RegisterModel(new(NameSpaceWritable))
	orm.RegisterModel(new(ConnectionManager))
	orm.RegisterModel(new(GroupMetatableMapping))
	orm.RegisterModel(new(GroupConnMapping))
	orm.RegisterModel(new(DagAirFlow))
	orm.RegisterModel(new(DagAirFlowImg))
	orm.RegisterModel(new(XiaopApply))
	orm.RegisterModel(new(SqlTemplateHeat))
	orm.RegisterModel(new(FrontRedashDashboard))
	orm.RegisterModel(new(FrontRedashQuery))
	orm.RegisterModel(new(UdfManager))
	orm.RegisterModel(new(indicator.Indicator))
	orm.RegisterModel(new(indicator.ServiceTag))
	orm.RegisterModel(new(indicator.IndicatorRelation))
	orm.RegisterModel(new(indicator.Dimension))
	orm.RegisterModel(new(indicator.DimensionFactRelation))
	orm.RegisterModel(new(indicator.IndicatorFactRelation))
	orm.RegisterModel(new(indicator.Qualifier))
	orm.RegisterModel(new(BodongMenuBar))
	orm.RegisterModel(new(DataRequirement))
	orm.RegisterModel(new(DemandAttachment))
	orm.RegisterModel(new(DemandAssetRelation))
	orm.RegisterModel(new(DemandUserDefineIndicator))
	orm.RegisterModel(new(BodongMenuMapping))
	orm.RegisterModel(new(DagTriggerHistory))
	orm.RegisterModel(new(DataRequirementComment))
	orm.RegisterModel(new(UserNameMapping))
	orm.RegisterModel(new(DataRequirementHistory))
	orm.RegisterModel(new(DomainLvpi))
	orm.RegisterModel(new(LvpiMapping))

	orm.Debug = true

	adapter := beegoormadapter.NewAdapter(beego.AppConfig.String("drivername"), beego.AppConfig.String("datasourcename"), true)
	Enforcer = casbin.NewEnforcer("conf/authz_model.conf", adapter)
	Enforcer.EnableAutoSave(true)
	_ = Enforcer.LoadPolicy()

	beego.InsertFilter("jones/*", beego.BeforeRouter, authz.NewAuthorizer(Enforcer))

	// 定时任务
	AddGlobalTask()

	SyncRemoteAirflow()

	// 全局的 Dag Channel
	//go GetAndRunTriggerDag()
	CreateWorkerPool(TriggerDagConfig.WorkerNum)
	Logger.Info("启动 %d 个worker执行任务", TriggerDagConfig.WorkerNum)

	// 注册数据源
	url := "adtd:noSafeNoWork2016@tcp(adtd.config.rds.sogou:3306)/adtd_config?charset=utf8&loc=Asia%2FShanghai"
	err := orm.RegisterDataBase("DimStyleClassify", "mysql", url, 0, 1) // SetMaxIdleConns, SetMaxOpenConns
	if err != nil {
		panic(err)
	}
}

func SyncRemoteAirflow() {
	airbase := beego.AppConfig.String("airflowFileRoot")
	utils.MkDir(airbase)
	airBaseArr := strings.Split(airbase, "/")
	airbase = strings.Join(airBaseArr[:len(airBaseArr)-1], "/")

	remoteBase := fmt.Sprintf("%s/*", beego.AppConfig.String("airflowRemoteFileRoot2"))
	//remoteBase := beego.AppConfig.String("airflowRemoteFileRoot2")
	utils.MkDir(airbase)
	cmd := fmt.Sprintf("rsync -aP --delete %s %s", remoteBase, airbase)
	fmt.Println("rsync 执行:", cmd)
	splits := strings.Split(cmd, " ")
	utils.ExecShell(splits[0], splits[1:])
}

// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"AStar/controllers"
	"AStar/controllers/indicator"
	"AStar/utils"
	"github.com/astaxie/beego"
)

var rootNameSpace = controllers.RootNameSpace

func init() {

	ns := beego.NewNamespace(rootNameSpace,
		beego.NSInclude(
			&controllers.AuthController{},
		),
		beego.NSInclude(
			&controllers.DWController{},
		),
		beego.NSInclude(
			&controllers.QueryWsController{},
		),
		beego.NSInclude(
			&controllers.QueryPcController{},
		),
		beego.NSInclude(
			&controllers.QueryYhController{},
		),
		beego.NSInclude(
			&controllers.TaskController{},
		),
		beego.NSInclude(
			&controllers.SmallFlowController{},
		),
		beego.NSInclude(
			&controllers.DocController{},
		),
		beego.NSInclude(
			&controllers.TopLabelController{},
		),
		beego.NSInclude(
			&controllers.SchedulerJobController{},
		),
	)
	beego.AddNamespace(ns)

	tns := beego.NewNamespace(controllers.XiaopNameSpace,
		beego.NSInclude(
			&controllers.XiaopController{},
		),
	)
	beego.AddNamespace(tns)

	lvpiNS := beego.NewNamespace(controllers.LvpiNameSpace,
		beego.NSInclude(
			&controllers.LvpiAuthController{},
		),
	)
	beego.AddNamespace(lvpiNS)

	testController := controllers.TestController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, controllers.TestNameSpace),
		beego.NSInclude(&testController)))

	jdbcController := controllers.JdbcController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, controllers.JdbcNameSpace),
		beego.NSInclude(&jdbcController)))

	tableController := controllers.MetadataTableController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, controllers.MetaDataNameSpace),
		beego.NSInclude(&tableController)))

	fieldController := controllers.FieldsMetaController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, controllers.MetaDataNameSpace),
		beego.NSInclude(&fieldController)))

	sqlQueryController := controllers.SqlQueryController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, controllers.SqlQueryNameSpace),
		beego.NSInclude(&sqlQueryController)))

	RedashController := controllers.RedashController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, controllers.RedashNameSpace),
		beego.NSInclude(&RedashController)))

	AirflowController := controllers.AirflowController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, controllers.AirflowNameSpace),
		beego.NSInclude(&AirflowController)))

	ConnMetaController := controllers.ConnManagerController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, controllers.ConnNameSpace),
		beego.NSInclude(&ConnMetaController)))

	DagController := controllers.DagController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, controllers.DagNameSpace),
		beego.NSInclude(&DagController)))

	MobileDashboardController := controllers.MobileDashboardController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, controllers.MobileDashboardNameSpace),
		beego.NSInclude(&MobileDashboardController)))

	UdfController := controllers.UdfController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, controllers.UdfNameSpace),
		beego.NSInclude(&UdfController)))

	IndicatorController := indicator.IndicatorController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, indicator.IndicatorNameSpace),
		beego.NSInclude(&IndicatorController)))

	DimensionController := indicator.DimensionController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, indicator.IndicatorNameSpace),
		beego.NSInclude(&DimensionController)))

	BodongController := controllers.BodongController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, controllers.BodongNameSpace),
		beego.NSInclude(&BodongController)))

	DataRequirementController := controllers.DataRequirementController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, controllers.DataRequirementSpace),
		beego.NSInclude(&DataRequirementController)))

	SingleTableController := controllers.SingleTableController{}
	beego.AddNamespace(beego.NewNamespace(utils.GenNameSpace(rootNameSpace, controllers.SingleNameSpace),
		beego.NSInclude(&SingleTableController)))
}

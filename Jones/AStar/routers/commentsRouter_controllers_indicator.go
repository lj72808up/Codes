package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "DelIndicator",
            Router: `/delIndicator/:id`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "CheckIndicatorDuplicate",
            Router: `/dictionary/checkIndicatorDuplicate`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "GetDictionary",
            Router: `/dictionary/getAll`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "SearchDictionary",
            Router: `/dictionary/search`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "GetDictionaryById",
            Router: `/dictionary/searchById`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "DelDimFactRelation",
            Router: `/dimension/delDimFactRelation`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "GetDimensions",
            Router: `/dimension/getAll`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "GetAllMainDimension",
            Router: `/dimension/getAllMainDimension`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "GetDimFactRelations",
            Router: `/dimension/getDimFactRelations/:dimId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "PostDimFactRelation",
            Router: `/dimension/postDimFactRelation`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "NewDimensions",
            Router: `/dimension/postNew`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "EditCombineIndicator",
            Router: `/editCombineIndicator`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "EditPreIndicator",
            Router: `/editPreIndicator`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "GetAllCombineIndicator",
            Router: `/getAllCombineIndicator`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "GetIndicator",
            Router: `/getAllIndicator`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "GetEdit",
            Router: `/getEdit`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "GetIndicatorDependenciesById",
            Router: `/getIndicatorDependenciesById/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "GetIndicatorFactRelation",
            Router: `/getIndicatorFactRelation/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "PostCombineIndicator",
            Router: `/postCombineIndicator`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "PostPreIndicator",
            Router: `/postPreIndicator`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "AddIndicatorFactRelation",
            Router: `/predefine/addIndicatorFactRelation`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "DelIndicatorFactRelation",
            Router: `/predefine/delIndicatorFactRelation/:id`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "AddQualifier",
            Router: `/qualifier/add`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "GetAllTimeQualifier",
            Router: `/qualifier/time/get`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "EditServiceTag",
            Router: `/serviceTag/editServiceTag`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "GetServiceTag",
            Router: `/serviceTag/getAll`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"] = append(beego.GlobalControllerRouter["AStar/controllers/indicator:IndicatorController"],
        beego.ControllerComments{
            Method: "PutServiceTag",
            Router: `/serviceTag/putServiceTag`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["AStar/controllers:AirflowController"] = append(beego.GlobalControllerRouter["AStar/controllers:AirflowController"],
        beego.ControllerComments{
            Method: "GetConnByPrivilege",
            Router: `/GetConnByPrivilege/:connType`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AirflowController"] = append(beego.GlobalControllerRouter["AStar/controllers:AirflowController"],
        beego.ControllerComments{
            Method: "GetAllAirflow",
            Router: `/airflowall`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AirflowController"] = append(beego.GlobalControllerRouter["AStar/controllers:AirflowController"],
        beego.ControllerComments{
            Method: "DeleteDag",
            Router: `/deletedag`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AirflowController"] = append(beego.GlobalControllerRouter["AStar/controllers:AirflowController"],
        beego.ControllerComments{
            Method: "MarkTaskState",
            Router: `/markTaskState`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AirflowController"] = append(beego.GlobalControllerRouter["AStar/controllers:AirflowController"],
        beego.ControllerComments{
            Method: "MarkDagState",
            Router: `/markdagstate`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AirflowController"] = append(beego.GlobalControllerRouter["AStar/controllers:AirflowController"],
        beego.ControllerComments{
            Method: "SwitchDag",
            Router: `/switchdag`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AirflowController"] = append(beego.GlobalControllerRouter["AStar/controllers:AirflowController"],
        beego.ControllerComments{
            Method: "GetAllTemplate",
            Router: `/template/getAllTemplate/:tasktype`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AirflowController"] = append(beego.GlobalControllerRouter["AStar/controllers:AirflowController"],
        beego.ControllerComments{
            Method: "GetJdbcTemplate",
            Router: `/template/getById/:templateId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "GetAllBodongTree",
            Router: `/auth/bodong/GetAllBodongTree`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "GetCheckedBodongNode",
            Router: `/auth/bodong/GetCheckedBodongNode`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "UpdateCheckedBodongNode",
            Router: `/auth/bodong/UpdateCheckedBodongNode`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "GetBodongTreeAccess",
            Router: `/auth/bodong/getTreeAccess/:role`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "BatchChangeWorkSheetMapping",
            Router: `/auth/jdbc/batchChangeWorkSheetMapping`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "ChangeWorkSheetMapping",
            Router: `/auth/jdbc/changeWorkSheetMapping/:roleName`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "UpdateSortId",
            Router: `/auth/label/updateSortId`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "PostAdministrator",
            Router: `/auth/lvpiDomain/postAdministrator`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "PostFirstDomain",
            Router: `/auth/lvpiDomain/postFirstDomain`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "GetAllMetaTableTreeAccess",
            Router: `/auth/metadata/GetAllMetaTableTreeAccess/:role`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "GetMetaTableAccess",
            Router: `/auth/metadata/GetMetaTableAccess/:role`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "ChangeMetaTableMapping",
            Router: `/auth/metadata/changeMetaTableMapping`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "GetAllMetaTableAccess",
            Router: `/auth/metadata/getAllMetaTableAccess/:role`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "GetRoles",
            Router: `/auth/roles`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "DelRole",
            Router: `/auth/roles/:rolename`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "PostNewRole",
            Router: `/auth/roles/addRole`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "AddReadOnlyNameSpace",
            Router: `/auth/roles/namespace/addReadOnly`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "AddRouteForRole",
            Router: `/auth/roles/route/:rolename`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "DelRouteForRole",
            Router: `/auth/roles/route/:rolename`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "AddUserForRole",
            Router: `/auth/roles/user/:rolename`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "DelRoleByUser",
            Router: `/auth/roles/user/:rolename`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "GetRolesByUser",
            Router: `/auth/roles/user/:username`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "GetRoutes",
            Router: `/auth/routes`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "GetRoutesForRole",
            Router: `/auth/routes/role/:rolename`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "GetAllWorksheetCategory",
            Router: `/auth/template/getAllWorksheetCategory/:role`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "PostOrUpdateUser",
            Router: `/auth/user/postUser`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "AddRolesForUser",
            Router: `/auth/users/:username`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "GetJdbcCheck",
            Router: `/auth/users/jdbcRole/:username`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:AuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:AuthController"],
        beego.ControllerComments{
            Method: "GetUsersByRole",
            Router: `/auth/users/role/:rolename`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:BodongController"] = append(beego.GlobalControllerRouter["AStar/controllers:BodongController"],
        beego.ControllerComments{
            Method: "GetChildOptions",
            Router: `/getChildOptions/:mid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:BodongController"] = append(beego.GlobalControllerRouter["AStar/controllers:BodongController"],
        beego.ControllerComments{
            Method: "GetFirstLevel",
            Router: `/getFirst`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"] = append(beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"],
        beego.ControllerComments{
            Method: "AddConn",
            Router: `/addConn`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"] = append(beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/get`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"] = append(beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"],
        beego.ControllerComments{
            Method: "GetById",
            Router: `/getById/:connId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"] = append(beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"],
        beego.ControllerComments{
            Method: "GetByOwner",
            Router: `/getByOwner`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"] = append(beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"],
        beego.ControllerComments{
            Method: "GetEdit",
            Router: `/getEdit`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"] = append(beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"],
        beego.ControllerComments{
            Method: "GetPrivilegeRolesById",
            Router: `/getPrivilegeRolesById/:dsId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"] = append(beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"],
        beego.ControllerComments{
            Method: "TestConn",
            Router: `/testConn`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"] = append(beego.GlobalControllerRouter["AStar/controllers:ConnManagerController"],
        beego.ControllerComments{
            Method: "UpdateConn",
            Router: `/updateConn`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DWController"] = append(beego.GlobalControllerRouter["AStar/controllers:DWController"],
        beego.ControllerComments{
            Method: "GetExploreData",
            Router: `/dw/explore`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DWController"] = append(beego.GlobalControllerRouter["AStar/controllers:DWController"],
        beego.ControllerComments{
            Method: "Meta",
            Router: `/dw/meta/get`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DWController"] = append(beego.GlobalControllerRouter["AStar/controllers:DWController"],
        beego.ControllerComments{
            Method: "MetaAll",
            Router: `/dw/metaAll`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DWController"] = append(beego.GlobalControllerRouter["AStar/controllers:DWController"],
        beego.ControllerComments{
            Method: "SearchField",
            Router: `/dw/searchField`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DWController"] = append(beego.GlobalControllerRouter["AStar/controllers:DWController"],
        beego.ControllerComments{
            Method: "QuerySQLTask",
            Router: `/dw/sql/custom`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DWController"] = append(beego.GlobalControllerRouter["AStar/controllers:DWController"],
        beego.ControllerComments{
            Method: "QueryID",
            Router: `/dw/sql/id`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DWController"] = append(beego.GlobalControllerRouter["AStar/controllers:DWController"],
        beego.ControllerComments{
            Method: "GetOptions",
            Router: `/dw/sql/options`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DWController"] = append(beego.GlobalControllerRouter["AStar/controllers:DWController"],
        beego.ControllerComments{
            Method: "GetOptionsPCDate",
            Router: `/dw/sql/options_pc/date`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DWController"] = append(beego.GlobalControllerRouter["AStar/controllers:DWController"],
        beego.ControllerComments{
            Method: "GetOptionsDate",
            Router: `/dw/sql/options_wap/date`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DWController"] = append(beego.GlobalControllerRouter["AStar/controllers:DWController"],
        beego.ControllerComments{
            Method: "GetOptionsYHDate",
            Router: `/dw/sql/options_yh/date`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DWController"] = append(beego.GlobalControllerRouter["AStar/controllers:DWController"],
        beego.ControllerComments{
            Method: "QuerySQL",
            Router: `/dw/sql/task`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DWController"] = append(beego.GlobalControllerRouter["AStar/controllers:DWController"],
        beego.ControllerComments{
            Method: "SqlVerify",
            Router: `/dw/sqlVerify`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DWController"] = append(beego.GlobalControllerRouter["AStar/controllers:DWController"],
        beego.ControllerComments{
            Method: "Upload",
            Router: `/dw/upload/file`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "PostDag",
            Router: `/addDag`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "AfterPostDag",
            Router: `/afterPostDag`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "AssignDag",
            Router: `/assignDag`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "CheckDagName",
            Router: `/checkDagName/:dagId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "CheckOutTable",
            Router: `/checkOutTable`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "GetAllDagImgs",
            Router: `/getAllDagImgs`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "GetDagByName",
            Router: `/getDag`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "GetTaskDetail",
            Router: `/getDetailByDagId`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "GetHiveStorage",
            Router: `/getHiveStorage`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "GetManualTriggerHistory",
            Router: `/getManualTriggerHistory`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "GetOutTables",
            Router: `/getOutTables`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "GetSqlFields",
            Router: `/getSqlFields`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "GetTaskCntInOneDag",
            Router: `/getTaskCntInOneDag/:dagId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "GetTriggerLog",
            Router: `/getTriggerLog`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "ManualTrigger",
            Router: `/manualTrigger`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "PrePostDag",
            Router: `/preAddDag`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "SaveDagWhenPostFail",
            Router: `/saveDagWhenPostFail`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DagController"] = append(beego.GlobalControllerRouter["AStar/controllers:DagController"],
        beego.ControllerComments{
            Method: "UpdateDag",
            Router: `/updateDag`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "GetDevStatusMapping",
            Router: `/GetDevStatusMapping`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "DelAttachmentById",
            Router: `/attachment/delete/:aid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "DownloadAttachment",
            Router: `/attachment/download/:aid`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "GetAttachmentList",
            Router: `/attachment/get/:rid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "UploadAttachment",
            Router: `/attachment/upload/:rid`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "GetAllDemand",
            Router: `/getAllDemand`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "GetAssetRelationById",
            Router: `/getAssetRelationById/:rid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "GetDemandById",
            Router: `/getDemandById/:rid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "GetHistory",
            Router: `/getHistory/:rid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "AddDataRequirement",
            Router: `/indicator/addDataRequirement`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "GetDataRequirement",
            Router: `/indicator/getDataRequirement`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "GetUserDefineIndicator",
            Router: `/indicator/userDefine/get/:rid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "PostUserDefineIndicator",
            Router: `/indicator/userDefine/post/:rid`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/post`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "PostAttachment",
            Router: `/postAttachment/:requestId`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "PostRequireAssetRelation",
            Router: `/postRequireAssetRelation/:requestId`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"] = append(beego.GlobalControllerRouter["AStar/controllers:DataRequirementController"],
        beego.ControllerComments{
            Method: "Update",
            Router: `/update/:rid`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DocController"] = append(beego.GlobalControllerRouter["AStar/controllers:DocController"],
        beego.ControllerComments{
            Method: "GetDoc",
            Router: `/doc/get`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DocController"] = append(beego.GlobalControllerRouter["AStar/controllers:DocController"],
        beego.ControllerComments{
            Method: "GetDocHtml",
            Router: `/doc/gethtml`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:DocController"] = append(beego.GlobalControllerRouter["AStar/controllers:DocController"],
        beego.ControllerComments{
            Method: "SubmitDoc",
            Router: `/doc/submit`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "UpdateFields",
            Router: `/field/update`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "UpdatePartitionField",
            Router: `/field/updatePartition`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "DeleteField",
            Router: `/fields/delete`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "GetFieldsMetaDataByVersion",
            Router: `/fields/getByVersion`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "GetFieldsListByVersion",
            Router: `/fields/getListByVersion`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "GetPartitionFieldsByVersion",
            Router: `/fields/getPartitionByVersion`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "PutFieldsMetaData",
            Router: `/fields/newVersion/put`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "DeleteFieldsRelation",
            Router: `/fields/relation/delete`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "GetFieldsRelation",
            Router: `/fields/relation/get/:fieldId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "PutFieldsRelation",
            Router: `/fields/relation/update`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "UpdateFieldsByVersion",
            Router: `/fields/updateByVersion`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "GetFieldAllVersion",
            Router: `/fields/version/get/:tableId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "GetDefaultNewVersion",
            Router: `/fields/version/getNewVersion/:tableId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "GetDataBaseByType",
            Router: `/getDataBaseByType/:tableType`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "GetTablesByDataBase",
            Router: `/getTablesByDataBase`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "Graph",
            Router: `/graph/:tableId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "CheckTablePairRelation",
            Router: `/graph/relation/check`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "GenerateFieldAuto",
            Router: `/import`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "GenerateFieldAutoSingle",
            Router: `/import/single`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "GetTableTypes",
            Router: `/tableTypes/get`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "UploadFieldRelation",
            Router: `/upload/relation/file`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"] = append(beego.GlobalControllerRouter["AStar/controllers:FieldsMetaController"],
        beego.ControllerComments{
            Method: "UpdateVersionDesc",
            Router: `/version/update`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "LvpiCreateQueryView",
            Router: `/sql/createLvpiView`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "SQLCreateQueryView",
            Router: `/sql/createQueryView`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "ChangeTemplateAffiliation",
            Router: `/template/ChangeTemplateAffiliation`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "GetCanModify",
            Router: `/template/canModify/:templateId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "CreateQueryView",
            Router: `/template/createQueryView`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "DeleteTemplateById",
            Router: `/template/deleteTemplateById`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "GetAllDataSources",
            Router: `/template/getAllDataSources`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "GetAllTemplate",
            Router: `/template/getAllTemplate`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "GetJdbcTemplate",
            Router: `/template/getById/:templateId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "GetChildNodes",
            Router: `/template/getChildNodes/:templateId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "GetEdit",
            Router: `/template/getEdit`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "GetJumpAccess",
            Router: `/template/getjumpaccess`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "GetViewAccess",
            Router: `/template/getviewaccess`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "PutJdbcTemplate",
            Router: `/template/put`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "RunJdbcTemplate",
            Router: `/template/run`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:JdbcController"] = append(beego.GlobalControllerRouter["AStar/controllers:JdbcController"],
        beego.ControllerComments{
            Method: "UpdateJdbcTemplate",
            Router: `/template/updateJdbcTemplate`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"],
        beego.ControllerComments{
            Method: "GetAllDomain",
            Router: `/getAllDomain`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"],
        beego.ControllerComments{
            Method: "GetDomainBySelfAdmin",
            Router: `/getDomainBySelfAdmin`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"],
        beego.ControllerComments{
            Method: "GetFirstDomainLvpiVisible",
            Router: `/getFirstDomainLvpiVisible`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"],
        beego.ControllerComments{
            Method: "PostLvpiMapping",
            Router: `/getFirstDomainLvpiVisible/uid`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"],
        beego.ControllerComments{
            Method: "GetLvpiBySuperior",
            Router: `/getLvpiBySuperior/:pid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"],
        beego.ControllerComments{
            Method: "GetLvpiInDomain",
            Router: `/getLvpiInDomain/:did`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"],
        beego.ControllerComments{
            Method: "GetLvpis",
            Router: `/getLvpis/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"],
        beego.ControllerComments{
            Method: "GetlvpiTree",
            Router: `/getlvpiTree`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"],
        beego.ControllerComments{
            Method: "PostSecondLvpi",
            Router: `/postSecondLvpi`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"] = append(beego.GlobalControllerRouter["AStar/controllers:LvpiAuthController"],
        beego.ControllerComments{
            Method: "PostSecondLvpiUpdate",
            Router: `/postSecondLvpiUpdate`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"],
        beego.ControllerComments{
            Method: "AddTableMetaData",
            Router: `/addtable`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"],
        beego.ControllerComments{
            Method: "GetEdit",
            Router: `/getEdit`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"],
        beego.ControllerComments{
            Method: "GetOwnerEdit",
            Router: `/getOwnerEdit`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"],
        beego.ControllerComments{
            Method: "GetTableByConditions",
            Router: `/getTableByConditions`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"],
        beego.ControllerComments{
            Method: "GetTablesById",
            Router: `/getTableById/:tid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"],
        beego.ControllerComments{
            Method: "GetTablePrivilegeByTid",
            Router: `/getTablePrivilegeByTid/:tid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"],
        beego.ControllerComments{
            Method: "GetAllDBMetaData",
            Router: `/getalldb`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"],
        beego.ControllerComments{
            Method: "GetAllTableMetaData",
            Router: `/getalltable`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"],
        beego.ControllerComments{
            Method: "GetTablesOfDb",
            Router: `/getdbtable/:dbname`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"],
        beego.ControllerComments{
            Method: "PostTables",
            Router: `/postTables`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"],
        beego.ControllerComments{
            Method: "DeleteTableById",
            Router: `/table/delete/:tid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"],
        beego.ControllerComments{
            Method: "GetTableNumber",
            Router: `/table/getNums`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"],
        beego.ControllerComments{
            Method: "UpdateTableById",
            Router: `/table/update`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"],
        beego.ControllerComments{
            Method: "GetConnectionIds",
            Router: `/tableImport/getConnections`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:MetadataTableController"],
        beego.ControllerComments{
            Method: "UpdateTableByVersion",
            Router: `/updateTableByVersion`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:MobileDashboardController"] = append(beego.GlobalControllerRouter["AStar/controllers:MobileDashboardController"],
        beego.ControllerComments{
            Method: "GetMobileDashboardUrl",
            Router: `/getDashboardUrl/:labelAbbreviation`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:QueryPcController"] = append(beego.GlobalControllerRouter["AStar/controllers:QueryPcController"],
        beego.ControllerComments{
            Method: "QuerySubmit",
            Router: `/jone/combination_pc`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:QueryWsController"] = append(beego.GlobalControllerRouter["AStar/controllers:QueryWsController"],
        beego.ControllerComments{
            Method: "QuerySubmit",
            Router: `/jone/combination`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:QueryYhController"] = append(beego.GlobalControllerRouter["AStar/controllers:QueryYhController"],
        beego.ControllerComments{
            Method: "QuerySubmit",
            Router: `/jone/combination_yh`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:RedashController"] = append(beego.GlobalControllerRouter["AStar/controllers:RedashController"],
        beego.ControllerComments{
            Method: "GetAllDashboard",
            Router: `/dashboardall`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:RedashController"] = append(beego.GlobalControllerRouter["AStar/controllers:RedashController"],
        beego.ControllerComments{
            Method: "DeleteDashboard",
            Router: `/dashboarddel`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:RedashController"] = append(beego.GlobalControllerRouter["AStar/controllers:RedashController"],
        beego.ControllerComments{
            Method: "InsertDashboard",
            Router: `/dashboardinsert`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:RedashController"] = append(beego.GlobalControllerRouter["AStar/controllers:RedashController"],
        beego.ControllerComments{
            Method: "UpdateDashboard",
            Router: `/dashboardupdate`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:RedashController"] = append(beego.GlobalControllerRouter["AStar/controllers:RedashController"],
        beego.ControllerComments{
            Method: "GetAllRedash",
            Router: `/redashall`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:RedashController"] = append(beego.GlobalControllerRouter["AStar/controllers:RedashController"],
        beego.ControllerComments{
            Method: "CheckRedash",
            Router: `/redashcheck`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:RedashController"] = append(beego.GlobalControllerRouter["AStar/controllers:RedashController"],
        beego.ControllerComments{
            Method: "DeleteRedash",
            Router: `/redashdel`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:RedashController"] = append(beego.GlobalControllerRouter["AStar/controllers:RedashController"],
        beego.ControllerComments{
            Method: "InsertRedash",
            Router: `/redashinsert`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:RedashController"] = append(beego.GlobalControllerRouter["AStar/controllers:RedashController"],
        beego.ControllerComments{
            Method: "GetAllRedahQuery",
            Router: `/redashqueryall`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:RedashController"] = append(beego.GlobalControllerRouter["AStar/controllers:RedashController"],
        beego.ControllerComments{
            Method: "DeleteRedashQuery",
            Router: `/redashquerydel`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:RedashController"] = append(beego.GlobalControllerRouter["AStar/controllers:RedashController"],
        beego.ControllerComments{
            Method: "GetTimeStamp",
            Router: `/timestamp`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:RedashController"] = append(beego.GlobalControllerRouter["AStar/controllers:RedashController"],
        beego.ControllerComments{
            Method: "ChangeUser",
            Router: `/user_change`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:RedashController"] = append(beego.GlobalControllerRouter["AStar/controllers:RedashController"],
        beego.ControllerComments{
            Method: "CheckWriteAccess",
            Router: `/write_access`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SchedulerJobController"] = append(beego.GlobalControllerRouter["AStar/controllers:SchedulerJobController"],
        beego.ControllerComments{
            Method: "GetSchedulerJob",
            Router: `/scheduler_job/getById/:jobId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SchedulerJobController"] = append(beego.GlobalControllerRouter["AStar/controllers:SchedulerJobController"],
        beego.ControllerComments{
            Method: "PutSchedulerJob",
            Router: `/scheduler_job/put`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SchedulerJobController"] = append(beego.GlobalControllerRouter["AStar/controllers:SchedulerJobController"],
        beego.ControllerComments{
            Method: "StartSchedulerJob",
            Router: `/scheduler_job/startById`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SchedulerJobController"] = append(beego.GlobalControllerRouter["AStar/controllers:SchedulerJobController"],
        beego.ControllerComments{
            Method: "StopSchedulerJob",
            Router: `/scheduler_job/stopJobById`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SingleTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:SingleTableController"],
        beego.ControllerComments{
            Method: "AddDim",
            Router: `/addDim`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SingleTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:SingleTableController"],
        beego.ControllerComments{
            Method: "DelDim",
            Router: `/delDim/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SingleTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:SingleTableController"],
        beego.ControllerComments{
            Method: "GetDims",
            Router: `/getDims`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SingleTableController"] = append(beego.GlobalControllerRouter["AStar/controllers:SingleTableController"],
        beego.ControllerComments{
            Method: "UpdateDim",
            Router: `/updateDim`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"] = append(beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"],
        beego.ControllerComments{
            Method: "Download",
            Router: `/small/download`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"] = append(beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"],
        beego.ControllerComments{
            Method: "GetIpPc",
            Router: `/small/ip/pc`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"] = append(beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"],
        beego.ControllerComments{
            Method: "GetIpWap",
            Router: `/small/ip/wap`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"] = append(beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"],
        beego.ControllerComments{
            Method: "GetAllDisable",
            Router: `/small/log/disable`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"] = append(beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"],
        beego.ControllerComments{
            Method: "GetAllEnable",
            Router: `/small/log/enable`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"] = append(beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"],
        beego.ControllerComments{
            Method: "GetJsonByTaskId",
            Router: `/small/log/query/:task_id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"] = append(beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"],
        beego.ControllerComments{
            Method: "GetAllReady",
            Router: `/small/log/ready`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"] = append(beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"],
        beego.ControllerComments{
            Method: "Query",
            Router: `/small/query`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"] = append(beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"],
        beego.ControllerComments{
            Method: "Submit2Pc",
            Router: `/small/submit2pc`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"] = append(beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"],
        beego.ControllerComments{
            Method: "Submit2Ws",
            Router: `/small/submit2ws`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"] = append(beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"],
        beego.ControllerComments{
            Method: "TaskOffline",
            Router: `/small/task/offline/:task_id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"] = append(beego.GlobalControllerRouter["AStar/controllers:SmallFlowController"],
        beego.ControllerComments{
            Method: "TaskOnline",
            Router: `/small/task/online/:task_id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SqlQueryController"] = append(beego.GlobalControllerRouter["AStar/controllers:SqlQueryController"],
        beego.ControllerComments{
            Method: "GetConnByPrivilege",
            Router: `/GetConnByPrivilege/:connType`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SqlQueryController"] = append(beego.GlobalControllerRouter["AStar/controllers:SqlQueryController"],
        beego.ControllerComments{
            Method: "GetSqlHint",
            Router: `/hint/getSqlHint`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:SqlQueryController"] = append(beego.GlobalControllerRouter["AStar/controllers:SqlQueryController"],
        beego.ControllerComments{
            Method: "GetSqlUDF",
            Router: `/hint/getSqlUDF`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "ChangeTaskUser",
            Router: `/task/changeTaskUser`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "CheckException",
            Router: `/task/checkException`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "Limit",
            Router: `/task/data/path`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "DeleteTaskRecoder",
            Router: `/task/deleteTaskRecoder/:taskId`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "File",
            Router: `/task/file`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "KillByJobId",
            Router: `/task/killByJobId`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/task/log`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "GetByConfig",
            Router: `/task/log/config`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "GetByDate",
            Router: `/task/log/date/:date`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "GetJsonById",
            Router: `/task/log/query/:query_id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "GetRsync",
            Router: `/task/log/rsync`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "GetByUser",
            Router: `/task/log/user/:username`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "GetSubmitQueryLog",
            Router: `/task/queryLog/:queryId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "ReRun",
            Router: `/task/reRun`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "GetByConfigGod",
            Router: `/task_god/log/config`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TaskController"] = append(beego.GlobalControllerRouter["AStar/controllers:TaskController"],
        beego.ControllerComments{
            Method: "GetTaskById",
            Router: `/task_god/query/:queryId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TestController"] = append(beego.GlobalControllerRouter["AStar/controllers:TestController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/testGet/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TestController"] = append(beego.GlobalControllerRouter["AStar/controllers:TestController"],
        beego.ControllerComments{
            Method: "GetCanWrite",
            Router: `/testGetPermission`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TestController"] = append(beego.GlobalControllerRouter["AStar/controllers:TestController"],
        beego.ControllerComments{
            Method: "PostAll",
            Router: `/testPost`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TopLabelController"] = append(beego.GlobalControllerRouter["AStar/controllers:TopLabelController"],
        beego.ControllerComments{
            Method: "GetLabels",
            Router: `/auth/labels/:groupName`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TopLabelController"] = append(beego.GlobalControllerRouter["AStar/controllers:TopLabelController"],
        beego.ControllerComments{
            Method: "AddNewLabels",
            Router: `/auth/labels/addNewLabels`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TopLabelController"] = append(beego.GlobalControllerRouter["AStar/controllers:TopLabelController"],
        beego.ControllerComments{
            Method: "GetCheckedLabels",
            Router: `/auth/labels/checked/:username`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TopLabelController"] = append(beego.GlobalControllerRouter["AStar/controllers:TopLabelController"],
        beego.ControllerComments{
            Method: "GetPMLabels",
            Router: `/auth/labels/pmlabels`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TopLabelController"] = append(beego.GlobalControllerRouter["AStar/controllers:TopLabelController"],
        beego.ControllerComments{
            Method: "UpdatePMLabels",
            Router: `/auth/labels/pmupdate`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TopLabelController"] = append(beego.GlobalControllerRouter["AStar/controllers:TopLabelController"],
        beego.ControllerComments{
            Method: "AddRDLabels",
            Router: `/auth/labels/rdadd`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TopLabelController"] = append(beego.GlobalControllerRouter["AStar/controllers:TopLabelController"],
        beego.ControllerComments{
            Method: "GetRDLabels",
            Router: `/auth/labels/rdlabels`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TopLabelController"] = append(beego.GlobalControllerRouter["AStar/controllers:TopLabelController"],
        beego.ControllerComments{
            Method: "UpdateRDLabels",
            Router: `/auth/labels/rdupdate`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:TopLabelController"] = append(beego.GlobalControllerRouter["AStar/controllers:TopLabelController"],
        beego.ControllerComments{
            Method: "UpdateLabels",
            Router: `/auth/labels/update`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:UdfController"] = append(beego.GlobalControllerRouter["AStar/controllers:UdfController"],
        beego.ControllerComments{
            Method: "DelById",
            Router: `/delete/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:UdfController"] = append(beego.GlobalControllerRouter["AStar/controllers:UdfController"],
        beego.ControllerComments{
            Method: "GetEdit",
            Router: `/getEdit`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:UdfController"] = append(beego.GlobalControllerRouter["AStar/controllers:UdfController"],
        beego.ControllerComments{
            Method: "GetUdfList",
            Router: `/getList`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:UdfController"] = append(beego.GlobalControllerRouter["AStar/controllers:UdfController"],
        beego.ControllerComments{
            Method: "UpdateFunctionDesc",
            Router: `/update/description`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:XiaopController"] = append(beego.GlobalControllerRouter["AStar/controllers:XiaopController"],
        beego.ControllerComments{
            Method: "Apply",
            Router: `/apply/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:XiaopController"] = append(beego.GlobalControllerRouter["AStar/controllers:XiaopController"],
        beego.ControllerComments{
            Method: "GetApplyStatus",
            Router: `/applyStatus/:uid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:XiaopController"] = append(beego.GlobalControllerRouter["AStar/controllers:XiaopController"],
        beego.ControllerComments{
            Method: "PostOrUpdateUser",
            Router: `/auth/user/postUser`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:XiaopController"] = append(beego.GlobalControllerRouter["AStar/controllers:XiaopController"],
        beego.ControllerComments{
            Method: "ParseToken",
            Router: `/jone/parseToken`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:XiaopController"] = append(beego.GlobalControllerRouter["AStar/controllers:XiaopController"],
        beego.ControllerComments{
            Method: "ParseLvpiToken",
            Router: `/lvpi/parseToken`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:XiaopController"] = append(beego.GlobalControllerRouter["AStar/controllers:XiaopController"],
        beego.ControllerComments{
            Method: "Robot",
            Router: `/robot/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:XiaopController"] = append(beego.GlobalControllerRouter["AStar/controllers:XiaopController"],
        beego.ControllerComments{
            Method: "GetRole",
            Router: `/role/:uid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:XiaopController"] = append(beego.GlobalControllerRouter["AStar/controllers:XiaopController"],
        beego.ControllerComments{
            Method: "AddRole",
            Router: `/role/add`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["AStar/controllers:XiaopController"] = append(beego.GlobalControllerRouter["AStar/controllers:XiaopController"],
        beego.ControllerComments{
            Method: "Roles",
            Router: `/roles/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}

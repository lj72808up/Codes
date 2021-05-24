package controllers

import (
	"AStar/models"
	"encoding/json"
	"strconv"
)

type SqlQueryController struct {
	BaseController
}

var SqlQueryNameSpace = "dw"

// @Title 获取sql输入的提示消息
// @Summary 获取sql输入的提示消息
// @Description 获取sql输入的提示消息
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /hint/getSqlHint [post]
func (this *SqlQueryController) GetSqlHint() {
	sqlQuery := models.SqlQueryModel{}
	var params map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	dsType := params["dsType"]
	dsId, _ := strconv.ParseInt(params["dsId"], 10, 64)

	sqlQuery.GetSqlHint(dsType, dsId)
	res, _ := json.Marshal(sqlQuery.SqlHint)
	this.Ctx.WriteString(string(res))
}

// @Title 获取所有UDF函数
// @Summary 获取所有UDF函数
// @Description 获取所有UDF函数
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /hint/getSqlUDF [get]
func (this *SqlQueryController) GetSqlUDF() {
	sqlQuery := models.SqlQueryModel{}
	sqlQuery.GetSqlUDF()
	res, _ := json.Marshal(sqlQuery.FunctionUDFs)
	this.Ctx.WriteString(string(res))
}

// @Title 获取自身角色可以看到的连接
// @Summary 获取自身角色可以看到的连接
// @Description 获取自身角色可以看到的连接
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /GetConnByPrivilege/:connType [get]
func (this *SqlQueryController) GetConnByPrivilege() {
	connType := this.Ctx.Input.Params()[":connType"]
	roles := models.Enforcer.GetRolesForUser(this.userName)
	conns := models.GetConnByPrivilege(connType, roles)
	res, _ := json.Marshal(conns)
	this.Ctx.WriteString(string(res))
}

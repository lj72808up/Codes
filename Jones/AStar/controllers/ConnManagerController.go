package controllers

import (
	"AStar/models"
	"AStar/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type ConnManagerController struct {
	BaseController
}

var ConnNameSpace = "connMeta"

// @Title 获取所有连接
// @Summary 获取所有连接
// @Description 获取所有连接
// @Success 200
// @Accept json
// @router /get [get]
func (this *ConnManagerController) Get() {
	connType := this.GetString("connType") // 连接类型
	roles := models.Enforcer.GetRolesForUser(this.userName)
	conns := models.GetConnByOwner(connType, roles)
	res, _ := json.Marshal(conns)
	this.Ctx.WriteString(string(res))
}

// @Title 获取owner下的连接
// @Summary 获取owner下的连接
// @Description 获取owner下的连接
// @Success 200
// @Accept json
// @router /getByOwner [get]
func (this *ConnManagerController) GetByOwner() {
	connType := this.GetString("connType") // 连接类型
	roles := models.Enforcer.GetRolesForUser(this.userName)
	conns := models.GetConnOnlyByOwner(connType, roles)
	res, _ := json.Marshal(conns)
	this.Ctx.WriteString(string(res))
}

// @Title 根据Id获取连接
// @Summary 根据Id获取连接
// @Description 根据Id获取连接
// @Success 200
// @Accept json
// @router /getById/:connId [get]
func (this *ConnManagerController) GetById() {
	connId, _ := strconv.ParseInt(this.Ctx.Input.Params()[":connId"], 10, 64)
	var conn models.ConnectionManager
	conn.GetById(connId)
	res, _ := json.Marshal(conn)
	this.Ctx.WriteString(string(res))
}

// @Title 测试连接
// @Summary 测试连接
// @Description 测试连接
// @Success 200
// @Accept json
// @router /testConn [post]
func (this *ConnManagerController) TestConn() {
	var params map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &params)
	var manager models.ConnectionManager
	connType := params["connType"]
	user := params["user"]
	passwd := params["passwd"]
	database := params["database"]
	host := params["host"]
	port := params["port"]
	// 连接失败捕获
	defer func() {
		if r := recover(); r != nil {
			//fmt.Println(r)
			response := utils.VueResponse{
				Res:  false,
				Info: fmt.Sprintf("%v", r),
			}
			this.Ctx.WriteString(response.String())
		}
	}()
	// 正常返回
	err := manager.TestConn(connType, user, passwd, database, host, port)
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	res, _ := json.Marshal(utils.VueResponse{
		Res:  err == nil,
		Info: msg,
	})
	this.Ctx.WriteString(string(res))
}

// @Title GetEdit
// @Summary 查看是否可以编辑连接管理
// @Description 查看是否可以编辑连接管理
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /getEdit [get]
func (this *ConnManagerController) GetEdit() {
	method := "POST"
	checkUrl := fmt.Sprintf("%s/%s_%s", RootNameSpace, ConnNameSpace, "authority")
	models.Logger.Info(checkUrl)
	flag := models.Enforcer.Enforce(this.userName, checkUrl, method)
	res, _ := json.Marshal(map[string]bool{
		"readOnly": !flag,
	})
	this.Ctx.WriteString(string(res))
}

// @Title 新增连接
// @Summary 新增连接
// @Description 新增连接
// @Success 200
// @Accept json
// @router /addConn [post]
func (this *ConnManagerController) AddConn() {
	var conn models.ConnectionManager
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &conn)
	defer func() {
		if r := recover(); r != nil {
			this.Ctx.WriteString(utils.VueResponse{
				Res:  false,
				Info: fmt.Sprintf("%v", r),
			}.String())
		}
	}()
	roles := models.Enforcer.GetRolesForUser(this.userName)         // 用户在casbin中的权限
	privilegeRoles := strings.Split(conn.PrivilegeRoles, ",")   // 页面都选的有连接权限的用户
	conn.Add(roles, privilegeRoles)
	this.Ctx.WriteString(utils.VueResponse{
		Res: true,
	}.String())
}

// @Title 修改连接
// @Summary 修改连接
// @Description 修改连接
// @Success 200
// @Accept json
// @router /updateConn [post]
func (this *ConnManagerController) UpdateConn() {
	var conn models.ConnectionManager
	roles := models.Enforcer.GetRolesForUser(this.userName)
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &conn)
	defer func() {
		if r := recover(); r != nil {
			this.Ctx.WriteString(utils.VueResponse{
				Res:  false,
				Info: fmt.Sprintf("%v", r),
			}.String())
		}
	}()
	conn.Owner = strings.Join(roles,",")
	conn.Update(roles)
	this.Ctx.WriteString(utils.VueResponse{
		Res: true,
	}.String())
}

// @Title 获取该链接有权访问的角色列表
// @Summary 获取该链接有权访问的角色列表
// @Description 获取该链接有权访问的角色列表
// @Success 200
// @Accept json
// @router /getPrivilegeRolesById/:dsId [get]
func (this *ConnManagerController) GetPrivilegeRolesById() {
	param := this.GetString(":dsId")
	dsId, _ := strconv.ParseInt(param, 10, 64)
	roles := models.GetPrivilegeRolesByConnId(dsId)
	res, _ := json.Marshal(roles)
	this.Ctx.WriteString(string(res))
}

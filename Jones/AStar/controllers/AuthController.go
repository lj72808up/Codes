package controllers

import (
	"AStar/models"
	"AStar/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/swagger"
	"io/ioutil"
	"strconv"
	"strings"
)

type AuthController struct {
	BaseController
}

// @Title GetRoutes
// @Summary 通过获取所有routes路径
// @Description 通过获取所有routes路径
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/routes [get]
func (this *AuthController) GetRoutes() {
	swagger := swagger.Swagger{}
	bytes, err := ioutil.ReadFile("swagger/swagger.json")
	if err != nil {
		panic("bee run -downdoc=true -gendoc=true")
	}
	json.Unmarshal(bytes, &swagger)
	var routes []string
	for k := range swagger.Paths {
		routes = append(routes, swagger.BasePath+"/"+k)
	}
	jsonData, _ := json.Marshal(routes)
	this.Ctx.WriteString(string(jsonData))
}

// @Title 通过获取角色的所有routes路径
// @Param	rolename path string true "角色名"
// @Summary 通过获取角色的所有routes路径
// @Description 通过获取角色的所有routes路径
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/routes/role/:rolename [get]
func (this *AuthController) GetRoutesForRole() {
	var routes []string
	rolename := this.Ctx.Input.Params()[":rolename"]
	a := models.Enforcer.GetPolicy()
	for _, resbin := range a {
		if resbin[0] == rolename {
			routes = append(routes, resbin[1])
		}
	}
	jsonData, _ := json.Marshal(routes)
	this.Ctx.WriteString(string(jsonData))
}

// @Title GetRole
// @Summary 通过获取所有角色路径
// @Description 通过获取所有routes路径
// @Success 200 {array[string]} []string
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/roles [get]
func (this *AuthController) GetRoles() {
	jsonStr, _ := json.Marshal(models.Enforcer.GetAllRoles())
	this.Ctx.WriteString(string(jsonStr))
}

// @Title 角色的用户添加
// @Param	rolename path string true "角色名"git
// @Param	username query string true "用户名"
// @Summary 角色的用户添加
// @Description 角色的用户添加
// @Success 200 {bool} bool
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/roles/user/:rolename [post]
func (this *AuthController) AddUserForRole() {
	rolename := this.Ctx.Input.Params()[":rolename"]
	username := this.GetString("username")
	checkDuplicate, _ := this.GetBool("checkDuplicate", false)
	if checkDuplicate {
		roles := models.Enforcer.GetRolesForUser(username)
		if len(roles) > 0 {
			rolesJson, _ := json.Marshal(roles)
			response := utils.VueResponse{
				Res:  false,
				Info: string(rolesJson),
			}
			this.Ctx.WriteString(response.String())
			return
		}
	}
	this.Data["json"] = models.Enforcer.AddRoleForUser(username, rolename)
	_ = models.Enforcer.LoadPolicy()
	//this.ServeJSON()
	response := utils.VueResponse{
		Res:  true,
		Info: "添加成功",
	}
	this.Ctx.WriteString(response.String())
}

// @Title 角色的许可路由添加
// @Param	rolename path string true "角色名"
// @Param	route query string true "许可路由"
// @Summary 角色的许可路由添加
// @Description 角色的许可路由添加
// @Success 200 {bool} bool
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/roles/route/:rolename [post]
func (this *AuthController) AddRouteForRole() {
	rolename := this.Ctx.Input.Params()[":rolename"]
	route := this.GetString("route")
	this.Data["json"] = models.Enforcer.AddPolicy(rolename, route, "*")
	models.Enforcer.LoadPolicy()
	this.ServeJSON()
}

// @Title 角色的许可路由删除
// @Param	rolename path string true "角色名"
// @Param	route query string true "许可路由"
// @Summary 角色的许可路由删除
// @Description 角色的许可路由删除
// @Success 200 {bool} bool
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /auth/roles/route/:rolename [delete]
func (this *AuthController) DelRouteForRole() {
	rolename := this.Ctx.Input.Params()[":rolename"]
	route := this.GetString("route")
	this.Data["json"] = models.Enforcer.RemovePolicy(rolename, route)
	models.Enforcer.LoadPolicy()
	this.ServeJSON()
}

// @Title GetRoleByUser
// @Param	username path string true "用于获取角色的用户名称"
// @Summary 获取该用户的所有角色
// @Description 获取该用户的所有角色
// @Success 200 {array[string]} []string
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/roles/user/:username [get]
func (this *AuthController) GetRolesByUser() {
	username := this.Ctx.Input.Params()[":username"]
	jsonStr, _ := json.Marshal(models.Enforcer.GetRolesForUser(username))
	this.Ctx.WriteString(string(jsonStr))
}

// @Title 删除角色
// @Param	rolename path string true "需删除的角色名"
// @Summary 删除角色
// @Description 删除角色
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /auth/roles/:rolename [delete]
func (this *AuthController) DelRole() {
	rolename := this.Ctx.Input.Params()[":rolename"]
	models.Enforcer.DeleteRole(rolename)
	models.Enforcer.LoadPolicy()
	this.ServeJSON()
}

// @Title 删除指定用户角色
// @Param	rolename path string true "需删除的角色名"
// @Param	username query string true "指定的用户名"
// @Summary 删除指定用户角色
// @Description 删除指定用户角色
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /auth/roles/user/:rolename [delete]
func (this *AuthController) DelRoleByUser() {
	rolename := this.Ctx.Input.Params()[":rolename"]
	username := this.GetString("username")
	models.Enforcer.DeleteRoleForUser(username, rolename)
	models.Enforcer.LoadPolicy()
	this.ServeJSON()
}

// @Title 为用户添加角色
// @Param	username path string true "用户名"
// @Param	rolename query string true "角色名"
// @Summary 为用户添加角色
// @Description 为用户添加角色
// @Success 200 {bool} bool
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/users/:username [post]
func (this *AuthController) AddRolesForUser() {
	username := this.Ctx.Input.Params()[":username"]
	rolename := this.GetString("rolename")
	this.Data["json"] = models.Enforcer.AddRoleForUser(username, rolename)
	models.Enforcer.LoadPolicy()
	this.ServeJSON()
}

// @Title 获取角色的所有用户
// @Param	rolename path string true "角色名"
// @Summary 获取角色的所有用户
// @Description 获取角色的所有用户
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/users/role/:rolename [get]
func (this *AuthController) GetUsersByRole() {
	rolename := this.Ctx.Input.Params()[":rolename"]
	jsonStr, _ := json.Marshal(models.Enforcer.GetUsersForRole(rolename))
	this.Ctx.WriteString(string(jsonStr))
}

// @Title 查看用户是否属于只能查看工单的用户组
// @Summary  查看用户是否属于只能查看工单的用户组
// @Description  查看用户是否属于只能查看工单的用户组
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/users/jdbcRole/:username [get]
func (this *AuthController) GetJdbcCheck() {
	userName := this.Ctx.Input.Params()[":username"]
	rolesForUser := models.Enforcer.GetRolesForUser(userName)
	flag := false
	for _, role := range rolesForUser {
		if role == "jdbcCheckGroup" {
			flag = true
		}
	}
	res, _ := json.Marshal(flag)
	this.Ctx.WriteString(string(res))
}

// @Title 更改用户和工单权限的配置
// @Summary  更改用户和工单权限的配置
// @Description  更改用户和工单权限的配置
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/jdbc/changeWorkSheetMapping/:roleName [post]
func (this *AuthController) ChangeWorkSheetMapping() {
	roleName := this.Ctx.Input.Params()[":roleName"]
	var checkedWorkSheets []models.SqlTemplateView
	// 前段参数中没有HasPrivilege属性, 数组中包含的都是用户选中的.
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &checkedWorkSheets)
	template := models.SqlTemplate{}
	template.AddworkSheetMapping(roleName, checkedWorkSheets) // 更新权限

}

// @Title 批量修改工单权限
// @Summary  批量修改工单权限
// @Description  批量修改工单权限
// @Accept json
// @router /auth/jdbc/batchChangeWorkSheetMapping [post]
func (this *AuthController) BatchChangeWorkSheetMapping() {
	attr := map[string]string{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	roles := strings.Split(attr["roles"],",")
	worksheetIds := strings.Split(attr["worksheetIds"],",")
	o := orm.NewOrm()

	for _, role := range roles {
		for _, wsId := range worksheetIds {
			wsIdInt,_ := strconv.ParseInt(wsId,10,64)
			_, _ = o.Raw("REPLACE INTO group_worksheet_mapping(group_name, worksheet_id, valid) VALUES(?,?,?)", role, wsIdInt, true).Exec()
		}
	}

}

// @Title 增加一个用户组(默认doc权限)
// @Summary  增加一个用户组(默认doc权限)
// @Description  增加一个用户组(默认doc权限)
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/roles/addRole [post]
func (this *AuthController) PostNewRole() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)

	role := attr["role"]
	if role == "" {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: "role name cannot empty",
		}.String())
	} else {
		e := models.Enforcer
		areRulesAdded := e.AddPolicy(role, "/jones/doc/*", "*")
		e.AddRoleForUser(role, role) // 增加一个和用户组同名的用户, 空用户的用户组在get时获取不到
		models.AddDefaultaCasbin(role)
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: fmt.Sprintf("%t", areRulesAdded),
		}.String())
	}
}

// @Title 对一个角色和一个页面限制只读权限 (命名空间_authority)
// @Summary  对一个角色和一个页面限制只读权限 (命名空间_authority)
// @Description  对一个角色和一个页面限制只读权限 (命名空间_authority)
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/roles/namespace/addReadOnly [post]
func (this *AuthController) AddReadOnlyNameSpace() {
	var nms []models.NameSpaceWritable
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &nms)

	for _, nm := range nms {
		// 1. 表里增加一项

		// 2. 增加一个只读标记
		e := models.Enforcer
		object := fmt.Sprintf("%s/%s_%s", RootNameSpace, nm.NameSpace, "authority")
		e.AddPolicy(nm.RoleName, object, "GET")
	}

	//this.Ctx.WriteString(fmt.Sprintf("%t",areRulesAdded))
}

// @Title getAllWorksheetCategory
// @Summary 查看所有工单目录
// @Description 查看所有工单目录
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /auth/template/getAllWorksheetCategory/:role [get]
func (this *AuthController) GetAllWorksheetCategory() {
	template := models.SqlTemplate{}
	//allCagtegory := template.GetAllSqlTemplateCategory() // 所有工单
	role := this.Ctx.Input.Params()[":role"]
	allCagtegory := template.GetCheckedCategory(role)
	res, _ := json.Marshal(allCagtegory)
	this.Ctx.WriteString(string(res))
}

// @Title GetAllMetaTableAccess
// @Summary 查看所有元数据表的权限
// @Description 查看所有元数据表的权限
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /auth/metadata/getAllMetaTableAccess/:role [get]
func (this *AuthController) GetAllMetaTableAccess() {
	metaaccess := models.TableMetaData{}
	role := this.Ctx.Input.Params()[":role"]
	allAccess := metaaccess.GetCheckedMetaTable(role)
	res, _ := json.Marshal(allAccess)
	this.Ctx.WriteString(string(res))
}

// @Title ChangeMetaTableMapping
// @Summary  更改用户和数据源权限的配置
// @Description  更改用户和数据源权限的配置
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/metadata/changeMetaTableMapping [post]
func (this *AuthController) ChangeMetaTableMapping() {
	var attr map[string][]string
	json.Unmarshal(this.Ctx.Input.RequestBody, &attr)

	role := attr["role"][0]
	workIds := attr["ids"]
	updateType := attr["type"][0]
	metaaccess := models.TableMetaData{}
	metaaccess.UpdateMetaTableMapping(role, workIds, updateType) // 更新权限

}

// @Title GetAllMetaTableTreeAccess
// @Summary 查看所有波动图表的树状形式的权限
// @Description 查看所有波动图表的树状形式的权限
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /auth/bodong/getTreeAccess/:role [get]
func (this *AuthController) GetBodongTreeAccess() {
	//role := this.Ctx.Input.Params()[":role"]
	o := orm.NewOrm()
	firstLevelNodes := []models.BodongMenuNode{}
	//nodes := []models.BodongMenuNode{}
	_, _ = o.Raw("SELECT * FROM bodong_menu_bar WHERE pid=0").QueryRows(&firstLevelNodes)
	for _, node := range firstLevelNodes {
		children := []models.BodongMenuNode{}
		_, _ = o.Raw("SELECT * FROM bodong_menu_bar WHERE pid=?", node.Mid).QueryRows(&children)
		node.Children = children
	}
	bytes,_ := json.Marshal(firstLevelNodes)
	this.Ctx.WriteString(string(bytes))
}

// @Title GetAllMetaTableTreeAccess
// @Summary 插卡所有数据表的树状形式的权限
// @Description 插卡所有数据表的树状形式的权限
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /auth/metadata/GetAllMetaTableTreeAccess/:role [get]
func (this *AuthController) GetAllMetaTableTreeAccess() {
	metaaccess := models.TableMetaData{}
	role := this.Ctx.Input.Params()[":role"]
	//allAccess    := metaaccess.GetCheckedMetaTable(role)
	allTableInfo := metaaccess.GetMetaTableFullInfoList(role)

	treeNodeList := make([]models.TreeViewNode, 0)

	lastConn := ""
	lastDatasource := ""
	curConnIdx := -1
	curDsInx := -1

	noLeafId := 1000001

	//循环遍历 获取 树的各个层次的 去重信息
	for _, tableInfo := range allTableInfo {
		//dsid 		:= tableInfo.Dsid
		conn := tableInfo.ConnName
		datasource := tableInfo.DatasourceName
		database := tableInfo.DatabaseName
		table := tableInfo.TableName
		tid := tableInfo.Tid

		if lastConn == "" || conn != lastConn {
			var conn_node models.TreeViewNode
			conn_node.Id = noLeafId
			noLeafId += 1
			curDsInx = -1
			conn_node.NodeName = conn
			conn_node.NodeChildren = []models.TreeViewNode{}

			treeNodeList = append(treeNodeList, conn_node)
			curConnIdx += 1
			lastConn = conn
		}

		if lastDatasource == "" || datasource != lastDatasource {
			var ds_node models.TreeViewNode
			ds_node.Id = noLeafId
			noLeafId += 1
			ds_node.NodeName = datasource
			ds_node.NodeChildren = make([]models.TreeViewNode, 0)

			treeNodeList[curConnIdx].NodeChildren = append(treeNodeList[curConnIdx].NodeChildren, ds_node)
			curDsInx += 1
			lastDatasource = datasource
		}

		if database == "" {
			continue
		}

		//添加表节点
		var db_tb_node models.TreeViewNode
		db_tb_node.Id = tid
		db_tb_node.NodeName = database + "." + table
		db_tb_node.NodeChildren = make([]models.TreeViewNode, 0)

		treeNodeList[curConnIdx].NodeChildren[curDsInx].NodeChildren = append(treeNodeList[curConnIdx].NodeChildren[curDsInx].NodeChildren, db_tb_node)
	}

	res, _ := json.Marshal(treeNodeList)
	this.Ctx.WriteString(string(res))

}

// @Title GetMetaTableAccess
// @Summary 获取当前角色数据表的读写权限
// @Description 获取当前角色数据表的读写权限
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /auth/metadata/GetMetaTableAccess/:role [get]
func (this *AuthController) GetMetaTableAccess() {
	metaaccess := models.TableMetaData{}
	role := this.Ctx.Input.Params()[":role"]
	allAccess := metaaccess.GetCheckedMetaTable(role)
	res, _ := json.Marshal(allAccess)
	this.Ctx.WriteString(string(res))
}

// @Title GetAllBodongTree
// @Summary 获取所有波动图表树
// @Description 获取所有波动图表树
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /auth/bodong/GetAllBodongTree [get]
func (this *AuthController) GetAllBodongTree() {
	treeNodes := models.GetBoDongMenuTree()
	res, _ := json.Marshal(treeNodes)
	this.Ctx.WriteString(string(res))
}

// @Title GetAllBodongTree
// @Summary 获取选中的id
// @Description 获取选中的id
// @router /auth/bodong/GetCheckedBodongNode [get]
func (this *AuthController) GetCheckedBodongNode() {
	role := this.Ctx.Input.Query("role")
	o := orm.NewOrm()
	checkedIds := []int{}
	_, _ = o.Raw("SELECT menu_id FROM bodong_menu_mapping WHERE group_name = ?", role).QueryRows(&checkedIds)
	bytes,_ := json.Marshal(checkedIds)
	this.Ctx.WriteString(string(bytes))
}

// @Title UpdateCheckedBodongNode
// @Summary 更新选中的id
// @Description 更新选中的id
// @router /auth/bodong/UpdateCheckedBodongNode [post]
func (this *AuthController) UpdateCheckedBodongNode() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	targetRole := attr["role"]
	checkedIds := strings.Split(attr["ids"],",")

	err := models.UpdateBodongMapping(targetRole, checkedIds)
	if err != nil {
		models.Logger.Error("%s", err.Error())
		this.Ctx.WriteString(utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}.String())
	} else {
		this.Ctx.WriteString(utils.VueResponse{
			Res:  true,
			Info: "",
		}.String())
	}
}

// @Title 更新标签的排序顺序
// @Summary  更新标签的排序顺序
// @Description  更新标签的排序顺序
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/label/updateSortId [post]
func (this *AuthController) UpdateSortId() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)

	sortId, _ := strconv.ParseInt(attr["sortId"], 10, 64)
	labelId, _ := strconv.ParseInt(attr["labelId"], 10, 64)

	o := orm.NewOrm()
	_, err := o.Raw("UPDATE top_label SET sort_id = ? WHERE lid = ?", sortId, labelId).Exec()
	msg := "success"
	if err != nil {
		msg = err.Error()
	}
	this.Ctx.WriteString(msg)
}

// @Title 登录时添加用户
// @Summary  登录时添加用户
// @Description  登录时添加用户
// @Accept json
// @router /auth/user/postUser [post]
func (this *AuthController) PostOrUpdateUser() {
	var userNameMapping models.UserNameMapping
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &userNameMapping)
	userNameMapping.AddOrUpdate()
}


// @Title 更改域的管理员
// @Summary  更改域的管理员
// @Description  更改域的管理员
// @Accept json
// @router /auth/lvpiDomain/postAdministrator [post]
func (this *AuthController) PostAdministrator() {
	// 目前只有ID字段有用
	var domainLvpi models.DomainLvpi
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &domainLvpi)
	o := orm.NewOrm()
	_, _ = o.Raw("UPDATE domain_lvpi SET administrator = ? WHERE did = ?",
		domainLvpi.Administrator, domainLvpi.Did).Exec()
}

// @Title 添加一级域
// @Summary  添加一级域
// @Description  添加一级域
// @Accept json
// @router /auth/lvpiDomain/postFirstDomain [post]
func (this *AuthController) PostFirstDomain() {
	var domainLvpi models.DomainLvpi
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &domainLvpi)
	o := orm.NewOrm()
	_, _ = o.Insert(&domainLvpi)
	this.Ctx.WriteString("添加完成")
}

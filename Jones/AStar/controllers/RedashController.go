package controllers

import (
	"AStar/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type RedashController struct {
	BaseController
}

var RedashNameSpace = "redash"

// @Title GetTimeStamp
// @Summary 获取服务器时间戳
// @Description 获取服务器时间戳
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /timestamp [get]
func (this *RedashController) GetTimeStamp() {
	cur_time := time.Now().Unix()
	marshal, _ := json.Marshal(cur_time)
	this.Ctx.WriteString(string(marshal))
}

// @Title InsertRedash
// @Summary 添加图表显示
// @Description 添加图表显示
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /redashinsert [post]
func (this *RedashController) InsertRedash() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	grpList := attr["privilegeRoles"]
	newLid := models.InsertRedash(attr)
	if(newLid != -1 && len(grpList) != 0) {
		model := models.LabelModel{}
		model.AddLabelToMultiGroup(grpList, newLid)
	}
	this.Ctx.WriteString(fmt.Sprintf("%t", newLid))
}


// @Title CheckRedash
// @Summary 查看图表url是否已存在
// @Description 查看图表url是否已存在
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /redashcheck [post]
func (this *RedashController) CheckRedash() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)

	resInfo := models.CheckRedash(attr)
	jsonData, _ := json.Marshal(resInfo)
	this.Ctx.WriteString(fmt.Sprintf("%s", jsonData))
}


// @Title GetAllRedash
// @Summary 获取所有图表信息
// @Description 获取所有图表信息
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /redashall [get]
func (this *RedashController) GetAllRedash() {


	resInfo := models.GetAllRedashInfo("redashshow")
	jsonData, _ := json.Marshal(resInfo)
	this.Ctx.WriteString(fmt.Sprintf("%s", jsonData))
}


// @Title DeleteRedash
// @Summary 删除图表
// @Description 删除图表
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /redashdel [post]
func (this *RedashController) DeleteRedash() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)

	delIndex := attr["redashIndex"]
	o := orm.NewOrm()

	urlSql :=fmt.Sprintf(" delete from top_label where `index`='%s' ",delIndex)
	res, err := o.Raw(urlSql).Exec()
	if(err != nil){
		fmt.Printf("Delete get err:%s | res:%s",err,res)
	}
	jsonData, _ := json.Marshal("DeleteSuccess")
	this.Ctx.WriteString(fmt.Sprintf("%s", jsonData))
}


// @Title InsertDashboard
// @Summary 添加Dashboard
// @Description 添加Dashboard
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dashboardinsert [post]
func (this *RedashController) InsertDashboard() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	dashname 	:= attr["dashboardname"]
	usrid 		:= attr["usrid"]
	_,dashurl := models.InsertDashBoard(dashname,usrid)
	this.Ctx.WriteString(fmt.Sprintf("%t", dashurl))
}


// @Title CheckWriteAccess
// @Summary 获取Dashboard和RedashQuery是否有写权限
// @Description 获取Dashboard和RedashQuery是否有写权限
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /write_access [get]
func (this *RedashController) CheckWriteAccess() {
	roles := models.Enforcer.GetRolesForUser(this.userName)
	dash_CanWrite := models.GetIndexIsWrite("dashboardoption",roles)
	query_CanWrite := models.GetIndexIsWrite("redashqueryoption",roles)

	var attr map[string]bool
	attr = make(map[string]bool)
	attr["dashboard_write"] = dash_CanWrite
	attr["redashquery_write"] = query_CanWrite

	jsonData, _ := json.Marshal(attr)
	this.Ctx.WriteString(fmt.Sprintf("%s", jsonData))
}


// @Title ChangeUser
// @Summary 指派权限
// @Description 指派权限
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /user_change [post]
func (this *RedashController) ChangeUser() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)

	assignType := attr["assignType"]
	assignIndex := attr["assignIndex"]
	assignUser := attr["assignUserName"]
	flag := models.ChangeAssignUser(assignType,assignIndex,assignUser)

	var res map[string]bool
	res = make(map[string]bool)
	res["Res"] = flag

	jsonData, _ := json.Marshal(res)
	this.Ctx.WriteString(fmt.Sprintf("%s", jsonData))
}


// @Title GetAllDashboard
// @Summary 获取所有Dashbord信息
// @Description 获取所有Dashbord信息
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dashboardall [get]
func (this *RedashController) GetAllDashboard() {
	//checkUrl := fmt.Sprintf("%s/%s_%s", RootNameSpace,RedashNameSpace, "authority")
	//models.Logger.Info(checkUrl)
	//get_all_flag := models.Enforcer.Enforce(this.userName, checkUrl, "POST")

	roles := models.Enforcer.GetRolesForUser(this.userName)
	canWrite := models.GetIndexIsWrite("dashboardoption",roles)

	resInfo := models.GetAllDashboardInfo(this.userName,canWrite)
	jsonData, _ := json.Marshal(resInfo)
	this.Ctx.WriteString(fmt.Sprintf("%s", jsonData))
}


// @Title DeleteDashboard
// @Summary 删除Dashboard
// @Description 删除Dashboard
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dashboarddel [post]
func (this *RedashController) DeleteDashboard() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)

	delIndex := attr["dashboardIndex"]
	flag := models.DeleteDashBoard(delIndex)
	jsonData, _ := json.Marshal(flag)
	this.Ctx.WriteString(fmt.Sprintf("%s", jsonData))
}

// @Title UpdateDashboard
// @Summary 修改Dashboard名字
// @Description 修改Dashboard名字
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /dashboardupdate [post]
func (this *RedashController) UpdateDashboard() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)


	upIndex := attr["dashboardIndex"]
	upName := attr["dashboardName"]
	//upUser := this.userName
	upCreator := ""

	flag := models.UpdateDashBoard(upIndex,upName,upCreator)
	jsonData, _ := json.Marshal(flag)
	this.Ctx.WriteString(fmt.Sprintf("%s", jsonData))
}


// @Title GetAllRedahQuery
// @Summary 获取所有Query信息
// @Description 获取所有Query信息
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /redashqueryall [get]
func (this *RedashController) GetAllRedahQuery() {
	roles := models.Enforcer.GetRolesForUser(this.userName)
	canWrite := models.GetIndexIsWrite("redashqueryoption",roles)

	resInfo := models.GetAllRedashQueryInfo(this.userName,canWrite)
	jsonData, _ := json.Marshal(resInfo)
	this.Ctx.WriteString(fmt.Sprintf("%s", jsonData))
}


// @Title DeleteRedashQuery
// @Summary 删除Query
// @Description 删除Query
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /redashquerydel [post]
func (this *RedashController) DeleteRedashQuery() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)

	delIndex := attr["queryIndex"]
	flag := models.DeleteRedashQuery(delIndex)
	jsonData, _ := json.Marshal(flag)
	this.Ctx.WriteString(fmt.Sprintf("%s", jsonData))
}

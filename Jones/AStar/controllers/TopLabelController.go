package controllers

import (
	"AStar/models"
	"AStar/utils"
	"encoding/json"
	"strings"
)

type TopLabelController struct {
	BaseController
}

// @Title GetLabels
// @Summary 获取所有的顶部标签
// @Description 获取所有的顶部标签
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/labels/:groupName [get]
func (this *TopLabelController) GetLabels() {
	groupName := this.Ctx.Input.Params()[":groupName"]
	model := models.LabelModel{}
	views := model.GetAllTopLabel(groupName)
	marshal, _ := json.Marshal(views)
	this.Ctx.WriteString(string(marshal))
}

// @Title UpdateLabels
// @Summary 更新顶部标签
// @Description 更新顶部标签
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/labels/update [post]
func (this *TopLabelController) UpdateLabels() {
	var attr map[string][]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)

	role := attr["role"][0]
	urlIds := attr["ids"]
	writableUrls := attr["writableUrls"]
	model := models.LabelModel{}
	writableUrls = utils.MakeUniqSet(writableUrls...)
	model.AddLabelMapping(role, urlIds, writableUrls)

}

// @Title GetCheckedLabels
// @Summary 获取有权限的顶栏标签
// @Description 获取有权限的顶栏标签
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/labels/checked/:username [get]
func (this *TopLabelController) GetCheckedLabels() {
	username := this.Ctx.Input.Params()[":username"]
	roles := models.Enforcer.GetRolesForUser(username)
	model := models.LabelModel{}
	views := model.GetCheckedTopLabel(roles)
	marshal, _ := json.Marshal(views)
	this.Ctx.WriteString(string(marshal))
}

// @Title GetPMLabels
// @Summary 获取产品关注的标签属性
// @Description 获取产品关注的标签属性
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/labels/pmlabels [get]
func (this *TopLabelController) GetPMLabels() {
	model := models.LabelModel{}
	views := model.GetPMLabelAttr()
	marshal, _ := json.Marshal(views)
	this.Ctx.WriteString(string(marshal))
}

// @Title UpdatePMLabels
// @Summary 更新PM关注的顶栏标签属性
// @Description 更新PM关注的顶栏标签属性
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/labels/pmupdate [post]
func (this *TopLabelController) UpdatePMLabels() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)

	lid := attr["lid"]
	title := attr["title"]
	class := attr["class"]
	pid := attr["pid"]

	model := models.LabelModel{}
	views := model.UpdatePMLabelAttr(lid, title, class, pid)
	marshal, _ := json.Marshal(views)
	this.Ctx.WriteString(string(marshal))

}

// @Title GetRDLabels
// @Summary 获取开发关注的标签属性
// @Description 获取开发关注的标签属性
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/labels/rdlabels [get]
func (this *TopLabelController) GetRDLabels() {
	model := models.LabelModel{}
	views := model.GetRDLabelAttr()
	marshal, _ := json.Marshal(views)
	this.Ctx.WriteString(string(marshal))
}

// @Title AddRDLabels
// @Summary 新增RD标签
// @Description 新增RD标签
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/labels/rdadd [post]
func (this *TopLabelController) AddRDLabels() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)

	title := attr["title"]
	index := attr["index"]
	namespace := attr["namespace"]

	model := models.LabelModel{}
	views := model.AddRDLabel(title, index, namespace)
	marshal, _ := json.Marshal(views)
	this.Ctx.WriteString(string(marshal))

}

// @Title UpdateRDLabels
// @Summary 更新RD关注的顶栏标签属性
// @Description 更新RD关注的顶栏标签属性
// @Success 200 {array[string]}
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /auth/labels/rdupdate [post]
func (this *TopLabelController) UpdateRDLabels() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)

	lid := attr["lid"]
	index := attr["index"]
	namespace := attr["namespace"]

	model := models.LabelModel{}
	views := model.UpdateRDLabelAttr(lid, index, namespace)
	marshal, _ := json.Marshal(views)
	this.Ctx.WriteString(string(marshal))

}

// @Title AddNewLabels
// @Summary 批量新智能标签权限
// @Description 批量新智能标签权限
// @Accept json
// @router /auth/labels/addNewLabels [post]
func (this *TopLabelController) AddNewLabels() {
	var attr map[string]string
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	roles := strings.Split(attr["roles"], ",")
	labelNames := strings.Split(attr["labelNames"], ",")
	models.AddLabelsForUser(roles, labelNames)
	this.Ctx.WriteString("添加成功")
}

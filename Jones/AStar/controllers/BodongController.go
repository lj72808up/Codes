package controllers

import (
	"AStar/models"
	"encoding/json"
	"strconv"
)

type BodongController struct {
	BaseController
}

var BodongNameSpace = "bodong"

// @Title 获取一级标签
// @Summary 获取一级标签
// @Description 获取一级标签
// @Success 200
// @Accept json
// @router /getFirst [get]
func (this *BodongController) GetFirstLevel() {
	bar := models.BodongMenuBar{}
	roles := models.Enforcer.GetRolesForUser(this.userName)
	menus := bar.GetFirstMenu(roles)
	bytes, _ := json.Marshal(menus)
	this.Ctx.WriteString(string(bytes))
}

// @Title 获取子菜单
// @Summary 获取子菜单
// @Description 获取子菜单
// @Success 200
// @Accept json
// @router /getChildOptions/:mid [get]
func (this *BodongController) GetChildOptions() {
	mid, _ := strconv.ParseInt(this.Ctx.Input.Params()[":mid"], 10, 64)
	bar := models.BodongMenuBar{}
	// 权限
	roles := models.Enforcer.GetRolesForUser(this.userName)
	menus := bar.GetChildMenus(mid,roles)
	bytes, _ := json.Marshal(menus)
	this.Ctx.WriteString(string(bytes))
}

/*// @Title 获取点击的菜单
// @Summary 获取点击的菜单
// @Description 获取点击的菜单
// @Success 200
// @Accept json
// @router /getCheckedMenuIds [get]
func (this *BodongController) GetCheckedMenuIds() {
	roles := models.Enforcer.GetRolesForUser(this.userName)
	bar := models.BodongMenuBar{}
	ids := bar.GetCheckedMenuIds(roles)
	bytes,_ := json.Marshal(ids)
	this.Ctx.WriteString(string(bytes))
}*/
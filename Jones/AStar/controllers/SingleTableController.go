package controllers

import (
	"AStar/models/singleTable"
	"encoding/json"
	"strconv"
)

var SingleNameSpace = "single"

type SingleTableController struct {
	BaseController
}

// @Title AddDim
// @Summary 添加 dim
// @Description 添加 dim
// @Success 200
// @router /addDim [post]
func (this *SingleTableController) AddDim() {
	var dim singleTable.DimStyleClassify
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &dim)
	singleTable.AddDimStyleClassify(dim)
}

// @Title DelDim
// @Summary 删除 dim
// @Description 删除 dim
// @Success 200
// @router /delDim/:id [get]
func (this *SingleTableController) DelDim() {
	id, _ := strconv.ParseInt(this.Ctx.Input.Params()[":id"], 10, 64)
	singleTable.DeleteDimStyleClassifySet(id)
}

// @Title UpdateDim
// @Summary 修改 dim
// @Description 修改 dim
// @Success 200
// @router /updateDim [post]
func (this *SingleTableController) UpdateDim() {
	var dim singleTable.DimStyleClassify
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &dim)
	singleTable.UpdateDimStyleClassifySet(dim)
}

// @Title GetDims
// @Summary 获取 dim
// @Description 获取 dim
// @Success 200
// @router /getDims [get]
func (this *SingleTableController) GetDims() {
	dims := singleTable.GetDimStyleClassifySet()
	bytes, _ := json.Marshal(dims)
	this.Ctx.WriteString(string(bytes))
}

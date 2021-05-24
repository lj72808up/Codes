package controllers

import (
	"AStar/models"
	"AStar/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type UdfController struct {
	BaseController
}

var UdfNameSpace = "udf"

// @Title 获取所有udf列表
// @Summary 获取所有udf列表
// @Description 获取所有udf列表
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /getList [post]
func (this *UdfController) GetUdfList() {
	var attr map[string]bool
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &attr)
	known := attr["known"]
	udfManager := models.UdfManager{}
	res, _ := json.Marshal(udfManager.GetUdfList(known))
	this.Ctx.WriteString(string(res))
}

// @Title 更新函数描述
// @Summary 更新函数描述
// @Description 更新函数描述
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /update/description [post]
func (this *UdfController) UpdateFunctionDesc() {
	var udfs []models.UdfManager
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &udfs)
	udfManager := models.UdfManager{}
	err := udfManager.UpdateDescription(udfs)
	vueRes := utils.VueResponse{
		Res:  true,
		Info: "更新成功",
	}
	if err != nil {
		vueRes = utils.VueResponse{
			Res:  false,
			Info: err.Error(),
		}
	}
	res, _ := json.Marshal(vueRes)
	this.Ctx.WriteString(string(res))
}

// @Title GetEdit
// @Summary 查看是否可以编辑UDF管理
// @Description 查看是否可以编辑UDF管理
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /getEdit [get]
func (this *UdfController) GetEdit() {
	method := "POST"
	checkUrl := fmt.Sprintf("%s/%s_%s", RootNameSpace, UdfNameSpace, "authority")
	models.Logger.Info(checkUrl)
	flag := models.Enforcer.Enforce(this.userName, checkUrl, method)
	res, _ := json.Marshal(map[string]bool{
		"readOnly": !flag,
	})
	this.Ctx.WriteString(string(res))
}

// @Title 删除
// @Summary 删除
// @Description 删除
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /delete/:id [get]
func (this *UdfController) DelById() {
	id, _ := strconv.ParseInt(this.Ctx.Input.Params()[":id"], 10, 64)
	o := orm.NewOrm()
	o.Raw("DELETE FROM udf_manager WHERE func_id = ?", id).Exec()
}

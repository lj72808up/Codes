package controllers

import (
	"AStar/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
)

type MobileDashboardController struct {
	BaseController
}


var  MobileDashboardNameSpace = "mobileDashboard"

// @Title GetMobileDashboardUrl
// @Summary 查看 redash url
// @Description 查看 redash url
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @router /getDashboardUrl/:labelAbbreviation [get]
func (this *MobileDashboardController) GetMobileDashboardUrl() {
	labelAbbreviation := this.Ctx.Input.Params()[":labelAbbreviation"]
	index := fmt.Sprintf("redashshow/%s",labelAbbreviation)
	var topLabel models.TopLabel

	o := orm.NewOrm()
	_ = o.Raw("SELECT * FROM top_label WHERE `index` = ?", index).QueryRow(&topLabel)
	var jsonParam map[string]string
	_ = json.Unmarshal([]byte(topLabel.Info), &jsonParam)

	redashUrl := jsonParam[labelAbbreviation]
	this.Ctx.WriteString(redashUrl)
}
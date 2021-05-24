package controllers

import (
	"AStar/models"
	"AStar/utils"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"time"
)

type DocController struct {
	BaseController
}

// @Title 查看文档内容
// @Summary 查看最新系统文档
// @Description 查看最新系统文档
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /doc/get [get]
func (this *DocController) GetDoc() {
	var QueryDoc []models.PlatformDoc
	orm.NewOrm().QueryTable("platform_doc").All(&QueryDoc)
	data, err := json.Marshal(QueryDoc[0].Content)
	if (err != nil) && (len(QueryDoc) > 0) {
		models.Logger.Error("数据库查询失败")
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte("查询失败"))
		return
	}
	this.Ctx.WriteString(string(data))
}


// @Title 查看文档内容HTML
// @Summary 查看最新系统文档HTML
// @Description 查看最新系统文档HTML
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /doc/gethtml [get]
func (this *DocController) GetDocHtml() {
	var QueryDoc []models.PlatformDoc
	orm.NewOrm().QueryTable("platform_doc").All(&QueryDoc)
	data, err := json.Marshal(QueryDoc[0].Html)
	if (err != nil) && (len(QueryDoc) > 0) {
		models.Logger.Error("数据库查询失败")
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte("查询失败"))
		return
	}
	this.Ctx.WriteString(string(data))
}

// @Title 提交文档内容
// @Summary 提交系统文档
// @Description 提交文档
// @Param para body models.SubmitDocJson true "提交信息"
// @Success 200
// @Failure 400 Bad request
// @Failure 403 Forbidden
// @Failure 404 Not found
// @Accept json
// @router /doc/submit [post]
func (this *DocController) SubmitDoc() {
	var queryJSON = models.SubmitDocJson{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &queryJSON)
	if err != nil {
		models.Logger.Error(err.Error() + "JSON解释失败")
		this.Ctx.ResponseWriter.WriteHeader(400)
		this.Ctx.ResponseWriter.Write([]byte(err.Error() + "JSON解释失败"))
		return
	}

	var submitter = queryJSON.Submitter
	var content = queryJSON.Content
	var html = queryJSON.Html

	models.PlatformDoc{1, "r", "r", "r","r"}.UpdateContent(utils.DateFormatWithoutLine(time.Now()), submitter, content,html)
	this.Ctx.WriteString("保存成功")
}

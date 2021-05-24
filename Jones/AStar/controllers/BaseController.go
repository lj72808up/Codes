package controllers

import (
	"AStar/models"
	"AStar/utils"
	"github.com/astaxie/beego"
)

var RootNameSpace = "/jones"
var XiaopNameSpace = "/xiaop"
var LvpiNameSpace = "/lvpi"

type BaseController struct {
	beego.Controller
	controllerName string
	actionName     string
	userName       string
}

func (this *BaseController) Prepare() {
	this.controllerName, this.actionName = this.GetControllerAndAction()
	this.userName, _, _ = this.Ctx.Request.BasicAuth()
	models.Logger.Info(this.userName +
		"\t" + this.controllerName +
		"\t" + this.actionName +
		"\t" + utils.GetHostFromAddr(this.Ctx.Request.RemoteAddr) +
		"\t" + this.Ctx.Request.RequestURI)
}

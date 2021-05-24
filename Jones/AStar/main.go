package main

import (
	"AStar/models"
	_ "AStar/routers"
	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	models.Init()
	beego.Run()
	models.Logger.Close()
}

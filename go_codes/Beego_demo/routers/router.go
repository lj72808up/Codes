package routers

import (
	"Beego_demo/controllers"
	"github.com/astaxie/beego"
)


// 映射URL到controller
func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/user", &controllers.UserController{})
}

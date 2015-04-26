package routers

import (
	"github.com/astaxie/beego"
	"updateHelper/controllers"
)

func init() {
	beego.Router("/StartUpdate", &controllers.MainController{}, "get:StartUpdate")
}

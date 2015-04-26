package routers

import (
	"github.com/astaxie/beego"
	"updateHelper/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:TestAlive")
}

package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	"xchango_mid/controllers"
)

func init() {
	beego.Router("/api/v1/reportes", &controllers.ReportesController{}, "post:Crear")
}

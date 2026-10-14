package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	"xchango_mid/controllers"
)

func init() {
	// panel de administracion
	beego.Router("/api/v1/admin/admins", &controllers.PanelController{}, "get:Admins;post:CrearAdmin")
	beego.Router("/api/v1/admin/admins/:id", &controllers.PanelController{}, "delete:EliminarAdmin")
	beego.Router("/api/v1/admin/historial", &controllers.PanelController{}, "get:Historial;post:RegistrarAccion")
}
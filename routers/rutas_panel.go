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
	beego.Router("/api/v1/admin/alertas", &controllers.PanelController{}, "get:Alertas")
	beego.Router("/api/v1/admin/alertas/leidas", &controllers.PanelController{}, "patch:TodasLeidas")
	beego.Router("/api/v1/admin/alertas/:id/leida", &controllers.PanelController{}, "patch:AlertaLeida")
	beego.Router("/api/v1/admin/metricas", &controllers.PanelController{}, "get:Metricas")
}
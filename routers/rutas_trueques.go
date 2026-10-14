package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	"xchango_mid/controllers"
)

func init() {
	// trueques y notificaciones
	beego.Router("/api/v1/trueques", &controllers.TruequesController{}, "get:Listar;post:Solicitar")
	beego.Router("/api/v1/trueques/:id/aceptar", &controllers.TruequesController{}, "patch:Aceptar")
	beego.Router("/api/v1/trueques/:id/rechazar", &controllers.TruequesController{}, "patch:Rechazar")
	beego.Router("/api/v1/trueques/:id/confirmar", &controllers.TruequesController{}, "patch:Confirmar")
	beego.Router("/api/v1/trueques/:id/calificar", &controllers.TruequesController{}, "patch:Calificar")
	beego.Router("/api/v1/notificaciones", &controllers.TruequesController{}, "get:Notificaciones")
	beego.Router("/api/v1/notificaciones/leidas", &controllers.TruequesController{}, "patch:MarcarTodasLeidas")
	beego.Router("/api/v1/notificaciones/:id/leida", &controllers.TruequesController{}, "patch:MarcarLeida")
}

package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"

	"xchango_mid/controllers"
)

func init() {
	// para que Angular pueda hacer peticiones al MID
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{beego.AppConfig.DefaultString("origen_cliente", "http://localhost:4200")},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// rutas de inicio de sesión, registro y usuarios
	beego.Router("/api/v1/auth/pin", &controllers.AuthController{}, "post:EnviarPin")
	beego.Router("/api/v1/auth/verificar", &controllers.AuthController{}, "post:VerificarPin")
	beego.Router("/api/v1/auth/registro", &controllers.AuthController{}, "post:Registro")
	beego.Router("/api/v1/auth/perfil", &controllers.AuthController{}, "get:Perfil")
	beego.Router("/api/v1/auth/admin/login", &controllers.AuthController{}, "post:LoginAdmin")
	beego.Router("/api/v1/usuarios", &controllers.AuthController{}, "get:ListarUsuarios")
	beego.Router("/api/v1/usuarios/perfil", &controllers.AuthController{}, "put:ActualizarPerfil")

	// publicaciones y favoritos
	beego.Router("/api/v1/categorias", &controllers.PublicacionesController{}, "get:Categorias")
	beego.Router("/api/v1/avisos-ubicacion", &controllers.PublicacionesController{}, "get:Avisos")
	beego.Router("/api/v1/publicaciones", &controllers.PublicacionesController{}, "get:Listar;post:Crear")
	beego.Router("/api/v1/publicaciones/:id", &controllers.PublicacionesController{}, "put:Actualizar;delete:Eliminar")
	beego.Router("/api/v1/publicaciones/:id/estado", &controllers.PublicacionesController{}, "patch:CambiarEstado")
	beego.Router("/api/v1/favoritos", &controllers.PublicacionesController{}, "get:Favoritos")
	beego.Router("/api/v1/favoritos/:id", &controllers.PublicacionesController{}, "post:AlternarFavorito")

	// solicitudes de trueque y notificaciones
	beego.Router("/api/v1/trueques", &controllers.TruequesController{}, "get:Listar;post:Solicitar")
	beego.Router("/api/v1/trueques/:id/aceptar", &controllers.TruequesController{}, "patch:Aceptar")
	beego.Router("/api/v1/trueques/:id/rechazar", &controllers.TruequesController{}, "patch:Rechazar")
	beego.Router("/api/v1/trueques/:id/confirmar", &controllers.TruequesController{}, "patch:Confirmar")
	beego.Router("/api/v1/trueques/:id/calificar", &controllers.TruequesController{}, "patch:Calificar")
	beego.Router("/api/v1/notificaciones", &controllers.TruequesController{}, "get:Notificaciones")
	beego.Router("/api/v1/notificaciones/leidas", &controllers.TruequesController{}, "patch:MarcarTodasLeidas")
	beego.Router("/api/v1/notificaciones/:id/leida", &controllers.TruequesController{}, "patch:MarcarLeida")

	// opciones del administrador para moderar usuarios y publicaciones
	beego.Router("/api/v1/admin/usuarios", &controllers.ModeracionController{}, "get:Usuarios")
	beego.Router("/api/v1/admin/usuarios/:id", &controllers.ModeracionController{}, "put:EditarUsuario")
	beego.Router("/api/v1/admin/usuarios/:id/estado", &controllers.ModeracionController{}, "patch:CambiarEstadoUsuario")
	beego.Router("/api/v1/admin/publicaciones", &controllers.ModeracionController{}, "get:Publicaciones")
	beego.Router("/api/v1/admin/reportes", &controllers.ModeracionController{}, "get:Reportes")
	beego.Router("/api/v1/admin/reportes/:id", &controllers.ModeracionController{}, "patch:ActualizarReporte")
	beego.Router("/api/v1/admin/reportes/:id/eliminar", &controllers.ModeracionController{}, "post:EliminarPublicacion")
	beego.Router("/api/v1/admin/eliminadas", &controllers.ModeracionController{}, "get:Eliminadas")
	beego.Router("/api/v1/admin/documentos", &controllers.ModeracionController{}, "get:Documentos")
	beego.Router("/api/v1/admin/documentos/:id", &controllers.ModeracionController{}, "patch:RevisarDocumento")

	// administración, historial y alertas del panel
	beego.Router("/api/v1/admin/admins", &controllers.PanelController{}, "get:Admins;post:CrearAdmin")
	beego.Router("/api/v1/admin/admins/:id", &controllers.PanelController{}, "delete:EliminarAdmin")
	beego.Router("/api/v1/admin/admins/:id/estado", &controllers.PanelController{}, "patch:EstadoAdmin")
	beego.Router("/api/v1/admin/historial", &controllers.PanelController{}, "get:Historial;post:RegistrarAccion")
	beego.Router("/api/v1/admin/alertas", &controllers.PanelController{}, "get:Alertas")
	beego.Router("/api/v1/admin/alertas/leidas", &controllers.PanelController{}, "patch:TodasLeidas")
	beego.Router("/api/v1/admin/alertas/:id/leida", &controllers.PanelController{}, "patch:AlertaLeida")
	beego.Router("/api/v1/admin/metricas", &controllers.PanelController{}, "get:Metricas")
}
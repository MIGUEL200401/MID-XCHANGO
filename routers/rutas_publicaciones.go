package routers

import (
	beego "github.com/beego/beego/v2/server/web"

	"xchango_mid/controllers"
)

func init() {
	// publicaciones y favoritos
	beego.Router("/api/v1/categorias", &controllers.PublicacionesController{}, "get:Categorias")
	beego.Router("/api/v1/avisos-ubicacion", &controllers.PublicacionesController{}, "get:Avisos")
	beego.Router("/api/v1/publicaciones", &controllers.PublicacionesController{}, "get:Listar;post:Crear")
	beego.Router("/api/v1/publicaciones/:id", &controllers.PublicacionesController{}, "put:Actualizar;delete:Eliminar")
	beego.Router("/api/v1/publicaciones/:id/estado", &controllers.PublicacionesController{}, "patch:CambiarEstado")
	beego.Router("/api/v1/favoritos", &controllers.PublicacionesController{}, "get:Favoritos")
	beego.Router("/api/v1/favoritos/:id", &controllers.PublicacionesController{}, "post:AlternarFavorito")

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
}

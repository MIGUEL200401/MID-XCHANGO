package controllers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"xchango_mid/models"
	"xchango_mid/services"
)

type ReportesController struct {
	beego.Controller
}

// POST /api/v1/reportes
// Un usuario reporta una publicación para que el administrador la revise
func (c *ReportesController) Crear() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	var solicitud models.SolicitudReporte

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	publicacionId := strings.TrimSpace(solicitud.PublicacionId)
	motivo := strings.TrimSpace(solicitud.Motivo)

	if publicacionId == "" {
		responderError(&c.Controller, 400, "Falta la publicación que quieres reportar.")
		return
	}

	if motivo == "" {
		responderError(&c.Controller, 400, "Escribe el motivo del reporte.")
		return
	}

	if largo(motivo) > 120 {
		responderError(&c.Controller, 400, "El motivo puede tener máximo 120 caracteres.")
		return
	}

	if !usuarioPuedeActuar(&c.Controller, sesion.ID) {
		return
	}

	crud := services.NewCrudService()

	var publicacion models.Publicacion
	err := crud.Obtener("publicaciones", publicacionId, &publicacion)
	if err != nil {
		if err == services.ErrNoEncontrado {
			responderError(&c.Controller, 404, "Esa publicación no existe.")
			return
		}
		errorCrud(&c.Controller, err)
		return
	}

	if publicacion.Estado == "eliminada" {
		responderError(&c.Controller, 404, "Esa publicación no existe.")
		return
	}

	if publicacion.UsuarioId == sesion.ID {
		responderError(&c.Controller, 400, "No puedes reportar tu propia publicación.")
		return
	}

	var reportes []models.Reporte
	err = crud.Listar("reportes", &reportes)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	// un usuario no reporta dos veces la misma publicación
	for _, r := range reportes {
		if r.PublicacionId == publicacionId && r.ReportanteId == sesion.ID {
			responderError(&c.Controller, 409, "Ya reportaste esta publicación.")
			return
		}
	}

	nuevo := models.Reporte{
		PublicacionId: publicacionId,
		ReportanteId:  sesion.ID,
		Motivo:        motivo,
		Fecha:         hoy(),
		Estado:        "pendiente",
		Gravedad:      "media",
	}

	var creado models.Reporte
	err = crud.Crear("reportes", nuevo, &creado)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responder(&c.Controller, 201, creado)
}

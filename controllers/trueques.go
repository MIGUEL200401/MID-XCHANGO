package controllers

import (
	"sort"
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"xchango_mid/models"
	"xchango_mid/services"
)

type TruequesController struct {
	beego.Controller
}

// GET /api/v1/trueques
func (c *TruequesController) Listar() {
	claims, ok := exigirUsuario(&c.Controller)
	if !ok {
		return
	}

	var todos []models.Trueque
	err := services.NewCrudService().Listar("trueques", &todos)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	mios := []models.Trueque{}
	for _, t := range todos {
		if t.SolicitanteId == claims.ID || t.PropietarioId == claims.ID {
			mios = append(mios, t)
		}
	}

	responder(&c.Controller, 200, mios)
}

// POST /api/v1/trueques
func (c *TruequesController) Solicitar() {
	claims, ok := exigirUsuario(&c.Controller)
	if !ok {
		return
	}

	var solicitud models.SolicitudTrueque

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	if !usuarioPuedeActuar(&c.Controller, claims.ID) {
		return
	}

	crud := services.NewCrudService()

	var publicacion models.Publicacion
	err := crud.Obtener("publicaciones", solicitud.PublicacionId, &publicacion)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	if publicacion.UsuarioId == claims.ID {
		responderError(&c.Controller, 400, "Esta publicación es tuya.")
		return
	}

	if publicacion.Estado != "activa" {
		responderError(&c.Controller, 400, "Esta publicación no está disponible.")
		return
	}

	var existentes []models.Trueque
	err = crud.Listar("trueques", &existentes)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	for _, t := range existentes {
		abierto := t.Estado == "pendiente" || t.Estado == "aceptado"
		if t.PublicacionId == publicacion.Id && t.SolicitanteId == claims.ID && abierto {
			responderError(&c.Controller, 409, "Ya tienes una propuesta abierta para esta publicación.")
			return
		}
	}

	ofrece := strings.TrimSpace(solicitud.Ofrece)
	if ofrece == "" {
		ofrece = "Propuesta por acordar"
	}
}
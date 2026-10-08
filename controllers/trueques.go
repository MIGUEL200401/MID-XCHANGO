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
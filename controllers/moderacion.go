package controllers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"MID-XCHANGO/models"
	"xchango_mid/services"
)

type ModeracionController struct {
	beego.Controller
}


//GET/api/v1/admin/usuarios
func (c *ModeracionController) Usuarios() {
	_, ok := exigirAdmin(&c.Controller, false)
	if !ok {
		return
	}

	usuarios, err := cargarUsuarios()
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.Controller, 200, usuarios)
}


//PATCH/api/v1/admin/usuarios/:id/estado
func (c *ModeracionController) CambiarEstadoUsuario() {
	_, ok := exigirAdmin(&c.Controller, false)
	if !ok {
		return
	}

	var cambio models.CambioEstado

	if !leerCuerpo(&c.Controller, &cambio) {
		return
	}

	if cambio.Estado != "activo" && cambio.Estado != "suspendido" && cambio.Estado != "advertido" {
		responderError(&c.Controller, 400, "El estado debe de ser activo, suspendido o advertido.")
		return
	}

	crud := services.NewCrudService()

	var usuario models.Usuario
	err := crud.Obtener("usuarios", c.Ctx.Input.Param(":id"), &usuario)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	usuario.Estado = cambio.Estado

	err = crud.Actualizar("usuarios", usuario.Id, usuario, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responderVacio(&c.Controller)
}
package controllers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"xchango_mid/models"
	"xchango_mid/services"
)

type PublicaionesController struct {
	beego.Controller
}

//Get/api/v1/categorias
func (c *PublicacionesController) Categorias() {
	var categorias []models.Categoria

	err := services.NewCrudServices().Listar("categorias", &categorias)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.Controller, 200, categorias)
}

//Get/api/v1/avisos-ubicacion
func (c *PublicaionesController) Avisos() {
	var avisos []models.AvisoUbicacion
	
	err := services.NewCrudService().Listar("avisos-ubicacion", &avisos)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.Controller, 200, avisos)
}

//Get/api/v1/publicaciones
func (c *PublicaionesController) Listar() {
	var todas []models.Publicacion

	err := services.NewCrudService().Listar("publicaciones", &todas)
	errorCrud(&c.Controller, err)
	return
}
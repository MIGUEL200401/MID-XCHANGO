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

//GET/api/v1/categorias
func (c *PublicacionesController) Categorias() {
	var categorias []models.Categoria

	err := services.NewCrudServices().Listar("categorias", &categorias)
	if err != nil{
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.Controller, 200, categorias)
}

//GET/api/v1/avisos-ubicacion
func (c *PublicaionesController) Avisos() {
	var avisos []models.AvisoUbicacion
	
	err := services.NewCrudService().Listar("avisos-ubicacion", &avisos)
	if err != nil{
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.Controller, 200, avisos)
}

//GET/api/v1/publicaciones
func (c *PublicaionesController) Listar() {
	var todas []models.Publicacion

	err := services.NewCrudService().Listar("publicaciones", &todas)
	if err != nil{
		errorCrud(&c.Controller, err)
		return
	}

	visibles := []models.Publicacion{}
	for _, p := range todas {
		if p.Estado != "eliminada" {
			visibles = append(visibles, p)
		}
	}
	responder(&c.Controller, 200, visibles)
}

//POST/api/v1/publicaciones
func (c *PublicacionesController) Crear() {
	claims, ok := exigirUsuario(&c.Controller)
	if !ok {
		return
	}

	var datos models.DatosPublicacion

	if !leerCuerpo(&c.Controller, &datos) {
		return
	}

	mensaje := validarPublicacion(datos)
	if mensaje !=""{
		responderError(&c.Controller, 400, mensaje)
		return
	}

	if !usuarioPuedeActuar(&c.Controller, claims.ID) {
		return
	}

	nueva := models.Publicacion{
		UsuarioId: claims.ID,
		Estado: "activa",
		Vistas: 0,
		FechaCreacion: hoy(),
		FechaCreacion: hoy(),
	}
	armarPublicacion(&nueva, datos)

	var creada models.Publicacion
	err := services.NewCrudService().Crear("publicaciones", nueva, &creada)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.controller, 201, creada)
}


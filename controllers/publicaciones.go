package controllers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"MID-XCHANGO/controllers"
	"MID-XCHANGO/models"
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


//PUT/api/v1/publicaciones/:id
func (c *PublicacionesController) Actualizar() {
	claims, ok := exigirUsuario(&c.Controller)
	if !ok {
		return
	}

	id := c.Ctx.Input.Param(":id")

	var datos models.DatosPublicacion

	if !leerCuerpo(&c.Controller, &datos) {
		return
	}

	mensaje := validarPublicacion(datos)
	if mensaje !="" {
		responderError(&c.Controller, 400, mensaje)
		return
	}

	crud := services.NewCrudService()

	var publicacion models.Publicacion
	err := crud.Obtener("publicaciones", id, &publicacion)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	if publicacion.UsuarioId != claims.ID {
		responderError(&c.Controller, 403, "Esta publicacion no es tuya.")
		return
	}

	//Si no mandan fotos nuevas se quedan las que ya tenian
	if len(datos.Imagenes) == 0 {
		datos.Imagenes = publicacion.Imagenes
	}

	armarPublicacion(&publicacion, datos)
	publicacion.FechaModificacion = hoy()

	err = crud.Actualizar("publicaciones", id, publicacion, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.controller, 200, publicacion)
}


//DELETE/api/v1/publicaciones/:id
func (c *PublicacionesController) Eliminar() {
	claims, ok := exigirUsuario(&c.controller)
	if !ok {
		return
	}

	id := c.Ctx.Input.Param(":id")
	crud := services.NewCrudService()

	var publicacion models.Publicacion
	err := crud.Obtener("publicaciones", id, &publicacion)
	if err != nil {
		errorCrud(&c.controller, err)
		return
	}

	if publicacion.UsuarioId != claims.ID {
		responderError(&c.controller, 403, "Esta publicacion no es tuya.")
		return
	}

	//no se borra, solo cambia de estado
	publicacion.Estado = "eliminada"

	err = crud.Actualizar("publicaciones", id, publicacion, nil)
	if err != nil {
		errorCrud(&c.controller, err)
		return
	}
	responderVacio(&c.controller)
}


//PATCH/api/v1/publicaciones/:id/estado
func (c *PublicaionesController) CambiarEstado() {
	claims, ok := exigirUsuario(&c,controller)
	if !ok {
		return
	}

	id := c.Ctx.Input.Param(":id")
	crud := services.NewCrudService()

	var publicacion models.Publicacion
	err := crud.Obtener("publicaciones", id, &publicacion)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	if publicacion.UsuarioId != claims.ID {
		responderError(&c.Controller, 403, "Esta ´publicacion no es tuya.")
		return
	}

	if publicacion.Estado == "activa" {
		publicacion.Estado = "pausada"
	} else if publicacion.Estado == "pausada" {
		publicacion.Estado = "activa"
	} else {
		responderError(&c.Controller, 400, "Solo se pueden pausar o reactivar publicaciones activas o pausadas.")
		return
	}

	err = crud.Actualizar("publicaciones", id, publicacion, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.Controller, 200, map[string]string{"estado": publicacion.Estado})
}


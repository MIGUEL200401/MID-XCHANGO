package controllers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"xchango_mid/models"
	"xchango_mid/services"
)

type PublicacionesController struct {
	beego.Controller
}

// GET/api/v1/categorias
func (c *PublicacionesController) Categorias() {
	var categorias []models.Categoria

	err := services.NewCrudService().Listar("categorias", &categorias)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.Controller, 200, categorias)
}

// GET/api/v1/avisos-ubicacion
func (c *PublicacionesController) Avisos() {
	var avisos []models.AvisoUbicacion

	err := services.NewCrudService().Listar("avisos-ubicacion", &avisos)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.Controller, 200, avisos)
}

// GET/api/v1/publicaciones
func (c *PublicacionesController) Listar() {
	var todas []models.Publicacion

	err := services.NewCrudService().Listar("publicaciones", &todas)
	if err != nil {
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

// POST/api/v1/publicaciones
func (c *PublicacionesController) Crear() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	var datos models.DatosPublicacion

	if !leerCuerpo(&c.Controller, &datos) {
		return
	}

	mensaje := validarPublicacion(datos)
	if mensaje != "" {
		responderError(&c.Controller, 400, mensaje)
		return
	}

	if !usuarioPuedeActuar(&c.Controller, sesion.ID) {
		return
	}

	nueva := models.Publicacion{
		UsuarioId:         sesion.ID,
		Estado:            "activa",
		Vistas:            0,
		FechaCreacion:     hoy(),
		FechaModificacion: hoy(),
	}
	armarPublicacion(&nueva, datos)

	var creada models.Publicacion
	err := services.NewCrudService().Crear("publicaciones", nueva, &creada)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.Controller, 201, creada)
}

// PUT/api/v1/publicaciones/:id
func (c *PublicacionesController) Actualizar() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	id := c.Ctx.Input.Param(":id")

	var datos models.DatosPublicacion

	if !leerCuerpo(&c.Controller, &datos) {
		return
	}

	mensaje := validarPublicacion(datos)
	if mensaje != "" {
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

	if publicacion.UsuarioId != sesion.ID {
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
	responder(&c.Controller, 200, publicacion)
}

// DELETE/api/v1/publicaciones/:id
func (c *PublicacionesController) Eliminar() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
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

	if publicacion.UsuarioId != sesion.ID {
		responderError(&c.Controller, 403, "Esta publicacion no es tuya.")
		return
	}

	//no se borra, solo cambia de estado
	publicacion.Estado = "eliminada"

	err = crud.Actualizar("publicaciones", id, publicacion, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responderVacio(&c.Controller)
}

// PATCH/api/v1/publicaciones/:id/estado
func (c *PublicacionesController) CambiarEstado() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
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

	if publicacion.UsuarioId != sesion.ID {
		responderError(&c.Controller, 403, "Esta publicacion no es tuya.")
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

// GET/api/v1/favoritos
func (c *PublicacionesController) Favoritos() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	var favoritos []models.Favoritos
	err := services.NewCrudService().Listar("favoritos", &favoritos)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	ids := []string{}
	for _, f := range favoritos {
		if f.UsuarioId == sesion.ID {
			ids = append(ids, f.PublicacionId)
		}
	}
	responder(&c.Controller, 200, ids)
}

// POST/api/v1/favoritos/:id (si ya era favorito lo quita)
func (c *PublicacionesController) AlternarFavorito() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	publicacionId := c.Ctx.Input.Param(":id")
	crud := services.NewCrudService()

	var favoritos []models.Favoritos
	err := crud.Listar("favoritos", &favoritos)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	for _, f := range favoritos {
		if f.UsuarioId == sesion.ID && f.PublicacionId == publicacionId {
			err = crud.Borrar("favoritos", f.Id)
			if err != nil {
				errorCrud(&c.Controller, err)
				return
			}
			responderVacio(&c.Controller)
			return
		}
	}

	nuevo := models.Favoritos{UsuarioId: sesion.ID, PublicacionId: publicacionId}

	err = crud.Crear("favoritos", nuevo, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responderVacio(&c.Controller)
}

func validarPublicacion(d models.DatosPublicacion) string {
	if d.Tipo != "bien_fisico" && d.Tipo != "servicio" && d.Tipo != "bien_digital" {
		return "El tipo de publicacion no es valido."
	}
	if strings.TrimSpace(d.Titulo) == "" {
		return "El titulo es obligatorio."
	}
	if strings.TrimSpace(d.Descripcion) == "" {
		return "La descripcion es obligatoria."
	}
	if strings.TrimSpace(d.CategoriaId) == "" {
		return "Elige una categoria."
	}
	if d.Tipo != "bien_digital" && strings.TrimSpace(d.Municipio) == "" {
		return "El municipio es obligatorio."
	}
	if d.Tipo == "bien _fisico" && strings.TrimSpace(d.Barrio) == "" {
		return "El barrio es obligatorio para un bien físico."
	}
	return ""
}

// El bien fisico lleva municipio y barrio, el servicio solo municipio
// Y el bien digital no lleva ubicacion
func armarPublicacion(p *models.Publicacion, d models.DatosPublicacion) {
	p.Tipo = d.Tipo
	p.CategoriaId = d.CategoriaId
	p.Titulo = strings.TrimSpace(d.Titulo)
	p.Descripcion = strings.TrimSpace(d.Descripcion)
	p.Ofreces = strings.TrimSpace(d.Ofreces)
	if p.Ofreces == "" {
		p.Ofreces = p.Titulo
	}
	p.Buscas = strings.TrimSpace(d.Buscas)
	p.Imagenes = d.Imagenes
	if p.Imagenes == nil {
		p.Imagenes = []string{}
	}

	p.Municipio = ""
	p.Barrio = ""
	p.CantidadDisponible = ""
	p.Disponibilidad = ""

	if d.Tipo == "bien_fisico" {
		p.Municipio = d.Municipio
		p.Barrio = d.Barrio
		p.CantidadDisponible = d.CantidadDisponible
		if p.CantidadDisponible == "" {
			p.CantidadDisponible = "1 unidad"
		}
		p.Disponibilidad = d.Disponibilidad
		if p.Disponibilidad == "" {
			p.Disponibilidad = "Por acordar"
		}
	}

	if d.Tipo == "servicio" {
		p.Municipio = d.Municipio
		p.CantidadDisponible = d.CantidadDisponible
		if p.CantidadDisponible == "" {
			p.CantidadDisponible = "1 sesion"
		}
		p.Disponibilidad = d.Disponibilidad
		if p.Disponibilidad == "" {
			p.Disponibilidad = "Por acordar"
		}
	}
}

package controllers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"xchango_mid/models"
	"xchango_mid/services"
)

type ModeracionController struct {
	beego.Controller
}

// GET/api/v1/admin/usuarios
func (c *ModeracionController) Usuarios() {
	_, valido := exigirAdmin(&c.Controller, false)
	if !valido {
		return
	}

	usuarios, err := cargarUsuarios()
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.Controller, 200, usuarios)
}

// PATCH/api/v1/admin/usuarios/:id/estado
func (c *ModeracionController) CambiarEstadoUsuario() {
	_, valido := exigirAdmin(&c.Controller, false)
	if !valido {
		return
	}

	var cambio models.CambioEstado

	if !leerCuerpo(&c.Controller, &cambio) {
		return
	}

	if cambio.Estado != "activo" && cambio.Estado != "suspendido" && cambio.Estado != "advertido" {
		responderError(&c.Controller, 400, "El estado debe ser activo, suspendido o advertido.")
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

// PUT/api/v1/admin/usuarios/:id
func (c *ModeracionController) EditarUsuario() {
	_, valido := exigirAdmin(&c.Controller, false)
	if !valido {
		return
	}

	var cambios models.Usuario

	if !leerCuerpo(&c.Controller, &cambios) {
		return
	}

	crud := services.NewCrudService()

	var usuario models.Usuario
	err := crud.Obtener("usuarios", c.Ctx.Input.Param(":id"), &usuario)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	//Solo se cambian algunos cambios, el correo y los contadores no
	if strings.TrimSpace(cambios.Nombre) != "" {
		usuario.Nombre = strings.TrimSpace(cambios.Nombre)
	}
	if strings.TrimSpace(cambios.Apellido) != "" {
		usuario.Apellido = strings.TrimSpace(cambios.Apellido)
	}
	if cambios.Telefono != "" {
		usuario.Telefono = soloNumeros(cambios.Telefono)
	}
	if cambios.Ubicacion != "" {
		usuario.Ubicacion = cambios.Ubicacion
	}
	if cambios.Descripcion != "" {
		usuario.Descripcion = cambios.Descripcion
	}
	if cambios.Verificacion == "verificado" || cambios.Verificacion == "no_verificado" {
		usuario.Verificacion = cambios.Verificacion
	}
	if cambios.Estado == "activo" || cambios.Estado == "suspendido" || cambios.Estado == "advertido" {
		usuario.Estado = cambios.Estado
	}

	err = crud.Actualizar("usuarios", usuario.Id, usuario, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responderVacio(&c.Controller)
}

// GET/api/v1/admin/publicaciones (tambien las eliminadas)
func (c *ModeracionController) Publicaciones() {
	_, valido := exigirAdmin(&c.Controller, false)
	if !valido {
		return
	}

	var publicaciones []models.Publicacion
	err := services.NewCrudService().Listar("publicaciones", &publicaciones)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.Controller, 200, publicaciones)
}

// GET/api/v1/admin/reportes
func (c *ModeracionController) Reportes() {
	_, valido := exigirAdmin(&c.Controller, false)
	if !valido {
		return
	}

	var reportes []models.Reporte
	err := services.NewCrudService().Listar("reportes", &reportes)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.Controller, 200, reportes)
}

// PATCH/api/v1/admin/reportes/:id
func (c *ModeracionController) ActualizarReporte() {
	_, valido := exigirAdmin(&c.Controller, false)
	if !valido {
		return
	}

	var cambio models.CambioEstado

	if !leerCuerpo(&c.Controller, &cambio) {
		return
	}

	if cambio.Estado != "pendiente" && cambio.Estado != "aprobada" && cambio.Estado != "eliminada" {
		responderError(&c.Controller, 400, "El estado debe ser pendiente, aprobada o eliminada.")
		return
	}

	crud := services.NewCrudService()

	var reporte models.Reporte
	err := crud.Obtener("reportes", c.Ctx.Input.Param(":id"), &reporte)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	reporte.Estado = cambio.Estado

	err = crud.Actualizar("reportes", reporte.Id, reporte, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responderVacio(&c.Controller)
}

// POST/api/v1/admin/reportes/:id/eliminar
func (c *ModeracionController) EliminarPublicacion() {
	sesion, valido := exigirAdmin(&c.Controller, false)
	if !valido {
		return
	}

	var solicitud models.SolicitudEliminar

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	motivo := strings.TrimSpace(solicitud.Motivo)
	if motivo == "" {
		responderError(&c.Controller, 400, "Escribe el motivo de la eliminacion.")
		return
	}

	crud := services.NewCrudService()

	var reporte models.Reporte
	err := crud.Obtener("reportes", c.Ctx.Input.Param(":id"), &reporte)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	var publicacion models.Publicacion
	err = crud.Obtener("publicaciones", reporte.PublicacionId, &publicacion)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	//Primero la publicacion, despues el reporte y al final el registro
	publicacion.Estado = "eliminada"
	err = crud.Actualizar("publicaciones", publicacion.Id, publicacion, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	reporte.Estado = "eliminada"
	err = crud.Actualizar("reportes", reporte.Id, reporte, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	registro := models.PublicacionEliminada{
		PublicacionId: publicacion.Id,
		UsuarioId:     publicacion.UsuarioId,
		AdminId:       sesion.ID,
		Motivo:        motivo,
		Fecha:         ahora(),
	}

	err = crud.Crear("publicaciones-eliminadas", registro, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responderVacio(&c.Controller)
}

// GET/api/v1/admin/eliminadas
func (c *ModeracionController) Eliminadas() {
	_, valido := exigirAdmin(&c.Controller, false)
	if !valido {
		return
	}

	var eliminadas []models.PublicacionEliminada
	err := services.NewCrudService().Listar("publicaciones-eliminadas", &eliminadas)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.Controller, 200, eliminadas)
}

// GET/api/v1/admin/documentos
func (c *ModeracionController) Documentos() {
	_, valido := exigirAdmin(&c.Controller, false)
	if !valido {
		return
	}

	var documentos []models.Documento
	err := services.NewCrudService().Listar("documentos", &documentos)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}
	responder(&c.Controller, 200, documentos)
}

// PATCH/api/v1/admin/documentos/:id
func (c *ModeracionController) RevisarDocumento() {
	sesion, valido := exigirAdmin(&c.Controller, false)
	if !valido {
		return
	}

	var revision models.RevisionDocumento

	if !leerCuerpo(&c.Controller, &revision) {
		return
	}

	if revision.Estado != "aprobado" && revision.Estado != "rechazado" {
		responderError(&c.Controller, 400, "El estado debe ser aprobado o rechazado.")
		return
	}

	crud := services.NewCrudService()

	var documento models.Documento
	err := crud.Obtener("documentos", c.Ctx.Input.Param(":id"), &documento)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	revisor := sesion.ID
	documento.Estado = revision.Estado
	documento.Observacion = revision.Observacion
	documento.RevisadoPor = &revisor

	err = crud.Actualizar("documentos", documento.Id, documento, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	// si lo aprueban la cuenta queda verificada
	if revision.Estado == "aprobado" {
		var usuario models.Usuario
		err = crud.Obtener("usuarios", documento.UsuarioId, &usuario)
		if err != nil {
			errorCrud(&c.Controller, err)
			return
		}

		usuario.Verificacion = "verificado"

		err = crud.Actualizar("usuarios", usuario.Id, usuario, nil)
		if err != nil {
			errorCrud(&c.Controller, err)
			return
		}
	}

	responderVacio(&c.Controller)
}

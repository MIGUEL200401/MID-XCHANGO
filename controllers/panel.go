package controllers

import (
	"math"
	"net/mail"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"

	"xchango_mid/models"
	"xchango_mid/services"
)

// Administra las metricas del historial, alertas y metricas del panel
type PanelController struct {
	beego.Controller
}

func comparar(actual int, anterior int) models.Comparativa {
	porcentaje := 0.0

	if anterior > 0 {
		porcentaje = math.Round(float64(actual-anterior)/float64(anterior)*1000) / 10
	} else if actual > 0 {
		porcentaje = 100
	}

	return models.Comparativa{Actual: actual, Anterior: anterior, Porcentaje: porcentaje}
}

// crea un nuevo administrador o moderador
func (c *PanelController) CrearAdmin() {
	_, ok := exigirAdmin(&c.Controller, true)
	if !ok {
		return
	}

	var nuevo models.AdminUsuario

	if !leerCuerpo(&c.Controller, &nuevo) {
		return
	}

	nuevo.Nombre = strings.TrimSpace(nuevo.Nombre)
	nuevo.Email = strings.ToLower(strings.TrimSpace(nuevo.Email))

	if nuevo.Nombre == "" {
		responderError(&c.Controller, 400, "El nombre es obligatorio.")
		return
	}

	_, err := mail.ParseAddress(nuevo.Email)
	if err != nil {
		responderError(&c.Controller, 400, "El correo no tiene un formato válido.")
		return
	}

	if len(nuevo.Password) < 8 {
		responderError(&c.Controller, 400, "La contraseña debe tener al menos 8 caracteres.")
		return
	}

	if nuevo.Rol != "admin" && nuevo.Rol != "moderador" {
		responderError(&c.Controller, 400, "El rol debe ser admin o moderador.")
		return
	}

	if nuevo.Estado != "activo" && nuevo.Estado != "inactivo" {
		nuevo.Estado = "activo"
	}

	crud := services.NewCrudService()

	var admins []models.AdminUsuario
	err = crud.Listar("admins", &admins)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	for _, a := range admins {
		if strings.ToLower(a.Email) == nuevo.Email {
			responderError(&c.Controller, 409, "Ya existe un usuario interno con ese correo.")
			return
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(nuevo.Password), bcrypt.DefaultCost)
	if err != nil {
		responderError(&c.Controller, 500, "No se pudo cifrar la contraseña.")
		return
	}

	nuevo.Password = string(hash)
	nuevo.Avatar = "https://i.pravatar.cc/80?u=" + nuevo.Email
	nuevo.UltimoAcceso = ahora()

	var creado models.AdminUsuario
	err = crud.Crear("admins", nuevo, &creado)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	creado.Password = ""

	responder(&c.Controller, 201, creado)
}

// aqui se cambia si un admin queda activo o inactivo
func (c *PanelController) EstadoAdmin() {
	claims, ok := exigirAdmin(&c.Controller, true)
	if !ok {
		return
	}

	var cambio models.CambioEstado

	if !leerCuerpo(&c.Controller, &cambio) {
		return
	}

	if cambio.Estado != "activo" && cambio.Estado != "inactivo" {
		responderError(&c.Controller, 400, "El estado debe ser activo o inactivo.")
		return
	}

	id := c.Ctx.Input.Param(":id")

	if id == claims.ID {
		responderError(&c.Controller, 400, "No puedes cambiar el estado de tu propia cuenta.")
		return
	}

	crud := services.NewCrudService()

	var admin models.AdminUsuario
	err := crud.Obtener("admins", id, &admin)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	admin.Estado = cambio.Estado

	err = crud.Actualizar("admins", id, admin, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responderVacio(&c.Controller)
}

// borrar un admin del panel
func (c *PanelController) EliminarAdmin() {
	claims, ok := exigirAdmin(&c.Controller, true)
	if !ok {
		return
	}

	id := c.Ctx.Input.Param(":id")

	if id == claims.ID {
		responderError(&c.Controller, 400, "No puedes eliminar tu propia cuenta.")
		return
	}

	err := services.NewCrudService().Borrar("admins", id)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responderVacio(&c.Controller)
}

// ver todo lo que han hecho los admins
func (c *PanelController) Historial() {
	_, ok := exigirAdmin(&c.Controller, false)
	if !ok {
		return
	}

	var acciones []models.AccionAdmin
	err := services.NewCrudService().Listar("historial-acciones", &acciones)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responder(&c.Controller, 200, acciones)
}

// guardar una accion nueva en el historial
func (c *PanelController) RegistrarAccion() {
	claims, ok := exigirAdmin(&c.Controller, false)
	if !ok {
		return
	}

	var accion models.AccionAdmin

	if !leerCuerpo(&c.Controller, &accion) {
		return
	}

	if strings.TrimSpace(accion.Accion) == "" {
		responderError(&c.Controller, 400, "La acción es obligatoria.")
		return
	}

	// el admin sale del token, no de lo que mande el cliente
	accion.Id = ""
	accion.AdminId = claims.ID
	accion.Fecha = ahora()

	var creada models.AccionAdmin
	err := services.NewCrudService().Crear("historial-acciones", accion, &creada)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responder(&c.Controller, 201, creada)
}
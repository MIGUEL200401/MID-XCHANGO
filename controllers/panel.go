package controllers

import (
	"math"
	"net/mail"
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"xchango_mid/models"
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
}

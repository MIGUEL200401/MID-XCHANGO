package controllers

import (
	"encoding/json"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"

	"xchango_mid/models"
	"xchango_mid/services"
)

func responder(c *beego.Controller, codigo int, datos interface{}) {
	c.Ctx.Output.SetStatus(codigo)
	c.Data["json"] = datos
	c.ServeJSON()
}

func responderError(c *beego.Controller, codigo int, mensaje string) {
	responder(c, codigo, map[string]interface{}{
		"ok":      false,
		"mensaje": mensaje,
	})
}

func errorCrud(c *beego.Controller, err error) {
	if err == services.ErrNoEncontrado {
		responderError(c, 404, "No se encontró lo que buscas.")
		return
	}
	responderError(c, 502, "El servicio de datos no respondió: "+err.Error())
}

// lee el JSON del cuerpo, si viene mal ya responde el error
func leerCuerpo(c *beego.Controller, destino interface{}) bool {
	cuerpo := c.Ctx.Input.RequestBody

	if len(cuerpo) == 0 {
		responderError(c, 400, "El cuerpo de la solicitud está vacío.")
		return false
	}

	err := json.Unmarshal(cuerpo, destino)
	if err != nil {
		responderError(c, 400, "El cuerpo de la solicitud no tiene un formato válido.")
		return false
	}

	return true
}

func hoy() string {
	return time.Now().Format("2006-01-02")
}

func ahora() string {
	return time.Now().Format("2006-01-02T15:04:05")
}

// lee el token pero no responde error, solo devuelve false si no hay token o es invalido
func leerToken(c *beego.Controller) (*services.Claims, bool) {
	cabecera := c.Ctx.Input.Header("Authorization")

	if !strings.HasPrefix(cabecera, "Bearer ") {
		return nil, false
	}

	claims, err := services.ValidarToken(strings.TrimPrefix(cabecera, "Bearer "))
	if err != nil {
		return nil, false
	}

	return claims, true
}

// deja pasar solo a un usuario con sesion
func exigirUsuario(c *beego.Controller) (*services.Claims, bool) {
	claims, ok := leerToken(c)

	if !ok || claims.Rol != "usuario" {
		responderError(c, 401, "Inicia sesión para continuar.")
		return nil, false
	}

	return claims, true
}

// soloAdmin = true deja pasar unicamente al administrador, no al moderador
func exigirAdmin(c *beego.Controller, soloAdmin bool) (*services.Claims, bool) {
	claims, ok := leerToken(c)

	if !ok || (claims.Rol != "admin" && claims.Rol != "moderador") {
		responderError(c, 401, "Necesitas iniciar sesión como administrador.")
		return nil, false
	}

	if soloAdmin && claims.Rol != "admin" {
		responderError(c, 403, "Solo un administrador puede hacer esto.")
		return nil, false
	}

	return claims, true
}

// un usuario suspendido no puede publicar ni proponer trueques
func usuarioPuedeActuar(c *beego.Controller, id string) bool {
	var usuario models.Usuario

	err := services.NewCrudService().Obtener("usuarios", id, &usuario)
	if err != nil {
		errorCrud(c, err)
		return false
	}

	if usuario.Estado == "suspendido" {
		responderError(c, 403, "Tu cuenta está suspendida.")
		return false
	}

	return true
}

package controllers

import (
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"

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
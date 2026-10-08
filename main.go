package main

import (
	"fmt"

	_ "xchango_mid/routers"
	"xchango_mid/services"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// si los PIN son por Gmail, se autoriza al inicio y no en medio de una peticion
	if beego.AppConfig.DefaultBool("enviar_correo", false) {
		_, err := services.NuevoServicioGmail()
		if err != nil {
			fmt.Println("Gmail no está listo:", err)
			return
		}
	}

	beego.Run()
}
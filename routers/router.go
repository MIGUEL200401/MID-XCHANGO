package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"

	"xchango_mid/controllers"
)

func init() {
	// funcion para permitir cors, para q el cliente puede hacer peticiones desde su puerto
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{beego.AppConfig.DefaultString("origen_cliente", "http://localhost:4200")},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	beego.Router("/api/v1/auth/pin", &controllers.AuthController{}, "post:EnviarPin")
}
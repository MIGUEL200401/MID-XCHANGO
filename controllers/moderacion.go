package controllers

import (
	"strings"

	"github.com/astaxie/beego"
	beego "github.com/beego/beego/v2/server/web"

	"xchango_mid/models"
	"xchango_mid/services"
)

type ModeracionController struct {
	beego.Controller
}
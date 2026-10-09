package controllers

import (
	"crypto/rand"
	"fmt"
	"html"
	"math/big"
	"net/mail"
	"os"
	"strings"
	"sync"
	"time"

	beego "github.com/beego/beego/v2/server/web"

	"xchango_mid/models"
	"xchango_mid/services"
)

// El pin se guarda en la ram, no en la base 
var pines = make(map[string]models.CodigoPin)

// candado para que no se crucen las peticiones 
var candadoPines sync.RWMutex

const maximoIntentos = 5

type AuthController struct {
	beego.Controller
}

// POST /api/v1/auth/pin
func (c *AuthController) EnviarPin() {
	var solicitud models.SolicitudPin

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	correo := strings.ToLower(strings.TrimSpace(solicitud.Correo))

	if correo == "" {
		responderError(&c.Controller, 400, "El correo es obligatorio.")
		return
	}

	_, err := mail.ParseAddress(correo)
	if err != nil {
		responderError(&c.Controller, 400, "El correo no tiene un formato válido.")
		return
	}

	// no se puede pedir otro PIN antes de 30 segundos
	candadoPines.RLock()
	anterior, hayAnterior := pines[correo]
	candadoPines.RUnlock()

	if hayAnterior && time.Since(anterior.Creado) < 30*time.Second {
		responderError(&c.Controller, 429, "Espera unos segundos para pedir otro PIN.")
		return
	}

	// miramos si el correo ya tiene cuenta
	var usuarios []models.Usuario
	err = services.NewCrudService().Listar("usuarios", &usuarios)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	requiereRegistro := true
	nombre := ""
	for _, u := range usuarios {
		if strings.ToLower(u.Email) == correo {
			requiereRegistro = false
			nombre = u.Nombre
			if u.Estado == "suspendido" {
				responderError(&c.Controller, 403, "Tu cuenta está suspendida.")
				return
			}
		}
	}

	pin, err := generarPin()
	if err != nil {
		responderError(&c.Controller, 500, "No se pudo generar el PIN.")
		return
	}

	candadoPines.Lock()
	pines[correo] = models.CodigoPin{
		Correo:   correo,
		Pin:      pin,
		Creado:   time.Now(),
		Vence:    time.Now().Add(10 * time.Minute),
		Intentos: 0,
		Activo:   true,
	}
	candadoPines.Unlock()

	err = mandarPin(correo, pin, nombre)
	if err != nil {
		candadoPines.Lock()
		delete(pines, correo)
		candadoPines.Unlock()

		responderError(&c.Controller, 500, "No se pudo enviar el correo: "+err.Error())
		return
	}

	responder(&c.Controller, 200, map[string]interface{}{
		"ok":               true,
		"mensaje":          "Te enviamos un PIN al correo.",
		"requiereRegistro": requiereRegistro,
	})
}

func generarPin() (string, error) {
	numero, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%04d", numero.Int64()), nil
}

func mandarPin(correo string, pin string, nombre string) error {
	if !beego.AppConfig.DefaultBool("enviar_correo", false) {
		fmt.Println("PIN para", correo, "->", pin)
		return nil
	}

	servicioGmail, err := services.NuevoServicioGmail()
	if err != nil {
		return err
	}

	plantilla, err := os.ReadFile("templates/correo_pin.html")
	if err != nil {
		return err
	}


	logo, _ := os.ReadFile("templates/logo.png")

	// si el usuario ya tiene cuenta se saluda por su nombre
	saludo := "Hola,"
	if strings.TrimSpace(nombre) != "" {
		saludo = "Hola, " + html.EscapeString(strings.Fields(nombre)[0]) + ","
	}

	cuerpo := strings.ReplaceAll(string(plantilla), "{{PIN}}", pin)
	cuerpo = strings.ReplaceAll(cuerpo, "{{SALUDO}}", saludo)

	_, err = servicioGmail.EnviarCorreo(correo, "Tu PIN de XchanGo", cuerpo, logo)
	return err
}

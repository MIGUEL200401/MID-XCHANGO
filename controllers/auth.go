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
	"golang.org/x/crypto/bcrypt"

	"xchango_mid/models"
	"xchango_mid/services"
)

// El pin se guarda en la ram, no en la base
var pines = make(map[string]models.CodigoPin)

// candado para que no se crucen las peticiones
var candadoPines sync.RWMutex

const maximoIntentos = 5

// la foto de perfil llega como texto (base64), este es el maximo de caracteres
const maximoAvatar = 300000

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

// POST /api/v1/auth/verificar
func (c *AuthController) VerificarPin() {
	var solicitud models.SolicitudVerificar

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	correo := strings.ToLower(strings.TrimSpace(solicitud.Correo))
	pin := strings.TrimSpace(solicitud.Pin)

	if correo == "" || pin == "" {
		responderError(&c.Controller, 400, "El correo y el PIN son obligatorios.")
		return
	}

	candadoPines.RLock()
	guardado, existe := pines[correo]
	candadoPines.RUnlock()

	if !existe {
		responderError(&c.Controller, 404, "No encontramos un PIN para ese correo. Pídelo de nuevo.")
		return
	}

	if !guardado.Activo {
		responderError(&c.Controller, 401, "Ese PIN ya no sirve. Pide uno nuevo.")
		return
	}

	if time.Now().After(guardado.Vence) {
		candadoPines.Lock()
		guardado.Activo = false
		pines[correo] = guardado
		candadoPines.Unlock()

		responderError(&c.Controller, 401, "El PIN se venció. Pide uno nuevo.")
		return
	}

	if pin != guardado.Pin {
		candadoPines.Lock()
		guardado.Intentos++
		if guardado.Intentos >= maximoIntentos {
			guardado.Activo = false
		}
		pines[correo] = guardado
		candadoPines.Unlock()

		if !guardado.Activo {
			responderError(&c.Controller, 401, "Demasiados intentos. Pide un PIN nuevo.")
			return
		}

		quedan := maximoIntentos - guardado.Intentos
		responderError(&c.Controller, 401, fmt.Sprintf("El PIN no coincide. Te quedan %d intentos.", quedan))
		return
	}

	// PIN bueno, ya no se puede volver a usar
	candadoPines.Lock()
	guardado.Activo = false
	pines[correo] = guardado
	candadoPines.Unlock()

	var usuarios []models.Usuario
	err := services.NewCrudService().Listar("usuarios", &usuarios)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	// si tiene cuenta el token es de usuario, si no solo sirve para registrarse
	id := ""
	rol := "registro"
	for _, u := range usuarios {
		if strings.ToLower(u.Email) == correo {
			id = u.Id
			rol = "usuario"
		}
	}

	token, err := services.GenerarToken(id, correo, rol)
	if err != nil {
		responderError(&c.Controller, 500, "No se pudo crear el token.")
		return
	}

	responder(&c.Controller, 200, map[string]interface{}{
		"ok":               true,
		"token":            token,
		"requiereRegistro": rol == "registro",
	})
}

// POST /api/v1/auth/registro
func (c *AuthController) Registro() {
	claims, ok := leerToken(&c.Controller)
	if !ok || claims.Rol != "registro" {
		responderError(&c.Controller, 401, "Primero verifica tu correo con el PIN.")
		return
	}

	var solicitud models.SolicitudRegistro

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	correo := strings.ToLower(strings.TrimSpace(solicitud.Email))
	nombre := strings.TrimSpace(solicitud.Nombre)
	apellido := strings.TrimSpace(solicitud.Apellido)

	if correo != claims.Correo {
		responderError(&c.Controller, 403, "El correo no es el que verificaste.")
		return
	}

	if nombre == "" {
		responderError(&c.Controller, 400, "El nombre es obligatorio.")
		return
	}

	if apellido == "" {
		responderError(&c.Controller, 400, "El apellido es obligatorio.")
		return
	}

	mensaje := validarEdad(solicitud.FechaNacimiento)
	if mensaje != "" {
		responderError(&c.Controller, 400, mensaje)
		return
	}

	direccion := strings.TrimSpace(solicitud.Direccion)
	barrio := strings.TrimSpace(solicitud.Barrio)

	if direccion == "" {
		responderError(&c.Controller, 400, "La dirección es obligatoria.")
		return
	}

	if barrio == "" {
		responderError(&c.Controller, 400, "El barrio es obligatorio.")
		return
	}

	crud := services.NewCrudService()

	var usuarios []models.Usuario
	err := crud.Listar("usuarios", &usuarios)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	for _, u := range usuarios {
		if strings.ToLower(u.Email) == correo {
			responderError(&c.Controller, 409, "Ese correo ya está registrado.")
			return
		}
	}

	ubicacion := strings.TrimSpace(solicitud.Ubicacion)
	if ubicacion == "" {
		ubicacion = "Yopal"
	}

	// si no mandan foto se le pone una al azar
	avatar := avatarAlAzar()
	if solicitud.Avatar != "" {
		if !strings.HasPrefix(solicitud.Avatar, "data:image/") {
			responderError(&c.Controller, 400, "La foto de perfil no es válida.")
			return
		}
		if len(solicitud.Avatar) > maximoAvatar {
			responderError(&c.Controller, 400, "La foto de perfil es muy pesada.")
			return
		}
		avatar = solicitud.Avatar
	}

	nuevo := models.Usuario{
		Nombre:          nombre,
		Apellido:        apellido,
		FechaNacimiento: strings.TrimSpace(solicitud.FechaNacimiento),
		Direccion:       direccion,
		Barrio:          barrio,
		Email:           correo,
		Telefono:        soloNumeros(solicitud.Telefono),
		Estado:          "activo",
		FechaRegistro:   hoy(),
		Ubicacion:       ubicacion,
		NivelActividad:  "bajo",
		Avatar:          avatar,
		Verificacion:    "no_verificado",
	}

	var creado models.Usuario
	err = crud.Crear("usuarios", nuevo, &creado)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	token, err := services.GenerarToken(creado.Id, creado.Email, "usuario")
	if err != nil {
		responderError(&c.Controller, 500, "No se pudo crear el token.")
		return
	}

	responder(&c.Controller, 201, map[string]interface{}{
		"token":   token,
		"usuario": creado,
	})
}

// revisa la fecha de nacimiento (aaaa-mm-dd) y que la persona tenga 18 años o mas
// devuelve el mensaje de error, o "" si esta bien
func validarEdad(fecha string) string {
	fecha = strings.TrimSpace(fecha)

	if fecha == "" {
		return "La fecha de nacimiento es obligatoria."
	}

	nacimiento, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return "La fecha de nacimiento debe tener el formato aaaa-mm-dd."
	}

	ahoraMismo := time.Now()

	if nacimiento.After(ahoraMismo) {
		return "La fecha de nacimiento no puede ser futura."
	}

	edad := ahoraMismo.Year() - nacimiento.Year()

	// si todavia no ha cumplido años este año, se resta uno
	if ahoraMismo.Month() < nacimiento.Month() ||
		(ahoraMismo.Month() == nacimiento.Month() && ahoraMismo.Day() < nacimiento.Day()) {
		edad--
	}

	if edad < 18 {
		return "Debes ser mayor de 18 años para registrarte."
	}

	return ""
}

// deja solo los digitos de un telefono
func soloNumeros(texto string) string {
	resultado := ""
	for _, letra := range texto {
		if letra >= '0' && letra <= '9' {
			resultado += string(letra)
		}
	}
	return resultado
}

// foto de perfil al azar mientras el usuario no sube la suya
func avatarAlAzar() string {
	fotos := []int{12, 5, 15, 33, 44, 45, 47, 48, 60, 65, 68, 32, 20, 25}
	numero, err := rand.Int(rand.Reader, big.NewInt(int64(len(fotos))))
	if err != nil {
		return "https://i.pravatar.cc/80?img=12"
	}
	return fmt.Sprintf("https://i.pravatar.cc/80?img=%d", fotos[numero.Int64()])
}

// GET /api/v1/auth/perfil
func (c *AuthController) Perfil() {
	claims, ok := exigirUsuario(&c.Controller)
	if !ok {
		return
	}

	usuarios, err := cargarUsuarios()
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	for _, u := range usuarios {
		if u.Id == claims.ID {
			responder(&c.Controller, 200, u)
			return
		}
	}

	responderError(&c.Controller, 404, "No encontramos tu cuenta.")
}

// PUT /api/v1/usuarios/perfil
func (c *AuthController) ActualizarPerfil() {
	claims, ok := exigirUsuario(&c.Controller)
	if !ok {
		return
	}

	var cambios models.CambiosPerfil

	if !leerCuerpo(&c.Controller, &cambios) {
		return
	}

	crud := services.NewCrudService()

	var usuario models.Usuario
	err := crud.Obtener("usuarios", claims.ID, &usuario)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	// solo se cambia lo que venga con algo escrito
	if strings.TrimSpace(cambios.Nombre) != "" {
		usuario.Nombre = strings.TrimSpace(cambios.Nombre)
	}
	if strings.TrimSpace(cambios.Apellido) != "" {
		usuario.Apellido = strings.TrimSpace(cambios.Apellido)
	}
	if strings.TrimSpace(cambios.Telefono) != "" {
		usuario.Telefono = soloNumeros(cambios.Telefono)
	}
	if strings.TrimSpace(cambios.Ubicacion) != "" {
		usuario.Ubicacion = strings.TrimSpace(cambios.Ubicacion)
	}
	if strings.TrimSpace(cambios.Direccion) != "" {
		usuario.Direccion = strings.TrimSpace(cambios.Direccion)
	}
	if strings.TrimSpace(cambios.Barrio) != "" {
		usuario.Barrio = strings.TrimSpace(cambios.Barrio)
	}
	if cambios.Descripcion != "" {
		usuario.Descripcion = cambios.Descripcion
	}
	if cambios.Avatar != "" {
		if len(cambios.Avatar) > maximoAvatar {
			responderError(&c.Controller, 400, "La foto de perfil es muy pesada.")
			return
		}
		usuario.Avatar = cambios.Avatar
	}

	err = crud.Actualizar("usuarios", usuario.Id, usuario, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	usuarios, err := cargarUsuarios()
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	for _, u := range usuarios {
		if u.Id == usuario.Id {
			responder(&c.Controller, 200, u)
			return
		}
	}

	responder(&c.Controller, 200, usuario)
}

// GET /api/v1/usuarios
func (c *AuthController) ListarUsuarios() {
	var usuarios []models.Usuario
	err := services.NewCrudService().Listar("usuarios", &usuarios)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	// esto lo ve cualquiera, asi que se esconden los datos privados
	claims, haySesion := leerToken(&c.Controller)

	for i := range usuarios {
		esMio := haySesion && claims.Rol == "usuario" && claims.ID == usuarios[i].Id

		if !esMio {
			usuarios[i].Email = ""
		}
		if !haySesion {
			usuarios[i].Telefono = ""
		}
	}

	responder(&c.Controller, 200, usuarios)
}

// POST /api/v1/auth/admin/login
func (c *AuthController) LoginAdmin() {
	var datos models.LoginAdmin

	if !leerCuerpo(&c.Controller, &datos) {
		return
	}

	email := strings.ToLower(strings.TrimSpace(datos.Email))

	crud := services.NewCrudService()

	var admins []models.AdminUsuario
	err := crud.Listar("admins", &admins)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	var admin *models.AdminUsuario
	for i := range admins {
		if strings.ToLower(admins[i].Email) == email {
			admin = &admins[i]
		}
	}

	if admin == nil {
		responderError(&c.Controller, 401, "Correo o contraseña incorrectos.")
		return
	}

	passwordBien := false
	if strings.HasPrefix(admin.Password, "$2") {
		// ya esta cifrada
		passwordBien = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(datos.Password)) == nil
	} else if admin.Password != "" && admin.Password == datos.Password {
		// estaba en texto plano en el CRUD, la guardamos cifrada
		passwordBien = true
		hash, err := bcrypt.GenerateFromPassword([]byte(datos.Password), bcrypt.DefaultCost)
		if err == nil {
			admin.Password = string(hash)
		}
	}

	if !passwordBien {
		responderError(&c.Controller, 401, "Correo o contraseña incorrectos.")
		return
	}

	if admin.Estado == "inactivo" {
		responderError(&c.Controller, 403, "Esta cuenta está inactiva.")
		return
	}

	admin.UltimoAcceso = ahora()
	err = crud.Actualizar("admins", admin.Id, admin, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	token, err := services.GenerarToken(admin.Id, admin.Email, admin.Rol)
	if err != nil {
		responderError(&c.Controller, 500, "No se pudo crear el token.")
		return
	}

	// la contraseña nunca sale del MID
	admin.Password = ""

	responder(&c.Controller, 200, map[string]interface{}{
		"token": token,
		"admin": admin,
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
	saludo := "Hola"
	if strings.TrimSpace(nombre) != "" {
		saludo = "Hola, " + html.EscapeString(strings.Fields(nombre)[0])
	}

	cuerpo := strings.ReplaceAll(string(plantilla), "{{PIN}}", pin)
	cuerpo = strings.ReplaceAll(cuerpo, "{{SALUDO}}", saludo)

	_, err = servicioGmail.EnviarCorreo(correo, "Tu PIN de XchanGo", cuerpo, logo)
	return err
}

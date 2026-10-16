package controllers

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"xchango_mid/models"
	"xchango_mid/services"
)

type DocumentosController struct {
	beego.Controller
}

// tipos de documento que se aceptan
var tiposDocumento = []string{"Cédula de ciudadanía", "Cédula de extranjería"}

// el archivo llega como texto (imagen en base64), este es el maximo de caracteres
const maximoArchivo = 3000000

// POST /api/v1/documentos
// El usuario sube su documento y queda pendiente hasta que un administrador lo revise
func (c *DocumentosController) Subir() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	var solicitud models.SolicitudDocumento

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	tipo := strings.TrimSpace(solicitud.TipoDocumento)
	numero := strings.TrimSpace(solicitud.Numero)
	archivo := strings.TrimSpace(solicitud.Archivo)

	tipoValido := false
	for _, t := range tiposDocumento {
		if t == tipo {
			tipoValido = true
		}
	}

	if !tipoValido {
		responderError(&c.Controller, 400, "El tipo de documento debe ser Cédula de ciudadanía o Cédula de extranjería.")
		return
	}

	if numero == "" {
		responderError(&c.Controller, 400, "El número del documento es obligatorio.")
		return
	}

	if len(numero) > 30 {
		responderError(&c.Controller, 400, "El número del documento es muy largo.")
		return
	}

	if archivo == "" {
		responderError(&c.Controller, 400, "Debes adjuntar la foto del documento.")
		return
	}

	if len(archivo) > maximoArchivo {
		responderError(&c.Controller, 400, "La foto del documento es muy pesada.")
		return
	}

	crud := services.NewCrudService()

	var usuario models.Usuario
	err := crud.Obtener("usuarios", sesion.ID, &usuario)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	if usuario.Verificacion == "verificado" {
		responderError(&c.Controller, 400, "Tu cuenta ya está verificada.")
		return
	}

	var documentos []models.Documento
	err = crud.Listar("documentos", &documentos)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	// solo se puede tener un documento en revision a la vez
	for _, d := range documentos {
		if d.UsuarioId == sesion.ID && d.Estado == "pendiente" {
			responderError(&c.Controller, 409, "Ya tienes un documento en revisión.")
			return
		}
	}

	nuevo := models.Documento{
		UsuarioId:     sesion.ID,
		TipoDocumento: tipo,
		Numero:        numero,
		Archivo:       archivo,
		FechaSubida:   hoy(),
		Estado:        "pendiente",
		Observacion:   "",
	}

	var creado models.Documento
	err = crud.Crear("documentos", nuevo, &creado)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	// la foto no se devuelve para no responder algo tan pesado
	creado.Archivo = ""

	responder(&c.Controller, 201, creado)
}

// GET /api/v1/documentos/mio
// El usuario ve en que estado van sus documentos (pendiente, aprobado o rechazado)
func (c *DocumentosController) Mio() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	var documentos []models.Documento

	err := services.NewCrudService().Listar("documentos", &documentos)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	mios := []models.Documento{}
	for _, d := range documentos {
		if d.UsuarioId == sesion.ID {
			d.Archivo = ""
			mios = append(mios, d)
		}
	}

	responder(&c.Controller, 200, mios)
}

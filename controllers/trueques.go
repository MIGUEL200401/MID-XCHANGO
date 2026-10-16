package controllers

import (
	"sort"
	"strings"

	beego "github.com/beego/beego/v2/server/web"

	"xchango_mid/models"
	"xchango_mid/services"
)

type TruequesController struct {
	beego.Controller
}

// GET /api/v1/trueques
func (c *TruequesController) Listar() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	var todos []models.Trueque
	err := services.NewCrudService().Listar("trueques", &todos)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	mios := []models.Trueque{}
	for _, t := range todos {
		if t.SolicitanteId == sesion.ID || t.PropietarioId == sesion.ID {
			mios = append(mios, t)
		}
	}

	responder(&c.Controller, 200, mios)
}

// POST /api/v1/trueques
func (c *TruequesController) Solicitar() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	var solicitud models.SolicitudTrueque

	if !leerCuerpo(&c.Controller, &solicitud) {
		return
	}

	if !usuarioPuedeActuar(&c.Controller, sesion.ID) {
		return
	}

	crud := services.NewCrudService()

	var publicacion models.Publicacion
	err := crud.Obtener("publicaciones", solicitud.PublicacionId, &publicacion)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	if publicacion.UsuarioId == sesion.ID {
		responderError(&c.Controller, 400, "Esta publicación es tuya.")
		return
	}

	if publicacion.Estado != "activa" {
		responderError(&c.Controller, 400, "Esta publicación no está disponible.")
		return
	}

	var existentes []models.Trueque
	err = crud.Listar("trueques", &existentes)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	for _, t := range existentes {
		abierto := t.Estado == "pendiente" || t.Estado == "aceptado"
		if t.PublicacionId == publicacion.Id && t.SolicitanteId == sesion.ID && abierto {
			responderError(&c.Controller, 409, "Ya tienes una propuesta abierta para esta publicación.")
			return
		}
	}

	ofrece := strings.TrimSpace(solicitud.Ofrece)
	if ofrece == "" {
		ofrece = "Propuesta por acordar"
	}

	nuevo := models.Trueque{
		SolicitanteId:  sesion.ID,
		PropietarioId:  publicacion.UsuarioId,
		PublicacionId:  publicacion.Id,
		Ofrece:         ofrece,
		Busca:          publicacion.Titulo,
		Estado:         "pendiente",
		FechaSolicitud: hoy(),
	}

	var creado models.Trueque
	err = crud.Crear("trueques", nuevo, &creado)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	var solicitante models.Usuario
	err = crud.Obtener("usuarios", sesion.ID, &solicitante)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	err = crearNotificacion(
		publicacion.UsuarioId,
		"solicitud",
		"Nueva solicitud de trueque",
		solicitante.Nombre+" "+solicitante.Apellido+" quiere truequear por tu publicación \""+publicacion.Titulo+"\".",
		creado.Id,
	)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responder(&c.Controller, 201, creado)
}

// PATCH /api/v1/trueques/:id/aceptar
func (c *TruequesController) Aceptar() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	crud := services.NewCrudService()

	trueque, valido := buscarTrueque(&c.Controller, c.Ctx.Input.Param(":id"))
	if !valido {
		return
	}

	if trueque.PropietarioId != sesion.ID {
		responderError(&c.Controller, 403, "Solo el dueño de la publicación puede aceptar.")
		return
	}

	if trueque.Estado != "pendiente" {
		responderError(&c.Controller, 400, "Esta solicitud ya fue respondida.")
		return
	}

	trueque.Estado = "aceptado"

	err := crud.Actualizar("trueques", trueque.Id, trueque, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	err = crearNotificacion(
		trueque.SolicitanteId,
		"aceptada",
		"Tu solicitud fue aceptada",
		"Aceptaron tu propuesta. Continúa la negociación por WhatsApp.",
		trueque.Id,
	)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responder(&c.Controller, 200, trueque)
}

// PATCH /api/v1/trueques/:id/rechazar
func (c *TruequesController) Rechazar() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	crud := services.NewCrudService()

	trueque, valido := buscarTrueque(&c.Controller, c.Ctx.Input.Param(":id"))
	if !valido {
		return
	}

	if trueque.PropietarioId != sesion.ID {
		responderError(&c.Controller, 403, "Solo el dueño de la publicación puede rechazar.")
		return
	}

	if trueque.Estado != "pendiente" {
		responderError(&c.Controller, 400, "Esta solicitud ya fue respondida.")
		return
	}

	fecha := hoy()
	trueque.Estado = "rechazado"
	trueque.FechaCierre = &fecha

	err := crud.Actualizar("trueques", trueque.Id, trueque, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	var publicacion models.Publicacion
	titulo := "la publicación"
	if crud.Obtener("publicaciones", trueque.PublicacionId, &publicacion) == nil {
		titulo = publicacion.Titulo
	}

	err = crearNotificacion(
		trueque.SolicitanteId,
		"rechazada",
		"Tu solicitud fue rechazada",
		"Tu propuesta por \""+titulo+"\" fue rechazada.",
		trueque.Id,
	)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responder(&c.Controller, 200, trueque)
}

// PATCH /api/v1/trueques/:id/confirmar
func (c *TruequesController) Confirmar() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	crud := services.NewCrudService()

	trueque, valido := buscarTrueque(&c.Controller, c.Ctx.Input.Param(":id"))
	if !valido {
		return
	}

	if trueque.Estado != "aceptado" {
		responderError(&c.Controller, 400, "Solo se confirman los trueques aceptados.")
		return
	}

	if trueque.SolicitanteId == sesion.ID {
		trueque.ConfirmadoSolicitante = true
	} else if trueque.PropietarioId == sesion.ID {
		trueque.ConfirmadoPropietario = true
	} else {
		responderError(&c.Controller, 403, "Este trueque no es tuyo.")
		return
	}

	// cuando confirman los dos queda completado
	if trueque.ConfirmadoSolicitante && trueque.ConfirmadoPropietario {
		fecha := hoy()
		trueque.Estado = "completado"
		trueque.FechaCierre = &fecha

		var publicacion models.Publicacion
		err := crud.Obtener("publicaciones", trueque.PublicacionId, &publicacion)
		if err == nil && publicacion.Tipo != "bien_digital" {
			publicacion.Estado = "finalizada"
			err = crud.Actualizar("publicaciones", publicacion.Id, publicacion, nil)
			if err != nil {
				errorCrud(&c.Controller, err)
				return
			}
		}
	}

	err := crud.Actualizar("trueques", trueque.Id, trueque, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responder(&c.Controller, 200, trueque)
}

// PATCH /api/v1/trueques/:id/calificar
func (c *TruequesController) Calificar() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	var calificacion models.Calificacion

	if !leerCuerpo(&c.Controller, &calificacion) {
		return
	}

	if calificacion.Nota < 1 || calificacion.Nota > 5 {
		responderError(&c.Controller, 400, "La nota debe ser de 1 a 5.")
		return
	}

	trueque, valido := buscarTrueque(&c.Controller, c.Ctx.Input.Param(":id"))
	if !valido {
		return
	}

	if trueque.Estado != "completado" {
		responderError(&c.Controller, 400, "Solo se califican los trueques completados.")
		return
	}

	nota := calificacion.Nota

	if trueque.SolicitanteId == sesion.ID {
		trueque.CalificacionSolicitante = &nota
	} else if trueque.PropietarioId == sesion.ID {
		trueque.CalificacionPropietario = &nota
	} else {
		responderError(&c.Controller, 403, "Este trueque no es tuyo.")
		return
	}

	err := services.NewCrudService().Actualizar("trueques", trueque.Id, trueque, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responder(&c.Controller, 200, trueque)
}

// GET /api/v1/notificaciones
func (c *TruequesController) Notificaciones() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	var todas []models.Notificacion
	err := services.NewCrudService().Listar("notificaciones", &todas)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	mias := []models.Notificacion{}
	for _, n := range todas {
		if n.UsuarioId == sesion.ID {
			mias = append(mias, n)
		}
	}

	// las mas nuevas primero
	sort.Slice(mias, func(i, j int) bool {
		return mias[i].Fecha > mias[j].Fecha
	})

	responder(&c.Controller, 200, mias)
}

// PATCH /api/v1/notificaciones/:id/leida
func (c *TruequesController) MarcarLeida() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	crud := services.NewCrudService()

	var notificacion models.Notificacion
	err := crud.Obtener("notificaciones", c.Ctx.Input.Param(":id"), &notificacion)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	if notificacion.UsuarioId != sesion.ID {
		responderError(&c.Controller, 403, "Esta notificación no es tuya.")
		return
	}

	notificacion.Leida = true

	err = crud.Actualizar("notificaciones", notificacion.Id, notificacion, nil)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	responderVacio(&c.Controller)
}

// PATCH /api/v1/notificaciones/leidas
func (c *TruequesController) MarcarTodasLeidas() {
	sesion, valido := exigirUsuario(&c.Controller)
	if !valido {
		return
	}

	crud := services.NewCrudService()

	var todas []models.Notificacion
	err := crud.Listar("notificaciones", &todas)
	if err != nil {
		errorCrud(&c.Controller, err)
		return
	}

	for _, n := range todas {
		if n.UsuarioId == sesion.ID && !n.Leida {
			n.Leida = true
			err = crud.Actualizar("notificaciones", n.Id, n, nil)
			if err != nil {
				errorCrud(&c.Controller, err)
				return
			}
		}
	}

	responderVacio(&c.Controller)
}

func buscarTrueque(c *beego.Controller, id string) (models.Trueque, bool) {
	var trueque models.Trueque

	err := services.NewCrudService().Obtener("trueques", id, &trueque)
	if err != nil {
		errorCrud(c, err)
		return trueque, false
	}

	return trueque, true
}

func crearNotificacion(usuarioId string, tipo string, titulo string, mensaje string, referenciaId string) error {
	nueva := models.Notificacion{
		UsuarioId:    usuarioId,
		Tipo:         tipo,
		Titulo:       titulo,
		Mensaje:      mensaje,
		ReferenciaId: &referenciaId,
		Leida:        false,
		Fecha:        ahora(),
	}

	return services.NewCrudService().Crear("notificaciones", nueva, nil)
}
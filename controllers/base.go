package controllers

import (
	"encoding/json"
	"math"
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

// para las respuestas sin cuerpo
func responderVacio(c *beego.Controller) {
	c.Ctx.Output.SetStatus(204)
	c.Ctx.Output.Body([]byte(""))
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

// las fechas del CRUD pueden venir de varias formas
func leerFecha(texto string) time.Time {
	formatos := []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02"}

	for _, formato := range formatos {
		fecha, err := time.Parse(formato, texto)
		if err == nil {
			return fecha
		}
	}

	return time.Time{}
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

// ---------- usuarios con sus contadores ----------

// publicaciones, intercambios, reportes y calificacion no se guardan,
// se cuentan cada vez con lo que hay en el CRUD
func cargarUsuarios() ([]models.Usuario, error) {
	crud := services.NewCrudService()

	var usuarios []models.Usuario
	var publicaciones []models.Publicacion
	var trueques []models.Trueque
	var reportes []models.Reporte

	err := crud.Listar("usuarios", &usuarios)
	if err != nil {
		return nil, err
	}
	err = crud.Listar("publicaciones", &publicaciones)
	if err != nil {
		return nil, err
	}
	err = crud.Listar("trueques", &trueques)
	if err != nil {
		return nil, err
	}
	err = crud.Listar("reportes", &reportes)
	if err != nil {
		return nil, err
	}

	for i := range usuarios {
		id := usuarios[i].Id
		totalPublicaciones := 0
		totalIntercambios := 0
		totalReportes := 0
		suma := 0
		notas := 0

		for _, p := range publicaciones {
			if p.UsuarioId == id && p.Estado != "eliminada" {
				totalPublicaciones++
			}
			if p.UsuarioId == id {
				for _, r := range reportes {
					if r.PublicacionId == p.Id {
						totalReportes++
					}
				}
			}
		}

		for _, t := range trueques {
			if t.Estado == "completado" && (t.SolicitanteId == id || t.PropietarioId == id) {
				totalIntercambios++
			}
			// la nota que pone uno la recibe el otro
			if t.PropietarioId == id && t.CalificacionSolicitante != nil {
				suma += *t.CalificacionSolicitante
				notas++
			}
			if t.SolicitanteId == id && t.CalificacionPropietario != nil {
				suma += *t.CalificacionPropietario
				notas++
			}
		}

		usuarios[i].Publicaciones = totalPublicaciones
		usuarios[i].Intercambios = totalIntercambios
		usuarios[i].Reportes = totalReportes
		usuarios[i].TotalCalificaciones = notas
		usuarios[i].Calificacion = 0
		if notas > 0 {
			usuarios[i].Calificacion = math.Round(float64(suma)/float64(notas)*10) / 10
		}

		total := totalPublicaciones + totalIntercambios
		if total >= 8 {
			usuarios[i].NivelActividad = "alto"
		} else if total >= 3 {
			usuarios[i].NivelActividad = "medio"
		} else {
			usuarios[i].NivelActividad = "bajo"
		}
	}

	return usuarios, nil
}

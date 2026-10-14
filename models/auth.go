package models

import "time"

type SolicitudPin struct {
	Correo string `json:"correo"`
}

type SolicitudVerificar struct {
	Correo string `json:"correo"`
	Pin    string `json:"pin"`
}

// el PIN se guarda en memoria mientras el MID este prendido
type CodigoPin struct {
	Correo   string
	Pin      string
	Creado   time.Time
	Vence    time.Time
	Intentos int
	Activo   bool
}

type SolicitudRegistro struct {
	Nombre    string `json:"nombre"`
	Email     string `json:"email"`
	Telefono  string `json:"telefono"`
	Ubicacion string `json:"ubicacion"`
}

type LoginAdmin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

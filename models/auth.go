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
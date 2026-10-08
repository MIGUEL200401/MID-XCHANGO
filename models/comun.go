package models

// cuerpo que usan moderacion (usuarios) y el panel (administradores)
type CambioEstado struct {
	Estado string `json:"estado"`
}
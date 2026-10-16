package models

// lo que manda el usuario para pedir que verifiquen su cuenta
type SolicitudDocumento struct {
	TipoDocumento string `json:"tipoDocumento"`
	Numero        string `json:"numero"`
	Archivo       string `json:"archivo"`
}

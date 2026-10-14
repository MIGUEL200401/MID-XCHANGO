package models

type AdminUsuario struct {
	Id           string `json:"id,omitempty"`
	Nombre       string `json:"nombre"`
	Email        string `json:"email"`
	Password     string `json:"password,omitempty"`
	Rol          string `json:"rol"`
	Estado       string `json:"estado"`
	Avatar       string `json:"avatar"`
	UltimoAcceso string `json:"ultimoAcceso"`
}

package models

type Trueque struct {
	Id  string  `json:"id,omitempty"`
	SolicitanteId  string  `json:"solicitanteId"`
	PropietarioId  string  `json:"propietarioId"`
	PublicacionId  string  `json:"publicacionId"`
	Ofrece  string  `json:"ofrece"`
	Busca   string  `json:"busca"`
	Estado  string  `json:"estado"`
	FechaSolicitud  string  `json:"fechaSolicitud"`
	FechaCierre  *string `json:"fechaCierre"`
	ConfirmadoSolicitante   bool   `json:"confirmadoSolicitante"`
	ConfirmadoPropietario   bool   `json:"confirmadoPropietario"`
	CalificacionSolicitante *int   `json:"calificacionSolicitante"`
	CalificacionPropietario *int   `json:"calificacionPropietario"`
}

type SolicitudTrueque struct {
	PublicacionId string `json:"publicacionId"`
	Ofrece        string `json:"ofrece"`
}

type Calificacion struct {
	Nota int `json:"nota"`
}


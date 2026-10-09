package models

type Reporte struct {
	Id string `json:"id,omitempty"`
	PublicacionId string `json:"publicacionId"`
	ReportanteId string `json:"reportanteId"`
	Motivo string `json:"motivo"`
	Fecha string `json:"fecha"`
	Estado string `json:"estado"`
	Gravedad string `json:"gravedad"`
}

type PublicacionEliminada struct {
	Id string `json:"id,omitempty"`
	PublicacionId string `json:"publicacionId"`
	UsuarioId string `json:"usuarioId"`
	AdminId string `json:"adminId"`
	Motivo string `json:"motivo"`
	Fecha string `json:"fecha"`
}

type Documento struct {
	Id string `json:"id,omitempty"`
	UsuarioId string `json:"usuarioId"`
	TipoDocumento string `json:"tipoDocumento"`
	Numero string `json:"numero"`
	Archivo string `json:"archivo"`
	FechaSubida string `json:"fechaSubida"`
	Estado string `json:"estado"`
	Observacion string `json:"observacion"`
	RevisadoPor *string `json:"revisadoPor"`
}
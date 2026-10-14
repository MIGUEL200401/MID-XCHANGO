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

type Reporte struct {
	Id            string `json:"id,omitempty"`
	PublicacionId string `json:"publicacionId"`
	ReportanteId  string `json:"reportanteId"`
	Motivo        string `json:"motivo"`
	Fecha         string `json:"fecha"`
	Estado        string `json:"estado"`
	Gravedad      string `json:"gravedad"`
}

type PublicacionEliminada struct {
	Id            string `json:"id,omitempty"`
	PublicacionId string `json:"publicacionId"`
	UsuarioId     string `json:"usuarioId"`
	AdminId       string `json:"adminId"`
	Motivo        string `json:"motivo"`
	Fecha         string `json:"fecha"`
}

type Documento struct {
	Id            string  `json:"id,omitempty"`
	UsuarioId     string  `json:"usuarioId"`
	TipoDocumento string  `json:"tipoDocumento"`
	Numero        string  `json:"numero"`
	Archivo       string  `json:"archivo"`
	FechaSubida   string  `json:"fechaSubida"`
	Estado        string  `json:"estado"`
	Observacion   string  `json:"observacion"`
	RevisadoPor   *string `json:"revisadoPor"`
}

type Alerta struct {
	Id        string `json:"id,omitempty"`
	Tipo      string `json:"tipo"`
	Titulo    string `json:"titulo"`
	Mensaje   string `json:"mensaje"`
	Fecha     string `json:"fecha"`
	Leida     bool   `json:"leida"`
	Prioridad string `json:"prioridad"`
}

type AccionAdmin struct {
	Id          string `json:"id,omitempty"`
	UsuarioId   string `json:"usuarioId"`
	AdminId     string `json:"adminId"`
	Accion      string `json:"accion"`
	Descripcion string `json:"descripcion"`
	Fecha       string `json:"fecha"`
}

// cuerpos de las peticiones del panel
type CambioEstado struct {
	Estado string `json:"estado"`
}

type SolicitudEliminar struct {
	Motivo string `json:"motivo"`
}

type RevisionDocumento struct {
	Estado      string `json:"estado"`
	Observacion string `json:"observacion"`
}

// la metrica del dashboard
type Punto struct {
	Label string `json:"label"`
	Valor int    `json:"valor"`
}

type Comparativa struct {
	Actual     int     `json:"actual"`
	Anterior   int     `json:"anterior"`
	Porcentaje float64 `json:"porcentaje"`
}

type Resumen struct {
	Usuarios                int `json:"usuarios"`
	UsuariosActivos         int `json:"usuariosActivos"`
	Publicaciones           int `json:"publicaciones"`
	Intercambios            int `json:"intercambios"`
	PublicacionesReportadas int `json:"publicacionesReportadas"`
	UsuariosSuspendidos     int `json:"usuariosSuspendidos"`
}

type Metricas struct {
	Resumen                   Resumen                `json:"resumen"`
	Comparativas              map[string]Comparativa `json:"comparativas"`
	ActividadMensual          []Punto                `json:"actividadMensual"`
	IntercambiosPorMes        []Punto                `json:"intercambiosPorMes"`
	PublicacionesPorCategoria []Punto                `json:"publicacionesPorCategoria"`
	ActividadPorUbicacion     []Punto                `json:"actividadPorUbicacion"`
	CrecimientoCategorias     []Punto                `json:"crecimientoCategorias"`
}

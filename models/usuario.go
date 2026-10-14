package models

type Usuario struct {
	Id                  string  `json:"id,omitempty"`
	Nombre              string  `json:"nombre"`
	Email               string  `json:"email"`
	Telefono            string  `json:"telefono"`
	Estado              string  `json:"estado"`
	FechaRegistro       string  `json:"fechaRegistro"`
	Ubicacion           string  `json:"ubicacion"`
	NivelActividad      string  `json:"nivelActividad"`
	Avatar              string  `json:"avatar"`
	Descripcion         string  `json:"descripcion"`
	Verificacion        string  `json:"verificacion"`
	Calificacion        float64 `json:"calificacion"`
	TotalCalificaciones int     `json:"totalCalificaciones"`
	Publicaciones       int     `json:"publicaciones"`
	Intercambios        int     `json:"intercambios"`
	Reportes            int     `json:"reportes"`
}
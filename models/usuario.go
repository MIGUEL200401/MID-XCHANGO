package models

type Usuario struct {
	Id                  string  `json:"id,omitempty"`
	Nombre              string  `json:"nombre"`
	Apellido            string  `json:"apellido"`
	FechaNacimiento     string  `json:"fechaNacimiento,omitempty"`
	Email               string  `json:"email"`
	Telefono            string  `json:"telefono"`
	Estado              string  `json:"estado"`
	FechaRegistro       string  `json:"fechaRegistro"`
	Ubicacion           string  `json:"ubicacion"`
	Direccion           string  `json:"direccion"`
	Barrio              string  `json:"barrio"`
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

// lo que se puede cambiar desde el perfil
type CambiosPerfil struct {
	Nombre      string `json:"nombre"`
	Apellido    string `json:"apellido"`
	Telefono    string `json:"telefono"`
	Ubicacion   string `json:"ubicacion"`
	Direccion   string `json:"direccion"`
	Barrio      string `json:"barrio"`
	Descripcion string `json:"descripcion"`
	Avatar      string `json:"avatar"`
}

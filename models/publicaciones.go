package models

type Publicacion struct {
	Id string `json:"id,omitempty"`
	UsuarioId string `json:"usuarioId"`
	CategoriaId string `json:"categoriaId"`
	Tipo string `json:"tipo"`
	Titulo string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	Ofreces string `json:"ofreces"`
	Buscas string `json:"buscas"`
	Imagenes []string `json:"imagenes"`
	Estado string `json:"estado"`
	Vistas int `json:"vistas"`
	FechaCreacion string `json:"fechaCreacion"`
	FechaModificacion string `json:"fechaModificacion"`
	Municipio string `json:"municipio,omitempty"`
	Barrio string `json:"barrio,omitempty"`
	CantidadDisponible string `json:"cantidadDisponible,omitempty"`
	Disponibilidad string `json:"disponibilidad,omitempty"`
}

type DatosPublicacion struct {
	Tipo string `json:"tipo"`
	CategoriaId string `json:"categoriaId"`
	Titulo string `json:"titulo"`
	Descripcion string `json:"descrpcion"`
	Ofreces string `json:"ofreces"`
	Buscas string `json:"buscas"`
	Municipio string `json:"municipio"`
	Barrio string `json:"barrio"`
	CantidadDisponible string `json:"cantidadDisponible"`
	Disponibilidad string `json:"disponibilidad"`
	Imagenes []string `json:"imagenes"`
	Estado string `json:"estado"`
}

type Categoria struct {
	Id string `json:"id,omitempty"`
	Nombre string `json:"nombre"`
	Icono string `json:"icono"`
}

type AvisoUbicacion struct {
	Id string `json:"id,omitempty"`
	Municipio string `json:"municipio"`
	Titulo string `json:"titulo"`
	Mensaje string `json:"mensaje"`
	Icono string `json:"icono"`
}

type Favoritos struct {
	Id string `json:"id,omitempty"`
	UsuarioId string `json:"usuarioId"`
	PublicacionId string `json:"publicacionId"`
}
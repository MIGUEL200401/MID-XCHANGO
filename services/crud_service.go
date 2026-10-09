package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

var ErrNoEncontrado = errors.New("no encontrado")

// habla con el CRUD, que es el que guarda en la base de datos
type CrudService struct {
	URL     string
	Cliente *http.Client
}

func NewCrudService() *CrudService {
	return &CrudService{
		URL:     beego.AppConfig.DefaultString("url_crud", "http://localhost:3000"),
		Cliente: &http.Client{Timeout: 10 * time.Second},
	}
}

// trae toda una coleccion
func (s *CrudService) Listar(coleccion string, destino interface{}) error {
	return s.pedir("GET", coleccion, nil, destino)
}

func (s *CrudService) Obtener(coleccion string, id string, destino interface{}) error {
	return s.pedir("GET", coleccion+"/"+id, nil, destino)
}

func (s *CrudService) Crear(coleccion string, dato interface{}, destino interface{}) error {
	return s.pedir("POST", coleccion, dato, destino)
}

func (s *CrudService) Actualizar(coleccion string, id string, dato interface{}, destino interface{}) error {
	return s.pedir("PUT", coleccion+"/"+id, dato, destino)
}

func (s *CrudService) Borrar(coleccion string, id string) error {
	return s.pedir("DELETE", coleccion+"/"+id, nil, nil)
}

func (s *CrudService) pedir(metodo string, ruta string, dato interface{}, destino interface{}) error {
	var cuerpo bytes.Buffer

	if dato != nil {
		err := json.NewEncoder(&cuerpo).Encode(dato)
		if err != nil {
			return fmt.Errorf("error armando el JSON: %w", err)
		}
	}

	peticion, err := http.NewRequest(metodo, s.URL+"/"+ruta, &cuerpo)
	if err != nil {
		return fmt.Errorf("error creando la peticion: %w", err)
	}
	peticion.Header.Set("Content-Type", "application/json")

	respuesta, err := s.Cliente.Do(peticion)
	if err != nil {
		return fmt.Errorf("no se pudo conectar con el CRUD: %w", err)
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode == http.StatusNotFound {
		return ErrNoEncontrado
	}

	if respuesta.StatusCode < 200 || respuesta.StatusCode > 299 {
		return fmt.Errorf("el CRUD respondio con codigo %d", respuesta.StatusCode)
	}

	if destino == nil || respuesta.StatusCode == http.StatusNoContent {
		return nil
	}

	err = json.NewDecoder(respuesta.Body).Decode(destino)
	if err != nil {
		return fmt.Errorf("error leyendo la respuesta del CRUD: %w", err)
	}

	return nil
}

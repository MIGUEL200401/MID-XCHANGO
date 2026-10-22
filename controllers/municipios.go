package controllers

import "strings"

// los municipios de Casanare, escritos como estan en la base de datos
var municipiosCasanare = []string{
	"Yopal", "Aguazul", "Chámeza", "Hato Corozal", "La Salina", "Maní",
	"Monterrey", "Nunchía", "Orocué", "Paz de Ariporo", "Pore", "Recetor",
	"Sabanalarga", "Sácama", "San Luis de Palenque", "Támara", "Tauramena",
	"Trinidad", "Villanueva",
}

// quita tildes y pasa a minusculas para poder comparar "mani" con "Maní"
func sinTildes(texto string) string {
	cambios := strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u")
	return cambios.Replace(strings.ToLower(strings.TrimSpace(texto)))
}

// devuelve el municipio como esta en la base, o "" si no existe
func buscarMunicipio(texto string) string {
	buscado := sinTildes(texto)
	for _, municipio := range municipiosCasanare {
		if sinTildes(municipio) == buscado {
			return municipio
		}
	}
	return ""
}

const mensajeMunicipio = "Ese municipio no existe en Casanare. Ejemplos: Yopal, Aguazul, Maní, Tauramena."

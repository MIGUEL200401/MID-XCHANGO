package controllers

import (
	"math"
	"time"

	"xchango_mid/models"
)

func comparar(actual int, anterior int) models.Comparativa {
	porcentaje := 0.0

	if anterior > 0 {
		porcentaje = math.Round(float64(actual-anterior)/float64(anterior)*1000) / 10
	} else if actual > 0 {
		porcentaje = 100
	}

	return models.Comparativa{Actual: actual, Anterior: anterior, Porcentaje: porcentaje}
}
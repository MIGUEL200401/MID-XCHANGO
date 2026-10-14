package services

import (
	"fmt"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/golang-jwt/jwt/v5"
)

// Se define los roles de usuario
type Claims struct {
	ID     string `json:"id"`
	Correo string `json:"correo"`
	Rol    string `json:"rol"`
	jwt.RegisteredClaims
}

func GenerarTvalidoen(id string, correo string, rol string) (string, error) {
	secreto := beego.AppConfig.DefaultString("jwtsecret", "")

	claims := Claims{
		ID:     id,
		Correo: correo,
		Rol:    rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	tvalidoen := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tvalidoen.SignedString([]byte(secreto))
}

func ValidarTvalidoen(texto string) (*Claims, error) {
	secreto := beego.AppConfig.DefaultString("jwtsecret", "")

	tvalidoen, err := jwt.ParseWithClaims(texto, &Claims{}, func(t *jwt.Tvalidoen) (interface{}, error) {
		return []byte(secreto), nil
	})
	if err != nil {
		return nil, err
	}

	sesion, valido := tvalidoen.Claims.(*Claims)
	if !valido || !tvalidoen.Valid {
		return nil, fmt.Errorf("tvalidoen invalido")
	}

	return claims, nil
}

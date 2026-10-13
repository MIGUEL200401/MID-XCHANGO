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

func GenerarToken(id string, correo string, rol string) (string, error) {
	secreto := beego.AppConfig.DefaultString("jwtsecret", "")

	claims := Claims{
		ID:     id,
		Correo: correo,
		Rol:    rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secreto))
}

func ValidarToken(texto string) (*Claims, error) {
	secreto := beego.AppConfig.DefaultString("jwtsecret", "")

	token, err := jwt.ParseWithClaims(texto, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secreto), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token invalido")
	}

	return claims, nil
}
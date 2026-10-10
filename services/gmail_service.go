package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type ServicioGmail struct {
	Servicio *gmail.Service
}

// se crea el servicio de Gmail, si no tiene token pide autorizacion y lo guarda en token/token.json
func NuevoServicioGmail() (*ServicioGmail, error) {
	contexto := context.Background()

	rutaCredenciales := "credentials/credentials.json"

	credenciales, err := os.ReadFile(rutaCredenciales)
	if err != nil {
		return nil, fmt.Errorf("no fue posible leer credentials.json: %w", err)
	}

	configuracion, err := google.ConfigFromJSON(
		credenciales,
		gmail.GmailSendScope,
	)
	if err != nil {
		return nil, fmt.Errorf("no fue posible configurar OAuth 2.0: %w", err)
	}

	cliente, err := obtenerClienteGmail(configuracion)
	if err != nil {
		return nil, err
	}

	servicio, err := gmail.NewService(
		contexto,
		option.WithHTTPClient(cliente),
	)
	if err != nil {
		return nil, fmt.Errorf("no fue posible crear el servicio de Gmail: %w", err)
	}

	return &ServicioGmail{
		Servicio: servicio,
	}, nil
}

// si ya hay un token guardado lo usa, si no pide autorizacion
func obtenerClienteGmail(configuracion *oauth2.Config) (*http.Client, error) {
	rutaToken := "token/token.json"

	token, err := leerTokenDeArchivo(rutaToken)
	if err != nil {
		token, err = pedirTokenPorNavegador(configuracion)
		if err != nil {
			return nil, err
		}

		err = guardarToken(rutaToken, token)
		if err != nil {
			return nil, err
		}
	}

	return configuracion.Client(context.Background(), token), nil
}

// muestra el link de Google y espera que se pegue el codigo
func pedirTokenPorNavegador(configuracion *oauth2.Config) (*oauth2.Token, error) {
	urlAutorizacion := configuracion.AuthCodeURL(
		"state-token",
		oauth2.AccessTypeOffline,
	)
	fmt.Println("===================")
	fmt.Println("Autorizacion de Gmail")
	fmt.Println("===================")
	fmt.Println("Abra la siguiente URL en su navegador")
	fmt.Println()
	fmt.Println(urlAutorizacion)
	fmt.Println()
	fmt.Println("===================")

	fmt.Print("Ingrese el codigo de autorizacion: ")

	var codigo string
	_, err := fmt.Scan(&codigo)
	if err != nil {
		return nil, fmt.Errorf("no fue posible leer el codigo: %w", err)
	}

	token, err := configuracion.Exchange(
		context.Background(),
		codigo,
	)
	if err != nil {
		return nil, fmt.Errorf("no fue posible obtener el token: %w", err)
	}

	return token, nil
}

// lee el token que se guardo la vez anterior
func leerTokenDeArchivo(ruta string) (*oauth2.Token, error) {
	archivo, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer archivo.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(archivo).Decode(token)
	if err != nil {
		return nil, fmt.Errorf("no fue posible leer token.json: %w", err)
	}

	return token, nil
}

func guardarToken(ruta string, token *oauth2.Token) error {
	err := os.MkdirAll("token", 0700)
	if err != nil {
		return fmt.Errorf("no fue posible crear la carpeta token: %w", err)
	}

	archivo, err := os.OpenFile(
		ruta,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		0600,
	)
	if err != nil {
		return fmt.Errorf("no fue posible crear el archivo del token: %w", err)
	}
	defer archivo.Close()

	return json.NewEncoder(archivo).Encode(token)
}

// divide el texto en líneas de máximo 76 caracteres
func partirLineas(texto string) string {
	resultado := ""
	for len(texto) > 76 {
		resultado += texto[:76] + "\r\n"
		texto = texto[76:]
	}
	return resultado + texto
}

// se arma el correo con el html y el logo 
func armarCorreo(destinatario string, asunto string, cuerpo string, logo []byte) string {
	limite := "limite-xchango"

	partes := []string{
		"To: " + destinatario,
		"Subject: " + mime.QEncoding.Encode("UTF-8", asunto),
		"MIME-Version: 1.0",
		"Content-Type: multipart/related; boundary=" + limite,
		"",
		"--" + limite,
		"Content-Type: text/html; charset=UTF-8",
		"Content-Transfer-Encoding: base64",
		"",
		partirLineas(base64.StdEncoding.EncodeToString([]byte(cuerpo))),
	}

	if len(logo) > 0 {
		partes = append(partes,
			"--"+limite,
			"Content-Type: image/png",
			"Content-Transfer-Encoding: base64",
			"Content-ID: <logo-xchango>",
			"Content-Disposition: inline; filename=\"logo.png\"",
			"",
			partirLineas(base64.StdEncoding.EncodeToString(logo)),
		)
	}

	partes = append(partes, "--"+limite+"--", "")

	return strings.Join(partes, "\r\n")
}

func (g *ServicioGmail) EnviarCorreo(
	destinatario string,
	asunto string,
	cuerpo string,
	logo []byte,
) (string, error) {
	if strings.TrimSpace(destinatario) == "" {
		return "", fmt.Errorf("el destinatario es obligatorio")
	}
	if strings.TrimSpace(asunto) == "" {
		return "", fmt.Errorf("el asunto es obligatorio")
	}

	correo := armarCorreo(destinatario, asunto, cuerpo, logo)
	mensajeCodificado := base64.URLEncoding.EncodeToString([]byte(correo))

	mensaje := &gmail.Message{
		Raw: mensajeCodificado,
	}

	respuesta, err := g.Servicio.Users.Messages.Send("me", mensaje).Do()
	if err != nil {
		return "", fmt.Errorf("no fue posible enviar el correo con Gmail: %w", err)
	}

	if respuesta == nil || respuesta.Id == "" {
		return "", fmt.Errorf("no fue posible enviar el correo con Gmail")
	}

	return respuesta.Id, nil
}

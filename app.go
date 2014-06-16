package main

import (
	"github.com/go-martini/martini"
	"github.com/handshakejs/handshakejserrors"
	"github.com/handshakejs/handshakejslogic"
	"github.com/handshakejs/handshakejstransport"
	"github.com/hoisie/mustache"
	"github.com/joho/godotenv"
	"github.com/martini-contrib/render"
	"net/http"
	"os"
)

const (
	LOGIC_ERROR_CODE_UNKNOWN = "unknown"
)

var (
	DB_ENCRYPTION_SALT   string
	REDIS_URL            string
	FROM                 string
	SMTP_ADDRESS         string
	SMTP_PORT            string
	SMTP_USERNAME        string
	SMTP_PASSWORD        string
	BACKUP_SMTP_ADDRESS  string
	BACKUP_SMTP_PORT     string
	BACKUP_SMTP_USERNAME string
	BACKUP_SMTP_PASSWORD string
	SUBJECT_TEMPLATE     string
	TEXT_TEMPLATE        string
	HTML_TEMPLATE        string
)

func CrossDomain() martini.Handler {
	return func(res http.ResponseWriter) {
		res.Header().Add("Access-Control-Allow-Origin", "*")
	}
}

func main() {
	loadEnvs()

	logic_options := handshakejslogic.Options{DbEncryptionSalt: DB_ENCRYPTION_SALT}
	transport_options := handshakejstransport.Options{SmtpAddress: SMTP_ADDRESS, SmtpPort: SMTP_PORT, SmtpUsername: SMTP_USERNAME, SmtpPassword: SMTP_PASSWORD}

	handshakejslogic.Setup(REDIS_URL, logic_options)
	handshakejstransport.Setup(transport_options)

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(CrossDomain())

	m.Any("/api/v1/apps/create.json", AppsCreate)
	m.Any("/api/v1/login/request.json", IdentitiesCreate)
	m.Any("/api/v1/login/confirm.json", IdentitiesConfirm)

	m.Run()
}

func ErrorPayload(logic_error *handshakejserrors.LogicError) map[string]interface{} {
	error_object := map[string]interface{}{"code": logic_error.Code, "field": logic_error.Field, "message": logic_error.Message}
	errors := []interface{}{}
	errors = append(errors, error_object)
	payload := map[string]interface{}{"errors": errors}

	return payload
}

func AppsPayload(app map[string]interface{}) map[string]interface{} {
	apps := []interface{}{}
	apps = append(apps, app)
	payload := map[string]interface{}{"apps": apps}

	return payload
}

func IdentitiesCreatePayload(identity map[string]interface{}) map[string]interface{} {
	email := identity["email"].(string)
	app_name := identity["app_name"].(string)
	authcode_expired_at := identity["authcode_expired_at"].(string)

	identities := []interface{}{}
	output_identity := map[string]interface{}{"email": email, "app_name": app_name, "authcode_expired_at": authcode_expired_at}
	identities = append(identities, output_identity)
	payload := map[string]interface{}{"identities": identities}

	return payload
}

func IdentitiesConfirmPayload(identity map[string]interface{}) map[string]interface{} {
	email := identity["email"].(string)
	app_name := identity["app_name"].(string)
	hash := identity["hash"].(string)

	identities := []interface{}{}
	output_identity := map[string]interface{}{"email": email, "app_name": app_name, "hash": hash}
	identities = append(identities, output_identity)
	payload := map[string]interface{}{"identities": identities}

	return payload
}

func AppsCreate(req *http.Request, r render.Render) {
	email := req.URL.Query().Get("email")
	app_name := req.URL.Query().Get("app_name")
	salt := req.URL.Query().Get("salt")

	app := map[string]interface{}{"email": email, "app_name": app_name, "salt": salt}
	result, logic_error := handshakejslogic.AppsCreate(app)
	if logic_error != nil {
		payload := ErrorPayload(logic_error)
		statuscode := determineStatusCodeFromLogicError(logic_error)
		r.JSON(statuscode, payload)
	} else {
		payload := AppsPayload(result)
		r.JSON(200, payload)
	}
}

func IdentitiesCreate(req *http.Request, r render.Render) {
	email := req.URL.Query().Get("email")
	app_name := req.URL.Query().Get("app_name")

	identity := map[string]interface{}{"email": email, "app_name": app_name}
	result, logic_error := handshakejslogic.IdentitiesCreate(identity)
	if logic_error != nil {
		payload := ErrorPayload(logic_error)
		statuscode := determineStatusCodeFromLogicError(logic_error)
		r.JSON(statuscode, payload)
	} else {
		logic_error := deliverAuthcodeEmail(result)
		if logic_error != nil {
			payload := ErrorPayload(logic_error)
			statuscode := determineStatusCodeFromLogicError(logic_error)
			r.JSON(statuscode, payload)
		} else {
			payload := IdentitiesCreatePayload(result)
			r.JSON(200, payload)
		}
	}
}

func IdentitiesConfirm(req *http.Request, r render.Render) {
	email := req.URL.Query().Get("email")
	app_name := req.URL.Query().Get("app_name")
	authcode := req.URL.Query().Get("authcode")

	identity := map[string]interface{}{"email": email, "app_name": app_name, "authcode": authcode}
	result, logic_error := handshakejslogic.IdentitiesConfirm(identity)
	if logic_error != nil {
		payload := ErrorPayload(logic_error)
		statuscode := determineStatusCodeFromLogicError(logic_error)
		r.JSON(statuscode, payload)
	} else {
		payload := IdentitiesConfirmPayload(result)
		r.JSON(200, payload)
	}
}

func determineStatusCodeFromLogicError(logic_error *handshakejserrors.LogicError) int {
	code := 400
	if logic_error.Code == LOGIC_ERROR_CODE_UNKNOWN {
		code = 500
	}

	return code
}

func deliverAuthcodeEmail(identity map[string]interface{}) *handshakejserrors.LogicError {
	email := identity["email"].(string)
	subject := renderTemplate(SUBJECT_TEMPLATE, identity)
	text := renderTemplate(TEXT_TEMPLATE, identity)
	html := renderTemplate(HTML_TEMPLATE, identity)

	logic_error := handshakejstransport.ViaEmail(email, FROM, subject, text, html)
	if logic_error != nil {
		if BACKUP_SMTP_ADDRESS != "" {
			backup_transport_options := handshakejstransport.Options{SmtpAddress: BACKUP_SMTP_ADDRESS, SmtpPort: BACKUP_SMTP_PORT, SmtpUsername: BACKUP_SMTP_USERNAME, SmtpPassword: BACKUP_SMTP_PASSWORD}
			handshakejstransport.Setup(backup_transport_options)

			logic_error := handshakejstransport.ViaEmail(email, FROM, subject, text, html)
			if logic_error != nil {
				return logic_error
			}
		} else {
			return logic_error
		}
	}

	return nil
}

func renderTemplate(template_string string, identity map[string]interface{}) string {
	data := mustache.Render(template_string, identity)

	return data
}

func loadEnvs() {
	godotenv.Load()

	DB_ENCRYPTION_SALT = os.Getenv("DB_ENCRYPTION_SALT")
	REDIS_URL = os.Getenv("REDIS_URL")
	FROM = os.Getenv("FROM")
	SMTP_ADDRESS = os.Getenv("SMTP_ADDRESS")
	SMTP_PORT = os.Getenv("SMTP_PORT")
	SMTP_USERNAME = os.Getenv("SMTP_USERNAME")
	SMTP_PASSWORD = os.Getenv("SMTP_PASSWORD")
	BACKUP_SMTP_ADDRESS = os.Getenv("BACKUP_SMTP_ADDRESS")
	BACKUP_SMTP_PORT = os.Getenv("BACKUP_SMTP_PORT")
	BACKUP_SMTP_USERNAME = os.Getenv("BACKUP_SMTP_USERNAME")
	BACKUP_SMTP_PASSWORD = os.Getenv("BACKUP_SMTP_PASSWORD")
	SUBJECT_TEMPLATE = os.Getenv("SUBJECT_TEMPLATE")
	TEXT_TEMPLATE = os.Getenv("TEXT_TEMPLATE")
	HTML_TEMPLATE = os.Getenv("HTML_TEMPLATE")
}

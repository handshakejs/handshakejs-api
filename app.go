package main

import (
	"github.com/go-martini/martini"
	"github.com/handshakejs/handshakejslogic"
	"github.com/martini-contrib/render"
	"net/http"
)

const (
	LOGIC_ERROR_CODE_UNKNOWN = "unkown"
)

func main() {
	handshakejslogic.Setup("redis://127.0.0.1:6379")

	m := martini.Classic()
	m.Use(martini.Logger())
	m.Use(render.Renderer())

	m.Any("/api/v1/apps/create.json", AppsCreate)
	m.Any("/api/v1/login/request.json", IdentitiesCreate)

	m.Run()
}

func ErrorPayload(logic_error *handshakejslogic.LogicError) map[string]interface{} {
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

func IdentitiesPayload(identity map[string]interface{}) map[string]interface{} {
	email := identity["email"].(string)
	app_name := identity["app_name"].(string)
	authcode_expired_at := identity["authcode_expired_at"].(string)

	identities := []interface{}{}
	output_identity := map[string]interface{}{"email": email, "app_name": app_name, "authcode_expired_at": authcode_expired_at}
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
		payload := IdentitiesPayload(result)
		r.JSON(200, payload)
	}
}

func determineStatusCodeFromLogicError(logic_error *handshakejslogic.LogicError) int {
	code := 400
	if logic_error.Code == LOGIC_ERROR_CODE_UNKNOWN {
		code = 500
	}

	return code
}

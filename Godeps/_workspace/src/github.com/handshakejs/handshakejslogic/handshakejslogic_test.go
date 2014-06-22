package handshakejslogic_test

import (
	"../handshakejslogic"
	"github.com/garyburd/redigo/redis"
	"github.com/stvp/tempredis"
	"log"
	"testing"
)

const (
	APP_NAME           = "app0"
	EMAIL              = "app0@mailinator.com"
	IDENTITY_EMAIL     = "identity0@mailinator.com"
	AUTHCODE           = "5678"
	SALT               = "1234"
	REDIS_URL          = "redis://127.0.0.1:11001"
	DB_ENCRYPTION_SALT = "somesecretsaltthatis32characters"
)

func defaultOptions() handshakejslogic.Options {
	value := handshakejslogic.Options{DbEncryptionSalt: DB_ENCRYPTION_SALT}
	return value
}

func tempredisConfig() tempredis.Config {
	config := tempredis.Config{
		"port":      "11001",
		"databases": "1",
	}
	return config
}

func TestAppsCreate(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		app := map[string]interface{}{"email": EMAIL, "app_name": APP_NAME}

		handshakejslogic.Setup(REDIS_URL, defaultOptions())
		result, logic_error := handshakejslogic.AppsCreate(app)
		if logic_error != nil {
			t.Errorf("Error", logic_error)
		}
		if result["email"] != EMAIL {
			t.Errorf("Incorrect email " + result["email"].(string))
		}
		if result["app_name"] != APP_NAME {
			t.Errorf("Incorrect app_name " + result["app_name"].(string))
		}
		if result["salt"] == nil {
			t.Errorf("Salt is nil and should not be.")
		}
	})
}

func TestAppsCreateCustomSalt(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		app := map[string]interface{}{"email": EMAIL, "app_name": APP_NAME, "salt": SALT}

		options := defaultOptions()
		handshakejslogic.Setup(REDIS_URL, options)
		result, logic_error := handshakejslogic.AppsCreate(app)
		if logic_error != nil {
			t.Errorf("Error", logic_error)
		}

		if result["salt"] != SALT {
			t.Errorf("Salt did not equal " + SALT)
		}
	})
}

func TestAppsCreateCustomBlankSalt(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		app := map[string]interface{}{"email": EMAIL, "app_name": APP_NAME, "salt": ""}

		handshakejslogic.Setup(REDIS_URL, defaultOptions())
		result, logic_error := handshakejslogic.AppsCreate(app)
		if logic_error != nil {
			t.Errorf("Error", logic_error)
		}

		if result["salt"] == nil || result["salt"].(string) == "" {
			t.Errorf("It should generate a salt if blank.")
		}
	})
}

func TestAppsCreateBlankAppName(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		app := map[string]interface{}{"email": EMAIL, "app_name": ""}

		handshakejslogic.Setup(REDIS_URL, defaultOptions())
		_, logic_error := handshakejslogic.AppsCreate(app)
		if logic_error.Code != "required" {
			t.Errorf("Error", err)
		}
	})
}

func TestAppsCreateNilAppName(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		app := map[string]interface{}{"email": EMAIL}

		handshakejslogic.Setup(REDIS_URL, defaultOptions())
		_, logic_error := handshakejslogic.AppsCreate(app)
		if logic_error.Code != "required" {
			t.Errorf("Error", err)
		}
	})
}

func TestAppsCreateSpacedAppName(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		app := map[string]interface{}{"email": EMAIL, "app_name": " "}

		handshakejslogic.Setup(REDIS_URL, defaultOptions())
		_, logic_error := handshakejslogic.AppsCreate(app)
		if logic_error.Code != "required" {
			t.Errorf("Error", err)
		}
	})
}

func TestAppsCreateAppNameWithSpaces(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		app := map[string]interface{}{"email": EMAIL, "app_name": "combine these"}

		handshakejslogic.Setup(REDIS_URL, defaultOptions())
		result, _ := handshakejslogic.AppsCreate(app)
		if result["app_name"] != "combinethese" {
			t.Errorf("Incorrect combining of app_name " + result["app_name"].(string))
		}
	})
}

func TestIdentitiesCreate(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)

		identity := map[string]interface{}{"app_name": APP_NAME, "email": IDENTITY_EMAIL}
		result, logic_error := handshakejslogic.IdentitiesCreate(identity)
		if logic_error != nil {
			t.Errorf("Error", logic_error)
		}
		if result["authcode_expired_at"] == nil {
			t.Errorf("Error", result)
		}
		if len(result["authcode"].(string)) < 4 {
			t.Errorf("Error", result)
		}
	})
}

func TestIdentitiesCreateBlankAppName(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)

		identity := map[string]interface{}{"app_name": "", "email": IDENTITY_EMAIL}
		_, logic_error := handshakejslogic.IdentitiesCreate(identity)
		if logic_error.Code != "required" {
			t.Errorf("Error", logic_error)
		}
		if logic_error.Field != "app_name" {
			t.Errorf("Error", logic_error)
		}
	})
}

func TestIdentitiesCreateNilAppName(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)

		identity := map[string]interface{}{"email": IDENTITY_EMAIL}
		_, logic_error := handshakejslogic.IdentitiesCreate(identity)
		if logic_error.Code != "required" {
			t.Errorf("Error", logic_error)
		}
		if logic_error.Field != "app_name" {
			t.Errorf("Error", logic_error)
		}
	})
}

func TestIdentitiesCreateNonExistingAppName(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)

		identity := map[string]interface{}{"app_name": "doesnotexist", "email": IDENTITY_EMAIL}
		_, logic_error := handshakejslogic.IdentitiesCreate(identity)
		if logic_error.Code != "not_found" {
			t.Errorf("Error", logic_error)
		}
		if logic_error.Field != "app_name" {
			t.Errorf("Error", logic_error)
		}
	})
}
func TestIdentitiesCreateBlankEmail(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)

		identity := map[string]interface{}{"app_name": APP_NAME, "email": ""}
		_, logic_error := handshakejslogic.IdentitiesCreate(identity)
		if logic_error.Code != "required" {
			t.Errorf("Error", logic_error)
		}
		if logic_error.Field != "email" {
			t.Errorf("Error", logic_error)
		}
	})
}

func TestIdentitiesCreateNilEmail(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)

		identity := map[string]interface{}{"app_name": APP_NAME}
		_, logic_error := handshakejslogic.IdentitiesCreate(identity)
		if logic_error.Code != "required" {
			t.Errorf("Error", logic_error)
		}
		if logic_error.Field != "email" {
			t.Errorf("Error", logic_error)
		}
	})
}

func TestIdentitiesConfirm(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)
		authcode := setupIdentity(t)

		identity_check := map[string]interface{}{"app_name": APP_NAME, "email": IDENTITY_EMAIL, "authcode": authcode}
		result, logic_error := handshakejslogic.IdentitiesConfirm(identity_check)
		if logic_error != nil {
			t.Errorf("Error", logic_error)
		}
		if result["hash"].(string) == "" {
			t.Errorf("Error", "missing hash in result")
		}
		if result["hash"].(string) != "2402d6b6008c2cd1a3c73db00d8bac8a" {
			t.Errorf("Error", result["hash"].(string)+" is the incorrect hash")
		}
	})
}

func TestIdentitiesConfirmIncorrectAuthcode(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)
		setupIdentity(t)

		identity_check := map[string]interface{}{"app_name": APP_NAME, "email": IDENTITY_EMAIL, "authcode": "1234"}
		_, logic_error := handshakejslogic.IdentitiesConfirm(identity_check)
		if logic_error.Code != "incorrect" {
			t.Errorf("Error", logic_error)
		}
		if logic_error.Field != "authcode" {
			t.Errorf("Error", logic_error)
		}
	})
}

func TestIdentitiesConfirmExpiredAuthcode(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupAppWithShortAuthcodeLife(t)
		authcode := setupIdentity(t)

		identity_check := map[string]interface{}{"app_name": APP_NAME, "email": IDENTITY_EMAIL, "authcode": authcode}
		_, logic_error := handshakejslogic.IdentitiesConfirm(identity_check)
		if logic_error.Code != "expired" {
			t.Errorf("Error", logic_error)
		}
		if logic_error.Field != "authcode" {
			t.Errorf("Error", logic_error)
		}
	})
}
func TestIdentitiesConfirmBlankAppName(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)
		setupIdentity(t)

		identity_check := map[string]interface{}{"app_name": "", "email": IDENTITY_EMAIL, "authcode": AUTHCODE}
		_, logic_error := handshakejslogic.IdentitiesConfirm(identity_check)
		if logic_error.Code != "required" {
			t.Errorf("Error", logic_error)
		}
		if logic_error.Field != "app_name" {
			t.Errorf("Error", logic_error)
		}
	})
}

func TestIdentitiesConfirmNilAppName(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)
		setupIdentity(t)

		identity_check := map[string]interface{}{"email": IDENTITY_EMAIL, "authcode": AUTHCODE}
		_, logic_error := handshakejslogic.IdentitiesConfirm(identity_check)
		if logic_error.Code != "required" {
			t.Errorf("Error", logic_error)
		}
		if logic_error.Field != "app_name" {
			t.Errorf("Error", logic_error)
		}
	})
}

func TestIdentitiesConfirmNonExistingAppName(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)
		setupIdentity(t)

		identity_check := map[string]interface{}{"app_name": "doesnotexist", "email": IDENTITY_EMAIL, "authcode": AUTHCODE}
		_, logic_error := handshakejslogic.IdentitiesConfirm(identity_check)
		if logic_error.Code != "not_found" {
			t.Errorf("Error", logic_error)
		}
		if logic_error.Field != "app_name" {
			t.Errorf("Error", logic_error)
		}
	})
}

func TestIdentitiesConfirmBlankEmail(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)
		setupIdentity(t)

		identity_check := map[string]interface{}{"app_name": APP_NAME, "email": "", "authcode": AUTHCODE}
		_, logic_error := handshakejslogic.IdentitiesConfirm(identity_check)
		if logic_error.Code != "required" {
			t.Errorf("Error", logic_error)
		}
		if logic_error.Field != "email" {
			t.Errorf("Error", logic_error)
		}
	})
}

func TestIdentitiesConfirmNilEmail(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)
		setupIdentity(t)

		identity_check := map[string]interface{}{"app_name": APP_NAME, "authcode": AUTHCODE}
		_, logic_error := handshakejslogic.IdentitiesConfirm(identity_check)
		if logic_error.Code != "required" {
			t.Errorf("Error", logic_error)
		}
		if logic_error.Field != "email" {
			t.Errorf("Error", logic_error)
		}
	})
}

func TestIdentitiesConfirmBlankAuthcode(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)
		setupIdentity(t)

		identity_check := map[string]interface{}{"app_name": APP_NAME, "email": EMAIL, "authcode": ""}
		_, logic_error := handshakejslogic.IdentitiesConfirm(identity_check)
		if logic_error.Code != "required" {
			t.Errorf("Error", logic_error)
		}
		if logic_error.Field != "authcode" {
			t.Errorf("Error", logic_error)
		}
	})
}

func TestIdentitiesConfirmNilAuthcode(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)
		setupIdentity(t)

		identity_check := map[string]interface{}{"app_name": APP_NAME, "email": EMAIL}
		_, logic_error := handshakejslogic.IdentitiesConfirm(identity_check)
		if logic_error.Code != "required" {
			t.Errorf("Error", logic_error)
		}
		if logic_error.Field != "authcode" {
			t.Errorf("Error", logic_error)
		}
	})
}

func TestIdentitiesConfirmNonExistingEmail(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupApp(t)
		setupIdentity(t)

		identity_check := map[string]interface{}{"app_name": APP_NAME, "email": "doenot@existe.com", "authcode": AUTHCODE}
		_, logic_error := handshakejslogic.IdentitiesConfirm(identity_check)
		if logic_error.Code != "not_found" {
			t.Errorf("Error", logic_error)
		}
		if logic_error.Field != "email" {
			t.Errorf("Error", logic_error)
		}
	})
}

func setupApp(t *testing.T) {
	app := map[string]interface{}{"email": EMAIL, "app_name": APP_NAME, "salt": SALT}

	handshakejslogic.Setup(REDIS_URL, defaultOptions())
	_, logic_error := handshakejslogic.AppsCreate(app)
	if logic_error != nil {
		t.Errorf("Error", logic_error)
	}
}

func setupAppWithShortAuthcodeLife(t *testing.T) {
	app := map[string]interface{}{"email": EMAIL, "app_name": APP_NAME}

	// set it negative for test purposes
	options := handshakejslogic.Options{DbEncryptionSalt: DB_ENCRYPTION_SALT, AuthcodeLifeInMs: -5, AuthcodeLength: 5}
	handshakejslogic.Setup(REDIS_URL, options)
	_, logic_error := handshakejslogic.AppsCreate(app)
	if logic_error != nil {
		t.Errorf("Error", logic_error)
	}
}

func setupIdentity(t *testing.T) string {
	identity := map[string]interface{}{"app_name": APP_NAME, "email": IDENTITY_EMAIL}
	_, logic_error := handshakejslogic.IdentitiesCreate(identity)
	if logic_error != nil {
		t.Errorf("Error", logic_error)
	}

	app_name_key := "apps/" + identity["app_name"].(string)
	key := app_name_key + "/identities/" + identity["email"].(string)
	authcode, err := redis.String(handshakejslogic.Conn().Do("HGET", key, "authcode"))
	if err != nil {
		t.Errorf("Error", err)
	}

	return authcode
}

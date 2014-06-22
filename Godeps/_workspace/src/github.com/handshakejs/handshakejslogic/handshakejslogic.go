package handshakejslogic

import (
	"bytes"
	"code.google.com/p/go.crypto/pbkdf2"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/dchest/uniuri"
	"github.com/garyburd/redigo/redis"
	"github.com/handshakejs/handshakejscrypter"
	"github.com/handshakejs/handshakejserrors"
	"github.com/scottmotte/redisurlparser"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	BASE_10                           = 10
	AUTHCODE_LIFE_IN_MS_DEFAULT       = 120000
	AUTHCODE_LENGTH_DEFAULT           = 4
	KEY_EXPIRATION_IN_SECONDS_DEFAULT = 86400 // 24 hours in seconds
	PBKDF2_HASH_ITERATIONS_DEFAULT    = 1000
	PBKDF2_HASH_BITES_DEFAULT         = 16
	DB_ENCRYPTION_ITERATIONS          = 1000
	DB_ENCRYPTION_BITES               = 16
)

var (
	DB_ENCRYPTION_SALT        string
	AUTHCODE_LIFE_IN_MS       int64
	AUTHCODE_LENGTH           int
	KEY_EXPIRATION_IN_SECONDS int
	PBKDF2_HASH_ITERATIONS    int
	PBKDF2_HASH_BITES         int
	redisurl                  redisurlparser.RedisURL
	pool                      *redis.Pool
)

type Options struct {
	DbEncryptionSalt       string
	AuthcodeLifeInMs       int64
	AuthcodeLength         int
	KeyExpirationInSeconds int
	Pbkdf2HashIterations   int
	Pbkdf2HashBites        int
}

func Setup(redis_url_string string, options Options) {
	if options.DbEncryptionSalt == "" {
		log.Fatal("You must specify DbEncryptionSalt for security reasons")
	} else {
		DB_ENCRYPTION_SALT = options.DbEncryptionSalt
	}
	if options.AuthcodeLifeInMs == 0 {
		AUTHCODE_LIFE_IN_MS = AUTHCODE_LIFE_IN_MS_DEFAULT
	} else {
		AUTHCODE_LIFE_IN_MS = options.AuthcodeLifeInMs
	}
	if options.AuthcodeLength == 0 {
		AUTHCODE_LENGTH = AUTHCODE_LENGTH_DEFAULT
	} else {
		AUTHCODE_LENGTH = options.AuthcodeLength
	}
	if options.KeyExpirationInSeconds == 0 {
		KEY_EXPIRATION_IN_SECONDS = KEY_EXPIRATION_IN_SECONDS_DEFAULT
	} else {
		KEY_EXPIRATION_IN_SECONDS = options.KeyExpirationInSeconds
	}
	if options.Pbkdf2HashIterations == 0 {
		PBKDF2_HASH_ITERATIONS = PBKDF2_HASH_ITERATIONS_DEFAULT
	} else {
		PBKDF2_HASH_ITERATIONS = options.Pbkdf2HashIterations
	}
	if options.Pbkdf2HashBites == 0 {
		PBKDF2_HASH_BITES = PBKDF2_HASH_BITES_DEFAULT
	} else {
		PBKDF2_HASH_BITES = options.Pbkdf2HashBites
	}
	if len(DB_ENCRYPTION_SALT) != 32 {
		log.Fatal("DbEncryptionSalt size must be 32 bits long")
	}

	handshakejscrypter.Setup(DB_ENCRYPTION_SALT)

	redisurl, err := redisurlparser.Parse(redis_url_string)
	if err != nil {
		log.Fatal(err)
	}

	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisurl.Host+":"+redisurl.Port)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}

			if redisurl.Password != "" {
				if _, err := c.Do("AUTH", redisurl.Password); err != nil {
					c.Close()
					log.Fatal(err)
					return nil, err
				}
			}
			return c, err
		},
	}

	//c, err = redis.Dial("tcp", redisurl.Host+":"+redisurl.Port)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//if redisurl.Password != "" {
	//	if _, err := conn.Do("AUTH", redisurl.Password); err != nil {
	//		conn.Close()
	//		log.Fatal(err)
	//	}
	//}
}

func AppsCreate(app map[string]interface{}) (map[string]interface{}, *handshakejserrors.LogicError) {
	var app_name string
	if str, ok := app["app_name"].(string); ok {
		app_name = strings.Replace(str, " ", "", -1)
	} else {
		app_name = ""
	}
	if app_name == "" {
		logic_error := &handshakejserrors.LogicError{"required", "app_name", "app_name cannot be blank"}
		return app, logic_error
	}
	app["app_name"] = app_name

	generated_salt := uniuri.NewLen(20)
	if app["salt"] == nil {
		app["salt"] = generated_salt
	}
	if app["salt"].(string) == "" {
		app["salt"] = generated_salt
	}

	key := "apps/" + app["app_name"].(string)
	err := validateAppDoesNotExist(key)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"not_unique", "app_name", "app_name must be unique"}
		return app, logic_error
	}
	err = addAppToApps(app["app_name"].(string))
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return nil, logic_error
	}
	err = saveApp(key, app)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return nil, logic_error
	}

	return app, nil
}

func IdentitiesConfirm(identity map[string]interface{}) (map[string]interface{}, *handshakejserrors.LogicError) {
	app_name, logic_error := checkAppNamePresent(identity)
	if logic_error != nil {
		return identity, logic_error
	}
	identity["app_name"] = app_name

	email, logic_error := checkEmailPresent(identity)
	if logic_error != nil {
		return identity, logic_error
	}
	identity["email"] = email

	authcode, logic_error := checkAuthcodePresent(identity)
	if logic_error != nil {
		return identity, logic_error
	}
	identity["authcode"] = authcode

	app_name_key := "apps/" + identity["app_name"].(string)
	key := app_name_key + "/identities/" + identity["email"].(string)

	err := validateAppExists(app_name_key)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"not_found", "app_name", "app_name could not be found"}
		return identity, logic_error
	}
	err = validateIdentityExists(key)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"not_found", "email", "email could not be found"}
		return identity, logic_error
	}

	var r struct {
		Email             string `redis:"email"`
		AppName           string `redis:"app_name"`
		Authcode          string `redis:"authcode"`
		AuthcodeExpiredAt string `redis:"authcode_expired_at"`
	}

	conn := Conn()
	defer conn.Close()
	values, err := redis.Values(conn.Do("HGETALL", key))
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return identity, logic_error
	}
	err = redis.ScanStruct(values, &r)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return identity, logic_error
	}

	email = r.Email
	res_authcode := r.Authcode
	res_authcode_expired_at := r.AuthcodeExpiredAt

	current_ms_epoch_time := (time.Now().Unix() * 1000)
	res_authcode_expired_at_int64, err := strconv.ParseInt(res_authcode_expired_at, 10, 64)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return identity, logic_error
	}

	if len(res_authcode) > 0 && res_authcode == authcode {
		if res_authcode_expired_at_int64 < current_ms_epoch_time {
			logic_error := &handshakejserrors.LogicError{"expired", "authcode", "authcode has expired. request another one."}
			return identity, logic_error
		}

		app_salt, err := redis.String(conn.Do("HGET", app_name_key, "salt"))
		if err != nil {
			logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
			return identity, logic_error
		}
		app_salt = decrypt(app_salt)

		hash := pbkdf2.Key([]byte(email), []byte(app_salt), PBKDF2_HASH_ITERATIONS, PBKDF2_HASH_BITES, sha1.New)
		identity["hash"] = hex.EncodeToString(hash)

		return identity, nil
	} else {
		logic_error := &handshakejserrors.LogicError{"incorrect", "authcode", "the authcode was incorrect"}
		return identity, logic_error
	}
}

func IdentitiesCreate(identity map[string]interface{}) (map[string]interface{}, *handshakejserrors.LogicError) {
	app_name, logic_error := checkAppNamePresent(identity)
	if logic_error != nil {
		return identity, logic_error
	}
	identity["app_name"] = app_name

	email, logic_error := checkEmailPresent(identity)
	if logic_error != nil {
		return identity, logic_error
	}
	identity["email"] = email

	app_name_key := "apps/" + identity["app_name"].(string)
	key := app_name_key + "/identities/" + identity["email"].(string)

	err := validateAppExists(app_name_key)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"not_found", "app_name", "app_name could not be found"}
		return identity, logic_error
	}
	err = addIdentityToIdentities(app_name_key, identity["email"].(string))
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return identity, logic_error
	}
	err = saveIdentity(key, identity)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return nil, logic_error
	}

	return identity, nil
}

func saveApp(key string, app map[string]interface{}) error {
	app_to_save := make(map[string]interface{})
	for k, v := range app {
		app_to_save[k] = v
	}

	app_to_save["salt"] = encrypt(app_to_save["salt"].(string))

	args := []interface{}{key}
	for k, v := range app_to_save {
		args = append(args, k, v)
	}
	conn := Conn()
	defer conn.Close()
	_, err := conn.Do("HMSET", args...)
	if err != nil {
		return err
	}

	return nil
}

func addAppToApps(app_name string) error {
	conn := Conn()
	defer conn.Close()
	_, err := conn.Do("SADD", "apps", app_name)
	if err != nil {
		return err
	}

	return nil
}

func validateAppDoesNotExist(key string) error {
	conn := Conn()
	defer conn.Close()
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		log.Printf("ERROR " + err.Error())
		return err
	}
	if exists == true {
		err = errors.New("That app_name already exists.")
		return err
	}

	return nil
}

func validateAppExists(key string) error {
	conn := Conn()
	defer conn.Close()
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		log.Printf("ERROR " + err.Error())
		return err
	}
	if !exists {
		err = errors.New("That app_name does not exist.")
		return err
	}

	return nil
}

func validateIdentityExists(key string) error {
	conn := Conn()
	defer conn.Close()
	res, err := conn.Do("EXISTS", key)
	if err != nil {
		log.Printf("ERROR " + err.Error())
		return err
	}
	if res.(int64) != 1 {
		err = errors.New("That identity does not exist.")
		return err
	}

	return nil
}
func addIdentityToIdentities(app_name_key string, email string) error {
	conn := Conn()
	defer conn.Close()
	_, err := conn.Do("SADD", app_name_key+"/identities", email)
	if err != nil {
		log.Printf("ERROR " + err.Error())
		return err
	}

	return nil
}

func saveIdentity(key string, identity map[string]interface{}) error {
	rand.Seed(time.Now().UnixNano())
	authcode, err := randomAuthCode()
	if err != nil {
		log.Printf("ERROR " + err.Error())
		return err
	}
	identity["authcode"] = authcode
	unixtime := (time.Now().Unix() * 1000) + AUTHCODE_LIFE_IN_MS
	identity["authcode_expired_at"] = strconv.FormatInt(unixtime, BASE_10)

	args := []interface{}{key}
	for k, v := range identity {
		args = append(args, k, v)
	}
	conn := Conn()
	defer conn.Close()
	_, err = conn.Do("HMSET", args...)
	if err != nil {
		log.Printf("ERROR " + err.Error())
		return err
	}
	_, err = conn.Do("EXPIRE", key, KEY_EXPIRATION_IN_SECONDS)
	if err != nil {
		log.Printf("ERROR " + err.Error())
		return err
	}

	return nil
}

func randomAuthCode() (string, error) {
	rand.Seed(time.Now().UnixNano())
	var buffer bytes.Buffer

	for i := 1; i <= AUTHCODE_LENGTH; i++ {
		random_number := int64(rand.Intn(10))
		number_as_string := strconv.FormatInt(random_number, BASE_10)
		buffer.WriteString(number_as_string)
	}

	return buffer.String(), nil
}

func checkAppNamePresent(identity map[string]interface{}) (string, *handshakejserrors.LogicError) {
	var app_name string
	if str, ok := identity["app_name"].(string); ok {
		app_name = strings.Replace(str, " ", "", -1)
	} else {
		app_name = ""
	}
	if app_name == "" {
		logic_error := &handshakejserrors.LogicError{"required", "app_name", "app_name cannot be blank"}
		return app_name, logic_error
	}

	return app_name, nil
}

func checkEmailPresent(identity map[string]interface{}) (string, *handshakejserrors.LogicError) {
	var email string
	if str, ok := identity["email"].(string); ok {
		email = strings.Replace(str, " ", "", -1)
	} else {
		email = ""
	}
	if email == "" {
		logic_error := &handshakejserrors.LogicError{"required", "email", "email cannot be blank"}
		return email, logic_error
	}

	return email, nil
}

func checkAuthcodePresent(identity map[string]interface{}) (string, *handshakejserrors.LogicError) {
	var authcode string
	if str, ok := identity["authcode"].(string); ok {
		authcode = strings.Replace(str, " ", "", -1)
	} else {
		authcode = ""
	}
	if authcode == "" {
		logic_error := &handshakejserrors.LogicError{"required", "authcode", "authcode cannot be blank"}
		return authcode, logic_error
	}

	return authcode, nil
}

func Conn() redis.Conn {
	return pool.Get()
}

func encrypt(text string) string {
	return handshakejscrypter.Encrypt(text)
}

func decrypt(text string) string {
	return handshakejscrypter.Decrypt(text)
}

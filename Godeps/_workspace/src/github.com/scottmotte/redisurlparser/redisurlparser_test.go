package redisurlparser_test

import (
	redisurlparser "github.com/scottmotte/redisurlparser"
	"testing"
)

const REDIS_URL = "redis://redistogo:64cfea0093507536a374ab6ad40f8263@angelfish.redistogo.com:10039/"
const LOCAL_REDIS_URL = "redis://127.0.0.1:6379"

func TestParse(t *testing.T) {
	result, err := redisurlparser.Parse(REDIS_URL)
	if err != nil {
		t.Errorf("Error", err)
	}
	if result.Password != "64cfea0093507536a374ab6ad40f8263" {
		t.Errorf("Password incorrect")
	}
	if result.Username != "redistogo" {
		t.Errorf("Username incorrect")
	}
	if result.Host != "angelfish.redistogo.com" {
		t.Errorf("Host incorrect")
	}
	if result.Port != "10039" {
		t.Errorf("Port incorrect")
	}
}

func TestParseOnLocal(t *testing.T) {
	result, err := redisurlparser.Parse(LOCAL_REDIS_URL)
	if err != nil {
		t.Errorf("Error", err)
	}

	if result.Password != "" {
		t.Errorf("Password incorrect")
	}
	if result.Username != "" {
		t.Errorf("Username incorrect")
	}
	if result.Host != "127.0.0.1" {
		t.Errorf("Host incorrect")
	}
	if result.Port != "6379" {
		t.Errorf("Port incorrect")
	}
}

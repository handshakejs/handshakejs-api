package handshakejserrors_test

import (
	"../handshakejserrors"
	"testing"
)

func TestSetup(t *testing.T) {
	logic_error := handshakejserrors.LogicError{"required", "authcode", "authcode cannot be blank"}

	if logic_error.Code != "required" {
		t.Errorf("Error", logic_error)
	}
	if logic_error.Field != "authcode" {
		t.Errorf("Error", logic_error)
	}
	if logic_error.Message != "authcode cannot be blank" {
		t.Errorf("Error", logic_error)
	}
}

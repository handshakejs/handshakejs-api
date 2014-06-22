package handshakejscrypter_test

import (
	"../handshakejscrypter"
	"log"
	"testing"
)

const (
	KEY                = "somesecretsaltthatis32characters" //32 bytes
	ORIGINAL_PLAINTEXT = "some really long plaintext"
)

func TestEncryptionAndDecryption(t *testing.T) {
	handshakejscrypter.Setup(KEY)

	cipher := handshakejscrypter.Encrypt(ORIGINAL_PLAINTEXT)
	log.Println(cipher)

	plaintext := handshakejscrypter.Decrypt(cipher)
	log.Println(plaintext)

	if plaintext != ORIGINAL_PLAINTEXT {
		t.Errorf("Incorrect decrypted plaintext: " + plaintext)
	}
}

func TestEncryptionAndDecryption2(t *testing.T) {
	handshakejscrypter.Setup(KEY)

	plaintext := handshakejscrypter.Decrypt("4c916c441a8004057e69183eb11f74bdf1a0c6c5383c774ccf8bd69dd2fd4dba98dc3b6f16afb44a79416d4c")
	log.Println(plaintext)

	if plaintext != "qt5JjmQWIR3j1pZKQYjA" {
		t.Errorf("Incorrect decrypted plaintext: " + plaintext)
	}
}

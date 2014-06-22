package handshakejscrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"log"
)

var (
	KEY []byte
)

func Setup(key string) {
	if len(key) != 32 {
		log.Fatal("Key size must be 32 bits long")
	}
	KEY = []byte(key)
}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Encrypt(content string) string {
	text := []byte(content)
	block, err := aes.NewCipher(KEY)
	if err != nil {
		panic(err)
	}
	b := encodeBase64(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	hextext := hex.EncodeToString(ciphertext)
	return hextext
}

func Decrypt(content string) string {
	text, err := hex.DecodeString(content)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(KEY)
	if err != nil {
		panic(err)
	}
	if len(text) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	return string(decodeBase64(string(text)))
}

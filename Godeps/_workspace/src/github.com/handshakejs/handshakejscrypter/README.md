# handshakejscrypter

<img src="https://raw.githubusercontent.com/handshakejs/handshakejscrypter/master/handshakejscrypter.gif" alt="handshakejscrypter" align="right" width="270" />

[![BuildStatus](https://travis-ci.org/handshakejs/handshakejscrypter.png?branch=master)](https://travis-ci.org/handshakejs/handshakejscrypter)

Utility to encrypt and decrypt sensitive data for handshakejs. Code largely taken from [here](http://stackoverflow.com/questions/18817336/golang-encrypting-a-string-with-aes-and-base64).

This library is part of the larger [Handshake.js ecosystem](https://github.com/handshakejs).

## Usage

```go
package main

import (
  "fmt"
  "github.com/handshakejs/handshakejscrypter"
)

func main() {
  handshakejscrypter.Setup("somesecretsaltthatis32characters") // 32 bytes

  ciphertext := handshakejscrypter.Encrypt("some text to encrypt")
  fmt.Println(ciphertext)

  plaintext := handshakejscrypter.Decrypt(ciphertext)
  fmt.Println(plaintext)
}
```

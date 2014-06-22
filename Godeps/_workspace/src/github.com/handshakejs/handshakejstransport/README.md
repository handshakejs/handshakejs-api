# handshakejstransport

<img src="https://raw.githubusercontent.com/scottmotte/handshakejstransport/master/handshakejslogictransport.gif" alt="handshakejstransport" align="right" />

[![BuildStatus](https://travis-ci.org/handshakejs/handshakejstransport.png?branch=master)](https://travis-ci.org/handshakejs/handshakejstransport)

Transport mechanisms for delivering handshakejs authcodes to people.

This library is part of the larger [Handshake.js ecosystem](https://github.com/handshakejs).

## Usage

```go
package main

import (
  "fmt"
  handshakejstransport "github.com/handshakejs/handshakejstransport"
)

func main() {
  options := handshakejstransport.Options{"smtp.sendgrid.net", "587", "username", "password"}
  handshakejstransport.Setup(options)

  logic_error := handshakejstransport.ViaEmail("person0@mailinator.com", "from@yourapp.com", "Your authcode is 1234", "This is the text of the email", "This is the <b>html</b> of the email")
  if logic_err != nil {
    fmt.Println(logic_error)
  }
}
```

### Setup

Sets up the configuration.

```go
options := handshakejstransport.Options{SmtpAddress: "smtp.sendgrid.net", SmtpPort: "587", SmtpUsername: "username", SmtpPassword: "password"}
handshakejstransport.Setup(options)
```

### ViaEmail

Deliver authcode by way of email.

```go
logic_error := handshakejstransport.ViaEmail(to, from, subject, text, html)
if logic_error != nil {
  fmt.Println(logic_error)
}
```

## Installation

```
go get github.com/handshakejs/handshakejstransport
```

## Running Tests

```
go test -v
```


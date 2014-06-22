# handshakejslogic

<img src="https://raw.githubusercontent.com/handshakejs/handshakejslogic/master/handshakejslogic.gif" alt="handshakejslogic" align="right" width="190" />

[![BuildStatus](https://travis-ci.org/handshakejs/handshakejslogic.png?branch=master)](https://travis-ci.org/handshakejs/handshakejslogic)

Logic for saving handshakejs data to the redis database.

This library is part of the larger [Handshake.js ecosystem](https://github.com/handshakejs).

## Usage

```go
package main

import (
  "fmt"
  "github.com/handshakejs/handshakejslogic"
)

func main() {
  handshakejslogic.Setup("redis://127.0.0.1:6379", handshakejslogic.Options{})

  app := map[string]interface{}{"email": "email@myapp.com", "app_name": "myapp"}
  result, logic_error := handshakejslogic.AppsCreate(app)
  if logic_error != nil {
    fmt.Println(logic_error)
  }
  fmt.Println(result)
}
```

### Setup

Connects to Redis.

```go
options := handshakejslogic.Options{}
handshakejslogic.Setup("redis://127.0.0.1.6379", options)
```

### AppsCreate

```go
app := map[string]interface{}{"email": "email@myapp.com", "app_name": "myapp"}
result, logic_error := handshakejslogic.AppsCreate(app)
```

### IdentitiesCreate

```go
identity := map[string]interface{}{"email": "user@email.com", "app_name": "myapp"}
result, logic_error := handshakejslogic.IdentitiesCreate(identity)
```

## Installation

```
go get github.com/handshakejs/handshakejslogic
```

## Running Tests

```
go test -v
```

## Database Schema Details (using Redis)

Handshakejslogic uses a purposely simple database schema - as simple as possible. If you know a simpler approach, even better, please let me know or share as a pull request. 

Handshakejslogic uses Redis because of its light footprint, ephemeral nature, and lack of migrations.

* apps - collection of keys with all the app_names in there. SADD
* apps/myappname - hash with all the data in there. HSET or HMSET
* apps/theappname/identities - collection of keys with all the identities' emails in there. SADD
* apps/theappname/identities/emailaddress HSET or HMSET

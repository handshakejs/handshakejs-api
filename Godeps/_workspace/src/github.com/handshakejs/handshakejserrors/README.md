# handshakejserrors

<img src="https://raw.githubusercontent.com/handshakejs/handshakejserrors/master/handshakejserrors.gif" alt="handshakejserrors" align="right" width="190" />

[![BuildStatus](https://travis-ci.org/handshakejs/handshakejserrors.png?branch=master)](https://travis-ci.org/handshakejs/handshakejserrors)

Handshakejs error handling. Re-use in various handshakejs minor libraries.

This library is part of the larger [Handshake.js ecosystem](https://github.com/handshakejs).

## Usage

```go
package main

import (
  "fmt"
  "github.com/handshakejs/handshakejserrors"
)

func SomeFunction() *handshakejserrors.LogicError {
  // do something here and return LogicError type
}
```

### LogicError

```
*handshakejserrors.LogicError
```



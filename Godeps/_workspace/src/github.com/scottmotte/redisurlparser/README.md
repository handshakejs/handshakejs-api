# redisurlparser

Redis url parser is a go library for parsing redis://.. URLs into struct friendly options.

[![BuildStatus](https://travis-ci.org/scottmotte/redisurlparser.png?branch=master)](https://travis-ci.org/scottmotte/redisurlparser)

## Usage

```go
package main

import (
  "fmt"
  "github.com/scottmotte/redisurlparser"
)

func main() {
  url := "redis://redistogo:64cfea0093507536a374ba6ad40f8463@angelfish.redistogo.com:10049/"

  result, err := redisurlparser.Parse(url)
  if err != nil {
    fmt.Println(err)
  }

  fmt.Println(result.Username)
  fmt.Println(result.Password)
  fmt.Println(result.Host)
  fmt.Println(result.Port)

  fmt.Println(result)
}
```






# A Go ZeroLog Module Wrapper

A Go wrapper module for Zerolog module with some default which I do repeatedly on all my applications.

`zerolog_wrapper` is a wrapper  modules with some default logs configuration which helpedso we can get started with 
logging in Golang right away. This will ensure we follow DRY principle not to duplicate the same code in multiple applications.

## Usage

```go
package main

import (
    log "github.com/ashokrajar/zerolog_wrapper"
)

func init() {
    log.InitLog(log.TraceLevel, "dev")
}

func main() {
    log.Info().Msg("hello world")
}
```
Output
```shell
{"time":1494567715,"level":"info","message":"hello world"}
```

### How to add fields
```go
log.Info().Str("foo", "bar").Msg("hello world")
```

Output
```shell
{"time":1494567715,"level":"info","message":"hello world","foo":"bar"}
```

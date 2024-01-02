# midfabs

[![GoDoc](https://godoc.org/github.com/novrin/midfabs?status.svg)](https://pkg.go.dev/github.com/novrin/midfabs) 
![tests](https://github.com/novrin/midfabs/workflows/tests/badge.svg) 
[![Go Report Card](https://goreportcard.com/badge/github.com/novrin/midfabs)](https://goreportcard.com/report/github.com/novrin/midfabs)

`midfabs` is a Go package for prefabricated HTTP middleware.

### Installation

```shell
go get github.com/novrin/midfabs
``` 

## Usage

```go
package main

import (
	"log/slog"
	"net/http"

	"github.com/novrin/midfabs"
	"github.com/novrin/midway" // package for arranging middleware
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func broker(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("broker caught panic"))
}

func main() {
	// SecureHeaders sets sane security header defaults. Give it a CSP policy.
	defaultSecured := midfabs.SecureHeaders("default-src 'self';")

	// AccessLogger uses a slog.Logger type to log request and response details.
	logged := midfabs.AccessLogger(slog.Default(), "ACCESS")

	// PanicBroker takes a broker to gracefully handle panics in the given handler.
	panicRecovery := midfabs.PanicBroker(http.HandlerFunc(broker))

	queued := midway.Queue(panicRecovery, defaultSecured, logged)
	app := http.HandlerFunc(hello)
	http.ListenAndServe(":1313", queued(app))
	// serves logged(defaultSecured(panicRecovery(app)))
}
```

## License

Copyright (c) 2023-present [novrin](https://github.com/novrin)

Licensed under [MIT License](./LICENSE)
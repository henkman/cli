package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	_port uint
)

func init() {
	flag.UintVar(&_port, "p", 8080, "port")
	flag.Parse()
}

func main() {
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", _port), http.FileServer(http.Dir("./")))
}

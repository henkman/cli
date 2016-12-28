package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/kardianos/osext"
)

var (
	_port  uint
	_https bool
)

func init() {
	flag.UintVar(&_port, "p", 8080, "port")
	flag.BoolVar(&_https, "s", false, "https")
	flag.Parse()
}

func main() {
	addr := fmt.Sprintf("0.0.0.0:%d", _port)
	fs := http.FileServer(http.Dir("./"))
	if !_https {
		log.Fatal(http.ListenAndServe(addr, fs))
	} else {
		exe, _ := osext.Executable()
		d := filepath.Dir(exe)
		log.Fatal(http.ListenAndServeTLS(addr, filepath.Join(d, "cert.pem"), filepath.Join(d, "key.pem"), fs))
	}

}

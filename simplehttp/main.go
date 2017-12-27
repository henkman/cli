package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/kardianos/osext"
	"github.com/valyala/fasthttp"
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
	fs := &fasthttp.FS{
		Root:               "./",
		GenerateIndexPages: true,
		Compress:           true,
		AcceptByteRange:    true,
	}
	h := fs.NewRequestHandler()
	if !_https {
		log.Fatal(fasthttp.ListenAndServe(addr, h))
	} else {
		exe, _ := osext.Executable()
		d := filepath.Dir(exe)
		log.Fatal(fasthttp.ListenAndServeTLS(addr,
			filepath.Join(d, "cert.pem"),
			filepath.Join(d, "key.pem"), h))
	}
}

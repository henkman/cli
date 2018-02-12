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
	_host  string
	_port  uint
	_dir   string
	_https bool
)

func init() {
	flag.StringVar(&_host, "h", "0.0.0.0", "host")
	flag.UintVar(&_port, "p", 8080, "port")
	flag.StringVar(&_dir, "d", "./", "directory to server")
	flag.BoolVar(&_https, "s", false, "https")
	flag.Parse()
}

func main() {
	addr := fmt.Sprintf("%s:%d", _host, _port)
	fs := &fasthttp.FS{
		Root:               _dir,
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

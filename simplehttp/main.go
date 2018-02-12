package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/valyala/fasthttp"
)

func main() {
	var opts struct {
		Host  string
		Port  uint
		Dir   string
		Https bool
	}
	flag.StringVar(&opts.Host, "h", "0.0.0.0", "host")
	flag.UintVar(&opts.Port, "p", 8080, "port")
	flag.StringVar(&opts.Dir, "d", "./", "directory to server")
	flag.BoolVar(&opts.Https, "s", false, "https")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	fs := &fasthttp.FS{
		Root:               opts.Dir,
		GenerateIndexPages: true,
		Compress:           true,
		AcceptByteRange:    true,
	}
	h := fs.NewRequestHandler()
	if !opts.Https {
		log.Fatal(fasthttp.ListenAndServe(addr, h))
	} else {
		exe, _ := os.Executable()
		d := filepath.Dir(exe)
		log.Fatal(fasthttp.ListenAndServeTLS(addr,
			filepath.Join(d, "cert.pem"),
			filepath.Join(d, "key.pem"), h))
	}
}

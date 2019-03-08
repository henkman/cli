package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	if !opts.Https {
		log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir(opts.Dir))))
	} else {
		exe, _ := os.Executable()
		d := filepath.Dir(exe)
		log.Fatal(http.ListenAndServeTLS(addr,
			filepath.Join(d, "cert.pem"),
			filepath.Join(d, "key.pem"), http.FileServer(http.Dir(opts.Dir))))
	}
}

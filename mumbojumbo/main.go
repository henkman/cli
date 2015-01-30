package main
 
import (
	"encoding/json"
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"os"
	"strings"
)
 
type Config struct {
	Proxy   string
	Port    int
	Allowed []string
}
 
func readConfig(file string) (*Config, error) {
	fd, err := os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
 
	c := new(Config)
	jdec := json.NewDecoder(fd)
	jdec.Decode(c)
 
	return c, nil
}
 
func main() {
	fd, err := os.OpenFile("./log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	logger := log.New(fd, "", log.LstdFlags)
 
	conf, err := readConfig("./config.js")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("%+v\n", *conf)
 
	isAllowed := func(RemoteAddr string) bool {
		ip := strings.Split(RemoteAddr, ":")[0]
		for _, al := range conf.Allowed {
			if ip == al {
				return true
			}
		}
		return false
	}
 
	if conf.Proxy != "" {
		os.Setenv("HTTP_PROXY", conf.Proxy)
	}
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			logger.Printf("%s %s from %s", r.Method, r.URL, r.RemoteAddr)
			if !isAllowed(r.RemoteAddr) {
				return r, goproxy.NewResponse(r,
					goproxy.ContentTypeText, http.StatusForbidden,
					"Not allowed")
			}
			return r, nil
		})
	logger.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), proxy))
}
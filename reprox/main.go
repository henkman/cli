package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func main() {
	var opts struct {
		Listen    string
		Remote    string
		BasicAuth struct {
			User string
			Pass string
		}
	}
	flag.StringVar(&opts.Listen, "l", "127.0.0.1:1234", "local listen address")
	flag.StringVar(&opts.Remote, "r", "", "remote address")
	flag.StringVar(&opts.BasicAuth.User, "u", "", "basic auth user")
	flag.StringVar(&opts.BasicAuth.Pass, "p", "", "basic auth pass")
	flag.Parse()
	if opts.Remote == "" {
		flag.Usage()
		return
	}
	remote, err := url.Parse(opts.Remote)
	if err != nil {
		fmt.Println(err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		query := remote.RawQuery
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = singleJoiningSlash(remote.Path, req.URL.Path)
		if query == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = query + req.URL.RawQuery
		} else {
			req.URL.RawQuery = query + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}
		if opts.BasicAuth.User != "" {
			req.SetBasicAuth(opts.BasicAuth.User, opts.BasicAuth.Pass)
		}
	}
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	log.Fatal(http.ListenAndServe(opts.Listen, proxy))
}

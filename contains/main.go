package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

func main() {
	var opts struct {
		File        string
		Regex       string
		Print       string
		Ignorecase  bool
		Stoponfirst bool
	}
	flag.StringVar(&opts.File, "f", "", "file to search in")
	flag.StringVar(&opts.Regex, "r", "", "search regex")
	flag.StringVar(&opts.Print, "p", "", "instead of whole match print this (can contain $0..$N to refer to groups)")
	flag.BoolVar(&opts.Ignorecase, "c", false, "ignore case (prepends (?i) to regex)")
	flag.BoolVar(&opts.Stoponfirst, "s", false, "stop when first match is found")
	flag.Parse()

	if opts.File == "" || opts.Regex == "" {
		flag.Usage()
		return
	}
	if opts.Ignorecase {
		opts.Regex = "(?mi)" + opts.Regex
	} else {
		opts.Regex = "(?m)" + opts.Regex
	}
	re, err := regexp.Compile(opts.Regex)
	if err != nil {
		fmt.Println(err)
		return
	}
	fd, err := os.OpenFile(opts.File, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fd.Close()
	c, err := ioutil.ReadAll(fd)
	if err != nil {
		fmt.Println(err)
		return
	}
	var dst []byte
	if opts.Print != "" {
		dst = make([]byte, 0, 1*1024)
	}
	for {
		m := re.FindSubmatchIndex(c)
		if m == nil {
			break
		}
		if opts.Print != "" {
			dst = dst[:0]
			fmt.Println(string(re.ExpandString(dst, opts.Print, string(c), m)))
		} else {
			fmt.Println(string(c[m[0]:m[1]]))
		}
		if opts.Stoponfirst {
			break
		}
		c = c[m[1]:]
	}
}

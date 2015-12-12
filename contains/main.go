package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

var (
	_file        string
	_regex       string
	_print       string
	_ignorecase  bool
	_stoponfirst bool
)

func init() {
	flag.StringVar(&_file, "f", "", "file to search in")
	flag.StringVar(&_regex, "r", "", "search regex")
	flag.StringVar(&_print, "p", "", "instead of whole match print this (can contain $0..$N to refer to groups)")
	flag.BoolVar(&_ignorecase, "c", false, "ignore case (prepends (?i) to regex)")
	flag.BoolVar(&_stoponfirst, "s", false, "stop when first match is found")
	flag.Parse()
}

func main() {
	if _file == "" || _regex == "" {
		flag.Usage()
		return
	}
	if _ignorecase {
		_regex = "(?mi)" + _regex
	} else {
		_regex = "(?m)" + _regex
	}
	re, err := regexp.Compile(_regex)
	if err != nil {
		fmt.Println(err)
		return
	}
	fd, err := os.OpenFile(_file, os.O_RDONLY, 0600)
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
	if _print != "" {
		dst = make([]byte, 0, 1*1024)
	}
	for {
		m := re.FindSubmatchIndex(c)
		if m == nil {
			break
		}
		if _print != "" {
			dst = dst[:0]
			fmt.Println(string(re.ExpandString(dst, _print, string(c), m)))
		} else {
			fmt.Println(string(c[m[0]:m[1]]))
		}
		if _stoponfirst {
			break
		}
		c = c[m[1]:]
	}
}

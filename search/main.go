/*
	Simple search program

*/
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	_base               string
	_search             string
	_contains           string
	_stoponfirst        bool
	_absolutepath       bool
	_searchignorecase   bool
	_containsignorecase bool
)

func init() {
	flag.StringVar(&_search, "s", "", "regex of file(s) to search for")
	flag.StringVar(&_contains, "c", "", "regex of thing(s) that should be in file(s)")
	flag.StringVar(&_base, "b", ".", "the base directory")
	flag.BoolVar(&_stoponfirst, "f", false, "stop on the first occurance")
	flag.BoolVar(&_absolutepath, "a", false, "only show absolute paths")
	flag.BoolVar(&_searchignorecase, "is", false, "ignore case in file regex")
	flag.BoolVar(&_containsignorecase, "ic", false, "ignore case in contains regex")
	flag.Parse()
}

func main() {
	if _search == "" {
		flag.Usage()
		return
	}

	if _searchignorecase {
		_search = "(?i)" + _search
	}

	fr, err := regexp.Compile(_search)
	if err != nil {
		fmt.Println(err)
		return
	}

	var cr *regexp.Regexp
	if _contains != "" {
		if _containsignorecase {
			_contains = "(?i)" + _contains
		}

		cr, err = regexp.Compile(_contains)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if _base == "." {
		_base, err = os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	searchVisit := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !fr.MatchString(info.Name()) {
			return nil
		}
		if cr != nil {
			fd, err := os.OpenFile(path, os.O_RDONLY, 0600)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			defer fd.Close()
			bf := bufio.NewReader(fd)
			if cr.FindReaderIndex(bf) == nil {
				return nil
			}
		}
		if _absolutepath {
			fmt.Println(path)
		} else {
			rp, _ := filepath.Rel(_base, path)
			fmt.Println(rp)
		}
		if _stoponfirst {
			return errors.New("")
		}
		return nil
	}

	filepath.Walk(_base, searchVisit)
}

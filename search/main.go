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

func main() {
	var opts struct {
		Base               string
		Search             string
		Contains           string
		Stoponfirst        bool
		Absolutepath       bool
		Searchignorecase   bool
		Containsignorecase bool
	}
	flag.StringVar(&opts.Search, "s", "", "regex of file(s) to search for")
	flag.StringVar(&opts.Contains, "c", "", "regex of thing(s) that should be in file(s)")
	flag.StringVar(&opts.Base, "b", ".", "the base directory")
	flag.BoolVar(&opts.Stoponfirst, "f", false, "stop on the first occurance")
	flag.BoolVar(&opts.Absolutepath, "a", false, "only show absolute paths")
	flag.BoolVar(&opts.Searchignorecase, "is", false, "ignore case in file regex")
	flag.BoolVar(&opts.Containsignorecase, "ic", false, "ignore case in contains regex")
	flag.Parse()

	if opts.Search == "" {
		flag.Usage()
		return
	}

	if opts.Searchignorecase {
		opts.Search = "(?i)" + opts.Search
	}

	fr, err := regexp.Compile(opts.Search)
	if err != nil {
		fmt.Println(err)
		return
	}

	var cr *regexp.Regexp
	if opts.Contains != "" {
		if opts.Containsignorecase {
			opts.Contains = "(?i)" + opts.Contains
		}

		cr, err = regexp.Compile(opts.Contains)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if opts.Base == "." {
		opts.Base, err = os.Getwd()
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
		if opts.Absolutepath {
			fmt.Println(path)
		} else {
			rp, _ := filepath.Rel(opts.Base, path)
			fmt.Println(rp)
		}
		if opts.Stoponfirst {
			return errors.New("")
		}
		return nil
	}

	filepath.Walk(opts.Base, searchVisit)
}

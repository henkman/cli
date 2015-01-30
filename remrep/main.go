package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

var (
	_regex string
	_repl  string
	_test  bool
)

func init() {
	flag.StringVar(&_regex, "r", "", "regex")
	flag.StringVar(&_repl, "p", "", "replacement")
	flag.BoolVar(&_test, "t", false, "do not move for realz, only print")
	flag.Parse()
}

func lsdir(dir string) ([]os.FileInfo, error) {
	fd, err := os.Open(".")
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return fd.Readdir(-1)
}

func remrep(regex, repl, dir string, test bool) error {
	fis, err := lsdir(dir)
	if err != nil {
		return err
	}
	reRemRep, err := regexp.Compile(regex)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		if fi.IsDir() || !reRemRep.MatchString(fi.Name()) {
			continue
		}
		nn := reRemRep.ReplaceAllString(fi.Name(), repl)
		if test {
			fmt.Println(fi.Name(), "->", nn)
		} else {
			os.Rename(fi.Name(), nn)
		}
	}
	return nil
}

func main() {
	if _regex == "" || _repl == "" {
		flag.Usage()
		return
	}
	err := remrep(_regex, _repl, ".", _test)
	if err != nil {
		fmt.Println(err)
	}
}

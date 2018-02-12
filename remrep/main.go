package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

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
	var opts struct {
		Regex string
		Repl  string
		Test  bool
	}
	flag.StringVar(&opts.Regex, "r", "", "regex")
	flag.StringVar(&opts.Repl, "p", "", "replacement")
	flag.BoolVar(&opts.Test, "t", false, "do not move for realz, only print")
	flag.Parse()

	if opts.Regex == "" || opts.Repl == "" {
		flag.Usage()
		return
	}
	err := remrep(opts.Regex, opts.Repl, ".", opts.Test)
	if err != nil {
		fmt.Println(err)
	}
}

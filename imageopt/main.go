package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

type Optimizer struct {
	CanOptimize    *regexp.Regexp
	OptimizeParams []string
	FileIndex      int
	Executable     string
	Flag           string
}

var (
	_dir   string
	_quiet bool
	_help  bool
	opts   = []Optimizer{
		{
			CanOptimize:    regexp.MustCompile("(?i)\\.jpe?g$"),
			OptimizeParams: []string{"-quiet", "-s", "-m", "80", ""},
			FileIndex:      4,
			Flag:           "jpegoptim",
		},
		{
			CanOptimize:    regexp.MustCompile("(?i)\\.(?:png|bmp|gif|pnm|tiff)$"),
			OptimizeParams: []string{"-q", "-o7", "-fix", ""},
			FileIndex:      3,
			Flag:           "optipng",
		},
		{
			CanOptimize:    regexp.MustCompile("(?i)\\.svg$"),
			OptimizeParams: []string{"-q", ""},
			FileIndex:      1,
			Flag:           "svgo",
		},
	}
)

func init() {
	flag.StringVar(&_dir, "d", ".", "directory")
	flag.BoolVar(&_quiet, "q", false, "quiet")
	flag.BoolVar(&_help, "h", false, "help")
	for i, opt := range opts {
		flag.StringVar(&opts[i].Executable, opt.Flag, opt.Flag, opt.Flag)
	}
	flag.Parse()
}

func list(dir string) ([]os.FileInfo, error) {
	fd, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	return fd.Readdir(-1)
}

func optimize(dir string, quiet bool) error {
	fis, err := list(dir)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		if fi.IsDir() {
			optimize(filepath.Join(dir, fi.Name()), quiet)
			continue
		}
		for _, opt := range opts {
			if opt.CanOptimize.FindString(fi.Name()) != "" {
				file := filepath.Join(dir, fi.Name())
				if !quiet {
					fmt.Println("optimizing", file)
				}
				opt.OptimizeParams[opt.FileIndex] = file
				cmd := exec.Command(opt.Executable, opt.OptimizeParams...)
				if err := cmd.Run(); err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

func main() {
	if _help {
		flag.Usage()
		return
	}
	if err := optimize(_dir, _quiet); err != nil {
		panic(err)
	}
}

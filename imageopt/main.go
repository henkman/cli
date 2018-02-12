package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
)

type Optimizer struct {
	CanOptimize    *regexp.Regexp
	OptimizeParams []string
	FileIndex      int
	Executable     string
	Flag           string
}

type Task struct {
	Optimizer *Optimizer
	File      string
}

var optimizers = []Optimizer{
	{
		CanOptimize:    regexp.MustCompile("(?i)\\.jpe?g$"),
		OptimizeParams: []string{"-quiet", "-s", "-m", "80", ""},
		FileIndex:      4,
		Flag:           "jpegoptim",
	},
	{
		CanOptimize:    regexp.MustCompile("(?i)\\.(?:png|bmp|gif|pnm|tiff)$"),
		OptimizeParams: []string{"-q", "-o4", "-fix", ""},
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

func list(dir string) ([]os.FileInfo, error) {
	fd, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	fis, err := fd.Readdir(-1)
	fd.Close()
	return fis, err
}

func optimize(dir string, quiet bool, tasks chan Task) error {
	fis, err := list(dir)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		if fi.IsDir() {
			optimize(filepath.Join(dir, fi.Name()), quiet, tasks)
			continue
		}
		for i, opt := range optimizers {
			if opt.CanOptimize.FindString(fi.Name()) != "" {
				file := filepath.Join(dir, fi.Name())
				tasks <- Task{&optimizers[i], file}
				if !quiet {
					fmt.Println("optimizing", file)
				}
				break
			}
		}
	}
	return nil
}

func main() {
	var opts struct {
		Dir     string
		Quiet   bool
		Help    bool
		Workers uint
	}
	flag.StringVar(&opts.Dir, "d", ".", "directory")
	flag.BoolVar(&opts.Quiet, "q", false, "quiet")
	flag.UintVar(&opts.Workers, "w", uint(runtime.NumCPU()), "number of workers")
	flag.BoolVar(&opts.Help, "h", false, "help")
	for i, opt := range optimizers {
		flag.StringVar(&optimizers[i].Executable, opt.Flag, opt.Flag, opt.Flag)
	}
	flag.Parse()

	if opts.Help {
		flag.Usage()
		return
	}
	var wg sync.WaitGroup
	tasks := make(chan Task)
	wg.Add(int(opts.Workers))
	for i := uint(0); i < opts.Workers; i++ {
		go func() {
			params := []string{}
			for task := range tasks {
				opt := task.Optimizer
				params = append(params, opt.OptimizeParams...)
				params[opt.FileIndex] = task.File
				exec.Command(opt.Executable, params...).Run()
				params = params[:0]
			}
			wg.Done()
		}()
	}
	if err := optimize(opts.Dir, opts.Quiet, tasks); err != nil {
		panic(err)
	}
	close(tasks)
	wg.Wait()
}

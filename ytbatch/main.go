package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"runtime"
	"sync"
)

var (
	reUrl = regexp.MustCompile("(?m)^https?://\\S+$")
)

var (
	_ytdl  string
	_in    string
	_procs int
)

func init() {
	flag.StringVar(&_ytdl, "y", "youtube-dl", "youtube-dl binary")
	flag.StringVar(&_in, "i", "", "file with urls")
	flag.IntVar(&_procs, "p", 8, "number of processes")
	flag.Parse()
}

func worker(urls chan string, wg *sync.WaitGroup) {
	for url := range urls {
		cmd := exec.Command(_ytdl, url)
		cmd.Run()
	}
	wg.Done()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if _in == "" {
		flag.Usage()
		return
	}
	data, err := ioutil.ReadFile(_in)
	if err != nil {
		fmt.Println(err)
		return
	}
	m := reUrl.FindAllString(string(data), -1)
	if m == nil {
		fmt.Println("no urls found")
		return
	}
	var wg sync.WaitGroup
	wg.Add(_procs)
	urls := make(chan string)
	for i := 0; i < _procs; i++ {
		go worker(urls, &wg)
	}
	for _, url := range m {
		urls <- url
		fmt.Println("downloading", url)
	}
	close(urls)
	wg.Wait()
}

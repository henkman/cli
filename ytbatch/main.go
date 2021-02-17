package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

var reUrl = regexp.MustCompile("(?m)^https?://\\S+$")

func worker(urls chan string, ytdl string, ad []string, wg *sync.WaitGroup) {
	for url := range urls {
		cmd := exec.Command(ytdl, append(ad, url)...)
		cmd.Run()
	}
	wg.Done()
}

func main() {
	var opts struct {
		Ytdl  string
		In    string
		Ad    string
		Procs int
	}
	flag.StringVar(&opts.Ytdl, "y", "youtube-dl", "youtube-dl binary")
	flag.StringVar(&opts.In, "i", "", "file with urls")
	flag.StringVar(&opts.Ad, "a", "", "additional parameters")
	flag.IntVar(&opts.Procs, "p", 8, "number of processes")
	flag.Parse()

	if opts.In == "" {
		flag.Usage()
		return
	}
	data, err := os.ReadFile(opts.In)
	if err != nil {
		fmt.Println(err)
		return
	}
	m := reUrl.FindAllString(string(data), -1)
	if m == nil {
		fmt.Println("no urls found")
		return
	}
	ad := strings.Split(opts.Ad, " ")
	var wg sync.WaitGroup
	wg.Add(opts.Procs)
	urls := make(chan string)
	for i := 0; i < opts.Procs; i++ {
		go worker(urls, opts.Ytdl, ad, &wg)
	}
	for _, url := range m {
		urls <- url
		fmt.Println("downloading", url)
	}
	close(urls)
	wg.Wait()
}

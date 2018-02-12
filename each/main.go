package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	var opts struct {
		Do        string
		Re        string
		Test      bool
		Errignore bool
	}
	flag.StringVar(&opts.Do, "do", "", "command for each element, $0 is complete match, $1..$n are groups")
	flag.StringVar(&opts.Re, "re", "", "regex to use for partitioning, if empty whole string is used")
	flag.BoolVar(&opts.Errignore, "e", false, "ignore errors")
	flag.BoolVar(&opts.Test, "t", false, "don't do, just print")
	flag.Parse()
	if opts.Do == "" {
		flag.Usage()
		return
	}
	var reMatch *regexp.Regexp
	if opts.Re != "" {
		treMatch, err := regexp.Compile(opts.Re)
		if err != nil {
			log.Fatal(err)
		}
		reMatch = treMatch
	}
	dst := make([]byte, 0, 512)
	bin := bufio.NewReader(os.Stdin)
	end := false
	for !end {
		line, err := bin.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				end = true
			} else {
				log.Fatal(err)
			}
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var do string
		if reMatch != nil {
			m := reMatch.FindStringSubmatchIndex(line)
			if m == nil {
				continue
			}
			dst = dst[:0]
			do = string(reMatch.ExpandString(dst, opts.Do, line, m))
		} else {
			do = strings.Replace(opts.Do, "$0", line, -1)
		}
		if opts.Test {
			fmt.Println(do)
			continue
		}
		parts := strings.Fields(do)
		if parts == nil {
			continue
		}
		prog := parts[0]
		args := parts[1:]
		cmd := exec.Command(prog, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); !opts.Errignore && err != nil {
			log.Fatal(err)
		}
	}
}

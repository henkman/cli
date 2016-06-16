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

var (
	_do        string
	_re        string
	_test      bool
	_errignore bool
)

func init() {
	flag.StringVar(&_do, "do", "", "command for each element, $0 is complete match, $1..$n are groups")
	flag.StringVar(&_re, "re", "", "regex to use for partitioning, if empty whole string is used")
	flag.BoolVar(&_errignore, "e", false, "ignore errors")
	flag.BoolVar(&_test, "t", false, "don't do, just print")
	flag.Parse()
}

func main() {
	if _do == "" {
		flag.Usage()
		return
	}
	var reMatch *regexp.Regexp
	if _re != "" {
		treMatch, err := regexp.Compile(_re)
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
			do = string(reMatch.ExpandString(dst, _do, line, m))
		} else {
			do = strings.Replace(_do, "$0", line, -1)
		}
		if _test {
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
		if err := cmd.Run(); !_errignore && err != nil {
			log.Fatal(err)
		}
	}
}

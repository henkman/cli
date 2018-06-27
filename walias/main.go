package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	var opts struct {
		Dir     string
		Command string
		Alias   string
		Force   bool
	}
	exe, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		return
	}
	flag.StringVar(&opts.Dir, "d", filepath.Dir(exe), "directory where to place the bat")
	flag.StringVar(&opts.Command, "c", "", "the command to alias")
	flag.StringVar(&opts.Alias, "a", ".", "the alias for the command")
	flag.BoolVar(&opts.Force, "f", false, "force overwrite")
	flag.Parse()
	if opts.Command == "" {
		flag.Usage()
		return
	}
	abs := filepath.Join(opts.Dir, opts.Alias+".bat")
	if !opts.Force {
		if _, err := os.Stat(abs); err == nil {
			fmt.Println("alias already exists. use -f to overwrite")
			return
		}
	}
	if err := ioutil.WriteFile(abs, []byte(fmt.Sprintf("@%s %%*", opts.Command)), 0660); err != nil {
		fmt.Println(err)
	}
}

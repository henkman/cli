package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var opts struct {
		Decode bool
		In     string
		Out    string
	}
	flag.BoolVar(&opts.Decode, "d", false, "decode mode. omit if you want to encode")
	flag.StringVar(&opts.In, "i", "", "input file. leave empty for stdin")
	flag.StringVar(&opts.Out, "o", "", "output file. leave empty for stdout")
	flag.Parse()

	var src io.Reader
	if opts.In != "" {
		fd, err := os.Open(opts.In)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer fd.Close()
		src = fd
	} else {
		src = os.Stdin
	}

	var dst io.WriteCloser
	if opts.Out != "" {
		fd, err := os.OpenFile(opts.Out, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0750)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer fd.Close()
		dst = fd
	} else {
		dst = os.Stdout
	}

	if opts.Decode {
		src = base64.NewDecoder(base64.RawStdEncoding, src)
	} else {
		dst = base64.NewEncoder(base64.RawStdEncoding, dst)
		defer dst.Close()
	}

	if _, err := io.Copy(dst, src); err != nil {
		fmt.Println(err)
	}
}

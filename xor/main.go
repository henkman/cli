package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

func xor(data io.Reader, key io.ReadSeeker, out io.Writer) error {
	const BUF_SIZE = 32 * 1024

	// short path for key files
	// smaller than buffer
	if kf, ok := key.(*os.File); ok {
		fi, err := kf.Stat()
		if err != nil {
			return err
		}
		if fi.Size() < BUF_SIZE {
			skbuf := make([]byte, fi.Size())
			_, err := key.Read(skbuf)
			if err != nil {
				return err
			}
			key = bytes.NewReader(skbuf)
		}
	}

	var dbuf, kbuf [BUF_SIZE]byte
	for {
		n, err := data.Read(dbuf[:])
		if n > 0 {
			{
				o := 0
				for o != n {
					kn, err := key.Read(kbuf[o:n])
					o += kn
					if err != nil {
						if err == io.EOF {
							key.Seek(0, 0)
							continue
						}
						return err
					}
				}
			}
			for i := 0; i < n; i++ {
				dbuf[i] ^= kbuf[i]
			}
			_, err := out.Write(dbuf[:n])
			if err != nil {
				return err
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

func main() {
	var opts struct {
		Datafile string
		Data     string
		Keyfile  string
		Key      string
		Outfile  string
	}
	flag.StringVar(&opts.Datafile, "df", "", "datafile")
	flag.StringVar(&opts.Data, "d", "", "data")
	flag.StringVar(&opts.Keyfile, "kf", "", "keyfile")
	flag.StringVar(&opts.Key, "k", "", "key")
	flag.StringVar(&opts.Outfile, "of", "", "outfile")
	flag.Parse()

	if opts.Keyfile == "" && opts.Key == "" {
		flag.Usage()
		return
	}

	var data io.Reader
	if opts.Data != "" {
		data = bytes.NewBufferString(opts.Data)
	} else if opts.Datafile != "" {
		fd, err := os.Open(opts.Datafile)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer fd.Close()
		data = fd
	} else {
		data = os.Stdin
	}

	var key io.ReadSeeker
	if opts.Key != "" {
		key = bytes.NewReader([]byte(opts.Key))
	} else {
		fd, err := os.Open(opts.Keyfile)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer fd.Close()
		key = fd
	}

	var out io.Writer
	if opts.Outfile != "" {
		fd, err := os.OpenFile(opts.Outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			0600)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer fd.Close()
		out = fd
	} else {
		out = os.Stdout
	}

	err := xor(data, key, out)
	if err != nil {
		fmt.Println(err)
	}
}

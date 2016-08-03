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

var (
	_datafile string
	_data     string
	_keyfile  string
	_key      string
	_outfile  string
)

func init() {
	flag.StringVar(&_datafile, "df", "", "datafile")
	flag.StringVar(&_data, "d", "", "data")
	flag.StringVar(&_keyfile, "kf", "", "keyfile")
	flag.StringVar(&_key, "k", "", "key")
	flag.StringVar(&_outfile, "of", "", "outfile")
	flag.Parse()
}

func main() {
	if _keyfile == "" && _key == "" {
		flag.Usage()
		return
	}

	var data io.Reader
	if _data != "" {
		data = bytes.NewBufferString(_data)
	} else if _datafile != "" {
		fd, err := os.Open(_datafile)
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
	if _key != "" {
		key = bytes.NewReader([]byte(_key))
	} else {
		fd, err := os.Open(_keyfile)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer fd.Close()
		key = fd
	}

	var out io.Writer
	if _outfile != "" {
		fd, err := os.OpenFile(_outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
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

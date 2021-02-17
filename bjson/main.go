package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
)

func main() {
	var raw []byte
	if len(os.Args) == 2 {
		_raw, err := os.ReadFile(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		raw = _raw
	} else {
		_raw, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		raw = _raw
	}
	var out bytes.Buffer
	json.Indent(&out, raw, "", "\t")
	out.WriteTo(os.Stdout)
}

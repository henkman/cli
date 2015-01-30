package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	_number uint64
	_string string
	_base   uint64
	_digits string
)

func init() {
	flag.Uint64Var(&_number, "n", 0, "number")
	flag.StringVar(&_string, "s", "", "string")
	flag.Uint64Var(&_base, "b", 0, "base <= length of digits string")
	flag.StringVar(&_digits, "d",
		"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"digits to be used")
	flag.Parse()
}

func tobase(n, b uint64) string {
	s := ""
	for {
		r := n % b
		s = _digits[r:r+1] + s
		n /= b
		if n == 0 {
			break
		}
	}
	return s
}

func pow(n, e uint64) uint64 {
	var i, x uint64
	x = 1
	for i = 0; i < e; i++ {
		x *= n
	}
	return x
}

func frombase(s string, b uint64) (uint64, error) {
	var n, e uint64
	for i := len(s) - 1; i >= 0; i-- {
		v := strings.Index(_digits[:b], string(s[i]))
		if v == -1 {
			return 0, fmt.Errorf("unknown digit %c at position %d", s[i], i)
		}
		n += uint64(v) * pow(b, uint64(e))
		e++
	}
	return n, nil
}

func main() {
	if (_number == 0 && _string == "") ||
		_base == 0 || _base > uint64(len(_digits)) {
		flag.Usage()
		return
	}
	if _string != "" {
		n, err := frombase(_string, _base)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(n)
	} else {
		fmt.Println(tobase(_number, _base))
	}
}

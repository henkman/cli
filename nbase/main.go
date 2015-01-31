package main

import (
	"flag"
	"fmt"
	"math/big"
	"strings"
)

var (
	_number string
	_string string
	_base   uint64
	_digits string
)

func init() {
	flag.StringVar(&_number, "n", "", "number")
	flag.StringVar(&_string, "s", "", "string")
	flag.Uint64Var(&_base, "b", 0, "base <= length of digits string")
	flag.StringVar(&_digits, "d",
		"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"digits to be used")
	flag.Parse()
}

func tobase(n, b *big.Int) string {
	zero := big.NewInt(0)
	r := new(big.Int)
	s := ""
	for {
		r.Mod(n, b)
		s = _digits[r.Uint64():r.Uint64()+1] + s
		n.Div(n, b)
		if n.Cmp(zero) == 0 {
			break
		}
	}
	return s
}

func frombase(s string, b *big.Int) (*big.Int, error) {
	var e uint64
	n := new(big.Int)
	eb := new(big.Int)
	vb := new(big.Int)
	for i := len(s) - 1; i >= 0; i-- {
		v := strings.Index(_digits[:b.Uint64()], string(s[i]))
		if v == -1 {
			return nil, fmt.Errorf("unknown digit %c at position %d", s[i], i)
		}
		vb.SetUint64(uint64(v))
		eb.SetUint64(e)
		n.Add(n, vb.Mul(vb, eb.Exp(b, eb, nil)))
		e++
	}
	return n, nil
}

func main() {
	if (_number == "" && _string == "") ||
		_base == 0 || _base > uint64(len(_digits)) {
		flag.Usage()
		return
	}
	base := new(big.Int).SetUint64(_base)
	if _string != "" {
		n, err := frombase(_string, base)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(n)
	} else {
		n := new(big.Int)
		if _, ok := n.SetString(_number, 10); !ok {
			fmt.Println("number is invalid")
			return
		}
		fmt.Println(tobase(n, base))
	}
}

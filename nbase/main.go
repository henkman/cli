package main

import (
	"flag"
	"fmt"
	"math/big"
	"strings"
)

func tobase(n, b *big.Int, digits string) string {
	zero := big.NewInt(0)
	r := new(big.Int)
	s := ""
	for {
		r.Mod(n, b)
		s = digits[r.Uint64():r.Uint64()+1] + s
		n.Div(n, b)
		if n.Cmp(zero) == 0 {
			break
		}
	}
	return s
}

func frombase(s string, b *big.Int, digits string) (*big.Int, error) {
	var e uint64
	n := new(big.Int)
	eb := new(big.Int)
	vb := new(big.Int)
	for i := len(s) - 1; i >= 0; i-- {
		v := strings.Index(digits[:b.Uint64()], string(s[i]))
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
	var opts struct {
		Number string
		String string
		Base   uint64
		Digits string
	}
	flag.StringVar(&opts.Number, "n", "", "number")
	flag.StringVar(&opts.String, "s", "", "string")
	flag.Uint64Var(&opts.Base, "b", 0, "base <= length of digits string")
	flag.StringVar(&opts.Digits, "d",
		"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"digits to be used")
	flag.Parse()

	if (opts.Number == "" && opts.String == "") ||
		opts.Base == 0 || opts.Base > uint64(len(opts.Digits)) {
		flag.Usage()
		return
	}
	base := new(big.Int).SetUint64(opts.Base)
	if opts.String != "" {
		n, err := frombase(opts.String, base, opts.Digits)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(n)
	} else {
		n := new(big.Int)
		if _, ok := n.SetString(opts.Number, 10); !ok {
			fmt.Println("number is invalid")
			return
		}
		fmt.Println(tobase(n, base, opts.Digits))
	}
}

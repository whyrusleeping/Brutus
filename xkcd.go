package main

import (
	"crypto/skein"
	"runtime"
	"fmt"
	"encoding/hex"
	"math/rand"
)

func DifHash(a, good []byte) int {
	sum := 0
	for i := 0; i < len(a); i++ {
		l := a[i]
		r := good[i]

		for i := 0; i < 8; i++ {
			lw := l & 1
			rw := r & 1
			if lw == rw {
				sum++
			}
			l = l >> 1
			r = r >> 1
		}
	}
	return 1024 - sum
}

func makeSampleString() []byte {
	b := make([]byte, 128)
	k := 0
	i := 'a'
	for ; i <= 'z'; i++ {
		b[k] = byte(i)
		k++
	}
	for i = 'A'; i <= 'Z'; i++ {
		b[k] = byte(i)
		k++
	}
	for i = '0'; i <= '9'; i++ {
		b[k] = byte(i)
		k++
	}
	b[k] = '!'
	k++
	b[k] = '@'
	k++
	b[k] = '#'
	k++
	b[k] = '$'
	k++
	b[k] = '%'
	k++
	b[k] = '^'
	k++
	b[k] = '&'
	k++
	b[k] = '*'
	k++
	b[k] = '('
	k++
	b[k] = ')'
	k++
	b[k] = '-'
	k++
	b[k] = '_'
	k++
	b[k] = '='
	k++
	b[k] = '+'
	k++
	b[k] = '\''
	k++
	b[k] = '.'
	k++
	b[k] = ','
	k++
	b[k] = '<'
	k++
	b[k] = '>'
	k++
	b[k] = '/'
	k++
	b[k] = '?'
	k++
	return b[:k]
}

func RandString(l int, dict []byte) []byte {
	s := make([]byte, l)
	for i := 0; i < l; i++ {
		s[i] = dict[rand.Intn(len(dict))]
	}
	return s
}

func DiffFromString(gs , s []byte) int {
	b, _ := skein.New(skein.Skein1024, 1024)
	b.Write(s)
	m := b.DoFinal()
	return DifHash(m, gs)
}

func main() {
	runtime.GOMAXPROCS(4)
	THEGOOD,_ := hex.DecodeString("5b4da95f5fa08280fc9879df44f418c8f9f12ba424b7757de02bbdfbae0d4c4fdf9317c80cc5fe04c6429073466cf29706b8c25999ddd2f6540d4475cc977b87f4757be023f19b8f4035d7722886b78869826de916a79cf9c94cc79cd4347d24b567aa3e2390a573a373a48a5e676640c79cc70197e1c5e7f902fb53ca1858b6")
	ds := makeSampleString()
	record := 1024
	for {
		rs := RandString(rand.Intn(16), ds)
		dnum := DiffFromString(THEGOOD, rs)
		if dnum < record {
			fmt.Printf("%s %d\n", rs, dnum)
			record = dnum
		}
	}	
	fmt.Println("Stop?")
	for {}
}

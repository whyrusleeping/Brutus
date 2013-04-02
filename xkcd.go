package main

import (
	"github.com/whyrusleeping/GoSkein"
	"fmt"
	"encoding/hex"
	"math/rand"
	"time"
	//"runtime"
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
	return b[:k]
}

func RandString(l int, dict []byte, ret []byte) {
	for i := 0; i < l; i++ {
		ret[i] = dict[rand.Intn(len(dict))]
	}
}

func DiffFromString(b *skein.Skein, gs , s []byte) int {
	b.Write(s)
	m := b.DoFinal()
	r := DifHash(m, gs)
	skein.FreeBuf(m)
	return r
}

func Brute(num int, check, dict []byte) {
	//runtime.LockOSThread()
	buff := make([]byte, 32)
	t := time.Now()
	record := 1024
	count := int64(1)
	b, _ := skein.New(skein.Skein1024, 1024)
	for {
		if count % 1000000 == 0 {
			fmt.Println(float64(1000000) / float64(time.Now().Unix() - t.Unix()))
			t = time.Now()
		}
		n := rand.Intn(32)
		buff = buff[:n]
		RandString(n, dict, buff)
		dnum := DiffFromString(b, check, buff)
		if dnum < record {
			fmt.Printf("%d: %s %d\n", count, buff, dnum)
			record = dnum
		}
		count++
	}
}

func main() {
	//runtime.GOMAXPROCS(8)
	rand.Seed(time.Now().Unix())
	THEGOOD,_ := hex.DecodeString("5b4da95f5fa08280fc9879df44f418c8f9f12ba424b7757de02bbdfbae0d4c4fdf9317c80cc5fe04c6429073466cf29706b8c25999ddd2f6540d4475cc977b87f4757be023f19b8f4035d7722886b78869826de916a79cf9c94cc79cd4347d24b567aa3e2390a573a373a48a5e676640c79cc70197e1c5e7f902fb53ca1858b6")
	ds := makeSampleString()
	Brute(0, THEGOOD, ds)
}

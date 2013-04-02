package main

import (
	"github.com/whyrusleeping/GoSkein"
	"net"
	"fmt"
	"encoding/hex"
	"math/rand"
	"time"
	"math/big"
	"strings"
	"encoding/base64"
	"bytes"
)

const (
	A = 1 << iota
	B
	C
	D
	E
	F
	G
	H
)
func DifHash(a, good []byte) int {
	sum := 0
	for i := 0; i < len(a); i++ {
		r := a[i] ^ good[i]

		if r & A != 0 {
			sum++
		}
		if r & B != 0 {
			sum++
		}
		if r & C != 0 {
			sum++
		}
		if r & D != 0 {
			sum++
		}
		if r & E != 0 {
			sum++
		}
		if r & F != 0 {
			sum++
		}
		if r & G != 0 {
			sum++
		}
		if r & H != 0 {
			sum++
		}
	}
	return sum
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
	b[k] = '.'
	k++
	return b[:k]
}

func RandString(l int, dict []byte, ret []byte) {
	for i := 0; i < l; i++ {
		ret[i] = dict[rand.Intn(len(dict))]
		//ret[i] = byte(rand.Intn(95) + 32)
	}
}

func DiffFromString(b *skein.Skein, gs , s []byte) int {
	b.Write(s)
	m := b.DoFinal()
	r := DifHash(m, gs)
	skein.FreeBuf(m)
	return r
}

func SendToEric(word string, num int) {
	addr, _ := net.ResolveTCPAddr("tcp","hobosteaux.dyndns.org:9000")
	conn, _ := net.DialTCP("tcp", nil,addr)
	conn.Write([]byte(fmt.Sprintf("update;%d;%s", num, word)))
	conn.Close()
}

func RequestNewRange() (*big.Int, *big.Int) {
	//return big.NewInt(5000000000), big.NewInt(5002000000)
	addr, _ := net.ResolveTCPAddr("tcp","hobosteaux.dyndns.org:9000")
	conn, _ := net.DialTCP("tcp", nil,addr)
	conn.Write([]byte("ask"))
	rbuf := make([]byte, 256)
	n,_ := conn.Read(rbuf)
	rbuf = rbuf[:n]
	conn.Close()
	rng := strings.Split(string(rbuf),";")
	lo := big.NewInt(0)
	lo.SetString(rng[0],10)
	hi := big.NewInt(0)
	hi.SetString(rng[1],10)
	fmt.Printf("Now checking range %s to %s\n.", lo.String(), hi.String())
	return lo, hi
}

func SchedBrute(check []byte) {
	
	buff := make([]byte, 64)
	t := time.Now()
	count := int64(1)
	record := 1024
	b, _ := skein.New(skein.Skein1024, 1024)
	lo, hi := RequestNewRange()
	one := big.NewInt(1)
	ttrim := make([]byte, 2)
	ttrim[0] = 0
	ttrim[1] = '='
	trimset := string(ttrim)
	for {
		if count % 1000000 == 0 {
			fmt.Println(1000000.0 / (float64(time.Now().UnixNano() - t.UnixNano()) / 1000000000.0))
			t = time.Now()
		}
		//MAKE STRING HERE
		by := lo.Bytes()
		base64.StdEncoding.Encode(buff, by)
		subbuff := bytes.TrimRight(buff, trimset)

		//
		dnum := DiffFromString(b, check, subbuff)
		if dnum < record {
			entry := string(subbuff)
			fmt.Printf("%d: %s %d\n", count, entry, dnum)
			SendToEric(entry, dnum)

			record = dnum
		}
		if lo.Cmp(hi) == 0 {
			lo, hi = RequestNewRange()
		} else {
			lo.Add(lo, one)
		}
		count++
	}
}

func Brute(num int, check, dict []byte) {
	//runtime.LockOSThread()
	buff := make([]byte, 64)
	t := time.Now()
	record := 1024
	count := int64(1)
	b, _ := skein.New(skein.Skein1024, 1024)
	for {
		if count % 1000000 == 0 {
			fmt.Println(1000000.0 / (float64(time.Now().UnixNano() - t.UnixNano()) / 1000000000.0))
			t = time.Now()
		}
		n := rand.Intn(64)
		buff = buff[:n]
		RandString(n, dict, buff)
		dnum := DiffFromString(b, check, buff)
		if dnum < record {
			entry := string(buff)
			fmt.Printf("%d: %s %d\n", count, entry, dnum)
			//SendToEric(entry, dnum)
			record = dnum
		}
		count++
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	THEGOOD,_ := hex.DecodeString("5b4da95f5fa08280fc9879df44f418c8f9f12ba424b7757de02bbdfbae0d4c4fdf9317c80cc5fe04c6429073466cf29706b8c25999ddd2f6540d4475cc977b87f4757be023f19b8f4035d7722886b78869826de916a79cf9c94cc79cd4347d24b567aa3e2390a573a373a48a5e676640c79cc70197e1c5e7f902fb53ca1858b6")
	ds := makeSampleString()
	Brute(0, THEGOOD, ds)
}

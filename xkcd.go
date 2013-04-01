package main

import (
	"crypto/skein"
	"fmt"
	"encoding/hex"
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
	return sum
}

func DiffFromString(sk *skein.Skein, b []byte, s string) int {

}
func main() {
	THEGOOD,_ := hex.DecodeString("5b4da95f5fa08280fc9879df44f418c8f9f12ba424b7757de02bbdfbae0d4c4fdf9317c80cc5fe04c6429073466cf29706b8c25999ddd2f6540d4475cc977b87f4757be023f19b8f4035d7722886b78869826de916a79cf9c94cc79cd4347d24b567aa3e2390a573a373a48a5e676640c79cc70197e1c5e7f902fb53ca1858b6")


	b,_ := skein.New(skein.Skein1024, 1024)
	br := []byte("FISH")
	b.Write(br)
	m := b.DoFinal()
	fmt.Println(hex.EncodeToString(m))
	fmt.Println(DifHash(m, THEGOOD))
}

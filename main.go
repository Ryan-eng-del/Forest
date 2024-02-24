package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"unicode/utf8"
)

func main() {
	r, s := utf8.DecodeRune([]byte("a你好"))
	fmt.Println(string(r), s)

	// a := []int{1, 2, 5}
	b := make([]byte, 512, 512)
	
	
	for i := 0; i < len(b); i++ {
		b[i] = byte(1)
	}

	fmt.Println(b)
	// b = b[len(b):cap(b)]

	rz := bytes.NewReader(b)
	ci, err := io.ReadAll(rz)

	if err != nil {
		log.Println(err)
	}

	log.Println(len(ci), ci)
	fmt.Println(len(b))
	// fmt.Println(b[:12])
	// n := copy(b, a)
	// b = b[:len(b)+n]

	// fmt.Println(n, b, len(b))
	// fmt.Println(len(a), cap(a))


	// a = a[:10]
	// fmt.Println(len(a), cap(a))
}

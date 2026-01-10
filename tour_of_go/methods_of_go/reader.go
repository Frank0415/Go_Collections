package main

import (
	"fmt"
	"io"
	"strings"
)

type MyReader struct{}

func (mr MyReader) Read(buf []byte) (int, error) {
	for i := range len(buf) {
		buf[i] = 'A'
	}
	return len(buf), nil
}

func A_reader() {
	var r = MyReader{}

	for i := range(10){
		b := make([]byte, i+1)
		n,_ := r.Read(b)
		fmt.Printf("%q\n",b[:n])
	}
	
	
}

func Strings_reader() {
	r := strings.NewReader("Hello, Reader!")

	b := make([]byte, 8)
	for {
		n, err := r.Read(b)
		fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
		fmt.Printf("b[:n] = %q\n", b[:n])
		if err == io.EOF {
			break
		}
	}
}

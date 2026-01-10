package main

import (
	"fmt"
	"io"
	"os"
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

	for i := range 10 {
		b := make([]byte, i+1)
		n, _ := r.Read(b)
		fmt.Printf("%q\n", b[:n])
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

type rot13Reader struct {
	r io.Reader
}

func (rR rot13Reader) Read(buf []byte) (int, error) {
	n, err := rR.r.Read(buf)
	if err != nil {
		return n, err
	}

	for i := range n {
		b := buf[i]
		switch {
		case (b >= 'A' && b <= 'M') || (b >= 'a' && b <= 'm'):
			buf[i] += 13
		case (b >= 'N' && b <= 'Z') || (b >= 'n' && b <= 'z'):
			buf[i] -= 13
		default:
		}
	}
	return n, nil
}

func Test_rot() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stderr,&r)
}

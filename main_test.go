package main

import (
	"encoding/base64"
	"fmt"
	"github.com/tmthrgd/go-hex"
	"strconv"
	"testing"
)

func TestGetCallerFuncName2(t *testing.T) {

	b, err := base64.StdEncoding.DecodeString("MNjRwBctB0XiG4zNJugkWg==")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("hex: ", hex.EncodeToString(b))

	h := "30d8d1c0172d0745e21b8ccd26e8245a"

	m, err := strconv.ParseUint(h[:16], 16, 64)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("prefix: ", m)

	n, err := strconv.ParseUint(h[16:], 16, 64)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("suffix: ", n)
}

func TestReadMovies(t *testing.T) {
	mov, err := readMovies()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(mov))
}

func TestGetCallerFuncName(t *testing.T) {

	for i := 0; i < 10; i++ {
		fmt.Println(getCurServID())
	}

	for i := 0; i < 10; i++ {
		fmt.Println(getNextServID())
	}

	fmt.Println(getCurServName())
	fmt.Println(getNextServName())

}

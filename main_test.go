package main

import (
	"fmt"
	"testing"
)

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

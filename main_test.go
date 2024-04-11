package main

import (
	"fmt"
	"testing"
)

func TestReadMovies(t *testing.T) {
	if _, err := readMovies(); err != nil {
		t.Fatal(err)
	}
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

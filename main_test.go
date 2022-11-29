package main

import "testing"

func TestReadMovies(t *testing.T) {
	if _, err := readMovies(); err != nil {
		t.Fatal(err)
	}
}

package main

import (
	"fmt"

	"github.com/shasderias/cryptopals/padding"
)

const (
	c9Input          = "YELLOW SUBMARINE"
	c9ExpectedOutput = "YELLOW SUBMARINE\x04\x04\x04\x04"
)

func runC9() error {
	padded := padding.PKCS7([]byte(c9Input), 20)
	if string(padded) != c9ExpectedOutput {
		fmt.Printf("got %s; want %s", padded, c9ExpectedOutput)
	}
	fmt.Println(string(padded))
	return nil
}

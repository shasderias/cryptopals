package main

import (
	"encoding/base64"
	"fmt"
)

const (
	c1Input          = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	c1ExpectedOutput = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
)

func runC1() error {
	b64Str := base64.StdEncoding.EncodeToString(mustDecodeHexString(c1Input))
	if b64Str != c1ExpectedOutput {
		fmt.Printf("got %s; want %s\n", b64Str, c1ExpectedOutput)
		return nil
	}
	fmt.Println(b64Str)
	return nil
}

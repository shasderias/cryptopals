package main

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/shasderias/cryptopals/cbc"
)

const (
	c10InputFile = "./set2/c10.txt"
	c10Key       = "YELLOW SUBMARINE"
	c10IV        = "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"
)

func runC10() error {
	b64CipherText, err := ioutil.ReadFile(c10InputFile)
	if err != nil {
		return err
	}
	cipherText, err := base64.StdEncoding.DecodeString(string(b64CipherText))
	if err != nil {
		return err
	}

	cip, err := aes.NewCipher([]byte(c10Key))
	if err != nil {
		return err
	}

	plainText := cbc.Decrypt(cip, []byte(c10IV), cipherText)

	fmt.Println(string(plainText))
	return nil
}

package main

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

const (
	c7InputFile = "./set1/c7.txt"
	c7Key       = "YELLOW SUBMARINE"
)

func runC7() error {
	b64CipherText, err := ioutil.ReadFile(c7InputFile)
	if err != nil {
		return err
	}
	cipherText, err := base64.StdEncoding.DecodeString(string(b64CipherText))
	if err != nil {
		return err
	}

	cipher, err := aes.NewCipher([]byte(c7Key))
	if err != nil {
		return err
	}

	plainText := make([]byte, len(cipherText))

	cipher.Decrypt(plainText, cipherText)

	keySize := len(c7Key)
	blocks := len(cipherText) / keySize
	lastBlockSize := len(cipherText) % keySize

	for i := 0; i < blocks; i++ {
		s := i * keySize
		e := s + keySize
		cipher.Decrypt(plainText[s:e], cipherText[s:e])
	}
	if lastBlockSize != 0 {
		s := blocks * keySize
		e := s + lastBlockSize
		cipher.Decrypt(plainText[s:e], cipherText[s:e])
	}

	fmt.Println(string(plainText))
	return nil
}

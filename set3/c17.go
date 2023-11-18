package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"math/rand"

	"github.com/shasderias/cryptopals/blocks"
	"github.com/shasderias/cryptopals/cbc"
)

const (
	c17Strings = `MDAwMDAwTm93IHRoYXQgdGhlIHBhcnR5IGlzIGp1bXBpbmc=
MDAwMDAxV2l0aCB0aGUgYmFzcyBraWNrZWQgaW4gYW5kIHRoZSBWZWdhJ3MgYXJlIHB1bXBpbic=
MDAwMDAyUXVpY2sgdG8gdGhlIHBvaW50LCB0byB0aGUgcG9pbnQsIG5vIGZha2luZw==
MDAwMDAzQ29va2luZyBNQydzIGxpa2UgYSBwb3VuZCBvZiBiYWNvbg==
MDAwMDA0QnVybmluZyAnZW0sIGlmIHlvdSBhaW4ndCBxdWljayBhbmQgbmltYmxl
MDAwMDA1SSBnbyBjcmF6eSB3aGVuIEkgaGVhciBhIGN5bWJhbA==
MDAwMDA2QW5kIGEgaGlnaCBoYXQgd2l0aCBhIHNvdXBlZCB1cCB0ZW1wbw==
MDAwMDA3SSdtIG9uIGEgcm9sbCwgaXQncyB0aW1lIHRvIGdvIHNvbG8=
MDAwMDA4b2xsaW4nIGluIG15IGZpdmUgcG9pbnQgb2g=
MDAwMDA5aXRoIG15IHJhZy10b3AgZG93biBzbyBteSBoYWlyIGNhbiBibG93`
)

var (
	c17IV         []byte
	c17Key        []byte
	c17Cipher     cipher.Block
	c17PlainTexts [][]byte
)

func init() {
	c17IV = make([]byte, aes.BlockSize)
	if _, err := rand.Read(c17IV); err != nil {
		panic(err)
	}
	c17Key = make([]byte, aes.BlockSize)
	if _, err := rand.Read(c17Key); err != nil {
		panic(err)
	}

	var err error
	c17Cipher, err = aes.NewCipher(c17Key)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBufferString(c17Strings)
	s := bufio.NewScanner(buf)
	for s.Scan() {
		decoded, err := base64.StdEncoding.DecodeString(s.Text())
		if err != nil {
			panic(err)
		}
		c17PlainTexts = append(c17PlainTexts, decoded)
	}
	if s.Err() != nil {
		panic(err)
	}
}

func runC17() error {
	iv, ct := c17GetCookie()

	ctBlocks := blocks.NewSlice(ct, aes.BlockSize)

	for i := 0; i < 256; i++ {
		ct[len(ct)-1] = byte(i)
		if c17PaddingValid(ct) {
			break
		}
	}

	return nil
}

func c17GetCookie() (iv, cipherText []byte) {
	pt := c17PlainTexts[rand.Intn(len(c17PlainTexts))]
	ct := cbc.Encrypt(c17Cipher, c17IV, pt)
	return c17IV, ct
}

func c17PaddingValid(cipherText []byte) bool {
	_, err := cbc.DecryptAndCheckPadding(c17Cipher, c17IV, cipherText)
	return err == nil
}

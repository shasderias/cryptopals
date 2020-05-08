package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"math/rand"
	"net/url"

	"github.com/shasderias/cryptopals/blocks"
	"github.com/shasderias/cryptopals/cbc"
)

const (
	c16BS           = 16
	c16StringPrefix = "comment1=cooking%20MCs;userdata="
	c16StringSuffix = ";comment2=%20like%20a%20pound%20of%20bacon"
)

var (
	c16IV     []byte
	c16Key    []byte
	c16Cipher cipher.Block
)

func init() {
	var err error
	c16IV = make([]byte, c16BS)
	if _, err = rand.Read(c16IV); err != nil {
		panic(err)
	}
	c16Key = make([]byte, c16BS)
	if _, err = rand.Read(c16Key); err != nil {
		panic(err)
	}

	c16Cipher, err = aes.NewCipher(c16Key)
	if err != nil {
		panic(err)
	}
}

func runC16() error {
	if c16DecryptAndCheckForAdmin(c16Encrypt([]byte(";admin=true;"))) {
		return errors.New("admin bit set, invalid escaping")
	} else {
		fmt.Println("not admin")
	}

	// c16Learn()

	const bs = c16BS

	colonBitFlip := byte(';') ^ byte('A')
	equalBitFlip := byte('=') ^ byte('A')

	payloadBlock := []byte("AadminAtrueA1234")

	payload := []byte{}

	payload = append(payload, fixedByteSlice(bs, 'A')...)
	payload = append(payload, payloadBlock...)

	cipherText := c16Encrypt(payload)
	fmt.Println("encrypt/decrypt test")
	c16DecryptAndCheckForAdmin(cipherText)
	fmt.Println("bit flipping")
	cipherText[bs*2+0] = cipherText[bs*2+0] ^ colonBitFlip
	cipherText[bs*2+6] = cipherText[bs*2+6] ^ equalBitFlip
	cipherText[bs*2+11] = cipherText[bs*2+11] ^ colonBitFlip

	fmt.Println("admin found?", c16DecryptAndCheckForAdmin(cipherText))

	return nil
}

func c16Learn() {
	const bs = 16
	key := make([]byte, bs)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	iv := make([]byte, bs)
	if _, err := rand.Read(iv); err != nil {
		panic(err)
	}
	cipher, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	plainText := fixedByteSlice(16*4, 0)

	cipherText := cbc.Encrypt(cipher, iv, plainText)
	ctBlks := blocks.NewSlice(cipherText, bs)
	fmt.Printf("%08b\n%08b\n", ctBlks.N(0), ctBlks.N(1))

	cipherText[0] = cipherText[0] ^ 0b11111111
	fmt.Printf("%08b\n", cipherText[0])
	decPt := cbc.Decrypt(cipher, iv, cipherText)
	decPtBlks := blocks.NewSlice(decPt, bs)
	fmt.Printf("%08b\n%08b\n%08b\n", decPtBlks.N(0), decPtBlks.N(1), decPtBlks.N(2))
}

func c16Encrypt(b []byte) []byte {
	escaped := []byte(url.QueryEscape(string(b)))

	buf := append([]byte{}, c16StringPrefix...)
	buf = append(buf, escaped...)
	buf = append(buf, c16StringSuffix...)

	cipherText := cbc.Encrypt(c16Cipher, c16IV, buf)

	return cipherText
}

func c16DecryptAndCheckForAdmin(cipherText []byte) bool {
	plainText := cbc.Decrypt(c16Cipher, c16IV, cipherText)

	return bytes.Contains(plainText, []byte(";admin=true;"))
}

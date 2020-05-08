package main

import (
	"crypto/aes"
	"fmt"
	"math/rand"

	"github.com/shasderias/cryptopals/cbc"
	"github.com/shasderias/cryptopals/ecb"
)

const (
	ECB = iota
	CBC
)

const (
	c11RunCount  = 100
	c11PlainText = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
)

func runC11() error {
	fmt.Printf("encrypting %d blocks and guessing block cipher mode\n", c11RunCount)

	correctGuesses := 0
	ecbEncoded := 0
	cbcEncoded := 0

	for i := 0; i < c11RunCount; i++ {
		res := c11Oracle([]byte(c11PlainText))
		switch res.mode {
		case ECB:
			ecbEncoded++
		case CBC:
			cbcEncoded++
		}

		cnt := ecb.MostCommonBlockCount(16, res.cipherText)
		var guess bool
		if cnt > 1 {
			guess = res.Guess(ECB)
		} else {
			guess = res.Guess(CBC)
		}
		if guess {
			correctGuesses++
		}
	}
	fmt.Println("ECB Encoded Count:", ecbEncoded)
	fmt.Println("CBC Encoded Count:", cbcEncoded)
	fmt.Println("Correct Guesses:", correctGuesses)
	return nil
}

type c11OracleResult struct {
	mode       int
	cipherText []byte
}

func (r c11OracleResult) Guess(mode int) bool {
	return r.mode == mode
}

func c11Oracle(plainText []byte) c11OracleResult {
	key := randAESKey()

	aLen := rand.Intn(5) + 5
	bLen := rand.Intn(5) + 5

	aRandBytes := make([]byte, aLen)
	bRandBytes := make([]byte, bLen)
	if _, err := rand.Read(aRandBytes); err != nil {
		panic(err)
	}
	if _, err := rand.Read(bRandBytes); err != nil {
		panic(err)
	}

	paddedPlainText := append([]byte{}, aRandBytes...)
	paddedPlainText = append(paddedPlainText, plainText...)
	paddedPlainText = append(paddedPlainText, bRandBytes...)

	cipher, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	switch rand.Intn(2) {
	case ECB: // ECB
		return c11OracleResult{
			mode:       ECB,
			cipherText: ecb.Encrypt(cipher, paddedPlainText),
		}
	case CBC: // CBC
		return c11OracleResult{
			mode:       CBC,
			cipherText: cbc.Encrypt(cipher, randAESKey(), paddedPlainText),
		}
	}
	panic("unreachable code")
}

func randAESKey() []byte {
	key := make([]byte, 16)

	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	return key
}

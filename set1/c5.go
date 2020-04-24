package main

import (
	"encoding/hex"
	"fmt"

	"github.com/shasderias/cryptopals/xor"
)

const (
	c5Input          = "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
	c5Key            = "ICE"
	c5ExpectedOutput = "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"
)

func runC5() error {
	cipherText := xor.RotatingKey([]byte(c5Key), []byte(c5Input))
	hexCipherText := hex.EncodeToString(cipherText)
	if hexCipherText != c5ExpectedOutput {
		fmt.Printf("got %s; want %s", hexCipherText, c5ExpectedOutput)
		return nil
	}
	fmt.Println(hexCipherText)
	return nil
}

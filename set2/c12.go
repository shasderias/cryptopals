package main

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"

	"github.com/google/go-cmp/cmp"

	"github.com/shasderias/cryptopals/ecb"
)

const (
	c12UnknownStringB64 = `Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkg
aGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBq
dXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUg
YnkK`
	c12MaxBlockSize = 32
	c12PadChar      = 'A'
)

var (
	c12UnknownString []byte
	c12Key           []byte
)

func init() {
	var err error
	if c12UnknownString, err = base64.StdEncoding.DecodeString(c12UnknownStringB64); err != nil {
		panic(err)
	}
	c12Key = randAESKey()
}

func runC12() error {
	bs := findBlockSize(c12Oracle)
	fmt.Println("Block Size:", bs)

	mcbc := ecb.MostCommonBlockCount(bs, c12Oracle(fixedByteSlice(128, 'A')))
	fmt.Println("Most Common Block Count:", mcbc)

	blockCount := len(c12UnknownString) / bs

	plainText := []byte{}

	padCharTable := makePadCharTable(bs)

	fmt.Println("starting decrypt of", blockCount, "blocks")

	for i := 0; i < blockCount; i++ {
		for j := bs; j > 0; j-- {
			var dict ecbDecryptDict

			switch {
			case i == 0:
				dict = buildFirstBlockECBDecryptDict(c12Oracle, bs, plainText)
			default:
				dict = buildSubBlockECBDecryptDict(c12Oracle, bs, plainText)
			}

			cipherText := c12Oracle(padCharTable[j-1])
			cipherText = cipherText[bs*i : bs*(i+1)]

			decryptedChar, ok := dict[string(cipherText)]
			if !ok {
				panic("crypted string not in dict")
			}
			plainText = append(plainText, decryptedChar)
		}
	}

	fmt.Println(string(plainText))

	return nil
}

type oracle func([]byte) []byte

func c12Oracle(plainText []byte) []byte {
	paddedPlainText := append([]byte{}, plainText...)
	paddedPlainText = append(paddedPlainText, c12UnknownString...)

	cipher, err := aes.NewCipher(c12Key)
	if err != nil {
		panic(err)
	}

	return ecb.Encrypt(cipher, paddedPlainText)
}

func findBlockSize(o oracle) int {
	preCipherText := o(fixedByteSlice(1, 'A'))

	for i := 2; i < c12MaxBlockSize; i++ {
		cipherText := o(fixedByteSlice(i, 'A'))
		if cmp.Equal(preCipherText[:i-1], cipherText[:i-1]) {
			return i - 1
		}
		preCipherText = cipherText
	}
	return -1
}

func fixedByteSlice(l int, b byte) []byte {
	buf := make([]byte, l)
	for i := 0; i < l; i++ {
		buf[i] = b
	}
	return buf
}

func makePadCharTable(keySize int) [][]byte {
	tbl := make([][]byte, keySize)
	for i := 0; i < keySize; i++ {
		tbl[i] = fixedByteSlice(i, c12PadChar)
	}
	return tbl
}

type ecbDecryptDict map[string]byte

func buildFirstBlockECBDecryptDict(o oracle, keySize int, plainText []byte) ecbDecryptDict {
	padChars := fixedByteSlice((keySize-1)-len(plainText), c12PadChar)
	ptFragLen := (keySize - 1) - len(padChars)
	ptFrag := []byte{}
	ptFrag = append(ptFrag, padChars...)
	ptFrag = append(ptFrag, plainText[len(plainText)-ptFragLen:]...)
	return buildECBDecryptDict(o, keySize, ptFrag)
}

func buildSubBlockECBDecryptDict(o oracle, keySize int, plainText []byte) ecbDecryptDict {
	ptFrag := plainText[len(plainText)-(keySize-1):]
	return buildECBDecryptDict(o, keySize, ptFrag)
}

func buildECBDecryptDict(o oracle, keySize int, plainTextFrag []byte) ecbDecryptDict {
	dict := ecbDecryptDict{}

	for i := 0; i < 256; i++ {
		plainText := []byte{}
		plainText = append(plainText, plainTextFrag...)
		plainText = append(plainText, byte(i))
		cipherText := o(plainText)
		cipherText = cipherText[:keySize]
		dict[string(cipherText)] = byte(i)
	}

	return dict
}

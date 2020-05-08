package ecb

import (
	"crypto/cipher"

	"github.com/shasderias/cryptopals/padding"
)

func Encrypt(c cipher.Block, plainText []byte) []byte {
	bs := c.BlockSize()
	plainText = padding.PKCS7(plainText, bs)
	l := len(plainText)
	blocks := l / bs

	cipherText := make([]byte, l)

	for i := 0; i < blocks; i++ {
		s := i * bs
		e := s + bs
		c.Encrypt(cipherText[s:e], plainText[s:e])
	}
	return cipherText
}

func Decrypt(c cipher.Block, cipherText []byte) []byte {
	bs := c.BlockSize()
	l := len(cipherText)
	if l%bs != 0 {
		panic("cipherText must be an integer multiple of bs")
	}
	blocks := l / bs

	plainText := make([]byte, l)

	for i := 0; i < blocks; i++ {
		s := i * bs
		e := s + bs
		c.Decrypt(plainText[s:e], cipherText[s:e])
	}
	plainText = padding.PKCS7Strip(plainText)
	return plainText
}

// MostCommonBlockCount is used to guess whether cipherText is encoded in ECB
// mode by returning the number of times the most common block appears.
func MostCommonBlockCount(keySize int, cipherText []byte) int {
	blocks := len(cipherText) / keySize
	counts := map[string]int{}

	for i := 0; i < blocks; i++ {
		block := string(cipherText[i*keySize : (i+1)*keySize])
		counts[block]++
	}

	max := 0
	for _, count := range counts {
		if count > max {
			max = count
		}
	}
	return max
}

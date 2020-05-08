package cbc

import (
	"crypto/cipher"

	"github.com/shasderias/cryptopals/blocks"
	"github.com/shasderias/cryptopals/padding"
	"github.com/shasderias/cryptopals/xor"
)

func Decrypt(cipher cipher.Block, iv []byte, cipherText []byte) []byte {
	bs := cipher.BlockSize()

	if len(iv) != cipher.BlockSize() {
		panic("iv must be the same size as cipher's block size")
	}
	if len(cipherText)%bs != 0 {
		panic("length of cipherText must be a multiple of block size")
	}

	plainText := make([]byte, len(cipherText))

	ptBlocks := blocks.NewSlice(plainText, bs)
	ctBlocks := blocks.NewSlice(cipherText, bs)
	prevBlock := iv

	for i := 0; i < ctBlocks.Len(); i++ {
		cipher.Decrypt(ptBlocks.N(i), ctBlocks.N(i))
		xor.Fixed(ptBlocks.N(i), ptBlocks.N(i), prevBlock)
		prevBlock = ctBlocks.N(i)
	}

	plainText = padding.PKCS7Strip(plainText)

	return plainText
}

func DecryptWithoutStrippingPadding(cipher cipher.Block, iv []byte, cipherText []byte) []byte {
	bs := cipher.BlockSize()

	if len(iv) != cipher.BlockSize() {
		panic("iv must be the same size as cipher's block size")
	}
	if len(cipherText)%bs != 0 {
		panic("length of cipherText must be a multiple of block size")
	}

	plainText := make([]byte, len(cipherText))

	ptBlocks := blocks.NewSlice(plainText, bs)
	ctBlocks := blocks.NewSlice(cipherText, bs)
	prevBlock := iv

	for i := 0; i < ctBlocks.Len(); i++ {
		cipher.Decrypt(ptBlocks.N(i), ctBlocks.N(i))
		xor.Fixed(ptBlocks.N(i), ptBlocks.N(i), prevBlock)
		prevBlock = ctBlocks.N(i)
	}

	return plainText
}

func Encrypt(cipher cipher.Block, iv []byte, plainText []byte) []byte {
	bs := cipher.BlockSize()

	if len(iv) != bs {
		panic("iv must be the same size as block size")
	}
	plainText = padding.PKCS7(plainText, bs)

	blocks := len(plainText) / bs

	cipherText := make([]byte, len(plainText))

	xor.Fixed(cipherText[:bs], iv, plainText[:bs])
	cipher.Encrypt(cipherText[:bs], cipherText[:bs])

	for i := 1; i < blocks; i++ {
		i0 := bs * (i - 1)
		i1 := bs * i
		i2 := bs * (i + 1)
		xor.Fixed(cipherText[i1:i2], cipherText[i0:i1], plainText[i1:i2])
		cipher.Encrypt(cipherText[i1:i2], cipherText[i1:i2])
	}

	return cipherText
}

//cipher, err := aes.NewCipher([]byte(c7Key))
//if err != nil {
//return err
//}
//
//plainText := make([]byte, len(cipherText))
//
//cipher.Decrypt(plainText, cipherText)
//
//keySize := len(c7Key)
//blocks := len(cipherText) / keySize
//lastBlockSize := len(cipherText) % keySize
//
//for i := 0; i < blocks; i++ {
//s := i * keySize
//e := s + keySize
//cipher.Decrypt(plainText[s:e], cipherText[s:e])
//}
//if lastBlockSize != 0 {
//s := blocks * keySize
//e := s + lastBlockSize
//cipher.Decrypt(plainText[s:e], cipherText[s:e])
//}

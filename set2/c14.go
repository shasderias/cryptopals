package main

import (
	"crypto/aes"
	"fmt"
	"math/rand"

	"github.com/shasderias/cryptopals/blocks"
	"github.com/shasderias/cryptopals/ecb"
)

var (
	c14Key         []byte
	c14RandomBytes []byte
)

func init() {
	c14Key = make([]byte, 16)
	if _, err := rand.Read(c14Key); err != nil {
		panic(err)
	}
	c14RandomBytes = make([]byte, rand.Intn(251)+5)
	if _, err := rand.Read(c14RandomBytes); err != nil {
		panic(err)
	}
}

func runC14() error {
	const bs = 16
	boundaryPrefix, targetBlock := findBoundaryPrefix(c14Oracle, bs)
	fmt.Println("Guessed Bytes Needed to Reach Block Boundary", boundaryPrefix)
	fmt.Println("Actual Bytes Needed to Reach Block Boundary", bs-len(c14RandomBytes)%bs)
	fmt.Println("Target Block", targetBlock)

	decrypter := ecbDecrypter{
		o:              c14Oracle,
		ks:             bs,
		boundaryPrefix: boundaryPrefix,
		padChar:        'A',
		targetBlock:    targetBlock,
	}

	fmt.Println(string(decrypter.decrypt()))
	return nil
}

func findBoundaryPrefix(o oracle, bs int) (prefixLen, blockN int) {
	ct := o(fixedByteSlice(bs*16, 'A'))
	blks := blocks.NewSlice(ct, bs)
	_, firstPos, _ := blks.MostCommonBlock()

	firstPos--

	prevBlock := string(blks.N(firstPos))
	for i := 0; i < bs; i++ {
		ct := o(fixedByteSlice(i, 'A'))
		blks := blocks.NewSlice(ct, bs)
		if prevBlock == string(blks.N(firstPos)) {
			return i, firstPos + 1
		}
	}
	return -1, -1
}

type ecbDecrypter struct {
	o              oracle
	ks             int
	boundaryPrefix int
	padChar        byte
	targetBlock    int
}

func (d *ecbDecrypter) buildFirstBlockDict(plainText []byte) ecbDecryptDict {
	padCharCount := (d.ks - 1) - len(plainText)
	padChars := fixedByteSlice(d.boundaryPrefix+padCharCount, d.padChar)
	ptFragLen := (d.ks - 1) - padCharCount
	ptFrag := append([]byte{}, padChars...)
	ptFrag = append(ptFrag, plainText[len(plainText)-ptFragLen:]...)
	return d.buildDict(ptFrag)
}

func (d *ecbDecrypter) buildSubsequentBlockDict(plainText []byte) ecbDecryptDict {
	boundaryPadChars := fixedByteSlice(d.boundaryPrefix, d.padChar)
	ptFrag := append(boundaryPadChars, plainText[len(plainText)-(d.ks-1):]...)
	return d.buildDict(ptFrag)
}

func (d *ecbDecrypter) buildDict(plainText []byte) ecbDecryptDict {
	dict := ecbDecryptDict{}
	for i := 0; i < 256; i++ {
		payload := append([]byte{}, plainText...)
		payload = append(payload, byte(i))
		cipherText := d.o(payload)
		cipherText = cipherText[d.targetBlock*d.ks : (d.targetBlock+1)*d.ks]
		dict[string(cipherText)] = byte(i)
	}
	return dict
}

func (d *ecbDecrypter) makePadCharTable() [][]byte {
	tbl := make([][]byte, d.ks)
	for i := 0; i < d.ks; i++ {
		tbl[i] = fixedByteSlice(i+d.boundaryPrefix, d.padChar)
	}
	return tbl
}

func (d *ecbDecrypter) decrypt() []byte {
	padCharTable := d.makePadCharTable()

	plainText := []byte{}

	for i := d.targetBlock; true; i++ {
		for j := d.ks; j > 0; j-- {
			var dict ecbDecryptDict
			if i == d.targetBlock {
				dict = d.buildFirstBlockDict(plainText)
			} else {
				dict = d.buildSubsequentBlockDict(plainText)
			}

			ct := d.o(padCharTable[j-1])
			ct = ct[d.ks*i : d.ks*(i+1)]

			decryptedChar, ok := dict[string(ct)]
			if !ok {
				goto DecryptComplete
			}
			plainText = append(plainText, decryptedChar)
		}
	}
DecryptComplete:
	return plainText
}

func c14Oracle(plainText []byte) []byte {
	buf := append([]byte{}, c14RandomBytes...)
	buf = append(buf, plainText...)
	buf = append(buf, c12UnknownString...)

	cipher, err := aes.NewCipher(c14Key)
	if err != nil {
		panic(err)
	}

	return ecb.Encrypt(cipher, buf)
}

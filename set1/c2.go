package main

import (
	"encoding/hex"
	"fmt"
)

const (
	c2Input1         = "1c0111001f010100061a024b53535009181c"
	c2Input2         = "686974207468652062756c6c277320657965"
	c2ExpectedOutput = "746865206b696420646f6e277420706c6179"
)

func runC2() error {
	in1 := mustDecodeHexString(c2Input1)
	in2 := mustDecodeHexString(c2Input2)

	out := fixedXOR(in1, in2)
	outHex := hex.EncodeToString(out)
	if outHex != c2ExpectedOutput {
		fmt.Printf("got %s; want %s;", outHex, c2ExpectedOutput)
		return nil
	}
	fmt.Println(outHex)
	return nil
}

func fixedXOR(a, b []byte) []byte {
	if len(a) != len(b) {
		panic("a and b must be equal in length")
	}

	buf := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		buf[i] = a[i] ^ b[i]
	}
	return buf
}

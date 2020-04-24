package main

import (
	"encoding/hex"
	"fmt"
)

type result struct {
	cipherText string
	plainText  string
	key        byte
	score      float64
}

func printTopNResults(n int, r []result) {
	if n > len(r) {
		n = len(r)
	}

	for i := 0; i < n; i++ {
		fmt.Printf("%0.2f - %d - %s - %v\n", r[i].score, r[i].key, r[i].cipherText, r[i].plainText)
	}
}

type byScore []result

func (s byScore) Len() int           { return len(s) }
func (s byScore) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byScore) Less(i, j int) bool { return s[i].score < s[j].score }

func singleCharXOR(src []byte, b byte) []byte {
	buf := make([]byte, len(src))
	copy(buf, src)
	for i := 0; i < len(buf); i++ {
		buf[i] = buf[i] ^ b
	}
	return buf
}

func mustDecodeHexString(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

type challengeFunc func() error

func main() {
	challenges := []challengeFunc{
		runC1,
		runC2,
		runC3,
		runC4,
		runC5,
		runC6,
		runC7,
		runC8,
	}

	for n, challengeFunc := range challenges {
		fmt.Println("Challenge", n+1)
		if err := challengeFunc(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("")
	}
}

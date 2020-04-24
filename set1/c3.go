package main

import (
	"sort"

	"github.com/shasderias/cryptopals/ptscore"
)

const (
	c3CipherText     = "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	c3ResultsToPrint = 3
)

func runC3() error {
	scorer := ptscore.NewScorer(ptscore.EngWithSpaceCharFreq)
	results := []result{}
	cipherText := mustDecodeHexString(c3CipherText)

	for i := 0; i < 256; i++ {
		plainText := string(singleCharXOR(cipherText, byte(i)))
		score := scorer.ChiSquare(plainText)
		results = append(results, result{
			cipherText: c3CipherText,
			plainText:  plainText,
			key:        byte(i),
			score:      score,
		})
	}

	sort.Sort(byScore(results))
	printTopNResults(c3ResultsToPrint, results)
	return nil
}

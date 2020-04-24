package main

import (
	"bufio"
	"os"
	"sort"

	"github.com/shasderias/cryptopals/ptscore"
)

const (
	c4InputFile      = "./set1/c4.txt"
	c4ResultsToPrint = 3
)

func runC4() error {
	f, err := os.Open(c4InputFile)
	if err != nil {
		panic(err)
	}
	s := bufio.NewScanner(f)

	scorer := ptscore.NewScorer(ptscore.EngWithSpaceCharFreq)
	results := []result{}

	for s.Scan() {
		cipherText := mustDecodeHexString(s.Text())
		for i := 0; i < 256; i++ {
			plainText := string(singleCharXOR(cipherText, byte(i)))
			score := scorer.ChiSquare(plainText)
			results = append(results, result{
				cipherText: s.Text(),
				plainText:  plainText,
				key:        byte(i),
				score:      score,
			})
		}
	}
	if err := s.Err(); err != nil {
		panic(err)
	}

	sort.Sort(byScore(results))
	printTopNResults(c4ResultsToPrint, results)

	return nil
}

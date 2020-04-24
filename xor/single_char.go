package xor

import (
	"sort"
)

type SingleCharResult struct {
	CipherText string
	PlainText  string
	Key        byte
	Score      float64
}

type SingleCharResults []SingleCharResult

func SolveSingleChar(scorer Scorer, cipherText []byte) SingleCharResults {
	results := SingleCharResults{}

	for i := 0; i < 256; i++ {
		plainText := SingleChar(cipherText, byte(i))
		score := scorer.Score(string(plainText))
		results = append(results, SingleCharResult{
			CipherText: string(cipherText),
			PlainText:  string(plainText),
			Key:        byte(i),
			Score:      score,
		})
	}

	sortSingleCharResultsByScore(results)

	return results
}

func sortSingleCharResultsByScore(r SingleCharResults) {
	sort.Slice(r, func(i, j int) bool {
		return r[i].Score < r[j].Score
	})
}

func SingleChar(src []byte, b byte) []byte {
	buf := make([]byte, len(src))
	copy(buf, src)
	for i := 0; i < len(buf); i++ {
		buf[i] = buf[i] ^ b
	}
	return buf
}

type Scorer interface {
	Score(s string) float64
}

package ptscore

import (
	"math"
	"strings"
	"unicode/utf8"
)

const (
	defaultInvalidUTF8Penalty      = 100
	defaultUnrecognizedCharRatio   = 0.7
	defaultUnrecognizedCharPenalty = 100
)

type Scorer struct {
	expectedCharFreq        CharFreq
	invalidUTF8Penalty      float64
	unrecognizedCharRatio   float64
	unrecognizedCharPenalty float64
}

func NewScorer(expectedCF CharFreq) Scorer {
	return Scorer{
		expectedCharFreq:        expectedCF,
		invalidUTF8Penalty:      defaultInvalidUTF8Penalty,
		unrecognizedCharRatio:   defaultUnrecognizedCharRatio,
		unrecognizedCharPenalty: defaultUnrecognizedCharPenalty,
	}
}

func (sc Scorer) ChiSquare(s string) float64 {
	sStats := calcStringStats(sc.expectedCharFreq, s)
	sLen := float64(len(s))

	chiSqStat := 0.0

	for char, expectedFreq := range sc.expectedCharFreq {
		actualObs := float64(sStats.CharObs[char])
		expectedObs := expectedFreq * sLen
		charChiSqStat := math.Pow(actualObs-expectedObs, 2) / expectedObs
		chiSqStat += charChiSqStat
	}

	if float64(sStats.CharCount) < sLen*sc.unrecognizedCharRatio {
		chiSqStat += sc.unrecognizedCharPenalty
	}

	if !utf8.ValidString(s) {
		chiSqStat += sc.invalidUTF8Penalty
	}

	return chiSqStat
}

type StringStat struct {
	CharObs
	CharCount int
}

type CharObs map[rune]int

func calcStringStats(expectedCharFreq CharFreq, s string) StringStat {
	s = strings.ToLower(s)

	obs := make(CharObs)
	count := 0

	for _, char := range s {
		if _, ok := expectedCharFreq[char]; !ok {
			continue
		}
		obs[char] += 1
		count++
	}

	return StringStat{
		CharObs:   obs,
		CharCount: count,
	}
}

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/shasderias/cryptopals/ptscore"
	"github.com/shasderias/cryptopals/strdist"
	"github.com/shasderias/cryptopals/xor"
)

const (
	c6InputFile           = "./set1/c6.txt"
	c6MinKeySize          = 2
	c6MaxKeySize          = 40
	c6KeySizeGuessesToUse = 5
)

type keySizeGuess struct {
	keySize            int
	normalizedEditDist float64
}

type c6Result struct {
	plainText []byte
	keySize   int
	key       []byte
	score     float64
}

type chiSqShim struct {
	ptscore.Scorer
}

func (s chiSqShim) Score(str string) float64 {
	return s.Scorer.ChiSquare(str)
}

func runC6() error {
	b64CipherText, err := ioutil.ReadFile(c6InputFile)
	if err != nil {
		return err
	}
	cipherText, err := base64.StdEncoding.DecodeString(string(b64CipherText))
	if err != nil {
		return err
	}

	guesses := []keySizeGuess{}

	for i := c6MinKeySize; i < c6MaxKeySize+1; i++ {
		guesses = append(guesses, keySizeGuess{
			keySize:            i,
			normalizedEditDist: strdist.AverageHamming(i, 8, cipherText),
		})
	}

	sort.Slice(guesses, func(i, j int) bool {
		return guesses[i].normalizedEditDist < guesses[j].normalizedEditDist
	})

	fmt.Println("attempting to decode with following key sizes")
	for i := 0; i < c6KeySizeGuessesToUse; i++ {
		fmt.Printf("#%d - %2d - %0.3f\n", i+1, guesses[i].keySize, guesses[i].normalizedEditDist)
	}

	scorer := chiSqShim{ptscore.NewScorer(ptscore.EngWithSpaceCharFreq)}
	results := []c6Result{}

	for i := 0; i < c6KeySizeGuessesToUse; i++ {
		keySize := guesses[i].keySize

		blocks := splitAndTranspose(keySize, cipherText)

		key := make([]byte, keySize)

		for j := 0; j < len(blocks); j++ {
			results := xor.SolveSingleChar(scorer, blocks[j])
			key[j] = results[0].Key
		}

		plainText := xor.RotatingKey(key, cipherText)
		results = append(results, c6Result{
			plainText: plainText,
			keySize:   keySize,
			key:       key,
			score:     scorer.Score(string(plainText)),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score < results[j].score
	})

	fmt.Println("Key:", string(results[0].key))
	fmt.Println("KeySize:", results[0].keySize)
	fmt.Println(string(results[0].plainText))
	return nil
}

func normalizedEditDistanceNBlocks(keySize, n int, cipherText []byte) float64 {
	totalEditDistance := 0
	for i := 1; i < n; i++ {
		preBlock := cipherText[keySize*(n-1) : keySize*n]
		curBlock := cipherText[keySize*n : keySize*(n+1)]
		totalEditDistance += strdist.Hamming(preBlock, curBlock)
	}
	return float64(totalEditDistance) / float64(keySize) / float64(n-1)
}

/*
0123456789 for KEYSIZE of 3

012 345 678 9

0369
147
258
*/
func splitAndTranspose(keySize int, src []byte) [][]byte {
	blocks := make([][]byte, keySize)
	for i := 0; i < keySize; i++ {
		blocks[i] = make([]byte, 0, len(src)/keySize)
		for j := 0; j*keySize+i < len(src); j++ {
			blocks[i] = append(blocks[i], src[j*keySize+i])
		}
	}
	return blocks
}

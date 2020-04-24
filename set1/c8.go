package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
)

const (
	c8InputFilename = "./set1/c8.txt"
	c8ResultsToShow = 5
)

type c8Result struct {
	cipherText      string
	editDistScore   float64
	repetitionScore int
}

func runC8() error {
	f, err := os.Open(c8InputFilename)
	if err != nil {
		return err
	}
	s := bufio.NewScanner(f)

	results := []c8Result{}

	for s.Scan() {
		cipherTextHex := s.Text()
		cipherText, err := hex.DecodeString(cipherTextHex)
		if err != nil {
			return err
		}

		results = append(results, c8Result{
			cipherText:      cipherTextHex,
			editDistScore:   normalizedEditDistanceNBlocks(16, 8, cipherText),
			repetitionScore: maxRepeatedBlocks(cipherText),
		})
	}

	if s.Err() != nil {
		return err
	}

	fmt.Println("Edit Distance Candidates")
	sort.Slice(results, func(i, j int) bool {
		return results[i].editDistScore < results[j].editDistScore
	})

	for i := 0; i < c8ResultsToShow; i++ {
		fmt.Println(results[i].editDistScore, "-", results[i].cipherText)
	}

	fmt.Println("\nRepeated Blocks Candidates")
	sort.Slice(results, func(i, j int) bool {
		return results[i].repetitionScore > results[j].repetitionScore
	})

	for i := 0; i < c8ResultsToShow; i++ {
		fmt.Println(results[i].repetitionScore, "-", results[i].cipherText)
	}
	return nil
}

func maxRepeatedBlocks(b []byte) int {
	const keySize = 16
	blocks := len(b) / keySize
	counts := map[[keySize]byte]int{}

	for i := 0; i < blocks; i++ {
		var block [keySize]byte
		copy(block[:], b[i*keySize:(i+1)*keySize])
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

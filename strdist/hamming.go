package strdist

import "math/bits"

// HammingStr returns the hamming distance between s1 and s2, in number of bits
func HammingStr(s1, s2 string) int {
	if len(s1) != len(s2) {
		panic("s1 and s2 must of equal lengths")
	}

	return Hamming([]byte(s1), []byte(s2))
}

// Hamming returns the hamming distance between b1 and b2, in number of bits
func Hamming(b1, b2 []byte) int {
	if len(b1) != len(b2) {
		panic("s1 and s2 must of equal lengths")
	}

	dist := 0

	for i := 0; i < len(b1); i++ {
		d := b1[i] ^ b2[i]
		dist += bits.OnesCount8(d)
	}

	return dist
}

// AverageHamming divides b into keySized blocks and calculates the hamming
// distance of the first n successive blocks, normalized by keySize
func AverageHamming(keySize, n int, b []byte) float64 {
	totalEditDistance := 0
	for i := 1; i < n; i++ {
		preBlock := b[keySize*(n-1) : keySize*n]
		curBlock := b[keySize*n : keySize*(n+1)]
		totalEditDistance += Hamming(preBlock, curBlock)
	}
	return float64(totalEditDistance) / float64(keySize) / float64(n-1)
}

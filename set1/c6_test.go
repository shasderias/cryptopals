package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSplitAndTranspose(t *testing.T) {
	testCases := []struct {
		keySize int
		in      string
		out     []string
	}{
		{
			keySize: 3,
			in:      "012345678",
			out:     []string{"036", "147", "258"},
		},
		{
			keySize: 3,
			in:      "0123456789",
			out:     []string{"0369", "147", "258"},
		},
		{
			keySize: 4,
			in:      "0123456789",
			out:     []string{"048", "159", "26", "37"},
		},
	}

	for _, tc := range testCases {
		out := splitAndTranspose(tc.keySize, []byte(tc.in))
		expectedOut := make([][]byte, tc.keySize)
		for i := 0; i < tc.keySize; i++ {
			expectedOut[i] = []byte(tc.out[i])
		}
		if diff := cmp.Diff(out, expectedOut); diff != "" {
			t.Fatal(diff)
		}
	}
}

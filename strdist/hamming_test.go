package strdist_test

import (
	"testing"

	"github.com/shasderias/cryptopals/strdist"
)

func TestHamming(t *testing.T) {
	testCases := []struct {
		s1           string
		s2           string
		expectedDist int
	}{
		{
			s1:           "this is a test",
			s2:           "wokka wokka!!!",
			expectedDist: 37,
		},
	}

	for _, tc := range testCases {
		dist := strdist.HammingStr(tc.s1, tc.s2)
		if dist != tc.expectedDist {
			t.Fatalf("got %d; want %d", dist, tc.expectedDist)
		}
	}
}

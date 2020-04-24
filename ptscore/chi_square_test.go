package ptscore

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCalcStringStats(t *testing.T) {
	testCases := []struct {
		str             string
		expectedCharObs CharObs
	}{
		{
			str: "ab",
			expectedCharObs: CharObs{
				'a': 1,
				'b': 1,
			},
		},
		{
			str: "aab",
			expectedCharObs: CharObs{
				'a': 2,
				'b': 1,
			},
		},
		{
			str: "abAA",
			expectedCharObs: CharObs{
				'a': 3,
				'b': 1,
			},
		},
		{
			str: "aAbB",
			expectedCharObs: CharObs{
				'a': 2,
				'b': 2,
			},
		},
	}

	for _, tc := range testCases {
		stats := calcStringStats(EngWithSpaceCharFreq, tc.str)
		if diff := cmp.Diff(stats.CharObs, tc.expectedCharObs); diff != "" {
			t.Fatal(diff)
		}
	}
}

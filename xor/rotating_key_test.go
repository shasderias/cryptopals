package xor_test

import (
	"encoding/hex"
	"testing"

	"github.com/shasderias/cryptopals/xor"
)

func TestRotatingKey(t *testing.T) {
	testCases := []struct {
		key        string
		plainText  string
		cipherText string
	}{
		{
			key:        "ICE",
			plainText:  "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal",
			cipherText: "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f",
		},
	}

	for _, tc := range testCases {
		cipherText := xor.RotatingKey([]byte(tc.key), []byte(tc.plainText))
		encCipherText := hex.EncodeToString(cipherText)
		if encCipherText != tc.cipherText {
			t.Fatalf("got %s;\nwant %s", encCipherText, tc.cipherText)
		}
	}
}

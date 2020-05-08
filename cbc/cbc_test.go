package cbc_test

import (
	"crypto/aes"
	"crypto/cipher"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/shasderias/cryptopals/cbc"
)

const (
	key       = "YELLOW SUBMARINE"
	iv        = "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"
	plainText = "the quick brown fox jumps over the lazy dog fish"
)

func TestEncryptCryptoPals(t *testing.T) {
	cip, err := aes.NewCipher([]byte(key))
	if err != nil {
		t.Fatal(err)
	}

	cipherText := cbc.Encrypt(cip, []byte(iv), []byte(plainText))
	decryptedPlainText := cbc.Decrypt(cip, []byte(iv), cipherText)

	if diff := cmp.Diff(string(decryptedPlainText), plainText); diff != "" {
		t.Fatal(diff)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	testCases := []struct {
		plainText string
	}{
		{
			plainText: "approximately two and some blocks long",
		},
	}

	for _, tc := range testCases {
		iv, _, c := aesCipherWithRandKey(t)
		ct := cbc.Encrypt(c, iv, []byte(tc.plainText))
		cbc.Decrypt(c, iv, ct)
		if diff := cmp.Diff(tc.plainText, string(ct)); diff != "" {
			t.Error(diff)
		}
	}
}

func aesCipherWithRandKey(t *testing.T) (iv, key []byte, c cipher.Block) {
	const bs = aes.BlockSize

	t.Helper()

	iv = make([]byte, bs)
	key = make([]byte, bs)

	if _, err := rand.Read(iv); err != nil {
		t.Fatal(err)
	}
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}

	return
}

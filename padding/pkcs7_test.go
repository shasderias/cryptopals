package padding_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/shasderias/cryptopals/padding"
)

func TestPKCS7(t *testing.T) {
	testCases := []struct {
		b     []byte
		bs    int
		wantB []byte
	}{
		{
			b:     []byte("YELLOW SUBMARINE"),
			bs:    20,
			wantB: []byte("YELLOW SUBMARINE\x04\x04\x04\x04"),
		},
		{
			b:     []byte("YELLOW SUBMARINE"),
			bs:    16,
			wantB: []byte("YELLOW SUBMARINE\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10"),
		},
	}

	for _, tc := range testCases {
		padded := padding.PKCS7(tc.b, tc.bs)
		if diff := cmp.Diff(padded, tc.wantB); diff != "" {
			t.Fatal(diff)
		}
		stripped := padding.PKCS7Strip(padded)
		if diff := cmp.Diff(stripped, tc.b); diff != "" {
			t.Fatal(diff)
		}

	}
}

func TestPKCS7ValidateAndStrip(t *testing.T) {
	testCases := []struct {
		b        []byte
		bs       int
		stripped []byte
		valid    bool
	}{
		{
			b:        []byte("YELLOW SUBMARINE\x04\x04\x04\x04"),
			bs:       20,
			stripped: []byte("YELLOW SUBMARINE"),
			valid:    true,
		},
		{
			b:     []byte("ICE ICE BABY\x05\x05\x05\x05"),
			bs:    16,
			valid: false,
		},
	}

	for _, tc := range testCases {
		err, stripped := padding.PKCS7ValidateAndStrip(tc.b, tc.bs)
		if !tc.valid && err != padding.ErrInvalidPKCS7Padding {
			t.Fatalf("got %v; want error", err)
		}
		if tc.valid {
			if diff := cmp.Diff(stripped, tc.stripped); diff != "" {
				t.Fatal(diff)
			}
		}
	}
}

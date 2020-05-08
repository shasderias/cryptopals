package main_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	main "github.com/shasderias/cryptopals/set2"
)

func TestProfileFor(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		wantProfile main.Profile
	}{
		{
			name:  "cryptopals",
			email: "foo@bar.com",
			wantProfile: main.Profile{
				Email: "foo@bar.com",
				UID:   10,
				Role:  "user",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prof := main.ProfileFor(tt.email)
			if diff := cmp.Diff(prof, tt.wantProfile); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestEncodedProfileFor(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		want    string
		wantErr bool
	}{
		{
			name:    "cryptopals",
			email:   "foo@bar.com",
			want:    "email=foo@bar.com&uid=10&role=user",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encProf, err := main.EncodedProfileFor(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodedProfileFor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if encProf != tt.want {
				t.Errorf("got: %v; want: %v", encProf, tt.want)
			}
		})
	}
}

func TestProfileEncryptDecrypt(t *testing.T) {
	wantProfile := main.Profile{
		Email: "abc@example.com",
		UID:   10,
		Role:  "user",
	}

	cipherText, err := main.EncryptedEncodedProfile(wantProfile.Email)
	if err != nil {
		t.Fatal(err)
	}

	prof, err := main.DecryptProfile(cipherText)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(prof, wantProfile); diff != "" {
		t.Fatal(diff)
	}
}

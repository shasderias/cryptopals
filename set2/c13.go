package main

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"

	"github.com/shasderias/cryptopals/ecb"
	"github.com/shasderias/cryptopals/kv"
)

const (
	c13EmailPadChar = 'a'
	c13Domain       = "@abc.com"
)

var (
	c13Key    []byte
	c13Cipher cipher.Block
)

func init() {
	c13Key = randAESKey()
	var err error
	c13Cipher, err = aes.NewCipher(c13Key)
	if err != nil {
		panic(err)
	}
}

func runC13() error {
	const (
		bs = 16
	)

	padCharCount, err := FindBlockStart(EncryptedEncodedProfile, bs)
	if err != nil {
		return err
	}
	adminPayload := fixedString(padCharCount, '0') + "admin" + fixedString(11, 11)
	adminPayloadCiphertext, err := EncryptedEncodedProfile(adminPayload)
	if err != nil {
		return err
	}
	adminPayloadCiphertext = adminPayloadCiphertext[bs : bs*2]

	for i := 1; i < 256; i++ {
		prof, err := CutAndPasteAndDecode(EncryptedEncodedProfile, nLenEmail(i), bs, adminPayloadCiphertext)
		if err != nil {
			continue
		}
		if prof.Role == "admin" {
			fmt.Println(prof)
			return nil
		}
	}

	return nil
}

func nLenEmail(n int) string {
	buf := make([]byte, n, n+len(c13Domain))
	for i := 0; i < n; i++ {
		buf[i] = c13EmailPadChar
	}
	buf = append(buf, []byte(c13Domain)...)
	return string(buf)
}

type Profile struct {
	Email string
	UID   int
	Role  string
}

func ProfileFor(email string) Profile {
	return Profile{
		Email: email,
		UID:   10,
		Role:  "user",
	}
}

func EncodedProfileFor(email string) (string, error) {
	encodedProfile, err := kv.Marshal(ProfileFor(email))
	if err != nil {
		return "", err
	}
	return string(encodedProfile), nil
}

func EncryptedEncodedProfile(email string) ([]byte, error) {
	prof, err := EncodedProfileFor(email)
	if err != nil {
		return nil, err
	}
	return ecb.Encrypt(c13Cipher, []byte(prof)), nil
}

func DecryptProfile(cipherText []byte) (Profile, error) {
	plainText := ecb.Decrypt(c13Cipher, cipherText)
	prof := Profile{}
	if err := kv.Unmarshal(plainText, &prof); err != nil {
		return Profile{}, err
	}
	return prof, nil
}

type strOracle func(string) ([]byte, error)

func fixedString(l int, b byte) string {
	buf := make([]byte, l)
	for i := 0; i < l; i++ {
		buf[i] = b
	}
	return string(buf)
}

func FindBlockStart(o strOracle, bs int) (int, error) {
	cipherText, err := o("")
	if err != nil {
		return 0, err
	}
	preFirstBlock := string(cipherText[:bs])

	for i := 1; i < bs; i++ {
		str := fixedString(i, 0)
		cipherText, err := o(str)
		if err != nil {
			return 0, err
		}
		firstBlock := string(cipherText[:bs])
		if preFirstBlock == firstBlock {
			return i - 1, nil
		}
		preFirstBlock = firstBlock
	}

	return 0, errors.New("cannot find start of block")
}

// email=foo@bar.com&uid=10&role=user
// email

func CutAndPasteAndDecode(o strOracle, email string, bs int, adminCipherText []byte) (Profile, error) {
	cipherText, err := o(email)
	if err != nil {
		return Profile{}, nil
	}
	blocks := len(cipherText) / bs
	payloadCipherText := append(cipherText[:bs*(blocks-1)], adminCipherText...)
	return DecryptProfile(payloadCipherText)
}

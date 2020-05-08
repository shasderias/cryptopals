package padding

import (
	"errors"
)

var (
	ErrInvalidPKCS7Padding = errors.New("invalid PKCS7 padding")
	ErrInvalidBlockSize    = errors.New("invalid block size")
)

func PKCS7(b []byte, bs int) []byte {
	l := len(b)
	bytesToAdd := bs - (l % bs)
	if bytesToAdd == 0 {
		bytesToAdd = bs
	}
	buf := make([]byte, l+bytesToAdd)
	copy(buf, b)

	for i := len(b); i < len(buf); i++ {
		buf[i] = byte(bytesToAdd)
	}
	return buf
}

func PKCS7Strip(b []byte) []byte {
	bytesToStrip := b[len(b)-1]
	return b[:len(b)-int(bytesToStrip)]
}

func PKCS7ValidateAndStrip(b []byte, bs int) (error, []byte) {
	if len(b)%bs != 0 {
		return ErrInvalidBlockSize, nil
	}
	l := len(b)
	bytesToStrip := b[l-1]
	for i := l - int(bytesToStrip); i < len(b); i++ {
		if b[i] != bytesToStrip {
			return ErrInvalidPKCS7Padding, nil
		}
	}
	return nil, PKCS7Strip(b)
}

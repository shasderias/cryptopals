package xor

func RotatingKey(key, plainText []byte) []byte {
	cipherText := make([]byte, len(plainText))

	nextKey := func() func() byte {
		i := 0
		return func() byte {
			if i == len(key) {
				i = 0
			}
			nk := key[i]
			i++
			return nk
		}
	}()

	for i := 0; i < len(plainText); i++ {
		cipherText[i] = plainText[i] ^ nextKey()
	}

	return cipherText
}

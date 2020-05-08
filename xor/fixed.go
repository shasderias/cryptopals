package xor

func Fixed(dst, a, b []byte) {
	if len(a) != len(b) {
		panic("a and b must be equal in length")
	}

	for i := 0; i < len(a); i++ {
		dst[i] = a[i] ^ b[i]
	}
}

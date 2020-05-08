package blocks

type Slice struct {
	bs  int
	buf []byte

	blockCount int
}

func NewSlice(b []byte, blockSize int) *Slice {
	if len(b)%blockSize != 0 {
		panic("blockSize must be an integer multiple of len(b)")
	}

	blockCount := len(b) / blockSize

	return &Slice{
		bs:         blockSize,
		buf:        b,
		blockCount: blockCount,
	}
}

func (s *Slice) N(n int) []byte {
	buf, bs := s.buf, s.bs
	return buf[bs*n : bs*(n+1)]
}

func (s *Slice) Len() int {
	return s.blockCount
}

// MostCommonBlockCount is used to guess whether cipherText is encoded in ECB
// mode by returning the number of times the most common block appears.
func (s *Slice) MostCommonBlock() (block []byte, firstPos, count int) {
	counter := map[string]struct {
		firstPos int
		count    int
	}{}

	for i := 0; i < s.blockCount; i++ {
		block := string(s.N(i))
		c := counter[block]
		c.count++
		if c.firstPos == 0 {
			c.firstPos = i
		}
		counter[block] = c
	}

	max := 0
	mostCommonBlock := counter[""]
	for _, c := range counter {
		if c.count > max {
			max = c.count
			mostCommonBlock = c
		}
	}

	return s.N(mostCommonBlock.firstPos), mostCommonBlock.firstPos, mostCommonBlock.count
}

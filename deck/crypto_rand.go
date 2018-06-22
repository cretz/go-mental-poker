package deck

import (
	"crypto/rand"
	"encoding/binary"
	insecureRand "math/rand"
)

type cryptoRandSource struct{}

func (cryptoRandSource) Int63() int64 {
	var b [8]byte
	rand.Read(b[:])
	// mask off sign bit to ensure positive number
	return int64(binary.LittleEndian.Uint64(b[:]) & (1<<63 - 1))
}

func (cryptoRandSource) Seed(int64) {}

func newCryptoRand() *insecureRand.Rand { return insecureRand.New(cryptoRandSource{}) }

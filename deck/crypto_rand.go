package deck

import (
	"crypto/rand"
	"encoding/binary"
	insecureRand "math/rand"
)

// cryptoRandSource impls math/rand.Source for crypto/rand
type cryptoRandSource struct{}

// Int63 impls math/rand.Source.Int63
func (cryptoRandSource) Int63() int64 {
	var b [8]byte
	rand.Read(b[:])
	// mask off sign bit to ensure positive number
	return int64(binary.LittleEndian.Uint64(b[:]) & (1<<63 - 1))
}

// Seed is a noop impl of math/rand.Source.Seed
func (cryptoRandSource) Seed(int64) {}

// newCryptoRand creates a math/rand.Rand for cryptoRandSource
func newCryptoRand() *insecureRand.Rand { return insecureRand.New(cryptoRandSource{}) }

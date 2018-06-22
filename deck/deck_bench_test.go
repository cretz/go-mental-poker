package deck_test

import (
	"crypto/rand"
	"log"
	"math/big"
	"testing"

	"github.com/cretz/go-mental-poker/deck"
)

// Testing has shown prime size doesn't matter that much
var prime = genPrime(1024)
var playersByBits = map[int][]deck.Player{}

func BenchmarkShuffle2Players32Bits52Cards(b *testing.B) {
	benchmarkShuffle(b, 2, 32, 52)
}
func BenchmarkShuffle2Players64Bits52Cards(b *testing.B) {
	benchmarkShuffle(b, 2, 64, 52)
}
func BenchmarkShuffle3Players32Bits52Cards(b *testing.B) {
	benchmarkShuffle(b, 3, 32, 52)
}
func BenchmarkShuffle3Players64Bits52Cards(b *testing.B) {
	benchmarkShuffle(b, 3, 64, 52)
}
func BenchmarkShuffle6Players32Bits52Cards(b *testing.B) {
	benchmarkShuffle(b, 6, 32, 52)
}
func BenchmarkShuffle6Players64Bits52Cards(b *testing.B) {
	benchmarkShuffle(b, 6, 64, 52)
}
func BenchmarkShuffle2Players32Bits104Cards(b *testing.B) {
	benchmarkShuffle(b, 2, 32, 104)
}
func BenchmarkShuffle2Players64Bits104Cards(b *testing.B) {
	benchmarkShuffle(b, 2, 64, 104)
}
func BenchmarkShuffle3Players32Bits104Cards(b *testing.B) {
	benchmarkShuffle(b, 3, 32, 104)
}

func BenchmarkShuffle3Players64Bits104Cards(b *testing.B) {
	benchmarkShuffle(b, 3, 64, 104)
}
func BenchmarkShuffle6Players32Bits104Cards(b *testing.B) {
	benchmarkShuffle(b, 6, 32, 104)
}
func BenchmarkShuffle6Players64Bits104Cards(b *testing.B) {
	benchmarkShuffle(b, 6, 64, 104)
}

var benchErr error

func benchmarkShuffle(b *testing.B, playerCount int, bits int, cardCount int) {
	var err error
	for i := 0; i < b.N; i++ {
		d := deck.New(playersByBits[bits][:playerCount], cardCount)
		err = d.ResetAndShuffle()
	}
	benchErr = err
}

func init() {
	for _, bits := range []int{32, 64} {
		players := make([]deck.Player, 40)
		for i := 0; i < len(players); i++ {
			players[i] = deck.NewMe(prime, bits)
		}
		playersByBits[bits] = players
	}
}

func genPrime(bits int) *big.Int {
	ret, err := rand.Prime(rand.Reader, bits)
	if err != nil {
		log.Fatal(err)
	}
	return ret
}

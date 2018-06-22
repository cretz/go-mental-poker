package sra_test

import (
	"crypto/rand"
	"log"
	"math/big"
	"testing"

	"github.com/cretz/go-mental-poker/sra"
)

// Testing has shown prime size doesn't matter that much
var smallPrime = genPrime(64)
var mediumPrime = genPrime(256)
var largePrime = genPrime(1024)

func BenchmarkGenerateKeyPairSmallPrime32Bit(b *testing.B) {
	benchmarkGenerateKeyPair(b, smallPrime, 32)
}
func BenchmarkGenerateKeyPairSmallPrime64Bit(b *testing.B) {
	benchmarkGenerateKeyPair(b, smallPrime, 64)
}
func BenchmarkGenerateKeyPairSmallPrime128Bit(b *testing.B) {
	benchmarkGenerateKeyPair(b, smallPrime, 128)
}

func BenchmarkGenerateKeyPairMediumPrime32Bit(b *testing.B) {
	benchmarkGenerateKeyPair(b, mediumPrime, 32)
}
func BenchmarkGenerateKeyPairMediumPrime64Bit(b *testing.B) {
	benchmarkGenerateKeyPair(b, mediumPrime, 64)
}
func BenchmarkGenerateKeyPairMediumPrime128Bit(b *testing.B) {
	benchmarkGenerateKeyPair(b, mediumPrime, 128)
}

func BenchmarkGenerateKeyPairLargePrime32Bit(b *testing.B) {
	benchmarkGenerateKeyPair(b, largePrime, 32)
}
func BenchmarkGenerateKeyPairLargePrime64Bit(b *testing.B) {
	benchmarkGenerateKeyPair(b, largePrime, 64)
}
func BenchmarkGenerateKeyPairLargePrime128Bit(b *testing.B) {
	benchmarkGenerateKeyPair(b, largePrime, 128)
}
func BenchmarkGenerateKeyPairLargePrime256Bit(b *testing.B) {
	benchmarkGenerateKeyPair(b, largePrime, 256)
}
func BenchmarkGenerateKeyPairLargePrime512Bit(b *testing.B) {
	benchmarkGenerateKeyPair(b, largePrime, 512)
}
func BenchmarkGenerateKeyPairLargePrime1024Bit(b *testing.B) {
	benchmarkGenerateKeyPair(b, largePrime, 1024)
}
func BenchmarkGenerateKeyPairLargePrime2048Bit(b *testing.B) {
	benchmarkGenerateKeyPair(b, largePrime, 2048)
}

var resultKp *sra.KeyPair

func benchmarkGenerateKeyPair(b *testing.B, prime *big.Int, bits int) {
	var kp *sra.KeyPair
	var err error
	for i := 0; i < b.N; i++ {
		kp, err = sra.GenerateKeyPair(rand.Reader, prime, bits)
		if err != nil {
			b.Fatal(err)
		}
	}
	resultKp = kp
}

func genPrime(bits int) *big.Int {
	ret, err := rand.Prime(rand.Reader, bits)
	if err != nil {
		log.Fatal(err)
	}
	return ret
}

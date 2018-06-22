package sra

import (
	"crypto/rand"
	"io"
	"math/big"
)

// KeyPair is a commutative SRA key pair used for encryption and decryption.
type KeyPair struct {
	// Prime is the prime that was originally provided and is usually shared.
	Prime *big.Int
	// Enc is the big int used to encrypt.
	Enc *big.Int
	// Dec is the big int used to decrypt.
	Dec *big.Int
}

var bigOne = big.NewInt(1)

// GenerateKeyPair generates a SRA key pair for the given prime and numBits.
func GenerateKeyPair(rnd io.Reader, prime *big.Int, numBits int) (kp *KeyPair, err error) {
	kp = &KeyPair{Prime: prime}
	phiP := new(big.Int).Sub(prime, bigOne)
	// TODO: max tries?
	for {
		if kp.Enc, err = rand.Prime(rnd, numBits); err != nil {
			return nil, err
		}
		if new(big.Int).GCD(nil, nil, kp.Enc, phiP).Cmp(bigOne) == 0 {
			break
		}
	}
	kp.Dec = new(big.Int).ModInverse(kp.Enc, phiP)
	// TODO: since direction doesn't matter, I could switch Enc and Dec here.
	// Should I swap em randomly as shown below?
	// if v, err := rand.Int(rand.Reader, big.NewInt(2)); err == nil && v.Cmp(bigOne) == 0 {
	// 	kp.Enc, kp.Dec = kp.Dec, kp.Enc
	// }
	return
}

// EncryptInt returns v encrypted with Enc.
func (k *KeyPair) EncryptInt(v *big.Int) *big.Int {
	return new(big.Int).Exp(v, k.Enc, k.Prime)
}

// DecryptInt returns v decrypted with Dec.
func (k *KeyPair) DecryptInt(v *big.Int) *big.Int {
	return new(big.Int).Exp(v, k.Dec, k.Prime)
}

package sra

import (
	"crypto/rand"
	"io"
	"math/big"
)

type KeyPair struct {
	Prime *big.Int
	Enc   *big.Int
	Dec   *big.Int
}

var bigOne = big.NewInt(1)

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
	return
}

func (k *KeyPair) EncryptInt(v *big.Int) *big.Int {
	return new(big.Int).Exp(v, k.Enc, k.Prime)
}

func (k *KeyPair) DecryptInt(v *big.Int) *big.Int {
	return new(big.Int).Exp(v, k.Dec, k.Prime)
}

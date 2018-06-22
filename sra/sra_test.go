package sra_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/cretz/go-mental-poker/sra"
	"github.com/stretchr/testify/require"
)

func TestSRACommutative(t *testing.T) {
	// Gen prime
	prime, err := rand.Prime(rand.Reader, 256)
	require.NoError(t, err)
	// Gen key pairs for alice, bob, and ted
	alice, err := sra.GenerateKeyPair(rand.Reader, prime, 32)
	require.NoError(t, err)
	bob, err := sra.GenerateKeyPair(rand.Reader, prime, 32)
	require.NoError(t, err)
	ted, err := sra.GenerateKeyPair(rand.Reader, prime, 32)
	require.NoError(t, err)

	// Make sure it can be encrypted by all the people in any order, and decrypted
	peoplePerms := [][]*sra.KeyPair{
		{alice, bob, ted},
		{alice, ted, bob},
		{bob, alice, ted},
		{bob, ted, alice},
		{ted, alice, bob},
		{ted, bob, alice},
	}
	// Go over each set of people in any order and make sure the decrypted value always comes out right
	for _, encPeople := range peoplePerms {
		superSecretInt := newSuperSecretInt(t, 32)
		encrypted := encryptMulti(t, superSecretInt, encPeople)
		for _, decPeople := range peoplePerms {
			decrypted := decryptMulti(t, encrypted, decPeople)
			require.Zero(t, superSecretInt.Cmp(decrypted))
		}
	}
}

func encryptMulti(t *testing.T, v *big.Int, people []*sra.KeyPair) *big.Int {
	orig := v
	for _, person := range people {
		v = person.EncryptInt(v)
		// Make sure does match orig at any time
		if v.Cmp(orig) == 0 {
			require.FailNow(t, "Was orig after enc", orig)
		}
	}
	return v
}

func decryptMulti(t *testing.T, v *big.Int, people []*sra.KeyPair) *big.Int {
	orig := v
	for _, person := range people {
		v = person.DecryptInt(v)
		// Make sure does match orig at any time
		if v.Cmp(orig) == 0 {
			require.FailNow(t, "Was orig after dec", orig)
		}
	}
	return v
}

func newSuperSecretInt(t *testing.T, maxBits int) *big.Int {
	bigTwo := big.NewInt(2)
	max := new(big.Int).Sub(new(big.Int).Exp(bigTwo, big.NewInt(int64(maxBits)), nil), big.NewInt(1))
	for {
		superSecretInt, err := rand.Int(rand.Reader, max)
		require.NoError(t, err)
		if superSecretInt.Cmp(bigTwo) >= 0 {
			return superSecretInt
		}
	}
}

package deck

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/google/uuid"

	"github.com/cretz/go-mental-poker/sra"
)

type Player interface {
	ID() uuid.UUID
	// This encrypts all and shuffles
	ShuffleStage1(cards []*big.Int) error
	// This decrypts from stage one and encrypts each individually
	ShuffleStage2(cards []*big.Int) error
	// This is the final deck in the same order as stage 2 that can be used for mapping
	ShuffleComplete(cards []*big.Int) error
	// Can be nil if not found
	DecryptCard(origEncryptedCard *big.Int, valToDecrypt *big.Int) *big.Int
}

type Me struct {
	id          uuid.UUID
	sharedPrime *big.Int
	keyBits     int
	// Only non-nil after stage 1 and before stage 2
	tempShuffleStage1Pair *sra.KeyPair
	// Only non-nil after stage 2 and before complete
	tempShuffleStage2Pairs []*sra.KeyPair
	// Only non-nil on complete. Keyed by the encrypted card string.
	cardKeys           map[string]*sra.KeyPair
	DecryptedCards     []*big.Int
	OrigEncryptedCards []*big.Int
}

func NewMe(sharedPrime *big.Int, keyBits int) *Me {
	ret := &Me{sharedPrime: sharedPrime, keyBits: keyBits}
	var err error
	if ret.id, err = uuid.NewRandom(); err != nil {
		panic(err)
	}
	return ret
}

func (m *Me) ID() uuid.UUID { return m.id }

func (m *Me) ShuffleStage1(cards []*big.Int) (err error) {
	if m.tempShuffleStage1Pair != nil || m.tempShuffleStage2Pairs != nil {
		return fmt.Errorf("Another stage was left incomplete")
	}
	m.cardKeys = nil
	// Create a key pair for the entire deck
	if m.tempShuffleStage1Pair, err = sra.GenerateKeyPair(rand.Reader, m.sharedPrime, m.keyBits); err != nil {
		return
	}
	// Encrypt each card
	for i, card := range cards {
		cards[i] = m.tempShuffleStage1Pair.EncryptInt(card)
	}
	// Shuffle em
	newCryptoRand().Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	return
}

func (m *Me) ShuffleStage2(cards []*big.Int) (err error) {
	// TODO: Could check things like count and what not here
	if m.tempShuffleStage1Pair == nil || m.tempShuffleStage2Pairs != nil || m.cardKeys != nil {
		return fmt.Errorf("Stage 1 not complete")
	}
	m.tempShuffleStage2Pairs = make([]*sra.KeyPair, len(cards))
	for i, card := range cards {
		// Generate key pair for just this card
		if m.tempShuffleStage2Pairs[i], err = sra.GenerateKeyPair(rand.Reader, m.sharedPrime, m.keyBits); err != nil {
			break
		}
		// Decrypt what we had before and re-encrypt with card-specific key pair
		cards[i] = m.tempShuffleStage2Pairs[i].EncryptInt(m.tempShuffleStage1Pair.DecryptInt(card))
	}
	m.tempShuffleStage1Pair = nil
	return
}

func (m *Me) ShuffleComplete(cards []*big.Int) error {
	if m.tempShuffleStage1Pair != nil || len(m.tempShuffleStage2Pairs) != len(cards) || m.cardKeys != nil {
		return fmt.Errorf("Stage 2 not complete")
	}
	// Just map the cards to their keys
	m.cardKeys = make(map[string]*sra.KeyPair, len(cards))
	for i, card := range cards {
		m.cardKeys[card.String()] = m.tempShuffleStage2Pairs[i]
	}
	m.tempShuffleStage2Pairs = nil
	return nil
}

func (m *Me) DecryptCard(origEncryptedCard *big.Int, valToDecrypt *big.Int) *big.Int {
	// TODO: in a real implementation, this player would have for more information to make sure they are ok with
	// giving this up in this situation (e.g. info could include the player asking and could encrypt for that player).
	// TODO: Also note, we'd probably remove the key once its used in real-world circumstances
	cardPair := m.cardKeys[origEncryptedCard.String()]
	if cardPair == nil {
		return nil
	}
	return cardPair.DecryptInt(valToDecrypt)
}

func (m *Me) DrawCard(deck *Deck) error {
	// Grab card decrypted by everyone but me
	origEncryptedCard, mostlyDecryptedCard, err := deck.DrawCard(m.id)
	if err != nil {
		return err
	}
	// Decrypt it for me
	decryptedCard := m.DecryptCard(origEncryptedCard, mostlyDecryptedCard)
	if decryptedCard == nil {
		return fmt.Errorf("Can't find card decryption key")
	}
	m.DecryptedCards = append(m.DecryptedCards, decryptedCard)
	m.OrigEncryptedCards = append(m.OrigEncryptedCards, origEncryptedCard)
	return nil
}

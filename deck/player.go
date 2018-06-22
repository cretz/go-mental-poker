package deck

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/google/uuid"

	"github.com/cretz/go-mental-poker/sra"
)

// Player is an interface implemented by all players. This means it could be
// implemented with a remote player or a local one.
type Player interface {
	// ID is the unique identifier for this player.
	ID() uuid.UUID

	// ShuffleStage1 encrypts all cards with a single encryption key, stores
	// that key for stage 2, and shuffles the slice. The cards may be encrypted
	// from another player's stage-1 run or not.
	ShuffleStage1(cards []*big.Int) error

	// ShuffleStage2 decrypts each card from stage 1, then re-encrypts it with
	// a new per-card key, and stores that key by index for use on complete.
	// The cards may be encrypted from another player's stage-2 run or not.
	ShuffleStage2(cards []*big.Int) error

	// ShuffleComplete provides the completed, fully encrypted deck after all
	// players' stage-2 runs are done. It is in the same order as stage 2 and
	// the per-card keys from stage 2 can now be mapped to the fully-encrypted
	// card values.
	ShuffleComplete(cards []*big.Int) error

	// DecryptCard locates the decryption key for origEncryptionCard and
	// returns valToDecrypt decrypted with it. The valToDecrypt value may be
	// some already-half-decrypted value from other players.
	DecryptCard(origEncryptedCard *big.Int, valToDecrypt *big.Int) *big.Int
}

// Me is an implementation of Player for a local user.
type Me struct {
	id          uuid.UUID
	sharedPrime *big.Int
	keyBits     int
	// Only non-nil after stage 1 and before stage 2
	tempShuffleStage1Pair *sra.KeyPair
	// Only non-nil after stage 2 and before complete
	tempShuffleStage2Pairs []*sra.KeyPair
	// Only non-nil on complete. Keyed by the encrypted card string.
	cardKeys map[string]*sra.KeyPair
	// DecryptedCards are the current, decrypted cards in my hand.
	DecryptedCards []*big.Int
	// OrigEncryptedCards are the fully-encrypted values for DecryptedCards.
	OrigEncryptedCards []*big.Int
}

// NewMe creates a new local player with the given sharedPrime and keyBits
// count used to create the SRA key pair. This is assigned a random UUID ID.
func NewMe(sharedPrime *big.Int, keyBits int) *Me {
	ret := &Me{sharedPrime: sharedPrime, keyBits: keyBits}
	var err error
	if ret.id, err = uuid.NewRandom(); err != nil {
		panic(err)
	}
	return ret
}

// ID impls Player.ID.
func (m *Me) ID() uuid.UUID { return m.id }

// ShuffleStage1 impls Player.ShuffleStage1.
func (m *Me) ShuffleStage1(cards []*big.Int) (err error) {
	if m.tempShuffleStage1Pair != nil || m.tempShuffleStage2Pairs != nil {
		return fmt.Errorf("Another stage was left incomplete")
	}
	m.cardKeys = nil
	m.DecryptedCards = nil
	m.OrigEncryptedCards = nil
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

// ShuffleStage2 impls Player.ShuffleStage2.
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

// ShuffleComplete impls Player.ShuffleComplete.
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

// DecryptCard impls Player.DecryptCard.
func (m *Me) DecryptCard(origEncryptedCard *big.Int, valToDecrypt *big.Int) *big.Int {
	// TODO: In a real implementation, this player would have for more
	// information to make sure they are ok with giving this up in this
	// situation (e.g. info could include the player asking or whether it was
	// their turn).
	// TODO: Also note, we'd probably remove the key once its used in
	// real-world circumstances, especially in games where discarded cards can
	// be reused.
	cardPair := m.cardKeys[origEncryptedCard.String()]
	if cardPair == nil {
		return nil
	}
	return cardPair.DecryptInt(valToDecrypt)
}

// DrawCard draws the next card off the deck and puts it in my hand.
func (m *Me) DrawCard(deck *Deck) error {
	// Grab card decrypted by everyone but me
	origEncryptedCard, mostlyDecryptedCard, err := deck.DrawCard(m.id)
	if err != nil {
		return err
	}
	// Decrypt it for me which means, as the last one to decrypt, that it is
	// fully decrypted.
	decryptedCard := m.DecryptCard(origEncryptedCard, mostlyDecryptedCard)
	if decryptedCard == nil {
		return fmt.Errorf("Can't find card decryption key")
	}
	m.DecryptedCards = append(m.DecryptedCards, decryptedCard)
	m.OrigEncryptedCards = append(m.OrigEncryptedCards, origEncryptedCard)
	return nil
}

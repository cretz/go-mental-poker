package deck

import (
	"fmt"
	"math/big"

	"github.com/google/uuid"
)

// Deck is collection of cards and players. Note, that cards within are from 2
// to count + 2 due to the requirement that there be no "0" or "1" card.
type Deck struct {
	players []Player
	count   int
	// Must be positive integers > 1
	cards []*big.Int
}

// New creates a new deck for the given player set and count. Note, the cards
// within the deck range from 2 to count + 2.
func New(players []Player, count int) *Deck {
	return &Deck{players: players, count: count}
}

// ResetAndShuffle first resets the deck to cards 2 to count + 2. Then the
// three shuffle steps are executed across the players for secure shuffling.
func (d *Deck) ResetAndShuffle() error {
	// First, reset to 2 to count + 2
	d.cards = make([]*big.Int, d.count)
	for i := 0; i < d.count; i++ {
		d.cards[i] = big.NewInt(int64(i + 2))
	}
	// Have each player run stage 1 of the shuffle which chains requests for
	// each to encrypt the entire deck and shuffle it.
	for _, player := range d.players {
		if err := player.ShuffleStage1(d.cards); err != nil {
			return err
		}
	}
	// Have each player run stage 2 of the shuffle which takes the shuffled
	// and everyone-encrypted deck and asks each player to re-encrypt their
	// cards with a key per card.
	for _, player := range d.players {
		if err := player.ShuffleStage2(d.cards); err != nil {
			return err
		}
	}
	// Tell each player what the completed deck looks like. This allows them
	// to map their per-card keys to the full-encrypted card values.
	for _, player := range d.players {
		if err := player.ShuffleComplete(d.cards); err != nil {
			return err
		}
	}
	return nil
}

// DrawCard takes a card off the end of the deck and decrypts it from all
// players except playerIDToLeaveEncryptedFor (usually the asking player). If
// playerIDToLeaveEncryptedFor is uuid.Nil or otherwise doesn't match any
// players, the mostlyDecryptedCard result value will be fully decrypted.
//
// Note, in a more serious implementation, checks would be done that confirm
// who is asking and that they can at a certain time (e.g. it is their turn).
// Also, the players would be smarter about validating the decryption requests.
func (d *Deck) DrawCard(
	playerIDToLeaveEncryptedFor uuid.UUID,
) (origEncryptedCard *big.Int, mostlyDecryptedCard *big.Int, err error) {
	origEncryptedCard = d.cards[len(d.cards)-1]
	d.cards = d.cards[:len(d.cards)-1]
	mostlyDecryptedCard, err = d.MostlyRevealCard(origEncryptedCard, playerIDToLeaveEncryptedFor)
	return
}

// MostlyRevealCard takes the given fully-encrypted card and decrypts it from
// all players except playerIDToLeaveEncryptedFor (usually the asking player).
// If playerIDToLeaveEncryptedFor is uuid.Nil or otherwise doesn't match any
// players, the mostlyDecryptedCard result value will be fully decrypted.
//
// Note, this is exposed for debugging purposes and in a more serious
// implementation it would not even exist and no reasonable-written player
// would let this caller decrypt a card whenever it wanted.
func (d *Deck) MostlyRevealCard(
	origEncryptedCard *big.Int,
	playerIDToLeaveEncryptedFor uuid.UUID,
) (mostlyDecryptedCard *big.Int, err error) {
	mostlyDecryptedCard = origEncryptedCard
	// Decrypt the card from all other players but the given one
	for _, player := range d.players {
		if player.ID() != playerIDToLeaveEncryptedFor {
			mostlyDecryptedCard = player.DecryptCard(origEncryptedCard, mostlyDecryptedCard)
			if mostlyDecryptedCard == nil {
				return nil, fmt.Errorf("No decrypted card from %v", player)
			}
		}
	}
	return
}

// RevealCards returns the revealed cards in the deck. This is only for
// debugging purposes and in a real-world implementation this would not exist
// and not be possible because the players would balk at decryption requests.
func (d *Deck) RevealCards() (revealed []*big.Int, err error) {
	revealed = make([]*big.Int, len(d.cards))
	for i, card := range d.cards {
		// Ask to reveal from all players
		if revealed[i], err = d.MostlyRevealCard(card, uuid.Nil); err != nil {
			break
		}
	}
	return
}

package deck

import (
	"fmt"
	"math/big"

	"github.com/google/uuid"
)

type Deck struct {
	players []Player
	count   int
	// Must be positive integers > 1
	cards []*big.Int
}

func New(players []Player, count int) *Deck {
	return &Deck{players: players, count: count}
}

func (d *Deck) ResetAndShuffle() error {
	// First, set em back to all their indices
	d.cards = make([]*big.Int, d.count)
	for i := 0; i < d.count; i++ {
		d.cards[i] = big.NewInt(int64(i + 2))
	}
	// Now run through shuffle stage 1 for all players
	for _, player := range d.players {
		if err := player.ShuffleStage1(d.cards); err != nil {
			return err
		}
	}
	// Now stage 2
	for _, player := range d.players {
		if err := player.ShuffleStage2(d.cards); err != nil {
			return err
		}
	}
	// Now complete
	for _, player := range d.players {
		if err := player.ShuffleComplete(d.cards); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deck) DrawCard(
	playerIDToLeaveEncryptedFor uuid.UUID,
) (origEncryptedCard *big.Int, mostlyDecryptedCard *big.Int, err error) {
	origEncryptedCard = d.cards[len(d.cards)-1]
	d.cards = d.cards[:len(d.cards)-1]
	mostlyDecryptedCard, err = d.MostlyRevealCard(origEncryptedCard, playerIDToLeaveEncryptedFor)
	return
}

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

func (d *Deck) RevealCards() (revealed []*big.Int, err error) {
	revealed = make([]*big.Int, len(d.cards))
	for i, card := range d.cards {
		if revealed[i], err = d.MostlyRevealCard(card, uuid.Nil); err != nil {
			break
		}
	}
	return
}

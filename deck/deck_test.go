package deck_test

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"testing"

	"github.com/cretz/go-mental-poker/deck"
	"github.com/stretchr/testify/require"
)

func TestSimpleDraw(t *testing.T) {
	// Get/print all cards in deck in order
	allCardsInOrder := make([]Card, 52)
	for i := 0; i < 52; i++ {
		allCardsInOrder[i] = CardFromInt(i)
	}
	fmt.Printf("%-19v %v\n", "All cards:", allCardsInOrder)

	// Create a prime everyone shares
	sharedPrime, err := rand.Prime(rand.Reader, 256)
	require.NoError(t, err)

	// Create three players
	alice := deck.NewMe(sharedPrime, 32)
	bob := deck.NewMe(sharedPrime, 32)
	ted := deck.NewMe(sharedPrime, 32)

	// Create a deck of cards
	d := deck.New([]deck.Player{alice, bob, ted}, 52)

	// Do a shuffle
	require.NoError(t, d.ResetAndShuffle())
	fmt.Printf("%-19v %v\n", "Deck after shuffle:", deckCards(t, d))

	// Give each player 7 cards
	for i := 0; i < 7; i++ {
		require.NoError(t, alice.DrawCard(d))
		require.NoError(t, bob.DrawCard(d))
		require.NoError(t, ted.DrawCard(d))
	}
	fmt.Printf("%-19v %v\n", "Deck after draws:", deckCards(t, d))
	fmt.Printf("%-19v %v\n", "Alice's draw:", playerCards(t, alice))
	fmt.Printf("%-19v %v\n", "Bob's draw:", playerCards(t, bob))
	fmt.Printf("%-19v %v\n", "Ted's draw:", playerCards(t, ted))

	// Combine em all and confirm it's the same elements as the full set
	finalCards := append([]Card{}, deckCards(t, d)...)
	finalCards = append(finalCards, playerCards(t, alice)...)
	finalCards = append(finalCards, playerCards(t, bob)...)
	finalCards = append(finalCards, playerCards(t, ted)...)
	require.ElementsMatch(t, allCardsInOrder, finalCards)
}

type Card struct {
	// 0 through 3
	Suit int
	// 2 through 14
	Card int
}

// v is 0-53, mod 4 is suit, div 4 is 0-aligned card (so we add 2 to get 2 through 14)
func CardFromInt(v int) Card { return Card{Suit: v % 4, Card: (v / 4) + 2} }

func CardFromBigInt(v *big.Int) (card Card, err error) {
	if cardInt64 := v.Int64(); !v.IsInt64() || cardInt64 > math.MaxInt32 {
		err = fmt.Errorf("Invalid card")
	} else {
		// The cardInt is from 2 through 55, so we subtract 2 to 0-align it
		card = CardFromInt(int(cardInt64 - 2))
	}
	return
}

func CardsFromBigInts(v []*big.Int) (cards []Card, err error) {
	cards = make([]Card, len(v))
	for i, card := range v {
		if cards[i], err = CardFromBigInt(card); err != nil {
			break
		}
	}
	return
}

func (c Card) SuitChar() rune {
	switch c.Suit {
	case 0:
		return '♠'
	case 1:
		return '♥'
	case 2:
		return '♦'
	default:
		return '♣'
	}
}

func (c Card) CardName() string {
	switch c.Card {
	case 11:
		return "J"
	case 12:
		return "Q"
	case 13:
		return "K"
	case 14:
		return "A"
	default:
		return strconv.Itoa(c.Card)
	}
}

func (c Card) String() string {
	return string(c.CardName()) + string(c.SuitChar())
}

func deckCards(t *testing.T, d *deck.Deck) []Card {
	revealed, err := d.RevealCards()
	require.NoError(t, err)
	cards, err := CardsFromBigInts(revealed)
	require.NoError(t, err)
	return cards
}

func playerCards(t *testing.T, me *deck.Me) []Card {
	cards, err := CardsFromBigInts(me.DecryptedCards)
	require.NoError(t, err)
	return cards
}

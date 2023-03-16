package utils

import (
	"math/rand"
	"strings"
	"time"
	"fmt"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/lando-ke/card-api/models"
)

var (
	values = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "JACK", "QUEEN", "KING", "ACE"}
	suits  = []string{"SPADES", "DIAMONDS", "CLUBS", "HEARTS"}
)

func NewDeck(db *gorm.DB, shuffled bool, cardsParam string) (models.Deck, error) {
	deck := models.Deck{
		DeckID:   uuid.New().String(),
		Shuffled: shuffled,
	}

	// Create cards and associate them with the deck
	cards := createCards(shuffled, cardsParam)
	deck.Remaining = len(cards) // Set the remaining count dynamically based on the created cards

	// Save the deck to the database
	if err := db.Create(&deck).Error; err != nil {
		return models.Deck{}, err
	}

	for _, card := range cards {
		card.DeckID = deck.DeckID
		if err := db.Create(&card).Error; err != nil {
			return models.Deck{}, err
		}
	}

	// Retrieve the cards associated with the deck and set the Cards field
	if err := db.Where("deck_id = ?", deck.DeckID).Find(&deck.Cards).Error; err != nil {
		return models.Deck{}, err
	}

	return deck, nil
}

func createCards(shuffled bool, cardsParam string) []models.Card {
	var cards []models.Card

	if cardsParam != "" {
		cards = CreatePartialDeck(cardsParam)
	} else {
		cards = CreateFullDeck()
	}

	if shuffled {
		cards = ShuffleCards(cards)
	}

	return cards
}

func CreateFullDeck() []models.Card {
	cards := []models.Card{}

	for _, suit := range suits {
		for _, value := range values {
			card := models.Card{
				Value: value,
				Suit:  suit,
				Code:  value[:1] + suit[:1],
			}
			cards = append(cards, card)
		}
	}

	return cards
}

func CreatePartialDeck(cardsParam string) []models.Card {
	cardCodes := strings.Split(cardsParam, ",")
	cards := []models.Card{}

	codeToValueSuit := make(map[string]struct {
		Value string
		Suit  string
	})

	for _, suit := range suits {
		for _, value := range values {
			code := value[:1] + suit[:1]
			if value == "10" {
				code = value + suit[:1]
			}
			codeToValueSuit[code] = struct {
				Value string
				Suit  string
			}{Value: value, Suit: suit}
		}
	}

	for _, cardCode := range cardCodes {
		cardCode = strings.ToUpper(cardCode)
		valueSuit, ok := codeToValueSuit[cardCode]
		if !ok {
			continue
		}

		card := models.Card{
			Value: valueSuit.Value,
			Suit:  valueSuit.Suit,
			Code:  cardCode,
		}

		cards = append(cards, card)
	}

	return cards
}


func ShuffleCards(cards []models.Card) []models.Card {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	return cards
}

func generateDeckID() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}

func ValidateCardsParam(cardsParam string) []string {
    cards := strings.Split(cardsParam, ",")
    partialDeck := CreatePartialDeck(cardsParam)
    invalidCards := []string{}

    cardMap := make(map[string]bool)
    cardCounts := make(map[string]int)

    for _, card := range partialDeck {
        cardMap[strings.ToUpper(card.Code)] = true
    }

    for _, card := range cards {
        card = strings.ToUpper(card)
        if _, exists := cardMap[card]; !exists {
            invalidCards = append(invalidCards, card)
        } else {
            cardCounts[card]++
            if cardCounts[card] > 1 {
                invalidCards = append(invalidCards, fmt.Sprintf("%s (duplicate)", card))
            }
        }
    }

    return invalidCards
}


func isValidCard(cardStr string) bool {
	if len(cardStr) < 2 || len(cardStr) > 3 {
		return false
	}

	value := cardStr[:len(cardStr)-1]
	suit := cardStr[len(cardStr)-1:]

	for _, validSuit := range suits {
		for _, validValue := range values {
			if validSuit == suit && validValue == value {
				return true
			}
		}
	}

	return false
}


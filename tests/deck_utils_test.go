package tests

import (
	"github.com/lando-ke/card-api/models"
	"github.com/lando-ke/card-api/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"reflect"
)

func setupDatabase(t *testing.T) *gorm.DB {
	// Setup the database connection
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to the in-memory database: %v", err)
	}

	// Migrate the models
	err = db.AutoMigrate(&models.Deck{}, &models.Card{})
	if err != nil {
		t.Fatalf("failed to migrate models: %v", err)
	}

	return db
}

func TestNewDeck_Shuffled(t *testing.T) {
	db := setupDatabase(t)

	// Create a shuffled deck
	shuffledDeck, err := utils.NewDeck(db, true, "")
	if err != nil {
		t.Fatalf("failed to create shuffled deck: %v", err)
	}

	// Create an unshuffled deck
	unshuffledDeck, err := utils.NewDeck(db, false, "")
	if err != nil {
		t.Fatalf("failed to create unshuffled deck: %v", err)
	}

	// Check if the shuffled deck is different from the unshuffled deck
	shuffledCount := 0
	for i, card := range shuffledDeck.Cards {
		if card.Code != unshuffledDeck.Cards[i].Code {
			shuffledCount++
		}
	}

	if shuffledCount <= len(shuffledDeck.Cards)*3/4 {
		t.Errorf("the shuffled deck is not sufficiently shuffled, only %d cards are in different positions", shuffledCount)
	}
}

func TestNewDeck_Unshuffled(t *testing.T) {
	db := setupDatabase(t)

	// Create an unshuffled deck
	unshuffledDeck, err := utils.NewDeck(db, false, "")
	if err != nil {
		t.Fatalf("failed to create unshuffled deck: %v", err)
	}

	// Create a reference deck
	referenceDeck := utils.CreateFullDeck()

	// Check if the unshuffled deck is the same as the reference deck
	for i, card := range unshuffledDeck.Cards {
		if card.Code != referenceDeck[i].Code {
			t.Errorf("the unshuffled deck is not in the correct order, card at position %d has code %s instead of %s", i, card.Code, referenceDeck[i].Code)
		}
	}
}

func TestCreateFullDeck(t *testing.T) {
	fullDeck := utils.CreateFullDeck()

	if len(fullDeck) != 52 {
		t.Errorf("full deck should have 52 cards, but has %d", len(fullDeck))
	}

	cardCount := make(map[string]int)
	for _, card := range fullDeck {
		cardCount[card.Code]++
	}

	for _, card := range fullDeck {
		if cardCount[card.Code] != 1 {
			t.Errorf("card with code %s should appear exactly once, but appears %d times", card.Code, cardCount[card.Code])
		}
	}
}

func TestCreatePartialDeck(t *testing.T) {
	cardsParam := "AS,KH,2D,JC,10C"
	partialDeck := utils.CreatePartialDeck(cardsParam)

	if len(partialDeck) != 5 {
		t.Errorf("partial deck should have 5 cards, but has %d", len(partialDeck))
	}

	expectedCards := []string{"AS", "KH", "2D", "JC", "10C"}
	for i, card := range partialDeck {
		if card.Code != expectedCards[i] {
			t.Errorf("card at position %d should have code %s, but has %s", i, expectedCards[i], card.Code)
		}
	}
}

func TestValidateCardsParam(t *testing.T) {
	cardsParam := "AS,KH,2D,JC,10C,XX,2D,KH"
	invalidCards := utils.ValidateCardsParam(cardsParam)

	expectedInvalidCards := []string{"XX", "2D (duplicate)", "KH (duplicate)"}
	if !reflect.DeepEqual(invalidCards, expectedInvalidCards) {
		t.Errorf("invalid cards should be %v, but got %v", expectedInvalidCards, invalidCards)
	}
}


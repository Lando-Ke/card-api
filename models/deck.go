package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Deck struct {
	gorm.Model
	DeckID    string `json:"deck_id" gorm:"type:uuid;primary_key"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
	Cards     []Card `json:"cards" gorm:"foreignKey:DeckID"`
}

func (deck Deck) MarshalJSON() ([]byte, error) {
	type Alias Deck
	return json.Marshal(&struct {
		*Alias
		DeckID    string `json:"deck_id"`
		Shuffled  bool   `json:"shuffled"`
		Remaining int    `json:"remaining"`
		Cards     []Card `json:"cards"`
	}{
		Alias:     (*Alias)(&deck),
		DeckID:    deck.DeckID,
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining,
		Cards:     deck.Cards,
	})
}


func (deck *Deck) BeforeCreate(tx *gorm.DB) (err error) {
	deck.DeckID = uuid.New().String()
	deck.CreatedAt = time.Now()
	deck.UpdatedAt = time.Now()
	return
}

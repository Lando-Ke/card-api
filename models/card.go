package models

import (
	"encoding/json"
	"gorm.io/gorm"
)

type Card struct {
	gorm.Model
	Value  string `json:"value" gorm:"type:varchar(255)"`
	Suit   string `json:"suit" gorm:"type:varchar(255)"`
	Code   string `json:"code" gorm:"type:varchar(255)"`
	DeckID string `json:"-" gorm:"index"`
}

func (card Card) MarshalJSON() ([]byte, error) {
	type Alias Card
	return json.Marshal(&struct {
		*Alias
		Value string `json:"value"`
		Suit  string `json:"suit"`
		Code  string `json:"code"`
	}{
		Alias: (*Alias)(&card),
		Value: card.Value,
		Suit:  card.Suit,
		Code:  card.Code,
	})
}

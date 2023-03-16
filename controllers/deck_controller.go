package controllers

import (
	"gorm.io/gorm"
	"net/http"
	"strings"
	"strconv"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lando-ke/card-api/models"
	"github.com/lando-ke/card-api/utils"
)

type DeckController struct {
	db *gorm.DB
}

type DeckResponse struct {
	DeckID    string         `json:"deck_id"`
	Shuffled  bool           `json:"shuffled"`
	Remaining int            `json:"remaining"`
	Cards     []CardResponse `json:"cards"`
}

type CardResponse struct {
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

func cardModelToResponse(card models.Card) CardResponse {
	return CardResponse{
		Value: card.Value,
		Suit:  card.Suit,
		Code:  card.Code,
	}
}

func NewDeckController(db *gorm.DB) *DeckController {
	return &DeckController{db}
}

func (dc *DeckController) CreateDeck(c *gin.Context) {
	shuffled := c.Query("shuffled") == "true"
	cardsParam := c.Query("cards")

	if cardsParam != "" {
		invalidCards := utils.ValidateCardsParam(cardsParam)
		if len(invalidCards) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid cards values: " + strings.Join(invalidCards, ", ")})
			return
		}
	}

	deck, err := utils.NewDeck(dc.db, shuffled, cardsParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating deck"})
		return
	}

	response := DeckResponse{
		DeckID:    deck.DeckID,
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining,
	}

	for _, card := range deck.Cards {
		response.Cards = append(response.Cards, cardModelToResponse(card))
	}

	c.JSON(http.StatusOK, response)
}

func (dc *DeckController) OpenDeck(c *gin.Context) {
	deckID := c.Param("deck_id")

	var deck models.Deck
	if err := dc.db.Where("deck_id = ?", deckID).First(&deck).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
		return
	}

	var cards []models.Card
	if err := dc.db.Model(&models.Card{}).Where("deck_id = ? AND deleted_at IS NULL", deckID).Find(&cards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading cards"})
		return
	}

	// Convert cards to CardResponse
	cardResponses := []CardResponse{}
	for _, card := range cards {
		cardResponses = append(cardResponses, cardModelToResponse(card))
	}

	response := DeckResponse{
		DeckID:    deck.DeckID,
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining,
		Cards:     cardResponses,
	}

	c.JSON(http.StatusOK, response)
}

func (dc *DeckController) DrawCard(c *gin.Context) {
	deckID := c.Param("deck_id")
	countStr := c.DefaultQuery("count", "1")

	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count parameter"})
		return
	}

	var deck models.Deck
	if err := dc.db.Where("deck_id = ?", deckID).First(&deck).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
		return
	}

	if count > deck.Remaining {
    	c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("not enough cards in deck, only %d remaining", deck.Remaining)})
    	return
	}

	var cards []models.Card
	if err := dc.db.Where("deck_id = ?", deckID).Order("id ASC").Limit(count).Find(&cards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error loading cards"})
		return
	}

	drawnCards := cards[:count]

	// Convert drawnCards to CardResponse
	drawnCardResponses := []CardResponse{}
	for _, card := range drawnCards {
		drawnCardResponses = append(drawnCardResponses, cardModelToResponse(card))
	}

	// Update the remaining card count
	deck.Remaining = deck.Remaining - count
	if err := dc.db.Model(&models.Deck{}).Where("deck_id = ?", deckID).Update("remaining", deck.Remaining).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating remaining card count"})
		return
	}

	// Soft delete the drawn cards
	for _, card := range drawnCards {
		dc.db.Delete(&models.Card{}, card.ID)
	}

	c.JSON(http.StatusOK, gin.H{"cards": drawnCardResponses})
}

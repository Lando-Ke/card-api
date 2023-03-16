package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lando-ke/card-api/models"
	"github.com/lando-ke/card-api/controllers"
	"github.com/lando-ke/card-api/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&models.Deck{})
	db.AutoMigrate(&models.Card{})

	return db
}

func TestCreateDeck(t *testing.T) {
	db := setupDB()
	deckController := controllers.NewDeckController(db)
	gin.SetMode(gin.TestMode)

	t.Run("create unshuffled deck", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/deck?shuffled=false", nil)
		c.Request = req

		deckController.CreateDeck(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response controllers.DeckResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		assert.Equal(t, 52, response.Remaining)
		assert.Equal(t, false, response.Shuffled)
	})

	t.Run("create shuffled deck", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest("GET", "/deck?shuffled=true", nil)
		c.Request = req

		deckController.CreateDeck(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response controllers.DeckResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		assert.Equal(t, 52, response.Remaining)
		assert.Equal(t, true, response.Shuffled)
	})

	t.Run("create partial deck", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
	
		req, _ := http.NewRequest("GET", "/deck?shuffled=false&cards=AS,KS,QS,JS,10S", nil)
		c.Request = req
	
		deckController.CreateDeck(c)
	
		assert.Equal(t, http.StatusOK, w.Code)
	
		var response controllers.DeckResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}
	
		assert.Equal(t, 5, response.Remaining)
		assert.Equal(t, false, response.Shuffled)
	
		expectedCards := []controllers.CardResponse{
			{Value: "ACE", Suit: "SPADES", Code: "AS"},
			{Value: "KING", Suit: "SPADES", Code: "KS"},
			{Value: "QUEEN", Suit: "SPADES", Code: "QS"},
			{Value: "JACK", Suit: "SPADES", Code: "JS"},
			{Value: "10", Suit: "SPADES", Code: "10S"},
		}
	
		for i, expectedCard := range expectedCards {
			assert.Equal(t, expectedCard.Value, response.Cards[i].Value)
			assert.Equal(t, expectedCard.Suit, response.Cards[i].Suit)
			assert.Equal(t, expectedCard.Code, response.Cards[i].Code)
		}
	})

	t.Run("create_partial_deck_with_duplicates", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/deck?cards=AS,AS,KS,QS,JS,10H,kh", nil)

		deckController.CreateDeck(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "invalid cards values: AS (duplicate)", response["message"])
	})

	t.Run("create_partial_deck_with_invalid_card", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/deck?cards=AS,KS,QS,JS,10H,kh,XX", nil)
		deckController.CreateDeck(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "invalid cards values: XX", response["message"])
	})

	t.Run("create_partial_deck_with_invalid_and_duplicate_cards", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/deck?cards=AS,KS,KS,QS,JS,10H,kh,KH,XX", nil)

		deckController.CreateDeck(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "invalid cards values: KS (duplicate), KH (duplicate), XX", response["message"])
	})
}

func TestOpenDeck(t *testing.T) {
	db := setupDB()

	dc := controllers.NewDeckController(db)

	t.Run("open_full_deck", func(t *testing.T) {
		// Create a deck first
		deck, _ := utils.NewDeck(db, false, "")
	
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/deck/"+deck.DeckID, nil)
		c.Params = []gin.Param{{Key: "deck_id", Value: deck.DeckID}}
	
		dc.OpenDeck(c)
	
		assert.Equal(t, http.StatusOK, w.Code)
		var response controllers.DeckResponse
		json.Unmarshal(w.Body.Bytes(), &response)
	
		assert.Equal(t, deck.DeckID, response.DeckID)
		assert.Equal(t, 52, response.Remaining)
		assert.Len(t, response.Cards, 52)
	})
	
	

	t.Run("open_partial_deck", func(t *testing.T) {
		// Create a partial deck first
		cardsParam := "AS,KH,2D,JC,10C"
		deck, _ := utils.NewDeck(db, false, cardsParam)
	
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/deck/"+deck.DeckID, nil)
		c.Params = []gin.Param{{Key: "deck_id", Value: deck.DeckID}} // Make sure we pass the deck_id in the context parameters
	
		dc.OpenDeck(c)
	
		assert.Equal(t, http.StatusOK, w.Code)
		var response controllers.DeckResponse
		json.Unmarshal(w.Body.Bytes(), &response)
	
		assert.Equal(t, deck.DeckID, response.DeckID)
		assert.Equal(t, 5, response.Remaining)
		assert.Len(t, response.Cards, 5)
	})
	

	t.Run("open_non_existent_deck", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/deck/nonexistentdeck123", nil)

		dc.OpenDeck(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Deck not found", response["error"])
	})
}


func TestDrawCard(t *testing.T) {
	db := setupDB()
	dc := controllers.NewDeckController(db)

	t.Run("draw_one_from_full_deck", func(t *testing.T) {
		deck, _ := utils.NewDeck(db, false, "")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/deck/"+deck.DeckID+"/draw?count=1", nil)
		c.Params = []gin.Param{{Key: "deck_id", Value: deck.DeckID}}

		dc.DrawCard(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string][]controllers.CardResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Len(t, response["cards"], 1)
	})

	t.Run("draw_multiple_from_full_deck", func(t *testing.T) {
		deck, _ := utils.NewDeck(db, false, "")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/deck/"+deck.DeckID+"/draw?count=5", nil)
		c.Params = []gin.Param{{Key: "deck_id", Value: deck.DeckID}}

		dc.DrawCard(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string][]controllers.CardResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Len(t, response["cards"], 5)
	})

	t.Run("draw_more_than_remaining", func(t *testing.T) {
		deck, _ := utils.NewDeck(db, false, "")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/deck/"+deck.DeckID+"/draw?count=55", nil)
		c.Params = []gin.Param{{Key: "deck_id", Value: deck.DeckID}}

		dc.DrawCard(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("draw_from_partial_deck", func(t *testing.T) {
		cardsParam := "AS,KH,2D,JC,10C"
		deck, _ := utils.NewDeck(db, false, cardsParam)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/deck/"+deck.DeckID+"/draw?count=3", nil)
		c.Params = []gin.Param{{Key: "deck_id", Value: deck.DeckID}}

		dc.DrawCard(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string][]controllers.CardResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Len(t, response["cards"], 3)
	})

	t.Run("invalid_count_parameter", func(t *testing.T) {
		deck, _ := utils.NewDeck(db, false, "")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/deck/"+deck.DeckID+"/draw?count=-1", nil)
		c.Params = []gin.Param{{Key: "deck_id", Value: deck.DeckID}}

		dc.DrawCard(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

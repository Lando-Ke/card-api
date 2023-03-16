package routes

import (
	"github.com/lando-ke/card-api/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterDeckRoutes(r *gin.Engine, db *gorm.DB) {
	deckController := controllers.NewDeckController(db)
	r.POST("/deck", deckController.CreateDeck)
	r.GET("/deck/:deck_id", deckController.OpenDeck)
	r.GET("/deck/:deck_id/draw", deckController.DrawCard)
}

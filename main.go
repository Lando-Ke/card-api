package main

import (
	"github.com/lando-ke/card-api/database"
	"github.com/lando-ke/card-api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	dbInstance, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	err = database.RunMigrations(dbInstance)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	routes.RegisterDeckRoutes(r, dbInstance)
	r.Run(":8080")
}

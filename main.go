package main

import (
	"github.com/gin-gonic/gin"
	"github.com/relumini/shortdl/database"
	"github.com/relumini/shortdl/routes"
)

func main() {
	_, err := database.ConnectDB()
	if err != nil {
		panic("failed to connect database")
	}

	// Auto Migrate the model
	r := gin.Default()
	routes.InitRoute(r)

	r.Run()
}

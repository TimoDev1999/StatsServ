package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	err := db.InitDB()
	if err != nil {
		log.Fatalf("Could not initialize the database: %v", err)
	}
	defer db.CloseDB()

	r := gin.Default()

	r.GET("/orderbook/:exchange_name/:pair", handlers.GetOrderBook)
	r.POST("/orderbook", handlers.SaveOrderBook)
	r.GET("/orderhistory", handlers.GetOrderHistory)
	r.POST("/order", handlers.SaveOrder)

	log.Println("Starting server on :8080")
	r.Run(":8080")
}

//

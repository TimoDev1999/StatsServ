package main

import (
	"log"
	"statsServTask/db"
	"statsServTask/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	if err := db.InitDB(); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	// Создаем роутер Gin
	r := gin.Default()

	// Определяем маршруты с замыканием для передачи базы данных
	r.GET("/orderbook/:exchange_name/:pair", func(c *gin.Context) {
		handlers.GetOrderBook(c, db.DB)
	})
	r.POST("/orderbook", func(c *gin.Context) {
		handlers.SaveOrderBook(c, db.DB)
	})
	r.GET("/orderhistory", func(c *gin.Context) {
		handlers.GetOrderHistory(c, db.DB)
	})
	r.POST("/order", func(c *gin.Context) {
		handlers.SaveOrder(c, db.DB)
	})

	// Запускаем сервер
	if err := r.Run(":8484"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

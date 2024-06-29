package handlers

import (
	"database/sql"
	"net/http"
	"stats_service/db"
	"stats_service/models"

	"github.com/gin-gonic/gin"
)

func GetOrderBook(c *gin.Context) {
	exchangeName := c.Param("exchange_name")
	pair := c.Param("pair")

	var orderBook models.OrderBook
	row := db.DB.QueryRow(`SELECT id, exchange, pair, asks, bids 
                       FROM order_book 
					   WHERE exchange=$1 AND pair=$2 LIMIT 1`, exchangeName, pair)
	err := row.Scan(&orderBook.Id, &orderBook.Exchange, &orderBook.Pair, &orderBook.Asks, &orderBook.Bids)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "OrderBook not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, orderBook)
}

func SaveOrderBook(c *gin.Context) {

	var orderBook models.OrderBook

	if err := c.ShouldBindJSON(&orderBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := db.DB.Exec(`INSERT INTO order_book (exchange, pair, asks, bids) VALUES ($1, $2, $3, $4)`,
		orderBook.Exchange, orderBook.Pair, orderBook.Asks, orderBook.Bids)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "OrderBook saved"})
}

func GetOrderHistory(c *gin.Context) {
	clientName := c.Query("client_name")

	exchangeName := c.Query("exchange_name")
	rows, err := db.DB.Query(`SELECT client_name, exchange_name, label, pair, side, type, base_qty, price, algorithm_name_placed, lowest_sell_prc, highest_buy_prc, commission_quote_qty, time_placed
                              FROM order_history
                              WHERE client_name=$1 AND exchange_name=$2`, clientName, exchangeName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var orderHistory []models.HistoryOrder

	for rows.Next() {
		var order models.HistoryOrder
		err := rows.Scan(&order.ClientName, &order.ExchangeName, &order.Label, &order.Pair, &order.Side, &order.Type, &order.BaseQty, &order.Price, &order.AlgorithmNamePlaced, &order.LowestSellPrc, &order.HighestBuyPrc, &order.CommissionQuoteQty, &order.TimePlaced)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orderHistory = append(orderHistory, order)
	}
	c.JSON(http.StatusOK, orderHistory)
}

func SaveOrder(c *gin.Context) {
	var order models.HistoryOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := db.DB.Exec(`INSERT INTO order_history (client_name, exchange_name, label, pair, side, type, base_qty, price, 
	                     algorithm_name_placed, lowest_sell_prc, highest_buy_prc, commission_quote_qty, time_placed)
                         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		order.ClientName, order.ExchangeName, order.Label, order.Pair, order.Side, order.Type, order.BaseQty, order.Price,
		order.AlgorithmNamePlaced, order.LowestSellPrc, order.HighestBuyPrc, order.CommissionQuoteQty, order.TimePlaced)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "Order saved"})
}

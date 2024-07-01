package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"statsServTask/handlers"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetOrderBook(t *testing.T) {
	// Создание мокированной базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock: %v", err)
	}
	defer db.Close()

	// Настройка моков
	rows := sqlmock.NewRows([]string{"id", "exchange", "pair", "asks", "bids"}).
		AddRow(1, "test_exchange", "BTC/USD", `[{"price":50000,"baseqty":1}]`, `[{"price":49000,"baseqty":1}]`)
	mock.ExpectQuery(`SELECT id, exchange, pair, asks, bids FROM order_book WHERE exchange=\$1 AND pair=\$2 LIMIT 1`).
		WithArgs("test_exchange", "BTC/USD").
		WillReturnRows(rows)

	// Настройка маршрутов
	router := gin.Default()
	router.GET("/orderbook/:exchange_name/:pair", func(c *gin.Context) {
		handlers.GetOrderBook(c, db)
	})

	// Создание запроса
	req, _ := http.NewRequest("GET", "/orderbook/test_exchange/BTC/USD", nil)
	w := httptest.NewRecorder()

	// Выполнение запроса
	router.ServeHTTP(w, req)

	// Проверка результатов
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test_exchange")
}

func TestSaveOrderBook(t *testing.T) {
	// Создание мокированной базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock: %v", err)
	}
	defer db.Close()

	// Настройка моков
	mock.ExpectExec(`INSERT INTO order_book \(exchange, pair, asks, bids\) VALUES \(\$1, \$2, \$3, \$4\)`).
		WithArgs("test_exchange", "BTC/USD", `[{"price":50000,"baseqty":1}]`, `[{"price":49000,"baseqty":1}]`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Настройка маршрутов
	router := gin.Default()
	router.POST("/orderbook", func(c *gin.Context) {
		handlers.SaveOrderBook(c, db)
	})

	// Создание запроса
	body := `{"exchange":"test_exchange","pair":"BTC/USD","asks":[{"price":50000,"baseqty":1}],"bids":[{"price":49000,"baseqty":1}]}`
	req, _ := http.NewRequest("POST", "/orderbook", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Выполнение запроса
	router.ServeHTTP(w, req)

	// Проверка результатов
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "OrderBook saved")
}

func TestGetOrderHistory(t *testing.T) {
	// Создание мокированной базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock: %v", err)
	}
	defer db.Close()

	// Настройка моков
	rows := sqlmock.NewRows([]string{"client_name", "exchange_name", "label", "pair", "side", "type", "base_qty", "price", "algorithm_name_placed", "lowest_sell_prc", "highest_buy_prc", "commission_quote_qty", "time_placed"}).
		AddRow("test_client", "test_exchange", "test_label", "BTC/USD", "buy", "limit", 1.0, 50000.0, "test_algo", 50000.0, 49000.0, 0.01, time.Now())
	mock.ExpectQuery(`SELECT client_name, exchange_name, label, pair, side, type, base_qty, price, algorithm_name_placed, lowest_sell_prc, highest_buy_prc, commission_quote_qty, time_placed FROM order_history WHERE client_name=\$1 AND exchange_name=\$2`).
		WithArgs("test_client", "test_exchange").
		WillReturnRows(rows)

	// Настройка маршрутов
	router := gin.Default()
	router.GET("/orderhistory", func(c *gin.Context) {
		handlers.GetOrderHistory(c, db)
	})

	// Создание запроса
	req, _ := http.NewRequest("GET", "/orderhistory?client_name=test_client&exchange_name=test_exchange", nil)
	w := httptest.NewRecorder()

	// Выполнение запроса
	router.ServeHTTP(w, req)

	// Проверка результатов
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test_client")
}

func TestSaveOrder(t *testing.T) {
	// Создание мокированной базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock: %v", err)
	}
	defer db.Close()

	// Настройка моков
	mock.ExpectExec(`INSERT INTO order_history \(client_name, exchange_name, label, pair, side, type, base_qty, price, algorithm_name_placed, lowest_sell_prc, highest_buy_prc, commission_quote_qty, time_placed\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10, \$11, \$12, \$13\)`).
		WithArgs("test_client", "test_exchange", "test_label", "BTC/USD", "buy", "limit", 1.0, 50000.0, "test_algo", 50000.0, 49000.0, 0.01, "2023-07-01T00:00:00Z").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Настройка маршрутов
	router := gin.Default()
	router.POST("/order", func(c *gin.Context) {
		handlers.SaveOrder(c, db)
	})

	// Создание запроса
	body := `{
		"client_name": "test_client",
		"exchange_name": "test_exchange",
		"label": "test_label",
		"pair": "BTC/USD",
		"side": "buy",
		"type": "limit",
		"base_qty": 1.0,
		"price": 50000.0,
		"algorithm_name_placed": "test_algo",
		"lowest_sell_prc": 50000.0,
		"highest_buy_prc": 49000.0,
		"commission_quote_qty": 0.01,
		"time_placed": "2023-07-01T00:00:00Z"
	}`
	req, _ := http.NewRequest("POST", "/order", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Выполнение запроса
	router.ServeHTTP(w, req)

	// Проверка результатов
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Order saved")
}

Сервер статистики торгов
Этот проект представляет собой сервер для работы с статистикой торговых ордеров. Он использует базу данных для хранения и доступа к данным.

Убедитесь, что у вас установлена и настроена PostgreSQL база данных. В файле db/db.go указаны параметры подключения к базе данных. Измените эти параметры под свои нужды перед инициализацией.

Сервер предоставляет следующие HTTP endpoints:

GET /orderbook/:exchange_name/:pair - Получение статистики ордеров для заданной биржи и пары.

POST /orderbook - Сохранение статистики ордеров.

GET /orderhistory - Получение истории ордеров для заданного клиента и биржи.

POST /order - Сохранение информации ордера.

Примеры использования

Получение статистики ордеров:

curl -X GET http://localhost:8484/orderbook/test_exchange/BTC/USD

Сохранение статистики ордеров:

curl -X POST -H "Content-Type: application/json" -d '{"exchange":"test_exchange","pair":"BTC/USD","asks":[{"price":50000,"baseqty":1}],"bids":[{"price":49000,"baseqty":1}]}' http://localhost:8484/orderbook

Получение истории ордеров:

curl -X GET http://localhost:8484/orderhistory?client_name=test_client&exchange_name=test_exchange

Сохранение информации ордера:

curl -X POST -H "Content-Type: application/json" -d '{"client_name": "test_client","exchange_name": "test_exchange","label": "test_label","pair": "BTC/USD","side": "buy","type": "limit","base_qty": 1.0,"price": 50000.0,"algorithm_name_placed": "test_algo","lowest_sell_prc": 50000.0,"highest_buy_prc": 49000.0,"commission_quote_qty": 0.01,"time_placed": "2023-07-01T00:00:00Z"}' http://localhost:8484/order

# KirBank - REST API Банковского Сервиса

## Содержание
- [Описание](#описание)
- [Технологии](#технологии)
- [Структура проекта](#структура-проекта)
- [Сборка и запуск](#сборка-и-запуск)
- [Использование API](#использование-api)
  - [Аутентификация](#аутентификация)
  - [Управление счетами](#управление-счетами)
  - [Операции с картами](#операции-с-картами)
  - [Кредитные операции](#кредитные-операции)
  - [Транзакции](#транзакции)
  - [Аналитика](#аналитика)
  

## Описание

**KirBank** — это RESTful API для управления банковскими операциями. Проект предоставляет функционал для регистрации и аутентификации пользователей, создания и управления счетами, выпуска карт, оформления кредитов, выполнения переводов и получения аналитики финансовых операций. API обеспечивает безопасность данных с использованием шифрования (bcrypt для паролей, PGP для карт) и интеграцию с внешними сервисами, такими как ЦБ РФ для получения ключевой ставки и SMTP для отправки уведомлений.

### Основные функции:

- Регистрация и аутентификация пользователей с использованием JWT (срок действия токена — 24 часа).
- Создание банковских счетов, пополнение и списание средств.
- Выпуск виртуальных карт с генерацией номеров по алгоритму Луна и шифрованием данных.
- Оформление кредитов с расчетом аннуитетных платежей и автоматическим списанием.
- Переводы между счетами и аналитика транзакций (доходы/расходы, кредитная нагрузка).
- Прогноз баланса счета на срок до 365 дней.


## Технологии
| Стек                    | Назначение                           |
|-------------------------|-----------------------------------------------|
| Go 1.23+                   | Язык программирования                                        |
| gorilla/mux            | Маршрутизация HTTP-запросов                                 |
| lib/pq                 | Драйвер для PostgreSQL             |
| golang-jwt/jwt        | Аутентификация с JWT                        |
| sirupsen/logrus              | Логирование                         |
| golang.org/x/crypto/bcrypt       | Хеширование паролей      |
| gopkg.in/gomail.v2                  | Отправка email-уведомлений                     |

## Структура проекта
```bash
├───cmd
├───config
├───internal
│   ├───config
│   ├───handlers
│   ├───middleware
│   ├───models
│   ├───repositories
│   ├───services
│   └───utils
└───sql
```

## Сборка и запуск

### 1. Клонировать репозиторий

```bash
git clone https://github.com/yourusername/banking-api.git
cd banking-api
```

### 2. Настройка базы данных

- Установите PostgreSQL 17 и создайте базу данных
- Настройте схему базы данных, выполнив SQL-скрипты (sql/db_completion)
- 
### 3. Настройка переменных окружения

Создайте файл config.yaml, пример файла config.yaml (вставьте свои данные):
```yaml
  server:
  port: 8080

database:
  url: postgres://username:password@localhost:5432/db_name?sslmode=disable

auth:
  jwt_secret: secret
  hmac_secret: secret
  encryption_key: key

smtp:
  host: smtp.yandex.com
  port: 587
  user: pochta@yandex.ru
  pass: password
  from: pochta@yandex.ru
```

### 4. Установка зависимостей

```bash
go mod tidy
```
Файл go.mod
```go
require (
	golang.org/x/crypto v0.38.0
	github.com/beevik/etree v1.5.1 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-chi/chi/v5 v5.2.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
```

### 5. Запуск приложения

```bash
go run cmd/server/main.go
```
API будет доступно по адресу http://localhost:8080 (или на другом порту, если указано в конфигурации).

## Использование API

### Аутентификация
|Метод |	URL|	Описание|	Тело запроса|	Ответ|
|------|----|-----------|--------------|------|
POST|	/register	|Регистрация нового пользователя|	{ "email": "string", "username": "string", "password": "string" }|	201 Created с деталями пользователя|
POST|	/login	|Вход и получение JWT-токена|	{ "email": "string", "password": "string" }	|200 OK с { "token": "string" }|

**Пример запроса POST /register**
```http
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "user123",
  "password": "securepassword"
}
```
**Ответ 201 Created**
```json
{
  "id": "123",
  "username": "user123",
  "email": "user@example.com"
}
```
**Пример запроса POST /login**
```http
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```
**Ответ 200 OK**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
Для всех защищенных эндпоинтов добавляйте заголовок `Authorization: Bearer <token>`.

### Управление счетами
|Метод |	URL|	Описание|	Тело запроса|	Ответ|
|------|----|-----------|--------------|------|
POST|	/accounts	|Создание нового счета|	-	|201 Created с деталями счета
GET|	/accounts/{accountId}/balance|	Получение баланса счета|	-	|200 OK с { "balance": float }
GET|	/accounts/{accountId}/predict|	Прогноз баланса счета|	-	|200 OK с { "predicted_balance": float, "key_rate": float }

**Пример запроса POST /accounts**
```http
POST /accounts
Authorization: Bearer <token>
Content-Type: application/json
```
**Ответ 201 Created**

```json
{
  "id": "456",
  "account_number": "1234567890",
  "balance": 0.0,
  "currency": "RUB"
}
```
**Пример запроса GET /accounts/{accountId}/balance**
```http
GET /accounts/456/balance
Authorization: Bearer <token>
```
**Ответ 200 OK**
```json
{
  "balance": 1000.50
}
```
**Пример запроса GET /accounts/{accountId}/predict**

```http
GET /accounts/456/predict
Authorization: Bearer <token>
```
**Ответ 200 OK**
```json
{
  "predicted_balance": 950.75,
  "key_rate": 7.5
}
```

### Операции с картами

|Метод |	URL|	Описание|	Тело запроса|	Ответ|
|------|----|-----------|--------------|------|
POST|	/cards|	Выпуск новой карты	|{ "account_id": "string" }|	201 Created с деталями карты

**Пример запроса POST /cards**
```http
POST /cards
Authorization: Bearer <token>
Content-Type: application/json

{
  "account_id": "456"
}
```
**Ответ 201 Created**
```json
{
  "card_number": "4111111111111111",
  "expiry": "12/27",
  "cvv": "123"
}
```
### Кредитные операции
|Метод |	URL|	Описание|	Тело запроса|	Ответ|
|------|----|-----------|--------------|------|
|POST|	/credits|	Оформление нового кредита|	{ "account_id": "string", "amount": float, "interest_rate": float, "term_months": int }	|201 Created с деталями кредита|
|GET|	/credits/{creditId}/schedule|	Получение графика платежей|	-|	200 OK с графиком платежей|

**Пример запроса POST /credits**
```http
POST /credits
Authorization: Bearer <token>
Content-Type: application/json

{
  "account_id": "456",
  "amount": 100000.0,
  "interest_rate": 10.0,
  "term_months": 12
}
```
**Ответ 201 Created**
```json
{
  "id": "789",
  "amount": 100000.0,
  "interest_rate": 10.0,
  "term_months": 12,
  "status": "active"
}
```
**Пример запроса GET /credits/{creditId}/schedule**
```http
GET /credits/789/schedule
Authorization: Bearer <token>
```
**Ответ 200 OK**
```json
[
  {
    "payment_date": "2025-06-01",
    "payment_amount": 8791.67,
    "principal_amount": 7916.67,
    "interest_amount": 875.0,
    "status": "pending"
  },
  ...
]
```
### Транзакции
|Метод |	URL|	Описание|	Тело запроса|	Ответ|
|------|----|-----------|--------------|------|
|POST|	/transfer	|Перевод между счетами|	{ "from_account": "string", "to_account": "string", "amount": float }	|200 OK с деталями транзакции|

**Пример запроса POST /transfer**
```http
POST /transfer
Authorization: Bearer <token>
Content-Type: application/json

{
  "from_account": "456",
  "to_account": "789",
  "amount": 500.0
}
```
**Ответ 200 OK**
```json
{
  "transaction_id": "101112",
  "status": "success"
}
```
### Аналитика
|Метод |	URL|	Описание|	Тело запроса|	Ответ|
|------|----|-----------|--------------|------|
|GET|	/analytics|	Получение аналитики транзакций|	-|	200 OK с данными аналитики|

**Пример запроса GET /analytics**
```http
GET /analytics
Authorization: Bearer <token>
```
**Ответ 200 OK**
```json
{
  "sent": {
    "count": 10,
    "total": 5000.0
  },
  "received": {
    "count": 5,
    "total": 3000.0
  }
}
```


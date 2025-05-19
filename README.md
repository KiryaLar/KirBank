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

### 1️⃣ Клонировать репозиторий

```bash
git clone https://github.com/username/defcode.git
cd defcode
```

### 2️⃣ Настроить окружение

- Установите PostgreSQL 15+ и создайте базу defcode_db.

- Создайти application.yaml на основе application.yaml.origin из проекта
### 3️⃣ Собрать и запустить
```bash
mvn clean install      # сборка и юнит‑тесты
mvn spring-boot:run    # запуск
```
После старта API доступно по адресу http://localhost:8080.

## Использование API

    Во всех примерах HOST = http://localhost:8080.
    Формат даты времени — ISO 8601.

## Аутентификация
| Метод | URL                  | Тело запроса                  | Описание                                |
|-------|-----------------------|-------------------------------|-----------------------------------------|
| POST  | `/auth/register`      | `{username, password, role}`  | Регистрация (роль: `ADMIN` или `USER`)  |
| POST  | `/auth/login`         | `{username, password}`        | Получить `accessToken` и `refreshToken` |
| POST  | `/auth/refresh-token` | `{refreshToken}`              | Обновить JWT‑пару                        |
| POST  | `/auth/logout`        | `{refreshToken}`              | Отозвать refresh‑токен                  |

### Пример запроса POST `/auth/register`
```http
POST /auth/register
Content-Type: application/json

{
  "username": "user",
  "password": "Secret123",
  "role": "user",
}
```
**Ответ 201 Created**
```json
{
  "message": "User alice has successfully registered as an USER",
  "timestamp": "Tue May 01 12:45:01 GMT+03:00 2025"
}
```

### Пример запроса POST `/auth/login`
```http
POST /auth/login
Content-Type: application/json

{
  "username": "user",
  "password": "Secret123"
}
```
**Ответ 200 OK**
```json
{
  "role": "USER",
  "accessToken": "<JWT>",
  "refreshToken": "<JWT>"
}
```
### Пример запроса POST `/auth/logout`
```http
POST /auth/logout
Content-Type: application/json

{
  "refreshToken": "<ваш_refresh_token>"
}
```
**Ответ 200 OK**
```json
{
  "message": "Success logged out",
  "timestamp": "Tue May 01 12:45:01 GMT+03:00 2025"
}
```
### Пример запроса POST `/auth/refresh-token`
```http
POST /auth/refresh-token
Content-Type: application/json

{
  "refreshToken": "<ваш_refresh_token>"
}
```
**Ответ 200 OK**
```json
{
  "role": "USER",
  "accessToken": "<новый_access_token>",
  "refreshToken": "<новый_refresh_token>"
}
```

Добавляйте заголовок Authorization: Bearer <accessToken> ко всем защищённым эндпоинтам.

## Пользовательские OTP-операции

| Метод | URL                   | Тело запроса                         | Назначение                     |
|-------|------------------------|--------------------------------------|--------------------------------|
| POST  | `/user/otp/generate`   | `{method, contact, operationType}`   | Сгенерировать и отправить код |
| POST  | `/user/otp/validate`   | `{code}`                             | Проверить введённый код       |

### Значения `operationType`:
1. Login Verification
2. Account Registration
3. Password Reset
4. Transaction Confirmation
5. Update Contact Information
6. Account Deletion

### Пример генерации кода по SMS
```http
POST /user/otp/generate
Authorization: Bearer <accessToken>
Content-Type: application/json

{
  "method": "sms",
  "contact": "+79179997799",
  "operationType": "4"
}
```
### Пример генерации кода по email
```http
POST /user/otp/generate
Authorization: Bearer <accessToken>
Content-Type: application/json

{
  "method": "email",
  "contact": "user@yande.ru",
  "operationType": "2"
}
```
### Пример генерации кода по telegram
```http
POST /user/otp/generate
Authorization: Bearer <accessToken>
Content-Type: application/json

{
  "method": "telegram",
  "contact": "<telegram bot chat id>",
  "operationType": "4"
}
```
**Ответ 200 OK**
```json
{
  "message": "Code was successfully sent",
  "timestamp": "Tue May 01 12:45:01 GMT+03:00 2025"
}
```

### Пример валидации кода
```http
POST /user/otp/validate
Content-Type: application/json

{
  "code": "123456"
}
```
**Ответ 200 OK**
```json
{
  "message": "Code is correct",
  "timestamp": "Tue May 01 12:45:01 GMT+03:00 2025"
}
```

## Администрирование
| Метод | URL                  | Тело запроса               | Функция                             |
|-------|-----------------------|----------------------------|-------------------------------------|
| PUT   | `/admin/otp-config`   | `{length, lifetime}`       | Изменить длину и TTL кода           |
| GET   | `/admin/users`        | –                          | Получить список пользователей (USER)|
| DELETE| `/admin/users/{id}`   | –                          | Удалить пользователя                |
### Значения `lifetime` в виде число + один символ из (smhd):
1. 30s
2. 2m
3. 3h
4. 1d

### Пример изменения конфигурации
```http
PUT /admin/otp-config
Authorization: Bearer <accessToken>
Content-Type: application/json

{
  "length": 8,
  "lifetime": "2m"   // 2 минуты 
}
```
**Ответ 200 OK**
```json
{
  "message": "OTP config was successfully changed: New duration: 2m New code length: 8",
  "timestamp": "Tue May 01 12:45:01 GMT+03:00 2025"
}
```

### Пример получения списка пользователей
```http
GET /admin/users
Authorization: Bearer <access_token>
```
**Ответ 200 OK**
```json
[
  {
    "id": 1,
    "username": "user1"
  },
  {
    "id": 2,
    "username": "user2"
  }
]
```

### Пример удаления пользователя
```http
GET /admin/users/{id}
Authorization: Bearer <access_token>
```
**Ответ 200 OK**
```json
{
  "message": "User with id 1 deleted",
  "timestamp": "Tue May 01 12:55:01 GMT+03:00 2025"
}
```

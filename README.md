# REST API Банковского Сервиса

Проект представляет собой REST API для банковского сервиса на языке Go. 

## Технологии

*   **Go 1.23+**
*   **PostgreSQL 17** + `pgcrypto`
*   **Gorilla Mux** - маршрутизация
*   **Logrus** - логирование
*   **JWT** - аутентификация
*   **Bcrypt & HMAC & PGP** - защита данных
*   **Etree & Gomail** - интеграции с ЦБ РФ и SMTP

## Запуск проекта

1.  **База данных:**
    Запустите PostgreSQL в Docker с помощью `docker-compose`:
    ```bash
    docker-compose up -d
    ```
    База данных будет доступна на порту 5432, а миграции из папки `migrations` выполнятся автоматически при старте.

2.  **Запуск приложения:**
    Установите зависимости и запустите `main.go`:
    ```bash
    go mod tidy
    go run cmd/api/main.go
    ```

## Доступные эндпоинты

### Публичные
*   `POST /register` — Регистрация пользователя (JSON: username, email, password)
*   `POST /login` — Аутентификация пользователя (JSON: username, password)

### Защищенные (требуют заголовок `Authorization: Bearer <jwt-token>`)
*   `POST /accounts` — Создать счет (JSON: currency)
*   `POST /accounts/deposit` — Пополнить счет (JSON: account_id, amount)
*   `POST /transfer` — Перевод между счетами (JSON: from_account_id, to_account_id, amount)
*   `POST /cards` — Выпустить виртуальную карту по алгоритму Луна (JSON: account_id)
*   `GET /accounts/{accountId}/cards` — Получить список карт (с расшифровкой номера)
*   `POST /credits` — Оформить кредит (расчет ставки по ЦБ РФ + аннуитетные платежи) (JSON: amount, term_months)
*   `GET /credits/{creditId}/schedule` — Посмотреть график платежей по кредиту
*   `GET /accounts/{accountId}/analytics` — Статистика доходов и расходов
*   `GET /accounts/{accountId}/predict?days=30` — Прогноз баланса

## Особенности реализации
- Планировщик платежей по кредитам реализован с помощью `time.Ticker` и запускается в фоне при старте приложения, проверяя просроченные платежи.
- PGP ключи генерируются в памяти при старте, так как это упрощает процесс тестирования без необходимости монтировать внешние ключи.
- SMTP настроен в мок-режиме, письма логируются в stdout.

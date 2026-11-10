# OrderFlow

OrderFlow — учебный backend-сервис для оформления заказов. Проект развивается последовательно: от простого HTTP API на Go до сервиса с PostgreSQL, Redis, Kafka, gRPC, Docker и Kubernetes.

Цель проекта — не собрать демонстрацию отдельных технологий, а постепенно построить один связный сервис и на его примере изучить архитектуру backend-приложений, работу с данными, интеграции и эксплуатацию.

## Что уже работает

- HTTP API на Go;
- проверка входных данных и расчёт стоимости заказа;
- создание и получение заказов;
- хранение заказов и позиций в PostgreSQL;
- транзакционное сохранение заказа;
- SQL-миграции;
- unit- и интеграционные тесты;
- сборка API через multi-stage Dockerfile;
- запуск API, PostgreSQL, тестовой PostgreSQL, миграций и Redis через Docker Compose.

Redis уже добавлен в инфраструктуру, но кеширование заказов пока не подключено к API.

## Как работает сервис

Текущая цепочка вызова:

```text
HTTP-клиент
    ↓
HTTP handler
    ↓
Order Service
    ↓
Repository interface
    ↓
PostgreSQL repository
    ↓
PostgreSQL
```

При создании заказа сервис:

1. декодирует JSON-запрос;
2. преобразует его в доменные структуры;
3. проверяет название, цену и количество каждой позиции;
4. создаёт заказ со статусом `new`;
5. в одной транзакции записывает заказ и его позиции;
6. возвращает ID, статус и общую стоимость.

При чтении заказа репозиторий сначала получает сам заказ, затем его позиции и собирает их обратно в один доменный объект.

## API

По умолчанию API доступен по адресу `http://localhost:8080`.

### Проверка состояния

```http
GET /health
```

```json
{"status":"ok"}
```

### Предварительный расчёт

```http
POST /orders/preview
Content-Type: application/json
```

```json
{
  "items": [
    {
      "name": "Keyboard",
      "price": 1500,
      "quantity": 2
    }
  ]
}
```

Ответ:

```json
{"total":3000}
```

### Создание заказа

```http
POST /orders
Content-Type: application/json
```

Тело запроса такое же, как у предварительного расчёта.

Ответ:

```json
{
  "id": 1,
  "status": "new",
  "total": 3000
}
```

### Получение заказа

```http
GET /orders/1
```

```json
{
  "id": 1,
  "status": "new",
  "total": 3000
}
```

## Технологии

| Технология | Назначение | Состояние |
|---|---|---|
| Go | HTTP API, доменная логика и сервисный слой | Работает |
| PostgreSQL | Основное хранилище заказов | Работает |
| pgx | Драйвер и пул соединений PostgreSQL | Работает |
| SQL migrations | Версионирование схемы БД | Работает |
| Docker | Сборка контейнера API | Работает |
| Docker Compose | Локальный запуск инфраструктуры | Работает |
| Redis | Кеш заказов | В разработке |
| gRPC | Взаимодействие между сервисами | Запланировано |
| Kafka | Асинхронные события заказов | Запланировано |
| Kubernetes | Развёртывание и управление сервисами | Запланировано |

## Архитектурные идеи

В проекте используются или будут добавлены:

- разделение HTTP, сервисного и инфраструктурного слоёв;
- Dependency Injection через конструкторы;
- Repository для отделения бизнес-логики от хранилища;
- транзакции PostgreSQL;
- Decorator и Cache-Aside для Redis;
- асинхронные события через Kafka;
- идемпотентность обработки сообщений;
- Outbox Pattern для надёжной публикации событий;
- graceful shutdown, таймауты и корректная обработка ошибок.

## Структура проекта

```text
.
├── compose.yaml                  # локальная инфраструктура
├── migrations/                   # миграции PostgreSQL
├── README.md
└── orderflow/
    ├── Dockerfile
    ├── go.mod
    ├── cmd/order-api/            # точка входа и HTTP handlers
    └── internal/
        ├── order/                # доменная модель, сервис, интерфейсы
        └── storage/
            ├── postgres/         # PostgreSQL repository
            └── redis/            # Redis cache repository
```

## Запуск через Docker Compose

Требования:

- Docker;
- Docker Compose.

Запустить весь стек:

```bash
docker compose up -d --build
```

Проверить состояние:

```bash
docker compose ps -a
curl http://localhost:8080/health
```

Контейнеры миграций выполняются один раз и должны завершиться со статусом `Exited (0)`.

Посмотреть логи API:

```bash
docker compose logs -f api
```

Остановить проект:

```bash
docker compose down
```

Команда сохраняет данные в Docker volumes. `docker compose down -v` дополнительно удалит данные обеих баз.

## Локальный запуск API

PostgreSQL можно оставить в Docker, а API запустить на компьютере:

```bash
docker compose up -d postgres migrate
cd orderflow
go run ./cmd/order-api
```

Для локального запуска используются переменные:

```env
DATABASE_URL=postgres://orderflow:orderflow@localhost:5432/orderflow?sslmode=disable
TEST_DATABASE_URL=postgres://orderflow:orderflow@localhost:5433/orderflow?sslmode=disable
```

Внутри Docker Compose API обращается к базе по адресу `postgres:5432`, потому что сервисы используют внутреннюю Docker-сеть.

## Тесты

Поднять тестовую БД и применить миграции:

```bash
docker compose up -d test-postgres migrate-test
```

Запустить тесты:

```bash
cd orderflow
go test ./...
```

Интеграционные тесты подключаются к отдельной PostgreSQL на `localhost:5433` и не используют основную базу.

## План развития

1. Закончить Redis-кеш по стратегии Cache-Aside.
2. Добавить изменение статуса и подтверждение заказа.
3. Добавить gRPC-взаимодействие с отдельным сервисом склада.
4. Публиковать события `order.created` и `order.confirmed` в Kafka.
5. Добавить надёжную доставку событий через Outbox Pattern.
6. Добавить логирование, метрики и трассировку.
7. Подготовить Kubernetes-манифесты и развернуть сервисы в кластере.


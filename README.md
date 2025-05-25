# ITS Reg: сервис для создания регистраций на мероприятия

## Как запускать

### Локальный запуск

Для локального запуска сервера и для проведения интеграционных тестов используется _локальное окружение_.

Поднимать при помощи `docker compose`:

```shell
docker compose -f deployment/docker-compose.local.yaml up -d
```

Запуск тестов:
- `go test -short ./internal/domain/...` - юнит-тесты домена;
- `go test -count=1 ./internal/infra/...` - интеграционные (инфраструктурные) тесты; также требуется задать переменную
  окружения `DATABASE_URI` (для local - `postgres://test-user:test-pass@localhost:30001/test-db?sslmode=disable`);

Запуск HTTP сервера API - `go run ./cmd/http/http.go` (требуется задать переменные окружения `PORT` и `DATABASE_URI`)

### Дев

Для запуска полной, рабочей копии, проекта, но для локального запуска используется _dev окружение_.

Поднимать при помощи `docker compose`:

```shell
docker compose -f deployment/docker-compose.dev.yaml up -d
```

Загружает окружение из `.env` файла. Пример содержимого `.env` представлен в файле 

### Продакшен

Поднимать при помощи `docker compose`:

```shell
docker compose -f deployment/docker-compose.prod.yaml up -d
```

Загружает окружение из `.env` файла. Пример содержимого `.env` представлен в файле 

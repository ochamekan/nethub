# Нетхаб - тестовое задание

### Требования

- Backend:
  - Go‑сервис с БД, который хранит устройства и их статусы.
  - REST‑эндпоинты для CRUD‑операций и смены статуса is_active.
  - Логирование основных операций (создание, обновление, удаление).
- Frontend:
  - TypeScript‑клиент, который:
    - показывает список устройств и их статус;
    - позволяет фильтровать и искать устройства;
    - позволяет создавать и редактировать устройства.
- Документация:
  - README с описанием:
    - как запустить backend и frontend;
    - какие зависимости нужны;
    - примеры запросов к API.
- Дополнительно (по желанию кандидата):
  - docker‑compose для запуска БД и сервиса.


## 🚀 Запуск проекта

### Backend

Запущен на `localhost:8000`

```bash
cd backend

# запуск проекта
./run.sh

# запуск тестов
./test.sh

# очистка контейнеров
./cleanup.sh
```

### Frontend

Запущен на `localhost:5173`

```bash
cd frontend

# запуск проекта
./run.sh

# очистка контейнеров
./cleanup.sh
```
## Стек приложения

**Backend:** Go, PostgreSQL (pgx), Zap (Для логирования), Goose (для миграций)

**Frontend:** React + TS, Vite, Tailwind + shadcn/ui, React Router, TanStack Query, React Hook Form

## 📦 Зависимости

* Docker
* Docker Compose
* (опционально) Go — для локальной разработки backend
* (опционально) Node.js + pnpm — для локальной разработки frontend

## Эндпоинты

| Метод | Endpoint                  | Описание |
|-------|---------------------------|----------|
| `POST`   | `/v1/devices`             | Создание нового устройства |
| `GET`    | `/v1/devices`             | Получение списка устройств |
| `GET`    | `/v1/devices/{id}`        | Получение устройства по ID |
| `PUT`    | `/v1/devices/{id}`        | Обновление устройства |
| `DELETE` | `/v1/devices/{id}`        | Мягкое удаление устройства (пометка флагом) |


## 📡 API примеры

### Получить список устройств

**GET** `/devices`

#### Query параметры:

* `search` — поиск по имени/описанию
* `is_active` — фильтр по активности (`true` / `false`)

#### Примеры запросов:

```bash
# все устройства
curl "http://localhost:8000/v1/devices"
```

```bash
# поиск устройств
curl "http://localhost:8000/v1/devices?search=foobar"
```

```bash
# только активные устройства
curl "http://localhost:8000/v1/devices?is_active=true"
```

```bash
# поиск + фильтр
curl "http://localhost:8000/v1/devices?search=foobar&is_active=true"
```

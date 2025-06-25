# WorkmateTestTask

HTTP-сервис для управления I/O-зависимыми задачами (время обработки 1–5 минут) с хранением в памяти.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
- [Installation](#installation)
- [Build & Run](#build--run)
- [Usage](#usage)
- [Logging & Graceful Shutdown](#logging--graceful-shutdown)
- [CI/CD (GitHub Actions)](#ci/cd-github-actions)

## Prerequisites

- Go 1.24 или выше

## Configuration

Сервис по умолчанию слушает порт `8080`. Можно изменить через переменную окружения:

```bash
export PORT=9090
```

Также можно ограничить число одновременно обрабатываемых задач через `MAX_CONCURRENT_TASKS` (по умолчанию 10):

```bash
export MAX_CONCURRENT_TASKS=5
```

## Installation

```bash
git clone https://github.com/nniicckk6/workmateTestTask.git
cd workmateTestTask
go mod download
```

## Build & Run

### Из исходников
```bash
go build -o workmateTestTask main.go
PORT=8080 ./workmateTestTask
```

### Через Docker
```bash
# Собрать образ
docker build -t workmate-test-task .
# Запустить контейнер
docker run -d -p 8080:8080 \
  -e PORT=8080 \
  -e MAX_CONCURRENT_TASKS=10 \
  workmate-test-task
```

### Через Docker Compose
```bash
docker-compose up -d
```

Сервис доступен по адресу http://localhost:${PORT}

## Usage

Ниже примеры работы с API через `curl`.

### Create Task
```bash
curl -X POST http://localhost:${PORT}/tasks \
  -H "Content-Type: application/json"
``` 
Ответ с кодом **201**:
```json
{ "id": "<uuid>", "status": "Pending", "created_at": "2025-06-25T12:34:56Z" }
```

### List Tasks
```bash
curl http://localhost:${PORT}/tasks
``` 
Ответ **200** – массив объектов.

### Get Task
```bash
curl http://localhost:${PORT}/tasks/<uuid>
``` 
Ответ **200**:
```json
{
  "id": "<uuid>",
  "status": "Completed",
  "created_at": "2025-06-25T12:34:56Z",
  "started_at": "2025-06-25T12:35:00Z",
  "finished_at": "2025-06-25T12:37:30Z",
  "duration": "2m30s",
  "result": "Обработано за 2m30s"
}
```

### Delete Task
```bash
curl -X DELETE http://localhost:${PORT}/tasks/<uuid>
``` 
Ответ **204 No Content**

## Logging & Graceful Shutdown

- Логи запросов и времени обработки выводятся в стандартный вывод.
- При получении сигналов SIGINT/SIGTERM сервер корректно завершается с таймаутом 5 секунд.

## CI/CD (GitHub Actions)

Ручной запуск CI доступен в разделе **Actions → Go CI** и **Publish Docker Image**. Для публикации Docker-образа:
1. В Actions выберите **Publish Docker Image** и нажмите **Run workflow**.
2. Укажите параметр `version` (например, `v1.0.0`), чтобы задать тег образа.
3. Нажмите **Run workflow** — GitHub Actions соберёт и опубликует образ в GHCR.

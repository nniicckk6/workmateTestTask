# WorkmateTestProject

HTTP-сервис для управления I/O-зависимыми задачами (время обработки 1–5 минут) с хранением в памяти.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
- [Installation](#installation)
- [Build & Run](#build--run)
- [Usage](#usage)
- [Logging & Graceful Shutdown](#logging--graceful-shutdown)
- [Contributing](#contributing)
- [License](#license)

## Prerequisites

- Go 1.24 или выше

## Configuration

Сервис по умолчанию слушает порт `8080`. Можно изменить через переменную окружения:

```bash
export PORT=9090
```

## Installation

```bash
git clone сейчас вставлю
cd workmateTestProject
go mod download
```

## Build & Run

```bash
go build -o workmateTestProject main.go
PORT=8080 ./workmateTestProject
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

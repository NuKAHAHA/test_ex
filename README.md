

````markdown
Task Management REST API

Простой REST API на Go для управления задачами.  
Хранение данных — в памяти, асинхронное логирование действий — через канал и горутину.  
Проект следует принципам чистой архитектуры и использует только стандартные библиотеки Go (без внешних зависимостей).

---

Возможности

Эндпоинты
| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/tasks` | Получить список задач (можно фильтровать по статусу, например `?status=pending`) |
| GET | `/tasks/{id}` | Получить задачу по ID |
| POST | `/tasks` | Создать новую задачу |

---

Структура задачи
```json
{
  "id": "abc123",
  "title": "Sample Task",
  "description": "This is a test",
  "status": "pending",
  "createdAt": "2025-08-15T12:00:00Z",
  "updatedAt": "2025-08-15T12:00:00Z"
}
````

**Статусы**:

* `pending`
* `in_progress`
* `done`
* `canceled`

---

### Основные особенности

* **Асинхронное логирование** в формате JSON (`stdout`)
* **Потокобезопасное** хранение (`sync.RWMutex`)
* **Graceful Shutdown** при `SIGINT` / `SIGTERM` (таймаут 5 секунд)
* **Чистая архитектура** с разделением на слои:

  * `Model` — структуры данных
  * `Repo` — доступ к данным
  * `Service` — бизнес-логика
  * `Handler` — HTTP-обработчики
  * `Utils` — вспомогательные функции

---

## Установка

```bash
git clone <https://github.com/NuKAHAHA/test_ex>
cd awesomeProject
```

Убедитесь, что Go установлен (**v1.22+**).

## Запуск сервера

```bash
go run ./cmd/main.go
```

По умолчанию сервер доступен на:

```
http://localhost:8080
```

В логах появится:

```
Server started on :8080
```

---

## Примеры использования API

### Создать задачу

```bash
curl -X POST http://localhost:8080/tasks \
-H "Content-Type: application/json" \
-d '{"title": "Sample Task", "description": "Test", "status": "pending"}'
```

**Ответ**: `201 Created` с JSON задачи.

---

### Получить список задач

```bash
curl http://localhost:8080/tasks?status=pending
```

---

### Получить задачу по ID

```bash
curl http://localhost:8080/tasks/<task-id>
```

---

## Пример логов

```json
{"time":"2025-08-15T12:34:56Z","action":"task_created","task_id":"abc123","meta":{"title":"Sample Task","status":"pending"}}
{"time":"2025-08-15T12:35:00Z","action":"http_request_start","meta":{"method":"GET","path":"/tasks"}}
```

---

## Архитектура

* **Логирование**: Буферизированный канал (2048 сообщений) + отдельная горутина
* **Потокобезопасность**: `RWMutex`
* **Генерация ID**: `crypto/rand` + timestamp
* **Middleware**: перехват паник с возвратом JSON ошибок
* **Ошибка → HTTP код**: автоматическое сопоставление

---

## Тестирование

Ручное тестирование:

```bash
curl http://localhost:8080/tasks
```

Автотесты — через `net/http/httptest` (юнит-тесты не включены, но легко добавить).

---

## Ограничения

* Данные **теряются при перезапуске**
* Нет аутентификации и rate limiting
* Логирование только в `stdout` (нет ротации файлов)
* Минимальная валидация

---

## Время разработки

\~5 часа (проектирование, кодинг, тестирование)

---

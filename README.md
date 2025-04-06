# S3 Compatible Storage API

Простой сервис для работы с файловым хранилищем, совместимым с S3 API (используется MinIO), построенный на трехслойной архитектуре с использованием веб-фреймворка Fiber.

## 🚀 Быстрый старт

1. Убедитесь, что установлены Docker и Docker Compose
2. Клонируйте репозиторий
3. Запустите сервис:
   ```bash
   docker-compose up -d
   ```
4. Сервис будет доступен по адресу: `http://localhost:8080`

## 📚 API Endpoints

| Метод | Endpoint                  | Описание                          |
|-------|---------------------------|-----------------------------------|
| GET   | `/storage`                | Получить список всех файлов       |
| POST  | `/storage`                | Загрузить текстовый файл          |
| DELETE| `/storage/{name}`         | Удалить файл                      |
| GET   | `/storage/{name}/content` | Получить содержимое файла         |

Base URL: `http://localhost:8080/api/v1`

## 🛠 Технологии

- **Язык**: Go
- **Фреймворк**: [Fiber](https://gofiber.io/)
- **Хранилище**: MinIO (S3-совместимое)
- **Документация**: Swagger/OpenAPI

## 📂 Структура проекта

```
docker-compose-training/
├── cmd/                # Точка входа
│   └── main.go
├── config/             # Конфигурация
│   ├── config.go
│   └── local.yaml
├── docs/               # Документация Swagger
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/           # Основная логика
│   ├── domain/         # Доменные модели и сервисы
│   ├── repository/     # Работа с хранилищем
│   ├── rest/           # HTTP handlers
│   └── server/         # Настройка сервера и MinIO
├── .gitignore
├── docker-compose.yml  # Конфигурация Docker
├── Dockerfile
├── go.mod              # Зависимости
└── README.md
```

## ⚙️ Конфигурация

Настройки в `config/local.yaml`:

```yaml
env: "local"
minio:
   host: "minio"
   port: "9000"
   access_key: "admin"
   secret_access_key: "qwerty_password"
   ssl_mode: false
   base_bucket: "files.bucket"
   base_path: "./tmp/"

server:
   address: ":8080"
   timeout: 4s
   idle_timeout: 60s
```

## 📌 Примеры использования

### Получить список файлов
```bash
curl -X GET http://localhost:8080/api/v1/storage
```

### Загрузить файл
```bash
curl -X POST -F "file=@notes.txt" http://localhost:8080/api/v1/storage
```

### Получить содержимое файла
```bash
curl -X GET http://localhost:8080/api/v1/storage/notes.txt/content
```

### Удалить файл
```bash
curl -X DELETE http://localhost:8080/api/v1/storage/notes.txt
```

## 🔍 Документация

Swagger UI доступен после запуска:
- Интерфейс: `http://localhost:8080/swagger/` (если настроен)

## 🛠 Локальная разработка

1. Соберите контейнеры:
   ```bash
   docker compose build 
   ```
2. Запустите сервис:
   ```bash
   docker compose up
   ```
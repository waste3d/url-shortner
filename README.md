# 🔗 URL Shortener — Сервис сокращения ссылок

Полноценный веб-сервис для сокращения ссылок с поддержкой статистики, QR-кодов, кастомных ссылок и срока жизни. Включает серверную часть на Go (Gin, GORM, PostgreSQL) и клиентскую часть на React + TailwindCSS.

---

## 🧰 Стек технологий

### Backend:
- [Go](https://golang.org/)
- [Gin](https://github.com/gin-gonic/gin) — веб-фреймворк
- [GORM](https://gorm.io/) — ORM для Go
- [PostgreSQL](https://www.postgresql.org/) — база данных

### Frontend:
- [React](https://reactjs.org/)
- [TailwindCSS](https://tailwindcss.com/) — CSS-фреймворк

---

## 🚀 Возможности

- Сокращение длинных ссылок
- Установка срока жизни (TTL) для каждой ссылки
- Создание кастмоных ссылок
- Статистика переходов: IP, User-Agent
- Подсчет общего количества переходов
- Генерация QR-кодов
- Настройка цвета QR-кода
- Кеширование QR-кодов
- Скачивание QR-кода

---

## 📦 Установка и запуск

### 1. Клонирование проекта

```bash
git clone https://github.com/waste3d/url-shortner.git
cd your-repo
````

---

### 2. Backend

#### 📄 Настройка `.env`

Создай `.env` в папке `backend/`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=shortener
SERVER_PORT=3000
```

#### ▶️ Запуск сервера

```bash
cd backend
go run main.go
```

---

### 3. Frontend

#### Установка и запуск

```bash
cd frontend
npm install
npm run dev
```

---

## 📡 API Эндпоинты

### 🔗 Ссылки

| Метод | Endpoint              | Описание                            |
| ----- | --------------------- | ----------------------------------- |
| GET   | `/links`              | Получить все ссылки                 |
| POST  | `/links`              | Создать новую короткую ссылку       |
| GET   | `/:shortened`         | Редирект по сокращённой ссылке      |
| GET   | `/links/:id`          | Получить информацию о ссылке        |
| GET   | `/links/:id/visitors` | Получить список переходов по ссылке |

#### Пример POST запроса:

```json
{
  "clicks":1,
  "created_at":"09.05.25 18:05",
  "expire_at":"31.12.25 01:44",
  "id":1,
  "original":"https://waste3d.su",
  "shortened":"flink"
}
```

---

### 🧾 QR-коды

| Метод | Endpoint                        | Описание                                    |
| ----- | ------------------------------- | ------------------------------------------- |
| GET   | `/qr/view?url=<url>&color=#000` | Отобразить QR-код (цвет задаётся hex-кодом) |
| GET   | `/qr/download?url=<url>`        | Скачать QR-код                              |

---

## 📌 Примечания

* Поддержка TTL — ссылки автоматически перестают работать по истечении срока.
* Кеш QR-кодов ускоряет повторную генерацию.
* QR можно стилизовать с помощью параметра `color` (`?color=#FF0000`).

---

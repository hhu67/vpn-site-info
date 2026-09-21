# VPN Site Info

Веб-приложение для управления VPN ссылками с аутентификацией на основе JWT токенов.

## Возможности

- 🔐 Аутентификация с использованием bcrypt и JWT
- 📝 Создание первичного пароля с подтверждением
- 🔄 Смена пароля с верификацией старого пароля
- ➕ Добавление, редактирование и удаление VPN ссылок
- 📋 Просмотр всех сохраненных ссылок
- 📋 Копирование ссылок в буфер обмена одной кнопкой
- 🎨 Современный адаптивный интерфейс на React
- 🔒 Защита всех API эндпоинтов с помощью JWT middleware

## Технологии

### Backend
- Go 1.23
- PostgreSQL с драйвером pgx/v5
- JWT для аутентификации (действителен 24 часа)
- bcrypt для хеширования паролей
- godotenv для управления переменными окружения

### Frontend
- React 18
- Vite
- Современный CSS с градиентами и анимациями

## Структура проекта

```
vpn-site-info/
├── cmd/
│   └── server/
│       └── main.go           # Точка входа приложения
├── internal/
│   ├── config/
│   │   └── config.go         # Загрузка конфигурации
│   ├── database/
│   │   └── database.go       # Подключение к БД и схема
│   ├── handlers/
│   │   ├── handler.go        # Базовый обработчик
│   │   ├── auth.go           # Аутентификация
│   │   └── vpn.go            # CRUD для VPN ссылок
│   └── middleware/
│       └── auth.go           # JWT middleware
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── CreatePassword.jsx
│   │   │   ├── Login.jsx
│   │   │   └── Dashboard.jsx
│   │   ├── App.jsx
│   │   ├── main.jsx
│   │   └── index.css
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
├── .env.example              # Пример конфигурации
├── go.mod
└── README.md
```

## Установка и запуск

### Предварительные требования

- Go 1.23 или выше
- PostgreSQL
- Node.js и npm

### 1. Настройка базы данных

Создайте базу данных PostgreSQL:

```bash
createdb vpndb
```

### 2. Настройка переменных окружения

Скопируйте файл с примером и отредактируйте его:

```bash
cp .env.example .env
```

Отредактируйте `.env`:

```env
PSQL=postgresql://user:password@localhost:5432/vpndb?sslmode=disable
JWT_SECRET=ваш-секретный-ключ-здесь
PORT=8080
```

### 3. Установка зависимостей Go

```bash
go mod download
```

### 4. Установка зависимостей фронтенда

```bash
cd frontend
npm install
```

### 5. Сборка фронтенда

```bash
npm run build
cd ..
```

### 6. Запуск сервера

```bash
go run cmd/server/main.go
```

Сервер будет доступен по адресу `http://localhost:8080`

## API Endpoints

### Публичные эндпоинты

- `GET /api/check-password-exists` - Проверка существования пароля
- `POST /api/create-password` - Создание первичного пароля
  ```json
  {
    "password": "ваш-пароль",
    "password_confirm": "ваш-пароль"
  }
  ```
- `POST /api/login` - Вход в систему
  ```json
  {
    "password": "ваш-пароль"
  }
  ```

### Защищенные эндпоинты (требуют JWT в cookies)

- `GET /api/list/vpn` - Получение списка всех ссылок
- `POST /api/insert/vpn` - Добавление новой ссылки
  ```json
  {
    "link": "vpn://ссылка"
  }
  ```
- `PUT /api/update/vpn` - Обновление ссылки
  ```json
  {
    "id": 1,
    "link": "vpn://новая-ссылка"
  }
  ```
- `DELETE /api/delete/vpn?id=1` - Удаление ссылки
- `POST /api/change-password` - Смена пароля (требует старый пароль)
  ```json
  {
    "old_password": "старый-пароль",
    "new_password": "новый-пароль",
    "password_confirm": "новый-пароль"
  }
  ```

## Схема базы данных

### Таблица `auth`
```sql
CREATE TABLE auth (
    id SERIAL PRIMARY KEY,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Таблица `vpn_links`
```sql
CREATE TABLE vpn_links (
    id SERIAL PRIMARY KEY,
    link VARCHAR(1024) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Разработка

### Запуск фронтенда в режиме разработки

```bash
cd frontend
npm run dev
```

Vite запустит dev-сервер с hot reload и проксированием API запросов на бэкенд.

### Сборка для продакшена

```bash
cd frontend
npm run build
cd ..
go build -o vpn-site-info cmd/server/main.go
```

## Безопасность

- Пароли хешируются с использованием bcrypt
- JWT токены действительны 24 часа
- Токены хранятся в HttpOnly cookies для защиты от XSS
- Все API эндпоинты управления данными защищены JWT middleware
- При первом запуске требуется создание пароля с подтверждением

## Логика работы

1. При первом запуске пользователь попадает на страницу создания пароля
2. После создания пароля автоматически выдается JWT токен
3. При повторном входе пользователь вводит пароль и получает JWT токен
4. JWT токен проверяется при каждом обращении к защищенным эндпоинтам
5. Если токен недействителен или отсутствует, пользователь перенаправляется на страницу входа
6. На панели управления пользователь может:
   - Просматривать все добавленные VPN ссылки
   - Копировать ссылку в буфер обмена одной кнопкой
   - Редактировать существующие ссылки
   - Удалять ссылки
   - Менять пароль (с верификацией старого пароля)
   - Выходить из системы

## Лицензия

MIT

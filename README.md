# Go Chat

Go Chat is a REST API backend for a chat and browser notification system. It is built with Go, Gin, GORM, and MySQL, and includes authentication, user management, private chat rooms, chat messages, browser Web Push subscriptions, notification delivery, and a small Gemini AI prompt testing endpoint.

## Features

- User registration, login, logout, and authenticated profile access
- Bearer token authentication using personal access tokens
- Role-based route protection for admin/root users
- User profile and PIN update endpoints
- Private chat room listing and paginated message retrieval
- Message read tracking and unread message counts
- Browser push subscription, unsubscribe, and notification sending
- VAPID key generator for Web Push setup
- Gemini-powered prompt testing endpoint
- Layered project structure with domain, app/service, port, adapter, and HTTP packages

## Tech Stack

- Go 1.25.8
- Gin HTTP framework
- GORM ORM
- MySQL
- Web Push with VAPID keys
- Google Gemini API

## Project Structure

```text
cmd/
  main.go              # API application entrypoint
  routes.go            # Route registration
  vapidkeygen/         # Helper command to generate Web Push VAPID keys

internal/
  adapter/
    db/                # GORM repository implementations
    dto/               # Request/response DTOs
    http/              # Gin HTTP handlers and middleware
  app/                 # Business logic services
  domain/              # Database/domain models
  port/                # Repository interfaces
  shared/              # Shared helper utilities
```

## Requirements

Before cloning and running this project, make sure you have:

- Git
- Go 1.25.8 or newer
- MySQL running locally or remotely
- A Gemini API key
- A database schema compatible with the models in `internal/domain`

> Note: this repository currently does not include migration files. Create/import the required MySQL tables before running the API.

## Clone and Run Tutorial

1. Clone the repository:

```bash
git clone https://github.com/rosyanxone/go-chat.git
```

2. Move into the project directory:

```bash
cd go-chat
```

3. Download Go dependencies:

```bash
go mod download
```

4. Create your environment file:

```bash
cp .env.example .env
```

On Windows PowerShell, use:

```powershell
Copy-Item .env.example .env
```

5. Configure `.env`:

```env
APP_PORT="8000"

DB_USER="root"
DB_PASS=""
DB_HOST="localhost"
DB_PORT="3306"
DB_NAME="chat"

GEMINI_API_KEY="your-gemini-api-key"

CHAT_DOMAIN="http://localhost:3000"

VAPID_PUBLIC_KEY=""
VAPID_PRIVATE_KEY=""
VAPID_SUBJECT="mailto:user@domain.com"
```

6. Generate VAPID keys for browser push notifications:

```bash
go run ./cmd/vapidkeygen
```

Copy the generated `VAPID_PUBLIC_KEY` and `VAPID_PRIVATE_KEY` values into your `.env` file.

7. Make sure the MySQL database exists:

```sql
CREATE DATABASE chat;
```

Then import or create the required tables for users, roles, chats, chat rooms, chat messages, personal access tokens, employees, and push subscriptions.

8. Run the API:

```bash
go run ./cmd
```

9. Test the health endpoint:

```bash
curl http://localhost:8000/health
```

Expected response:

```json
{
  "status": "ok"
}
```

## Main API Routes

Public routes:

- `GET /health`
- `POST /api/register`
- `POST /api/login`
- `POST /api/validate/phone-number`
- `GET /api/users`
- `GET /api/user/email?email=user@example.com`
- `GET /api/web/vapid-public-key`
- `POST /api/prompt`
- `POST /api/test`

Authenticated routes:

- `GET /api/user`
- `POST /api/logout`
- `POST /api/user/update`
- `POST /api/user/update/pin`
- `GET /api/chat/rooms`
- `POST /api/chat/messages`
- `POST /api/web/subscribe`
- `POST /api/web/unsubscribe`
- `POST /api/notify`

Authenticated requests should include this header:

```http
Authorization: Bearer <token>
```

## Notes

- The API runs on the port defined by `APP_PORT`.
- CORS is configured for local frontend ports and the production domains listed in `cmd/main.go`.
- `GEMINI_API_KEY` is required at startup because the Gemini client is initialized when the app boots.
- Push notifications require valid VAPID keys and frontend service worker support.

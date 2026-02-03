# GitHub Resume Builder - Backend

Production-grade backend for automatically generating resumes from GitHub profiles.

## Architecture

```
Layered Architecture:
├── Handlers (HTTP layer)
├── Services (Business logic)
├── Repositories (Data access)
└── Clients (External APIs)
```

## Tech Stack

- **Language**: Go 1.22
- **Router**: chi
- **Database**: PostgreSQL
- **Auth**: GitHub OAuth2
- **Encryption**: AES-256-GCM

## Setup

### Prerequisites

- Go 1.22+ OR Docker
- PostgreSQL 14+ (or use Docker Compose)
- GitHub OAuth App

### Option 1: Docker (Recommended)

1. Create `.env` file:
```bash
cp .env.example .env
# Edit .env with your GitHub OAuth credentials
```

2. Start all services:
```bash
make docker-up
```

3. Run migrations:
```bash
docker-compose exec app migrate -path migrations -database "postgres://postgres:postgres@postgres/resume_builder?sslmode=disable" up
```

4. Access at `http://localhost:8080`

### Option 2: Local Development

### 1. Create GitHub OAuth App

1. Go to GitHub Settings → Developer settings → OAuth Apps
2. Create new OAuth App
3. Set Authorization callback URL: `http://localhost:8080/auth/callback`
4. Note your Client ID and Client Secret

### 2. Database Setup

```bash
createdb resume_builder
```

### 3. Environment Variables

```bash
cp .env.example .env
# Edit .env with your values
```

Generate encryption key:
```bash
openssl rand -base64 32
```

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run Migrations

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -path migrations -database "postgres://localhost/resume_builder?sslmode=disable" up
```

### 6. Run Server

```bash
go run cmd/api/main.go
```

## API Endpoints

### Authentication

**Login**
```
GET /auth/login
```
Redirects to GitHub OAuth.

**Callback**
```
GET /auth/callback?code=xxx&state=xxx
```
Handles OAuth callback.

### Resumes (Protected)

**Generate Resume**
```
POST /resumes/generate
Authorization: Bearer <token>

{
  "target_role": "Backend Engineer"
}
```

**List Resumes**
```
GET /resumes
Authorization: Bearer <token>
```

**Get Resume**
```
GET /resumes/{id}
Authorization: Bearer <token>
```

**Update Resume**
```
PUT /resumes/{id}
Authorization: Bearer <token>

{
  "title": "Updated Resume",
  "target_role": "Senior Backend Engineer",
  "summary": "...",
  "projects": [...],
  "skills": [...]
}
```

**Delete Resume**
```
DELETE /resumes/{id}
Authorization: Bearer <token>
```

## Repository Ranking Algorithm

Repositories are scored based on:

- **Stars** (30%): Logarithmic scale
- **Recency** (25%): Last commit date
- **Language** (20%): Has primary language
- **Topics** (15%): Number of topics
- **Description** (10%): Has description

Forks are excluded from ranking.

## Security Features

- AES-256-GCM token encryption
- JWT authentication with 24-hour expiry
- CSRF protection via state parameter
- Secure cookie handling
- Context-aware timeouts
- Input validation
- Rate limiting (100 req/min public, 50 req/min authenticated)

## Performance Features

- Redis caching for GitHub API responses (1-hour TTL)
- LLM-powered resume summaries (OpenAI GPT-3.5-turbo)
- Automatic fallback to rule-based summaries
- Connection pooling for database
- Graceful shutdown
- Request timeouts

## Development

### Project Structure

```
resume-builder/
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── handler/         # HTTP handlers
│   ├── service/         # Business logic
│   ├── repository/      # Database operations
│   ├── client/          # External API clients
│   ├── model/           # Domain models
│   └── crypto/          # Encryption utilities
└── migrations/          # Database migrations
```

### Adding New Features

1. Define models in `internal/model/`
2. Create repository methods in `internal/repository/`
3. Implement business logic in `internal/service/`
4. Add HTTP handlers in `internal/handler/`
5. Register routes in `cmd/api/main.go`

## TODO

- [x] Implement JWT/session-based authentication
- [x] Add Redis caching for GitHub API responses
- [x] Implement rate limiting
- [x] Add comprehensive logging
- [ ] Add unit and integration tests
- [ ] Add PDF export functionality
- [ ] Implement webhook for auto-refresh
- [ ] Add metrics and monitoring

## License

MIT

## Deployment

See [RENDER.md](RENDER.md) for Render deployment instructions.

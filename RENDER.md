# Render Deployment Guide

## Prerequisites

1. **GitHub Repository**: Push your code to GitHub
2. **Render Account**: Sign up at https://render.com
3. **PostgreSQL Database**: Use Render PostgreSQL or external (Aiven, etc.)

## Deployment Steps

### 1. Create PostgreSQL Database (if needed)

1. Go to Render Dashboard → New → PostgreSQL
2. Choose a name: `resume-builder-db`
3. Select free tier or paid plan
4. Note the connection details

### 2. Deploy Application

1. Go to Render Dashboard → New → Web Service
2. Connect your GitHub repository
3. Configure:
   - **Name**: `resume-builder`
   - **Environment**: `Go`
   - **Build Command**: `go build -o bin/main cmd/api/main.go`
   - **Start Command**: `./bin/main`

### 3. Set Environment Variables

Add these in Render Dashboard → Environment:

```
PORT=8080
ENV=production

# Database (from Render PostgreSQL or external)
DB_HOST=<your-db-host>
DB_PORT=5432
DB_USER=<your-db-user>
DB_PASSWORD=<your-db-password>
DB_NAME=resume_builder
DB_SSLMODE=require

# GitHub OAuth
GITHUB_CLIENT_ID=<your-github-client-id>
GITHUB_CLIENT_SECRET=<your-github-client-secret>
GITHUB_REDIRECT_URL=https://your-app.onrender.com/auth/callback

# Security (generate new keys for production)
ENCRYPTION_KEY=<32-character-key>
JWT_SECRET=<32-character-key>

# Optional: Redis (use Render Redis or disable)
REDIS_ENABLED=false

# Optional: OpenAI
OPENAI_ENABLED=false
OPENAI_API_KEY=<your-openai-key>
```

### 4. Run Migrations

After first deploy, run migrations via Render Shell:

```bash
migrate -path migrations -database "postgres://<connection-string>" up
```

Or use Render's PostgreSQL connection string from environment.

### 5. Update GitHub OAuth

1. Go to GitHub Settings → Developer settings → OAuth Apps
2. Update Authorization callback URL to:
   ```
   https://your-app.onrender.com/auth/callback
   ```

## Auto-Deploy

Render automatically deploys on every push to your main branch.

## Health Check

Render will use: `GET /health`

## Scaling

- Free tier: Spins down after inactivity
- Paid tier: Always on, auto-scaling available

## Monitoring

- View logs in Render Dashboard
- Structured JSON logs for easy parsing
- Set up alerts for errors

## Cost Optimization

- **Free tier**: $0/month (with limitations)
- **Starter**: $7/month (always on)
- **PostgreSQL**: Free tier available
- **Redis**: Optional, disable if not needed
- **OpenAI**: Pay per use, disable if not needed

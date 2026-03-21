# Envoy | PostgreSQL Schema Management Tool

<p align="center">
  <img src="readme/logo.png" alt="Envoy Logo" width="160px"/>
  <br>
  <b>The Unified Database Control Plane for Startups.</b>
  <br>
  <i>Stop juggling terminals. Handle multiple environments with speed and compliance.</i>
</p>

---

## The Problem

Managing database migrations across **Staging, QA, and Production** is a fragmented mess:

* **Switching .env files** manually is dangerous and error-prone
* **Forgetting GRANT permissions** in pgAdmin causes production downtime  
* **Audit logs** are non-existent, scattered, or hard to verify
* **Multiple tools** required for simple database operations

### Example Situation (Every Friday at 10am)
1. Get production database credentials
2. Replace environment variable in .env file
3. Run migration command
4. Verify migration was successful
5. Try changes in production
6. **Realize you forgot to give permissions to new tables**
7. Go back to pgAdmin/psql to fix permissions
8. Try again in production
9. Repeat next week across all environments

## The Envoy Solution

Envoy sits between your ORM (Prisma, Gorm, Drizzle) and your infrastructure to provide:

* **One Interface:** Switch environments without touching config files
* **Permissions Audit:** Automatic GRANT verification per environment
* **Audit History:** GitHub-style log of every query executed
* **Security First:** Encrypted secrets and audit trails

---

### Who is this for?

* Founding Engineers building scalable DB processes from Day 1.
* DevOps standardizing how teams interact with production data.
* Solo Devs tired of manual migration PTSD.

---

## Quick Start

### Prerequisites
- Docker and Docker Compose installed
- Git

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-org/envoy.git
   cd envoy
   ```

2. **Start the development environment**
   ```bash
   docker-compose up -d
   ```

3. **Access the application**
   - **Web Interface**: http://localhost:5173 (Frontend)
   - **API Server**: http://localhost:8080 (Backend)
   - **Database Admin**: http://localhost:5050 (pgAdmin)

That's it! Envoy is now running with both frontend and backend services.

---

## Configuration

Copy the environment file and update with your values:

```bash
cp server/.env.example server/.env
```

### Required Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | Envoy's internal database (SQLite) | `file:./data/envoy.db` |
| `JWT_SECRET` | Authentication secret | `your-secret-key` |
| `ENCRYPTION_KEY` | 32-character encryption key | `your-encryption-key-32-chars` |
| `CHECKSUM_KEY` | 32-character checksum key | `your-checksum-key-32-chars` |
| `PORT` | API server port | `8080` |
| `APP_ENV` | Environment mode | `development` |

### Example `.env` File
```env
DATABASE_URL=file:./data/envoy.db
JWT_SECRET=your-jwt-secret-64-chars-ultra-secure-random-key
ENCRYPTION_KEY=your-encryption-key-32-chars-123
CHECKSUM_KEY=your-checksum-key-32-chars-12345
PORT=8080
APP_ENV=development
```

---

## How It Works

### Architecture
```
[React Frontend] → [Go API] → [Target Databases]
                    ↓
                [SQLite Storage]
```

### Key Features
- **Environment Management**: Add/switch between dev, staging, prod databases
- **Migration Editor**: Write SQL with syntax highlighting
- **Schema Preview**: See changes before applying
- **Permission Auditing**: Automatic GRANT verification
- **Audit Trail**: Complete history of all database operations

---

## Development

### Local Development

#### Backend (Go)
```bash
cd server
go mod download
go run cmd/main.go
```

#### Frontend (React)
```bash
cd client
npm install
npm run dev
```

### Project Structure
```
envoy/
├── server/          # Go backend API
├── client/          # React frontend
├── docker-compose.yml  # Development environment
└── readme.md        # This file
```

---

## Troubleshooting

### Common Issues

**Database connection failed**
```bash
docker-compose ps postgres-dev
docker-compose logs postgres-dev
```

**API not starting**
```bash
cat server/.env
docker-compose logs envoy-api
```

**Port conflicts**
```bash
netstat -tulpn | grep :8080
netstat -tulpn | grep :5432
```

---

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

<p align="center">
  <b>Built with ❤️ by Francisco Luna, for developers who value speed and security.</b>
  <br>
  <i>Stop managing databases. Start shipping features.</i>
</p>

---

<p align="center">
  <b>Open for Collaborations</b>
  <br>
  <i>Need help with your database infrastructure? I'm available for Node.js, Golang and Postgres backend and infrastructure consulting. Want to discuss your project?</i>
  <br>
  <br>
  <a href="mailto:franciscolunadev@gmail.com">Let's Talk About Your Environment</a>
</p>
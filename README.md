# GoBan

A minimalist, self-hosted Kanban board built with Go (Fiber) and React. Single-binary deploy, SQLite/Postgres support, and clean UI.

## Features

- Single binary deployment with embedded frontend
- SQLite (default) or PostgreSQL database support
- JWT authentication with HTTPOnly cookies
- Drag & drop card management
- Dark/Light mode toggle
- User isolation - each user sees only their boards
- Default columns ("To Do", "In Progress", "Done") on new boards
- Card priority levels (low, medium, high)

## Tech Stack

**Backend:**
- Go 1.23+ with Fiber framework
- GORM for database operations
- JWT for authentication

**Frontend:**
- React 19 with TypeScript
- Vite build tool
- Tailwind CSS 4
- @dnd-kit for drag and drop
- Lucide icons

## Quick Start

### Using Docker (Recommended)

```bash
# Build and run
docker build -t goban:latest .
docker run -p 8080:8080 -v goban-data:/app/data -e JWT_SECRET=your-secret-key goban:latest

# Or using docker-compose
docker-compose up -d
```

### Development

```bash
# Install dependencies
cd web && bun install && cd ..
go mod download

# Run backend (includes embedded frontend if built)
go run ./cmd/server

# Run frontend dev server (separate terminal)
cd web && bun run dev
```

### Build from Source

```bash
# Build frontend
cd web && bun install && bun run build && cd ..

# Build backend with embedded frontend
go build -o goban ./cmd/server

# Run
./goban
```

## Configuration

Environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_DRIVER` | Database driver (`sqlite` or `postgres`) | `sqlite` |
| `DATABASE_URL` | Database connection string | `./goban.db` |
| `JWT_SECRET` | Secret key for JWT tokens | `default-secret-change-me` |

Example `.env` file:

```env
PORT=8080
DB_DRIVER=sqlite
DATABASE_URL=./goban.db
JWT_SECRET=your-super-secret-key
```

### PostgreSQL Configuration

```env
DB_DRIVER=postgres
DATABASE_URL=postgres://user:password@localhost:5432/goban?sslmode=disable
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Create account
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/auth/me` - Get current user

### Boards
- `GET /api/v1/boards` - List boards
- `POST /api/v1/boards` - Create board
- `GET /api/v1/boards/:id` - Get board with columns/cards
- `PUT /api/v1/boards/:id` - Update board
- `DELETE /api/v1/boards/:id` - Delete board

### Columns
- `POST /api/v1/boards/:boardId/columns` - Create column
- `PUT /api/v1/columns/:id` - Update column
- `DELETE /api/v1/columns/:id` - Delete column
- `PUT /api/v1/columns/reorder` - Reorder columns

### Cards
- `POST /api/v1/columns/:columnId/cards` - Create card
- `GET /api/v1/cards/:id` - Get card
- `PUT /api/v1/cards/:id` - Update card
- `DELETE /api/v1/cards/:id` - Delete card
- `PUT /api/v1/cards/:id/move` - Move card to column
- `PUT /api/v1/cards/reorder` - Reorder cards

## Project Structure

```
goban/
├── cmd/server/          # Application entrypoint
├── internal/
│   ├── config/          # Configuration
│   ├── database/        # Database connection
│   ├── dto/             # Data transfer objects
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Auth middleware
│   ├── models/          # GORM models
│   ├── repository/      # Data access layer
│   ├── router/          # Route definitions
│   ├── services/        # Business logic
│   └── utils/           # Utilities
├── web/                 # React frontend
│   └── src/
│       ├── components/  # UI components
│       ├── context/     # React contexts
│       ├── lib/         # Utilities
│       ├── pages/       # Page components
│       └── types/       # TypeScript types
├── embed.go             # Go embed directive
├── Dockerfile           # Multi-stage build
└── docker-compose.yml   # Docker Compose config
```

## License

MIT

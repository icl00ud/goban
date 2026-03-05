# Goban

Your kanban board, your server, your rules.

Goban is a self-hosted kanban board that lives wherever you want — a VPS, a Raspberry Pi, a home server. Your tasks stay private, no subscriptions, no limits.

Beyond the web interface, Goban ships with a **CLI**, an **MCP server** (so any AI assistant can manage your boards), and a **Telegram bot**.

---

## What you can do

- Organize work across multiple boards with columns and cards
- Drag and drop cards between columns
- Set priorities (low, medium, high) on cards
- Switch between dark and light mode
- Manage everything from the terminal with the CLI
- Let your AI assistant (Claude, etc.) create and move cards for you via MCP
- Control your boards from Telegram on your phone

---

## Get started

The fastest way is Docker:

```bash
docker run -p 8080:8080 \
  -v goban-data:/app/data \
  -e JWT_SECRET=change-me \
  ghcr.io/icl00ud/goban:latest
```

Open `http://localhost:8080`, create an account and you're in.

Or with Docker Compose — create a `docker-compose.yml`:

```yaml
services:
  goban:
    image: ghcr.io/icl00ud/goban:latest
    ports:
      - "8080:8080"
    volumes:
      - goban-data:/app/data
    environment:
      JWT_SECRET: change-me

volumes:
  goban-data:
```

```bash
docker-compose up -d
```

---

## CLI

Manage your boards from the terminal.

**Install:**

```bash
# macOS
brew tap icl00ud/goban https://github.com/icl00ud/goban
brew install goban

# npm / bun / yarn / pnpm
npm install -g goban-cli
bun add -g goban-cli

# Linux — download from GitHub Releases
# https://github.com/icl00ud/goban/releases/latest
# .deb, .rpm and .apk available
```

**Usage:**

```bash
# Authenticate (connects to your self-hosted instance)
goban login --server https://your-goban.com

# Boards
goban boards list
goban boards create --name "Sprint 4" --color "#8b5cf6"
goban boards get 1

# Cards
goban cards list --board 1
goban cards create --column 2 --title "Fix login bug" --priority high
goban cards move 15 --to-column 3

# Save the token to use later without --server every time
# The config lives at ~/.config/goban/config.json
```

---

## AI integration (MCP)

Goban exposes an MCP server so any AI assistant that supports the Model Context Protocol can read and manage your boards.

**Setup with Claude Desktop** — add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "goban": {
      "command": "goban",
      "args": ["mcp"]
    }
  }
}
```

Make sure you've run `goban login` first so the token is saved. Then restart Claude Desktop and ask it things like:

> "Show me all my boards"
> "Create a card called 'Write tests' in the To Do column of Sprint 4"
> "Move card 12 to Done"

---

## Telegram bot

Control your boards from Telegram.

```bash
# Set your bot token (from @BotFather) and start
TELEGRAM_BOT_TOKEN=your-token goban bot
```

**Available commands:**

```
/auth <token>          Authenticate (get token with: goban login --print-token)
/boards                List your boards
/board <id>            See a board with all columns and cards
/add <column> <title>  Create a card
/done <card>           Move a card to the last column
/move <card> <column>  Move a card to any column
/cards <board>         List all cards in a board
/card <id>             See card details
/newboard <name>       Create a board
/newcol <board> <name> Create a column
```

---

## Configuration

| Variable | What it does | Default |
|----------|-------------|---------|
| `PORT` | Port the server listens on | `8080` |
| `JWT_SECRET` | Secret for signing sessions — **change this** | `default-secret-change-me` |
| `DATABASE_URL` | Path to the SQLite file, or a Postgres connection string | `./goban.db` |
| `DB_DRIVER` | `sqlite` or `postgres` | `sqlite` |

**PostgreSQL:**

```env
DB_DRIVER=postgres
DATABASE_URL=postgres://user:password@localhost:5432/goban?sslmode=disable
```

---

## License

MIT

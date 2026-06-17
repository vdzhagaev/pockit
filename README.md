# pockit

> A **pocket toolkit bot** for Telegram & Slack — small, genuinely useful tools you can call from any chat, built on a pluggable `Tool` architecture. Designed to run as a hosted, multi-tenant service.

🚧 **Status: in progress / pivoting.** This repo started life as a Go URL-shortener service and is being evolved into `pockit`. The URL-shortener below is the **first tool** and the working core; the messenger adapters and additional tools are on the roadmap.

## The idea

Instead of one app per tiny task, `pockit` is a single backend that exposes a set of small tools through chat. Call a command in any Telegram or Slack chat where the bot lives and get an instant result:

```
/short https://example.com/very/long/link   →  https://pk.it/x7Qa2
/remind 2h ping the deploy channel
/share  (attach a file)                      →  temporary share link
```

Because every tool is backed by **one shared service**, state is synced across messengers automatically — shorten a link in Telegram, see it in Slack.

### Why a pluggable architecture

Each tool implements a common `Tool` interface (command, input parsing, handler) and is registered into a router. Messenger integrations (Telegram, Slack) are **thin adapters** that translate a chat command into a tool invocation. Adding a tool means writing one package — not touching the transport layer. This seam is intentionally the same shape as a general workflow/automation engine, so the design scales beyond a bot.

## Current features (working today)

The URL-shortener core is implemented as an HTTP service:

- Create short aliases (custom or auto-generated) for URLs
- Redirect by alias
- Delete aliases
- Persistence behind a storage interface (SQLite today; **migrating to PostgreSQL** for the hosted, multi-tenant deployment — see roadmap)
- Structured logging (`slog`) with environment-aware output (pretty for local, JSON for dev/prod)
- Config via YAML + environment variables
- HTTP Basic Auth on write endpoints
- chi middleware stack (request ID, real IP, panic recovery)

## Tech stack

| Area | Choice |
|------|--------|
| Language | Go 1.25 |
| HTTP router | [chi](https://github.com/go-chi/chi) v5 |
| Storage | SQLite today → **PostgreSQL** ([`pgx`](https://github.com/jackc/pgx)) for the hosted deployment |
| Config | [cleanenv](https://github.com/ilyakaznacheev/cleanenv) (YAML + env) |
| Logging | stdlib `log/slog` + custom pretty handler |
| Validation | [validator/v10](https://github.com/go-playground/validator) |
| Tests | [testify](https://github.com/stretchr/testify) + [mockery](https://github.com/vektra/mockery) |

## Getting started

### Prerequisites
- Go 1.25+

### Configuration

The service reads its config file from the `CONFIG_PATH` environment variable. A sample lives at `config/local.yaml`:

```yaml
env: "local"            # local | dev | prod
storage_path: "./storage/storage.db"
http_server:
  address: "localhost:8082"
  timeout: 4s
  idle_timeout: 60s
  user: "admin"
```

The Basic Auth password is supplied separately via the `HTTP_SERVER_PASSWORD` env var (never commit it). Copy `local.env.example` to `local.env` and fill it in:

```env
CONFIG_PATH=config/local.yaml
HTTP_SERVER_PASSWORD=change-me
```

### Run

```bash
make run_local   # go fmt + go vet + run
```

The server starts on the configured address (default `localhost:8082`).

## API

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/url` | Basic | Save a URL; optional custom `alias`, otherwise auto-generated |
| `GET` | `/{alias}` | — | Redirect to the original URL |
| `DELETE` | `/url/{alias}` | Basic | Remove an alias |

`POST /url` body:
```json
{ "url": "https://example.com/long", "alias": "optional" }
```

A [Bruno](https://www.usebruno.com/) collection for these endpoints is included under `bruno/`.

## Roadmap

- [x] URL-shortener core (HTTP API + SQLite)
- [ ] Migrate storage to **PostgreSQL** (pgx) for the hosted deployment
- [ ] Extract a generic `Tool` interface + tool router
- [ ] Rename Go module to `github.com/vdzhagaev/pockit`
- [ ] Telegram adapter (`/short`, etc.)
- [ ] Slack app (slash commands; workspace-level install)
- [ ] Reminder tool
- [ ] Quick-share tool (temporary file/text links with expiry)
- [ ] Per-user / per-chat scoping and cross-messenger sync
- [ ] Deployment (Docker + a hosted demo)

---

*Personal project — built to explore clean, pluggable service design in Go.*

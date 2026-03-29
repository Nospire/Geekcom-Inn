# Contributing

## Local testing

```bash
# Terminal 1 — server
go run ./cmd/tavrn-admin

# Terminal 2 — connect via SSH
ssh localhost -p 2222
```

## Branch workflow

```
feature/* ──PR──> dev ──merge──> main (deploy)
```

| Branch | Purpose |
|--------|---------|
| `main` | Production. Runs on the VPS. |
| `dev` | All development. PRs target here. |
| `feature/*` | Short-lived feature branches created from dev. |

1. Create a feature branch from `dev`
2. Open a PR targeting `dev`
3. Test locally
4. When dev is stable, merge to `main` during planned downtime

## Project structure

```
cmd/
  tavrn-admin/     Server binary (SSH server, jukebox engine)
internal/
  chat/            Message parsing and storage types
  hub/             Connection management, broadcasting
  identity/        Nickname generation, flair, color assignment
  jukebox/         Track catalog, engine, streamer (web audio)
  ratelimit/       Chat rate limiting
  room/            Room definitions
  sanitize/        Input sanitization
  server/          Wish SSH server setup
  session/         Session state, message types
  store/           SQLite persistence
  sudoku/          Multiplayer sudoku game logic
  webstream/       Web audio streaming handler
ui/
  app.go           Main Bubble Tea model
  modal.go         Modal system (help, nick, rooms)
  topbar.go        Top bar with room and stats
  sidebar.go       Rooms panel, online users
  chatview.go      Chat message rendering
  gallery.go       Sticky note board
  sudoku_view.go   Multiplayer sudoku view
  overlay.go       Modal overlay compositor
  styles.go        Cantina color palette
  splash.go        Welcome screen
```

## Architecture

**Server** — Wish-based SSH server. Each connection gets a Bubble Tea TUI. A shared hub broadcasts messages between sessions. The jukebox engine manages track playback state for web streaming.

**Web audio** — When started with `--web-audio`, the server runs an HTTP endpoint on `:8090` serving `/stream` (continuous MP3) and `/now-playing` (JSON metadata). Caddy reverse-proxies these to the public domain.

## Admin commands

```bash
# Send banner to all connected users (server-side only)
./tavrn-admin --message "Maintenance in 10 minutes"

# Purge all data
./tavrn-admin purge

# Update server from main branch
./tavrn-admin --update
```

## Tests

```bash
go test ./...
```

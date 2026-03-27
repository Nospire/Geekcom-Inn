# Client Binary + Audio Streaming Design

The `tavrn` client binary gives users actual audio playback. Bare SSH users get the full TUI but no sound. The client binary adds a second SSH channel for MP3 streaming.

## Architecture

```
Server (internal/jukebox/)
  Engine → Streamer
    fetches MP3 from track URL
    writes chunks to all tavrn-audio channels

SSH connection (two channels)
  session channel    → Bubble Tea TUI (same as bare SSH)
  tavrn-audio channel → MP3 byte stream

Client (cmd/tavrn-client/)
  crypto/ssh → session channel → terminal
             → tavrn-audio channel → beep v2 → oto v3 → speakers
```

## Server-Side Streamer

Lives in `internal/jukebox/streamer.go`:

- Engine calls `streamer.TrackChanged(track)` on track switch
- Streamer HTTP GETs the track URL, reads 8KB chunks, writes to all audio conns
- `AddConn(io.WriteCloser)` / `RemoveConn()` for client lifecycle
- Slow clients dropped via non-blocking writes
- Track change cancels current fetch, starts new one

### Wire Protocol

Per-track:
1. `[4 bytes: JSON length][JSON metadata]\n` — title, artist, source, duration
2. Raw MP3 bytes until next track header

## Client Binary

`cmd/tavrn-client/main.go` — single command entry:

```
tavrn              connect to tavrn.sh with audio
tavrn --no-audio   connect without audio
tavrn --update     self-update from GitHub releases
tavrn --version    print version
```

### Startup Flow

1. Parse flags
2. Dial `tavrn.sh:22` via `crypto/ssh` with user's SSH key
3. Open `session` channel → wire to terminal (stdin/stdout/resize)
4. Open `tavrn-audio` channel (unless `--no-audio`)
5. Audio goroutine: read metadata header → beep MP3 decode → oto speaker
6. Block until session ends

### Dependencies

- `golang.org/x/crypto/ssh` — SSH client
- `gopxl/beep` v2 — MP3 decoding
- `ebitengine/oto` v3 — cross-platform audio output
- `creativeprojects/go-selfupdate` — self-update

### Distribution

- Same repo: `cmd/tavrn-client/` alongside `cmd/tavrn/`
- Binary renamed to `tavrn` via GoReleaser
- Install: `go install tavrn.sh/cmd/tavrn-client@latest`
- Targets: darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64

## Server-Side Changes

### Wish Channel Handler

In `internal/server/server.go`, add a custom channel handler for `"tavrn-audio"`:
- When client requests channel type `"tavrn-audio"`, accept and pass to streamer
- Bare SSH clients never request this — no change for them

### Engine Integration

- Engine notifies streamer on track change via callback
- Streamer registered in server config alongside engine

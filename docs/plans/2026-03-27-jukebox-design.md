# Jukebox Design

Shared music jukebox for tavrn.sh — everyone listens together, users search and vote for the next track.

## Architecture Overview

```
                      tavrn server
  ┌───────────┐  ┌───────────┐  ┌─────────────┐
  │  Jamendo   │  │  Radio    │  │  YouTube    │
  │  Backend   │  │  Browser  │  │  (opt-in)   │
  └─────┬─────┘  └─────┬─────┘  └──────┬──────┘
        └───────┬───────┴───────────────┘
          ┌─────▼─────┐
          │  Jukebox   │  single source of truth
          │  Engine    │  current track, queue,
          │            │  votes, audio stream
          └──┬─────┬───┘
             │     │
     ┌───────▼┐  ┌─▼──────────┐
     │ TUI    │  │ Audio      │
     │ "Now   │  │ Channel    │
     │ Playing│  │ (SSH mux)  │
     └────┬───┘  └─────┬──────┘
           │            │
    ┌──────▼──────┐  ┌──▼─────────────┐
    │ bare SSH    │  │ tavrn binary   │
    │ sees UI +   │  │ sees UI +      │
    │ votes, no   │  │ HEARS audio    │
    │ audio       │  │ via beep/oto   │
    └─────────────┘  └────────────────┘
```

## Music Backends

Each source implements a pluggable interface:

```go
type MusicBackend interface {
    Name() string
    Search(ctx context.Context, query string) ([]Track, error)
    StreamURL(ctx context.Context, track Track) (string, error)
    Enabled() bool
}
```

### Jamendo (default, always enabled)

- 600K+ CC-licensed tracks from 40K+ independent artists
- Strong in electronic, ambient, lo-fi, jazz
- REST API: `GET /v3.0/tracks/?client_id=XXX&search=QUERY`
- Returns direct MP3/OGG streaming URLs
- Free tier: 35K API requests/month
- Register at devportal.jamendo.com for client_id

### Radio Browser (always enabled)

- 30K+ internet radio stations worldwide
- Search by name, genre, country, language
- Returns resolved Icecast/Shoutcast stream URLs
- No auth required
- Go library: `randomtoy/radiobrowser-go`
- Provides "radio mode" complement to on-demand tracks

### YouTube (opt-in, disabled by default)

- Server admin enables via config flag
- Uses Piped API for search + stream URL extraction
- Falls back to yt-dlp if installed on server
- Legally gray — ToS violation risk, kept as user-responsible opt-in
- Go library: `lrstanley/go-ytdlp` for yt-dlp bindings
- Not bundled by default, keeping the project legally clean

### Search Merging

When a user searches, the engine queries all enabled backends in parallel, merges results, and tags each with its source: `[jamendo]`, `[radio]`, `[youtube]`.

## Jukebox Engine

Single server-side goroutine managing all shared state.

### State

- `currentTrack` — what's playing (title, artist, source, duration, position)
- `queue []Track` — upcoming tracks added by users
- `requestPool map[TrackID]int` — tracks users have requested, with count
- `shortlist [5]Track` — top 5 most requested, frozen when voting opens
- `votes map[TrackID]map[UserID]bool` — one vote per user per round
- `phase` — `playing`, `requesting`, or `voting`

### Playback Lifecycle

```
PLAYING current track
    │
    │  track hits 75% duration
    ▼
REQUESTING phase opens
    │  users search + add songs
    │  requestPool tallies popularity
    │
    │  track hits 90% duration OR 30s before end
    ▼
VOTING phase opens
    │  top 5 from requestPool → frozen shortlist
    │  users get one vote each
    │
    │  current track ends
    ▼
Winner plays → clear votes, clear requestPool
    │  tie → random among tied
    │  no votes → random from shortlist
    │  no requests → random from Jamendo popular/trending
    │
    └──→ back to PLAYING
```

### Rules

- One vote per user per voting round
- Top 5 requests make the shortlist (ties broken by earliest request time)
- Tie in votes → random among tied
- No votes cast → random from shortlist
- Empty request pool → engine picks random from Jamendo popular endpoint
- Skip requires >50% of connected users

## Audio Streaming Over SSH

Bare `ssh tavrn.sh` cannot play audio — terminals have no audio API, only `\a` bell. The `tavrn` binary exists to solve this.

### SSH Channel Multiplexing

The `tavrn` binary is a custom SSH client using `crypto/ssh` that opens two channels over one connection:

1. `session` channel — standard terminal session, Bubble Tea TUI
2. `tavrn-audio` channel — custom channel type for MP3 byte streaming

```
tavrn binary (client side)
├── ssh session channel  →  terminal (Bubble Tea TUI)
└── tavrn-audio channel  →  beep v2 decoder → oto speaker output
```

### Wire Format

- MP3 stream — universally decodable, beep v2 handles natively
- Simple header per track: `[4 bytes length][JSON metadata]\n[MP3 bytes...]`
- Server fetches track URL from backend, proxies MP3 bytes to all connected tavrn-audio channels

### Sync Strategy

- Server maintains canonical byte offset / timestamp for current track
- New client joins mid-track → server sends from current position (misses beginning, stays in sync)
- No complex clock sync — everyone gets the same bytes at the same time

### Fallback

If `tavrn-audio` channel isn't requested (bare SSH), nothing breaks. TUI works identically.

## `tavrn` CLI Binary

### Commands

```
tavrn              connect to tavrn.sh with audio enabled
tavrn --update     self-update to latest version
tavrn --version    print version
tavrn --no-audio   connect without audio (same as bare SSH)
```

No subcommands, no config files, no setup. Single command entry.

### What It Does

1. Dials `tavrn.sh:22` via `crypto/ssh` with keyboard-interactive auth (zero-signup)
2. Opens `session` channel → hooks to user's terminal (raw mode, resize events)
3. Opens `tavrn-audio` channel → goroutine pipes MP3 bytes into beep v2 → oto speaker
4. User sees identical TUI to bare SSH, but also hears music

### Self-Update

- `tavrn --update` hits GitHub releases API (or `tavrn.sh/releases/latest`)
- Downloads binary for user's OS/arch
- Replaces itself in-place
- Go library: `creativeprojects/go-selfupdate`

### Distribution

```
go install tavrn.sh/client@latest
brew install tavrn                    # future
curl -sSL tavrn.sh/install | sh      # one-liner
```

### Build Targets

`darwin/amd64`, `darwin/arm64`, `linux/amd64`, `linux/arm64`, `windows/amd64` — GoReleaser handles cross-compilation.

## Jukebox Modal (F4)

Three tabs within one modal, follows existing modal patterns in `ui/modal.go`.

> Post Note modal moves from F4 to F5.

### Keybinds

- `F4` — open jukebox modal
- `Tab` — switch between tabs
- `j/k` — navigate lists
- `Enter` — add track (search tab) or cast vote (vote tab)
- `Esc` — close modal

### Tab 1: Now Playing

```
╱╱╱╱╱╱╱ ♪ Jukebox ╱╱╱╱╱╱╱
│                            │
│  [Now Playing] Search  Vote│
│                            │
│  Lo-fi Sunset              │
│  ambient_collective        │
│  [jamendo]                 │
│  ▓▓▓▓▓▓▓░░░░ 2:34 / 3:45 │
│                            │
│  Queue:                    │
│  1. Midnight Jazz Club     │
│  2. Chill Waves            │
│                            │
│  5 listening · 3 requests  │
│                            │
╱╱ tab: switch · esc: close ╱╱
```

### Tab 2: Search

```
╱╱╱╱╱╱╱ ♪ Jukebox ╱╱╱╱╱╱╱
│                            │
│  Now Playing  [Search] Vote│
│                            │
│  > midnight jazz_          │
│                            │
│  1. Midnight Jazz  [jamendo│]
│  2. Jazz After Dark [radio]│
│  3. Late Night     [youtube│]
│                            │
│  j/k: navigate             │
│  enter: add to requests    │
│                            │
╱╱ tab: switch · esc: close ╱╱
```

### Tab 3: Vote

```
╱╱╱╱╱╱╱ ♪ Jukebox ╱╱╱╱╱╱╱
│                            │
│  Now Playing  Search [Vote]│
│                            │
│  VOTE FOR NEXT TRACK       │
│                            │
│  › 1. Midnight Jazz    ▓▓▓ │
│    2. Chill Waves      ▓▓  │
│    3. Ambient Rain     ▓▓  │
│    4. Deep Focus       ▓   │
│    5. Lo-fi Morning    ▓   │
│                            │
│  j/k: navigate             │
│  enter: cast vote (1 only) │
│                            │
╱╱ tab: switch · esc: close ╱╱
```

### Implementation Pattern

Follows existing modal conventions in `ui/modal.go`:

- Add `ModalJukebox` to `ModalType` enum
- `JukeboxModal` struct with `tab int`, `cursor int`, `searchInput textinput.Model`, phase state
- `Update(msg tea.Msg) (JukeboxModal, tea.Cmd)` handles Tab, j/k, Enter, typing
- `View(width, height int) string` renders with `╱╱╱` headers, `ColorBorder`/`ColorHighlight` palette
- Messages: `JukeboxSearchMsg`, `JukeboxAddMsg`, `JukeboxVoteMsg` flow back to App
- Add `jukeboxModal JukeboxModal` field to `App` struct
- Route in `updateModal()` and `View()` overlay switch

## Now Playing Bar

Persistent in the top bar, centered between room name and online count.

```
#tavern  │  ♪ Lo-fi Sunset ░▁▃▅▇▅▃▁░ 2:34  │  5 online
```

### Wave Animation

Unicode block characters cycling every ~200ms:

```
░▁▃▅▇▅▃▁░
░▃▅▇▅▃▁░▁
░▅▇▅▃▁░▁▃
```

Purely cosmetic — not tied to actual audio levels. Implemented as a `tea.Tick` command rotating a frame index through the block character array `▁▂▃▄▅▆▇█`.

When no track is playing: `♪ --` or hidden entirely.

## Go Dependencies

### Server

- `gopxl/beep` v2 — audio decode (MP3/OGG) for re-streaming
- `randomtoy/radiobrowser-go` — Radio Browser API client
- `lrstanley/go-ytdlp` — yt-dlp bindings (optional)

### Client (`tavrn` binary)

- `golang.org/x/crypto/ssh` — SSH client
- `gopxl/beep` v2 — MP3 decoding
- `ebitengine/oto` v3 — cross-platform audio output (used by beep)
- `creativeprojects/go-selfupdate` — self-update mechanism

## Key Design Decisions

| Decision | Rationale |
|----------|-----------|
| Separate `tavrn` binary for audio | SSH protocol has no audio channel; terminals cannot play sound |
| Jamendo as default backend | Fully legal, CC-licensed, good ambient/lo-fi catalog, simple REST API |
| YouTube opt-in only | ToS gray area; server admin's choice, not baked into project |
| Top 5 shortlist → vote | Prevents ballot splitting, keeps voting focused |
| MP3 over custom SSH channel | Universal format, beep v2 decodes natively, simple framing |
| No clock sync | Byte-level streaming keeps everyone in sync without complexity |
| F4 keybind | Consistent with existing F-key modal pattern |

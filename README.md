# tavrn.sh

A public SSH terminal tavern. Connect, chat, listen to music, and hang out — all from your terminal.

```
ssh tavrn.sh
```

No signup. No account. Your SSH key is your identity.

## What is this

tavrn.sh is a multi-user terminal application accessible over SSH. It runs a TUI built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) served through [Wish](https://github.com/charmbracelet/wish), so anyone with an SSH key can connect and interact in real time.

**Rooms** — Chat in the lounge, post sticky notes on the gallery board, or drop ideas in suggestions.

**Jukebox** — A shared music player powered by [Jamendo](https://www.jamendo.com/). Search for tracks, request songs, and vote on what plays next. Everyone in the tavern sees the same Now Playing bar.

**Audio** — Install the `tavrn` client binary to hear music through your speakers. Bare SSH gives you the full TUI without audio.

**Gallery** — A collaborative sticky note board. Post notes, drag them around, read what others left behind.

Everything resets weekly. Nothing is permanent.

## Connect

SSH in directly:

```
ssh tavrn.sh
```

Or install the client binary for audio:

```
go install tavrn.sh/cmd/tavrn-client@latest
tavrn
```

The client requires [mpv](https://mpv.io/) for audio playback.

## Keybinds

| Key | Action |
|-----|--------|
| F1 | Help |
| F2 | Change nickname |
| F3 | Switch rooms |
| F4 | Jukebox (search, request, vote) |
| F5 | Post a note (gallery) |
| ESC | Close modal |

## Stack

- Go
- [Bubble Tea v2](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Wish v2](https://github.com/charmbracelet/wish) — SSH server
- [Lipgloss v2](https://github.com/charmbracelet/lipgloss) — Terminal styling
- [Jamendo API](https://developer.jamendo.com/) — Music catalog
- SQLite — Data storage

## Self-hosting

```bash
git clone https://github.com/youruser/tavrn.git
cd tavrn
JAMENDO_CLIENT_ID=your_key go run ./cmd/tavrn
```

Get a free Jamendo API key at [devportal.jamendo.com](https://devportal.jamendo.com/). The server runs on port 2222 by default.

## License

MIT

# tavrn.sh

A public terminal tavern. Chat, listen to music, and hang out — all from your terminal.

```
ssh tavrn.sh
```

No signup. No account. Your SSH key is your identity.

---

### The tavern

Connect and you're in. Chat with strangers in the lounge, post sticky notes on the gallery board, or search for music on the shared jukebox. Everyone sees the same Now Playing bar. Vote on what plays next.

Everything resets weekly. Nothing is permanent.

### Two ways in

**Bare SSH** — Works everywhere. You get the full TUI: chat, gallery, jukebox controls, voting. No sound.

```
ssh tavrn.sh
```

**tavrn binary** — The full experience. Same TUI, plus music through your speakers.

```
brew install tavrn
tavrn
```

or

```
go install tavrn.sh/cmd/tavrn-client@latest
tavrn
```

Requires [mpv](https://mpv.io/) for audio. The binary checks on launch and tells you how to install it if missing.

### Keybinds

```
F1  help          F2  nickname       F3  rooms
F4  jukebox       F5  post note      ESC close
```

### Self-hosting

```bash
git clone https://github.com/neur0map/tavrn.git
cd tavrn
JAMENDO_CLIENT_ID=your_key go run ./cmd/tavrn
```

Free API key from [devportal.jamendo.com](https://devportal.jamendo.com/).

### Contributing

See [docs/CONTRIBUTING.md](docs/CONTRIBUTING.md) for development setup, branch workflow, and architecture.

### Built with

[Bubble Tea](https://github.com/charmbracelet/bubbletea) · [Wish](https://github.com/charmbracelet/wish) · [Lipgloss](https://github.com/charmbracelet/lipgloss) · [Jamendo](https://www.jamendo.com/) · Go · SQLite

### License

MIT

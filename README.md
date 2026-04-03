<div align="center">

# tavrn.sh

A terminal tavern over SSH.

Chat and hang out with strangers — right from your terminal.

No signup. No account. Your SSH key is your identity.

Everything resets weekly. Nothing is permanent.

</div>

---

<div align="center">

### Connect

```
ssh tavrn.sh
```

</div>

---

### What's inside

**Rooms** — Lounge, gallery, games, suggestions, and wargame CTF rooms.

**Gallery** — A shared sticky note board. Post, drag, read what others left behind.

**Wargames** — OverTheWire CTF rooms with flag submission, leaderboard, and points.

**GIFs** — Search and send animated GIFs inline in chat with `/gif`.

**Bartender** — The lounge bartender. Tag with `@bartender`.

**Music** — 24/7 radio streaming.

### Keybinds

```
F1  help        F2  nickname     F3  rooms       F4  mentions
F5  post note   F6  tankard      F7  leaderboard
```

---

### Run your own

Fork this and run your own tavern. The `tavrn` name is reserved (see [TRADEMARK.md](TRADEMARK.md)), but the engine is yours.

**Quick start:**

```bash
cp tavern.yaml.example tavern.yaml   # configure your tavern name, domain, rooms
cp .env.example .env                  # add API keys (all optional)
make run                              # build and connect
```

**Docker:**

```bash
cp tavern.yaml.example tavern.yaml
cp .env.example .env
docker compose up
```

**Environment variables** (all optional):

| Variable | What it does |
|---|---|
| `OPENAI_API_KEY` | Bartender character |
| `KLIPY_API_KEY` | GIF search (`/gif`) |
| `EXA_API_KEY` | Bartender web search |
| `TAVRN_PORT` | SSH port (default: 2222) |

**Admin CLI:**

```
tavrn --help                         Full command list
tavrn --set-flag bandit 1 "flag"     Set a wargame flag
tavrn --bartender-off                Disable bartender live
```

See [deploy/SETUP.md](deploy/SETUP.md) for VPS deployment with systemd + Caddy.

### Built with

[Bubble Tea](https://github.com/charmbracelet/bubbletea) · [Wish](https://github.com/charmbracelet/wish) · [Lipgloss](https://github.com/charmbracelet/lipgloss) · Go · SQLite

### License

MIT — see [LICENSE](LICENSE)

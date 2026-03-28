# Roadmap

## Theater Room — Watch YouTube Together

A shared room where users paste YouTube links, vote on what plays next, and watch together.

- Users with the client binary (mpv installed) get video in a separate mpv window
- Bare SSH users get audio-only via the existing streamer
- Server uses yt-dlp for metadata and audio fallback
- New wire protocol message type for "play this URL"
- Voting system: paste link, upvote, most votes plays next
- TUI shows now-playing info, queue, and vote counts while users chat

## Tavern Games

Multiplayer terminal games that fit the tavern theme.

- `/roll 2d6` — dice rolling for tabletop RPG sessions
- Trivia — timed questions, scoreboard, themed rounds
- Word games — hangman, word chains
- Tic-tac-toe, connect four — challenge another user
- Text adventure — room votes on choices, story unfolds together

## Hacker News / Reddit Reader

A room where users browse and discuss threads together.

- HN: public JSON API, no auth needed
- Reddit: scrape old.reddit.com/.json, no API key
- Scrollable thread view in the TUI
- Everyone reads the same thread and discusses in chat
- `/hn top` `/hn new` `/reddit r/golang` to navigate

## Mastodon Feed

Public Mastodon timeline in a dedicated room.

- Public API, no auth required for public posts
- Render toots with author, content, boosts
- Users discuss posts in real-time chat
- `/fedi trending` `/fedi local instance.social`

## Radio Requests + Voting

Let users browse the catalog and queue tracks.

- Browse tracks by genre in the jukebox modal
- Request a track — goes into the queue
- Other users upvote requests
- Most-voted track plays next
- Falls back to random if queue is empty

## DMs

Private messages between users.

- `/dm @nickname message` to whisper
- Conversation appears in a side panel or modal
- Routed by SSH fingerprint, no accounts needed

## Collaborative ASCII Canvas

A shared drawing room — r/place in the terminal.

- Grid canvas, users move cursor and place characters
- Color support via ANSI
- Canvas persists until weekly reset
- Watch others draw in real-time

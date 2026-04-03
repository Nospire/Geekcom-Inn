package dm

import (
	"database/sql"
	"sort"
	"sync"
	"time"
)

// DirectMessage represents a single DM row.
type DirectMessage struct {
	ID        int
	FromFP    string
	ToFP      string
	FromNick  string
	Text      string
	Read      bool
	CreatedAt time.Time
}

// Conversation is a summary for the inbox list.
type Conversation struct {
	PeerFP      string
	PeerNick    string
	LastMessage string
	LastTime    time.Time
	UnreadCount int
}

// Store handles direct message persistence.
type Store struct {
	db *sql.DB
	mu sync.Mutex
}

// New creates a DM store using an existing DB handle.
func New(db *sql.DB) (*Store, error) {
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS direct_messages (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		from_fp    TEXT NOT NULL,
		to_fp      TEXT NOT NULL,
		from_nick  TEXT NOT NULL,
		text       TEXT NOT NULL,
		read       INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_dm_to ON direct_messages(to_fp, created_at);
	CREATE INDEX IF NOT EXISTS idx_dm_conv ON direct_messages(from_fp, to_fp, created_at);
	`
	_, err := s.db.Exec(schema)
	return err
}

// Send stores a new DM.
func (s *Store) Send(fromFP, toFP, fromNick, text string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	res, err := s.db.Exec(`
		INSERT INTO direct_messages (from_fp, to_fp, from_nick, text)
		VALUES (?, ?, ?, ?)
	`, fromFP, toFP, fromNick, text)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

// Messages returns the conversation between two users, newest last.
func (s *Store) Messages(fpA, fpB string, limit int) ([]DirectMessage, error) {
	rows, err := s.db.Query(`
		SELECT id, from_fp, to_fp, from_nick, text, read, created_at
		FROM direct_messages
		WHERE (from_fp = ? AND to_fp = ?) OR (from_fp = ? AND to_fp = ?)
		ORDER BY created_at DESC
		LIMIT ?
	`, fpA, fpB, fpB, fpA, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []DirectMessage
	for rows.Next() {
		var m DirectMessage
		var readInt int
		var ts string
		if err := rows.Scan(&m.ID, &m.FromFP, &m.ToFP, &m.FromNick, &m.Text, &readInt, &ts); err != nil {
			continue
		}
		m.Read = readInt != 0
		m.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", ts)
		msgs = append(msgs, m)
	}
	// Reverse to oldest-first
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, nil
}

// MarkRead marks all messages TO this user FROM the peer as read.
func (s *Store) MarkRead(myFP, peerFP string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(`
		UPDATE direct_messages SET read = 1
		WHERE to_fp = ? AND from_fp = ? AND read = 0
	`, myFP, peerFP)
	return err
}

// UnreadCount returns total unread DMs for a user.
func (s *Store) UnreadCount(fp string) int {
	row := s.db.QueryRow(`SELECT COUNT(*) FROM direct_messages WHERE to_fp = ? AND read = 0`, fp)
	var count int
	row.Scan(&count)
	return count
}

// Conversations returns the inbox for a user: one entry per peer, sorted by most recent.
func (s *Store) Conversations(fp string) []Conversation {
	// Get all peers this user has DM'd with (sent or received)
	rows, err := s.db.Query(`
		SELECT DISTINCT CASE WHEN from_fp = ? THEN to_fp ELSE from_fp END AS peer_fp
		FROM direct_messages
		WHERE from_fp = ? OR to_fp = ?
	`, fp, fp, fp)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var peers []string
	for rows.Next() {
		var peerFP string
		if err := rows.Scan(&peerFP); err != nil {
			continue
		}
		peers = append(peers, peerFP)
	}

	var convos []Conversation
	for _, peerFP := range peers {
		c := Conversation{PeerFP: peerFP}

		// Last message in this conversation
		row := s.db.QueryRow(`
			SELECT from_nick, text, created_at
			FROM direct_messages
			WHERE (from_fp = ? AND to_fp = ?) OR (from_fp = ? AND to_fp = ?)
			ORDER BY created_at DESC LIMIT 1
		`, fp, peerFP, peerFP, fp)
		var ts string
		row.Scan(&c.PeerNick, &c.LastMessage, &ts)
		c.LastTime, _ = time.Parse("2006-01-02 15:04:05", ts)

		// Get the actual peer nickname (from their last sent message to us, or our last to them)
		nickRow := s.db.QueryRow(`
			SELECT from_nick FROM direct_messages
			WHERE from_fp = ?
			ORDER BY created_at DESC LIMIT 1
		`, peerFP)
		var peerNick string
		if nickRow.Scan(&peerNick) == nil && peerNick != "" {
			c.PeerNick = peerNick
		}

		// Unread count from this peer
		unreadRow := s.db.QueryRow(`
			SELECT COUNT(*) FROM direct_messages
			WHERE from_fp = ? AND to_fp = ? AND read = 0
		`, peerFP, fp)
		unreadRow.Scan(&c.UnreadCount)

		convos = append(convos, c)
	}

	// Sort by most recent
	sort.Slice(convos, func(i, j int) bool {
		return convos[i].LastTime.After(convos[j].LastTime)
	})

	return convos
}

// Purge deletes all DMs. Called during weekly reset.
func (s *Store) Purge() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(`DELETE FROM direct_messages`)
	return err
}

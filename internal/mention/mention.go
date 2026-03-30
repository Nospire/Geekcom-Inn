package mention

import (
	"strings"
	"time"
	"unicode"
)

// Mention represents an unread @mention.
type Mention struct {
	Room       string
	Author     string
	ColorIndex int
	Text       string
	Timestamp  time.Time
	Read       bool
}

// ExtractTokens finds all @name tokens in text.
// Only matches @ at the start of the string or after whitespace.
func ExtractTokens(text string) []string {
	var tokens []string
	words := strings.Fields(text)
	for _, word := range words {
		if len(word) > 1 && word[0] == '@' {
			// Strip trailing punctuation
			name := strings.TrimRightFunc(word[1:], func(r rune) bool {
				return unicode.IsPunct(r) && r != '_' && r != '~'
			})
			if name != "" {
				tokens = append(tokens, name)
			}
		}
	}
	return tokens
}

// IsMentioned checks if nickname is @mentioned in text (case-insensitive, full word match).
func IsMentioned(text, nickname string) bool {
	tokens := ExtractTokens(text)
	lower := strings.ToLower(nickname)
	for _, tok := range tokens {
		if strings.ToLower(tok) == lower {
			return true
		}
	}
	return false
}

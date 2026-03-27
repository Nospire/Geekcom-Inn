package identity

import (
	"crypto/sha256"
	"encoding/hex"
)

// DefaultNickname returns the first 8 hex chars of SHA256(fingerprint).
func DefaultNickname(fingerprint string) string {
	hash := sha256.Sum256([]byte(fingerprint))
	return hex.EncodeToString(hash[:])[:8]
}

// ColorIndex returns 0-11 for mapping a fingerprint to one of 12 cantina colors.
func ColorIndex(fingerprint string) int {
	hash := sha256.Sum256([]byte(fingerprint))
	return int(hash[0]) % 12
}

// HasFlair returns true if the user has connected 3+ times this week.
func HasFlair(visitCount int) bool {
	return visitCount >= 3
}

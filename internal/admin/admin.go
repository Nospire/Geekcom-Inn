package admin

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"tavrn/internal/store"
)

type Admin struct {
	store            *store.Store
	adminFingerprint string
}

func New(store *store.Store, adminFP string) *Admin {
	return &Admin{
		store:            store,
		adminFingerprint: adminFP,
	}
}

func (a *Admin) IsAdmin(fingerprint string) bool {
	return a.adminFingerprint != "" && fingerprint == a.adminFingerprint
}

// HandleCommand processes an admin command. Returns a status message.
func (a *Admin) HandleCommand(callerFP, command, args string) (string, error) {
	if !a.IsAdmin(callerFP) {
		return "", errors.New("not authorized")
	}

	switch command {
	case "ban":
		return a.handleBan(args)
	case "unban":
		return a.handleUnban(args)
	case "purge":
		return a.handlePurge()
	default:
		return "", fmt.Errorf("unknown admin command: %s", command)
	}
}

func (a *Admin) handleBan(args string) (string, error) {
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return "", errors.New("usage: /ban <fingerprint> [duration]")
	}

	fp := parts[0]
	var expiresAt *time.Time

	if len(parts) > 1 {
		d, err := time.ParseDuration(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid duration: %s", parts[1])
		}
		t := time.Now().Add(d)
		expiresAt = &t
	} else {
		t := time.Now().Add(24 * time.Hour)
		expiresAt = &t
	}

	if err := a.store.Ban(fp, "admin ban", expiresAt); err != nil {
		return "", err
	}

	short := fp
	if len(short) > 8 {
		short = short[:8]
	}
	return fmt.Sprintf("Banned %s", short), nil
}

func (a *Admin) handleUnban(args string) (string, error) {
	fp := strings.TrimSpace(args)
	if fp == "" {
		return "", errors.New("usage: /unban <fingerprint>")
	}
	if err := a.store.Unban(fp); err != nil {
		return "", err
	}

	short := fp
	if len(short) > 8 {
		short = short[:8]
	}
	return fmt.Sprintf("Unbanned %s", short), nil
}

func (a *Admin) handlePurge() (string, error) {
	if err := a.store.PurgeAll(); err != nil {
		return "", err
	}
	return "Purge complete", nil
}

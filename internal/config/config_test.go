package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "tavern.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadValid(t *testing.T) {
	path := writeConfig(t, `
tavern:
  name: "Test Tavern"
  domain: "test.sh"
  tagline: "a test place"

owner:
  name: "testowner"
  fingerprint: "SHA256:abc123"

rooms:
  - name: "lobby"
    type: chat
  - name: "art"
    type: gallery
  - name: "arcade"
    type: games
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Tavern.Name != "Test Tavern" {
		t.Errorf("name = %q", cfg.Tavern.Name)
	}
	if cfg.Tavern.Domain != "test.sh" {
		t.Errorf("domain = %q", cfg.Tavern.Domain)
	}
	if cfg.Tavern.Tagline != "a test place" {
		t.Errorf("tagline = %q", cfg.Tavern.Tagline)
	}
	if cfg.Owner.Name != "testowner" {
		t.Errorf("owner = %q", cfg.Owner.Name)
	}
	if cfg.Owner.Fingerprint != "SHA256:abc123" {
		t.Errorf("fingerprint = %q", cfg.Owner.Fingerprint)
	}
	if len(cfg.Rooms) != 3 {
		t.Fatalf("rooms = %d, want 3", len(cfg.Rooms))
	}
	if cfg.Rooms[0].Name != "lobby" || cfg.Rooms[0].Type != "chat" {
		t.Errorf("room[0] = %+v", cfg.Rooms[0])
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/tavern.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadMissingName(t *testing.T) {
	path := writeConfig(t, `
tavern:
  domain: "test.sh"
owner:
  name: "testowner"
  fingerprint: "SHA256:abc"
rooms:
  - name: "lobby"
    type: chat
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing tavern.name")
	}
}

func TestLoadMissingOwnerFingerprint(t *testing.T) {
	path := writeConfig(t, `
tavern:
  name: "Test"
  domain: "test.sh"
owner:
  name: "testowner"
rooms:
  - name: "lobby"
    type: chat
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing owner.fingerprint")
	}
}

func TestLoadNoRooms(t *testing.T) {
	path := writeConfig(t, `
tavern:
  name: "Test"
  domain: "test.sh"
owner:
  name: "testowner"
  fingerprint: "SHA256:abc"
rooms: []
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for empty rooms")
	}
}

func TestLoadInvalidRoomType(t *testing.T) {
	path := writeConfig(t, `
tavern:
  name: "Test"
  domain: "test.sh"
owner:
  name: "testowner"
  fingerprint: "SHA256:abc"
rooms:
  - name: "lobby"
    type: invalid_type
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid room type")
	}
}

func TestRoomNames(t *testing.T) {
	path := writeConfig(t, `
tavern:
  name: "Test"
  domain: "test.sh"
owner:
  name: "testowner"
  fingerprint: "SHA256:abc"
rooms:
  - name: "lobby"
    type: chat
  - name: "art"
    type: gallery
  - name: "games"
    type: games
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	names := cfg.RoomNames()
	if len(names) != 3 || names[0] != "lobby" || names[1] != "art" || names[2] != "games" {
		t.Errorf("RoomNames() = %v", names)
	}
}

func TestFirstRoom(t *testing.T) {
	path := writeConfig(t, `
tavern:
  name: "Test"
  domain: "test.sh"
owner:
  name: "testowner"
  fingerprint: "SHA256:abc"
rooms:
  - name: "lobby"
    type: chat
  - name: "art"
    type: gallery
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.FirstRoom() != "lobby" {
		t.Errorf("FirstRoom() = %q, want lobby", cfg.FirstRoom())
	}
}

func TestRoomIsType(t *testing.T) {
	path := writeConfig(t, `
tavern:
  name: "Test"
  domain: "test.sh"
owner:
  name: "testowner"
  fingerprint: "SHA256:abc"
rooms:
  - name: "lobby"
    type: chat
  - name: "art"
    type: gallery
  - name: "arcade"
    type: games
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if !cfg.RoomIsType("art", "gallery") {
		t.Error("art should be type gallery")
	}
	if cfg.RoomIsType("lobby", "gallery") {
		t.Error("lobby should not be type gallery")
	}
	if !cfg.RoomIsType("arcade", "games") {
		t.Error("arcade should be type games")
	}
	if cfg.RoomIsType("unknown", "chat") {
		t.Error("unknown room should return false")
	}
}

func TestRoomTypeMap(t *testing.T) {
	path := writeConfig(t, `
tavern:
  name: "Test"
  domain: "test.sh"
owner:
  name: "testowner"
  fingerprint: "SHA256:abc"
rooms:
  - name: "lobby"
    type: chat
  - name: "art"
    type: gallery
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	m := cfg.RoomTypeMap()
	if m["lobby"] != "chat" {
		t.Errorf("lobby type = %q", m["lobby"])
	}
	if m["art"] != "gallery" {
		t.Errorf("art type = %q", m["art"])
	}
}

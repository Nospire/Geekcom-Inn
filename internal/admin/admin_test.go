package admin

import (
	"testing"

	"tavrn/internal/store"
)

func tempStore(t *testing.T) *store.Store {
	t.Helper()
	path := t.TempDir() + "/test.db"
	s, err := store.New(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestBanCommand(t *testing.T) {
	st := tempStore(t)
	a := New(st, "admin_fp")

	msg, err := a.HandleCommand("admin_fp", "ban", "target_fp")
	if err != nil {
		t.Fatalf("ban: %v", err)
	}
	if msg == "" {
		t.Error("expected confirmation message")
	}

	banned, _ := st.IsBanned("target_fp")
	if !banned {
		t.Error("target should be banned")
	}
}

func TestBanCommandWithDuration(t *testing.T) {
	st := tempStore(t)
	a := New(st, "admin_fp")

	_, err := a.HandleCommand("admin_fp", "ban", "target_fp 48h")
	if err != nil {
		t.Fatalf("ban with duration: %v", err)
	}

	banned, _ := st.IsBanned("target_fp")
	if !banned {
		t.Error("target should be banned")
	}
}

func TestUnbanCommand(t *testing.T) {
	st := tempStore(t)
	a := New(st, "admin_fp")

	a.HandleCommand("admin_fp", "ban", "target_fp")
	a.HandleCommand("admin_fp", "unban", "target_fp")

	banned, _ := st.IsBanned("target_fp")
	if banned {
		t.Error("target should be unbanned")
	}
}

func TestNonAdminRejected(t *testing.T) {
	st := tempStore(t)
	a := New(st, "admin_fp")

	_, err := a.HandleCommand("not_admin", "ban", "target_fp")
	if err == nil {
		t.Error("non-admin should be rejected")
	}
}

func TestIsAdmin(t *testing.T) {
	a := New(nil, "admin_fp")
	if !a.IsAdmin("admin_fp") {
		t.Error("should be admin")
	}
	if a.IsAdmin("other_fp") {
		t.Error("should not be admin")
	}
}

package identity

import (
	"strings"
	"testing"
)

func TestDefaultNickname(t *testing.T) {
	fp := "SHA256:abc123def456"
	n1 := DefaultNickname(fp)
	n2 := DefaultNickname(fp)
	if n1 != n2 {
		t.Errorf("nondeterministic: %q != %q", n1, n2)
	}
	if len(n1) < 5 {
		t.Errorf("too short: %q", n1)
	}
	// Should contain underscore (adj_noun format)
	if !contains(n1, "_") {
		t.Errorf("expected adj_noun format, got %q", n1)
	}
	// Should contain #XXXX discriminator
	if !strings.Contains(n1, "#") {
		t.Errorf("expected #XXXX discriminator, got %q", n1)
	}
	parts := strings.SplitN(n1, "#", 2)
	if len(parts) != 2 || len(parts[1]) != 4 {
		t.Errorf("discriminator should be 4 digits, got %q", n1)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && findSub(s, sub)
}

func findSub(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func TestDefaultNicknameDifferentKeys(t *testing.T) {
	n1 := DefaultNickname("SHA256:aaa")
	n2 := DefaultNickname("SHA256:bbb")
	if n1 == n2 {
		t.Error("different fingerprints should produce different nicks")
	}
}

func TestColorIndex(t *testing.T) {
	idx := ColorIndex("SHA256:abc123")
	if idx < 0 || idx > 11 {
		t.Errorf("color index %d out of range 0-11", idx)
	}
	if ColorIndex("SHA256:abc123") != idx {
		t.Error("nondeterministic color index")
	}
}

func TestIsOwnerFingerprint(t *testing.T) {
	if !IsOwnerFingerprint("abc123", "abc123") {
		t.Error("matching fingerprint should be owner")
	}
	if IsOwnerFingerprint("abc123", "xyz789") {
		t.Error("non-matching fingerprint should not be owner")
	}
	if IsOwnerFingerprint("", "abc123") {
		t.Error("empty fingerprint should not be owner")
	}
}

func TestOwnerDisplayName(t *testing.T) {
	if OwnerDisplayName("neur0map") != "★ neur0map" {
		t.Errorf("got %q", OwnerDisplayName("neur0map"))
	}
	if OwnerDisplayName("alice") != "★ alice" {
		t.Errorf("got %q", OwnerDisplayName("alice"))
	}
}

func TestHasFlair(t *testing.T) {
	if HasFlair(2) {
		t.Error("2 visits should not have flair")
	}
	if !HasFlair(3) {
		t.Error("3 visits should have flair")
	}
	if !HasFlair(10) {
		t.Error("10 visits should have flair")
	}
}

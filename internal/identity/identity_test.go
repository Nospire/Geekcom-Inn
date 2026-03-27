package identity

import "testing"

func TestDefaultNickname(t *testing.T) {
	fp := "SHA256:abc123def456"
	n1 := DefaultNickname(fp)
	n2 := DefaultNickname(fp)
	if n1 != n2 {
		t.Errorf("nondeterministic: %q != %q", n1, n2)
	}
	if len(n1) != 8 {
		t.Errorf("len = %d, want 8", len(n1))
	}
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

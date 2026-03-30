package mention

import "testing"

func TestExtractTokens_SingleMention(t *testing.T) {
	tokens := ExtractTokens("hey @deadmau5 check this")
	if len(tokens) != 1 || tokens[0] != "deadmau5" {
		t.Errorf("got %v", tokens)
	}
}

func TestExtractTokens_MultipleMentions(t *testing.T) {
	tokens := ExtractTokens("@alice and @bob are here")
	if len(tokens) != 2 {
		t.Errorf("expected 2, got %v", tokens)
	}
}

func TestExtractTokens_StartOfMessage(t *testing.T) {
	tokens := ExtractTokens("@neur0map hello")
	if len(tokens) != 1 || tokens[0] != "neur0map" {
		t.Errorf("got %v", tokens)
	}
}

func TestExtractTokens_NoMention(t *testing.T) {
	tokens := ExtractTokens("hello world")
	if len(tokens) != 0 {
		t.Errorf("expected none, got %v", tokens)
	}
}

func TestExtractTokens_EmailNotMention(t *testing.T) {
	tokens := ExtractTokens("email me@example.com")
	// "me@example.com" — the @ is mid-word, not a mention
	// Only @ after space or at start counts
	if len(tokens) != 0 {
		t.Errorf("email should not be a mention, got %v", tokens)
	}
}

func TestExtractTokens_AtSignAlone(t *testing.T) {
	tokens := ExtractTokens("hey @ what")
	if len(tokens) != 0 {
		t.Errorf("bare @ should not match, got %v", tokens)
	}
}

func TestIsMentioned_CaseInsensitive(t *testing.T) {
	if !IsMentioned("hey @DeadMau5 check this", "deadmau5") {
		t.Error("should match case-insensitively")
	}
}

func TestIsMentioned_NotMentioned(t *testing.T) {
	if IsMentioned("hey @alice", "bob") {
		t.Error("bob should not be mentioned")
	}
}

func TestIsMentioned_PartialNickNotMatch(t *testing.T) {
	if IsMentioned("hey @dead", "deadmau5") {
		t.Error("partial nick should not count as mention")
	}
}

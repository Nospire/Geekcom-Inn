package chat

import "testing"

func TestParseInput_ChatMessage(t *testing.T) {
	result := ParseInput("hello everyone")
	if result.IsCommand {
		t.Error("plain text should not be a command")
	}
	if result.Text != "hello everyone" {
		t.Errorf("text = %q", result.Text)
	}
}

func TestParseInput_Command(t *testing.T) {
	result := ParseInput("/nick dusty_dev")
	if !result.IsCommand {
		t.Error("should be a command")
	}
	if result.Command != "nick" {
		t.Errorf("command = %q, want nick", result.Command)
	}
	if result.Args != "dusty_dev" {
		t.Errorf("args = %q, want dusty_dev", result.Args)
	}
}

func TestParseInput_CommandNoArgs(t *testing.T) {
	result := ParseInput("/help")
	if !result.IsCommand {
		t.Error("should be a command")
	}
	if result.Command != "help" {
		t.Errorf("command = %q", result.Command)
	}
	if result.Args != "" {
		t.Errorf("args = %q, want empty", result.Args)
	}
}

func TestParseInput_SlashMidMessage(t *testing.T) {
	result := ParseInput("use /nick to change name")
	if result.IsCommand {
		t.Error("slash mid-message is not a command")
	}
}

func TestParseInput_Empty(t *testing.T) {
	result := ParseInput("")
	if result.IsCommand {
		t.Error("empty should not be command")
	}
}

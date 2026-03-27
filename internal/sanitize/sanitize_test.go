package sanitize

import "testing"

func TestClean(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"printable ascii passes through", "hello world", "hello world"},
		{"strips control chars", "hello\x00\x01\x02world", "helloworld"},
		{"strips escape sequences", "hello\x1b[31mworld", "hello[31mworld"},
		{"strips null bytes", "he\x00llo", "hello"},
		{"preserves space", "foo bar", "foo bar"},
		{"preserves punctuation", "!@#$%^&*()", "!@#$%^&*()"},
		{"strips high bytes", "hello\x80\xff", "hello"},
		{"strips tabs and newlines", "hello\t\nworld", "helloworld"},
		{"empty string", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Clean(tt.input)
			if got != tt.want {
				t.Errorf("Clean(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCleanNick(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"valid nick", "dusty_dev", "dusty_dev", false},
		{"strips bad chars", "dust\x00y", "dusty", false},
		{"too short", "a", "", true},
		{"too long", "abcdefghijklmnopqrstu", "", true},
		{"min length", "ab", "ab", false},
		{"max length", "abcdefghijklmnopqrst", "abcdefghijklmnopqrst", false},
		{"empty after clean", "\x00\x01", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CleanNick(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CleanNick(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CleanNick(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCleanChat(t *testing.T) {
	t.Run("normal message", func(t *testing.T) {
		got := CleanChat("hello everyone")
		if got != "hello everyone" {
			t.Errorf("got %q", got)
		}
	})
	t.Run("truncates at 500", func(t *testing.T) {
		input := make([]byte, 600)
		for i := range input {
			input[i] = 'a'
		}
		got := CleanChat(string(input))
		if len(got) != 500 {
			t.Errorf("len = %d, want 500", len(got))
		}
	})
	t.Run("strips escape codes", func(t *testing.T) {
		got := CleanChat("\x1b[0mhello")
		if got != "[0mhello" {
			t.Errorf("got %q", got)
		}
	})
}

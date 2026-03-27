package sanitize

import "errors"

// Clean strips all bytes outside printable ASCII range (0x20-0x7E).
func Clean(input string) string {
	buf := make([]byte, 0, len(input))
	for i := 0; i < len(input); i++ {
		b := input[i]
		if b >= 0x20 && b <= 0x7E {
			buf = append(buf, b)
		}
	}
	return string(buf)
}

// CleanNick sanitizes a nickname. Must be 2-20 printable ASCII chars after cleaning.
func CleanNick(nick string) (string, error) {
	cleaned := Clean(nick)
	if len(cleaned) < 2 || len(cleaned) > 20 {
		return "", errors.New("nickname must be 2-20 characters")
	}
	return cleaned, nil
}

// CleanChat sanitizes a chat message. Strips non-printable bytes, caps at 500 chars.
func CleanChat(msg string) string {
	cleaned := Clean(msg)
	if len(cleaned) > 500 {
		cleaned = cleaned[:500]
	}
	return cleaned
}

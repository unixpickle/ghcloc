package ghcloc

import "testing"

func TestParseAndString(t *testing.T) {
	res := NewRepository("unixpickle", "ghcloc").String()
	if res != "unixpickle/ghcloc" {
		t.Error("Failed to stringify 'unixpickle/ghcloc': " + res)
	}
	parsed, err := ParseRepository(res)
	if err != nil || parsed == nil {
		t.Error("Failed to parse 'unixpickle/ghcloc': " + err.Error())
	} else if parsed.User != "unixpickle" || parsed.Name != "ghcloc" {
		t.Error("Parse expected 'unixpickle/ghcloc': " + parsed.String())
	}
	parsed, err = ParseRepository("hey")
	if err == nil {
		t.Error("Expected error when parsing 'hey'")
	}
	parsed, err = ParseRepository("hey/yo/yo")
	if err == nil {
		t.Error("Expected error when parsing 'hey/yo/yo'")
	}
}

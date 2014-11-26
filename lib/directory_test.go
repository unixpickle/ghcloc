package ghcloc

import "testing"

func TestReadDir(t *testing.T) {
	repo := NewRepository("unixpickle", "ghcloc")
	listing, err := repo.ReadDir("/")
	if err != nil || listing == nil {
		t.Error("Failed to ReadDir of 'unixpickle/ghcloc':", err)
	}
	listing, err = repo.ReadDir("/FILE_IS_NOT_HERE")
	if err == nil {
		t.Error("Expected ReadDir of '/FILE_IS_NOT_HERE' to fail.")
	}
}

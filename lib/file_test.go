package ghcloc

import "testing"

func TestReadFile(t *testing.T) {
	repo := NewRepository("unixpickle", "ghcloc")
	file, err := repo.ReadFile("/README.md")
	if err != nil || file == nil {
		t.Error("Failed to ReadFile '/README.md'.")
	} else if data, err := file.Bytes(); data == nil || err != nil {
		t.Error("Failed to get data of '/README.md'")
	}
	file, err = repo.ReadFile("/FILE_IS_NOT_HERE")
	if file != nil || err == nil {
		t.Error("Expected ReadFile '/FILE_IS_NOT_HERE' to fail.")
	}
}

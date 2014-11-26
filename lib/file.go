package ghcloc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

type File struct {
	Entity
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func (self *File) Bytes() ([]byte, error) {
	if self.Encoding == "base64" {
		return base64.StdEncoding.DecodeString(self.Content)
	} else {
		return nil, errors.New("Unknown encoding: " + self.Encoding)
	}
}

func (self *Repository) ReadFile(path string) (*File, error) {
	url := "https://api.github.com/repos/" + self.String() + "/contents" + path
	contents, err := self.Request(url)
	if err != nil {
		return nil, err
	}

	var result File
	if err := json.Unmarshal(contents, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

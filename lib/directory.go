package ghcloc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	FILE_TYPE = "file"
	DIR_TYPE = "dir"
)

type Entity struct {
	Name    string            `json:"name"`
	Path    string            `json:"path"`
	SHA     string            `json:"sha"`
	Size    int               `json:"size"`
	URL     string            `json:"url"`
	HTMLURL string            `json:"html_url"`
	GitURL  string            `json:"git_url"`
	Type    string            `json:"type"`
	Links   map[string]string `json:"_links"`
}

func (self *Entity) IsFile() bool {
	return self.Type == FILE_TYPE
}

func (self *Entity) IsDir() bool {
	return self.Type == DIR_TYPE
}

func (self *Repository) ReadDir(path string) ([]Entity, error) {
	if path == "" {
		return self.ReadDir("/")
	}
	url := "https://api.github.com/repos/" + self.String() + "/contents" + path
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	
	// Check for raw error value
	var raw map[string]string
	if err := json.Unmarshal(contents, &raw); err == nil {
		if message, ok := raw["message"]; ok {
			return nil, errors.New(message)
		}
	}
	
	var result []Entity
	if err := json.Unmarshal(contents, &result); err != nil {
		return nil, err
	}
	return result, nil
}

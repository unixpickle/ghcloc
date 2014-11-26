package ghcloc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type FileInfo struct {
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

func (self *Repository) ReadDir(path string) ([]FileInfo, error) {
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
	var result []FileInfo
	if err := json.Unmarshal(contents, &result); err != nil {
		return nil, err
	}
	return result, nil
}

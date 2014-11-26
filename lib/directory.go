package ghcloc

import "encoding/json"

const (
	FILE_TYPE = "file"
	DIR_TYPE  = "dir"
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
	contents, err := self.Request(url)
	if err != nil {
		return nil, err
	}

	var result []Entity
	if err := json.Unmarshal(contents, &result); err != nil {
		return nil, err
	}
	return result, nil
}

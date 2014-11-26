package ghcloc

import (
	pathlib "path"
	"strings"
)

type Counts struct {
	TotalLines map[string]int
}

func NewCounts() *Counts {
	return &Counts{map[string]int{}}
}

func (self *Counts) Add(other *Counts) {
	for language, count := range other.TotalLines {
		self.TotalLines[language] += count
	}
}

func CountInEntity(repo *Repository, entity Entity) *Counts {
	if entity.IsFile() {
		return CountInFile(repo, entity.Path)
	} else if entity.IsDir() {
		return CountInDir(repo, entity.Path)
	} else {
		return NewCounts()
	}
}

func CountInFile(repo *Repository, path string) *Counts {
	extension := pathlib.Ext(path)
	supported := map[string]string{"go": "Go", "cpp": "C++",
		"hpp": "C++", "c": "C", "h": "C", "c++": "C++", "cc": "C++",
		"js": "JavaScript", "java": "Java", "coffee": "CoffeeScript",
		"coffeescript": "CoffeeScript", "m": "Objective-C",
		"mm": "Objective-C++", "md": "Markdown", "makefile": "Makefile",
		"mk": "Makefile", "rs": "Rust", "s": "Assembly", "asm": "Assembly",
		"php": "PHP", "html": "HTML", "css": "CSS", "py": "Python",
		"rb": "Ruby"}
	result := NewCounts()
	
	// Detect the language or return nothing.
	language, ok := supported[strings.ToLower(extension)]
	if !ok {
		return result
	}
	
	// Read the file or return nothing.
	file, err := repo.ReadFile(path)
	if err != nil {
		return result
	}
	
	// Get the contents and count the lines.
	if data, err := file.Bytes(); err == nil {
		textContents := string(data)
		lines := strings.Count(textContents, "\n") + 1
		result.TotalLines[language] = lines
	}
	
	return result
}

func CountInDir(repo *Repository, path string) *Counts {
	result := NewCounts()
	listing, err := repo.ReadDir(path)
	if err != nil {
		return result
	}
	remaining := 0
	response := make(chan *Counts)
	for _, entity := range listing {
		go func(entity Entity) {
			response <- CountInEntity(repo, entity)
		}(entity)
		remaining++
	}
	for remaining > 0 {
		result.Add(<-response)
		remaining--
	}
	return result
}
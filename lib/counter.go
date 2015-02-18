package ghcloc

import (
	pathlib "path"
	"strings"
)

type Counts struct {
	FileCount  int
	TotalLines map[string]int
}

func NewCounts() *Counts {
	return &Counts{0, map[string]int{}}
}

func (self *Counts) Add(other *Counts) {
	self.FileCount += other.FileCount
	for language, count := range other.TotalLines {
		self.TotalLines[language] += count
	}
}

func CountInEntity(repo *Repository, entity Entity) (*Counts, error) {
	if entity.IsFile() {
		return CountInFile(repo, entity.Path)
	} else if entity.IsDir() {
		return CountInDir(repo, entity.Path)
	} else {
		return NewCounts(), nil
	}
}

func CountInFile(repo *Repository, path string) (*Counts, error) {
	extension := pathlib.Ext(path)
	if len(extension) > 0 {
		extension = extension[1:]
	}
	supported := map[string]string{"go": "Go", "cpp": "C++",
		"hpp": "C++", "c": "C", "h": "C", "c++": "C++", "cc": "C++",
		"js": "JavaScript", "json": "JSON", "java": "Java",
		"coffee": "CoffeeScript", "coffeescript": "CoffeeScript",
		"m": "Objective-C", "mm": "Objective-C++", "md": "Markdown",
		"makefile": "Makefile", "rb": "Ruby",
		"mk": "Makefile", "rs": "Rust", "s": "Assembly", "asm": "Assembly",
		"php": "PHP", "html": "HTML", "css": "CSS", "py": "Python",
		"pde": "Processing", "pragmash": "Pragmash", "sh": "Shell"}
	result := NewCounts()

	// Detect the language or return nothing.
	language, ok := supported[strings.ToLower(extension)]
	if !ok {
		return result, nil
	}

	// Read the file or return nothing.
	file, err := repo.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Get the contents and count the lines.
	if data, err := file.Bytes(); err == nil {
		textContents := string(data)
		lines := strings.Count(textContents, "\n") + 1
		result.TotalLines[language] = lines
		result.FileCount = 1
	} else {
		return nil, err
	}

	return result, nil
}

func CountInDir(repo *Repository, path string) (*Counts, error) {
	listing, err := repo.ReadDir(path)
	if err != nil {
		return nil, err
	}
	remaining := 0
	response := make(chan bgTaskResult)
	for _, entity := range listing {
		go func(entity Entity) {
			counts, err := CountInEntity(repo, entity)
			response <- bgTaskResult{counts, err}
		}(entity)
		remaining++
	}
	result := NewCounts()
	for remaining > 0 {
		remaining--
		res := <-response
		if res.err != nil {
			return nil, res.err
		} else {
			result.Add(res.counts)
		}
	}
	return result, nil
}

type bgTaskResult struct {
	counts *Counts
	err    error
}

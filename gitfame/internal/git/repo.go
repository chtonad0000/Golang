package git

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetFiles(repoPath, extensions, languages, excludePatterns, restrictPatterns string) ([]string, error) {
	if _, err := exec.Command("git", "-C", repoPath, "rev-parse").Output(); err != nil {
		return nil, errors.New("not a git repository")
	}

	output, err := exec.Command("git", "-C", repoPath, "ls-tree", "-r", "--name-only", "HEAD").Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")

	if extensions != "" {
		extFilters := strings.Split(extensions, ",")
		files = filterByExtensions(files, extFilters)
	}

	if languages != "" {
		langFilters := strings.Split(languages, ",")
		files, err = filterByLanguages(files, langFilters, "C:\\Users\\d_chu\\go-yandex\\dan305305\\gitfame\\configs\\language_extensions.json")
		if err != nil {
			return nil, err
		}
	}

	if excludePatterns != "" {
		excludeFilters := strings.Split(excludePatterns, ",")
		files = excludeByPatterns(files, excludeFilters)
	}

	if restrictPatterns != "" {
		restrictFilters := strings.Split(restrictPatterns, ",")
		files = restrictByPatterns(files, restrictFilters)
	}

	return files, nil
}

func filterByExtensions(files, extensions []string) []string {
	var filtered []string
	for _, file := range files {
		for _, ext := range extensions {
			if strings.HasSuffix(file, ext) {
				filtered = append(filtered, file)
				break
			}
		}
	}
	return filtered
}
func filterByLanguages(files, languages []string, languagesFile string) ([]string, error) {
	extensionToLanguage, err := LoadLanguages(languagesFile)
	if err != nil {
		return nil, err
	}

	var filtered []string
	for _, file := range files {
		ext := filepath.Ext(file)
		if lang, found := extensionToLanguage[ext]; found {
			for _, requestedLang := range languages {
				if lang == requestedLang {
					filtered = append(filtered, file)
					break
				}
			}
		}
	}
	return filtered, nil
}
func excludeByPatterns(files, patterns []string) []string {
	var filtered []string
	for _, file := range files {
		exclude := false
		for _, pattern := range patterns {
			match, _ := filepath.Match(pattern, file)
			if match {
				exclude = true
				break
			}
		}
		if !exclude {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func restrictByPatterns(files, patterns []string) []string {
	var filtered []string
	for _, file := range files {
		for _, pattern := range patterns {
			match, _ := filepath.Match(pattern, file)
			if match {
				filtered = append(filtered, file)
				break
			}
		}
	}
	return filtered
}

type Language struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Extensions []string `json:"extensions"`
}

func LoadLanguages(filepath string) (map[string]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	var languages []Language
	if err := json.NewDecoder(file).Decode(&languages); err != nil {
		return nil, err
	}

	extensionToLanguage := make(map[string]string)
	for _, lang := range languages {
		for _, ext := range lang.Extensions {
			extensionToLanguage[ext] = lang.Name
		}
	}

	return extensionToLanguage, nil
}

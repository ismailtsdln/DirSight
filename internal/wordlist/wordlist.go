package wordlist

import (
	"bufio"
	"os"
	"strings"

	"github.com/ismailtsdln/DirSight/internal/bypass"
)

// Loader handles loading and processing of wordlists
type Loader struct{}

// LoadFromFile reads a wordlist file and returns a slice of paths
func (l *Loader) LoadFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var wordlist []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" && !strings.HasPrefix(word, "#") {
			wordlist = append(wordlist, word)
		}
	}

	return wordlist, scanner.Err()
}

// ExpandWithBypasses takes a base wordlist and adds path-based bypass variations
func (l *Loader) ExpandWithBypasses(baseList []string) []string {
	var expanded []string
	seen := make(map[string]bool)

	for _, word := range baseList {
		if !seen[word] {
			expanded = append(expanded, word)
			seen[word] = true
		}

		variations := bypass.GeneratePathBypasses(word)
		for _, v := range variations {
			if !seen[v] {
				expanded = append(expanded, v)
				seen[v] = true
			}
		}
	}

	return expanded
}

package pipeline

import (
	"bufio"
	"strings"
)

// ReflowTxt takes raw whisper-cli txt output (one line per segment) and
// reflows it into clean paragraphs. Consecutive segment lines are joined
// with a space. Blank lines (gaps between whisper segments) are preserved
// as paragraph breaks.
func ReflowTxt(input string) string {
	if strings.TrimSpace(input) == "" {
		return ""
	}

	scanner := bufio.NewScanner(strings.NewReader(input))
	var paragraphs []string
	var current []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if len(current) > 0 {
				paragraphs = append(paragraphs, strings.Join(current, " "))
				current = nil
			}
			continue
		}
		current = append(current, line)
	}

	if len(current) > 0 {
		paragraphs = append(paragraphs, strings.Join(current, " "))
	}

	return strings.Join(paragraphs, "\n\n")
}

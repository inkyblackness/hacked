package text

import "strings"

const (
	// textLineLimit is the lowest number found by a cursory search.
	// It appears that various texts have different limits.
	textLineLimit = 79
)

// Blocked splits the given input string into a series of lines, usable for blocked serialization.
// It returns an array of strings that each are below a maximum line length. The list is terminated with
// one empty line (necessary for compound strings).
func Blocked(input string) []string {
	var result []string
	textLines := strings.Split(input, "\n")
	resultLine := ""
	addBlock := func() {
		if len(resultLine) > 0 {
			result = append(result, resultLine)
			resultLine = ""
		}
	}

	for textLineIndex, textLine := range textLines {
		words := strings.Split(textLine, " ")
		for wordIndex, word := range words {
			if wordIndex > 0 {
				resultLine += " "
			}
			if (len(resultLine) + len(word)) > textLineLimit {
				addBlock()
			}
			resultLine += word
		}
		if textLineIndex < (len(textLines) - 1) {
			resultLine += "\n"
		}
	}
	addBlock()
	return append(result, "")
}

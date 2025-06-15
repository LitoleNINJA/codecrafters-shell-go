package main

import (
	"os"
	"strings"
)

func parseInput(input string) (string, []string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil
	}

	var args []string
	var current strings.Builder
	inSingleQuotes, inDoubleQuotes := false, false

	for i := 0; i < len(input); i++ {
		char := input[i]

		switch {
		case char == '\'' && !inDoubleQuotes:
			inSingleQuotes = !inSingleQuotes
		case char == '"' && !inSingleQuotes:
			inDoubleQuotes = !inDoubleQuotes
		case char == ' ' && !inSingleQuotes && !inDoubleQuotes:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
			// Skip multiple spaces
			for i+1 < len(input) && input[i+1] == ' ' {
				i++
			}
		case char == '\\' && !inSingleQuotes:
			if nextPos, ok := handleEscape(input, i, &current, inDoubleQuotes); ok {
				i = nextPos
				continue
			} else {
				return "", nil
			}
		default:
			current.WriteByte(char)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	if len(args) == 0 {
		return "", nil
	}

	return args[0], args[1:]
}

func handleEscape(input string, i int, current *strings.Builder, inDoubleQuotes bool) (int, bool) {
	if i+1 >= len(input) {
		return i, false
	}

	nextChar := input[i+1]

	if inDoubleQuotes {
		// In double quotes, only escape ", \, and $
		// Backslashes preceding characters without a special meaning are left unmodified
		if !(nextChar == '"' || nextChar == '\\' || nextChar == '$') {
			current.WriteByte('\\')
		}
	}

	current.WriteByte(nextChar)
	return i + 1, true
}

func isOnPath(cmd string) (string, bool) {
	osPath := os.Getenv("PATH")
	paths := strings.Split(osPath, ":")

	for _, path := range paths {
		if _, err := os.Stat(path + "/" + cmd); err == nil {
			return path, true
		}
	}

	return "", false
}

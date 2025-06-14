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
    inQuotes := false
    quoteChar := byte(0)

    for i := 0; i < len(input); i++ {
        char := input[i]

        if !inQuotes && (char == '\'' || char == '"') {
            inQuotes = true
            quoteChar = char
        } else if inQuotes && char == quoteChar {
            inQuotes = false
            quoteChar = 0
        } else if !inQuotes && char == ' ' {
            if current.Len() > 0 {
                args = append(args, current.String())
                current.Reset()
            }

            for i+1 < len(input) && input[i+1] == ' ' {
                i++
            }
        } else {
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
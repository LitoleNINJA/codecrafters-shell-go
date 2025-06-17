package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	NoMatch int = iota
	FullMatch
	PartialMatch
	MultipleMatch
)

var tabCount = 0

func autoComplete(input *strings.Builder) {
	tabCount++
	currentInput := input.String()
	completed, matchType := tryAutoComplete(currentInput, tabCount)

	if completed != currentInput && matchType == FullMatch {
		// Clear current line and rewrite with completed text
		fmt.Print("\r\033[K$ ")
		fmt.Printf("%s ", completed)
		os.Stdout.Sync() // Force flush
		input.Reset()
		input.WriteString(completed + " ")

		tabCount = 0 // Reset tab count after full match
	} else if matchType == NoMatch {
		// if no completion, print bell sound
		fmt.Print("\x07") // ASCII Bell
	} else if matchType == MultipleMatch && tabCount > 1 {
		// if multiple completions
		fmt.Printf("$ %s", currentInput)
		tabCount = 0 // Reset tab count after multiple matches
	} else if matchType == PartialMatch {
		// Clear current line and rewrite with completed text
		fmt.Print("\r\033[K$ ")
		fmt.Printf("%s", completed)
		os.Stdout.Sync() // Force flush
		input.Reset()
		input.WriteString(completed)

		tabCount = 0 // Reset tab count after partial match
	}
}

func tryAutoComplete(input string, tabCount int) (string, int) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", NoMatch
	}

	var matches []string

	// check if it is a built-in command
	for _, cmd := range builtInCommands {
		if strings.HasPrefix(cmd, input) {
			matches = append(matches, cmd)
		}
	}

	// check if it an executable in PATH directory
	if len(matches) == 0 {
		executables, err := checkPATH()
		if err != nil {
			fmt.Printf("Error checking PATH: %v\n", err)
			return input, NoMatch
		}

		for _, exec := range executables {
			if strings.HasPrefix(exec, input) {
				matches = append(matches, exec)
			}
		}
	}

	if len(matches) == 1 {
		return matches[0], FullMatch
	} else if len(matches) > 1 {
		// Multiple matches - Ring bell for 1st tab, print all for 2nd tab
		if tabCount == 1 {
			fmt.Print("\a") // ASCII Bell

			longestPrefix := longestCommonPrefix(matches)

			if longestPrefix != input {
				return longestPrefix, PartialMatch
			} else {
				return input, MultipleMatch
			}
		} else if tabCount > 1 {
			slices.Sort(matches)
			fmt.Printf("\n\r")
			fmt.Printf("%s", strings.Join(matches, "  "))
			fmt.Printf("\n\r")

			return input, MultipleMatch
		}
	}

	return input, NoMatch
}

func checkPATH() ([]string, error) {
	var executables []string

	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return nil, fmt.Errorf("PATH environment variable is not set")
	}

	pathDirs := strings.Split(pathEnv, pathSeparator)
	for _, dir := range pathDirs {
		// read directory
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, file := range files {
			if !file.IsDir() {
				fileName := file.Name()

				if isExecutable(filepath.Join(dir, fileName)) {
					executables = append(executables, fileName)
				}
			}
		}
	}

	return executables, nil
}

func isExecutable(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	// Check if file has execute permission
	mode := info.Mode()
	return mode&0111 != 0 // Check if any execute bit is set
}

func longestCommonPrefix(matches []string) string {
	longestPrefix := matches[len(matches)-1]

	for _, match := range matches {
		for i := 0; i < len(longestPrefix) && i < len(match); i++ {
			if longestPrefix[i] != match[i] {
				longestPrefix = longestPrefix[:i]
				break
			}
		}
	}

	return longestPrefix
}

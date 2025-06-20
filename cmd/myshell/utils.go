package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

const (
	singleQuote   = '\''
	doubleQuote   = '"'
	backslash     = '\\'
	whitespace    = ' '
	redirOut      = '>'
	redirIn       = '<'
	pathSeparator = string(os.PathListSeparator)
)

type RedirectionType int

const (
	NoRedirection RedirectionType = iota
	OutputRedirection
	InputRedirection
	ErrorRedirection
	AppendOutRedirection
	AppendErrRedirection
)

type ParsedCommand struct {
	Cmd       string
	Args      []string
	RedirType RedirectionType
	RedirFile string
}

func parseInput(input string) ParsedCommand {
	input = strings.TrimSpace(input)
	if input == "" {
		return ParsedCommand{}
	}

	redirType := NoRedirection
	redirFile := ""
	for i := range len(input) {
		if input[i] == redirOut || input[i] == redirIn {
			leftInput, err := processRedirection(input, i, &redirType, &redirFile)
			if err != nil {
				fmt.Printf("Error processing redirection: %v\n", err)
				return ParsedCommand{}
			}
			input = leftInput
			break
		}
	}

	parsedArgs := []string{}
	currentArg := strings.Builder{}
	isInSingleQuotes := false
	isInDoubleQuotes := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		switch {
		case char == singleQuote && !isInDoubleQuotes:
			isInSingleQuotes = !isInSingleQuotes

		case char == doubleQuote && !isInSingleQuotes:
			isInDoubleQuotes = !isInDoubleQuotes

		case char == whitespace && !isInSingleQuotes && !isInDoubleQuotes:
			if currentArg.Len() > 0 {
				parsedArgs = append(parsedArgs, currentArg.String())
				currentArg.Reset()
			}
			i = skipConsecutiveSpaces(input, i)

		case char == backslash && !isInSingleQuotes:
			nextIndex, success := processEscapeSequence(input, i, &currentArg, isInDoubleQuotes)
			if !success {
				return ParsedCommand{}
			}
			i = nextIndex

		default:
			currentArg.WriteByte(char)
		}
	}

	// Add the last argument if it exists
	if currentArg.Len() > 0 {
		parsedArgs = append(parsedArgs, currentArg.String())
	}

	if len(parsedArgs) == 0 {
		return ParsedCommand{}
	}

	return ParsedCommand{
		Cmd:       parsedArgs[0],
		Args:      parsedArgs[1:],
		RedirType: redirType,
		RedirFile: redirFile,
	}
}

func skipConsecutiveSpaces(input string, currentIndex int) int {
	for currentIndex+1 < len(input) && input[currentIndex+1] == whitespace {
		currentIndex++
	}
	return currentIndex
}

func processEscapeSequence(input string, currentIndex int, builder *strings.Builder, isInDoubleQuotes bool) (newIndex int, success bool) {
	if currentIndex+1 >= len(input) {
		return currentIndex, false
	}

	escapedChar := input[currentIndex+1]

	if isInDoubleQuotes {
		// In double quotes, only ", \, and $ need escaping
		// Other backslashes are preserved literally
		if !(escapedChar == doubleQuote || escapedChar == backslash || escapedChar == '$') {
			builder.WriteByte(backslash)
		}
	}

	builder.WriteByte(escapedChar)
	return currentIndex + 1, true
}

func processRedirection(input string, pos int, redirType *RedirectionType, redirFile *string) (string, error) {
	if pos >= len(input) {
		return "", fmt.Errorf("redirection operator at end of input")
	}

	if input[pos] == redirOut {
		// Check for append redirection '>>' or '1>>' or '2>>'
		if pos+1 < len(input) && input[pos+1] == redirOut {
			if pos > 0 && input[pos-1] == '2' {
				*redirType = AppendErrRedirection
			} else {
				*redirType = AppendOutRedirection
			}

			*redirFile = strings.TrimSpace(input[pos+2:])
			return strings.TrimSpace(input[:pos-1]), nil
		}

		*redirFile = strings.TrimSpace(input[pos+1:])

		// Handle stdout redirection '1>'
		if pos > 0 && input[pos-1] == '1' {
			*redirType = OutputRedirection
			return strings.TrimSpace(input[:pos-1]), nil
		} else if pos > 0 && input[pos-1] == '2' {
			// Handle stderr redirection '2>'
			*redirType = ErrorRedirection
			return strings.TrimSpace(input[:pos-1]), nil
		}

		// Handle regular stdout redirection '>'
		*redirType = OutputRedirection
		return strings.TrimSpace(input[:pos]), nil

	} else if input[pos] == redirIn {
		*redirType = InputRedirection
		*redirFile = strings.TrimSpace(input[pos+1:])
		return strings.TrimSpace(input[:pos]), nil
	}

	return input, nil
}

func isOnPath(command string) (foundPath string, exists bool) {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return "", false
	}

	pathDirectories := strings.Split(pathEnv, pathSeparator)

	for _, directory := range pathDirectories {
		if directory == "" {
			continue
		}

		fullCommandPath := filepath.Join(directory, command)
		if commandExists(fullCommandPath) {
			return directory, true
		}
	}

	return "", false
}

func commandExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readUserInput() string {
	var input strings.Builder

	// Get the file descriptor for stdin
	fd := int(os.Stdin.Fd())

	// Save the original terminal state
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Printf("Error setting raw mode: %v\n", err)
		os.Exit(1)
	}

	// Restore terminal state
	defer term.Restore(fd, oldState)

	for {
		// Read one byte at a time
		var buf [1]byte
		n, err := os.Stdin.Read(buf[:])
		if err != nil || n == 0 {
			continue
		}

		char := rune(buf[0])

		switch char {
		case '\n', '\r': //  Enter key
			fmt.Print("\n\r")
			return input.String()

		case '\t': // Tab : autocomplete
			autoComplete(&input)

		case 127, 8: // Backspace (127 is DEL, 8 is BS)
			if input.Len() > 0 {
				// Remove last character from input
				currentStr := input.String()
				input.Reset()
				input.WriteString(currentStr[:len(currentStr)-1])
				// Move cursor back, print space, move back again
				fmt.Print("\b \b")
			}

		case 3: // Ctrl+C
			fmt.Print("\n\r")
			term.Restore(fd, oldState) // Restore before exit
			os.Exit(0)

		case 4: // Ctrl+D (EOF)
			fmt.Print("\n\r")
			term.Restore(fd, oldState) // Restore before exit
			os.Exit(0)

		case 27: // Escape sequence
			// Read the next two bytes for arrow keys or other escape sequences
			var buf2 [2]byte
			n, err := os.Stdin.Read(buf2[:])
			if err != nil || n < 2 {
				continue
			}

			// Check for arrow keys
			if buf2[0] == '[' {
				switch buf2[1] {
				case 'A': // Up arrow - get the previos cmd
					previousCmd := getPreviousCommand()

					if previousCmd.Cmd == "" {
						continue
					}

					// clear the input and terminal
					fmt.Print("\r\033[K")
					input.Reset()

					fmt.Printf("$ %s", previousCmd.Cmd)
					input.WriteString(previousCmd.Cmd)

					for _, arg := range previousCmd.Args {
						fmt.Printf(" %s", arg)
						input.WriteString(" " + arg)
					}
				}
			}

		default:
			if char >= 32 && char < 127 {
				input.WriteRune(char)
				fmt.Print(string(char))
			}
		}
	}
}

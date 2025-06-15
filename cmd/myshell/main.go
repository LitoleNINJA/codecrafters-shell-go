package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var builtInCommands = []string{"exit", "echo", "type", "pwd", "cd"}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		userInput, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("error reading input:", err)
			os.Exit(1)
		}

		parsedCommand := parseInput(userInput[:len(userInput)-1])

		handleCommand(&parsedCommand)
	}
}

func handleCommand(parsedCmd *ParsedCommand) {
	var outputFile 		*os.File
	var originalStd 	*os.File
	var err				error

	if parsedCmd.RedirType != NoRedirection {
		outputFile, originalStd, err = handleRedirection(parsedCmd.RedirType, parsedCmd.RedirFile)
		if err != nil {
			fmt.Println("Error handling redirection:", err)
			return
		}

		// Ensure we close the file and restore original std streams
		defer func() {
			if outputFile != nil {
				outputFile.Close()
			}
			if originalStd != nil {
				switch parsedCmd.RedirType {
				case OutputRedirection, AppendOutRedirection:
					os.Stdout = originalStd
				case InputRedirection:
					os.Stdin = originalStd
				case ErrorRedirection, AppendErrRedirection:
					os.Stderr = originalStd
				}
			}
		}()
	}

	switch parsedCmd.Cmd {
	case "exit":
		os.Exit(0)
	case "echo":
		if len(parsedCmd.Args) > 0 {
			fmt.Println(strings.Join(parsedCmd.Args, " "))
		}
	case "type":
		handleTypeCmd(parsedCmd.Args)
	case "pwd":
		handlePwdCmd()
	case "cd":
		handleCdCmd(parsedCmd.Args)
	default:
		runCommand(parsedCmd.Cmd, parsedCmd.Args)
	}
}

func handleRedirection(redirType RedirectionType, redirFile string) (*os.File, *os.File, error) {
	var outputFile *os.File
	var originalStd *os.File
	var err error

	if redirType == OutputRedirection {
		outputFile, err = os.OpenFile(redirFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return nil, nil, fmt.Errorf("could not create output file: %w", err)
		}

		originalStd = os.Stdout
		os.Stdout = outputFile
	} else if redirType == InputRedirection {
		inputFile, err := os.Open(redirFile)
		if err != nil {
			return nil, nil, fmt.Errorf("could not open input file: %w", err)
		}

		originalStd = os.Stdin
		os.Stdin = inputFile
	} else if redirType == ErrorRedirection {
		outputFile, err = os.OpenFile(redirFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return nil, nil, fmt.Errorf("could not create error file: %w", err)
		}

		originalStd = os.Stderr
		os.Stderr = outputFile
	} else if redirType == AppendOutRedirection {
		outputFile, err = os.OpenFile(redirFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, nil, fmt.Errorf("could not open file for appending: %w", err)
		}

		originalStd = os.Stdout
		os.Stdout = outputFile
	} else if redirType == AppendErrRedirection {
		outputFile, err = os.OpenFile(redirFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, nil, fmt.Errorf("could not open file for appending: %w", err)
		}
		
		originalStd = os.Stderr
		os.Stderr = outputFile
	} else {
		return nil, nil, fmt.Errorf("unknown redirection type")
	}

	return outputFile, originalStd, nil
}

func runCommand(cmd string, args []string) {
	command := exec.Command(cmd, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin  = os.Stdin

	err := command.Run()
    if err != nil {
        // Check if it's an ExitError (command found but exited with non-zero status)
        if _, ok := err.(*exec.ExitError); ok {
            // Command was found and ran, but exited with error
            // Don't print "command not found" - the command already printed its error
            return
        }
        // This is likely a "command not found" or similar startup error
        fmt.Fprintf(os.Stderr, "%s: command not found\n", cmd)
    }
}

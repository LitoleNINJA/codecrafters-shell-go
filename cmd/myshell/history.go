package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var HISTORY []ParsedCommand

const MAX_HISTORY = 100

var lastCommandPos int = -1

func addCmdToHistory(cmd ParsedCommand) {
	if len(HISTORY) >= MAX_HISTORY {
		HISTORY = HISTORY[1:] // Remove the oldest command
	}
	HISTORY = append(HISTORY, cmd)

	lastCommandPos = len(HISTORY)
}

func displayCmdHistory(args []string) {
	limit := len(HISTORY)
	if len(args) > 0 {
		parsedLimit, _ := strconv.ParseInt(args[0], 10, 64)
		limit = int(parsedLimit)
	}

	for i := len(HISTORY) - limit; i < len(HISTORY); i++ {
		cmd := HISTORY[i]
		fmt.Printf("\t%d  %s ", i+1, cmd.Cmd)
		for _, arg := range cmd.Args {
			fmt.Printf("%s ", arg)
		}
		fmt.Println()
	}
}

func getPreviousCommand() ParsedCommand {
	if lastCommandPos <= 0 || lastCommandPos > len(HISTORY) {
		return ParsedCommand{}
	}

	lastCommandPos--
	return HISTORY[lastCommandPos]
}

func getNextCommand() ParsedCommand {
	if lastCommandPos >= len(HISTORY)-1 {
		return ParsedCommand{}
	}

	lastCommandPos++
	return HISTORY[lastCommandPos]
}

func addContentsToHistory(fileName string) {
	historyFile, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening history file: %s\n", err)
		return
	}
	defer historyFile.Close()

	scanner := bufio.NewScanner(historyFile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Parse the line into command and arguments
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		cmd := ParsedCommand{
			Cmd:  parts[0],
			Args: parts[1:],
		}

		addCmdToHistory(cmd)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading history file: %v\n", err)
	}
}

func writeHistoryToFile(fileName string, append bool) {
	if fileName == "" {
		fmt.Println("No history file set")
		return
	}

	var historyFile *os.File
	if append {
		historyFile, _ = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	} else {
		historyFile, _ = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	}

	defer historyFile.Close()

	for _, cmd := range HISTORY {
		line := cmd.Cmd
		if len(cmd.Args) > 0 {
			line += " " + strings.Join(cmd.Args, " ")
		}

		_, err := historyFile.WriteString(line + "\n")
		if err != nil {
			fmt.Printf("Error writing to history file: %s\n", err)
			historyFile.Close()
			return
		}
	}

	// clear HISTORY
	HISTORY = []ParsedCommand{}
}

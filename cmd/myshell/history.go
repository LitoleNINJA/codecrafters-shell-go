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

func setHistoryFile(fileName string) {
	historyFile, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening history file: %s\n", err)
		return
	}

	addContentsToHistory(historyFile)
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

func addContentsToHistory(historyFile *os.File) {
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

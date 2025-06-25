package main

import (
	"fmt"
	"strconv"
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

	for i:=len(HISTORY) - limit; i < len(HISTORY); i++ {
		cmd := HISTORY[i]
		fmt.Printf("%d  %s ", i+1, cmd.Cmd)
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

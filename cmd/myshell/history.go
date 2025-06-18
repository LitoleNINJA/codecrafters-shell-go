package main

import "fmt"

var HISTORY []ParsedCommand
const MAX_HISTORY = 100

func addCmdToHistory(cmd ParsedCommand) {
	if len(HISTORY) >= MAX_HISTORY {
		HISTORY = HISTORY[1:] // Remove the oldest command
	}
	HISTORY = append(HISTORY, cmd)
}

func displayCmdHistory() {
	for i, cmd := range HISTORY {
		fmt.Printf("%d  %s ", i+1, cmd.Cmd)
		for _, arg := range cmd.Args {
			fmt.Printf("%s ", arg)
		}
		fmt.Println()
	}
}
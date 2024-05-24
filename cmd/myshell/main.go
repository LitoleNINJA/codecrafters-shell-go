package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var builtInCommands = []string{"exit", "echo", "type"}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "$ ")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error reading input:", err)
			continue
		}

		userInput = strings.ReplaceAll(userInput, "\n", "")
		cmd, args := parseInput(userInput)

		handleCommand(cmd, args)
	}
}

func handleCommand(cmd string, args []string) {
	switch cmd {
	case "exit":
		os.Exit(0)
	case "echo":
		fmt.Println(strings.Join(args, " "))
	case "type":
		if len(args) == 0 {
			fmt.Println("type: missing argument")
			return
		}

		if listContains(builtInCommands, args[0]) {
			fmt.Printf("%s is a shell builtin\n", args[0])
		} else {
			fmt.Printf("%s not found\n", args[0])
		}
	default:
		fmt.Printf("%s: command not found\n", cmd)
	}
}

func parseInput(input string) (string, []string) {
	fields := strings.Fields(input)
	cmd := fields[0]

	if len(fields) > 1 {
		return cmd, fields[1:]
	} else {
		return cmd, nil
	}
}

func listContains(list []string, elem string) bool {
	for _, e := range list {
		if e == elem {
			return true
		}
	}

	return false
}
